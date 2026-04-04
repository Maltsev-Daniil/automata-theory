package dfapckg

import (
	"slices"
	"strconv"
)

type MinDFA struct {
	alphabet      map[string]struct{}
	mstates       []set
	mtran         map[int]map[string]int
	start_state   int
	accept_states set
}

func extractSortedAlphabet(alph map[string]struct{}) (result []string) {
	result = make([]string, 0, len(alph))
	for k := range alph {
		result = append(result, k)
	}
	slices.Sort(result)
	return result
}

func firstPartition(dstates []set, acc_states set) []set {
	result := []set{
		make(set), // non-accepting
		make(set), // accepting
	}
	for id := range dstates {
		if _, ok := acc_states[id]; ok {
			result[1][id] = struct{}{}
		} else {
			result[0][id] = struct{}{}
		}
	}
	return result
}

func buildStateToGroup(partition []set) map[int]int {
	result := make(map[int]int)

	for group_id, group := range partition {
		for state := range group {
			result[state] = group_id
		}
	}

	return result
}

func computeSignature(
	state int,
	alphabet []string,
	dtran map[int]map[string]int,
	state_to_group map[int]int,
) string {

	sig := ""

	for _, a := range alphabet {
		next := dtran[state][a]
		group := state_to_group[next]

		sig += strconv.Itoa(group) + "|"
	}

	return sig
}

func splitGroup(
	group set,
	alphabet []string,
	dtran map[int]map[string]int,
	state_to_group map[int]int,
) []set {

	sig_map := make(map[string]set)

	for state := range group {
		sig := computeSignature(state, alphabet, dtran, state_to_group)

		if sig_map[sig] == nil {
			sig_map[sig] = make(set)
		}

		sig_map[sig][state] = struct{}{}
	}

	result := []set{}
	for _, subset := range sig_map {
		result = append(result, subset)
	}

	return result
}

func partitionsEqual(a, b []set) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !areStatesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func refine(partition []set, dfa *DFA) []set {
	alphabet := extractSortedAlphabet(dfa.alphabet)

	for {
		state_to_group := buildStateToGroup(partition)

		new_partition := []set{}

		for _, group := range partition {
			splits := splitGroup(group, alphabet, dfa.dtran, state_to_group)
			new_partition = append(new_partition, splits...)
		}

		if partitionsEqual(partition, new_partition) {
			break
		}

		partition = new_partition
	}

	return partition
}

func DfaToMinDFA(dfa *DFA) *MinDFA {
	init_partition := firstPartition(dfa.dstates, dfa.accept_states)
	mstates := refine(init_partition, dfa)
	state_to_group := buildStateToGroup(mstates)
	mtran := make(map[int]map[string]int)

	for i, group := range mstates {
		mtran[i] = make(map[string]int)

		var rep int
		for s := range group {
			rep = s
			break
		}

		for symbol := range dfa.alphabet {
			next := dfa.dtran[rep][symbol]
			mtran[i][symbol] = state_to_group[next]
		}
	}

	start := state_to_group[dfa.start_state]

	// 4. новые accepting
	accept := make(map[int]struct{})

	for i, group := range mstates {
		for s := range group {
			if _, ok := dfa.accept_states[s]; ok {
				accept[i] = struct{}{}
				break
			}
		}
	}

	return &MinDFA{
		alphabet:      dfa.alphabet,
		mstates:       mstates,
		mtran:         mtran,
		start_state:   start,
		accept_states: accept,
	}
}
