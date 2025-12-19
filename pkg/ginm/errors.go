package ginm

import (
	"fmt"
	"net/http"
)

// APIError 表示结构化的 API 错误。
type APIError struct {
	Err        error
	Message    string
	HTTPStatus int
	Code       int
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// NewAPIError 创建一个新的 API 错误。
func NewAPIError(httpStatus, code int, message string) *APIError {
	return &APIError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
	}
}

// WrapAPIError 创建一个包装其他错误的 API 错误。
func WrapAPIError(httpStatus, code int, message string, err error) *APIError {
	return &APIError{
		HTTPStatus: httpStatus,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

// 常用错误构造函数

// ErrBadRequest 创建 400 错误请求错误。
func ErrBadRequest(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, http.StatusBadRequest, message)
}

// ErrUnauthorized 创建 401 未授权错误。
func ErrUnauthorized(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, http.StatusUnauthorized, message)
}

// ErrForbidden 创建 403 禁止访问错误。
func ErrForbidden(message string) *APIError {
	return NewAPIError(http.StatusForbidden, http.StatusForbidden, message)
}

// ErrNotFound 创建 404 未找到错误。
func ErrNotFound(message string) *APIError {
	return NewAPIError(http.StatusNotFound, http.StatusNotFound, message)
}

// ErrConflict 创建 409 冲突错误。
func ErrConflict(message string) *APIError {
	return NewAPIError(http.StatusConflict, http.StatusConflict, message)
}

// ErrInternal 创建 500 内部服务器错误。
func ErrInternal(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, http.StatusInternalServerError, message)
}

// ErrInternalWrap 创建包装其他错误的 500 内部服务器错误。
func ErrInternalWrap(message string, err error) *APIError {
	return WrapAPIError(http.StatusInternalServerError, http.StatusInternalServerError, message, err)
}

// ErrNotImplemented 创建 501 未实现错误。
func ErrNotImplemented(method string) *APIError {
	return NewAPIError(http.StatusNotImplemented, http.StatusNotImplemented, method+" not implemented")
}

// BindError 表示请求绑定错误。
type BindError struct {
	Err    error
	Source string
}

func (e *BindError) Error() string {
	return fmt.Sprintf("binding error (%s): %v", e.Source, e.Err)
}

func (e *BindError) Unwrap() error {
	return e.Err
}

// NewBindError 创建一个新的绑定错误。
func NewBindError(source string, err error) *BindError {
	return &BindError{
		Source: source,
		Err:    err,
	}
}

// ValidationError 表示字段级验证错误。
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors 包含多个验证错误。
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %d errors", len(e.Errors))
}

// Add 添加一个验证错误。
func (e *ValidationErrors) Add(field, message string) {
	e.Errors = append(e.Errors, ValidationError{Field: field, Message: message})
}

// HasErrors 返回是否有验证错误。
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}
