package nfapckg

import (
	"errors"
	"reglib/synttree-pckg"
)

type Epsilon_tran struct {
	To  *State
	cap *cap_event // nil if none
}
type cap_event struct {
	name   string
	finish bool // 0 - start; 1 - finish
}
type State struct {
	Id      int
	Ntran   map[string][]*State
	Epsilon []Epsilon_tran
}
type NfaNode struct {
	Start  *State
	Finish *State
}
type NFA struct {
	Head NfaNode

	group_order []string
	name_to_id  map[string]int
}

var state_id int
var group_order []string
var name_to_id map[string]int

func GenNewState() *State {
	s := &State{
		Id:      state_id,
		Ntran:   make(map[string][]*State),
		Epsilon: make([]Epsilon_tran, 0),
	}
	state_id++
	return s
}

func buildLiteral(lit string) NfaNode {
	s := GenNewState()
	f := GenNewState()
	s.Ntran[lit] = []*State{f}
	return NfaNode{
		Start:  s,
		Finish: f,
	}
}

func buildConc(lhs, rhs NfaNode) NfaNode {
	lhs.Finish.Epsilon = append(lhs.Finish.Epsilon, Epsilon_tran{
		To:  rhs.Start,
		cap: nil,
	})
	return NfaNode{
		lhs.Start,
		rhs.Finish,
	}
}

func buildOr(lhs, rhs NfaNode) NfaNode {
	new_s := GenNewState()
	new_f := GenNewState()

	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To:  lhs.Start,
		cap: nil,
	})
	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To:  rhs.Start,
		cap: nil,
	})

	lhs.Finish.Epsilon = append(lhs.Finish.Epsilon, Epsilon_tran{
		To:  new_f,
		cap: nil,
	})
	rhs.Finish.Epsilon = append(rhs.Finish.Epsilon, Epsilon_tran{
		To:  new_f,
		cap: nil,
	})

	return NfaNode{
		Start:  new_s,
		Finish: new_f,
	}
}

func buildKlini(node NfaNode) NfaNode {
	new_s := GenNewState()
	new_f := GenNewState()

	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To:  node.Start,
		cap: nil,
	})
	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To:  new_f,
		cap: nil,
	})

	node.Finish.Epsilon = append(node.Finish.Epsilon, Epsilon_tran{
		To:  new_f,
		cap: nil,
	})
	node.Finish.Epsilon = append(node.Finish.Epsilon, Epsilon_tran{
		To:  node.Start,
		cap: nil,
	})

	return NfaNode{
		Start:  new_s,
		Finish: new_f,
	}
}

func buildOptional(node NfaNode) NfaNode {
	new_s := GenNewState()
	new_f := GenNewState()

	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To:  node.Start,
		cap: nil,
	})
	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To:  new_f,
		cap: nil,
	})

	node.Finish.Epsilon = append(node.Finish.Epsilon, Epsilon_tran{
		To:  new_f,
		cap: nil,
	})

	return NfaNode{
		Start:  new_s,
		Finish: new_f,
	}
}

func buildCapture(node *synttree.Node) NfaNode {
	child := buildNFA(node.Left)

	new_s := GenNewState()
	new_f := GenNewState()

	// чтобы можно было по айди находить
	if name_to_id == nil {
		name_to_id = make(map[string]int)
	}
	if _, exists := name_to_id[node.Value]; !exists {
		name_to_id[node.Value] = len(group_order)
		group_order = append(group_order, node.Value)
	}

	new_s.Epsilon = append(new_s.Epsilon, Epsilon_tran{
		To: child.Start,
		cap: &cap_event{
			name:   node.Value,
			finish: false,
		},
	})
	child.Finish.Epsilon = append(child.Finish.Epsilon, Epsilon_tran{
		To: new_f,
		cap: &cap_event{
			name:   node.Value,
			finish: true,
		}})
	return NfaNode{
		Start:  new_s,
		Finish: new_f,
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
