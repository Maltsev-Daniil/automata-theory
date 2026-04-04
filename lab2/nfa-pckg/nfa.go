package nfapckg

import (
	"errors"
	"reglib/synttree-pckg"
)

type epsilon_tran struct {
	to  *state
	cap *cap_event // nil if none
}
type cap_event struct {
	name   string
	finish bool // 0 - start; 1 - finish
}
type state struct {
	id      int
	ntran   map[string][]*state
	epsilon []epsilon_tran
}
type NfaNode struct {
	start  *state
	finish *state
}
type NFA struct {
	head NfaNode

	group_order []string
	name_to_id  map[string]int
}

var state_id int
var group_order []string
var name_to_id map[string]int

func genNewState() *state {
	s := &state{
		id:      state_id,
		ntran:   make(map[string][]*state),
		epsilon: make([]epsilon_tran, 0),
	}
	state_id++
	return s
}

func buildLiteral(lit string) NfaNode {
	s := genNewState()
	f := genNewState()
	s.ntran[lit] = []*state{f}
	return NfaNode{
		start:  s,
		finish: f,
	}
}

func buildConc(lhs, rhs NfaNode) NfaNode {
	lhs.finish.epsilon = append(lhs.finish.epsilon, epsilon_tran{
		to:  rhs.start,
		cap: nil,
	})
	return NfaNode{
		lhs.start,
		rhs.finish,
	}
}

func buildOr(lhs, rhs NfaNode) NfaNode {
	new_s := genNewState()
	new_f := genNewState()

	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to:  lhs.start,
		cap: nil,
	})
	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to:  rhs.start,
		cap: nil,
	})

	lhs.finish.epsilon = append(lhs.finish.epsilon, epsilon_tran{
		to:  new_f,
		cap: nil,
	})
	rhs.finish.epsilon = append(rhs.finish.epsilon, epsilon_tran{
		to:  new_f,
		cap: nil,
	})

	return NfaNode{
		start:  new_s,
		finish: new_f,
	}
}

func buildKlini(node NfaNode) NfaNode {
	new_s := genNewState()
	new_f := genNewState()

	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to:  node.start,
		cap: nil,
	})
	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to:  new_f,
		cap: nil,
	})

	node.finish.epsilon = append(node.finish.epsilon, epsilon_tran{
		to:  new_f,
		cap: nil,
	})
	node.finish.epsilon = append(node.finish.epsilon, epsilon_tran{
		to:  node.start,
		cap: nil,
	})

	return NfaNode{
		start:  new_s,
		finish: new_f,
	}
}

func buildOptional(node NfaNode) NfaNode {
	new_s := genNewState()
	new_f := genNewState()

	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to:  node.start,
		cap: nil,
	})
	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to:  new_f,
		cap: nil,
	})

	node.finish.epsilon = append(node.finish.epsilon, epsilon_tran{
		to:  new_f,
		cap: nil,
	})

	return NfaNode{
		start:  new_s,
		finish: new_f,
	}
}

func buildCapture(node *synttree.Node) NfaNode {
	child := buildNFA(node.Left)

	new_s := genNewState()
	new_f := genNewState()

	// чтобы можно было по айди находить
	if name_to_id == nil {
		name_to_id = make(map[string]int)
	}
	if _, exists := name_to_id[node.Value]; !exists {
		name_to_id[node.Value] = len(group_order)
		group_order = append(group_order, node.Value)
	}

	new_s.epsilon = append(new_s.epsilon, epsilon_tran{
		to: child.start,
		cap: &cap_event{
			name:   node.Value,
			finish: false,
		},
	})
	child.finish.epsilon = append(child.finish.epsilon, epsilon_tran{
		to: new_f,
		cap: &cap_event{
			name:   node.Value,
			finish: true,
		}})
	return NfaNode{
		start:  new_s,
		finish: new_f,
	}
}

func buildNFA(root *synttree.Node) NfaNode {
	switch root.Type_node {
	case synttree.LITERAL:
		return buildLiteral(root.Value)
	case synttree.OP_CONC:
		left := buildNFA(root.Left)
		right := buildNFA(root.Right)
		return buildConc(left, right)
	case synttree.OP_OR:
		left := buildNFA(root.Left)
		right := buildNFA(root.Right)
		return buildOr(left, right)
	case synttree.OP_KLINI:
		left := buildNFA(root.Left)
		return buildKlini(left)
	case synttree.OP_QUESTION:
		left := buildNFA(root.Left)
		return buildOptional(left)
	case synttree.CAPTURE_GROUP:
		// передаем ноду целиком тк нам нужна
		// информация по группе захвата
		return buildCapture(root)
	default:
		panic(errors.New("buildNFA: invalid node type"))
	}
}

func treeToNFA(tree *synttree.Tree) *NFA {
	return &NFA{
		buildNFA(tree.Root),
		group_order,
		name_to_id,
	}
}
