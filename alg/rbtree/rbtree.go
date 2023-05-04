package rbtree

import (
	"log"
	"strings"
)

const (
	RED   = true
	BLACK = false
)

type rbValue interface {
	uint64 | uint32 | uint16 | int64 | int32 | int16
}

type (
	Tree[T rbValue] struct {
		root *Node[T]
	}
)

func (t *Tree[T]) Clear() {
	t.root = nil
}

func (t *Tree[T]) Root() *Node[T] {
	return t.root
}

func (t *Tree[T]) Count() int {
	if t.root != nil {
		return t.root.size
	}

	return 0
}

func (t *Tree[T]) Rank(rank int) (id uint64, node *Node[T]) {
	return lookupNode(t.root, rank)
}

// --------------------------------------------------------- Lookup by score
func (t *Tree[T]) lookupScore(score T) (rank int, n *Node[T]) {
	n = t.root

	if n == nil {
		return -1, nil
	}

	base := 0
	for n != nil {
		if score == n.score {
			rank = base + _nodesize(n.left) + 1 // start rank
			return rank, n
		} else if score > n.score {
			n = n.left
		} else {
			base += _nodesize(n.left) + len(n.ids)
			n = n.right
		}
	}

	return -1, nil
}

// Locate
// ---------------------------------------------------------- locate a score & id
func (t *Tree[T]) Locate(score T, id uint64) (uint32, *Node[T]) {
	rank, node := t.lookupScore(score)

	if node == nil { // no such score exists
		return 0, nil
	}

	// find the id in all ids
	for k, v := range node.ids {
		if v == id {
			// current rank plus the order in the ids
			return uint32(rank + k), node
		}
	}

	return 0, nil
}

// Insert
// ---------------------------------------------------------- Insert an element
// WRITE-LOCK
func (t *Tree[T]) Insert(score T, id uint64) {
	inserted_node := newNode(score, id, RED, nil, nil)
	if t.root == nil {
		t.root = inserted_node
	} else {
		n := t.root
		for {
			n.size++              // the size of these nodes on the way will be increased by 1
			if score == n.score { // same score, just append the new id in the []ids then return, no structure changes.
				n.ids = append(n.ids, id)
				return
			} else if score > n.score { // find higher score in left subtree
				if n.left == nil {
					n.left = inserted_node
					break
				} else {
					n = n.left
				}
			} else if score < n.score { // find lower score in right subtree
				if n.right == nil {
					n.right = inserted_node
					break
				} else {
					n = n.right
				}
			}
		}
		inserted_node.parent = n
	}

	t.insertCase1(inserted_node)
}

// Delete
// ---------------------------------------------------------- Delete an id from a node
// WRITE-LOCK
func (t *Tree[T]) Delete(id uint64, n *Node[T]) {
	// just delete the given id in []ids if the id is not the only one in this node
	if len(n.ids) > 1 {
		for k, v := range n.ids {
			if v == id {
				n.ids = append(n.ids[:k], n.ids[k+1:]...)
				// decrease size by 1 from this node to the top
				fixupSize(n)
				return
			}
		}
	} else { // the only id in this node, node will be deleted, and the structure will change
		// just decrease size by 1 from N to the root
		fixupSize(n)

		// handle red-black properties, and deletion work.
		if n.left != nil && n.right != nil {
			/* Copy fields from predecessor and then delete it instead */
			pred := maximumNode(n.left)
			// copy score, id
			n.score = pred.score
			n.ids = pred.ids

			// decrease size by pred.size from pred to N
			tmp := pred
			for tmp != n {
				tmp.size -= len(pred.ids)
				tmp = tmp.parent
			}

			// deal with predecessor after.
			n = pred
		}

		var child *Node[T]
		if n.right == nil {
			child = n.left
		} else {
			child = n.right
		}

		if nodeColor(n) == BLACK {
			n.color = nodeColor(child)
			t.deleteCase1(n)
		}

		t.replaceNode(n, child)

		if n.parent == nil && child != nil {
			child.color = BLACK
		}
	}
}

func (t *Tree[T]) rotateLeft(n *Node[T]) {
	r := n.right
	t.replaceNode(n, r)
	n.right = r.left
	if r.left != nil {
		r.left.parent = n
	}
	r.left = n
	n.parent = r

	rotateLeftCallback(n, r)
}

func (t *Tree[T]) rotateRight(n *Node[T]) {
	L := n.left
	t.replaceNode(n, L)
	n.left = L.right
	if L.right != nil {
		L.right.parent = n
	}
	L.right = n
	n.parent = L

	rotateRightCallback(n, L)
}

func (t *Tree[T]) replaceNode(oldn, newn *Node[T]) {
	if oldn.parent == nil {
		t.root = newn
	} else {
		if oldn == oldn.parent.left {
			oldn.parent.left = newn
		} else {
			oldn.parent.right = newn
		}
	}
	if newn != nil {
		newn.parent = oldn.parent
	}
}

func (t *Tree[T]) insertCase1(n *Node[T]) {
	if n.parent == nil {
		n.color = BLACK
	} else {
		t.insertCase2(n)
	}
}

func (t *Tree[T]) insertCase2(n *Node[T]) {
	if nodeColor(n.parent) == BLACK {
		return /* Tree is still valid */
	} else {
		t.insertCase3(n)
	}
}

func (t *Tree[T]) insertCase3(n *Node[T]) {
	if nodeColor(uncle(n)) == RED {
		n.parent.color = BLACK
		uncle(n).color = BLACK
		grandparent(n).color = RED
		t.insertCase1(grandparent(n))
	} else {
		t.insertCase4(n)
	}
}

func (t *Tree[T]) insertCase4(n *Node[T]) {
	if n == n.parent.right && n.parent == grandparent(n).left {
		t.rotateLeft(n.parent)
		n = n.left
	} else if n == n.parent.left && n.parent == grandparent(n).right {
		t.rotateRight(n.parent)
		n = n.right
	}
	t.insertCase5(n)
}

func (t *Tree[T]) insertCase5(n *Node[T]) {
	n.parent.color = BLACK
	grandparent(n).color = RED
	if n == n.parent.left && n.parent == grandparent(n).left {
		t.rotateRight(grandparent(n))
	} else {
		t.rotateLeft(grandparent(n))
	}
}

func (t *Tree[T]) deleteCase1(n *Node[T]) {
	if n.parent == nil {
		return
	} else {
		t.deleteCase2(n)
	}
}

func (t *Tree[T]) deleteCase2(n *Node[T]) {
	if nodeColor(sibling(n)) == RED {
		n.parent.color = RED
		sibling(n).color = BLACK
		if n == n.parent.left {
			t.rotateLeft(n.parent)
		} else {
			t.rotateRight(n.parent)
		}
	}
	t.deleteCase3(n)
}

func (t *Tree[T]) deleteCase3(n *Node[T]) {
	if nodeColor(n.parent) == BLACK &&
		nodeColor(sibling(n)) == BLACK &&
		nodeColor(sibling(n).left) == BLACK &&
		nodeColor(sibling(n).right) == BLACK {
		sibling(n).color = RED
		t.deleteCase1(n.parent)
	} else {
		t.deleteCase4(n)
	}
}

func (t *Tree[T]) deleteCase4(n *Node[T]) {
	if nodeColor(n.parent) == RED &&
		nodeColor(sibling(n)) == BLACK &&
		nodeColor(sibling(n).left) == BLACK &&
		nodeColor(sibling(n).right) == BLACK {
		sibling(n).color = RED
		n.parent.color = BLACK
	} else {
		t.deleteCase5(n)
	}
}

func (t *Tree[T]) deleteCase5(n *Node[T]) {
	if n == n.parent.left &&
		nodeColor(sibling(n)) == BLACK &&
		nodeColor(sibling(n).left) == RED &&
		nodeColor(sibling(n).right) == BLACK {
		sibling(n).color = RED
		sibling(n).left.color = BLACK
		t.rotateRight(sibling(n))
	} else if n == n.parent.right &&
		nodeColor(sibling(n)) == BLACK &&
		nodeColor(sibling(n).right) == RED &&
		nodeColor(sibling(n).left) == BLACK {
		sibling(n).color = RED
		sibling(n).right.color = BLACK
		t.rotateLeft(sibling(n))
	}
	t.deleteCase6(n)
}

func (t *Tree[T]) deleteCase6(n *Node[T]) {
	sibling(n).color = nodeColor(n.parent)
	n.parent.color = BLACK
	if n == n.parent.left {
		sibling(n).right.color = BLACK
		t.rotateLeft(n.parent)
	} else {
		sibling(n).left.color = BLACK
		t.rotateRight(n.parent)
	}
}

// ---------------------------------------------------------- tree print
const INDENT_STEP = 4

func Print_helper[T rbValue](n *Node[T], indent int) {
	if n == nil {
		log.Printf("<empty tree>")
		return
	}
	if n.right != nil {
		Print_helper(n.right, indent+INDENT_STEP)
	}
	if n.color == BLACK {
		log.Printf(strings.Repeat(" ", indent)+"[score:%v size:%v id:%v len:%v]\n", n.score, n.size, n.ids, len(n.ids))
	} else {
		log.Printf(strings.Repeat(" ", indent)+"*[score:%v size:%v id:%v len:%v]\n", n.score, n.size, n.ids, len(n.ids))
	}

	if n.left != nil {
		Print_helper(n.left, indent+INDENT_STEP)
	}
}
