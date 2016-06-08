package relations

import (
	"github.com/oleiade/lane"
)

type treeNode struct {
	data     string
	left     *treeNode
	right    *treeNode
	index    int // only for leaf nodes
	nullable bool
	firstPos *set
	lastPos  *set
}

func (node *treeNode) updateFollowPos(followPos map[int]*set) map[int]*set {
	switch node.data {
	case ".":
		for position := range (*node.left.lastPos) {
			followPos[position] = node.right.firstPos.union(followPos[position])
		}
	case "*":
		for position := range (*node.lastPos) {
			followPos[position] = node.firstPos.union(followPos[position])
		}
	}
	return followPos
}

func newLeafNode(data string, index int) *treeNode {
	return &treeNode{
		data:     data,
		index:    index,
		firstPos: newSet(index),
		lastPos:  newSet(index),
	}
}

func newOperatorNode(operator string, left, right *treeNode) *treeNode {
	node := &treeNode{data: operator, left: left, right: right}

	switch operator {
	case "*":
		node.nullable = true
		node.firstPos = node.left.firstPos.clone()
		node.lastPos = node.left.lastPos.clone()
	case "+":
		node.nullable = node.right.nullable || node.left.nullable
		node.firstPos = node.left.firstPos.union(node.right.firstPos)
		node.lastPos = node.left.lastPos.union(node.right.lastPos)
	case ".":
		node.nullable = node.right.nullable && node.left.nullable
		if node.left.nullable {
			node.firstPos = node.left.firstPos.union(node.right.firstPos)
		} else {
			node.firstPos = node.left.firstPos.clone()
		}
		if node.right.nullable {
			node.lastPos = node.left.lastPos.union(node.right.lastPos)
		} else {
			node.lastPos = node.right.lastPos.clone()
		}
	}
	return node
}

func parseTree(raw string) (*treeNode, map[int]*set) {
	followPos := make(map[int]*set)

	nodeStack := lane.NewStack()
	operatorStack := lane.NewStack()

	position := 0
	index := 1

	for position < len(raw) {
		char := string(raw[position])
		switch char {
		case "<":
			// TODO: separate pairs and position
			// check if both are epsilon - in that case the node is nullable
			var leaf string
			leaf += char
			for char != ">" {
				position += 1
				char = string(raw[position])
				leaf += char
			}

			nodeStack.Push(newLeafNode(leaf, index))
			index += 1

		case "(", "+", ".":
			operatorStack.Push(char)

		case ")":
			operator := operatorStack.Pop().(string)
			for operator != "(" {
				right := nodeStack.Pop().(*treeNode)
				left := nodeStack.Pop().(*treeNode)
				newNode := newOperatorNode(operator, left, right)
				nodeStack.Push(newNode)
				followPos = newNode.updateFollowPos(followPos)
				operator = operatorStack.Pop().(string)
			}

		case "*":
			operand := nodeStack.Pop().(*treeNode)
			newNode := newOperatorNode(char, operand, nil)
			nodeStack.Push(newNode)
			followPos = newNode.updateFollowPos(followPos)
		}
		position += 1
	}

	// add endmarker character
	operatorStack.Push(".")
	nodeStack.Push(newLeafNode("!", index))

	for !operatorStack.Empty() {
		operator := operatorStack.Pop().(string)
		right := nodeStack.Pop().(*treeNode)
		left := nodeStack.Pop().(*treeNode)
		newNode := newOperatorNode(operator, left, right)
		nodeStack.Push(newNode)
		followPos = newNode.updateFollowPos(followPos)
	}

	root := nodeStack.Pop()
	return root.(*treeNode), followPos
}
