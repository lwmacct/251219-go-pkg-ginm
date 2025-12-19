package gox

import "cmp"

// Signed 是有符号整数类型的约束。
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 是无符号整数类型的约束。
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer 是所有整数类型的约束。
type Integer interface {
	Signed | Unsigned
}

// Float 是浮点类型的约束。
type Float interface {
	~float32 | ~float64
}

// Numeric 是所有数值类型的约束。
type Numeric interface {
	Integer | Float
}

// Ordered 是支持排序的类型约束。
// 使用标准库的 cmp.Ordered。
type Ordered = cmp.Ordered

// --- 数值运算 ---

// Sum 返回所有元素的和。
func Sum[T Numeric](items []T) T {
	var sum T
	for _, item := range items {
		sum += item
	}
	return sum
}

// SumBy 返回从元素提取的值的和。
func SumBy[T any, N Numeric](items []T, fn func(T) N) N {
	var sum N
	for _, item := range items {
		sum += fn(item)
	}
	return sum
}

// Average 返回所有元素的算术平均值。
// 空切片返回 0。
func Average[T Numeric](items []T) float64 {
	if len(items) == 0 {
		return 0
	}
	var sum float64
	for _, item := range items {
		sum += float64(item)
	}
	return sum / float64(len(items))
}

// Max 返回参数中的最大值。
// 如果没有提供参数则 panic。
func Max[T Ordered](items ...T) T {
	if len(items) == 0 {
		panic("Max requires at least one argument")
	}
	max := items[0]
	for _, item := range items[1:] {
		if item > max {
			max = item
		}
	}
	return max
}

// MaxSlice 返回切片中的最大值。
// 空切片返回零值和 false。
func MaxSlice[T Ordered](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	max := items[0]
	for _, item := range items[1:] {
		if item > max {
			max = item
		}
	}
	return max, true
}

// MaxBy 返回选择器函数返回最大值的元素。
func MaxBy[T any, K Ordered](items []T, fn func(T) K) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	maxItem := items[0]
	maxKey := fn(items[0])
	for _, item := range items[1:] {
		key := fn(item)
		if key > maxKey {
			maxKey = key
			maxItem = item
		}
	}
	return maxItem, true
}

// Min 返回参数中的最小值。
// 如果没有提供参数则 panic。
func Min[T Ordered](items ...T) T {
	if len(items) == 0 {
		panic("Min requires at least one argument")
	}
	min := items[0]
	for _, item := range items[1:] {
		if item < min {
			min = item
		}
	}
	return min
}

// MinSlice 返回切片中的最小值。
// 空切片返回零值和 false。
func MinSlice[T Ordered](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	min := items[0]
	for _, item := range items[1:] {
		if item < min {
			min = item
		}
	}
	return min, true
}

// MinBy 返回选择器函数返回最小值的元素。
func MinBy[T any, K Ordered](items []T, fn func(T) K) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	minItem := items[0]
	minKey := fn(items[0])
	for _, item := range items[1:] {
		key := fn(item)
		if key < minKey {
			minKey = key
			minItem = item
		}
	}
	return minItem, true
}

// Clamp 将值限制在指定范围 [min, max] 内。
func Clamp[T Ordered](value, minVal, maxVal T) T {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// Abs 返回绝对值。
func Abs[T Signed | Float](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Range 生成从 start 到 end（不包含）的整数切片。
func Range(start, end int) []int {
	if end <= start {
		return []int{}
	}
	result := make([]int, end-start)
	for i := range result {
		result[i] = start + i
	}
	return result
}

// RangeStep 生成从 start 到 end 的整数切片，步长为 step。
func RangeStep(start, end, step int) []int {
	if step == 0 || (step > 0 && end <= start) || (step < 0 && end >= start) {
		return []int{}
	}
	size := (end - start) / step
	if size < 0 {
		size = -size
	}
	result := make([]int, 0, size)
	if step > 0 {
		for i := start; i < end; i += step {
			result = append(result, i)
		}
	} else {
		for i := start; i > end; i += step {
			result = append(result, i)
		}
	}
	return result
}
