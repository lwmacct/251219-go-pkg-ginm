package gox

import (
	"fmt"
	"strconv"
)

// --- 字符串转类型，返回 Result ---

// ParseInt 将字符串解析为 int。
func ParseInt(s string) Result[int] {
	v, err := strconv.Atoi(s)
	if err != nil {
		return RErr[int](err)
	}
	return ROk(v)
}

// ParseInt64 将字符串解析为 int64。
func ParseInt64(s string) Result[int64] {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return RErr[int64](err)
	}
	return ROk(v)
}

// ParseInt32 将字符串解析为 int32。
func ParseInt32(s string) Result[int32] {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return RErr[int32](err)
	}
	return ROk(int32(v))
}

// ParseUint64 将字符串解析为 uint64。
func ParseUint64(s string) Result[uint64] {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return RErr[uint64](err)
	}
	return ROk(v)
}

// ParseFloat 将字符串解析为 float64。
func ParseFloat(s string) Result[float64] {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return RErr[float64](err)
	}
	return ROk(v)
}

// ParseFloat32 将字符串解析为 float32。
func ParseFloat32(s string) Result[float32] {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return RErr[float32](err)
	}
	return ROk(float32(v))
}

// ParseBool 将字符串解析为 bool。
// 接受: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func ParseBool(s string) Result[bool] {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return RErr[bool](err)
	}
	return ROk(v)
}

// --- 字符串转类型，返回 Optional ---

// ParseIntO 将字符串解析为 int，返回 Optional。
func ParseIntO(s string) Optional[int] {
	v, err := strconv.Atoi(s)
	if err != nil {
		return ONone[int]()
	}
	return OSome(v)
}

// ParseInt64O 将字符串解析为 int64，返回 Optional。
func ParseInt64O(s string) Optional[int64] {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ONone[int64]()
	}
	return OSome(v)
}

// ParseFloatO 将字符串解析为 float64，返回 Optional。
func ParseFloatO(s string) Optional[float64] {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return ONone[float64]()
	}
	return OSome(v)
}

// ParseBoolO 将字符串解析为 bool，返回 Optional。
func ParseBoolO(s string) Optional[bool] {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return ONone[bool]()
	}
	return OSome(v)
}

// --- 类型转字符串 ---

// ToString 使用 fmt.Sprint 将任意值转换为字符串。
func ToString[T any](v T) string {
	return fmt.Sprint(v)
}

// ToStringf 使用格式字符串将任意值转换为字符串。
func ToStringf[T any](format string, v T) string {
	return fmt.Sprintf(format, v)
}

// IntToString 将整数转换为字符串。
func IntToString[T Integer](v T) string {
	return strconv.FormatInt(int64(v), 10)
}

// FloatToString 将浮点数转换为字符串（默认精度）。
func FloatToString[T Float](v T) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}

// FloatToStringPrec 将浮点数转换为字符串（指定精度）。
func FloatToStringPrec[T Float](v T, prec int) string {
	return strconv.FormatFloat(float64(v), 'f', prec, 64)
}

// BoolToString 将 bool 转换为 "true" 或 "false"。
func BoolToString(v bool) string {
	return strconv.FormatBool(v)
}

// --- 安全类型转换 ---

// IntToInt64 安全地将任意整数类型转换为 int64。
func IntToInt64[T Integer](v T) int64 {
	return int64(v)
}

// IntToInt 安全地将任意整数类型转换为 int。
func IntToInt[T Integer](v T) int {
	return int(v)
}

// FloatToFloat64 安全地将任意浮点类型转换为 float64。
func FloatToFloat64[T Float](v T) float64 {
	return float64(v)
}
