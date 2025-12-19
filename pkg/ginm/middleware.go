package ginm

import "github.com/gin-gonic/gin"

// Extractor 是从请求中提取类型化值的函数。
type Extractor[T any] func(c *gin.Context) (T, error)

// WithContext 创建一个中间件，用于提取值并存储到上下文中。
// 如果提取失败，返回错误响应并中止请求。
func WithContext[T any](key ContextKey[T], extractor Extractor[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, err := extractor(c)
		if err != nil {
			handleError(c, err)
			c.Abort()
			return
		}
		Set(c, key, value)
		c.Next()
	}
}

// WithContextOptional 类似 WithContext，但不会在错误时中止。
// 如果提取失败，设置默认值并继续。
func WithContextOptional[T any](key ContextKey[T], extractor Extractor[T], defaultValue T) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, err := extractor(c)
		if err != nil {
			Set(c, key, defaultValue)
		} else {
			Set(c, key, value)
		}
		c.Next()
	}
}

// RequireContext 创建一个确保上下文值存在的中间件。
// 如果未找到值，返回未授权错误。
func RequireContext[T any](key ContextKey[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := Get(c, key); !ok {
			handleError(c, ErrUnauthorized("authentication required"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireContextWithError 创建一个确保上下文值存在的中间件。
// 如果未找到值，返回指定的错误。
func RequireContextWithError[T any](key ContextKey[T], err *APIError) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := Get(c, key); !ok {
			handleError(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}

// Transform 创建一个将一个上下文值转换为另一个的中间件。
func Transform[From, To any](
	fromKey ContextKey[From],
	toKey ContextKey[To],
	transformer func(From) (To, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		from, ok := Get(c, fromKey)
		if !ok {
			c.Next()
			return
		}

		to, err := transformer(from)
		if err != nil {
			handleError(c, err)
			c.Abort()
			return
		}

		Set(c, toKey, to)
		c.Next()
	}
}

// TransformOptional 类似 Transform，但不会在错误时中止。
func TransformOptional[From, To any](
	fromKey ContextKey[From],
	toKey ContextKey[To],
	transformer func(From) (To, error),
	defaultValue To,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		from, ok := Get(c, fromKey)
		if !ok {
			c.Next()
			return
		}

		to, err := transformer(from)
		if err != nil {
			Set(c, toKey, defaultValue)
		} else {
			Set(c, toKey, to)
		}
		c.Next()
	}
}

// Validator 是验证条件并在无效时返回错误的函数。
type Validator func(c *gin.Context) error

// Validate 创建一个在继续之前运行验证器的中间件。
func Validate(validators ...Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range validators {
			if err := v(c); err != nil {
				handleError(c, err)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequireHeader 创建一个检查必需请求头的验证器。
func RequireHeader(header string) Validator {
	return func(c *gin.Context) error {
		if c.GetHeader(header) == "" {
			return ErrBadRequest("missing required header: " + header)
		}
		return nil
	}
}

// RequireQuery 创建一个检查必需查询参数的验证器。
func RequireQuery(param string) Validator {
	return func(c *gin.Context) error {
		if c.Query(param) == "" {
			return ErrBadRequest("missing required query parameter: " + param)
		}
		return nil
	}
}
