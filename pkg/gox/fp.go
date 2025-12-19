package gox

// Map 对切片中每个元素应用函数，返回转换后的新切片。
func Map[T, R any](items []T, fn func(T) R) []R {
	if items == nil {
		return nil
	}
	result := make([]R, len(items))
	for i, item := range items {
		result[i] = fn(item)
	}
	return result
}

// Filter 返回满足条件的元素组成的新切片。
func Filter[T any](items []T, fn func(T) bool) []T {
	if items == nil {
		return nil
	}
	result := make([]T, 0)
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

// Reduce 使用累加函数将切片归约为单个值。
func Reduce[T, R any](items []T, init R, fn func(R, T) R) R {
	result := init
	for _, item := range items {
		result = fn(result, item)
	}
	return result
}

// Find 返回第一个满足条件的元素。
func Find[T any](items []T, fn func(T) bool) (T, bool) {
	for _, item := range items {
		if fn(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex 返回第一个满足条件的元素的索引。
// 未找到时返回 -1。
func FindIndex[T any](items []T, fn func(T) bool) int {
	for i, item := range items {
		if fn(item) {
			return i
		}
	}
	return -1
}

// Every 检查所有元素是否都满足条件。
func Every[T any](items []T, fn func(T) bool) bool {
	for _, item := range items {
		if !fn(item) {
			return false
		}
	}
	return true
}

// Some 检查是否至少有一个元素满足条件。
func Some[T any](items []T, fn func(T) bool) bool {
	for _, item := range items {
		if fn(item) {
			return true
		}
	}
	return false
}

// Contains 检查切片是否包含指定元素。
func Contains[T comparable](items []T, item T) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

// Unique 返回去重后的新切片。
func Unique[T comparable](items []T) []T {
	if items == nil {
		return nil
	}
	seen := make(map[T]struct{})
	result := make([]T, 0)
	for _, item := range items {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// GroupBy 按键函数对元素分组。
func GroupBy[T any, K comparable](items []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range items {
		key := fn(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Chunk 将切片分割成指定大小的块。
func Chunk[T any](items []T, size int) [][]T {
	if size <= 0 || len(items) == 0 {
		return nil
	}
	result := make([][]T, 0, (len(items)+size-1)/size)
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		result = append(result, items[i:end])
	}
	return result
}

// Flatten 将二维切片展平为一维切片。
func Flatten[T any](items [][]T) []T {
	if items == nil {
		return nil
	}
	result := make([]T, 0)
	for _, inner := range items {
		result = append(result, inner...)
	}
	return result
}

// First 返回切片的第一个元素。
func First[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	return items[0], true
}

// Last 返回切片的最后一个元素。
func Last[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T
		return zero, false
	}
	return items[len(items)-1], true
}

// Reverse 返回元素顺序相反的新切片。
func Reverse[T any](items []T) []T {
	if items == nil {
		return nil
	}
	result := make([]T, len(items))
	for i, j := 0, len(items)-1; i < len(items); i, j = i+1, j-1 {
		result[i] = items[j]
	}
	return result
}

// --- 指针工具 ---

// Ptr 返回给定值的指针。
func Ptr[T any](v T) *T {
	return &v
}

// Val 安全地解引用指针，如果为 nil 则返回默认值。
func Val[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// ValOrZero 安全地解引用指针，如果为 nil 则返回零值。
func ValOrZero[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// Coalesce 返回第一个非零值。
func Coalesce[T comparable](vals ...T) T {
	var zero T
	for _, v := range vals {
		if v != zero {
			return v
		}
	}
	return zero
}

// CoalescePtr 返回第一个非 nil 指针的值。
func CoalescePtr[T any](ptrs ...*T) (T, bool) {
	for _, p := range ptrs {
		if p != nil {
			return *p, true
		}
	}
	var zero T
	return zero, false
}

// --- 三元运算符 ---

// If 根据条件返回 trueVal 或 falseVal。
func If[T any](cond bool, trueVal, falseVal T) T {
	if cond {
		return trueVal
	}
	return falseVal
}

// IfFn 根据条件调用 trueFn 或 falseFn。
// 适用于值计算开销较大的场景。
func IfFn[T any](cond bool, trueFn, falseFn func() T) T {
	if cond {
		return trueFn()
	}
	return falseFn()
}

// --- 转换工具 ---

// Keys 返回 map 的所有键。
func Keys[K comparable, V any](m map[K]V) []K {
	if m == nil {
		return nil
	}
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// Values 返回 map 的所有值。
func Values[K comparable, V any](m map[K]V) []V {
	if m == nil {
		return nil
	}
	result := make([]V, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

// Entries 返回 map 的键值对切片。
func Entries[K comparable, V any](m map[K]V) []struct {
	Key   K
	Value V
} {
	if m == nil {
		return nil
	}
	result := make([]struct {
		Key   K
		Value V
	}, 0, len(m))
	for k, v := range m {
		result = append(result, struct {
			Key   K
			Value V
		}{Key: k, Value: v})
	}
	return result
}

// FromEntries 从键值对切片创建 map。
func FromEntries[K comparable, V any](entries []struct {
	Key   K
	Value V
}) map[K]V {
	result := make(map[K]V)
	for _, e := range entries {
		result[e.Key] = e.Value
	}
	return result
}
