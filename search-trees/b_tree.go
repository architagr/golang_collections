// Package searchtrees provides generic, ordered search-tree data structures.
//
// All types are safe for concurrent reads but require external synchronisation
// for concurrent writes.
package searchtrees

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrKeyNotFound is returned by Find when the requested key does not exist in
// the tree.
var ErrKeyNotFound = errors.New("key not found")

// BTree is a generic B-tree indexed by keys of any ordered type.
//
// A B-tree of degree d satisfies the following invariants:
//   - Every node holds at most 2d-1 keys.
//   - Every non-root internal node holds at least d-1 keys.
//   - Every internal node with m keys has exactly m+1 children.
//   - All leaves reside at the same depth.
//
// Degree 2 produces a 2-3-4 tree (each node holds 1–3 keys).
//
// Time complexity (n = number of stored keys):
//
//	Find:   O(log n · log d)   tree height × binary search per node
//	Insert: O(log n · d)       splits along root-to-leaf path
//	Delete: O(log n · d)       merges / borrows along root-to-leaf path
type BTree[k cmp.Ordered, v any] struct {
	root        *node[k, v]
	degree      int
	maxChildren int // 2 * degree
	maxItems    int // 2*degree - 1
	minItems    int // degree - 1   (not enforced on root)
}

// NewBTree creates an empty B-tree with the given degree.
//
// degree must be ≥ 2. Panics otherwise.
func NewBTree[k cmp.Ordered, v any](degree int) *BTree[k, v] {
	if degree < 2 {
		panic(fmt.Sprintf("btree: degree must be >= 2, got %d", degree))
	}
	maxChildren := 2 * degree   // max child pointers per node
	maxItems := maxChildren - 1 // max keys per node
	minItems := degree - 1      // min keys per non-root node

	return &BTree[k, v]{
		root:        newNode[k, v](maxItems, minItems, maxChildren),
		degree:      degree,
		maxChildren: maxChildren,
		maxItems:    maxItems,
		minItems:    minItems,
	}
}

// Find returns the value stored under key, or ErrKeyNotFound if key is absent.
//
// Traversal: at each node perform a binary search over the sorted items.
//   - Exact match → return the associated value immediately.
//   - No match on a leaf → key does not exist; return ErrKeyNotFound.
//   - No match on an internal node → descend into the child subtree indicated
//     by the index returned from searchNode and repeat.
func (bTree *BTree[k, v]) Find(key k) (v, error) {
	for next := bTree.root; next != nil; {
		index, found := next.searchNode(key)
		if found {
			return next.items[index].val, nil
		}
		// Leaf with no match: key is absent.
		if next.isLeaf() {
			break
		}
		next = next.children[index]
	}

	var zero v
	return zero, fmt.Errorf("key %+v: %w", key, ErrKeyNotFound)
}

// splitRoot handles the only case where tree height increases: the current root
// is full (numItems == maxItems) and must be split before a new key can be
// inserted.
//
// A brand-new empty root is created and the old root is split in two. The
// promoted median key is placed into the new root, which adopts the two halves
// as its children. bTree.root then points to this new root.
//
// # Dry run — degree=2 (maxItems=3, minItems=1)
//
// Before: root = [&{1}, &{3}, &{5}],  numItems=3  (full leaf)
//
//	newRoot = newNode(maxItems=3, minItems=1, maxChildren=4)
//	split(root):
//	  mid=1  →  midItem=&{3}
//	  right  = [&{5}],  right.numItems=1
//	  left   = [&{1}],  left.numItems=1    (root truncated in place)
//
//	newRoot.insertItemAt(0, &{3})           →  newRoot.items   = [&{3}]
//	newRoot.insertChildAt(0, left=[&{1}])   →  newRoot.children = [left, nil, ...]
//	newRoot.insertChildAt(1, right=[&{5}])  →  newRoot.children = [left, right, ...]
//	newRoot.numItems=1, newRoot.numChildrens=2
//
// After:
//
//	   [3]                ← new root
//	  /   \
//	[1]   [5]
func (bTree *BTree[k, v]) splitRoot() {
	newRoot := newNode[k, v](bTree.maxItems, bTree.minItems, bTree.maxChildren)
	midItem, newNode := bTree.root.split()
	newRoot.insertItemAt(0, midItem)
	newRoot.insertChildAt(0, bTree.root)
	newRoot.insertChildAt(1, newNode)
	bTree.root = newRoot
}

// Insert stores val under key, overwriting the existing value if key is already
// present.
//
// Strategy: top-down proactive split. Before descending, if the root is full,
// splitRoot is called — this is the only path that increases tree height.
// Thereafter, node.insert handles all child splits recursively without
// backtracking.
//
// # Dry run — degree=2 (maxItems=3), inserting keys 5, 3, 7, 1 in order
//
// Insert(5): root=[]  →  root=[5]
// Insert(3): root=[5]  →  root=[3,5]
// Insert(7): root=[3,5]  →  root=[3,5,7]   ← root now full
//
// Insert(1) — root is full (numItems=3 == maxItems):
//
//	splitRoot():
//	  oldRoot = [3,5,7]
//	  split → left=[3], midItem={5}, right=[7]
//	  newRoot.items    = [{5}]
//	  newRoot.children = [[3], [7]]
//	  bTree.root = newRoot
//
//	newRoot.insert({1}):
//	  searchNode([5]): 1<5 → (0, false)
//	  not leaf; children[0]=[3] not full
//	  recurse into [3].insert({1}):
//	    searchNode([3]): 1<3 → (0, false); isLeaf → insertItemAt(0,{1}) → [1,3]
//
//	tree after Insert(1):
//	        [5]
//	       /   \
//	    [1,3]  [7]
func (bTree *BTree[k, v]) Insert(key k, val v) {
	i := &item[k, v]{key, val}

	// Tree is empty: NewBTree pre-allocates the root, but guard nil just in case.
	if bTree.root == nil {
		bTree.root = newNode[k, v](bTree.maxItems, bTree.minItems, bTree.maxChildren)
	}

	// Root full: height must increase before descent.
	if bTree.root.numItems >= bTree.maxItems {
		bTree.splitRoot()
	}

	bTree.root.insert(i)
}

// Delete removes the entry for key from the tree.
// Returns true if the key was found and removed, false if the key did not exist.
//
// Strategy: bottom-up repair. node.delete descends to the target, removes it
// (replacing an internal key with its in-order successor if needed), then on
// the way back up calls fillChildAt on any child that underflowed. After the
// recursive call returns, if the root has no items left it is collapsed:
//   - Empty leaf root  → tree becomes empty (root = nil).
//   - Empty internal root → the sole remaining child becomes the new root,
//     shrinking tree height by one.
func (t *BTree[k, v]) Delete(key k) bool {
	if t.root == nil {
		return false
	}
	deletedItem := t.root.delete(key, false)

	// Collapse the root if the deletion drained it.
	if t.root.numItems == 0 {
		if t.root.isLeaf() {
			t.root = nil
		} else {
			t.root = t.root.children[0]
		}
	}

	return deletedItem != nil
}
