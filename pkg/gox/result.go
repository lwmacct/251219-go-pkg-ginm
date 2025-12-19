package gox

// Result 表示一个可能成功（Ok）或失败（Err）的值。
// 灵感来自 Rust 的 Result 类型，提供了一种无需多返回值的错误处理方式。
type Result[T any] struct {
	data T
	err  error
}

// ROk 创建一个成功的 Result。
func ROk[T any](data T) Result[T] {
	return Result[T]{data: data}
}

// RErr 创建一个失败的 Result。
func RErr[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Try 执行函数并将结果包装为 Result 类型。
func Try[T any](fn func() (T, error)) Result[T] {
	data, err := fn()
	if err != nil {
		return RErr[T](err)
	}
	return ROk(data)
}

// TryE 执行仅返回 error 的函数。
func TryE(fn func() error) Result[struct{}] {
	if err := fn(); err != nil {
		return RErr[struct{}](err)
	}
	return ROk(struct{}{})
}

// IsOk 返回 Result 是否成功。
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr 返回 Result 是否失败。
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Unwrap 返回数据，如果是 Err 则 panic。
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.data
}

// UnwrapErr 返回错误，如果是 Ok 则 panic。
func (r Result[T]) UnwrapErr() error {
	if r.err == nil {
		panic("called UnwrapErr on an Ok value")
	}
	return r.err
}

// UnwrapOr 返回数据，如果是 Err 则返回默认值。
func (r Result[T]) UnwrapOr(def T) T {
	if r.err != nil {
		return def
	}
	return r.data
}

// UnwrapOrElse 返回数据，如果是 Err 则调用 fn 并返回其结果。
func (r Result[T]) UnwrapOrElse(fn func() T) T {
	if r.err != nil {
		return fn()
	}
	return r.data
}

// UnwrapOrDefault 返回数据，如果是 Err 则返回零值。
func (r Result[T]) UnwrapOrDefault() T {
	if r.err != nil {
		var zero T
		return zero
	}
	return r.data
}

// Get 返回数据和表示成功的布尔值。
func (r Result[T]) Get() (T, bool) {
	return r.data, r.err == nil
}

// GetWithError 返回数据和错误。
func (r Result[T]) GetWithError() (T, error) {
	return r.data, r.err
}

// Error 返回错误，如果成功则返回 nil。
func (r Result[T]) Error() error {
	return r.err
}

// Map 如果是 Ok 则对数据应用函数，Err 保持不变。
func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return ROk(fn(r.data))
}

// MapTo 如果是 Ok 则应用类型转换函数。
func MapTo[T, R any](r Result[T], fn func(T) R) Result[R] {
	if r.err != nil {
		return RErr[R](r.err)
	}
	return ROk(fn(r.data))
}

// MapErr 如果是 Err 则对错误应用函数，Ok 保持不变。
func (r Result[T]) MapErr(fn func(error) error) Result[T] {
	if r.err == nil {
		return r
	}
	return RErr[T](fn(r.err))
}

// AndThen 链接另一个返回 Result 的操作。
func (r Result[T]) AndThen(fn func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return fn(r.data)
}

// AndThenTo 链接另一个类型转换操作。
func AndThenTo[T, R any](r Result[T], fn func(T) Result[R]) Result[R] {
	if r.err != nil {
		return RErr[R](r.err)
	}
	return fn(r.data)
}

// OrElse 如果当前是 Err 则提供替代 Result。
func (r Result[T]) OrElse(fn func(error) Result[T]) Result[T] {
	if r.err == nil {
		return r
	}
	return fn(r.err)
}

// Inspect 如果是 Ok 则用数据调用 fn，如果是 Err 则不做任何事。
// 返回 Result 本身用于链式调用。
func (r Result[T]) Inspect(fn func(T)) Result[T] {
	if r.err == nil {
		fn(r.data)
	}
	return r
}

// InspectErr 如果是 Err 则用错误调用 fn，如果是 Ok 则不做任何事。
func (r Result[T]) InspectErr(fn func(error)) Result[T] {
	if r.err != nil {
		fn(r.err)
	}
	return r
}

// Match 如果是 Ok 执行 okFn，如果是 Err 执行 errFn。
func Match[T, R any](r Result[T], okFn func(T) R, errFn func(error) R) R {
	if r.err != nil {
		return errFn(r.err)
	}
	return okFn(r.data)
}

// FlattenResult 展平嵌套的 Result。
func FlattenResult[T any](r Result[Result[T]]) Result[T] {
	if r.err != nil {
		return RErr[T](r.err)
	}
	return r.data
}

// Collect 将 Result 切片收集为切片的 Result。
// 遇到第一个错误就返回，或者返回包含所有值的 Ok。
func Collect[T any](results []Result[T]) Result[[]T] {
	data := make([]T, 0, len(results))
	for _, r := range results {
		if r.err != nil {
			return RErr[[]T](r.err)
		}
		data = append(data, r.data)
	}
	return ROk(data)
}
