package searchtrees_test

import (
	"errors"
	"fmt"

	searchtrees "github.com/architagr/golang_collections/search-trees"
)

// ExampleNewBTree demonstrates creating a B-tree and performing basic operations.
func ExampleNewBTree() {
	bt := searchtrees.NewBTree[int, string](2)

	bt.Insert(5, "five")
	bt.Insert(3, "three")
	bt.Insert(7, "seven")

	val, _ := bt.Find(5)
	fmt.Println(val)
	// Output:
	// five
}

// ExampleBTree_Insert shows inserting new keys and updating an existing key.
func ExampleBTree_Insert() {
	bt := searchtrees.NewBTree[string, int](2)

	bt.Insert("apple", 1)
	bt.Insert("banana", 2)
	bt.Insert("cherry", 3)

	// Update an existing key.
	bt.Insert("banana", 99)

	val, _ := bt.Find("banana")
	fmt.Println(val)
	// Output:
	// 99
}

// ExampleBTree_Find demonstrates Find returning a value and the not-found error.
func ExampleBTree_Find() {
	bt := searchtrees.NewBTree[int, string](3)
	for _, kv := range []struct {
		k int
		v string
	}{{1, "one"}, {2, "two"}, {3, "three"}} {
		bt.Insert(kv.k, kv.v)
	}

	val, err := bt.Find(2)
	fmt.Println(val, err)

	_, err = bt.Find(99)
	fmt.Println(errors.Is(err, searchtrees.ErrKeyNotFound))
	// Output:
	// two <nil>
	// true
}

// ExampleBTree_Delete shows removing a key and verifying it is gone.
func ExampleBTree_Delete() {
	bt := searchtrees.NewBTree[int, string](2)
	bt.Insert(10, "ten")
	bt.Insert(20, "twenty")
	bt.Insert(30, "thirty")

	found := bt.Delete(20)
	fmt.Println(found)

	_, err := bt.Find(20)
	fmt.Println(errors.Is(err, searchtrees.ErrKeyNotFound))
	// Output:
	// true
	// true
}
