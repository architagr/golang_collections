package searchtrees

import (
	"errors"
	"fmt"
	"testing"
)

// ── helpers ──────────────────────────────────────────────────────────────────

func intTree(t *testing.T, degree int) *BTree[int, string] {
	t.Helper()
	return NewBTree[int, string](degree)
}

// ── NewBTree ─────────────────────────────────────────────────────────────────

func TestNewBTree_ValidDegree(t *testing.T) {
	for _, degree := range []int{2, 3, 5, 10, 100} {
		bt := NewBTree[int, string](degree)
		if bt == nil {
			t.Fatalf("degree %d: got nil BTree", degree)
		}
		if bt.degree != degree {
			t.Fatalf("degree %d: bt.degree = %d", degree, bt.degree)
		}
		if bt.maxItems != 2*degree-1 {
			t.Fatalf("degree %d: maxItems = %d, want %d", degree, bt.maxItems, 2*degree-1)
		}
		if bt.maxChildren != 2*degree {
			t.Fatalf("degree %d: maxChildren = %d, want %d", degree, bt.maxChildren, 2*degree)
		}
		if bt.minItems != degree-1 {
			t.Fatalf("degree %d: minItems = %d, want %d", degree, bt.minItems, degree-1)
		}
	}
}

func TestNewBTree_PanicsOnLowDegree(t *testing.T) {
	for _, degree := range []int{0, 1, -1, -100} {
		degree := degree
		t.Run("degree", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("degree %d: expected panic, got none", degree)
				}
			}()
			NewBTree[int, string](degree)
		})
	}
}

// ── Find ─────────────────────────────────────────────────────────────────────

func TestFind_EmptyTree(t *testing.T) {
	bt := intTree(t, 2)
	_, err := bt.Find(1)
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("empty tree: got err=%v, want ErrKeyNotFound", err)
	}
}

func TestFind_KeyNotFound(t *testing.T) {
	bt := intTree(t, 2)
	// Manually place one item into the pre-allocated root slice.
	bt.root.items[0] = &item[int, string]{key: 10, val: "ten"}
	bt.root.numItems = 1

	_, err := bt.Find(99)
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("got err=%v, want ErrKeyNotFound", err)
	}
}

func TestFind_RootLeaf_Found(t *testing.T) {
	bt := intTree(t, 2)
	// Place items directly into the pre-allocated root slice in sorted order.
	bt.root.items[0] = &item[int, string]{key: 1, val: "one"}
	bt.root.items[1] = &item[int, string]{key: 3, val: "three"}
	bt.root.items[2] = &item[int, string]{key: 5, val: "five"}
	bt.root.numItems = 3

	tests := []struct {
		key     int
		wantVal string
	}{
		{1, "one"},
		{3, "three"},
		{5, "five"},
	}
	for _, tc := range tests {
		got, err := bt.Find(tc.key)
		if err != nil {
			t.Errorf("key=%d: unexpected error: %v", tc.key, err)
		}
		if got != tc.wantVal {
			t.Errorf("key=%d: got %q, want %q", tc.key, got, tc.wantVal)
		}
	}
}

func TestFind_MultiLevel_Found(t *testing.T) {
	// Manually build a 2-level B-tree (degree 2) to test internal-node traversal.
	//
	//          [5]
	//         /   \
	//      [3]    [7]
	//
	// root: items=[{5,"five"}], children=[left, right]
	left := &node[int, string]{
		items:    []*item[int, string]{{3, "three"}},
		numItems: 1,
	}
	right := &node[int, string]{
		items:    []*item[int, string]{{7, "seven"}},
		numItems: 1,
	}
	root := &node[int, string]{
		items:        []*item[int, string]{{5, "five"}},
		children:     []*node[int, string]{left, right},
		numItems:     1,
		numChildrens: 2,
	}

	bt := intTree(t, 2)
	bt.root = root

	tests := []struct {
		key     int
		wantVal string
		wantErr error
	}{
		{5, "five", nil},
		{3, "three", nil},
		{7, "seven", nil},
		{1, "", ErrKeyNotFound},
		{4, "", ErrKeyNotFound},
		{9, "", ErrKeyNotFound},
	}

	for _, tc := range tests {
		got, err := bt.Find(tc.key)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("key=%d: got err=%v, want %v", tc.key, err, tc.wantErr)
		}
		if err == nil && got != tc.wantVal {
			t.Errorf("key=%d: got %q, want %q", tc.key, got, tc.wantVal)
		}
	}
}

// ── splitRoot ────────────────────────────────────────────────────────────────

func TestSplitRoot_CorrectNewRoot(t *testing.T) {
	// degree=2: maxItems=3, minItems=1, maxChildren=4
	// root = [1,3,5]  (full leaf)
	// After splitRoot:  new root = [3]  children=[[1],[5]]
	bt := intTree(t, 2)
	bt.root.items[0] = &item[int, string]{1, "a"}
	bt.root.items[1] = &item[int, string]{3, "b"}
	bt.root.items[2] = &item[int, string]{5, "c"}
	bt.root.numItems = 3

	bt.splitRoot()

	if bt.root.numItems != 1 {
		t.Fatalf("new root numItems = %d, want 1", bt.root.numItems)
	}
	if bt.root.items[0].key != 3 {
		t.Fatalf("new root items[0].key = %d, want 3", bt.root.items[0].key)
	}
	if bt.root.numChildrens != 2 {
		t.Fatalf("new root numChildrens = %d, want 2", bt.root.numChildrens)
	}
}

func TestSplitRoot_LeftChildCorrect(t *testing.T) {
	bt := intTree(t, 2)
	bt.root.items[0] = &item[int, string]{1, "a"}
	bt.root.items[1] = &item[int, string]{3, "b"}
	bt.root.items[2] = &item[int, string]{5, "c"}
	bt.root.numItems = 3

	bt.splitRoot()

	left := bt.root.children[0]
	if left == nil {
		t.Fatal("left child is nil")
	}
	if left.numItems != 1 {
		t.Fatalf("left.numItems = %d, want 1", left.numItems)
	}
	if left.items[0].key != 1 {
		t.Fatalf("left.items[0].key = %d, want 1", left.items[0].key)
	}
}

func TestSplitRoot_RightChildCorrect(t *testing.T) {
	bt := intTree(t, 2)
	bt.root.items[0] = &item[int, string]{1, "a"}
	bt.root.items[1] = &item[int, string]{3, "b"}
	bt.root.items[2] = &item[int, string]{5, "c"}
	bt.root.numItems = 3

	bt.splitRoot()

	right := bt.root.children[1]
	if right == nil {
		t.Fatal("right child is nil")
	}
	if right.numItems != 1 {
		t.Fatalf("right.numItems = %d, want 1", right.numItems)
	}
	if right.items[0].key != 5 {
		t.Fatalf("right.items[0].key = %d, want 5", right.items[0].key)
	}
}

func TestSplitRoot_NewRootCapacityParams(t *testing.T) {
	// New root must carry correct minItems so node.split() can use n.minItems
	// as an array index without panicking.
	bt := intTree(t, 2)
	bt.root.items[0] = &item[int, string]{1, "a"}
	bt.root.items[1] = &item[int, string]{3, "b"}
	bt.root.items[2] = &item[int, string]{5, "c"}
	bt.root.numItems = 3

	bt.splitRoot()

	if bt.root.minItems != bt.minItems {
		t.Fatalf("new root minItems = %d, want %d", bt.root.minItems, bt.minItems)
	}
	if bt.root.maxItems != bt.maxItems {
		t.Fatalf("new root maxItems = %d, want %d", bt.root.maxItems, bt.maxItems)
	}
	if bt.root.maxChildren != bt.maxChildren {
		t.Fatalf("new root maxChildren = %d, want %d", bt.root.maxChildren, bt.maxChildren)
	}
}

func TestSplitRoot_HeightIncreases(t *testing.T) {
	// Trigger two successive root splits to ensure height grows correctly and
	// minItems propagates so no index-out-of-range panic occurs.
	bt := intTree(t, 2) // maxItems=3
	for i := 1; i <= 8; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}
	if bt.root.isLeaf() {
		t.Fatal("root must not be leaf after multiple splits")
	}
	for i := 1; i <= 8; i++ {
		if _, err := bt.Find(i); err != nil {
			t.Errorf("key=%d: %v", i, err)
		}
	}
}

// ── Insert (TDD — tests written first, implementation pending) ───────────────
//
// These tests define the expected contract for Insert. They will fail until
// Insert is implemented.

func TestInsert_NilRoot(t *testing.T) {
	// Cover the nil-root guard inside Insert (NewBTree always pre-allocates,
	// but the guard exists for safety; force it by setting root to nil directly).
	bt := intTree(t, 2)
	bt.root = nil
	bt.Insert(42, "forty-two")

	got, err := bt.Find(42)
	if err != nil {
		t.Fatalf("after Insert on nil root: Find returned error: %v", err)
	}
	if got != "forty-two" {
		t.Fatalf("got %q, want %q", got, "forty-two")
	}
}

func TestDelete_NilRoot(t *testing.T) {
	// Cover the nil-root early-return inside Delete.
	bt := intTree(t, 2)
	bt.root = nil
	if bt.Delete(1) {
		t.Fatal("Delete on nil root: got true, want false")
	}
}

func TestInsert_SingleKey(t *testing.T) {
	bt := intTree(t, 2)
	bt.Insert(1, "one")

	got, err := bt.Find(1)
	if err != nil {
		t.Fatalf("after Insert(1): Find returned error: %v", err)
	}
	if got != "one" {
		t.Fatalf("after Insert(1): got %q, want %q", got, "one")
	}
}

func TestInsert_UpdatesExistingKey(t *testing.T) {
	bt := intTree(t, 2)
	bt.Insert(1, "one")
	bt.Insert(1, "ONE")

	got, err := bt.Find(1)
	if err != nil {
		t.Fatalf("after second Insert(1): Find returned error: %v", err)
	}
	if got != "ONE" {
		t.Fatalf("after second Insert(1): got %q, want %q", got, "ONE")
	}
}

func TestInsert_MultipleKeys_FindAll(t *testing.T) {
	bt := intTree(t, 2)
	pairs := []struct {
		key int
		val string
	}{
		{5, "five"}, {3, "three"}, {7, "seven"},
		{1, "one"}, {4, "four"}, {6, "six"}, {9, "nine"},
	}

	for _, p := range pairs {
		bt.Insert(p.key, p.val)
	}

	for _, p := range pairs {
		got, err := bt.Find(p.key)
		if err != nil {
			t.Errorf("key=%d: Find returned error: %v", p.key, err)
			continue
		}
		if got != p.val {
			t.Errorf("key=%d: got %q, want %q", p.key, got, p.val)
		}
	}
}

func TestInsert_TriggersSplit(t *testing.T) {
	// Degree 2 → maxItems = 3. Insert 4 keys to force a root split.
	bt := intTree(t, 2)
	for i := 1; i <= 4; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	// After split, root must no longer be a leaf.
	if bt.root.isLeaf() {
		t.Fatal("root should not be a leaf after a split")
	}

	// All keys must still be findable.
	for i := 1; i <= 4; i++ {
		if _, err := bt.Find(i); err != nil {
			t.Errorf("key=%d: Find returned error after split: %v", i, err)
		}
	}
}

func TestInsert_LargeDataset(t *testing.T) {
	bt := intTree(t, 3)
	const n = 1000
	for i := 0; i < n; i++ {
		bt.Insert(i, fmt.Sprintf("val-%d", i))
	}
	for i := 0; i < n; i++ {
		got, err := bt.Find(i)
		if err != nil {
			t.Errorf("key=%d: Find error: %v", i, err)
			continue
		}
		want := fmt.Sprintf("val-%d", i)
		if got != want {
			t.Errorf("key=%d: got %q, want %q", i, got, want)
		}
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDelete_KeyNotFound(t *testing.T) {
	bt := intTree(t, 2)
	if bt.Delete(99) {
		t.Fatal("Delete on empty tree: got true, want false")
	}
}

func TestDelete_SingleKey(t *testing.T) {
	bt := intTree(t, 2)
	bt.Insert(1, "one")

	if !bt.Delete(1) {
		t.Fatal("Delete(1): got false, want true")
	}
	if _, err := bt.Find(1); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("after Delete(1): Find returned %v, want ErrKeyNotFound", err)
	}
}

func TestDelete_LeafKey(t *testing.T) {
	bt := intTree(t, 2)
	for _, k := range []int{5, 3, 7} {
		bt.Insert(k, fmt.Sprintf("%d", k))
	}

	if !bt.Delete(3) {
		t.Fatal("Delete(3): got false, want true")
	}
	if _, err := bt.Find(3); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("after Delete(3): Find returned %v, want ErrKeyNotFound", err)
	}
	for _, k := range []int{5, 7} {
		if _, err := bt.Find(k); err != nil {
			t.Errorf("key=%d still expected after Delete(3): %v", k, err)
		}
	}
}

func TestDelete_InternalKey(t *testing.T) {
	// Force enough insertions to create an internal node, then delete the
	// internal key to exercise in-order successor replacement.
	bt := intTree(t, 2)
	for i := 1; i <= 7; i++ {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	if !bt.Delete(4) {
		t.Fatal("Delete(4): got false, want true")
	}
	if _, err := bt.Find(4); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("after Delete(4): Find returned %v, want ErrKeyNotFound", err)
	}
	for _, k := range []int{1, 2, 3, 5, 6, 7} {
		if _, err := bt.Find(k); err != nil {
			t.Errorf("key=%d still expected: %v", k, err)
		}
	}
}

func TestDelete_TriggersRebalance(t *testing.T) {
	bt := intTree(t, 2)
	keys := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, k := range keys {
		bt.Insert(k, fmt.Sprintf("%d", k))
	}

	deleteOrder := []int{1, 5, 10, 3, 7}
	remaining := map[int]bool{2: true, 4: true, 6: true, 8: true, 9: true}

	for _, k := range deleteOrder {
		if !bt.Delete(k) {
			t.Errorf("Delete(%d): got false, want true", k)
		}
	}

	for k := range remaining {
		if _, err := bt.Find(k); err != nil {
			t.Errorf("key=%d should still exist: %v", k, err)
		}
	}
	for _, k := range deleteOrder {
		if _, err := bt.Find(k); !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("key=%d should be deleted", k)
		}
	}
}
