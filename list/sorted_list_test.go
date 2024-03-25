package list

import (
	"fmt"
	"testing"
)

// #region test initilization
// test if we are able to initilize a blank sorted list and have no element
func TestInitilizationOfBlankSortedItratorList(t *testing.T) {
	list := InitSortedList[user, *user](func(left, right *user) bool {
		return left.id < right.id
	})

	if list.Count() != 0 {
		t.Fatalf("sorted list not intilized")
	}
}

// #endregion

// #region test Add function
func TestAddMultipleRecordsToSortedList(t *testing.T) {
	testCases := []struct {
		input, expected []*user
	}{
		{
			input: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   20,
					name: "test20",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
			},
			expected: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   2,
					name: "test2",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   20,
					name: "test20",
				},
			},
		},
		{
			input: []*user{
				{
					id:   20,
					name: "test20",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
			},
			expected: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   2,
					name: "test2",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   20,
					name: "test20",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(tb *testing.T) {
			obj := InitSortedList(func(a, b *user) bool {
				return a.id < b.id
			})

			for _, usr := range tc.input {
				obj.Add(usr)
			}
			for j, expectedUser := range tc.expected {
				got, err := obj.Get(j)
				if err != nil {
					tb.Fatalf("error when trying to access index %d, error: %s", j, err.Error())
				}
				if !expectedUser.Equal(got) {
					tb.Fatalf("expected %+v at index %d but got %+v", expectedUser, j, got)
				}
			}
		})
	}
}

func TestAddMultipleRecordsToSortedListDecendingOrder(t *testing.T) {
	testCases := []struct {
		input, expected []*user
	}{
		{
			input: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   20,
					name: "test20",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
			},
			expected: []*user{
				{
					id:   20,
					name: "test20",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
				{
					id:   1,
					name: "test1",
				},
			},
		},
		{
			input: []*user{
				{
					id:   20,
					name: "test20",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
			},
			expected: []*user{
				{
					id:   20,
					name: "test20",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   1,
					name: "test1",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(tb *testing.T) {
			obj := InitSortedList(func(a, b *user) bool {
				return a.id > b.id
			})

			for _, usr := range tc.input {
				obj.Add(usr)
			}
			for j, expectedUser := range tc.expected {
				got, err := obj.Get(j)
				if err != nil {
					tb.Fatalf("error when trying to access index %d, error: %s", j, err.Error())
				}
				if !expectedUser.Equal(got) {
					tb.Fatalf("expected %+v at index %d but got %+v", expectedUser, j, got)
				}
			}
		})
	}
}

// #endregion

// #region test find
func TestFindDataInSortedList(t *testing.T) {
	initData := getSampleDataSortedList()
	testCases := []struct {
		input         *user
		expectedIndex int
	}{
		{
			input: &user{
				id:   2,
				name: "test2",
			},
			expectedIndex: 1,
		},
		{
			input: &user{
				id:   1,
				name: "test1",
			},
			expectedIndex: 0,
		},
		{
			input: &user{
				id:   2,
				name: "test21",
			},
			expectedIndex: -1,
		},
	}
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(tb *testing.T) {
			got := obj.Find(tc.input)

			if tc.expectedIndex != got {
				tb.Fatalf("expected index %d but got %d", tc.expectedIndex, got)
			}

		})
	}
}

func TestFindDataInSortedListDecending(t *testing.T) {
	initData := getSampleDataSortedList()
	testCases := []struct {
		input         *user
		expectedIndex int
	}{
		{
			input: &user{
				id:   2,
				name: "test2",
			},
			expectedIndex: 4,
		},
		{
			input: &user{
				id:   1,
				name: "test1",
			},
			expectedIndex: 5,
		},
		{
			input: &user{
				id:   2,
				name: "test21",
			},
			expectedIndex: -1,
		},
	}
	obj := InitSortedList(func(a, b *user) bool {
		return a.id > b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(tb *testing.T) {
			got := obj.Find(tc.input)

			if tc.expectedIndex != got {
				tb.Fatalf("expected index %d but got %d", tc.expectedIndex, got)
			}

		})
	}
}

// #endregion

// #region test get
func TestGetSortedListNegative(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	_, err := obj.Get(-1)
	if err == nil {
		t.Fatalf("trying to access negative index")
	}
}

// #endregion

// #region test Set
func TestSetSortedListNegative(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	err := obj.Set(-1, &user{
		id:   100,
		name: "test 100",
	})
	if err == nil {
		t.Fatalf("trying to set negative index")
	}
}

func TestSetSortedList(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	err := obj.Set(1, &user{
		id:   100,
		name: "test 100",
	})
	if err != nil {
		t.Fatalf("trying to set valid index")
	}
	index := obj.Find(&user{
		id:   100,
		name: "test 100",
	})

	if index == -1 {
		t.Fatalf("value was not correctly set")
	}
}

// #endregion

// #region test RemoveAtIndex
func TestRemoveAtIndexSortedListNegative(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	_, err := obj.RemoveAtIndex(-1)
	if err == nil {
		t.Fatalf("trying to remove negative index")
	}
}

func TestRemoveAtIndexSortedList(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	originalLen := obj.Count()
	data, err := obj.RemoveAtIndex(1)
	if err != nil {
		t.Fatalf("trying to remove valid index")
	}

	if originalLen <= obj.Count() || !data.Equal(&user{
		id:   2,
		name: "test2",
	}) {
		t.Fatalf("value was not correctly removed")
	}
}

func TestRemoveAtIndexSortedListDecending(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id > b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	originalLen := obj.Count()
	data, err := obj.RemoveAtIndex(1)
	if err != nil {
		t.Fatalf("trying to remove valid index")
	}

	if originalLen <= obj.Count() || !data.Equal(&user{
		id:   6,
		name: "test6",
	}) {
		t.Fatalf("value was not correctly removed")
	}
}

// #endregion

// #region test Remove
func TestRemoveSortedListNegative(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	_, err := obj.Remove(&user{
		id:   1,
		name: "test11",
	})
	if err == nil {
		t.Fatalf("trying to remove element that does not exist")
	}
}

func TestRemoveSortedList(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	originalLen := obj.Count()
	_, err := obj.Remove(&user{
		id:   20,
		name: "test20",
	})
	if err != nil {
		t.Fatalf("trying to remove valid value")
	}

	if originalLen <= obj.Count() {
		t.Fatalf("value was not correctly removed")
	}
}

// #endregion
func TestDeepCopySortedList(t *testing.T) {
	testCases := []struct {
		input, expected []*user
	}{
		{
			input: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   20,
					name: "test20",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
			},
			expected: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   2,
					name: "test2",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   20,
					name: "test20",
				},
			},
		},
		{
			input: []*user{
				{
					id:   20,
					name: "test20",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   2,
					name: "test2",
				},
			},
			expected: []*user{
				{
					id:   1,
					name: "test1",
				},
				{
					id:   1,
					name: "test1",
				},
				{
					id:   2,
					name: "test2",
				},
				{
					id:   3,
					name: "test3",
				},
				{
					id:   5,
					name: "test5",
				},
				{
					id:   6,
					name: "test6",
				},
				{
					id:   20,
					name: "test20",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d", i), func(tb *testing.T) {
			obj := InitSortedList(func(a, b *user) bool {
				return a.id < b.id
			})

			for _, usr := range tc.input {
				obj.Add(usr)
			}
			x := obj.DeepCopy()
			for j, expectedUser := range tc.expected {
				got := x[j]
				if !expectedUser.Equal(got) {
					tb.Fatalf("expected %+v at index %d but got %+v", expectedUser, j, got)
				}
			}
		})
	}
}

func getSampleDataSortedList() []*user {
	return []*user{
		{
			id:   1,
			name: "test1",
		},
		{
			id:   20,
			name: "test20",
		},
		{
			id:   5,
			name: "test5",
		},
		{
			id:   6,
			name: "test6",
		},
		{
			id:   3,
			name: "test3",
		},
		{
			id:   2,
			name: "test2",
		},
	}
}

// #region test filter
func TestFilterSortedListNegative(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	response := obj.Filter(func(val *user) bool {
		return val.id == 1 && val.name == "test"
	})
	if len(response) != 0 {
		t.Fatalf("Sorted list should return no data as a response of filter")
	}
}

func TestFilterSortedList(t *testing.T) {
	initData := getSampleDataSortedList()
	obj := InitSortedList(func(a, b *user) bool {
		return a.id < b.id
	})
	for _, usr := range initData {
		obj.Add(usr)
	}
	response := obj.Filter(func(val *user) bool {
		return val.id == 1
	})
	if len(response) != 1 {
		t.Fatalf("Sorted list should only return 1 data")
	}
}

// #endregion
