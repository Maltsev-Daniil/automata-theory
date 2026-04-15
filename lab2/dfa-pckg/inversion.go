package dfapckg

import (
	"reglib/nfa-pckg"
	"sort"
	"strconv"
)

func DfaToNfaNReverse(dfa *DFA) *nfapckg.NFA {
	dfa_to_nfa := make(map[int]*nfapckg.State)
	for state := range dfa.dstates {
		dfa_to_nfa[state] = &nfapckg.State{
			Id:    state,
			Ntran: make(map[string][]*nfapckg.State), // заполняем потом
		}
	}

	for state := range dfa.dstates {
		for lit := range dfa.alphabet {
			if next_state, ok := dfa.dtran[state][lit]; ok {
				dfa_to_nfa[next_state].Ntran[lit] = append(
					dfa_to_nfa[next_state].Ntran[lit],
					dfa_to_nfa[state])
			}
		}
	}

	// теперь задаем начало и конец
	new_start := nfapckg.GenNewState()
	for acc := range dfa.accept_states {
		new_start.Epsilon = append(
			new_start.Epsilon,
			nfapckg.Epsilon_tran{
				To: dfa_to_nfa[acc],
			})
	}

	new_finish := dfa_to_nfa[dfa.start_state]

	return &nfapckg.NFA{
		Head: nfapckg.NfaNode{
			new_start,
			new_finish,
		},
	}
}

func epsilonClosure(input []*nfapckg.State) []*nfapckg.State {
	// стек для обхода
	stack := make([]*nfapckg.State, len(input))
	copy(stack, input)

	visited := make(map[*nfapckg.State]struct{})

	for len(stack) > 0 {
		// pop
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// если уже были в этом (state, pos)
		if _, ok := visited[cur]; ok {
			continue
		}
		visited[cur] = struct{}{}

		// идём по ε-переходам
		for _, tr := range cur.Epsilon {
			stack = append(stack, tr.To)
		}
	}

	// собираем результат
	result := make([]*nfapckg.State, 0, len(visited))
	for state := range visited {
		result = append(result, state)
	}

	return result
}

func moveLiteral(states []*nfapckg.State, symbol string) []*nfapckg.State {
	visited := make(map[*nfapckg.State]struct{})
	for _, state := range states {
		next_state := state.Ntran[symbol]
		for _, st := range next_state {
			visited[st] = struct{}{}
		}
	}

	result := make([]*nfapckg.State, 0, len(visited))
	for state := range visited {
		result = append(result, state)
	}

	return result
}

func keySet(states []*nfapckg.State) string {
	var sorted []*nfapckg.State
	for _, state := range states {
		sorted = append(sorted, state)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Id < sorted[j].Id
	})

	var key string
	for _, state := range sorted {
		key += strconv.Itoa(state.Id) + "|"
	}
	return key
}

func calcAlphabetForNfa(nfa *nfapckg.NFA) map[string]struct{} {
	// стек для обхода
	stack := []*nfapckg.State{nfa.Head.Start}
	visited := make(map[*nfapckg.State]struct{})

	alphabet := make(map[string]struct{})

	for len(stack) > 0 {
		// pop
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// если уже были в этом (state, pos)
		if _, ok := visited[cur]; ok {
			continue
		}
		visited[cur] = struct{}{}

		for lit, tr := range cur.Ntran {
			alphabet[lit] = struct{}{}
			for _, st := range tr {
				stack = append(stack, st)
			}
		}
		for _, tr := range cur.Epsilon {
			stack = append(stack, tr.To)
		}
	}

	return alphabet
}

func NfaToDFA(nfa *nfapckg.NFA) *DFA {
	start := []*nfapckg.State{nfa.Head.Start}
	start_state := epsilonClosure(start)

	stack := [][]*nfapckg.State{start_state}
	set_to_id := make(map[string]int)
	dtran := make(map[int]map[string]int)
	accept_states := make(set)

	alphabet := calcAlphabetForNfa(nfa)

	set_to_id[keySet(start_state)] = 0
	for len(stack) > 0 {
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		s_id := set_to_id[keySet(state)]

		if _, ok := dtran[s_id]; !ok {
			dtran[s_id] = make(map[string]int)
		}

		for _, st := range state {
			if st == nfa.Head.Finish {
				accept_states[s_id] = struct{}{}
				break
			}
		}

		for lit := range alphabet {
			r_states := epsilonClosure(moveLiteral(state, lit))
			if len(r_states) == 0 {
				continue
			}
			r_states_key := keySet(r_states)
			if _, ok := set_to_id[r_states_key]; !ok {
				new_id := len(set_to_id)
				set_to_id[r_states_key] = new_id
				stack = append(stack, r_states)
			}
			dtran[s_id][lit] = set_to_id[r_states_key]
		}
	}

	dstates := make([]set, len(set_to_id))
	// ВАЖНО мы не сохраняем множества нка, но вроде нам они не нужны
	for _, id := range set_to_id {
		dstates[id] = make(set)
	}

	// building dead_states
	dead_id := len(dstates)
	dstates = append(dstates, make(set))
	dtran[dead_id] = make(map[string]int)
	for state_id := 0; state_id < dead_id; state_id++ {
		if dtran[state_id] == nil {
			dtran[state_id] = make(map[string]int)
		}
		for lit := range alphabet {
			if _, ok := dtran[state_id][lit]; !ok {
				dtran[state_id][lit] = dead_id
			}
		}
	}
	for lit := range alphabet {
		dtran[dead_id][lit] = dead_id
	}
	return &DFA{
		alphabet:      alphabet,
		dstates:       dstates,
		dtran:         dtran,
		start_state:   0,
		accept_states: accept_states,
	}
}

func Inversion(dfa *DFA) *DFA {
	nfa := DfaToNfaNReverse(dfa)
	return NfaToDFA(nfa)
}
