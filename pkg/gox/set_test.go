package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersect_ReturnsCommonElements(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{3, 4, 5, 6}
	result := Intersect(a, b)
	assert.Equal(t, []int{3, 4}, result)
}

func TestIntersect_ReturnsEmptyForNoCommon(t *testing.T) {
	a := []int{1, 2}
	b := []int{3, 4}
	result := Intersect(a, b)
	assert.Empty(t, result)
}

func TestIntersect_ReturnsEmptyForEmptyInput(t *testing.T) {
	a := []int{1, 2}
	var b []int
	assert.Empty(t, Intersect(a, b))
	assert.Empty(t, Intersect(b, a))
}

func TestIntersect_RemovesDuplicates(t *testing.T) {
	a := []int{1, 1, 2, 2}
	b := []int{1, 2, 2, 3}
	result := Intersect(a, b)
	assert.Equal(t, []int{1, 2}, result)
}

func TestUnion_CombinesElements(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	result := Union(a, b)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestUnion_RemovesDuplicates(t *testing.T) {
	a := []int{1, 1, 2}
	b := []int{2, 3, 3}
	result := Union(a, b)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDifference_ReturnsElementsOnlyInA(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{3, 4, 5}
	result := Difference(a, b)
	assert.Equal(t, []int{1, 2}, result)
}

func TestDifference_ReturnsEmptyIfAIsEmpty(t *testing.T) {
	var a []int
	b := []int{1, 2}
	result := Difference(a, b)
	assert.Empty(t, result)
}

func TestDifference_ReturnsCopyOfAIfBIsEmpty(t *testing.T) {
	a := []int{1, 2, 3}
	var b []int
	result := Difference(a, b)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestSymmetricDifference_ReturnsElementsInEitherButNotBoth(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}
	result := SymmetricDifference(a, b)
	assert.ElementsMatch(t, []int{1, 4}, result)
}

func TestPartition_SplitsCorrectly(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	evens, odds := Partition(nums, func(n int) bool { return n%2 == 0 })
	assert.Equal(t, []int{2, 4}, evens)
	assert.Equal(t, []int{1, 3, 5}, odds)
}

func TestIsSubset_ReturnsTrueWhenSubset(t *testing.T) {
	a := []int{1, 2}
	b := []int{1, 2, 3, 4}
	assert.True(t, IsSubset(a, b))
}

func TestIsSubset_ReturnsFalseWhenNotSubset(t *testing.T) {
	a := []int{1, 2, 5}
	b := []int{1, 2, 3, 4}
	assert.False(t, IsSubset(a, b))
}

func TestIsSubset_EmptyIsSubsetOfAny(t *testing.T) {
	var a []int
	b := []int{1, 2, 3}
	assert.True(t, IsSubset(a, b))
}

func TestIsSubset_NonEmptyIsNotSubsetOfEmpty(t *testing.T) {
	a := []int{1}
	var b []int
	assert.False(t, IsSubset(a, b))
}

func TestIsSuperset_ReturnsTrueWhenSuperset(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{1, 2}
	assert.True(t, IsSuperset(a, b))
}

func TestIsDisjoint_ReturnsTrueWhenNoCommon(t *testing.T) {
	a := []int{1, 2}
	b := []int{3, 4}
	assert.True(t, IsDisjoint(a, b))
}

func TestIsDisjoint_ReturnsFalseWhenCommon(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	assert.False(t, IsDisjoint(a, b))
}

func TestIsDisjoint_ReturnsTrueForEmptyInput(t *testing.T) {
	a := []int{1, 2}
	var b []int
	assert.True(t, IsDisjoint(a, b))
	assert.True(t, IsDisjoint(b, a))
}

func TestCount_CountsMatchingElements(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	count := Count(nums, func(n int) bool { return n%2 == 0 })
	assert.Equal(t, 3, count)
}

func TestCountBy_CountsByKey(t *testing.T) {
	words := []string{"a", "bb", "ccc", "dd", "eee"}
	counts := CountBy(words, func(s string) int { return len(s) })
	assert.Equal(t, 1, counts[1])
	assert.Equal(t, 2, counts[2])
	assert.Equal(t, 2, counts[3])
}

func TestZip_CombinesTwoSlices(t *testing.T) {
	a := []int{1, 2, 3}
	b := []string{"a", "b", "c"}
	result := Zip(a, b)
	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0].First)
	assert.Equal(t, "a", result[0].Second)
}

func TestZip_UsesMinLength(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []string{"a", "b"}
	result := Zip(a, b)
	assert.Len(t, result, 2)
}

func TestUnzip_SplitsPairs(t *testing.T) {
	pairs := []struct {
		First  int
		Second string
	}{
		{1, "a"},
		{2, "b"},
	}
	a, b := Unzip(pairs)
	assert.Equal(t, []int{1, 2}, a)
	assert.Equal(t, []string{"a", "b"}, b)
}
