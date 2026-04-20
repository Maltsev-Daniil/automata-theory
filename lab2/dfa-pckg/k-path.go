package dfapckg

import "sort"

type GNFA struct {
	start_state  int
	finish_state int
	gstates      set
	gtrans       map[int]map[int]string
}

func fromDfaToGnfa(dfa *DFA) *GNFA {
	gstates := make(set)
	gtrans := make(map[int]map[int]string)

	for state := range dfa.dstates {
		gstates[state] = struct{}{}
	}

	var max_id int
	for i := range dfa.dstates {
		max_id = max(max_id, i)
	}

	// создаем единственные новые
	// нач и кон состояния
	start_state := max_id + 1
	gstates[start_state] = struct{}{}

	finish_state := max_id + 2
	gstates[finish_state] = struct{}{}

	// init
	for i := range gstates {
		gtrans[i] = make(map[int]string)
		for j := range gstates {
			gtrans[i][j] = "∅"
		}
	}

	for id, tran := range dfa.dtran {
		for lit, state := range tran {
			if gtrans[id][state] == "∅" {
				gtrans[id][state] = lit
			} else {
				gtrans[id][state] = "(" + gtrans[id][state] + "|" + lit + ")"
			}
		}
	}

	gtrans[start_state][dfa.start_state] = "ε"
	for state := range dfa.accept_states {
		gtrans[state][finish_state] = "ε"
	}

	return &GNFA{
		start_state:  start_state,
		finish_state: finish_state,
		gstates:      gstates,
		gtrans:       gtrans,
	}
}

func union(a, b string) string {
	if a == b {
		return a
	} else if a == "∅" {
		return b
	} else if b == "∅" {
		return a
	}
	return "(" + a + "|" + b + ")"
}

func klini(a string) string {
	if a == "∅" || a == "ε" {
		return "ε"
	}
	return "(" + a + ")*"
}

func conc(a, b string) string {
	if a == "∅" || b == "∅" {
		return "∅"
	} else if a == "ε" {
		return b
	} else if b == "ε" {
		return a
	}
	return "(" + a + b + ")"
}

func Kpath(dfa *DFA) (res string) {
	gnfa := fromDfaToGnfa(dfa)
	states := make([]int, 0, len(gnfa.gstates))
	for state := range gnfa.gstates {
		states = append(states, state)
	}
	sort.Ints(states)
	for _, k := range states {
		for _, i := range states {
			for _, j := range states {
				rik := gnfa.gtrans[i][k]
				rkk := gnfa.gtrans[k][k]
				rkj := gnfa.gtrans[k][j]

				via := conc(rik, conc(klini(rkk), rkj))

				gnfa.gtrans[i][j] = union(gnfa.gtrans[i][j], via)
			}
		}
	}
	return gnfa.gtrans[gnfa.start_state][gnfa.finish_state]
}
