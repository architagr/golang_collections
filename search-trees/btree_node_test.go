package searchtrees

import "testing"

// ── helpers ───────────────────────────────────────────────────────────────────

// leafNode builds a leaf node with the given key-value pairs pre-loaded.
// Keys must be provided in ascending order.
func leafNode(t *testing.T, maxItems, minItems, maxChildren int, pairs ...struct {
	key int
	val string
}) *node[int, string] {
	t.Helper()
	n := newNode[int, string](maxItems, minItems, maxChildren)
	for i, p := range pairs {
		n.items[i] = &item[int, string]{key: p.key, val: p.val}
	}
	n.numItems = len(pairs)
	return n
}

// ── newNode ───────────────────────────────────────────────────────────────────

func TestNewNode_Fields(t *testing.T) {
	n := newNode[int, string](3, 1, 4)

	if n.maxItems != 3 {
		t.Errorf("maxItems = %d, want 3", n.maxItems)
	}
	if n.minItems != 1 {
		t.Errorf("minItems = %d, want 1", n.minItems)
	}
	if n.maxChildren != 4 {
		t.Errorf("maxChildren = %d, want 4", n.maxChildren)
	}
	if n.numItems != 0 {
		t.Errorf("numItems = %d, want 0", n.numItems)
	}
	if n.numChildrens != 0 {
		t.Errorf("numChildrens = %d, want 0", n.numChildrens)
	}
}

func TestNewNode_SliceLengths(t *testing.T) {
	n := newNode[int, string](3, 1, 4)

	if got := len(n.items); got != 3 {
		t.Errorf("len(items) = %d, want 3", got)
	}
	if got := len(n.children); got != 4 {
		t.Errorf("len(children) = %d, want 4", got)
	}
}

func TestNewNode_SlicesAllNil(t *testing.T) {
	n := newNode[int, string](3, 1, 4)

	for i, it := range n.items {
		if it != nil {
			t.Errorf("items[%d] = %v, want nil", i, it)
		}
	}
	for i, c := range n.children {
		if c != nil {
			t.Errorf("children[%d] = %v, want nil", i, c)
		}
	}
}

// ── isLeaf ────────────────────────────────────────────────────────────────────

func TestIsLeaf_NewNode(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	if !n.isLeaf() {
		t.Error("fresh node should be a leaf (numChildrens=0)")
	}
}

func TestIsLeaf_WithChildren(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	n.numChildrens = 1
	if n.isLeaf() {
		t.Error("node with numChildrens=1 should not be a leaf")
	}
}

// ── searchNode ────────────────────────────────────────────────────────────────

func TestSearchNode_EmptyNode(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	idx, found := n.searchNode(42)
	if found {
		t.Fatal("empty node: expected not found")
	}
	if idx != 0 {
		t.Fatalf("empty node: index = %d, want 0", idx)
	}
}

func TestSearchNode_SingleItem(t *testing.T) {
	n := &node[int, string]{
		items:    []*item[int, string]{{key: 5, val: "five"}},
		numItems: 1,
	}

	tests := []struct {
		key       int
		wantIdx   int
		wantFound bool
	}{
		{key: 5, wantIdx: 0, wantFound: true},  // exact match
		{key: 3, wantIdx: 0, wantFound: false}, // smaller → child[0]
		{key: 9, wantIdx: 1, wantFound: false}, // larger  → child[1]
	}

	for _, tc := range tests {
		idx, found := n.searchNode(tc.key)
		if found != tc.wantFound || idx != tc.wantIdx {
			t.Errorf("key=%d: got (%d, %v), want (%d, %v)",
				tc.key, idx, found, tc.wantIdx, tc.wantFound)
		}
	}
}

func TestSearchNode_MultipleItems(t *testing.T) {
	// Mirrors the dry-run in searchNode's documentation.
	// items = [&{1,"a"}, &{3,"b"}, &{5,"c"}, &{7,"d"}]
	n := &node[int, string]{
		items: []*item[int, string]{
			{1, "a"}, {3, "b"}, {5, "c"}, {7, "d"},
		},
		numItems: 4,
	}

	tests := []struct {
		key       int
		wantIdx   int
		wantFound bool
	}{
		{key: 1, wantIdx: 0, wantFound: true},  // leftmost item
		{key: 3, wantIdx: 1, wantFound: true},  // middle item
		{key: 7, wantIdx: 3, wantFound: true},  // rightmost item
		{key: 0, wantIdx: 0, wantFound: false}, // smaller than all → child[0]
		{key: 4, wantIdx: 2, wantFound: false}, // between 3 and 5 → child[2]
		{key: 8, wantIdx: 4, wantFound: false}, // larger than all  → child[4]
	}

	for _, tc := range tests {
		idx, found := n.searchNode(tc.key)
		if found != tc.wantFound || idx != tc.wantIdx {
			t.Errorf("key=%d: got (%d, %v), want (%d, %v)",
				tc.key, idx, found, tc.wantIdx, tc.wantFound)
		}
	}
}

// ── insertItemAt ──────────────────────────────────────────────────────────────

func TestInsertItemAt_AppendToEmpty(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	n.insertItemAt(0, &item[int, string]{key: 5, val: "five"})

	if n.numItems != 1 {
		t.Fatalf("numItems = %d, want 1", n.numItems)
	}
	if n.items[0].key != 5 {
		t.Errorf("items[0].key = %d, want 5", n.items[0].key)
	}
}

func TestInsertItemAt_AppendAtEnd(t *testing.T) {
	// Insert at tail (pos == numItems) — no shift needed.
	n := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "one"},
		struct{ key int; val string }{3, "three"},
	)
	n.insertItemAt(2, &item[int, string]{key: 5, val: "five"})

	if n.numItems != 3 {
		t.Fatalf("numItems = %d, want 3", n.numItems)
	}
	if n.items[2].key != 5 {
		t.Errorf("items[2].key = %d, want 5", n.items[2].key)
	}
}

func TestInsertItemAt_InsertAtHead(t *testing.T) {
	// Insert at pos=0 must shift all existing items right by one.
	n := leafNode(t, 4, 1, 5,
		struct{ key int; val string }{3, "three"},
		struct{ key int; val string }{5, "five"},
	)
	n.insertItemAt(0, &item[int, string]{key: 1, val: "one"})

	if n.numItems != 3 {
		t.Fatalf("numItems = %d, want 3", n.numItems)
	}
	wantKeys := []int{1, 3, 5}
	for i, wk := range wantKeys {
		if n.items[i].key != wk {
			t.Errorf("items[%d].key = %d, want %d", i, n.items[i].key, wk)
		}
	}
}

func TestInsertItemAt_InsertInMiddle(t *testing.T) {
	// Insert at pos=1 shifts items[1..] right.
	n := leafNode(t, 4, 1, 5,
		struct{ key int; val string }{1, "one"},
		struct{ key int; val string }{5, "five"},
	)
	n.insertItemAt(1, &item[int, string]{key: 3, val: "three"})

	if n.numItems != 3 {
		t.Fatalf("numItems = %d, want 3", n.numItems)
	}
	wantKeys := []int{1, 3, 5}
	for i, wk := range wantKeys {
		if n.items[i].key != wk {
			t.Errorf("items[%d].key = %d, want %d", i, n.items[i].key, wk)
		}
	}
}

// ── insertChildAt ─────────────────────────────────────────────────────────────

func TestInsertChildAt_AppendToEmpty(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	child := newNode[int, string](3, 1, 4)
	n.insertChildAt(0, child)

	if n.numChildrens != 1 {
		t.Fatalf("numChildrens = %d, want 1", n.numChildrens)
	}
	if n.children[0] != child {
		t.Error("children[0] should point to the inserted child")
	}
}

func TestInsertChildAt_AppendAtEnd(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	c0 := newNode[int, string](3, 1, 4)
	c1 := newNode[int, string](3, 1, 4)
	n.children[0] = c0
	n.numChildrens = 1

	n.insertChildAt(1, c1)

	if n.numChildrens != 2 {
		t.Fatalf("numChildrens = %d, want 2", n.numChildrens)
	}
	if n.children[1] != c1 {
		t.Error("children[1] should point to c1")
	}
}

func TestInsertChildAt_InsertAtHead(t *testing.T) {
	// Inserting at pos=0 must shift existing children right.
	n := newNode[int, string](3, 1, 4)
	c0 := newNode[int, string](3, 1, 4)
	c1 := newNode[int, string](3, 1, 4)
	n.children[0] = c0
	n.children[1] = c1
	n.numChildrens = 2

	cNew := newNode[int, string](3, 1, 4)
	n.insertChildAt(0, cNew)

	if n.numChildrens != 3 {
		t.Fatalf("numChildrens = %d, want 3", n.numChildrens)
	}
	if n.children[0] != cNew {
		t.Error("children[0] should be cNew")
	}
	if n.children[1] != c0 {
		t.Error("children[1] should be c0 (shifted right)")
	}
	if n.children[2] != c1 {
		t.Error("children[2] should be c1 (shifted right)")
	}
}

// ── split ─────────────────────────────────────────────────────────────────────

func TestSplit_LeafNode(t *testing.T) {
	// degree=2: maxItems=3, minItems=1, maxChildren=4
	// Full leaf: items = [&{1,a}, &{2,b}, &{3,c}], numItems=3
	//
	// Expected after split:
	//   left:  items[0] = &{1,a},  numItems=1
	//   mid:                &{2,b}
	//   right: items[0] = &{3,c},  numItems=1
	n := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{2, "b"},
		struct{ key int; val string }{3, "c"},
	)

	mid, right := n.split()

	// Check promoted key.
	if mid == nil {
		t.Fatal("mid item is nil")
	}
	if mid.key != 2 {
		t.Errorf("mid.key = %d, want 2", mid.key)
	}

	// Check left node (receiver).
	if n.numItems != 1 {
		t.Errorf("left.numItems = %d, want 1", n.numItems)
	}
	if n.items[0].key != 1 {
		t.Errorf("left.items[0].key = %d, want 1", n.items[0].key)
	}
	// Promoted + moved slots must be nil.
	if n.items[1] != nil || n.items[2] != nil {
		t.Error("left: items[1] and items[2] should be nil after split")
	}

	// Check right node.
	if right == nil {
		t.Fatal("right node is nil")
	}
	if right.numItems != 1 {
		t.Errorf("right.numItems = %d, want 1", right.numItems)
	}
	if right.items[0].key != 3 {
		t.Errorf("right.items[0].key = %d, want 3", right.items[0].key)
	}
}

func TestSplit_InternalNode(t *testing.T) {
	// degree=2: maxItems=3, minItems=1, maxChildren=4
	// Full internal node: items=[&{1}, &{2}, &{3}], 4 children c0..c3
	//
	// Expected after split:
	//   left:  items=[&{1}],     children=[c0,c1],  numItems=1, numChildrens=2
	//   mid:   &{2}
	//   right: items=[&{3}],     children=[c2,c3],  numItems=1, numChildrens=2
	n := newNode[int, string](3, 1, 4)
	n.items[0] = &item[int, string]{key: 1, val: "a"}
	n.items[1] = &item[int, string]{key: 2, val: "b"}
	n.items[2] = &item[int, string]{key: 3, val: "c"}
	n.numItems = 3

	c0 := newNode[int, string](3, 1, 4)
	c1 := newNode[int, string](3, 1, 4)
	c2 := newNode[int, string](3, 1, 4)
	c3 := newNode[int, string](3, 1, 4)
	n.children[0], n.children[1], n.children[2], n.children[3] = c0, c1, c2, c3
	n.numChildrens = 4

	mid, right := n.split()

	if mid.key != 2 {
		t.Errorf("mid.key = %d, want 2", mid.key)
	}

	// Left node checks.
	if n.numItems != 1 {
		t.Errorf("left.numItems = %d, want 1", n.numItems)
	}
	if n.numChildrens != 2 {
		t.Errorf("left.numChildrens = %d, want 2", n.numChildrens)
	}
	if n.children[0] != c0 || n.children[1] != c1 {
		t.Error("left should retain c0 and c1")
	}
	if n.children[2] != nil || n.children[3] != nil {
		t.Error("left: children[2] and children[3] should be nil after split")
	}

	// Right node checks.
	if right.numItems != 1 {
		t.Errorf("right.numItems = %d, want 1", right.numItems)
	}
	if right.numChildrens != 2 {
		t.Errorf("right.numChildrens = %d, want 2", right.numChildrens)
	}
	if right.children[0] != c2 || right.children[1] != c3 {
		t.Error("right should contain c2 and c3")
	}
}

func TestSplit_CapacityPreserved(t *testing.T) {
	// After split, both nodes must have the same maxItems/minItems/maxChildren
	// as the original so further inserts and splits work correctly.
	n := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{2, "b"},
		struct{ key int; val string }{3, "c"},
	)
	_, right := n.split()

	if n.maxItems != 3 || n.minItems != 1 || n.maxChildren != 4 {
		t.Errorf("left capacity changed: maxItems=%d minItems=%d maxChildren=%d",
			n.maxItems, n.minItems, n.maxChildren)
	}
	if right.maxItems != 3 || right.minItems != 1 || right.maxChildren != 4 {
		t.Errorf("right capacity wrong: maxItems=%d minItems=%d maxChildren=%d",
			right.maxItems, right.minItems, right.maxChildren)
	}
}

// ── node.insert ───────────────────────────────────────────────────────────────

func TestNodeInsert_IntoEmptyLeaf(t *testing.T) {
	// degree=2 leaf, no items yet.
	// insert({5,"five"}) → items=[{5}], numItems=1, returns true (new key).
	n := newNode[int, string](3, 1, 4)
	added := n.insert(&item[int, string]{key: 5, val: "five"})

	if !added {
		t.Error("expected true (new key), got false")
	}
	if n.numItems != 1 {
		t.Fatalf("numItems = %d, want 1", n.numItems)
	}
	if n.items[0].key != 5 || n.items[0].val != "five" {
		t.Errorf("items[0] = %+v, want {5 five}", n.items[0])
	}
}

func TestNodeInsert_LeafMaintainsSortOrder(t *testing.T) {
	// Insert keys 5, 1, 3 (out of order) into a leaf with capacity 3.
	// Each insert returns true; final order must be [1, 3, 5].
	//
	// Step-by-step:
	//   insert({5}): leaf=[] → [5]
	//   insert({1}): searchNode([5]): 1<5 → pos=0; insertItemAt(0,{1}) → [1,5]
	//   insert({3}): searchNode([1,5]): 3>1,3<5 → pos=1; insertItemAt(1,{3}) → [1,3,5]
	n := newNode[int, string](3, 1, 4)
	for _, kv := range []struct{ k int; v string }{{5, "e"}, {1, "a"}, {3, "c"}} {
		if !n.insert(&item[int, string]{key: kv.k, val: kv.v}) {
			t.Errorf("insert(%d): expected true", kv.k)
		}
	}

	if n.numItems != 3 {
		t.Fatalf("numItems = %d, want 3", n.numItems)
	}
	wantKeys := []int{1, 3, 5}
	for i, wk := range wantKeys {
		if n.items[i].key != wk {
			t.Errorf("items[%d].key = %d, want %d", i, n.items[i].key, wk)
		}
	}
}

func TestNodeInsert_UpdateExistingKey(t *testing.T) {
	// Inserting a key that already exists must overwrite the value and return false.
	//
	// insert({5,"first"}):  returns true,  items=[{5,"first"}]
	// insert({5,"second"}): returns false, items=[{5,"second"}]
	n := newNode[int, string](3, 1, 4)
	n.insert(&item[int, string]{key: 5, val: "first"})

	added := n.insert(&item[int, string]{key: 5, val: "second"})

	if added {
		t.Error("expected false (update), got true")
	}
	if n.numItems != 1 {
		t.Errorf("numItems = %d, want 1 (no extra item)", n.numItems)
	}
	if n.items[0].val != "second" {
		t.Errorf("val = %q, want %q", n.items[0].val, "second")
	}
}

func TestNodeInsert_InternalNodeNoSplit(t *testing.T) {
	// Tree (degree=2):
	//        [5]
	//       /   \
	//     [3]   [7]
	//
	// insert({2}): 2<5 → children[0]=[3] not full → recurse → [2,3]
	//
	// Tree after:
	//        [5]
	//       /   \
	//     [2,3] [7]
	left := leafNode(t, 3, 1, 4, struct{ key int; val string }{3, "c"})
	right := leafNode(t, 3, 1, 4, struct{ key int; val string }{7, "g"})

	root := newNode[int, string](3, 1, 4)
	root.items[0] = &item[int, string]{key: 5, val: "e"}
	root.numItems = 1
	root.children[0], root.children[1] = left, right
	root.numChildrens = 2

	added := root.insert(&item[int, string]{key: 2, val: "b"})

	if !added {
		t.Error("expected true (new key)")
	}
	if left.numItems != 2 {
		t.Fatalf("left.numItems = %d, want 2", left.numItems)
	}
	if left.items[0].key != 2 || left.items[1].key != 3 {
		t.Errorf("left items = [%d,%d], want [2,3]",
			left.items[0].key, left.items[1].key)
	}
	// Root must be unchanged.
	if root.numItems != 1 || root.items[0].key != 5 {
		t.Error("root should be unchanged")
	}
}

func TestNodeInsert_TriggersChildSplit_KeyGoesRight(t *testing.T) {
	// Tree (degree=2, maxItems=3):
	//          [5]
	//         /   \
	//   [1,2,3]   [7]       ← children[0] is full
	//
	// insert({4}): 4<5 → children[0] full → split([1,2,3]):
	//   left=[1], midItem={2}, right=[3]
	//   root becomes [{2},{5}], children=[[1],[3],[7]]
	//   4 > 2  →  pos=1  →  recurse into [3]: insert({4}) → [3,4]
	//
	// Tree after:
	//       [2, 5]
	//      /  |   \
	//    [1] [3,4] [7]
	left := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{2, "b"},
		struct{ key int; val string }{3, "c"},
	)
	right := leafNode(t, 3, 1, 4, struct{ key int; val string }{7, "g"})

	root := newNode[int, string](3, 1, 4)
	root.items[0] = &item[int, string]{key: 5, val: "e"}
	root.numItems = 1
	root.children[0], root.children[1] = left, right
	root.numChildrens = 2

	root.insert(&item[int, string]{key: 4, val: "d"})

	// Root now has two items.
	if root.numItems != 2 {
		t.Fatalf("root.numItems = %d, want 2", root.numItems)
	}
	if root.items[0].key != 2 || root.items[1].key != 5 {
		t.Errorf("root items = [%d,%d], want [2,5]",
			root.items[0].key, root.items[1].key)
	}
	if root.numChildrens != 3 {
		t.Fatalf("root.numChildrens = %d, want 3", root.numChildrens)
	}
	// children[0] is the left half after split → only key 1.
	if root.children[0].numItems != 1 || root.children[0].items[0].key != 1 {
		t.Error("children[0] should be [1]")
	}
	// children[1] is the right half after split, then got key 4 inserted.
	if root.children[1].numItems != 2 {
		t.Fatalf("children[1].numItems = %d, want 2", root.children[1].numItems)
	}
	if root.children[1].items[0].key != 3 || root.children[1].items[1].key != 4 {
		t.Errorf("children[1] items = [%d,%d], want [3,4]",
			root.children[1].items[0].key, root.children[1].items[1].key)
	}
}

func TestNodeInsert_TriggersChildSplit_KeyGoesLeft(t *testing.T) {
	// Tree (degree=2):
	//          [5]
	//         /   \
	//   [1,2,3]   [7]
	//
	// insert({0}): 0<5 → children[0] full → split([1,2,3]):
	//   left=[1], midItem={2}, right=[3]
	//   root becomes [{2},{5}]
	//   0 < 2  →  pos stays 0  →  recurse into [1]: insert({0}) → [0,1]
	left := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{2, "b"},
		struct{ key int; val string }{3, "c"},
	)
	right := leafNode(t, 3, 1, 4, struct{ key int; val string }{7, "g"})

	root := newNode[int, string](3, 1, 4)
	root.items[0] = &item[int, string]{key: 5, val: "e"}
	root.numItems = 1
	root.children[0], root.children[1] = left, right
	root.numChildrens = 2

	root.insert(&item[int, string]{key: 0, val: "z"})

	if root.numItems != 2 {
		t.Fatalf("root.numItems = %d, want 2", root.numItems)
	}
	// children[0] is [1] after split, then {0} inserted → [0,1]
	if root.children[0].numItems != 2 {
		t.Fatalf("children[0].numItems = %d, want 2", root.children[0].numItems)
	}
	if root.children[0].items[0].key != 0 || root.children[0].items[1].key != 1 {
		t.Errorf("children[0] items = [%d,%d], want [0,1]",
			root.children[0].items[0].key, root.children[0].items[1].key)
	}
}

// ── removeItemAt ─────────────────────────────────────────────────────────────

func TestRemoveItemAt_Tail(t *testing.T) {
	// Remove last item — no shift needed.
	// items=[{1},{3},{5}], numItems=3, remove pos=2 → [{1},{3}]
	n := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{3, "b"},
		struct{ key int; val string }{5, "c"},
	)

	got := n.removeItemAt(2)

	if got.key != 5 {
		t.Fatalf("removed item key = %d, want 5", got.key)
	}
	if n.numItems != 2 {
		t.Fatalf("numItems = %d, want 2", n.numItems)
	}
	if n.items[2] != nil {
		t.Errorf("items[2] should be nil after removal")
	}
}

func TestRemoveItemAt_Head(t *testing.T) {
	// Remove first item — remaining items shift left by one.
	// items=[{1},{3},{5}], remove pos=0 → [{3},{5}]
	n := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{3, "b"},
		struct{ key int; val string }{5, "c"},
	)

	got := n.removeItemAt(0)

	if got.key != 1 {
		t.Fatalf("removed item key = %d, want 1", got.key)
	}
	if n.numItems != 2 {
		t.Fatalf("numItems = %d, want 2", n.numItems)
	}
	if n.items[0].key != 3 {
		t.Errorf("items[0].key = %d, want 3", n.items[0].key)
	}
	if n.items[1].key != 5 {
		t.Errorf("items[1].key = %d, want 5", n.items[1].key)
	}
	if n.items[2] != nil {
		t.Errorf("items[2] should be nil after removal")
	}
}

func TestRemoveItemAt_Middle(t *testing.T) {
	// Remove middle item — items to its right shift left.
	// items=[{1},{3},{5},{7}], numItems=4, remove pos=1 → [{1},{5},{7}]
	n := newNode[int, string](5, 2, 6)
	n.items[0] = &item[int, string]{1, "a"}
	n.items[1] = &item[int, string]{3, "b"}
	n.items[2] = &item[int, string]{5, "c"}
	n.items[3] = &item[int, string]{7, "d"}
	n.numItems = 4

	got := n.removeItemAt(1)

	if got.key != 3 {
		t.Fatalf("removed item key = %d, want 3", got.key)
	}
	if n.numItems != 3 {
		t.Fatalf("numItems = %d, want 3", n.numItems)
	}
	if n.items[0].key != 1 {
		t.Errorf("items[0].key = %d, want 1", n.items[0].key)
	}
	if n.items[1].key != 5 {
		t.Errorf("items[1].key = %d, want 5", n.items[1].key)
	}
	if n.items[2].key != 7 {
		t.Errorf("items[2].key = %d, want 7", n.items[2].key)
	}
	if n.items[3] != nil {
		t.Errorf("items[3] should be nil after removal")
	}
}

func TestRemoveItemAt_SingleItem(t *testing.T) {
	n := leafNode(t, 3, 1, 4, struct{ key int; val string }{42, "x"})

	got := n.removeItemAt(0)

	if got.key != 42 {
		t.Fatalf("removed item key = %d, want 42", got.key)
	}
	if n.numItems != 0 {
		t.Fatalf("numItems = %d, want 0", n.numItems)
	}
	if n.items[0] != nil {
		t.Errorf("items[0] should be nil after removal")
	}
}

// ── removeChildAt ─────────────────────────────────────────────────────────────

func TestRemoveChildAt_Tail(t *testing.T) {
	// Remove last child — no shift needed.
	// children=[c0,c1,c2], numChildrens=3, remove pos=2 → [c0,c1]
	n := newNode[int, string](3, 1, 4)
	c0 := newNode[int, string](3, 1, 4)
	c1 := newNode[int, string](3, 1, 4)
	c2 := newNode[int, string](3, 1, 4)
	n.children[0], n.children[1], n.children[2] = c0, c1, c2
	n.numChildrens = 3

	got := n.removeChildAt(2)

	if got != c2 {
		t.Fatal("returned wrong child")
	}
	if n.numChildrens != 2 {
		t.Fatalf("numChildrens = %d, want 2", n.numChildrens)
	}
	if n.children[2] != nil {
		t.Errorf("children[2] should be nil after removal")
	}
}

func TestRemoveChildAt_Head(t *testing.T) {
	// Remove first child — remaining shift left.
	// children=[c0,c1,c2], remove pos=0 → [c1,c2]
	n := newNode[int, string](3, 1, 4)
	c0 := newNode[int, string](3, 1, 4)
	c1 := newNode[int, string](3, 1, 4)
	c2 := newNode[int, string](3, 1, 4)
	n.children[0], n.children[1], n.children[2] = c0, c1, c2
	n.numChildrens = 3

	got := n.removeChildAt(0)

	if got != c0 {
		t.Fatal("returned wrong child")
	}
	if n.numChildrens != 2 {
		t.Fatalf("numChildrens = %d, want 2", n.numChildrens)
	}
	if n.children[0] != c1 {
		t.Errorf("children[0] should be c1 after head removal")
	}
	if n.children[1] != c2 {
		t.Errorf("children[1] should be c2 after head removal")
	}
	if n.children[2] != nil {
		t.Errorf("children[2] should be nil after removal")
	}
}

func TestRemoveChildAt_Middle(t *testing.T) {
	// Remove middle child — children to its right shift left.
	// children=[c0,c1,c2,c3], numChildrens=4, remove pos=1 → [c0,c2,c3]
	n := newNode[int, string](5, 2, 6)
	c0 := newNode[int, string](5, 2, 6)
	c1 := newNode[int, string](5, 2, 6)
	c2 := newNode[int, string](5, 2, 6)
	c3 := newNode[int, string](5, 2, 6)
	n.children[0], n.children[1], n.children[2], n.children[3] = c0, c1, c2, c3
	n.numChildrens = 4

	got := n.removeChildAt(1)

	if got != c1 {
		t.Fatal("returned wrong child")
	}
	if n.numChildrens != 3 {
		t.Fatalf("numChildrens = %d, want 3", n.numChildrens)
	}
	if n.children[0] != c0 {
		t.Errorf("children[0] should be c0")
	}
	if n.children[1] != c2 {
		t.Errorf("children[1] should be c2")
	}
	if n.children[2] != c3 {
		t.Errorf("children[2] should be c3")
	}
	if n.children[3] != nil {
		t.Errorf("children[3] should be nil after removal")
	}
}

func TestRemoveChildAt_SingleChild(t *testing.T) {
	n := newNode[int, string](3, 1, 4)
	c0 := newNode[int, string](3, 1, 4)
	n.children[0] = c0
	n.numChildrens = 1

	got := n.removeChildAt(0)

	if got != c0 {
		t.Fatal("returned wrong child")
	}
	if n.numChildrens != 0 {
		t.Fatalf("numChildrens = %d, want 0", n.numChildrens)
	}
	if n.children[0] != nil {
		t.Errorf("children[0] should be nil after removal")
	}
}

// ── fillChildAt ──────────────────────────────────────────────────────────────

func TestFillChildAt_RotateRight_InternalNode(t *testing.T) {
	// Covers the !right.isLeaf() branch in Case 1 (rotate right).
	//
	// Tree before fillChildAt(1):
	//
	//              [10]
	//             /    \
	//         [3,7]   [13]          ← right is at minItems, but we force
	//         / | \   /  \            fillChildAt to trigger Case 1 since
	//        c0 c1 c2 c3  c4          left has surplus (numItems=2 > minItems=1)
	//
	// After fillChildAt(1) (rotate right):
	//
	//              [7]
	//             /   \
	//           [3]  [10,13]
	//           / \  / |  \
	//          c0 c1 c2 c3 c4
	//
	// Specifically: c2 (left's rightmost child) moves to right's leftmost slot.

	// Build leaf grandchildren.
	c0 := leafNode(t, 3, 1, 4, struct{ key int; val string }{1, "a"})
	c1 := leafNode(t, 3, 1, 4, struct{ key int; val string }{5, "b"})
	c2 := leafNode(t, 3, 1, 4, struct{ key int; val string }{6, "c"})
	c3 := leafNode(t, 3, 1, 4, struct{ key int; val string }{11, "d"})
	c4 := leafNode(t, 3, 1, 4, struct{ key int; val string }{15, "e"})

	// Left internal child: items=[{3},{7}], children=[c0,c1,c2].
	left := newNode[int, string](3, 1, 4)
	left.items[0] = &item[int, string]{3, "three"}
	left.items[1] = &item[int, string]{7, "seven"}
	left.numItems = 2
	left.children[0], left.children[1], left.children[2] = c0, c1, c2
	left.numChildrens = 3

	// Right internal child (deficient in practice but we test the rotate path):
	// items=[{13}], children=[c3,c4].
	right := newNode[int, string](3, 1, 4)
	right.items[0] = &item[int, string]{13, "thirteen"}
	right.numItems = 1
	right.children[0], right.children[1] = c3, c4
	right.numChildrens = 2

	// Parent: items=[{10}], children=[left, right].
	parent := newNode[int, string](3, 1, 4)
	parent.items[0] = &item[int, string]{10, "ten"}
	parent.numItems = 1
	parent.children[0], parent.children[1] = left, right
	parent.numChildrens = 2

	parent.fillChildAt(1)

	// Parent separator must now be {7}.
	if parent.items[0].key != 7 {
		t.Errorf("parent.items[0].key = %d, want 7", parent.items[0].key)
	}

	// Left must have shrunk to [{3}] with children [c0,c1].
	if left.numItems != 1 || left.items[0].key != 3 {
		t.Errorf("left: numItems=%d items[0].key=%d, want 1/3", left.numItems, left.items[0].key)
	}
	if left.numChildrens != 2 || left.children[0] != c0 || left.children[1] != c1 {
		t.Errorf("left children wrong after rotate right")
	}

	// Right must now have [{10},{13}] with children [c2,c3,c4].
	if right.numItems != 2 || right.items[0].key != 10 || right.items[1].key != 13 {
		t.Errorf("right items wrong: numItems=%d", right.numItems)
	}
	if right.numChildrens != 3 || right.children[0] != c2 || right.children[1] != c3 || right.children[2] != c4 {
		t.Errorf("right children wrong after rotate right: numChildrens=%d", right.numChildrens)
	}
}

func TestNodeInsert_TriggersChildSplit_KeyEqualsMid(t *testing.T) {
	// Tree (degree=2):
	//          [5]
	//         /   \
	//   [1,2,3]   [7]
	//
	// insert({2,"NEW"}): 2<5 → children[0] full → split([1,2,3]):
	//   midItem = {2,"b"}
	//   root: insertItemAt(0, {2,"b"}) → root.items=[{2},{5}]
	//   key=2 == root.items[0].key=2 → default case:
	//     root.items[0] = {2,"NEW"}; return true
	left := leafNode(t, 3, 1, 4,
		struct{ key int; val string }{1, "a"},
		struct{ key int; val string }{2, "b"},
		struct{ key int; val string }{3, "c"},
	)
	right := leafNode(t, 3, 1, 4, struct{ key int; val string }{7, "g"})

	root := newNode[int, string](3, 1, 4)
	root.items[0] = &item[int, string]{key: 5, val: "e"}
	root.numItems = 1
	root.children[0], root.children[1] = left, right
	root.numChildrens = 2

	root.insert(&item[int, string]{key: 2, val: "NEW"})

	// The promoted mid-key (2) should now hold "NEW".
	if root.numItems != 2 {
		t.Fatalf("root.numItems = %d, want 2", root.numItems)
	}
	if root.items[0].key != 2 || root.items[0].val != "NEW" {
		t.Errorf("root.items[0] = %+v, want {2 NEW}", root.items[0])
	}
}
