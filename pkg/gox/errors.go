package gox

import "fmt"

// MultiError 聚合多个错误为一个。
// 适用于批量操作、并行处理或收集所有验证失败。
type MultiError struct {
	errors []error
}

// NewMultiError 创建一个新的空 MultiError。
func NewMultiError() *MultiError {
	return &MultiError{}
}

// Add 添加一个错误到集合。nil 错误会被忽略。
func (m *MultiError) Add(err error) {
	if err != nil {
		m.errors = append(m.errors, err)
	}
}

// AddAll 添加多个错误。nil 错误会被忽略。
func (m *MultiError) AddAll(errs ...error) {
	for _, err := range errs {
		if err != nil {
			m.errors = append(m.errors, err)
		}
	}
}

// Errors 返回所有收集的错误。
func (m *MultiError) Errors() []error {
	return m.errors
}

// HasErrors 返回是否有任何错误。
func (m *MultiError) HasErrors() bool {
	return len(m.errors) > 0
}

// Len 返回错误数量。
func (m *MultiError) Len() int {
	return len(m.errors)
}

// Error 实现 error 接口。
// 返回所有错误的合并消息。
func (m *MultiError) Error() string {
	if len(m.errors) == 0 {
		return ""
	}
	if len(m.errors) == 1 {
		return m.errors[0].Error()
	}
	msg := fmt.Sprintf("%d errors: ", len(m.errors))
	for i, err := range m.errors {
		if i > 0 {
			msg += "; "
		}
		msg += err.Error()
	}
	return msg
}

// ErrorOrNil 如果没有错误返回 nil，否则返回自身。
// 方便函数返回使用。
func (m *MultiError) ErrorOrNil() error {
	if !m.HasErrors() {
		return nil
	}
	return m
}

// First 返回第一个错误或 nil。
func (m *MultiError) First() error {
	if len(m.errors) == 0 {
		return nil
	}
	return m.errors[0]
}

// Unwrap 返回错误列表，供 errors.Is/As 使用。
// 兼容 Go 1.20+ 的多错误解包。
func (m *MultiError) Unwrap() []error {
	return m.errors
}
