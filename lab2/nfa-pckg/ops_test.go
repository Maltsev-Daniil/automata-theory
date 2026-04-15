package nfapckg

import (
	"testing"
)

func TestOpKlini(t *testing.T) {
	nfa := CompileNFA("a...")

	if !nfa.matchNFA("aaa") {
		t.Fatalf("expected matchNFA to return true for OP_KLINI")
	}
	if nfa.matchNFA("b") {
		t.Fatalf("expected matchNFA to return false for OP_KLINI")
	}
}

func TestOpConc(t *testing.T) {
	nfa := CompileNFA("ab")

	if !nfa.matchNFA("ab") {
		t.Fatalf("expected matchNFA to return true for OP_CONC")
	}
	if nfa.matchNFA("a") {
		t.Fatalf("expected matchNFA to return false for OP_CONC")
	}
}

func TestOpQuestion(t *testing.T) {
	nfa := CompileNFA("a?")

	if !nfa.matchNFA("a") {
		t.Fatalf("expected matchNFA to return true for OP_QUESTION")
	}
	if nfa.matchNFA("aa") {
		t.Fatalf("expected matchNFA to return false for OP_QUESTION")
	}
}

func TestOpRepeat(t *testing.T) {
	nfa := CompileNFA("a{3}")

	if !nfa.matchNFA("aaa") {
		t.Fatalf("expected matchNFA to return true for OP_REPEAT")
	}
	if nfa.matchNFA("aa") {
		t.Fatalf("expected matchNFA to return false for OP_REPEAT")
	}
}

func TestOpOr(t *testing.T) {
	nfa := CompileNFA("a|b")

	if !nfa.matchNFA("b") {
		t.Fatalf("expected matchNFA to return true for OP_OR")
	}
	if nfa.matchNFA("c") {
		t.Fatalf("expected matchNFA to return false for OP_OR")
	}
}

func TestCaptureGroupById(t *testing.T) {
	nfa := CompileNFA("(<grp>ab)")
	matchRes := nfa.matchNFAWithCapture("ab")

	got, err := matchRes.GroupById(1)
	if err != nil {
		t.Fatalf("expected GroupById(1) to succeed: %v", err)
	}
	if got != "ab" {
		t.Fatalf("GroupById(1) = %q, want %q", got, "ab")
	}

	_, err = matchRes.GroupById(2)
	if err == nil {
		t.Fatalf("expected GroupById(2) to fail")
	}
}

func TestCaptureGroupByName(t *testing.T) {
	nfa := CompileNFA("(<grp>ab)")
	matchRes := nfa.matchNFAWithCapture("ab")

	got, err := matchRes.GroupByName("grp")
	if err != nil {
		t.Fatalf("expected GroupByName(%q) to succeed: %v", "grp", err)
	}
	if got != "ab" {
		t.Fatalf("GroupByName(%q) = %q, want %q", "grp", got, "ab")
	}

	_, err = matchRes.GroupByName("missing")
	if err == nil {
		t.Fatalf("expected GroupByName(%q) to fail", "missing")
	}
}

func TestCombinedOperations(t *testing.T) {
	nfa := CompileNFA("((<outer>(a|b{2}))c?)...d")

	if !nfa.matchNFA("abbcad") {
		t.Fatalf("expected matchNFA to return true for combined operations")
	}
	if nfa.matchNFA("abbcabb") {
		t.Fatalf("expected matchNFA to return false for combined operations")
	}

	matchRes := nfa.matchNFAWithCapture("abbcad")
	got, err := matchRes.GroupByName("outer")
	if err != nil {
		t.Fatalf("expected GroupByName(%q) to succeed: %v", "outer", err)
	}
	if got != "a" {
		t.Fatalf("GroupByName(%q) = %q, want %q", "outer", got, "a")
	}
}
