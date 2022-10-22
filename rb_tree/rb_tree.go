// Implementation of Red-Black Trees in a Functional Setting
//
// References:
// - Insertion: https://www.cs.tufts.edu/comp/150FP/archive/chris-okasaki/redblack99.pdf
// - Deletion: https://matt.might.net/articles/red-black-delete/
package immutable_rb_tree

import (
	"github.com/freebirdljj/immutable/comparator"
)

type (
	// The zero value of `RBTree` makes nonsense.
	RBTree[Value any] struct {
		cmp             comparator.Comparator[Value]
		doubleBlackLeaf *node[Value]
		cnt             int
		root            *node[Value]
	}

	node[Value any] struct {
		children [directionNum]*node[Value]
		color    color
		value    Value
	}

	color     int8
	direction int8
)

const (
	colorNegativeBlack color = iota - 1
	colorRed
	colorBlack
	colorDoubleBlack
)

const (
	directionLeft direction = iota
	directionRight
	directionNum
)

func NewRBTree[Value any](cmp comparator.Comparator[Value]) *RBTree[Value] {
	return &RBTree[Value]{
		cmp: cmp,
		doubleBlackLeaf: &node[Value]{
			color: colorDoubleBlack,
		},
	}
}

func NewRBTreeFromValues[Value any](cmp comparator.Comparator[Value], values ...Value) *RBTree[Value] {
	rbTree := NewRBTree[Value](cmp)
	for _, value := range values {
		rbTree, _ = rbTree.Insert(value)
	}
	return rbTree
}

func (rbTree *RBTree[Value]) Count() int {
	return rbTree.cnt
}

func (rbTree *RBTree[Value]) Lookup(value Value) *Value {
	return rbTree.root.lookup(rbTree.cmp, value)
}

func (rbTree *RBTree[Value]) Values() []Value {
	values := []Value(nil)
	rbTree.InorderTraversal(
		func(value Value) {
			values = append(values, value)
		},
	)
	return values
}

func (rbTree *RBTree[Value]) InorderTraversal(visitor func(value Value)) {
	rbTree.root.inorderTraversal(
		func(n *node[Value]) {
			visitor(n.value)
		},
	)
}

// `newTree` returned by `Insert()` is always different from the original one.
// `affected` is true, meaning an actual insertion occurred; otherwise, a replacement occurred.
func (rbTree *RBTree[Value]) Insert(value Value) (newTree *RBTree[Value], affected bool) {

	rbTreeCopy := *rbTree

	newRoot, affected := rbTreeCopy.root.insert(rbTree.cmp, value)
	if affected {
		rbTreeCopy.cnt++
	}

	rbTreeCopy.root = newRoot
	return &rbTreeCopy, affected
}

// `affected` is true, meaning that a real deletion occurred, `newTree` will be different from the original;
// otherwise nothing happens, `newTree` is the original one.
func (rbTree *RBTree[Value]) Delete(value Value) (newTree *RBTree[Value], affected bool) {

	newRoot, affected := rbTree.root.delete(rbTree.cmp, value, rbTree.doubleBlackLeaf)
	if !affected {
		return rbTree, false
	}

	rbTreeCopy := *rbTree
	rbTreeCopy.cnt--
	rbTreeCopy.root = newRoot
	return &rbTreeCopy, true
}

func (n *node[Value]) getColor() color {
	if n == nil {
		return colorBlack
	}
	return n.color
}

func (n *node[Value]) makeRed() *node[Value] {
	nCopy := *n
	nCopy.color = colorRed
	return &nCopy
}

func (n *node[Value]) makeBlack() *node[Value] {

	if n == nil {
		return nil
	}

	nCopy := *n
	nCopy.color = colorBlack
	return &nCopy
}

// Only used by `bubble()`.
func (n *node[Value]) makeRedder(doubleBlackLeaf *node[Value]) *node[Value] {

	if n == doubleBlackLeaf {
		return nil
	}

	nCopy := *n
	nCopy.color = redder(n.color)
	return &nCopy
}

func (n *node[Value]) withChildren(children [directionNum]*node[Value]) *node[Value] {
	nCopy := *n
	nCopy.children = children
	return &nCopy
}

func (n *node[Value]) withEqualValue(value Value) *node[Value] {
	nCopy := *n
	nCopy.value = value
	return &nCopy
}

func (n *node[Value]) lookup(cmp comparator.Comparator[Value], value Value) *Value {

	if n == nil {
		return nil
	}

	switch sign(cmp(value, n.value)) {
	case -1:
		return n.children[directionLeft].lookup(cmp, value)
	case 1:
		return n.children[directionRight].lookup(cmp, value)
	default:
		return &n.value
	}
}

func (n *node[Value]) inorderTraversal(visitor func(n *node[Value])) {

	if n == nil {
		return
	}

	n.children[directionLeft].inorderTraversal(visitor)
	visitor(n)
	n.children[directionRight].inorderTraversal(visitor)
}

// Only `balance` non-leaf node.
func (n *node[Value]) balance() *node[Value] {

	// try to find a red child with a red grandchild.
	if n.color == colorBlack || n.color == colorDoubleBlack {
		color := redder(n.color)
		for _, dir := range []direction{
			directionLeft,
			directionRight,
		} {

			child := n.children[dir]
			if child.getColor() != colorRed {
				continue
			}

			oppositeDir := directionLeft + directionRight - dir

			switch {
			case child.children[dir].getColor() == colorRed:
				grandchildren := [directionNum]*node[Value]{}
				grandchildren[dir] = child.children[oppositeDir]
				grandchildren[oppositeDir] = n.children[oppositeDir]

				newChildren := [directionNum]*node[Value]{}
				newChildren[dir] = child.children[dir].makeBlack()
				newChildren[oppositeDir] = n.withChildren(grandchildren).makeBlack()

				return &node[Value]{
					children: newChildren,
					color:    color,
					value:    child.value,
				}
			case child.children[oppositeDir].getColor() == colorRed:
				newChildren := [directionNum]*node[Value]{}
				{
					grandchildren := [directionNum]*node[Value]{}
					grandchildren[dir] = child.children[dir]
					grandchildren[oppositeDir] = child.children[oppositeDir].children[dir]
					newChildren[dir] = &node[Value]{
						children: grandchildren,
						color:    colorBlack,
						value:    child.value,
					}
				}
				{
					grandchildren := [directionNum]*node[Value]{}
					grandchildren[dir] = child.children[oppositeDir].children[oppositeDir]
					grandchildren[oppositeDir] = n.children[oppositeDir]
					newChildren[oppositeDir] = n.withChildren(grandchildren).makeBlack()
				}
				return &node[Value]{
					children: newChildren,
					color:    color,
					value:    child.children[oppositeDir].value,
				}
			}
		}
	}

	if n.color == colorDoubleBlack {
		for _, dir := range []direction{
			directionLeft,
			directionRight,
		} {

			child := n.children[dir]
			if child.getColor() != colorNegativeBlack {
				continue
			}

			oppositeDir := directionLeft + directionRight - dir

			newChildren := [directionNum]*node[Value]{}
			{
				grandchildren := [directionNum]*node[Value]{}
				grandchildren[dir] = child.children[dir].makeRed()
				grandchildren[oppositeDir] = child.children[oppositeDir].children[dir]
				newChildren[dir] = (&node[Value]{
					children: grandchildren,
					color:    colorBlack,
					value:    child.value,
				}).balance()
			}
			{
				grandchildren := [directionNum]*node[Value]{}
				grandchildren[dir] = child.children[oppositeDir].children[oppositeDir]
				grandchildren[oppositeDir] = n.children[oppositeDir]
				newChildren[oppositeDir] = &node[Value]{
					children: grandchildren,
					color:    colorBlack,
					value:    n.value,
				}
			}

			return &node[Value]{
				children: newChildren,
				color:    colorBlack,
				value:    child.children[oppositeDir].value,
			}
		}
	}

	return n
}

// `newNode` returned by `ins()` is always different from the original one.
// `affected` is true, meaning an actual insertion occurred; otherwise, a replacement occurred.
func (n *node[Value]) ins(cmp comparator.Comparator[Value], value Value) (*node[Value], bool) {

	if n == nil {
		return &node[Value]{
			value: value,
			color: colorRed,
		}, true
	}

	switch sign(cmp(value, n.value)) {
	case -1:
		newLeftChild, affected := n.children[directionLeft].ins(cmp, value)
		return n.withChildren([directionNum]*node[Value]{
			directionLeft:  newLeftChild,
			directionRight: n.children[directionRight],
		}).balance(), affected
	case 1:
		newRightChild, affected := n.children[directionRight].ins(cmp, value)
		return n.withChildren([directionNum]*node[Value]{
			directionLeft:  n.children[directionLeft],
			directionRight: newRightChild,
		}).balance(), affected
	default:
		return n.withEqualValue(value), false
	}
}

// `newNode` returned by `insert()` is always different from the original one.
// `affected` is true, meaning an actual insertion occurred; otherwise, a replacement occurred.
func (n *node[Value]) insert(cmp comparator.Comparator[Value], value Value) (newNode *node[Value], affected bool) {
	result, affected := n.ins(cmp, value)
	return result.makeBlack(), affected
}

// Only `bubble` non-leaf node.
func (n *node[Value]) bubble(doubleBlackLeaf *node[Value]) *node[Value] {
	for _, child := range n.children {
		if child.getColor() == colorDoubleBlack {
			n = &node[Value]{
				children: [directionNum]*node[Value]{
					directionLeft:  n.children[directionLeft].makeRedder(doubleBlackLeaf),
					directionRight: n.children[directionRight].makeRedder(doubleBlackLeaf),
				},
				color: blacker(n.color),
				value: n.value,
			}
			break
		}
	}
	return n.balance()
}

func (n *node[Value]) remove(doubleBlackLeaf *node[Value]) *node[Value] {

	// all children are leaves
	if n.children == [directionNum]*node[Value]{nil, nil} {
		switch n.color {
		case colorRed:
			return nil
		case colorBlack:
			return doubleBlackLeaf
		}
	}

	// exact one child is leaf, node color must be black
	for _, dir := range []direction{
		directionLeft,
		directionRight,
	} {
		if n.children[dir] == nil {
			oppositeDir := directionLeft + directionRight - dir
			nonLeafChild := n.children[oppositeDir]
			return nonLeafChild.makeBlack()
		}
	}

	newLeftChild, maxInLeft := n.children[directionLeft].removeMax(doubleBlackLeaf)
	return (&node[Value]{
		children: [directionNum]*node[Value]{
			directionLeft:  newLeftChild,
			directionRight: n.children[directionRight],
		},
		color: n.color,
		value: maxInLeft,
	}).bubble(doubleBlackLeaf)
}

func (n *node[Value]) removeMax(doubleBlackLeaf *node[Value]) (newNode *node[Value], max Value) {

	rightChild := n.children[directionRight]
	if rightChild == nil {
		return n.remove(doubleBlackLeaf), n.value
	}

	newRightChild, max := rightChild.removeMax(doubleBlackLeaf)
	return n.withChildren([directionNum]*node[Value]{
		directionLeft:  n.children[directionLeft],
		directionRight: newRightChild,
	}).bubble(doubleBlackLeaf), max
}

// `affected` is true, meaning that a real deletion occurred, `newNode` will be different from the original;
// otherwise nothing happens, `newNode` is the original one.
func (n *node[Value]) del(cmp comparator.Comparator[Value], value Value, doubleBlackLeaf *node[Value]) (newNode *node[Value], affected bool) {

	if n == nil {
		return nil, false
	}

	switch sign(cmp(value, n.value)) {
	case -1:

		newLeftChild, affected := n.children[directionLeft].del(cmp, value, doubleBlackLeaf)
		if !affected {
			return n, false
		}

		return n.withChildren([directionNum]*node[Value]{
			directionLeft:  newLeftChild,
			directionRight: n.children[directionRight],
		}).bubble(doubleBlackLeaf), true
	case 1:

		newRightChild, affected := n.children[directionRight].del(cmp, value, doubleBlackLeaf)
		if !affected {
			return n, false
		}

		return n.withChildren([directionNum]*node[Value]{
			directionLeft:  n.children[directionLeft],
			directionRight: newRightChild,
		}).bubble(doubleBlackLeaf), true
	default:
		return n.remove(doubleBlackLeaf), true
	}
}

// `affected` is true, meaning that a real deletion occurred, `newNode` will be different from the original;
// otherwise nothing happens, `newNode` is the original one.
func (n *node[Value]) delete(cmp comparator.Comparator[Value], value Value, doubleBlackLeaf *node[Value]) (newNode *node[Value], affected bool) {

	result, affected := n.del(cmp, value, doubleBlackLeaf)
	if !affected {
		return n, false
	}

	if result == doubleBlackLeaf {
		return nil, true
	}

	return result.makeBlack(), true
}

func blacker(c color) color {
	return map[color]color{
		colorNegativeBlack: colorRed,
		colorRed:           colorBlack,
		colorBlack:         colorDoubleBlack,
	}[c]
}

func redder(c color) color {
	return map[color]color{
		colorRed:         colorNegativeBlack,
		colorBlack:       colorRed,
		colorDoubleBlack: colorBlack,
	}[c]
}

func sign(x int) int {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	default:
		return 0
	}
}
