# B-Tree

A generic, ordered B-tree (`BTree[K, V]`) backed by pre-allocated fixed-length node slices and explicit occupancy counters. Keys must satisfy `cmp.Ordered`; values can be any type.

## Properties

| Property | Value |
|---|---|
| Degree `d` | configurable (min 2) |
| Keys per node | `d-1` (min, non-root) … `2d-1` (max) |
| Children per internal node | `numItems + 1` |
| Insertion strategy | top-down proactive split (no backtracking) |
| Deletion strategy | bottom-up repair (borrow / merge) |

## Time complexity

| Operation | Average / Worst |
|---|---|
| `Find` | `O(log n · log d)` |
| `Insert` | `O(log n · d)` |
| `Delete` | `O(log n · d)` |

## Quick Start

```go
package main

import (
    "fmt"

    searchtrees "github.com/architagr/golang_collections/search-trees"
)

func main() {
    // Create a B-tree of degree 3 (each node holds 2–5 keys).
    bt := searchtrees.NewBTree[int, string](3)

    // Insert key-value pairs.
    bt.Insert(10, "ten")
    bt.Insert(5, "five")
    bt.Insert(20, "twenty")
    bt.Insert(1, "one")
    bt.Insert(15, "fifteen")

    // Find a key.
    val, err := bt.Find(15)
    if err != nil {
        panic(err)
    }
    fmt.Println(val) // fifteen

    // Update an existing key — Insert overwrites.
    bt.Insert(10, "TEN")
    val, _ = bt.Find(10)
    fmt.Println(val) // TEN

    // Delete a key.
    found := bt.Delete(5)
    fmt.Println(found) // true

    _, err = bt.Find(5)
    fmt.Println(err) // key 5: key not found
}
```

## API Reference

### `NewBTree[K, V](degree int) *BTree[K, V]`

Creates an empty B-tree. Panics if `degree < 2`.

```go
bt := searchtrees.NewBTree[string, int](2) // 2-3-4 tree
```

### `Insert(key K, val V)`

Inserts or updates the entry for `key`. Tree height increases only when the root is full (top-down proactive split).

```go
bt.Insert("hello", 42)
bt.Insert("hello", 99) // overwrites; Find("hello") → 99
```

### `Find(key K) (V, error)`

Returns the value for `key`, or `ErrKeyNotFound` if absent.

```go
val, err := bt.Find("hello")
if errors.Is(err, searchtrees.ErrKeyNotFound) {
    fmt.Println("not found")
}
```

### `Delete(key K) bool`

Removes `key` from the tree. Returns `true` if found and removed, `false` if the key did not exist. Tree height may decrease when the root is drained.

```go
deleted := bt.Delete("hello") // true
deleted  = bt.Delete("hello") // false — already gone
```

### `ErrKeyNotFound`

Sentinel error returned (wrapped) by `Find` when a key is absent. Test with `errors.Is`.

```go
_, err := bt.Find(99)
errors.Is(err, searchtrees.ErrKeyNotFound) // true
```

## Choosing a Degree

- **Degree 2** — 2-3-4 tree; smallest node size, most splits/merges per operation. Good for small datasets or when minimising memory per node matters.
- **Degree 3–5** — balanced trade-off for general-purpose use.
- **Degree 10+** — wide nodes reduce tree height and are cache-friendly for large in-memory datasets.

## Error Handling

```go
bt := searchtrees.NewBTree[int, string](2)

_, err := bt.Find(42)
if errors.Is(err, searchtrees.ErrKeyNotFound) {
    fmt.Println("key 42 is not in the tree")
}
```
