package dfapckg

import (
	"slices"
	"testing"
)

func groupIDByState(groups []set, state int) int {
	for id, g := range groups {
		if _, ok := g[state]; ok {
			return id
		}
	}
	return -1
}

func TestExtractSortedAlphabet(t *testing.T) {
	alphabet := map[string]struct{}{
		"b": {},
		"a": {},
		"c": {},
	}

	got := extractSortedAlphabet(alphabet)
	want := []string{"a", "b", "c"}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected sorted alphabet: got=%v want=%v", got, want)
	}
}

func TestRefine_SplitsByTransitions(t *testing.T) {
	dfa := &DFA{
		alphabet: map[string]struct{}{"a": {}},
		dstates: []set{
			{0: {}}, // start
			{1: {}}, // accepting
			{2: {}}, // non-accepting dead
		},
		dtran: map[int]map[string]int{
			0: {"a": 1},
			1: {"a": 1},
			2: {"a": 2},
		},
		start_state:   0,
		accept_states: set{1: {}},
	}

	partition := firstPartition(dfa.dstates, dfa.accept_states)
	refined := refine(partition, dfa)

	if len(refined) != 3 {
		t.Fatalf("unexpected partition count after refine: got=%d want=%d (%v)", len(refined), 3, refined)
	}

	g0 := groupIDByState(refined, 0)
	g1 := groupIDByState(refined, 1)
	g2 := groupIDByState(refined, 2)
	if g0 == -1 || g1 == -1 || g2 == -1 {
		t.Fatalf("state missing in refined partition: g0=%d g1=%d g2=%d partition=%v", g0, g1, g2, refined)
	}
	if g0 == g2 {
		t.Fatalf("states 0 and 2 must be split by transitions, but got same group id=%d", g0)
	}
}

func TestDfaToMinDFA_MergesEquivalentAcceptingStates(t *testing.T) {
	dfa := &DFA{
		alphabet: map[string]struct{}{"a": {}, "b": {}},
		dstates: []set{
			{0: {}}, // start, non-accepting
			{1: {}}, // accepting
			{2: {}}, // accepting (equivalent to state 1)
			{3: {}}, // dead, non-accepting
		},
		dtran: map[int]map[string]int{
			0: {"a": 1, "b": 2},
			1: {"a": 1, "b": 1},
			2: {"a": 2, "b": 2},
			3: {"a": 3, "b": 3},
		},
		start_state:   0,
		accept_states: set{1: {}, 2: {}},
	}

	min := DfaToMinDFA(dfa)

	if len(min.mstates) != 3 {
		t.Fatalf("unexpected number of min states: got=%d want=%d (%v)", len(min.mstates), 3, min.mstates)
	}

	g0 := groupIDByState(min.mstates, 0)
	g1 := groupIDByState(min.mstates, 1)
	g2 := groupIDByState(min.mstates, 2)
	g3 := groupIDByState(min.mstates, 3)
	if g0 == -1 || g1 == -1 || g2 == -1 || g3 == -1 {
		t.Fatalf("state missing in min dfa groups: g0=%d g1=%d g2=%d g3=%d groups=%v", g0, g1, g2, g3, min.mstates)
	}

	if g1 != g2 {
		t.Fatalf("equivalent accepting states must be merged, got group(1)=%d group(2)=%d", g1, g2)
	}
	if g0 == g3 {
		t.Fatalf("start and dead states must not be merged, but both are in group %d", g0)
	}

	if min.start_state != g0 {
		t.Fatalf("unexpected min start state: got=%d want=%d", min.start_state, g0)
	}
	if len(min.accept_states) != 1 {
		t.Fatalf("unexpected number of min accept states: got=%d want=%d (%v)", len(min.accept_states), 1, min.accept_states)
	}
	if _, ok := min.accept_states[g1]; !ok {
		t.Fatalf("merged accepting group must be accepting: group=%d accept=%v", g1, min.accept_states)
	}

	if got := min.mtran[g0]["a"]; got != g1 {
		t.Fatalf("unexpected transition from start by 'a': got=%d want=%d", got, g1)
	}
	if got := min.mtran[g0]["b"]; got != g1 {
		t.Fatalf("unexpected transition from start by 'b': got=%d want=%d", got, g1)
	}
	if got := min.mtran[g3]["a"]; got != g3 {
		t.Fatalf("unexpected dead transition by 'a': got=%d want=%d", got, g3)
	}
	if got := min.mtran[g3]["b"]; got != g3 {
		t.Fatalf("unexpected dead transition by 'b': got=%d want=%d", got, g3)
	}
}
