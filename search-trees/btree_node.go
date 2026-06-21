package searchtrees

import (
	"cmp"
)

// item holds a single key-value pair stored inside a B-tree node.
type item[k cmp.Ordered, v any] struct {
	key k
	val v
}

// node is an internal or leaf node of a B-tree.
//
// items and children are pre-allocated slices of fixed length (maxItems and
// maxChildren respectively). numItems and numChildrens are explicit occupancy
// counters; only indices [0..numItems-1] and [0..numChildrens-1] hold live data.
//
// Invariants (maintained by Insert / Delete, not enforced here):
//   - Leaf node:     numChildrens == 0
//   - Internal node: numChildrens == numItems + 1
//   - items[0..numItems-1] are kept in ascending key order.
type node[k cmp.Ordered, v any] struct {
	items        []*item[k, v]
	children     []*node[k, v]
	maxItems     int
	minItems     int
	maxChildren  int
	numItems     int
	numChildrens int
}

// newNode allocates a node with pre-allocated slices of length maxItems (items)
// and maxChildren (children). minItems is stored on the node so that Insert /
// Delete can enforce the minimum-occupancy invariant without a back-reference
// to BTree. numItems and numChildrens both start at 0.
func newNode[k cmp.Ordered, v any](maxItems, minItems, maxChildren int) *node[k, v] {
	return &node[k, v]{
		items:       make([]*item[k, v], maxItems),
		children:    make([]*node[k, v], maxChildren),
		minItems:    minItems,
		maxItems:    maxItems,
		maxChildren: maxChildren,
	}
}

// isLeaf reports whether n is a leaf node (has no child pointers).
func (n *node[k, v]) isLeaf() bool {
	return n.numChildrens == 0
}

// insertItemAt inserts i at position pos in the sorted items array, shifting
// items[pos..numItems-1] one slot to the right to make room.
//
// Precondition: numItems < maxItems (node is not full).
//
// Dry run — items=[{1},{5}], numItems=2, inserting {3} at pos=1:
//
//	Before: [&{1}, &{5}, nil]
//	copy items[2:3] ← items[1:2]  →  [&{1}, &{5}, &{5}]
//	items[1] = &{3}               →  [&{1}, &{3}, &{5}]
//	numItems = 3
func (n *node[k, v]) insertItemAt(pos int, i *item[k, v]) {
	if pos < n.numItems {
		// Shift items[pos..numItems-1] right by one to open slot at pos.
		copy(n.items[pos+1:n.numItems+1], n.items[pos:n.numItems])
	}
	n.items[pos] = i
	n.numItems++
}

// insertChildAt inserts child pointer c at position pos, shifting
// children[pos..numChildrens-1] one slot to the right.
//
// Precondition: numChildrens < maxChildren.
//
// Dry run — children=[c0, c1, nil, nil], numChildrens=2, inserting cNew at pos=1:
//
//	Before: [c0, c1, nil, nil]
//	copy children[2:3] ← children[1:2]  →  [c0, c1, c1, nil]
//	children[1] = cNew                  →  [c0, cNew, c1, nil]
//	numChildrens = 3
func (n *node[k, v]) insertChildAt(pos int, c *node[k, v]) {
	if pos < n.numChildrens {
		// Shift children[pos..numChildrens-1] right by one.
		copy(n.children[pos+1:n.numChildrens+1], n.children[pos:n.numChildrens])
	}
	n.children[pos] = c
	n.numChildrens++
}

// searchNode locates key within the node's sorted items using binary search.
//
// Return values:
//   - (index, true)  — key found at n.items[index].
//   - (index, false) — key absent; n.children[index] is the subtree that
//     could contain key (valid only when n is an internal node).
//
// # Binary search dry run
//
// Assume n.items = [{1,"a"}, {3,"b"}, {5,"c"}, {7,"d"}]
//
// Example 1 — key = 4 (not present):
//
//	low=0, high=4
//	  iter 1: mid=2, items[2].key=5 > 4  → high=2
//	  iter 2: mid=1, items[1].key=3 < 4  → low=2
//	  low==high=2 → return (2, false)     ← descend into children[2]
//
// Example 2 — key = 3 (present):
//
//	low=0, high=4
//	  iter 1: mid=2, items[2].key=5 > 3  → high=2
//	  iter 2: mid=1, items[1].key=3 == 3 → return (1, true)
//
// Example 3 — key = 0 (smaller than all):
//
//	low=0, high=4
//	  iter 1: mid=2, items[2].key=5 > 0  → high=2
//	  iter 2: mid=1, items[1].key=3 > 0  → high=1
//	  iter 3: mid=0, items[0].key=1 > 0  → high=0
//	  low==high=0 → return (0, false)     ← descend into children[0]
func (n *node[k, v]) searchNode(key k) (index int, found bool) {
	low, high := 0, n.numItems
	for low < high {
		mid := (low + high) / 2
		switch {
		case key > n.items[mid].key:
			low = mid + 1
		case key < n.items[mid].key:
			high = mid
		default: // key == n.items[mid].key
			return mid, true
		}
	}
	return low, false
}

// split divides a full node at its middle key (index n.minItems) into two
// half-full nodes and returns the median key to be promoted to the parent.
//
// Precondition: n.numItems == n.maxItems (node is full).
//
// After the call:
//   - The receiver (left) retains items[0..minItems-1] and (if internal)
//     children[0..minItems].
//   - The returned right node holds items[minItems+1..maxItems-1] and (if
//     internal) children[minItems+1..maxItems].
//   - The returned midItem is the promoted separator; it leaves both nodes.
//
// # Dry run — degree=2 (maxItems=3, minItems=1)
//
// Before (leaf):
//
//	left.items = [&1, &2, &3],  left.numItems = 3
//
//	  mid     = 1
//	  midItem = items[1] = &2
//
//	After:
//	  right.items = [&3],  right.numItems = 1
//	  left.items  = [&1],  left.numItems  = 1    (items[1..2] nil-ed out)
//	  promoted    = &2
//
// Before (internal, 4 children):
//
//	left.items    = [&1, &2, &3],          left.numItems     = 3
//	left.children = [c0, c1, c2, c3],     left.numChildrens = 4
//
//	After:
//	  right.items    = [&3],      right.numItems     = 1
//	  right.children = [c2, c3],  right.numChildrens = 2
//	  left.items     = [&1],      left.numItems      = 1
//	  left.children  = [c0, c1],  left.numChildrens  = 2
//	  promoted       = &2
func (n *node[k, v]) split() (*item[k, v], *node[k, v]) {
	mid := n.minItems
	midItem := n.items[mid]

	// Create right node and copy the upper half of items into it.
	newNode := newNode[k, v](n.maxItems, n.minItems, n.maxChildren)
	copy(newNode.items, n.items[mid+1:n.numItems])
	newNode.numItems = n.numItems - mid - 1

	if !n.isLeaf() {
		copy(newNode.children, n.children[mid+1:n.numChildrens])
		newNode.numChildrens = n.numChildrens - (mid + 1)
	}

	// Nil out promoted + moved items and children from left node for GC.
	for i, l := mid, n.numItems; i < l; i++ {
		n.items[i] = nil
		if !n.isLeaf() {
			n.children[i+1] = nil
			n.numChildrens--
		}
	}
	n.numItems = mid

	return midItem, newNode
}

// insert places item into the subtree rooted at n using a top-down,
// proactive-split strategy: every full child is split before descent so that
// insertion always terminates at a non-full leaf without backtracking.
//
// Returns true if a new key was added, false if an existing key was updated.
//
// Precondition: n itself is not full (guaranteed by BTree.Insert, which
// pre-splits the root before calling this method).
//
// # Algorithm
//
//  1. Binary-search items for the key.
//  2. Found → overwrite value in-place; return false.
//  3. Leaf  → insertItemAt(pos, item); return true.
//  4. children[pos] full → split it, promote midItem into n, adjust pos.
//  5. Recurse into children[pos].
//
// # Dry run — degree=2 (maxItems=3, minItems=1)
//
// Tree before inserting key=4:
//
//	       [5]                      root.items=[{5}], numItems=1
//	      /   \                     root.children=[[1,2,3],[7]], numChildrens=2
//	[1,2,3]   [7]
//
// root.insert({4}):
//
//	Step 1 — searchNode([5]): 4 < 5  →  (pos=0, found=false)
//	Step 2 — not found.
//	Step 3 — not leaf (numChildrens=2).
//	Step 4 — children[0]=[1,2,3] has numItems=3 == maxItems → split:
//	           split([1,2,3]):  left=[1]  midItem={2}  right=[3]
//	           insertItemAt(0, {2})  on root  →  root.items=[{2},{5}]
//	           insertChildAt(1,[3])  on root  →  root.children=[[1],[3],[7]]
//	           key=4 > root.items[pos=0].key=2  →  pos++ = 1
//	Step 5 — recurse: children[1]=[3].insert({4})
//	           searchNode([3]): 4>3 → (1, false); isLeaf → insertItemAt(1,{4}) → [3,4]
//
// Tree after:
//
//	   [2, 5]
//	  /  |   \
//	[1] [3,4] [7]
func (n *node[k, v]) insert(item *item[k, v]) bool {
	pos, found := n.searchNode(item.key)
	// pos is the index where item.key sits (found=true) or should go (found=false).

	// Key already exists — overwrite value, no structural change.
	if found {
		n.items[pos] = item
		return false
	}

	// Leaf with room — place the item here and we are done.
	if n.isLeaf() {
		n.insertItemAt(pos, item)
		return true
	}

	// Proactive split: if the child we need to descend into is full, split it
	// now so it can accept the promoted key without overflowing.
	if n.children[pos].numItems >= n.maxItems {
		// Dry run (continued from above):
		//   children[0]=[1,2,3] is full  →  split produces midItem={2}, right=[3]
		//   After insertItemAt+insertChildAt:
		//     n.items    = [{2}, {5}]
		//     n.children = [[1], [3], [7]]
		midItem, newNode := n.children[pos].split()
		n.insertItemAt(pos, midItem)
		n.insertChildAt(pos+1, newNode)

		// After promotion, n.items[pos] is the midItem just inserted.
		// Re-evaluate which child to follow based on item.key vs midItem.key.
		switch {
		case item.key < n.items[pos].key:
			// item.key is still in the left subtree — pos unchanged.
		case item.key > n.items[pos].key:
			// midItem.key < item.key → descend into the right (newly created) child.
			// Dry run: 4 > 2  →  pos = 1  (right child [3])
			pos++
		default:
			// item.key == midItem.key — the key we inserted matches the promoted
			// separator; overwrite it in the parent and stop.
			n.items[pos] = item
			return true
		}
	}

	return n.children[pos].insert(item)
}

// removeItemAt removes and returns the item at position pos, shifting
// items[pos+1..numItems-1] one slot left to close the gap.
//
// Dry run — items=[&{1},&{3},&{5},&{7}], numItems=4, removing pos=1:
//
//	removedItem = items[1] = &{3}
//	items[1] = nil                              →  [&{1}, nil,  &{5}, &{7}]
//	lastPos = 3;  pos(1) < lastPos(3)  → shift
//	copy items[1:3] ← items[2:4]               →  [&{1}, &{5}, &{7}, &{7}]
//	items[3] = nil                              →  [&{1}, &{5}, &{7}, nil ]
//	numItems = 3
//	return &{3}
func (n *node[k, v]) removeItemAt(pos int) *item[k, v] {
	removedItem := n.items[pos]
	n.items[pos] = nil
	// Shift items[pos+1..numItems-1] left by one to close the gap.
	if lastPos := n.numItems - 1; pos < lastPos {
		copy(n.items[pos:lastPos], n.items[pos+1:lastPos+1])
		n.items[lastPos] = nil
	}
	n.numItems--

	return removedItem
}

// removeChildAt removes and returns the child pointer at position pos, shifting
// children[pos+1..numChildrens-1] one slot left to close the gap.
//
// Dry run — children=[c0,c1,c2,c3], numChildrens=4, removing pos=1:
//
//	removedChild = children[1] = c1
//	children[1] = nil                            →  [c0, nil, c2, c3]
//	lastPos = 3;  pos(1) < lastPos(3)  → shift
//	copy children[1:3] ← children[2:4]           →  [c0, c2, c3, c3]
//	children[3] = nil                            →  [c0, c2, c3, nil]
//	numChildrens = 3
//	return c1
func (n *node[k, v]) removeChildAt(pos int) *node[k, v] {
	removedChild := n.children[pos]
	n.children[pos] = nil
	// Shift children[pos+1..numChildrens-1] left by one to close the gap.
	if lastPos := n.numChildrens - 1; pos < lastPos {
		copy(n.children[pos:lastPos], n.children[pos+1:lastPos+1])
		n.children[lastPos] = nil
	}
	n.numChildrens--

	return removedChild
}

// fillChildAt restores the minimum-occupancy invariant for children[pos] after
// a deletion left it with fewer than minItems keys. Three strategies are tried
// in priority order:
//
//  1. Rotate right — borrow the right-most key from the left sibling.
//  2. Rotate left  — borrow the left-most key from the right sibling.
//  3. Merge        — fuse children[pos] with a sibling around the parent
//     separator, shrinking the parent by one key and one child pointer.
//
// # Dry run — degree=2 (maxItems=3, minItems=1)
//
// ## Case 1 — rotate right (left sibling has surplus)
//
// Before:
//
//	parent.items    = [{5}],  numItems=1
//	parent.children = [left=[{1,3}], deficient=[{7}]]   ← deficient at pos=1
//
// fillChildAt(1): children[0].numItems=2 > minItems=1 → rotate right
//
//	Shift deficient items right: copy items[1:2] ← items[0:1]  →  deficient=[nil,{7}]
//	deficient.items[0] = parent.items[0] = {5}                 →  deficient=[{5},{7}]
//	deficient.numItems = 2
//	Pull left's right-most to parent:
//	  parent.items[0] = left.removeItemAt(1) = {3}             →  left=[{1}]
//
// After:
//
//	parent.items    = [{3}]
//	parent.children = [left=[{1}],  right=[{5},{7}]]
//
// ## Case 2 — rotate left (right sibling has surplus)
//
// Before:
//
//	parent.items    = [{5}],  numItems=1
//	parent.children = [deficient=[{3}], right=[{7,9}]]   ← deficient at pos=0
//
// fillChildAt(0): children[1].numItems=2 > minItems=1 → rotate left
//
//	Append parent separator to deficient: deficient.items[1]={5};  deficient=[{3},{5}]
//	deficient.numItems = 2
//	Pull right's left-most to parent:
//	  parent.items[0] = right.removeItemAt(0) = {7}            →  right=[{9}]
//
// After:
//
//	parent.items    = [{7}]
//	parent.children = [left=[{3},{5}],  right=[{9}]]
//
// ## Case 3 — merge (both siblings at minimum occupancy)
//
// Before:
//
//	parent.items    = [{5}],  numItems=1
//	parent.children = [left=[{3}], deficient=[{7}]]   ← deficient at pos=1
//
// fillChildAt(1): no surplus sibling → merge; pos(1) >= numItems(1) → pos=0
//
//	left=children[0]=[{3}],  right=children[1]=[{7}]
//	Pull parent separator into left: left.items[1]=parent.removeItemAt(0)={5}; left.numItems=2
//	Copy right's items into left:    left.items[2]={7};  left.numItems=3  →  left=[{3},{5},{7}]
//	parent.removeChildAt(1)  →  parent.numChildrens=1,  parent.numItems=0
//	right discarded (GC-eligible)
//
// After:
//
//	parent.items    = []               (caller collapses root if now empty)
//	parent.children = [left=[{3},{5},{7}]]
func (n *node[k, v]) fillChildAt(pos int) {
	switch {
	// Case 1: rotate right — borrow right-most key from left sibling.
	case pos > 0 && n.children[pos-1].numItems > n.minItems:
		left, right := n.children[pos-1], n.children[pos]
		// Shift right's items one slot right to open position 0 for the parent separator.
		copy(right.items[1:right.numItems+1], right.items[:right.numItems])
		right.items[0] = n.items[pos-1]
		right.numItems++
		// For internal nodes, adopt left's rightmost child as right's new leftmost child.
		if !right.isLeaf() {
			right.insertChildAt(0, left.removeChildAt(left.numChildrens-1))
		}
		// Pull left's rightmost key up to the parent separator slot.
		n.items[pos-1] = left.removeItemAt(left.numItems - 1)

	// Case 2: rotate left — borrow left-most key from right sibling.
	case pos < n.numChildrens-1 && n.children[pos+1].numItems > n.minItems:
		left, right := n.children[pos], n.children[pos+1]
		// Append parent separator key to left's end.
		left.items[left.numItems] = n.items[pos]
		left.numItems++
		// For internal nodes, adopt right's leftmost child as left's new rightmost child.
		if !left.isLeaf() {
			left.insertChildAt(left.numChildrens, right.removeChildAt(0))
		}
		// Pull right's leftmost key up to the parent separator slot.
		n.items[pos] = right.removeItemAt(0)

	// Case 3: merge — no sibling has a surplus key.
	default:
		// Prefer merging with right sibling; if pos points past the last separator,
		// step left so children[pos] and children[pos+1] are always valid.
		if pos >= n.numItems {
			pos = n.numItems - 1
		}
		left, right := n.children[pos], n.children[pos+1]
		// Pull the parent separator down into left as its new rightmost key.
		left.items[left.numItems] = n.removeItemAt(pos)
		left.numItems++
		// Move all of right's keys into left.
		copy(left.items[left.numItems:], right.items[:right.numItems])
		left.numItems += right.numItems
		// For internal nodes, move all of right's children into left.
		if !left.isLeaf() {
			copy(left.children[left.numChildrens:], right.children[:right.numChildrens])
			left.numChildrens += right.numChildrens
		}
		// Detach right from the parent; right is now unreachable and eligible for GC.
		n.removeChildAt(pos + 1)
		right = nil
	}
}

// delete removes the item with the given key from the subtree rooted at n and
// returns it, or returns nil if the key is absent. isSeekingSuccessor signals
// that the current descent is hunting for the in-order successor of a key that
// was found in an ancestor internal node.
//
// # Algorithm
//
//  1. Binary-search for key in n.items.
//  2. Found on leaf  → removeItemAt(pos); return item.
//  3. Found on internal node → set isSeekingSuccessor=true, descend into
//     children[pos+1] to find the leftmost key of the right subtree.
//  4. isSeekingSuccessor on leaf → removeItemAt(0) returns the successor.
//  5. Back at the internal node that matched → replace items[pos] with the
//     returned successor (logical deletion of the original key).
//  6. After recursion: if the child underflowed (numItems < minItems),
//     call fillChildAt to borrow from a sibling or merge.
//
// # Dry run — degree=2 (maxItems=3, minItems=1)
//
// Tree:
//
//	       [5]
//	      /   \
//	  [1,3]   [7,9]
//
// delete(5, false) at root:
//
//	searchNode([5]): found=true, pos=0
//	not leaf → isSeekingSuccessor=true, next=children[1]=[7,9]
//	recurse: [7,9].delete(5, true)
//	  isLeaf && isSeekingSuccessor → removeItemAt(0)={7}; return {7}
//	back at root: found && seeking → items[0]={7}          (root=[7])
//	next.numItems=1 == minItems → no underflow
//
// Tree after delete(5):
//
//	    [7]
//	   /   \
//	[1,3]  [9]
//
// delete(1, false) on that tree:
//
//	searchNode([7]): 1<7 → (0, false); next=children[0]=[1,3]
//	recurse: [1,3].delete(1, false)
//	  searchNode([1,3]): found=true, pos=0; isLeaf → removeItemAt(0)={1}; return {1}
//	back at root: not seeking → no items overwrite
//	next.numItems=1 == minItems → no underflow
//
// Tree after delete(1):
//
//	  [7]
//	 /   \
//	[3]  [9]
func (n *node[k, v]) delete(key k, isSeekingSuccessor bool) *item[k, v] {
	pos, found := n.searchNode(key)

	var next *node[k, v]

	if found {
		// Key lives here; leaf → direct removal, internal → chase successor.
		if n.isLeaf() {
			return n.removeItemAt(pos)
		}
		next, isSeekingSuccessor = n.children[pos+1], true
	} else {
		next = n.children[pos]
	}

	// Reached the leftmost leaf while chasing the successor — it sits at pos 0.
	if n.isLeaf() && isSeekingSuccessor {
		return n.removeItemAt(0)
	}

	// Key absent and no child to descend into.
	if next == nil {
		return nil
	}

	deletedItem := next.delete(key, isSeekingSuccessor)

	// Replace the original key with its in-order successor now that we are back
	// at the internal node that held it.
	if found && isSeekingSuccessor {
		n.items[pos] = deletedItem
	}

	// Repair underflow in the child we just descended into.
	if next.numItems < n.minItems {
		if found && isSeekingSuccessor {
			n.fillChildAt(pos + 1)
		} else {
			n.fillChildAt(pos)
		}
	}

	return deletedItem
}
