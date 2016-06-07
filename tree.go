package relations

import (
	"github.com/deckarep/golang-set"
	"github.com/oleiade/lane"
)

type treeNode struct {
	data     string
	left     *treeNode
	right    *treeNode
	index    int // only for leaf nodes
	nullable bool
	firstPos mapset.Set
	lastPos  mapset.Set
}

func (node *treeNode) updateFollowPos(followPos map[int]mapset.Set) map[int]mapset.Set {
	switch node.data {
	case ".":
		for pos := range node.left.lastPos.Iter() {
			position := pos.(int)
			if followSet, ok := followPos[position]; ok {
				followPos[position] = followSet.Union(node.right.firstPos)
			} else {
				followPos[position] = node.right.firstPos.Clone()
			}
		}
	case "*":
		for pos := range node.lastPos.Iter() {
			position := pos.(int)
			if followSet, ok := followPos[position]; ok {
				followPos[position] = followSet.Union(node.firstPos)
			} else {
				followPos[position] = node.firstPos.Clone()
			}
		}
	}
	return followPos
}

func newLeafNode(data string, index int) *treeNode {
	return &treeNode{
		data:     data,
		index:    index,
		firstPos: mapset.NewSetWith(index),
		lastPos:  mapset.NewSetWith(index),
	}
}

func newOperatorNode(operator string, left, right *treeNode) *treeNode {
	node := &treeNode{data: operator, left: left, right: right}

	switch operator {
	case "*":
		node.nullable = true
		node.firstPos = node.left.firstPos.Clone()
		node.lastPos = node.left.lastPos.Clone()
	case "+":
		node.nullable = node.right.nullable || node.left.nullable
		node.firstPos = node.left.firstPos.Union(node.right.firstPos)
		node.lastPos = node.left.lastPos.Union(node.right.lastPos)
	case ".":
		node.nullable = node.right.nullable && node.left.nullable
		if node.left.nullable {
			node.firstPos = node.left.firstPos.Union(node.right.firstPos)
		} else {
			node.firstPos = node.left.firstPos.Clone()
		}
		if node.right.nullable {
			node.lastPos = node.left.lastPos.Union(node.right.lastPos)
		} else {
			node.lastPos = node.right.lastPos.Clone()
		}
	}
	return node
}

func parseTree(raw string) (*treeNode, map[int]mapset.Set) {
	followPos := make(map[int]mapset.Set)

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
