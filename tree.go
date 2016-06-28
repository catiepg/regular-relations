package relations

import (
	"bufio"
	"io"

	"github.com/oleiade/lane"
)

type ParseError struct {
	msg string
}

func (err ParseError) Error() string {
	return err.msg
}

type element struct {
	in  rune
	out string
}

func (p *element) contain(in rune, out string) bool {
	return p.in == in && p.out == out
}

type node struct {
	operator rune
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
	case '*':
		n.nullable = true
		n.first = n.left.first.clone()
		n.last = n.left.last.clone()
	case '+':
		n.nullable = n.right.nullable || n.left.nullable
		n.first = n.left.first.union(n.right.first)
		n.last = n.left.last.union(n.right.last)
	case '.':
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

type parseTree struct {
	alphabet  []*element
	rootFirst *set
	follow    map[int]*set
	elements  map[int]*element
	final     int
}

func newParseTree() *parseTree {
	return &parseTree{follow: make(map[int]*set), elements: make(map[int]*element)}
}

func (t *parseTree) newLeafNode(in rune, out string, index int) *node {
	t.elements[index] = t.updateAlphabet(in, out)
	newNode := &node{index: index, first: newSet(index), last: newSet(index)}
	// TODO: if in == "" && out == "" {
	if in == 0 && out == "" {
		newNode.nullable = true
	}
	return newNode
}

func (t *parseTree) newOperatorNode(operator rune, left, right *node) *node {
	newNode := &node{operator: operator, left: left, right: right}
	newNode.annotate()
	t.updateFollow(newNode)
	return newNode
}

func (t *parseTree) updateAlphabet(in rune, out string) *element {
	for _, p := range t.alphabet {
		if p.contain(in, out) {
			return p
		}
	}
	newPair := &element{in: in, out: out}
	// TODO: ???
	if in != '!' {
		t.alphabet = append(t.alphabet, newPair)
	}
	return newPair
}

// Updates followPos with information from newly created node
func (t *parseTree) updateFollow(n *node) {
	switch n.operator {
	case '*':
		for position := range *n.last {
			t.follow[position] = n.first.union(t.follow[position])
		}
	case '.':
		for position := range *n.left.last {
			t.follow[position] = n.right.first.union(t.follow[position])
		}
	}
}

// Builds parse tree from regular expression while computing
// nullable, firstPos, lastPos and followPos
func NewTree(source io.Reader) (*parseTree, error) {
	t := newParseTree()

	nodeStack := lane.NewStack()
	operatorStack := lane.NewStack()
	nodeIndex := 1

	reader := bufio.NewReader(source)

	for {
		char, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch char {
		case '<':
			var inputCharacters, characters []rune
			for {
				char, _, err = reader.ReadRune()
				if err == io.EOF {
					return nil, ParseError{"Incorrect input"}
				} else if err != nil {
					return nil, err
				}

				if char == ',' {
					inputCharacters = characters
					characters = []rune{}
					continue
				} else if char == '>' {
					break
				}

				characters = append(characters, char)
			}

			// Add first with output
			nodeStack.Push(t.newLeafNode(inputCharacters[0],
				string(characters), nodeIndex))
			inputCharacters = inputCharacters[1:]
			nodeIndex += 1

			// Add the rest with concatenation
			for _, c := range inputCharacters {
				right := t.newLeafNode(c, "", nodeIndex)
				left := nodeStack.Pop().(*node)
				nodeStack.Push(t.newOperatorNode('.', left, right))
				nodeIndex += 1
			}

		case '(', '+', '.':
			operatorStack.Push(char)

		case ')':
			operator := operatorStack.Pop().(rune)
			for operator != '(' {
				right := nodeStack.Pop().(*node)
				left := nodeStack.Pop().(*node)
				nodeStack.Push(t.newOperatorNode(operator, left, right))
				operator = operatorStack.Pop().(rune)
			}

		case '*':
			operand := nodeStack.Pop().(*node)
			nodeStack.Push(t.newOperatorNode(char, operand, nil))
		}
	}

	for !operatorStack.Empty() {
		operator := operatorStack.Pop().(rune)
		right := nodeStack.Pop().(*node)
		left := nodeStack.Pop().(*node)
		nodeStack.Push(t.newOperatorNode(operator, left, right))
	}

	// Add endmarker character
	right := t.newLeafNode('!', "", nodeIndex)
	left := nodeStack.Pop().(*node)
	root := t.newOperatorNode('.', left, right)

	t.final = nodeIndex

	t.rootFirst = root.first
	return t, nil
}
