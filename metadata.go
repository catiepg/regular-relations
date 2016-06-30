package relations

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/oleiade/lane"
)

// Regular expression operators
const (
	union  = '+'
	concat = '.'
	repeat = '*'
	end    = '!'
)

// rule is a basic relation element in a regular expression
type rule struct {
	in  rune
	out string
}

// node is an element in a parse tree
type node interface {
	base() *baseNode
}

// baseNode presents a common structure for a node element
type baseNode struct {
	nullable bool
	first    set
	last     set
}

// ruleNode is a leaf node element for a regular expression rule
type ruleNode struct {
	baseNode
	index int
}

func (rn *ruleNode) base() *baseNode {
	return &rn.baseNode
}

// operatorNode joins node elements via regular expression operation
type operatorNode struct {
	baseNode
	kind  rune
	left  node
	right node
}

func (on *operatorNode) base() *baseNode {
	return &on.baseNode
}

// metadata contains information derived from a regular expression which is
// necessary for the construction of a transducer
type metadata struct {
	rootFirst  set
	follow     map[int]set
	rules      map[int]rule
	finalIndex int
}

func (m *metadata) newRuleNode(in rune, out string) *ruleNode {
	// Unique index for each rule
	m.finalIndex++

	m.rules[m.finalIndex] = rule{in, out}
	node := &ruleNode{
		baseNode{first: newSet(m.finalIndex), last: newSet(m.finalIndex)},
		m.finalIndex,
	}

	// Mark if language of the new node accepts empty string
	node.nullable = (in == 0 && out == "")

	return node
}

// Create new parse tree node and calculate:
//   first  - first rules of the language of the node
//   last   - last rules of the language of the node
//   follow - rules that can follow this node
func (m *metadata) newOperatorNode(operator rune, left, right node) *operatorNode {
	node := &operatorNode{kind: operator, left: left, right: right}

	leftBase := node.left.base()

	// node.right == nil for `*` operator
	var rightBase *baseNode
	if node.right != nil {
		rightBase = node.right.base()
	}

	switch node.kind {
	case repeat:
		node.nullable = true

		node.first = leftBase.first.clone()
		node.last = leftBase.last.clone()

		for position := range node.last {
			m.follow[position] = node.first.union(m.follow[position])
		}
	case union:
		node.nullable = rightBase.nullable || leftBase.nullable

		node.first = leftBase.first.union(rightBase.first)
		node.last = leftBase.last.union(rightBase.last)
	case concat:
		node.nullable = rightBase.nullable && leftBase.nullable

		if leftBase.nullable {
			node.first = leftBase.first.union(rightBase.first)
		} else {
			node.first = leftBase.first.clone()
		}
		if rightBase.nullable {
			node.last = leftBase.last.union(rightBase.last)
		} else {
			node.last = rightBase.last.clone()
		}

		for position := range leftBase.last {
			m.follow[position] = rightBase.first.union(m.follow[position])
		}
	}

	return node
}

// Builds parse tree from regular expression while computing
// nullable, firstPos, lastPos and followPos
func ComputeRegExpMetadata(source io.Reader) (*metadata, error) {
	meta := &metadata{
		follow: map[int]set{},
		rules:  map[int]rule{},
	}

	nodes := lane.NewStack()
	operators := lane.NewStack()

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
			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			for {
				char, _, err = reader.ReadRune()
				if err != nil {
					return nil, errors.New("Incorrect input")
				}

				if char == ',' {
					in = out
					out = &bytes.Buffer{}
					continue
				} else if char == '>' {
					break
				}

				out.WriteRune(char)
			}

			// Add first with output
			first, _, _ := in.ReadRune()
			nodes.Push(meta.newRuleNode(first, out.String()))

			// Add the rest with concatenation
			for {
				c, _, err := in.ReadRune()
				if err == io.EOF {
					break
				}

				right := meta.newRuleNode(c, "")
				left := nodes.Pop().(node)
				nodes.Push(meta.newOperatorNode(concat, left, right))
			}

		case '(', union, concat:
			operators.Push(char)

		case ')':
			for {
				operator := operators.Pop().(rune)
				if operator == '(' {
					break
				}

				right := nodes.Pop().(node)
				left := nodes.Pop().(node)
				nodes.Push(meta.newOperatorNode(operator, left, right))
			}

		case repeat:
			operand := nodes.Pop().(node)
			nodes.Push(meta.newOperatorNode(char, operand, nil))
		}
	}

	for !operators.Empty() {
		operator := operators.Pop().(rune)

		right := nodes.Pop().(node)
		left := nodes.Pop().(node)
		nodes.Push(meta.newOperatorNode(operator, left, right))
	}

	// Add endmarker character
	right := meta.newRuleNode(end, "")
	left := nodes.Pop().(node)
	root := meta.newOperatorNode(concat, left, right)

	meta.rootFirst = root.first

	return meta, nil
}
