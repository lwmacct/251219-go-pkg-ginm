package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap_TransformsElements(t *testing.T) {
	nums := []int{1, 2, 3}
	doubled := Map(nums, func(n int) int { return n * 2 })
	assert.Equal(t, []int{2, 4, 6}, doubled)
}

func TestMap_ReturnsNilForNilInput(t *testing.T) {
	var nums []int
	result := Map(nums, func(n int) int { return n * 2 })
	assert.Nil(t, result)
}

func TestMap_TypeConversion(t *testing.T) {
	nums := []int{1, 2, 3}
	strs := Map(nums, func(n int) string { return string(rune('a' + n - 1)) })
	assert.Equal(t, []string{"a", "b", "c"}, strs)
}

func TestFilter_SelectsMatchingElements(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	assert.Equal(t, []int{2, 4}, evens)
}

func TestFilter_ReturnsEmptyForNoMatch(t *testing.T) {
	nums := []int{1, 3, 5}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	assert.Empty(t, evens)
}

func TestFilter_ReturnsNilForNilInput(t *testing.T) {
	var nums []int
	result := Filter(nums, func(n int) bool { return true })
	assert.Nil(t, result)
}

func TestReduce_SumsNumbers(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	assert.Equal(t, 15, sum)
}

func TestReduce_EmptySliceReturnsInit(t *testing.T) {
	var nums []int
	sum := Reduce(nums, 100, func(acc, n int) int { return acc + n })
	assert.Equal(t, 100, sum)
}

func TestFind_ReturnsFirstMatch(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	val, ok := Find(nums, func(n int) bool { return n > 2 })
	require.True(t, ok)
	assert.Equal(t, 3, val)
}

func TestFind_ReturnsFalseForNoMatch(t *testing.T) {
	nums := []int{1, 2, 3}
	_, ok := Find(nums, func(n int) bool { return n > 10 })
	assert.False(t, ok)
}

func TestFindIndex_ReturnsCorrectIndex(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	idx := FindIndex(nums, func(n int) bool { return n == 3 })
	assert.Equal(t, 2, idx)
}

func TestFindIndex_ReturnsNegativeOneForNoMatch(t *testing.T) {
	nums := []int{1, 2, 3}
	idx := FindIndex(nums, func(n int) bool { return n == 10 })
	assert.Equal(t, -1, idx)
}

func TestEvery_ReturnsTrueWhenAllMatch(t *testing.T) {
	nums := []int{2, 4, 6}
	result := Every(nums, func(n int) bool { return n%2 == 0 })
	assert.True(t, result)
}

func TestEvery_ReturnsFalseWhenAnyFails(t *testing.T) {
	nums := []int{2, 3, 4}
	result := Every(nums, func(n int) bool { return n%2 == 0 })
	assert.False(t, result)
}

func TestEvery_ReturnsTrueForEmptySlice(t *testing.T) {
	var nums []int
	result := Every(nums, func(n int) bool { return false })
	assert.True(t, result)
}

func TestSome_ReturnsTrueWhenAnyMatch(t *testing.T) {
	nums := []int{1, 2, 3}
	result := Some(nums, func(n int) bool { return n%2 == 0 })
	assert.True(t, result)
}

func TestSome_ReturnsFalseWhenNoneMatch(t *testing.T) {
	nums := []int{1, 3, 5}
	result := Some(nums, func(n int) bool { return n%2 == 0 })
	assert.False(t, result)
}

func TestContains_ReturnsTrueWhenPresent(t *testing.T) {
	nums := []int{1, 2, 3}
	assert.True(t, Contains(nums, 2))
}

func TestContains_ReturnsFalseWhenAbsent(t *testing.T) {
	nums := []int{1, 2, 3}
	assert.False(t, Contains(nums, 5))
}

func TestUnique_RemovesDuplicates(t *testing.T) {
	nums := []int{1, 2, 2, 3, 1, 4}
	result := Unique(nums)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestUnique_ReturnsNilForNilInput(t *testing.T) {
	var nums []int
	result := Unique(nums)
	assert.Nil(t, result)
}

func TestGroupBy_GroupsCorrectly(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	groups := GroupBy(nums, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})
	assert.Equal(t, []int{2, 4}, groups["even"])
	assert.Equal(t, []int{1, 3, 5}, groups["odd"])
}

func TestChunk_SplitsIntoCorrectSizes(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	chunks := Chunk(nums, 2)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, chunks)
}

func TestChunk_ReturnsNilForInvalidSize(t *testing.T) {
	nums := []int{1, 2, 3}
	assert.Nil(t, Chunk(nums, 0))
	assert.Nil(t, Chunk(nums, -1))
}

func TestChunk_ReturnsNilForEmptySlice(t *testing.T) {
	var nums []int
	result := Chunk(nums, 2)
	assert.Nil(t, result)
}

func TestFlatten_FlattensNestedSlices(t *testing.T) {
	nested := [][]int{{1, 2}, {3, 4}, {5}}
	result := Flatten(nested)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestFlatten_ReturnsNilForNilInput(t *testing.T) {
	var nested [][]int
	result := Flatten(nested)
	assert.Nil(t, result)
}

func TestFirst_ReturnsFirstElement(t *testing.T) {
	nums := []int{1, 2, 3}
	val, ok := First(nums)
	require.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestFirst_ReturnsFalseForEmpty(t *testing.T) {
	var nums []int
	_, ok := First(nums)
	assert.False(t, ok)
}

func TestLast_ReturnsLastElement(t *testing.T) {
	nums := []int{1, 2, 3}
	val, ok := Last(nums)
	require.True(t, ok)
	assert.Equal(t, 3, val)
}

func TestLast_ReturnsFalseForEmpty(t *testing.T) {
	var nums []int
	_, ok := Last(nums)
	assert.False(t, ok)
}

func TestReverse_ReversesSlice(t *testing.T) {
	nums := []int{1, 2, 3}
	result := Reverse(nums)
	assert.Equal(t, []int{3, 2, 1}, result)
}

func TestReverse_ReturnsNilForNilInput(t *testing.T) {
	var nums []int
	result := Reverse(nums)
	assert.Nil(t, result)
}

func TestPtr_ReturnsPointer(t *testing.T) {
	p := Ptr(42)
	require.NotNil(t, p)
	assert.Equal(t, 42, *p)
}

func TestVal_ReturnsValueOrDefault(t *testing.T) {
	p := Ptr(42)
	assert.Equal(t, 42, Val(p, 0))
	assert.Equal(t, 100, Val[int](nil, 100))
}

func TestValOrZero_ReturnsValueOrZero(t *testing.T) {
	p := Ptr(42)
	assert.Equal(t, 42, ValOrZero(p))
	assert.Equal(t, 0, ValOrZero[int](nil))
}

func TestCoalesce_ReturnsFirstNonZero(t *testing.T) {
	result := Coalesce("", "", "hello", "world")
	assert.Equal(t, "hello", result)
}

func TestCoalesce_ReturnsZeroIfAllZero(t *testing.T) {
	result := Coalesce("", "", "")
	assert.Empty(t, result)
}

func TestCoalescePtr_ReturnsFirstNonNil(t *testing.T) {
	a := Ptr(1)
	b := Ptr(2)
	val, ok := CoalescePtr(nil, nil, a, b)
	require.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestCoalescePtr_ReturnsFalseIfAllNil(t *testing.T) {
	_, ok := CoalescePtr[int](nil, nil)
	assert.False(t, ok)
}

func TestIf_ReturnsTrueValue(t *testing.T) {
	result := If(true, "yes", "no")
	assert.Equal(t, "yes", result)
}

func TestIf_ReturnsFalseValue(t *testing.T) {
	result := If(false, "yes", "no")
	assert.Equal(t, "no", result)
}

func TestIfFn_CallsCorrectFunction(t *testing.T) {
	trueCalled := false
	falseCalled := false
	result := IfFn(true,
		func() string { trueCalled = true; return "yes" },
		func() string { falseCalled = true; return "no" },
	)
	assert.Equal(t, "yes", result)
	assert.True(t, trueCalled)
	assert.False(t, falseCalled)
}

func TestKeys_ReturnsMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	keys := Keys(m)
	assert.Len(t, keys, 2)
	assert.True(t, Contains(keys, "a"))
	assert.True(t, Contains(keys, "b"))
}

func TestKeys_ReturnsNilForNilMap(t *testing.T) {
	var m map[string]int
	result := Keys(m)
	assert.Nil(t, result)
}

func TestValues_ReturnsMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	values := Values(m)
	assert.Len(t, values, 2)
	assert.True(t, Contains(values, 1))
	assert.True(t, Contains(values, 2))
}

func TestValues_ReturnsNilForNilMap(t *testing.T) {
	var m map[string]int
	result := Values(m)
	assert.Nil(t, result)
}
