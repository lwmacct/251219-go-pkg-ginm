package gox

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestROk_CreatesSuccessResult(t *testing.T) {
	r := ROk(42)
	assert.True(t, r.IsOk())
	assert.False(t, r.IsErr())
}

func TestRErr_CreatesErrorResult(t *testing.T) {
	r := RErr[int](assert.AnError)
	assert.True(t, r.IsErr())
	assert.False(t, r.IsOk())
}

func TestTry_ReturnsOkOnSuccess(t *testing.T) {
	r := Try(func() (int, error) { return 42, nil })
	require.True(t, r.IsOk())
	assert.Equal(t, 42, r.Unwrap())
}

func TestTry_ReturnsErrOnError(t *testing.T) {
	r := Try(func() (int, error) { return 0, assert.AnError })
	assert.True(t, r.IsErr())
}

func TestTryE_ReturnsOkOnNoError(t *testing.T) {
	r := TryE(func() error { return nil })
	assert.True(t, r.IsOk())
}

func TestTryE_ReturnsErrOnError(t *testing.T) {
	r := TryE(func() error { return assert.AnError })
	assert.True(t, r.IsErr())
}

func TestResult_Unwrap_ReturnsValue(t *testing.T) {
	r := ROk(42)
	assert.Equal(t, 42, r.Unwrap())
}

func TestResult_Unwrap_PanicsOnErr(t *testing.T) {
	r := RErr[int](assert.AnError)
	assert.Panics(t, func() { r.Unwrap() })
}

func TestResult_UnwrapErr_ReturnsError(t *testing.T) {
	r := RErr[int](assert.AnError)
	assert.Equal(t, assert.AnError, r.UnwrapErr())
}

func TestResult_UnwrapErr_PanicsOnOk(t *testing.T) {
	r := ROk(42)
	assert.Panics(t, func() { _ = r.UnwrapErr() })
}

func TestResult_UnwrapOr_ReturnsValueOnOk(t *testing.T) {
	r := ROk(42)
	assert.Equal(t, 42, r.UnwrapOr(100))
}

func TestResult_UnwrapOr_ReturnsDefaultOnErr(t *testing.T) {
	r := RErr[int](assert.AnError)
	assert.Equal(t, 100, r.UnwrapOr(100))
}

func TestResult_UnwrapOrElse_CallsFnOnlyOnErr(t *testing.T) {
	called := false
	r := ROk(42)
	result := r.UnwrapOrElse(func() int { called = true; return 100 })
	assert.Equal(t, 42, result)
	assert.False(t, called)
}

func TestResult_UnwrapOrElse_CallsFnOnErr(t *testing.T) {
	r := RErr[int](assert.AnError)
	result := r.UnwrapOrElse(func() int { return 100 })
	assert.Equal(t, 100, result)
}

func TestResult_Get_ReturnsValueAndTrue(t *testing.T) {
	r := ROk(42)
	val, ok := r.Get()
	assert.True(t, ok)
	assert.Equal(t, 42, val)
}

func TestResult_Get_ReturnsZeroAndFalse(t *testing.T) {
	r := RErr[int](assert.AnError)
	val, ok := r.Get()
	assert.False(t, ok)
	assert.Equal(t, 0, val)
}

func TestResult_GetWithError_ReturnsValueAndNil(t *testing.T) {
	r := ROk(42)
	val, err := r.GetWithError()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestResult_GetWithError_ReturnsZeroAndError(t *testing.T) {
	r := RErr[int](assert.AnError)
	val, err := r.GetWithError()
	assert.Error(t, err)
	assert.Equal(t, 0, val)
}

func TestResult_Map_TransformsValue(t *testing.T) {
	r := ROk(21)
	result := r.Map(func(n int) int { return n * 2 })
	assert.Equal(t, 42, result.Unwrap())
}

func TestResult_Map_PreservesError(t *testing.T) {
	r := RErr[int](assert.AnError)
	result := r.Map(func(n int) int { return n * 2 })
	assert.True(t, result.IsErr())
}

func TestMapTo_TransformsType(t *testing.T) {
	r := ROk(42)
	result := MapTo(r, func(n int) string { return "value" })
	assert.Equal(t, "value", result.Unwrap())
}

func TestResult_MapErr_TransformsError(t *testing.T) {
	r := RErr[int](errors.New("original"))
	result := r.MapErr(func(err error) error { return errors.New("wrapped: " + err.Error()) })
	assert.Equal(t, "wrapped: original", result.Error().Error())
}

func TestResult_MapErr_PreservesOk(t *testing.T) {
	r := ROk(42)
	result := r.MapErr(func(err error) error { return errors.New("wrapped") })
	assert.True(t, result.IsOk())
	assert.Equal(t, 42, result.Unwrap())
}

func TestResult_AndThen_ChainsOperations(t *testing.T) {
	r := ROk(21)
	result := r.AndThen(func(n int) Result[int] { return ROk(n * 2) })
	assert.Equal(t, 42, result.Unwrap())
}

func TestResult_AndThen_StopsOnError(t *testing.T) {
	r := RErr[int](assert.AnError)
	called := false
	result := r.AndThen(func(n int) Result[int] { called = true; return ROk(n * 2) })
	assert.True(t, result.IsErr())
	assert.False(t, called)
}

func TestAndThenTo_ChainsWithTypeChange(t *testing.T) {
	r := ROk(42)
	result := AndThenTo(r, func(n int) Result[string] { return ROk("number") })
	assert.Equal(t, "number", result.Unwrap())
}

func TestResult_OrElse_ReturnsOriginalOnOk(t *testing.T) {
	r := ROk(42)
	called := false
	result := r.OrElse(func(err error) Result[int] { called = true; return ROk(100) })
	assert.Equal(t, 42, result.Unwrap())
	assert.False(t, called)
}

func TestResult_OrElse_CallsFnOnErr(t *testing.T) {
	r := RErr[int](assert.AnError)
	result := r.OrElse(func(err error) Result[int] { return ROk(100) })
	assert.Equal(t, 100, result.Unwrap())
}

func TestResult_Inspect_CallsFnOnOk(t *testing.T) {
	called := false
	r := ROk(42)
	result := r.Inspect(func(n int) { called = true })
	assert.True(t, called)
	assert.True(t, result.IsOk())
}

func TestResult_Inspect_DoesNotCallFnOnErr(t *testing.T) {
	called := false
	r := RErr[int](assert.AnError)
	result := r.Inspect(func(n int) { called = true })
	assert.False(t, called)
	assert.True(t, result.IsErr())
}

func TestResult_InspectErr_CallsFnOnErr(t *testing.T) {
	called := false
	r := RErr[int](assert.AnError)
	result := r.InspectErr(func(err error) { called = true })
	assert.True(t, called)
	assert.True(t, result.IsErr())
}

func TestMatch_CallsCorrectFunction(t *testing.T) {
	ok := ROk(42)
	result := Match(ok, func(n int) string { return "ok" }, func(err error) string { return "err" })
	assert.Equal(t, "ok", result)

	err := RErr[int](assert.AnError)
	result = Match(err, func(n int) string { return "ok" }, func(err error) string { return "err" })
	assert.Equal(t, "err", result)
}

func TestCollect_CollectsAllValues(t *testing.T) {
	results := []Result[int]{ROk(1), ROk(2), ROk(3)}
	collected := Collect(results)
	require.True(t, collected.IsOk())
	assert.Equal(t, []int{1, 2, 3}, collected.Unwrap())
}

func TestCollect_ReturnsFirstError(t *testing.T) {
	results := []Result[int]{ROk(1), RErr[int](assert.AnError), ROk(3)}
	collected := Collect(results)
	assert.True(t, collected.IsErr())
}

func TestFlattenResult_FlattensNestedResult(t *testing.T) {
	nested := ROk(ROk(42))
	result := FlattenResult(nested)
	assert.Equal(t, 42, result.Unwrap())
}

func TestFlattenResult_PropagatesOuterError(t *testing.T) {
	nested := RErr[Result[int]](assert.AnError)
	result := FlattenResult(nested)
	assert.True(t, result.IsErr())
}
