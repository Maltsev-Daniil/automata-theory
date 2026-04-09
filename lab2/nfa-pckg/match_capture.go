package nfapckg

import "sort"

type config struct {
	state *State
	pos   int
	cap   *capture_state
}

type capture_state struct {
	start map[string]int
	end   map[string]int
}

type state_pos_key struct {
	state_id int
	pos      int
}

func cloneCaptureState(cap *capture_state) *capture_state {
	result := &capture_state{
		start: make(map[string]int),
		end:   make(map[string]int),
	}

	if cap == nil {
		return result
	}

	for k, v := range cap.start {
		result.start[k] = v
	}
	for k, v := range cap.end {
		result.end[k] = v
	}

	return result
}

func applyCapEvent(cap *capture_state, event *cap_event, pos int) {
	if event == nil {
		return
	}

	if !event.finish {
		cap.start[event.name] = pos
	} else {
		cap.end[event.name] = pos
	}
}

func betterCapture(new_cap, old_cap *capture_state) bool {
	if old_cap == nil {
		return true
	}
	if new_cap == nil {
		return false
	}

	// собираем все имена групп
	// чтобы сортнуть для детерминизма
	names_map := make(map[string]struct{})

	for k := range new_cap.start {
		names_map[k] = struct{}{}
	}
	for k := range old_cap.start {
		names_map[k] = struct{}{}
	}

	// превращаем в слайс
	names := make([]string, 0, len(names_map))
	for k := range names_map {
		names = append(names, k)
	}

	sort.Strings(names)

	for _, name := range names {
		new_start, new_has_start := new_cap.start[name]
		old_start, old_has_start := old_cap.start[name]

		new_end, new_has_end := new_cap.end[name]
		old_end, old_has_end := old_cap.end[name]

		if new_has_start && !old_has_start {
			return true
		}
		if !new_has_start && old_has_start {
			return false
		}

		if !new_has_start && !old_has_start {
			continue
		}

		if new_start < old_start {
			return true
		}
		if new_start > old_start {
			return false
		}

		if new_has_end && old_has_end {
			new_len := new_end - new_start
			old_len := old_end - old_start

			if new_len > old_len {
				return true
			}
			if new_len < old_len {
				return false
			}
		}

		if new_has_end && !old_has_end {
			return true
		}
		if !new_has_end && old_has_end {
			return false
		}
	}

	return false
}

func epsilonClosure(input []config) []config {
	// стек для обхода
	stack := make([]config, len(input))
	copy(stack, input)

	visited := make(map[state_pos_key]config)

	for len(stack) > 0 {
		// pop
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		key := state_pos_key{
			state_id: cur.state.Id,
			pos:      cur.pos,
		}

		// если уже были в этом (state, pos)
		if old, ok := visited[key]; ok {
			if !betterCapture(cur.cap, old.cap) {
				continue
			}
		}

		// сохраняем лучший
		visited[key] = cur

		// идём по ε-переходам
		for _, tr := range cur.state.Epsilon {
			var new_cap *capture_state

			if tr.cap != nil {
				new_cap = cloneCaptureState(cur.cap)
				applyCapEvent(new_cap, tr.cap, cur.pos)
			} else {
				new_cap = cur.cap
			}

			next := config{
				state: tr.To,
				pos:   cur.pos, // ВАЖНО pos не меняется
				cap:   new_cap,
			}

			stack = append(stack, next)
		}
	}

	// собираем результат
	result := make([]config, 0, len(visited))
	for _, cfg := range visited {
		result = append(result, cfg)
	}

	return result
}

func moveLiteral(states []config, symbol string) []config {
	visited := make(map[state_pos_key]config)
	for _, cfg := range states {
		next_state := cfg.state.Ntran[symbol]
		for _, st := range next_state {
			new_cfg := config{
				state: st,
				pos:   cfg.pos + 1,
				cap:   cfg.cap, // no touching do it in epsilonClosure
			}

			key := state_pos_key{
				state_id: st.Id,
				pos:      new_cfg.pos,
			}

			if old, ok := visited[key]; ok {
				if !betterCapture(new_cfg.cap, old.cap) {
					continue
				}
			}
			visited[key] = new_cfg
		}
	}

	result := make([]config, 0, len(visited))
	for _, cfg := range visited {
		result = append(result, cfg)
	}
	return result
}

func matchNFAImpl(nfa *NFA, input string) []config {
	start_config := config{
		nfa.Head.Start,
		0,
		nil,
	}
	configs := epsilonClosure([]config{start_config})

	for _, ch := range input {
		symbol := string(ch)

		configs = moveLiteral(configs, symbol)
		configs = epsilonClosure(configs)
	}
	return configs
}

func (nfa *NFA) matchNFA(input string) bool {
	configs := matchNFAImpl(nfa, input)

	for _, cfg := range configs {
		if cfg.state == nfa.Head.Finish {
			return true
		}
	}
	return false
}

func (nfa *NFA) matchNFAWithCapture(input string) *capture_state {
	configs := matchNFAImpl(nfa, input)

	var best *capture_state
	for _, cfg := range configs {
		if cfg.state == nfa.Head.Finish {
			if best == nil || betterCapture(best, cfg.cap) {
				best = cfg.cap
			}
		}
	}
	return best
}
