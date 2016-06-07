package relations

import (
	"github.com/oleiade/lane"
)

type treeNode struct {
	data  string
	left  *treeNode
	right *treeNode
}

func parseTree(raw string) *treeNode {
	nodeStack := lane.NewStack()
	operatorStack := lane.NewStack()

	index := 0

	for index < len(raw) {
		char := string(raw[index])
		switch char {
		case "<":
			// TODO: separate pairs and index
			var number string
			number += char
			for char != ">" {
				index += 1
				char = string(raw[index])
				number += char
			}
			node := &treeNode{data: number}
			nodeStack.Push(node)

		case "(":
			operatorStack.Push(char)

		case ")":
			operator := operatorStack.Pop()
			for operator != "(" {
				right := nodeStack.Pop()
				left := nodeStack.Pop()
				node := &treeNode{
					data: operator.(string),
					right: right.(*treeNode),
					left: left.(*treeNode),
				}
				nodeStack.Push(node)
				operator = operatorStack.Pop()
			}

		case "*":
			operand := nodeStack.Pop()
			node := &treeNode{data: char, left: operand.(*treeNode)}
			nodeStack.Push(node)

		case "+", ".":
			operatorStack.Push(char)
		}
		index += 1
	}

	for !operatorStack.Empty() {
		operator := operatorStack.Pop()
		right := nodeStack.Pop()
		left := nodeStack.Pop()
		node := &treeNode{
			data: operator.(string),
			left: left.(*treeNode),
			right: right.(*treeNode),
		}
		nodeStack.Push(node)
	}

	root := nodeStack.Pop()
	return root.(*treeNode)
}
