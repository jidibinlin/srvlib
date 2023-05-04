package rbtree

type Node[T rbValue] struct {
	left   *Node[T]
	right  *Node[T]
	parent *Node[T]

	size  int // the size of this subtree
	color bool

	score T        // the score
	ids   []uint64 // associated ids
}

func newNode[T rbValue](score T, id uint64, color bool, left, right *Node[T]) *Node[T] {
	n := Node[T]{score: score, color: color, left: left, right: right, size: 1, ids: []uint64{id}}
	return &n
}

func (n *Node[T]) Ids() []uint64 {
	return n.ids
}

func (n *Node[T]) Score() T {
	return n.score
}

func _nodesize[T rbValue](n *Node[T]) int {
	if n == nil {
		return 0
	}

	return n.size
}

func lookupNode[T rbValue](n *Node[T], rank int) (id uint64, node *Node[T]) {
	if n == nil {
		return 0, nil // beware of nil pointer
	}

	start := _nodesize(n.left) + 1
	end := _nodesize(n.left) + len(n.ids)

	if rank >= start && rank <= end {
		return n.ids[rank-start], n
	}

	if rank < start {
		return lookupNode(n.left, rank)
	}
	return lookupNode(n.right, rank-end)
}

/**
* left/right rotation call back function
 */
func rotateLeftCallback[T rbValue](n, parent *Node[T]) {
	parent.size = _nodesize(n)
	n.size = _nodesize(n.left) + _nodesize(n.right) + len(n.ids)
}

func rotateRightCallback[T rbValue](n, parent *Node[T]) {
	rotateLeftCallback(n, parent)
}

func fixupSize[T rbValue](n *Node[T]) {
	for n != nil {
		n.size--
		n = n.parent
	}
}

// --------------------------------------------------------- Tree part
func grandparent[T rbValue](n *Node[T]) *Node[T] {
	return n.parent.parent
}

func sibling[T rbValue](n *Node[T]) *Node[T] {
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func uncle[T rbValue](n *Node[T]) *Node[T] {
	return sibling(n.parent)
}

func nodeColor[T rbValue](n *Node[T]) bool {
	if n == nil {
		return BLACK
	}
	return n.color
}

func maximumNode[T rbValue](n *Node[T]) *Node[T] {
	for n.right != nil {
		n = n.right
	}
	return n
}
