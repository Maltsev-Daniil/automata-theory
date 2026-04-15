package dfapckg

func Isomorphic(dfa1, dfa2 *DFA) bool {
	alphabet1 := extractSortedAlphabet(dfa1.alphabet)
	alphabet2 := extractSortedAlphabet(dfa2.alphabet)

	if len(alphabet1) != len(alphabet2) {
		return false
	}
	for i := 0; i < len(alphabet1); i++ {
		if alphabet1[i] != alphabet2[i] {
			return false
		}
	}

	mapping := make(map[int]int)

	stack1 := make([]int, 0)
	stack2 := make([]int, 0)

	stack1 = append(stack1, dfa1.start_state)
	stack2 = append(stack2, dfa2.start_state)
	// достаточно только соответсвие вершинам проверить
	// тк все остальные варианты мы отсекли на проверке
	// равенства алфавита!!!
	for len(stack1) > 0 {
		v1 := stack1[len(stack1)-1]
		v2 := stack2[len(stack2)-1]

		stack1 = stack1[:len(stack1)-1]
		stack2 = stack2[:len(stack2)-1]

		if is_v2, ok := mapping[v1]; ok {
			if is_v2 != v2 {
				return false
			}
			continue
		}
		mapping[v1] = v2

		_, acc1 := dfa1.accept_states[v1]
		_, acc2 := dfa2.accept_states[v2]

		if acc1 != acc2 {
			return false
		}

		for _, lit := range alphabet1 {
			// тк граф у нас полный то у нас есть переходы
			// для каждого символа
			stack1 = append(stack1, dfa1.dtran[v1][lit])
			stack2 = append(stack2, dfa2.dtran[v2][lit])
		}
	}
	return true
}
