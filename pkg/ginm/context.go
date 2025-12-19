package ginm

import "github.com/gin-gonic/gin"

// ContextKey 是类型安全的上下文键。
type ContextKey[T any] string

// NewContextKey 创建一个新的类型化上下文键。
func NewContextKey[T any](name string) ContextKey[T] {
	return ContextKey[T](name)
}

// Set 在上下文中存储类型化的值。
func Set[T any](c *gin.Context, key ContextKey[T], value T) {
	c.Set(string(key), value)
}

// Get 从上下文中获取类型化的值。
func Get[T any](c *gin.Context, key ContextKey[T]) (T, bool) {
	value, exists := c.Get(string(key))
	if !exists {
		var zero T
		return zero, false
	}
	typed, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return typed, true
}

// MustGet 获取类型化的值，如果不存在则 panic。
func MustGet[T any](c *gin.Context, key ContextKey[T]) T {
	value, ok := Get(c, key)
	if !ok {
		panic("context value not found: " + string(key))
	}
	return value
}

// GetOrDefault 获取类型化的值，如果不存在则返回默认值。
func GetOrDefault[T any](c *gin.Context, key ContextKey[T], defaultValue T) T {
	value, ok := Get(c, key)
	if !ok {
		return defaultValue
	}
	return value
}

// Clear 将上下文值设为 nil。
// 注意: Gin 不支持从上下文中真正删除键。
// 此函数将值设为 nil，这会导致 Get 返回 (零值, false)。
// 键仍然存在于上下文 map 中，但值为 nil。
func Clear[T any](c *gin.Context, key ContextKey[T]) {
	c.Set(string(key), nil)
}

// 常用上下文键

var (
	// UserIDKey 用于存储当前用户的 ID。
	UserIDKey = NewContextKey[int64]("ginm:user_id")
	// RequestIDKey 用于存储请求 ID。
	RequestIDKey = NewContextKey[string]("ginm:request_id")
	// TenantIDKey 用于存储多租户应用中的租户 ID。
	TenantIDKey = NewContextKey[string]("ginm:tenant_id")
)

// SetUserID 是设置用户 ID 的便捷函数。
func SetUserID(c *gin.Context, userID int64) {
	Set(c, UserIDKey, userID)
}

// GetUserID 是获取用户 ID 的便捷函数。
func GetUserID(c *gin.Context) (int64, bool) {
	return Get(c, UserIDKey)
}

// SetRequestID 是设置请求 ID 的便捷函数。
func SetRequestID(c *gin.Context, requestID string) {
	Set(c, RequestIDKey, requestID)
}

// GetRequestID 是获取请求 ID 的便捷函数。
func GetRequestID(c *gin.Context) (string, bool) {
	return Get(c, RequestIDKey)
}

// SetTenantID 是设置租户 ID 的便捷函数。
func SetTenantID(c *gin.Context, tenantID string) {
	Set(c, TenantIDKey, tenantID)
}

// GetTenantID 是获取租户 ID 的便捷函数。
func GetTenantID(c *gin.Context) (string, bool) {
	return Get(c, TenantIDKey)
}
