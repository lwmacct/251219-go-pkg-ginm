// Package gox 提供 Go 1.18+ 泛型工具函数。
//
// 本包包含纯工具函数，不依赖任何特定框架。包括函数式编程辅助、
// 单子类型、数值运算、集合操作和类型转换。
//
// # 函数式编程
//
// 受函数式编程启发的集合操作：
//
//	nums := []int{1, 2, 3, 4, 5}
//	doubled := gox.Map(nums, func(n int) int { return n * 2 })
//	evens := gox.Filter(nums, func(n int) bool { return n%2 == 0 })
//	sum := gox.Reduce(nums, 0, func(acc, n int) int { return acc + n })
//
// # Result 类型
//
// Rust 风格的 Result 类型，用于显式错误处理：
//
//	result := gox.Try(func() (int, error) {
//	    return strconv.Atoi("42")
//	})
//	value := result.UnwrapOr(0)
//
// # Optional 类型
//
// 空值安全的 Optional 类型：
//
//	opt := gox.OSome(42)
//	value := opt.OrElse(0)
//
// # 数值运算
//
// 泛型数值工具：
//
//	max := gox.Max(1, 5, 3)           // 5
//	sum := gox.Sum([]int{1, 2, 3})    // 6
//	clamped := gox.Clamp(15, 0, 10)   // 10
//
// # 集合操作
//
// 切片的集合运算：
//
//	a := []int{1, 2, 3}
//	b := []int{2, 3, 4}
//	gox.Intersect(a, b)  // [2, 3]
//	gox.Union(a, b)      // [1, 2, 3, 4]
//
// # 类型转换
//
// 安全的类型转换，返回 Result 或 Optional：
//
//	gox.ParseInt("42").UnwrapOr(0)      // 42
//	gox.ParseIntO("invalid").OrElse(0)  // 0
package gox
