package match

type node struct {
	Left    *node
	right   *node
	content *SinglePerson
}

// time: O(logN)
// space: O(M) M = deep of tree
func insertNode(root *node, newNode *node) *node {
	if root == nil {
		return newNode
	}

	if newNode.content.Gender < root.content.Gender || (newNode.content.Gender == root.content.Gender && newNode.content.Height < root.content.Height) {
		root.Left = insertNode(root.Left, newNode)
	} else {
		root.right = insertNode(root.right, newNode)
	}

	return root
}

// time: O(logN)
// space: O(M) M = deep of tree
func deleteNode(root *node, nodeToDelete *node) *node {
	if root == nil {
		return root
	}

	if root.content.ID == nodeToDelete.content.ID {
		if root.Left == nil {
			return root.right
		} else if root.right == nil {
			return root.Left
		}

		successor := findSuccessor(root.right)
		root = successor
	} else if nodeToDelete.content.Gender < root.content.Gender || (nodeToDelete.content.Gender == root.content.Gender && nodeToDelete.content.Height < root.content.Height) {
		root.Left = deleteNode(root.Left, nodeToDelete)
	} else {
		root.right = deleteNode(root.right, nodeToDelete)
	}

	return root
}

func findSuccessor(root *node) *node {
	for root.Left != nil {
		root = root.Left
	}
	return root
}

// time: O(M*logN) M = limit
// space: O(D) D = deep of tree
func findMatchesWithLimitAndCondition(root *node, match *node, limit int) ([]*node, int) {
	var res []*node

	if root == nil || limit <= 0 {
		return res, 0
	}

	if match.content.Gender == female {
		if root.content.Gender > match.content.Gender && root.content.Height > match.content.Height {
			// exclude that has been matched person
			if _, ok := match.content.MatchPerson[root.content.ID]; !ok {
				if _, ok := match.content.NewMatchPerson[root.content.ID]; !ok {
					res = append(res, root)
					limit--
				}
			}
			rightNodes, remaining := findMatchesWithLimitAndCondition(root.right, match, limit)
			limit = remaining
			res = append(res, rightNodes...)
			leftNodes, remaining := findMatchesWithLimitAndCondition(root.Left, match, limit)
			res = append(res, leftNodes...)
			limit = remaining
		} else {
			rightNodes, remaining := findMatchesWithLimitAndCondition(root.right, match, limit)
			limit = remaining
			res = append(res, rightNodes...)
		}
	} else if match.content.Gender == male {
		if root.content.Gender < match.content.Gender && root.content.Height < match.content.Height {
			// exclude that has been matched person
			if _, ok := match.content.MatchPerson[root.content.ID]; !ok {
				if _, ok := match.content.NewMatchPerson[root.content.ID]; !ok {
					res = append(res, root)
					limit--
				}
			}
			rightNodes, remaining := findMatchesWithLimitAndCondition(root.right, match, limit)
			limit = remaining
			res = append(res, rightNodes...)
			leftNodes, remaining := findMatchesWithLimitAndCondition(root.Left, match, limit)
			res = append(res, leftNodes...)
			limit = remaining
		} else {
			leftNodes, remaining := findMatchesWithLimitAndCondition(root.Left, match, limit)
			res = append(res, leftNodes...)
			limit = remaining
		}
	}

	return res, limit
}
