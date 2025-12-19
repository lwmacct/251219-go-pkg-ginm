package gox

// Optional 表示一个可能存在或不存在的值。
// 灵感来自 Java 的 Optional 和 Rust 的 Option 类型。
type Optional[T any] struct {
	value T
	valid bool
}

// OSome 创建一个包含值的 Optional。
func OSome[T any](v T) Optional[T] {
	return Optional[T]{value: v, valid: true}
}

// ONone 创建一个空的 Optional。
func ONone[T any]() Optional[T] {
	return Optional[T]{valid: false}
}

// OFromPtr 从指针创建 Optional。
// 如果指针为 nil 返回 None，否则返回 Some(*p)。
func OFromPtr[T any](p *T) Optional[T] {
	if p == nil {
		return ONone[T]()
	}
	return OSome(*p)
}

// OFromOk 从 (value, ok) 对创建 Optional。
func OFromOk[T any](v T, ok bool) Optional[T] {
	if !ok {
		return ONone[T]()
	}
	return OSome(v)
}

// OFromErr 从 (value, error) 对创建 Optional。
// 如果 error 不为 nil 则返回 None。
func OFromErr[T any](v T, err error) Optional[T] {
	if err != nil {
		return ONone[T]()
	}
	return OSome(v)
}

// IsSome 返回 Optional 是否包含值。
func (o Optional[T]) IsSome() bool {
	return o.valid
}

// IsNone 返回 Optional 是否为空。
func (o Optional[T]) IsNone() bool {
	return !o.valid
}

// Get 返回值和表示是否存在的布尔值。
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.valid
}

// MustGet 返回值，如果为空则 panic。
func (o Optional[T]) MustGet() T {
	if !o.valid {
		panic("called MustGet on None")
	}
	return o.value
}

// OrElse 如果有值则返回值，否则返回默认值。
func (o Optional[T]) OrElse(def T) T {
	if !o.valid {
		return def
	}
	return o.value
}

// OrElseFn 如果有值则返回值，否则调用 fn。
func (o Optional[T]) OrElseFn(fn func() T) T {
	if !o.valid {
		return fn()
	}
	return o.value
}

// OrElseZero 如果有值则返回值，否则返回零值。
func (o Optional[T]) OrElseZero() T {
	if !o.valid {
		var zero T
		return zero
	}
	return o.value
}

// ToPtr 如果有值则返回值的指针，否则返回 nil。
func (o Optional[T]) ToPtr() *T {
	if !o.valid {
		return nil
	}
	return &o.value
}

// Map 如果有值则对其应用函数。
func (o Optional[T]) Map(fn func(T) T) Optional[T] {
	if !o.valid {
		return o
	}
	return OSome(fn(o.value))
}

// OMapTo 应用类型转换函数。
func OMapTo[T, R any](o Optional[T], fn func(T) R) Optional[R] {
	if !o.valid {
		return ONone[R]()
	}
	return OSome(fn(o.value))
}

// FlatMap 应用返回 Optional 的函数。
func (o Optional[T]) FlatMap(fn func(T) Optional[T]) Optional[T] {
	if !o.valid {
		return o
	}
	return fn(o.value)
}

// OFlatMapTo 应用转换为不同 Optional 类型的函数。
func OFlatMapTo[T, R any](o Optional[T], fn func(T) Optional[R]) Optional[R] {
	if !o.valid {
		return ONone[R]()
	}
	return fn(o.value)
}

// Filter 如果值不满足条件则返回 None。
func (o Optional[T]) Filter(fn func(T) bool) Optional[T] {
	if !o.valid || !fn(o.value) {
		return ONone[T]()
	}
	return o
}

// Or 如果当前有值则返回当前，否则返回替代值。
func (o Optional[T]) Or(alt Optional[T]) Optional[T] {
	if o.valid {
		return o
	}
	return alt
}

// OrFn 如果当前有值则返回当前，否则调用 fn。
func (o Optional[T]) OrFn(fn func() Optional[T]) Optional[T] {
	if o.valid {
		return o
	}
	return fn()
}

// And 如果两者都有值则返回 other，否则返回 None。
func (o Optional[T]) And(other Optional[T]) Optional[T] {
	if !o.valid {
		return ONone[T]()
	}
	return other
}

// Xor 如果恰好一个有值则返回 Some，否则返回 None。
func (o Optional[T]) Xor(other Optional[T]) Optional[T] {
	if o.valid && !other.valid {
		return o
	}
	if !o.valid && other.valid {
		return other
	}
	return ONone[T]()
}

// Inspect 如果有值则用值调用 fn，返回 Optional 本身。
func (o Optional[T]) Inspect(fn func(T)) Optional[T] {
	if o.valid {
		fn(o.value)
	}
	return o
}

// OMatch 如果有值执行 someFn，否则执行 noneFn。
func OMatch[T, R any](o Optional[T], someFn func(T) R, noneFn func() R) R {
	if o.valid {
		return someFn(o.value)
	}
	return noneFn()
}

// ToResult 将 Optional 转换为带自定义错误的 Result。
func (o Optional[T]) ToResult(err error) Result[T] {
	if !o.valid {
		return RErr[T](err)
	}
	return ROk(o.value)
}

// OZip 将两个 Optional 组合为一个。
func OZip[T, U any](a Optional[T], b Optional[U]) Optional[struct {
	First  T
	Second U
}] {
	if !a.valid || !b.valid {
		return ONone[struct {
			First  T
			Second U
		}]()
	}
	return OSome(struct {
		First  T
		Second U
	}{First: a.value, Second: b.value})
}

// OUnzip 将 Optional 的对拆分为两个 Optional。
func OUnzip[T, U any](o Optional[struct {
	First  T
	Second U
}]) (Optional[T], Optional[U]) {
	if !o.valid {
		return ONone[T](), ONone[U]()
	}
	return OSome(o.value.First), OSome(o.value.Second)
}
