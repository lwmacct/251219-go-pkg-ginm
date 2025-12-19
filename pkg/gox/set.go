package gox

// --- 切片的集合操作 ---

// Intersect 返回同时存在于两个切片中的元素。
func Intersect[T comparable](a, b []T) []T {
	if len(a) == 0 || len(b) == 0 {
		return []T{}
	}

	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}

	result := make([]T, 0)
	seen := make(map[T]struct{})
	for _, item := range a {
		if _, inB := set[item]; inB {
			if _, alreadySeen := seen[item]; !alreadySeen {
				result = append(result, item)
				seen[item] = struct{}{}
			}
		}
	}
	return result
}

// Union 返回存在于任一切片中的元素（去重）。
func Union[T comparable](a, b []T) []T {
	seen := make(map[T]struct{})
	result := make([]T, 0, len(a)+len(b))

	for _, item := range a {
		if _, ok := seen[item]; !ok {
			result = append(result, item)
			seen[item] = struct{}{}
		}
	}
	for _, item := range b {
		if _, ok := seen[item]; !ok {
			result = append(result, item)
			seen[item] = struct{}{}
		}
	}
	return result
}

// Difference 返回在 a 中但不在 b 中的元素。
func Difference[T comparable](a, b []T) []T {
	if len(a) == 0 {
		return []T{}
	}
	if len(b) == 0 {
		return append([]T{}, a...)
	}

	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}

	result := make([]T, 0)
	for _, item := range a {
		if _, inB := set[item]; !inB {
			result = append(result, item)
		}
	}
	return result
}

// SymmetricDifference 返回恰好存在于一个切片中的元素。
func SymmetricDifference[T comparable](a, b []T) []T {
	return Union(Difference(a, b), Difference(b, a))
}

// Partition 根据条件将元素分成两个切片。
// 返回 (满足条件的, 不满足条件的)。
func Partition[T any](items []T, fn func(T) bool) ([]T, []T) {
	matching := make([]T, 0)
	notMatching := make([]T, 0)
	for _, item := range items {
		if fn(item) {
			matching = append(matching, item)
		} else {
			notMatching = append(notMatching, item)
		}
	}
	return matching, notMatching
}

// IsSubset 检查 a 的所有元素是否都在 b 中。
func IsSubset[T comparable](a, b []T) bool {
	if len(a) == 0 {
		return true
	}
	if len(b) == 0 {
		return false
	}

	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}

	for _, item := range a {
		if _, ok := set[item]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset 检查 b 的所有元素是否都在 a 中。
func IsSuperset[T comparable](a, b []T) bool {
	return IsSubset(b, a)
}

// IsDisjoint 检查 a 和 b 是否没有共同元素。
func IsDisjoint[T comparable](a, b []T) bool {
	if len(a) == 0 || len(b) == 0 {
		return true
	}

	set := make(map[T]struct{}, len(b))
	for _, item := range b {
		set[item] = struct{}{}
	}

	for _, item := range a {
		if _, ok := set[item]; ok {
			return false
		}
	}
	return true
}

// Count 返回满足条件的元素数量。
func Count[T any](items []T, fn func(T) bool) int {
	count := 0
	for _, item := range items {
		if fn(item) {
			count++
		}
	}
	return count
}

// CountBy 返回每个键的计数 map。
func CountBy[T any, K comparable](items []T, fn func(T) K) map[K]int {
	result := make(map[K]int)
	for _, item := range items {
		key := fn(item)
		result[key]++
	}
	return result
}

// Zip 将两个切片组合成对的切片。
// 结果长度为两个输入长度的最小值。
func Zip[T, U any](a []T, b []U) []struct {
	First  T
	Second U
} {
	length := min(len(a), len(b))

	result := make([]struct {
		First  T
		Second U
	}, length)

	for i := range length {
		result[i] = struct {
			First  T
			Second U
		}{First: a[i], Second: b[i]}
	}
	return result
}

// Unzip 将对的切片拆分为两个切片。
func Unzip[T, U any](pairs []struct {
	First  T
	Second U
}) ([]T, []U) {
	a := make([]T, len(pairs))
	b := make([]U, len(pairs))
	for i, pair := range pairs {
		a[i] = pair.First
		b[i] = pair.Second
	}
	return a, b
}
