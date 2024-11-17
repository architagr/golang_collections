package benchmark

import (
	"fmt"
	"testing"

	"github.com/architagr/golang_collections/list"
	arraylist "github.com/emirpasic/gods/lists/arraylist"
)

type user struct {
	id   int
	name string
}

func initUser(a int) *user {
	return &user{
		id:   a,
		name: fmt.Sprintf("user_%d", a),
	}
}
func (usr *user) Copy() interface{} {
	cpy := new(user)
	cpy.id = usr.id
	cpy.name = usr.name
	return cpy
}
func (usr *user) Equal(val interface{}) bool {
	data, ok := val.(*user)
	if !ok {
		return false
	}
	return data.id == usr.id && data.name == usr.name
}

func Benchmark_ArrayListAdd_golang_collections(b *testing.B) {
	b.StopTimer()
	userList := list.InitArrayList[*user]()
	user := &user{
		id:   1,
		name: fmt.Sprintf("user_%d", 1),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		userList.Add(user)
	}
}

func Benchmark_ArrayListAdd_emirpasicGods(b *testing.B) {
	b.StopTimer()
	userList := arraylist.New()
	user := &user{
		id:   1,
		name: fmt.Sprintf("user_%d", 1),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		userList.Add(user)
	}
}

func Benchmark_ArrayListAddAtIndex_golang_collections(b *testing.B) {
	b.StopTimer()
	userList := list.InitArrayList[*user]()
	user := &user{
		id:   1,
		name: fmt.Sprintf("user_%d", 1),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		userList.AddAtIndex(0, user)
	}
}

func Benchmark_ArrayListAddAtIndex_emirpasicGods(b *testing.B) {
	b.StopTimer()
	userList := arraylist.New()
	user := &user{
		id:   1,
		name: fmt.Sprintf("user_%d", 1),
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		userList.Insert(0, user)
	}
}

func Benchmark_ArrayListFind_golang_collections(b *testing.B) {
	b.StopTimer()
	userList := list.InitArrayList(getUsersList(b.N)...)
	user := initUser(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		userList.Find(user)
	}
}

func Benchmark_ArrayListFind_emirpasicGods(b *testing.B) {
	b.StopTimer()
	userList := arraylist.New(getUsersList(b.N))
	user := initUser(1)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		userList.IndexOf(user)
	}
}

func getUsersList(n int) []*user {
	list := make([]*user, n)
	for i := 0; i < n; i++ {
		list[i] = initUser(i)
	}
	return list
}
