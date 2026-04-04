package dfapckg

import (
	synttree "reglib/synttree-pckg"
	"testing"
)

func mustSetEqual(t *testing.T, got, want set) {
	if len(got) != len(want) {
		t.Fatalf("set size mismatch: got=%v want=%v", got, want)
	}
	for k := range want {
		if _, ok := got[k]; !ok {
			t.Fatalf("missing element %d in set: got=%v want=%v", k, got, want)
		}
	}
}

func TestCalcInfoNodes(t *testing.T) {
	// a + (b?)
	a := &synttree.Node{Type_node: synttree.LITERAL, Value: "a"}
	b := &synttree.Node{Type_node: synttree.LITERAL, Value: "b"}
	q := &synttree.Node{Type_node: synttree.OP_QUESTION, Value: "?", Left: b}
	root := &synttree.Node{Type_node: synttree.OP_CONC, Value: "+", Left: a, Right: q}

	info := make(map[*synttree.Node]*InfoNodes)
	pos_to_literal := make(map[int]string)
	pos := 0

	calcInfoNodes(root, info, pos_to_literal, &pos)

	if pos != 2 {
		t.Fatalf("expected 2 positions, got %d", pos)
	}
	want_pos_to_literal := map[int]string{
		0: "a",
		1: "b",
	}
	if len(pos_to_literal) != len(want_pos_to_literal) {
		t.Fatalf("unexpected pos_to_literal size: got %d want %d (%#v)", len(pos_to_literal), len(want_pos_to_literal), pos_to_literal)
	}
	for p, lit := range want_pos_to_literal {
		if got_lit, ok := pos_to_literal[p]; !ok || got_lit != lit {
			t.Fatalf("unexpected pos_to_literal[%d]: got %q want %q (full=%#v)", p, got_lit, lit, pos_to_literal)
		}
	}

	if info[a].nullable {
		t.Fatalf("literal node should be non-nullable")
	}
	mustSetEqual(t, info[a].firstpos, set{0: {}})
	mustSetEqual(t, info[a].lastpos, set{0: {}})

	if !info[q].nullable {
		t.Fatalf("question node should be nullable")
	}
	mustSetEqual(t, info[q].firstpos, set{1: {}})
	mustSetEqual(t, info[q].lastpos, set{1: {}})

	if info[root].nullable {
		t.Fatalf("concat root should be non-nullable")
	}
	mustSetEqual(t, info[root].firstpos, set{0: {}})
	mustSetEqual(t, info[root].lastpos, set{0: {}, 1: {}})
}

func TestCalcFollowPos(t *testing.T) {
	// a + (b...)
	a := &synttree.Node{Type_node: synttree.LITERAL, Value: "a"}
	b := &synttree.Node{Type_node: synttree.LITERAL, Value: "b"}
	k := &synttree.Node{Type_node: synttree.OP_KLINI, Value: "...", Left: b}
	root := &synttree.Node{Type_node: synttree.OP_CONC, Value: "+", Left: a, Right: k}

	info := make(map[*synttree.Node]*InfoNodes)
	pos_to_literal := make(map[int]string)
	follow_pos := make(map[int]set)
	pos := 0

	calcInfoNodes(root, info, pos_to_literal, &pos)
	calcFollowPos(root, info, follow_pos)

	// From concat: follow(0) contains firstpos(right) -> {1}
	mustSetEqual(t, follow_pos[0], set{1: {}})
	// From klini: follow(1) contains firstpos(child) -> {1}
	mustSetEqual(t, follow_pos[1], set{1: {}})
}

func TestCalcAlphabet(t *testing.T) {
	// (a | c) + b + #
	a := &synttree.Node{Type_node: synttree.LITERAL, Value: "a"}
	c := &synttree.Node{Type_node: synttree.LITERAL, Value: "c"}
	b := &synttree.Node{Type_node: synttree.LITERAL, Value: "b"}
	hash := &synttree.Node{Type_node: synttree.SHEBANG, Value: "#"}
	or := &synttree.Node{Type_node: synttree.OP_OR, Value: "|", Left: a, Right: c}
	concat_1 := &synttree.Node{Type_node: synttree.OP_CONC, Value: "+", Left: or, Right: b}
	root := &synttree.Node{Type_node: synttree.OP_CONC, Value: "+", Left: concat_1, Right: hash}

	alphabet := make(map[string]struct{})
	calcAlphabet(root, alphabet)

	if len(alphabet) != 3 {
		t.Fatalf("expected 3 literals in alphabet, got %d (%v)", len(alphabet), alphabet)
	}
	for _, lit := range []string{"a", "b", "c"} {
		if _, ok := alphabet[lit]; !ok {
			t.Fatalf("expected literal %q in alphabet, got %v", lit, alphabet)
		}
	}
	if _, ok := alphabet["#"]; ok {
		t.Fatalf("shebang should not be part of alphabet: %v", alphabet)
	}
}

func TestTreeToDFA_Linear(t *testing.T) {
	// a + #
	a := &synttree.Node{Type_node: synttree.LITERAL, Value: "a"}
	hash := &synttree.Node{Type_node: synttree.SHEBANG, Value: "#"}
	root := &synttree.Node{Type_node: synttree.OP_CONC, Value: "+", Left: a, Right: hash}
	tree := &synttree.Tree{Root: root}

	dfa := treeToDFA(tree)

	if dfa.start_state != 0 {
		t.Fatalf("unexpected start_state: got %d want %d", dfa.start_state, 0)
	}
	if len(dfa.alphabet) != 1 {
		t.Fatalf("unexpected alphabet size: got %d want %d (%v)", len(dfa.alphabet), 1, dfa.alphabet)
	}
	if _, ok := dfa.alphabet["a"]; !ok {
		t.Fatalf("alphabet must contain literal %q: %v", "a", dfa.alphabet)
	}
	if len(dfa.dstates) != 2 {
		t.Fatalf("unexpected number of dstates: got %d want %d", len(dfa.dstates), 2)
	}

	dead_id := len(dfa.dstates)

	if got := dfa.dtran[0]["a"]; got != 1 {
		t.Fatalf("unexpected transition dtran[0][a]: got %d want %d", got, 1)
	}
	if got := dfa.dtran[1]["a"]; got != dead_id {
		t.Fatalf("unexpected transition dtran[1][a]: got %d want %d", got, dead_id)
	}
	if got := dfa.dtran[dead_id]["a"]; got != dead_id {
		t.Fatalf("unexpected transition dtran[dead][a]: got %d want %d", got, dead_id)
	}

	mustSetEqual(t, dfa.accept_states, set{1: {}})
}

func TestTreeToDFA_Or(t *testing.T) {
	// (a | b) + #
	a := &synttree.Node{Type_node: synttree.LITERAL, Value: "a"}
	b := &synttree.Node{Type_node: synttree.LITERAL, Value: "b"}
	hash := &synttree.Node{Type_node: synttree.SHEBANG, Value: "#"}
	or := &synttree.Node{Type_node: synttree.OP_OR, Value: "|", Left: a, Right: b}
	root := &synttree.Node{Type_node: synttree.OP_CONC, Value: "+", Left: or, Right: hash}
	tree := &synttree.Tree{Root: root}

	dfa := treeToDFA(tree)

	if dfa.start_state != 0 {
		t.Fatalf("unexpected start_state: got %d want %d", dfa.start_state, 0)
	}
	if len(dfa.alphabet) != 2 {
		t.Fatalf("unexpected alphabet size: got %d want %d (%v)", len(dfa.alphabet), 2, dfa.alphabet)
	}
	for _, lit := range []string{"a", "b"} {
		if _, ok := dfa.alphabet[lit]; !ok {
			t.Fatalf("alphabet must contain literal %q: %v", lit, dfa.alphabet)
		}
	}

	// start state must branch by both literals to the same accept state in this construction
	next_by_a, ok_a := dfa.dtran[0]["a"]
	next_by_b, ok_b := dfa.dtran[0]["b"]
	if !ok_a || !ok_b {
		t.Fatalf("missing transitions from start state: %v", dfa.dtran[0])
	}
	if next_by_a != next_by_b {
		t.Fatalf("expected same target for a and b from start, got a->%d b->%d", next_by_a, next_by_b)
	}
	if _, ok := dfa.accept_states[next_by_a]; !ok {
		t.Fatalf("target state for start transitions must be accepting: target=%d accept=%v", next_by_a, dfa.accept_states)
	}

	dead_id := len(dfa.dstates)
	if got := dfa.dtran[dead_id]["a"]; got != dead_id {
		t.Fatalf("unexpected dead transition for a: got %d want %d", got, dead_id)
	}
	if got := dfa.dtran[dead_id]["b"]; got != dead_id {
		t.Fatalf("unexpected dead transition for b: got %d want %d", got, dead_id)
	}
}
