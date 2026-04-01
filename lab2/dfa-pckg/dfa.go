package dfapckg

import (
	"fmt"
	synttree "reglib/synttree-pckg"
)

type set map[uint]struct{}

type DFA struct {
	s set
}

type InfoNodes struct {
	nullable bool
	firstpos set
	lastpos  set
}

// follow-pos will calc later

func calcInfoNodes(node *synttree.Node, info_nodes map[*synttree.Node]*InfoNodes, pos *uint) {
	if node == nil {
		return
	}
	calcInfoNodes(node.Left, info_nodes, pos)
	calcInfoNodes(node.Right, info_nodes, pos)

	switch node.Type_node {
	case synttree.LITERAL, synttree.SHEBANG:
		info_nodes[node] = &InfoNodes{
			nullable: false,
			firstpos: map[uint]struct{}{*pos: {}},
			lastpos:  map[uint]struct{}{*pos: {}},
		}
		*pos++
	case synttree.OP_OR:
		info_nodes[node] = &InfoNodes{
			nullable: info_nodes[node.Left].nullable || info_nodes[node.Right].nullable,
			firstpos: mergeSets(info_nodes[node.Left].firstpos, info_nodes[node.Right].firstpos),
			lastpos:  mergeSets(info_nodes[node.Left].lastpos, info_nodes[node.Right].lastpos),
		}
	case synttree.OP_CONC:
		if info_nodes[node.Left].nullable {
			info_nodes[node] = &InfoNodes{
				nullable: info_nodes[node.Left].nullable && info_nodes[node.Right].nullable,
				firstpos: mergeSets(info_nodes[node.Left].firstpos, info_nodes[node.Right].firstpos),
				lastpos:  mergeSets(info_nodes[node.Left].lastpos, info_nodes[node.Right].lastpos),
			}
		} else {
			info_nodes[node] = &InfoNodes{
				nullable: info_nodes[node.Left].nullable && info_nodes[node.Right].nullable,
				firstpos: info_nodes[node.Left].firstpos,
				lastpos:  info_nodes[node.Right].lastpos,
			}
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
		panic(fmt.Sprintf("calcInfoNodes: operand %v is not supported", node.type_node))
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

func treeToDFA(tree *synttree.Tree) *DFA {
	info_nodes := make(map[*synttree.Node]*InfoNodes)
	follow_pos := make(map[uint]set)
	alphabet := make(map[string]struct{})
	var pos uint = 0

	calcInfoNodes(tree.Root, info_nodes, &pos)
	calcFollowPos(tree.Root, info_nodes, follow_pos)
	calcAlphabet(tree.Root, alphabet)

	dstates := []set{}
	marked := make(map[uint]bool)
	dtran := make(map[uint]map[string]uint)

	dstates = append(dstates, info_nodes[tree.Root].firstpos)
	marked[0] = true

	for {
		all_marked := true
		var non_marked_id uint = 0
		for id, mark := range marked {
			if !mark {
				all_marked = false
				non_marked_id = id
			}
		}
		if all_marked {
			break
		}
		marked[non_marked_id] = true
		for
	}

	return DFA{}
}
