package ginm

import "github.com/gin-gonic/gin"

// HandlerChain 提供组合中间件和处理器的流式 API。
type HandlerChain struct {
	middlewares []gin.HandlerFunc
}

// Chain 使用给定的中间件创建新的 HandlerChain。
func Chain(middlewares ...gin.HandlerFunc) *HandlerChain {
	return &HandlerChain{
		middlewares: append([]gin.HandlerFunc{}, middlewares...),
	}
}

// Use 将中间件添加到链中，返回链本身以支持方法链式调用。
func (c *HandlerChain) Use(middleware gin.HandlerFunc) *HandlerChain {
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Then 是 Use 的别名，使代码读起来更流畅。
func (c *HandlerChain) Then(middleware gin.HandlerFunc) *HandlerChain {
	return c.Use(middleware)
}

// UseIf 根据条件添加中间件到链中。
func (c *HandlerChain) UseIf(cond bool, middleware gin.HandlerFunc) *HandlerChain {
	if cond {
		c.middlewares = append(c.middlewares, middleware)
	}
	return c
}

// UseMany 添加多个中间件到链中。
func (c *HandlerChain) UseMany(middlewares ...gin.HandlerFunc) *HandlerChain {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

// Handle 使用链中所有中间件包装处理器。
func (c *HandlerChain) Handle(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 按顺序执行中间件
		for _, middleware := range c.middlewares {
			middleware(ctx)
			if ctx.IsAborted() {
				return
			}
		}
		// 执行最终处理器
		handler(ctx)
	}
}

// HandleFunc 是 Handle 的别名。
func (c *HandlerChain) HandleFunc(handler gin.HandlerFunc) gin.HandlerFunc {
	return c.Handle(handler)
}

// Handlers 返回所有中间件加上处理器作为切片。
// 适用于 gin.RouterGroup.Handle()。
func (c *HandlerChain) Handlers(handler gin.HandlerFunc) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, 0, len(c.middlewares)+1)
	result = append(result, c.middlewares...)
	result = append(result, handler)
	return result
}

// Clone 创建链的副本，可以独立修改。
func (c *HandlerChain) Clone() *HandlerChain {
	return &HandlerChain{
		middlewares: append([]gin.HandlerFunc{}, c.middlewares...),
	}
}

// Len 返回链中中间件的数量。
func (c *HandlerChain) Len() int {
	return len(c.middlewares)
}

// --- 创建链的便捷函数 ---

// ChainOf 从多个中间件函数创建链。
func ChainOf(middlewares ...gin.HandlerFunc) *HandlerChain {
	return Chain(middlewares...)
}

// --- 路由辅助 ---

// RouterChain 使用中间件链包装 gin.RouterGroup。
type RouterChain struct {
	group *gin.RouterGroup
	chain *HandlerChain
}

// WithChain 创建用于流式路由注册的 RouterChain。
func WithChain(group *gin.RouterGroup, middlewares ...gin.HandlerFunc) *RouterChain {
	return &RouterChain{
		group: group,
		chain: Chain(middlewares...),
	}
}

// Use 添加中间件到路由链。
func (rc *RouterChain) Use(middleware gin.HandlerFunc) *RouterChain {
	rc.chain.Use(middleware)
	return rc
}

// GET 注册带链中中间件的 GET 路由。
func (rc *RouterChain) GET(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.GET(path, rc.chain.Handle(handler))
	return rc
}

// POST 注册带链中中间件的 POST 路由。
func (rc *RouterChain) POST(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.POST(path, rc.chain.Handle(handler))
	return rc
}

// PUT 注册带链中中间件的 PUT 路由。
func (rc *RouterChain) PUT(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.PUT(path, rc.chain.Handle(handler))
	return rc
}

// DELETE 注册带链中中间件的 DELETE 路由。
func (rc *RouterChain) DELETE(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.DELETE(path, rc.chain.Handle(handler))
	return rc
}

// PATCH 注册带链中中间件的 PATCH 路由。
func (rc *RouterChain) PATCH(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.PATCH(path, rc.chain.Handle(handler))
	return rc
}

// OPTIONS 注册带链中中间件的 OPTIONS 路由。
func (rc *RouterChain) OPTIONS(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.OPTIONS(path, rc.chain.Handle(handler))
	return rc
}

// HEAD 注册带链中中间件的 HEAD 路由。
func (rc *RouterChain) HEAD(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.HEAD(path, rc.chain.Handle(handler))
	return rc
}

// Any 为所有 HTTP 方法注册路由。
func (rc *RouterChain) Any(path string, handler gin.HandlerFunc) *RouterChain {
	rc.group.Any(path, rc.chain.Handle(handler))
	return rc
}

// Group 使用相同的链创建子分组。
func (rc *RouterChain) Group(path string) *RouterChain {
	return &RouterChain{
		group: rc.group.Group(path),
		chain: rc.chain.Clone(),
	}
}
