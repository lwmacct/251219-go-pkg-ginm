package gox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOSome_CreatesValidOptional(t *testing.T) {
	opt := OSome(42)
	assert.True(t, opt.IsSome())
	assert.False(t, opt.IsNone())
}

func TestONone_CreatesEmptyOptional(t *testing.T) {
	opt := ONone[int]()
	assert.True(t, opt.IsNone())
	assert.False(t, opt.IsSome())
}

func TestOFromPtr_ReturnsNoneForNil(t *testing.T) {
	opt := OFromPtr[int](nil)
	assert.True(t, opt.IsNone())
}

func TestOFromPtr_ReturnsSomeForValue(t *testing.T) {
	p := Ptr(42)
	opt := OFromPtr(p)
	assert.True(t, opt.IsSome())
	assert.Equal(t, 42, opt.MustGet())
}

func TestOFromOk_ReturnsNoneWhenNotOk(t *testing.T) {
	opt := OFromOk(0, false)
	assert.True(t, opt.IsNone())
}

func TestOFromOk_ReturnsSomeWhenOk(t *testing.T) {
	opt := OFromOk(42, true)
	assert.True(t, opt.IsSome())
	assert.Equal(t, 42, opt.MustGet())
}

func TestOFromErr_ReturnsNoneOnError(t *testing.T) {
	opt := OFromErr(0, assert.AnError)
	assert.True(t, opt.IsNone())
}

func TestOFromErr_ReturnsSomeOnNoError(t *testing.T) {
	opt := OFromErr(42, nil)
	assert.True(t, opt.IsSome())
	assert.Equal(t, 42, opt.MustGet())
}

func TestOptional_Get_ReturnsValueAndTrue(t *testing.T) {
	opt := OSome(42)
	val, ok := opt.Get()
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestOptional_Get_ReturnsZeroAndFalse(t *testing.T) {
	opt := ONone[int]()
	val, ok := opt.Get()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

func TestOptional_MustGet_PanicsOnNone(t *testing.T) {
	opt := ONone[int]()
	assert.Panics(t, func() { opt.MustGet() })
}

func TestOptional_OrElse_ReturnsValueWhenSome(t *testing.T) {
	opt := OSome(42)
	assert.Equal(t, 42, opt.OrElse(100))
}

func TestOptional_OrElse_ReturnsDefaultWhenNone(t *testing.T) {
	opt := ONone[int]()
	assert.Equal(t, 100, opt.OrElse(100))
}

func TestOptional_OrElseFn_CallsFnOnlyWhenNone(t *testing.T) {
	called := false
	opt := OSome(42)
	result := opt.OrElseFn(func() int { called = true; return 100 })
	assert.Equal(t, 42, result)
	assert.False(t, called)
}

func TestOptional_OrElseFn_CallsFnWhenNone(t *testing.T) {
	opt := ONone[int]()
	result := opt.OrElseFn(func() int { return 100 })
	assert.Equal(t, 100, result)
}

func TestOptional_ToPtr_ReturnsPointerWhenSome(t *testing.T) {
	opt := OSome(42)
	p := opt.ToPtr()
	require.NotNil(t, p)
	assert.Equal(t, 42, *p)
}

func TestOptional_ToPtr_ReturnsNilWhenNone(t *testing.T) {
	opt := ONone[int]()
	assert.Nil(t, opt.ToPtr())
}

func TestOptional_Map_TransformsValue(t *testing.T) {
	opt := OSome(21)
	result := opt.Map(func(n int) int { return n * 2 })
	assert.Equal(t, 42, result.MustGet())
}

func TestOptional_Map_ReturnsNoneForNone(t *testing.T) {
	opt := ONone[int]()
	result := opt.Map(func(n int) int { return n * 2 })
	assert.True(t, result.IsNone())
}

func TestOMapTo_TransformsType(t *testing.T) {
	opt := OSome(42)
	result := OMapTo(opt, func(n int) string { return "value" })
	assert.Equal(t, "value", result.MustGet())
}

func TestOptional_Filter_KeepsMatchingValue(t *testing.T) {
	opt := OSome(42)
	result := opt.Filter(func(n int) bool { return n%2 == 0 })
	assert.True(t, result.IsSome())
}

func TestOptional_Filter_RemovesNonMatchingValue(t *testing.T) {
	opt := OSome(41)
	result := opt.Filter(func(n int) bool { return n%2 == 0 })
	assert.True(t, result.IsNone())
}

func TestOptional_Or_ReturnsCurrentWhenSome(t *testing.T) {
	a := OSome(1)
	b := OSome(2)
	result := a.Or(b)
	assert.Equal(t, 1, result.MustGet())
}

func TestOptional_Or_ReturnsAlternativeWhenNone(t *testing.T) {
	a := ONone[int]()
	b := OSome(2)
	result := a.Or(b)
	assert.Equal(t, 2, result.MustGet())
}

func TestOptional_And_ReturnsOtherWhenBothSome(t *testing.T) {
	a := OSome(1)
	b := OSome(2)
	result := a.And(b)
	assert.Equal(t, 2, result.MustGet())
}

func TestOptional_And_ReturnsNoneWhenFirstIsNone(t *testing.T) {
	a := ONone[int]()
	b := OSome(2)
	result := a.And(b)
	assert.True(t, result.IsNone())
}

func TestOptional_Xor_ReturnsSomeWhenOnlyOneHasValue(t *testing.T) {
	a := OSome(1)
	b := ONone[int]()
	result := a.Xor(b)
	assert.Equal(t, 1, result.MustGet())

	result = b.Xor(a)
	assert.Equal(t, 1, result.MustGet())
}

func TestOptional_Xor_ReturnsNoneWhenBothHaveValues(t *testing.T) {
	a := OSome(1)
	b := OSome(2)
	result := a.Xor(b)
	assert.True(t, result.IsNone())
}

func TestOptional_Inspect_CallsFnOnSome(t *testing.T) {
	called := false
	opt := OSome(42)
	result := opt.Inspect(func(n int) { called = true })
	assert.True(t, called)
	assert.True(t, result.IsSome())
}

func TestOptional_Inspect_DoesNotCallFnOnNone(t *testing.T) {
	called := false
	opt := ONone[int]()
	result := opt.Inspect(func(n int) { called = true })
	assert.False(t, called)
	assert.True(t, result.IsNone())
}

func TestOMatch_CallsCorrectFunction(t *testing.T) {
	some := OSome(42)
	result := OMatch(some, func(n int) string { return "has value" }, func() string { return "empty" })
	assert.Equal(t, "has value", result)

	none := ONone[int]()
	result = OMatch(none, func(n int) string { return "has value" }, func() string { return "empty" })
	assert.Equal(t, "empty", result)
}

func TestOptional_ToResult_ConvertsCorrectly(t *testing.T) {
	some := OSome(42)
	result := some.ToResult(assert.AnError)
	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())

	none := ONone[int]()
	result = none.ToResult(assert.AnError)
	assert.True(t, result.IsErr())
}

func TestOZip_CombinesTwoOptionals(t *testing.T) {
	a := OSome(1)
	b := OSome("hello")
	result := OZip(a, b)
	require.True(t, result.IsSome())
	pair := result.MustGet()
	assert.Equal(t, 1, pair.First)
	assert.Equal(t, "hello", pair.Second)
}

func TestOZip_ReturnsNoneIfAnyIsNone(t *testing.T) {
	a := OSome(1)
	b := ONone[string]()
	result := OZip(a, b)
	assert.True(t, result.IsNone())
}
