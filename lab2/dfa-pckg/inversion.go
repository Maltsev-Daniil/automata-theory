package dfapckg

import "reglib/nfa-pckg"

//type epsilon_tran struct {
//	to  *state
//	cap *cap_event // nil if none
//}
//type cap_event struct {
//	name   string
//	finish bool // 0 - start; 1 - finish
//}
//type state struct {
//	id      int
//	ntran   map[string][]*state
//	epsilon []epsilon_tran
//}
//type NfaNode struct {
//	start  *state
//	finish *state
//}
//type NFA struct {
//	head NfaNode
//
//	group_order []string
//	name_to_id  map[string]int
//}
//
//var state_id int
//var group_order []string
//var name_to_id map[string]int

//type set map[int]struct{}
//
//type DFA struct {
//	alphabet      map[string]struct{}
//	dstates       []set
//	dtran         map[int]map[string]int
//	start_state   int
//	accept_states set
//}
//
//type InfoNodes struct {
//	nullable bool
//	firstpos set
//	lastpos  set
//}

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

func Inversion(dfa *DFA) *nfapckg.NFA {
	nfa := DfaToNfaNReverse(dfa)
	return nfa
}
