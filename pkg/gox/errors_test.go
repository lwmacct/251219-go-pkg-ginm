package gox

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiError_NewMultiError_CreatesEmpty(t *testing.T) {
	m := NewMultiError()
	require.NotNil(t, m)
	assert.False(t, m.HasErrors())
	assert.Equal(t, 0, m.Len())
}

func TestMultiError_Add_AddsError(t *testing.T) {
	m := NewMultiError()
	m.Add(errors.New("error 1"))
	assert.True(t, m.HasErrors())
	assert.Equal(t, 1, m.Len())
}

func TestMultiError_Add_IgnoresNil(t *testing.T) {
	m := NewMultiError()
	m.Add(nil)
	assert.False(t, m.HasErrors())
}

func TestMultiError_AddAll_AddsMultiple(t *testing.T) {
	m := NewMultiError()
	m.AddAll(errors.New("e1"), nil, errors.New("e2"), nil)
	assert.Equal(t, 2, m.Len())
}

func TestMultiError_Errors_ReturnsAll(t *testing.T) {
	m := NewMultiError()
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	m.AddAll(e1, e2)
	errs := m.Errors()
	assert.Len(t, errs, 2)
	assert.Contains(t, errs, e1)
	assert.Contains(t, errs, e2)
}

func TestMultiError_Error_ReturnsEmptyForNoErrors(t *testing.T) {
	m := NewMultiError()
	assert.Equal(t, "", m.Error())
}

func TestMultiError_Error_ReturnsSingleError(t *testing.T) {
	m := NewMultiError()
	m.Add(errors.New("single error"))
	assert.Equal(t, "single error", m.Error())
}

func TestMultiError_Error_CombinesMultipleErrors(t *testing.T) {
	m := NewMultiError()
	m.Add(errors.New("e1"))
	m.Add(errors.New("e2"))
	m.Add(errors.New("e3"))
	result := m.Error()
	assert.Equal(t, "3 errors: e1; e2; e3", result)
}

func TestMultiError_ErrorOrNil_ReturnsNilForNoErrors(t *testing.T) {
	m := NewMultiError()
	assert.NoError(t, m.ErrorOrNil())
}

func TestMultiError_ErrorOrNil_ReturnsSelfForErrors(t *testing.T) {
	m := NewMultiError()
	m.Add(errors.New("error"))
	err := m.ErrorOrNil()
	require.Error(t, err)
	assert.Equal(t, m, err)
}

func TestMultiError_First_ReturnsNilForNoErrors(t *testing.T) {
	m := NewMultiError()
	assert.NoError(t, m.First())
}

func TestMultiError_First_ReturnsFirstError(t *testing.T) {
	m := NewMultiError()
	e1 := errors.New("first")
	e2 := errors.New("second")
	m.AddAll(e1, e2)
	assert.Equal(t, e1, m.First())
}

func TestMultiError_Unwrap_ReturnsErrorList(t *testing.T) {
	m := NewMultiError()
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	m.AddAll(e1, e2)
	unwrapped := m.Unwrap()
	assert.Len(t, unwrapped, 2)
}

func TestMultiError_ImplementsErrorInterface(t *testing.T) {
	m := NewMultiError()
	m.Add(errors.New("test"))
	var err error = m
	assert.Error(t, err)
}

func TestMultiError_WorksWithErrorsIs(t *testing.T) {
	target := errors.New("target error")
	m := NewMultiError()
	m.Add(errors.New("other"))
	m.Add(target)
	assert.ErrorIs(t, m, target)
}
