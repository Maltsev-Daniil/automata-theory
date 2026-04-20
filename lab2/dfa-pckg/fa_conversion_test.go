package dfapckg

import "testing"

func mustCompileDFA(t *testing.T, expr string) *DFA {
	t.Helper()

	dfa := CompileDFA(expr)
	if dfa == nil {
		t.Fatalf("CompileDFA(%q) returned nil", expr)
	}
	return dfa
}

func TestComplementIsomorphic(t *testing.T) {
	source := mustCompileDFA(t, "a?")
	got := Complement(source)

	right := mustCompileDFA(t, "aaa...")
	if !Isomorphic(got, right) {
		t.Fatalf("Complement DFA is not isomorphic to expected DFA")
	}

	wrong := mustCompileDFA(t, "a...")
	if Isomorphic(got, wrong) {
		t.Fatalf("Complement DFA unexpectedly isomorphic to wrong expected DFA")
	}
}

func TestInversionIsomorphic(t *testing.T) {
	source := mustCompileDFA(t, "ab")
	got := Inversion(source)

	right := mustCompileDFA(t, "ba")
	if !Isomorphic(got, right) {
		t.Fatalf("Inversion DFA is not isomorphic to expected DFA")
	}

	wrong := mustCompileDFA(t, "ab")
	if Isomorphic(got, wrong) {
		t.Fatalf("Inversion DFA unexpectedly isomorphic to wrong expected DFA")
	}
}

func TestKPathRecompileIsomorphic(t *testing.T) {
	source := mustCompileDFA(t, "a")
	regex := Kpath(source)
	got := mustCompileDFA(t, regex)

	if !Isomorphic(got, source) {
		t.Fatalf("DFA compiled from Kpath regex is not isomorphic to source DFA; regex=%q", regex)
	}

	wrong := mustCompileDFA(t, "aa")
	if Isomorphic(got, wrong) {
		t.Fatalf("DFA compiled from Kpath regex unexpectedly isomorphic to wrong DFA; regex=%q", regex)
	}
}
