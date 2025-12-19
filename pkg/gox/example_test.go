package gox_test

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lwmacct/251219-go-pkg-ginm/pkg/gox"
)

func ExampleMap() {
	nums := []int{1, 2, 3, 4, 5}
	doubled := gox.Map(nums, func(n int) int { return n * 2 })
	fmt.Println(doubled)
	// Output:
	// [2 4 6 8 10]
}

func ExampleFilter() {
	nums := []int{1, 2, 3, 4, 5}
	evens := gox.Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens)
	// Output:
	// [2 4]
}

func ExampleReduce() {
	nums := []int{1, 2, 3, 4, 5}
	sum := gox.Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Println(sum)
	// Output:
	// 15
}

func ExampleFind() {
	users := []string{"alice", "bob", "charlie"}
	user, found := gox.Find(users, func(s string) bool {
		return strings.HasPrefix(s, "b")
	})
	fmt.Println(user, found)
	// Output:
	// bob true
}

func ExampleContains() {
	nums := []int{1, 2, 3}
	fmt.Println(gox.Contains(nums, 2))
	fmt.Println(gox.Contains(nums, 5))
	// Output:
	// true
	// false
}

func ExampleChunk() {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	chunks := gox.Chunk(nums, 3)
	fmt.Println(chunks)
	// Output:
	// [[1 2 3] [4 5 6] [7]]
}

func ExampleFlatten() {
	nested := [][]int{{1, 2}, {3, 4}, {5}}
	flat := gox.Flatten(nested)
	fmt.Println(flat)
	// Output:
	// [1 2 3 4 5]
}

func ExampleOSome() {
	opt := gox.OSome(42)
	fmt.Println(opt.IsSome())
	fmt.Println(opt.OrElse(0))
	// Output:
	// true
	// 42
}

func ExampleONone() {
	opt := gox.ONone[int]()
	fmt.Println(opt.IsNone())
	fmt.Println(opt.OrElse(100))
	// Output:
	// true
	// 100
}

func ExampleOptional_Map() {
	result := gox.OSome(10).
		Map(func(n int) int { return n * 2 }).
		OrElse(0)
	fmt.Println(result)
	// Output:
	// 20
}

func ExampleROk() {
	r := gox.ROk(42)
	fmt.Println(r.IsOk())
	fmt.Println(r.Unwrap())
	// Output:
	// true
	// 42
}

func ExampleRErr() {
	r := gox.RErr[int](errors.New("failed"))
	fmt.Println(r.IsErr())
	fmt.Println(r.UnwrapOr(0))
	// Output:
	// true
	// 0
}

func ExampleTry() {
	r := gox.Try(func() (int, error) { return 42, nil })
	fmt.Println(r.Unwrap())
	// Output:
	// 42
}

func ExampleIntersect() {
	a := []int{1, 2, 3, 4}
	b := []int{3, 4, 5, 6}
	fmt.Println(gox.Intersect(a, b))
	// Output:
	// [3 4]
}

func ExampleUnion() {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	fmt.Println(gox.Union(a, b))
	// Output:
	// [1 2 3 4 5]
}

func ExampleDifference() {
	a := []int{1, 2, 3, 4}
	b := []int{3, 4, 5}
	fmt.Println(gox.Difference(a, b))
	// Output:
	// [1 2]
}

func ExampleMultiError() {
	m := gox.NewMultiError()
	m.Add(errors.New("error 1"))
	m.Add(errors.New("error 2"))
	fmt.Println(m.Len())
	fmt.Println(m.Error())
	// Output:
	// 2
	// 2 errors: error 1; error 2
}

func ExamplePtr() {
	p := gox.Ptr(42)
	fmt.Println(*p)
	// Output:
	// 42
}

func ExampleVal() {
	p := gox.Ptr(42)
	fmt.Println(gox.Val(p, 0))
	fmt.Println(gox.Val[int](nil, 100))
	// Output:
	// 42
	// 100
}

func ExampleIf() {
	result := gox.If(true, "yes", "no")
	fmt.Println(result)
	// Output:
	// yes
}

func ExampleCoalesce() {
	result := gox.Coalesce("", "", "hello")
	fmt.Println(result)
	// Output:
	// hello
}
