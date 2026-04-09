package dfapckg

// работаем не с minDFA
// тк нам нужны dead_states
// Dfa must be complete
func Complement(dfa *DFA) *DFA {
	// разворачиваем принимающие состояния
	// остальное мы просто дип копи
	c_accept_states := make(set)
	for id := range dfa.dstates {
		if _, ok := dfa.accept_states[id]; !ok {
			c_accept_states[id] = struct{}{}
		}
	}

	c_dstates := make([]set, len(dfa.dstates))
	for id_new := range c_dstates {
		c_dstates[id_new] = make(set)
		for nfa_state := range dfa.dstates[id_new] {
			c_dstates[id_new][nfa_state] = struct{}{}
		}
	}

	c_dtran := make(map[int]map[string]int)
	for state, trans := range dfa.dtran {
		if c_dtran[state] == nil {
			c_dtran[state] = make(map[string]int)
		}
		for lit, to := range trans {
			// дипкопи, тк инт по значению копируется
			c_dtran[state][lit] = to
		}
	}
	// алфавит конст
	c_alphabet := dfa.alphabet

	return &DFA{
		alphabet:      c_alphabet,
		dstates:       c_dstates,
		dtran:         c_dtran,
		start_state:   dfa.start_state,
		accept_states: c_accept_states,
	}
}
