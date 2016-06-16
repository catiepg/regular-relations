package relations

import (
	"github.com/oleiade/lane"
)

type pair struct {
	in  string
	out string
}

func (p *pair) contain(in, out string) bool {
	return p.in == in && p.out == out
}

func (p *pair) equal(o *pair) bool {
	return p.contain(o.in, o.out)
}

type node struct {
	operator string
	left     *node
	right    *node
	index    int
	nullable bool
	first    *set
	last     *set
}

// Computes nullable, firstPos and lastPos
func (n *node) annotate() {
	switch n.operator {
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
	alphabet  []*pair
	rootFirst *set
	follow    map[int]*set
	symbols   map[int]*pair
	final     int
}

func newTree() *tree {
	return &tree{follow: make(map[int]*set), symbols: make(map[int]*pair)}
}

func (t *tree) updateAlphabet(in, out string) *pair {
	for _, p := range t.alphabet {
		if p.contain(in, out) {
			return p
		}
	}
	newPair := &pair{in: in, out: out}
	// TODO: ???
	if in != "!" {
		t.alphabet = append(t.alphabet, newPair)
	}
	return newPair
}

func (t *tree) newLeafNode(in, out string, index int) *node {
	t.symbols[index] = t.updateAlphabet(in, out)
	newNode := &node{index: index, first: newSet(index), last: newSet(index)}
	if in == "" && out == "" {
		newNode.nullable = true
	}
	return newNode
}

func (t *tree) newOperatorNode(operator string, left, right *node) *node {
	newNode := &node{operator: operator, left: left, right: right}
	newNode.annotate()
	t.updateFollow(newNode)
	return newNode
}

// Updates followPos with information from newly created node
func (t *tree) updateFollow(n *node) {
	switch n.operator {
	case "*":
		for position := range *n.last {
			t.follow[position] = n.first.union(t.follow[position])
		}
	case ".":
		for position := range *n.left.last {
			t.follow[position] = n.right.first.union(t.follow[position])
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
			var s, in, out string
			for {
				position += 1
				char = string(raw[position])
				if char == "," {
					in = s
					s = ""
					continue
				}
				if char == ">" {
					out = s
					break
				}
				s += char
			}

			nodeStack.Push(t.newLeafNode(in, out, index))
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

	for !operatorStack.Empty() {
		operator := operatorStack.Pop().(string)
		right := nodeStack.Pop().(*node)
		left := nodeStack.Pop().(*node)
		nodeStack.Push(t.newOperatorNode(operator, left, right))
	}

	// Add endmarker character
	right := t.newLeafNode("!", "", index)
	left := nodeStack.Pop().(*node)
	root := t.newOperatorNode(".", left, right)

	t.final = index

	t.rootFirst = root.first
	return t
}
