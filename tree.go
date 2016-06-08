package relations

import (
	"github.com/oleiade/lane"
)

type node struct {
	data     string
	left     *node
	right    *node
	index    int
	nullable bool
	first *set
	last  *set
}

// Computes nullable, firstPos and lastPos
func (n *node) annotate() {
	switch n.data {
	case "*":
		n.nullable = true
		n.first = n.left.first.clone()
		n.last = n.left.last.clone()
	case "+":
		n.nullable = n.right.nullable || n.left.nullable
		n.first = n.left.first.union(n.right.first)
		n.last = n.left.last.union(n.right.last)
	case ".":
		n.nullable = n.right.nullable && n.left.nullable
		if n.left.nullable {
			n.first = n.left.first.union(n.right.first)
		} else {
			n.first = n.left.first.clone()
		}
		if n.right.nullable {
			n.last = n.left.last.union(n.right.last)
		} else {
			n.last = n.right.last.clone()
		}
	}
}

type tree struct {
	root *node
	follow map[int]*set
	symbols map[int]string
}

func newTree() *tree {
	return &tree{follow: make(map[int]*set), symbols: make(map[int]string)}
}

func (t *tree) newLeafNode(data string, index int) *node {
	t.symbols[index] = data
	return &node{data: data, index: index,
		first: newSet(index), last: newSet(index)}
}

func (t *tree) newOperatorNode(operator string, left, right *node) *node {
	newNode := &node{data: operator, left: left, right: right}
	newNode.annotate()
	t.updateFollow(newNode)
	return newNode
}

// Updates followPos with information from newly created node
func (t *tree) updateFollow(n *node) {
	switch n.data {
	case ".":
		for position := range (*n.left.last) {
			t.follow[position] = n.right.first.union(t.follow[position])
		}
	case "*":
		for position := range (*n.last) {
			t.follow[position] = n.first.union(t.follow[position])
		}
	}
}

// Builds parse tree from regular expression while computing
// nullable, firstPos, lastPos and followPos
func buildTree(raw string) *tree {
	t := newTree()

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

			nodeStack.Push(t.newLeafNode(leaf, index))
			index += 1

		case "(", "+", ".":
			operatorStack.Push(char)

		case ")":
			operator := operatorStack.Pop().(string)
			for operator != "(" {
				right := nodeStack.Pop().(*node)
				left := nodeStack.Pop().(*node)
				nodeStack.Push(t.newOperatorNode(operator, left, right))
				operator = operatorStack.Pop().(string)
			}

		case "*":
			operand := nodeStack.Pop().(*node)
			nodeStack.Push(t.newOperatorNode(char, operand, nil))
		}
		position += 1
	}

	// Add endmarker character
	operatorStack.Push(".")
	nodeStack.Push(t.newLeafNode("!", index))

	for !operatorStack.Empty() {
		operator := operatorStack.Pop().(string)
		right := nodeStack.Pop().(*node)
		left := nodeStack.Pop().(*node)
		nodeStack.Push(t.newOperatorNode(operator, left, right))
	}

	t.root = nodeStack.Pop().(*node)
	return t
}
