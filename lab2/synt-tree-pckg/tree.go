package synttree

import (
	"errors"
	"unicode"
)

type TypeNode int
const (
	LEFT_PAR TypeNode = iota
	RIGHT_PAR
	OP_KLINI
	OP_CONC
	OP_QUESTION
	OP_REPEAT
	OP_OR
	CAPTURE_GROUP
	LITERAL
)

type Node struct {
	type_node TypeNode
	value string
	left *Node
	right *Node
	parent *Node	
}

type Tree struct {
	root *Node
}

func popToken(stack *[]Token) Token {
	if len(*stack) == 0 {
		return Token{}
	}
	token := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return token
}

func peekToken(stack *[]Token) Token {
	if len(*stack) == 0 {
		return Token{}
	}
	return (*stack)[len(*stack)-1]
}

func popNode(stack *[]*Node) *Node {
	if len(*stack) == 0 {
		return nil
	}
	node := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return node
}

type Token struct {
	value string
	type_token TypeNode
}

// работаем с байтами а не рунами, но не критично тк 
// на вход допустимы только цифры и {}
func expectRepeat(str string, i int) (res string, new_i int, err error) {
	for j := i + 1; j < len(str); j++ {
		if str[j] == '}' {
			if len(res) == 0 {
				return "", i, errors.New("expect_repeat: no repeat value was given")
			}
			return res, j, nil
		} else if unicode.IsDigit(rune(str[j])) {
			res += string(str[j])
		} else {
			return "", i, errors.New("expect_repeat: invalid character in repeat value")
		}
	}
	return "", i, errors.New("expect_repeat: no closing '}' found for repeat value")
}

func expectKlini(str string, i int) (res string, new_i int, err error) {
    if i+2 < len(str) && str[i+1] == '.' && str[i+2] == '.' {
        return "...", i + 2, nil // Возвращаем i+2, после i++ в tokenize индекс станет i+3
    }
    return "", i, errors.New("expect_klini: klini is not valid, expected '...'")
}

func expectCaptureGroup(str string, i int) (res string, new_i int, err error) {
	for j := i + 1; j < len(str); j++ {
		if str[j] == '>' {
			if len(res) == 0 {
				return "", i, errors.New("expect_capture_group: no capture group name was given")
			}
			return res, j, nil
		} else {
			res += string(str[j])
		}
	}
	return "", i, errors.New("expect_capture_group: no closing > was given")
}

func tokenize(str string) (result []Token, err error) {
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case '(':
			result = append(result, Token{value: "(", type_token: LEFT_PAR})
		case ')':
			result = append(result, Token{value: ")", type_token: RIGHT_PAR})
		case '?':
			result = append(result, Token{value: "?", type_token: OP_QUESTION})
		case '|':
			result = append(result, Token{value: "|", type_token: OP_OR})
		case '{':
			el, new_i, err := expectRepeat(str, i)
			if err != nil {
				return nil, err
			}
			i = new_i
			result = append(result, Token{value: el, type_token: OP_REPEAT})
		case '.':
			el, new_i, err := expectKlini(str, i)
			if err != nil {
				return nil, err
			}
			i = new_i
			result = append(result, Token{value: el, type_token: OP_KLINI})
		case '<':
			el, new_i, err := expectCaptureGroup(str, i)
			if err != nil {
				return nil, err
			}
			i = new_i
			result = append(result, Token{value: el, type_token: CAPTURE_GROUP})
		default:
			result = append(result, Token{value: string(str[i]), type_token: LITERAL})
		}
	}
	return result, nil
}

func canBeLeftFromConc(token Token) bool {
	switch token.type_token {
	case LITERAL,
	     RIGHT_PAR,
	     OP_KLINI,
	     OP_QUESTION,
	     OP_REPEAT:
		return true
	}
	return false
}

func canBeRightFromConc(token Token) bool {
	switch token.type_token {
		case LITERAL,
			 LEFT_PAR,
			 CAPTURE_GROUP:
			return true
	}
	return false
}

func addConcat(tokens []Token) (result []Token, err error) {
	if len(tokens) == 0 {
		return nil, nil
	}
	for i := 0; i < len(tokens) - 1; i++ {
		result = append(result, tokens[i])
		if canBeLeftFromConc(tokens[i]) && canBeRightFromConc(tokens[i + 1]) {
			result = append(result, Token{value: "+", type_token: OP_CONC})
		}
	}
	result = append(result, tokens[len(tokens)-1])
	return result, nil
}

func precedence(op TypeNode) int {
	switch op {
	case OP_KLINI, OP_REPEAT, OP_QUESTION:
		return 3
	case OP_CONC:
		return 2
	case OP_OR:
		return 1
	default:
		return 0
	}
}

func isOperator(token Token) bool {
	switch token.type_token {
	case OP_KLINI,
	     OP_CONC,
	     OP_QUESTION,
	     OP_REPEAT,
	     OP_OR:
		return true
	}
	return false
}

func makeNode(op Token, node_stack *[]*Node) (*Node, error) {
	node := &Node{
		type_node: op.type_token,
		value: op.value,
	}

	switch op.type_token {

	// 2 children
	case OP_CONC, OP_OR:
		if len(*node_stack) < 2 {
			return nil, errors.New("make_node: not enough operands for binary operator")
		}

		right := popNode(node_stack)
		left := popNode(node_stack)

		node.left = left
		node.right = right

		left.parent = node
		right.parent = node

	// 1 child
	case OP_KLINI, OP_QUESTION, OP_REPEAT:
		if len(*node_stack) < 1 {
			return nil, errors.New("make_node: not enough operands for unary operator")
		}

		child := popNode(node_stack)

		node.left = child
		child.parent = node

	case CAPTURE_GROUP:
		// capture group — это не оператор, он уже обрабатывается при построении дерева, так что сюда он не должен попадать
		// должно быть unreachable
		return nil, errors.New("make_node: capture group doesn't support referencing")

	default:
		return nil, errors.New("make_node: unknown operator type")
	}

	return node, nil
}

func buildTree(tokens []Token) (Tree, error) {
	tokens = append([]Token{Token{"(", LEFT_PAR}}, tokens...)	
	tokens = append(tokens, Token{")", RIGHT_PAR})

	var stack_ops []Token
	var stack_nodes []*Node

	for _, token := range tokens {
		if token.type_token == LITERAL {
			stack_nodes = append(stack_nodes, &Node{type_node: LITERAL, value: token.value})
			continue
		} else if token.type_token == CAPTURE_GROUP {
			stack_ops = append(stack_ops, token)
			continue
		}
		if len(stack_ops) == 0 || token.type_token == LEFT_PAR {
			stack_ops = append(stack_ops, token)
		} else if token.type_token == RIGHT_PAR {
			for {
				if len(stack_ops) == 0 {
					return Tree{}, errors.New("build_tree: mismatched parentheses")
				}
				top := popToken(&stack_ops)
				if top.type_token == LEFT_PAR {
					break
				} else if top.type_token == CAPTURE_GROUP {
					if peekToken(&stack_ops).type_token != LEFT_PAR {
						return Tree{}, errors.New("build_tree: capture group must be immediately after left parenthesis")
					}
					expr := popNode(&stack_nodes)
					if expr == nil {
						return Tree{}, errors.New("build_tree: no expression found for capture group")
					}
					name_node := &Node{type_node: CAPTURE_GROUP, value: top.value, left: expr}
					expr.parent = name_node
					stack_nodes = append(stack_nodes, name_node)
					popToken(&stack_ops)
					break
				}
				node, err := makeNode(top, &stack_nodes)
				if err != nil {
					return Tree{}, err
				}
				stack_nodes = append(stack_nodes, node)
			}
		} else if isOperator(token) {
			for len(stack_ops) > 0 &&
				peekToken(&stack_ops).type_token != LEFT_PAR &&
				precedence(peekToken(&stack_ops).type_token) >= precedence(token.type_token) {
				top := popToken(&stack_ops)
				node, err := makeNode(top, &stack_nodes)
				if err != nil {
					return Tree{}, err
				}
				stack_nodes = append(stack_nodes, node)
			}
			stack_ops = append(stack_ops, token)
		}
	}
	if len(stack_nodes) != 1 {
		return Tree{}, errors.New("build_tree: invalid expression, more than one node left in stack")
	}
	return Tree{root: stack_nodes[0]}, nil
}

func process(input string) (Tree, error) {
	tokens, err := tokenize(input)	
	if err != nil {
		return Tree{}, err
	}
	tokens, err = addConcat(tokens)
	if err != nil {
		return Tree{}, err
	}
	tree, err := buildTree(tokens)
	if err != nil {
		return Tree{}, err
	}
	return tree, nil
}