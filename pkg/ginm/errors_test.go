package ginm

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(http.StatusBadRequest, 400, "bad request")
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus)
	assert.Equal(t, 400, err.Code)
	assert.Equal(t, "bad request", err.Message)
	assert.Nil(t, err.Err)
}

func TestWrapAPIError(t *testing.T) {
	inner := errors.New("inner error")
	err := WrapAPIError(http.StatusInternalServerError, 500, "wrapper", inner)

	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
	assert.Equal(t, "wrapper", err.Message)
	assert.Equal(t, inner, err.Err)
}

func TestAPIError_Error(t *testing.T) {
	err := NewAPIError(400, 400, "bad request")
	assert.Equal(t, "bad request", err.Error())

	wrapped := WrapAPIError(500, 500, "wrapper", errors.New("cause"))
	assert.Equal(t, "wrapper: cause", wrapped.Error())
}

func TestAPIError_Unwrap(t *testing.T) {
	inner := errors.New("inner")
	err := WrapAPIError(500, 500, "outer", inner)
	assert.Equal(t, inner, err.Unwrap())
	assert.True(t, errors.Is(err, inner))
}

func TestErrBadRequest(t *testing.T) {
	err := ErrBadRequest("invalid input")
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus)
	assert.Equal(t, "invalid input", err.Message)
}

func TestErrUnauthorized(t *testing.T) {
	err := ErrUnauthorized("token expired")
	assert.Equal(t, http.StatusUnauthorized, err.HTTPStatus)
}

func TestErrForbidden(t *testing.T) {
	err := ErrForbidden("access denied")
	assert.Equal(t, http.StatusForbidden, err.HTTPStatus)
}

func TestErrNotFound(t *testing.T) {
	err := ErrNotFound("user not found")
	assert.Equal(t, http.StatusNotFound, err.HTTPStatus)
}

func TestErrConflict(t *testing.T) {
	err := ErrConflict("already exists")
	assert.Equal(t, http.StatusConflict, err.HTTPStatus)
}

func TestErrInternal(t *testing.T) {
	err := ErrInternal("database error")
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
}

func TestErrInternalWrap(t *testing.T) {
	cause := errors.New("db connection failed")
	err := ErrInternalWrap("database error", cause)
	assert.Equal(t, http.StatusInternalServerError, err.HTTPStatus)
	assert.True(t, errors.Is(err, cause))
}

func TestErrNotImplemented(t *testing.T) {
	err := ErrNotImplemented("PATCH")
	assert.Equal(t, http.StatusNotImplemented, err.HTTPStatus)
	assert.Contains(t, err.Message, "PATCH")
}

func TestNewBindError(t *testing.T) {
	cause := errors.New("validation failed")
	err := NewBindError("body", cause)

	assert.Equal(t, "body", err.Source)
	assert.Equal(t, cause, err.Err)
}

func TestBindError_Error(t *testing.T) {
	err := NewBindError("json", errors.New("invalid syntax"))
	assert.Contains(t, err.Error(), "json")
	assert.Contains(t, err.Error(), "invalid syntax")
}

func TestBindError_Unwrap(t *testing.T) {
	cause := errors.New("cause")
	err := NewBindError("query", cause)
	assert.Equal(t, cause, err.Unwrap())
	assert.True(t, errors.Is(err, cause))
}

func TestValidationErrors_Add(t *testing.T) {
	ve := &ValidationErrors{}
	ve.Add("email", "invalid format")
	ve.Add("age", "must be positive")

	assert.Len(t, ve.Errors, 2)
	assert.Equal(t, "email", ve.Errors[0].Field)
	assert.Equal(t, "invalid format", ve.Errors[0].Message)
}

func TestValidationErrors_HasErrors(t *testing.T) {
	ve := &ValidationErrors{}
	assert.False(t, ve.HasErrors())

	ve.Add("field", "error")
	assert.True(t, ve.HasErrors())
}

func TestValidationErrors_Error(t *testing.T) {
	ve := &ValidationErrors{}
	ve.Add("a", "error1")
	ve.Add("b", "error2")

	assert.Contains(t, ve.Error(), "2 errors")
}
