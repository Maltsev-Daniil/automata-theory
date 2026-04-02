package dfapckg

import (
	"fmt"
	synttree "reglib/synttree-pckg"
)

type set map[uint]struct{}

type DFA struct {
	alphabet      map[string]struct{}
	dstates       []set
	dtran         map[uint]map[string]uint
	start_state   uint
	accept_states set
}

type InfoNodes struct {
	nullable bool
	firstpos set
	lastpos  set
}

func mergeSets(a, b set) set {
	result := make(set)
	for k := range a {
		result[k] = struct{}{}
	}
	for k := range b {
		result[k] = struct{}{}
	}
	return result
}

func calcInfoNodes(node *synttree.Node, info_nodes map[*synttree.Node]*InfoNodes, pos_to_literal map[uint]string, pos *uint) {
	if node == nil {
		return
	}
	calcInfoNodes(node.Left, info_nodes, pos_to_literal, pos)
	calcInfoNodes(node.Right, info_nodes, pos_to_literal, pos)

	switch node.Type_node {
	case synttree.LITERAL, synttree.SHEBANG:
		info_nodes[node] = &InfoNodes{
			nullable: false,
			firstpos: map[uint]struct{}{*pos: {}},
			lastpos:  map[uint]struct{}{*pos: {}},
		}
		pos_to_literal[*pos] = node.Value
		*pos++
	case synttree.OP_OR:
		info_nodes[node] = &InfoNodes{
			nullable: info_nodes[node.Left].nullable || info_nodes[node.Right].nullable,
			firstpos: mergeSets(info_nodes[node.Left].firstpos, info_nodes[node.Right].firstpos),
			lastpos:  mergeSets(info_nodes[node.Left].lastpos, info_nodes[node.Right].lastpos),
		}
	case synttree.OP_CONC:
		info_nodes[node] = &InfoNodes{
			nullable: info_nodes[node.Left].nullable && info_nodes[node.Right].nullable,
		}

		if info_nodes[node.Left].nullable {
			info_nodes[node].firstpos = mergeSets(
				info_nodes[node.Left].firstpos,
				info_nodes[node.Right].firstpos,
			)
		} else {
			info_nodes[node].firstpos = info_nodes[node.Left].firstpos
		}

		if info_nodes[node.Right].nullable {
			info_nodes[node].lastpos = mergeSets(
				info_nodes[node.Left].lastpos,
				info_nodes[node.Right].lastpos,
			)
		} else {
			info_nodes[node].lastpos = info_nodes[node.Right].lastpos
		}
	case synttree.OP_KLINI:
		info_nodes[node] = &InfoNodes{
			nullable: true,
			firstpos: info_nodes[node.Left].firstpos,
			lastpos:  info_nodes[node.Left].lastpos,
		}
	case synttree.OP_QUESTION:
		info_nodes[node] = &InfoNodes{
			nullable: true,
			firstpos: info_nodes[node.Left].firstpos,
			lastpos:  info_nodes[node.Left].lastpos,
		}

	// we don't have capture-groups in DFA so we ignore
	case synttree.CAPTURE_GROUP:
		info_nodes[node] = &InfoNodes{
			nullable: info_nodes[node.Left].nullable,
			firstpos: info_nodes[node.Left].firstpos,
			lastpos:  info_nodes[node.Left].lastpos,
		}

	default:
		panic(fmt.Sprintf("calcInfoNodes: operand %v is not supported", node.Type_node))
	}
}

func cycleForFollowPos(from set, to set, follow_pos map[uint]set) {
	for i := range from {
		if follow_pos[i] == nil {
			follow_pos[i] = make(set)
		}
		for j := range to {
			// нас интересуют только наличие ключа
			// на сами значения пофик
			// тк это множества
			follow_pos[i][j] = struct{}{}
		}
	}
}

func calcFollowPos(root *synttree.Node, info_nodes map[*synttree.Node]*InfoNodes, follow_pos map[uint]set) {
	if root == nil {
		return
	}
	calcFollowPos(root.Left, info_nodes, follow_pos)
	calcFollowPos(root.Right, info_nodes, follow_pos)

	switch root.Type_node {
	case synttree.OP_CONC:
		cycleForFollowPos(
			info_nodes[root.Left].lastpos,
			info_nodes[root.Right].firstpos,
			follow_pos,
		)
	case synttree.OP_KLINI:
		cycleForFollowPos(
			info_nodes[root.Left].lastpos,
			info_nodes[root.Left].firstpos,
			follow_pos,
		)
	}
}

func calcAlphabet(root *synttree.Node, alphabet map[string]struct{}) {
	if root == nil {
		return
	}
	calcAlphabet(root.Left, alphabet)
	calcAlphabet(root.Right, alphabet)

	if root.Type_node == synttree.LITERAL {
		alphabet[root.Value] = struct{}{}
	}
}

func areStatesEqual(a, b set) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if _, ok := b[i]; !ok {
			return false
		}
	}
	return true
}

func findState(dstates []set, trans_by_lit set) int {
	for i := range dstates {
		if areStatesEqual(dstates[i], trans_by_lit) {
			return i
		}
	}
	return -1
}

func treeToDFA(tree *synttree.Tree) *DFA {
	info_nodes := make(map[*synttree.Node]*InfoNodes)
	follow_pos := make(map[uint]set)
	pos_to_literal := make(map[uint]string)
	alphabet := make(map[string]struct{})
	var pos uint = 0

	calcInfoNodes(tree.Root, info_nodes, pos_to_literal, &pos)
	calcFollowPos(tree.Root, info_nodes, follow_pos)
	calcAlphabet(tree.Root, alphabet)

	// слайс сетов
	dstates := []set{}
	marked := make(map[uint]bool)
	dtran := make(map[uint]map[string]uint)

	dstates = append(dstates, info_nodes[tree.Root].firstpos)
	marked[0] = false

	for {
		var non_marked_id uint = 0
		found := false
		for id := range dstates {
			if !marked[uint(id)] {
				non_marked_id = uint(id)
				found = true
				break
			}
		}
		if !found {
			break
		}
		marked[non_marked_id] = true
		// множество firstpos
		S := dstates[non_marked_id]
		for literal := range alphabet {
			trans_by_lit := make(set)
			for p := range S {
				if pos_to_literal[p] == literal {
					for fp := range follow_pos[p] {
						trans_by_lit[fp] = struct{}{}
					}
				}
			}
			if len(trans_by_lit) == 0 {
				continue
			}
			equ_id := findState(dstates, trans_by_lit)
			if equ_id == -1 {
				dstates = append(dstates, trans_by_lit)
				// выдаем уникальный айди
				equ_id = len(dstates) - 1
				marked[uint(equ_id)] = false
			}
			if dtran[non_marked_id] == nil {
				dtran[non_marked_id] = make(map[string]uint)
			}
			dtran[non_marked_id][literal] = uint(equ_id)
		}
	}

	// building dead_states
	dead_id := len(dstates)
	for state_id := range dstates {
		if dtran[uint(dead_id)] == nil {
			dtran[uint(dead_id)] = make(map[string]uint)
		}
		for lit := range alphabet {
			if _, ok := dtran[uint(state_id)][lit]; !ok {
				dtran[uint(state_id)][lit] = uint(dead_id)
			}
		}
	}
	for lit := range alphabet {
		dtran[uint(dead_id)][lit] = uint(dead_id)
	}

	// building accept_states
	var hash_pos uint
	for p, lit := range pos_to_literal {
		if lit == "#" {
			hash_pos = p
			break
		}
	}

	accept_states := make(set)
	for id, state := range dstates {
		if _, ok := state[hash_pos]; ok {
			accept_states[uint(id)] = struct{}{}
		}
	}

	return &DFA{
		alphabet,
		dstates,
		dtran,
		0,
		accept_states,
	}
}
