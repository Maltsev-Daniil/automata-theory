package synttree

import "testing"

func TestTokenize(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{name: "empty", in: "", want: []string{}},
		{name: "single literal", in: "x", want: []string{"x"}},
		{name: "ab", in: "ab", want: []string{"a", "b"}},
		{name: "look", in: "look", want: []string{"l", "o", "o", "k"}},
		{name: "parentheses", in: "(a)", want: []string{"(", "a", ")"}},
		{name: "question", in: "a?", want: []string{"a", "?"}},
		{name: "or", in: "a|b", want: []string{"a", "|", "b"}},
		{name: "repeat", in: "a{12}", want: []string{"a", "12"}},
		{name: "klini", in: "...a", want: []string{"...", "a"}},
		{name: "capture group", in: "<name>a", want: []string{"name", "a"}},
		{name: "capture group unicode", in: "<имя>a", want: []string{"имя", "a"}},
		{name: "escaped ascii literal", in: "%a%", want: []string{"a"}},
		{name: "escaped unicode literal", in: "%я%", want: []string{"я"}},
		{name: "escaped percent literal", in: "%%%", want: []string{"%"}},
		{name: "unicode literals", in: "ёж", want: []string{"ё", "ж"}},
		{name: "mixed simple", in: "(a|b)?", want: []string{"(", "a", "|", "b", ")", "?"}},
		{name: "mixed with capture and repeat", in: "<id>{3}|z", want: []string{"id", "3", "|", "z"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := tokenize(tc.in)
			if err != nil {
				t.Fatalf("tokenize(%q) error: %v", tc.in, err)
			}

			if len(tokens) != len(tc.want) {
				t.Fatalf("token count mismatch: got %d, want %d", len(tokens), len(tc.want))
			}

			for i, tok := range tokens {
				if tok.value != tc.want[i] {
					t.Fatalf("token[%d] mismatch: got %q, want %q", i, tok.value, tc.want[i])
				}
			}
		})
	}
}

func TestBuildTree(t *testing.T) {
	type expectedNode struct {
		typ   TypeNode
		value string
		left  *expectedNode
		right *expectedNode
	}

	cases := []struct {
		name    string
		in      []Token
		want    *expectedNode
		wantErr bool
	}{
		{
			name: "single literal",
			in:   []Token{{value: "a", type_token: LITERAL}},
			want: &expectedNode{typ: LITERAL, value: "a"},
		},
		{
			name: "simple concat",
			in: []Token{
				{value: "a", type_token: LITERAL},
				{value: "+", type_token: OP_CONC},
				{value: "b", type_token: LITERAL},
			},
			want: &expectedNode{
				typ:   OP_CONC,
				value: "+",
				left:  &expectedNode{typ: LITERAL, value: "a"},
				right: &expectedNode{typ: LITERAL, value: "b"},
			},
		},
		{
			name: "or with concat precedence",
			in: []Token{
				{value: "a", type_token: LITERAL},
				{value: "|", type_token: OP_OR},
				{value: "b", type_token: LITERAL},
				{value: "+", type_token: OP_CONC},
				{value: "c", type_token: LITERAL},
			},
			want: &expectedNode{
				typ:   OP_OR,
				value: "|",
				left:  &expectedNode{typ: LITERAL, value: "a"},
				right: &expectedNode{
					typ:   OP_CONC,
					value: "+",
					left:  &expectedNode{typ: LITERAL, value: "b"},
					right: &expectedNode{typ: LITERAL, value: "c"},
				},
			},
		},
		{
			name: "question over grouped expression",
			in: []Token{
				{value: "(", type_token: LEFT_PAR},
				{value: "a", type_token: LITERAL},
				{value: "|", type_token: OP_OR},
				{value: "b", type_token: LITERAL},
				{value: ")", type_token: RIGHT_PAR},
				{value: "?", type_token: OP_QUESTION},
			},
			want: &expectedNode{
				typ:   OP_QUESTION,
				value: "?",
				left: &expectedNode{
					typ:   OP_OR,
					value: "|",
					left:  &expectedNode{typ: LITERAL, value: "a"},
					right: &expectedNode{typ: LITERAL, value: "b"},
				},
			},
		},
		{
			name: "repeat unfolds to concat chain",
			in: []Token{
				{value: "a", type_token: LITERAL},
				{value: "3", type_token: OP_REPEAT},
			},
			want: &expectedNode{
				typ:   OP_CONC,
				value: "+",
				left: &expectedNode{
					typ:   OP_CONC,
					value: "+",
					left:  &expectedNode{typ: LITERAL, value: "a"},
					right: &expectedNode{typ: LITERAL, value: "a"},
				},
				right: &expectedNode{typ: LITERAL, value: "a"},
			},
		},
		{
			name: "capture group in parentheses",
			in: []Token{
				{value: "(", type_token: LEFT_PAR},
				{value: "name", type_token: CAPTURE_GROUP},
				{value: "a", type_token: LITERAL},
				{value: ")", type_token: RIGHT_PAR},
			},
			want: &expectedNode{
				typ:   CAPTURE_GROUP,
				value: "name",
				left:  &expectedNode{typ: LITERAL, value: "a"},
			},
		},
		{
			name: "invalid capture placement",
			in: []Token{
				{value: "(", type_token: LEFT_PAR},
				{value: "a", type_token: LITERAL},
				{value: "|", type_token: OP_OR},
				{value: "name", type_token: CAPTURE_GROUP},
				{value: "b", type_token: LITERAL},
				{value: ")", type_token: RIGHT_PAR},
			},
			wantErr: true,
		},
		{
			name: "invalid repeat value",
			in: []Token{
				{value: "a", type_token: LITERAL},
				{value: "0", type_token: OP_REPEAT},
			},
			wantErr: true,
		},
	}

	var assertNode func(t *testing.T, got *Node, want *expectedNode)
	assertNode = func(t *testing.T, got *Node, want *expectedNode) {
		if want == nil {
			if got != nil {
				t.Fatalf("expected nil node, got %v (%q)", got.Type_node, got.Value)
			}
			return
		}
		if got == nil {
			t.Fatalf("expected node %v (%q), got nil", want.typ, want.value)
		}
		if got.Type_node != want.typ || got.Value != want.value {
			t.Fatalf("expected node %v (%q), got %v (%q)", want.typ, want.value, got.Type_node, got.Value)
		}
		assertNode(t, got.Left, want.left)
		assertNode(t, got.Right, want.right)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := buildTree(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("buildTree() unexpected error: %v", err)
			}
			assertNode(t, tree.Root, tc.want)
		})
	}
}
