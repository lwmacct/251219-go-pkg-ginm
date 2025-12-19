package ginm

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/lwmacct/251219-go-pkg-ginm/pkg/gox"
)

// Bind 根据 Content-Type 自动绑定请求体到类型化结构体。
func Bind[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBind(&req); err != nil {
		return nil, NewBindError("body", err)
	}
	return &req, nil
}

// BindJSON 绑定 JSON 请求体到类型化结构体。
func BindJSON[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, NewBindError("json", err)
	}
	return &req, nil
}

// BindXML 绑定 XML 请求体到类型化结构体。
func BindXML[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindXML(&req); err != nil {
		return nil, NewBindError("xml", err)
	}
	return &req, nil
}

// BindQuery 绑定查询参数到类型化结构体。
func BindQuery[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindQuery(&req); err != nil {
		return nil, NewBindError("query", err)
	}
	return &req, nil
}

// BindURI 绑定 URI 参数到类型化结构体。
func BindURI[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindUri(&req); err != nil {
		return nil, NewBindError("uri", err)
	}
	return &req, nil
}

// BindHeader 绑定请求头到类型化结构体。
func BindHeader[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindHeader(&req); err != nil {
		return nil, NewBindError("header", err)
	}
	return &req, nil
}

// BindForm 绑定表单数据到类型化结构体。
func BindForm[T any](c *gin.Context) (*T, error) {
	var req T
	if err := c.ShouldBindWith(&req, binding.Form); err != nil {
		return nil, NewBindError("form", err)
	}
	return &req, nil
}

// MustBind 绑定请求，出错时 panic（配合 recovery 中间件使用）。
func MustBind[T any](c *gin.Context) *T {
	req, err := Bind[T](c)
	if err != nil {
		panic(err)
	}
	return req
}

// MustBindJSON 绑定 JSON，出错时 panic。
func MustBindJSON[T any](c *gin.Context) *T {
	req, err := BindJSON[T](c)
	if err != nil {
		panic(err)
	}
	return req
}

// MustBindQuery 绑定查询参数，出错时 panic。
func MustBindQuery[T any](c *gin.Context) *T {
	req, err := BindQuery[T](c)
	if err != nil {
		panic(err)
	}
	return req
}

// BindConfig 指定绑定数据来源。
type BindConfig struct {
	URI   bool
	Query bool
	Body  bool
}

// BindAll 从多个来源绑定到单个结构体。
// 优先级: URI > Query > Body
func BindAll[T any](c *gin.Context, config BindConfig) (*T, error) {
	var req T

	if config.URI {
		if err := c.ShouldBindUri(&req); err != nil {
			return nil, NewBindError("uri", err)
		}
	}

	if config.Query {
		if err := c.ShouldBindQuery(&req); err != nil {
			return nil, NewBindError("query", err)
		}
	}

	if config.Body {
		if err := c.ShouldBind(&req); err != nil {
			return nil, NewBindError("body", err)
		}
	}

	return &req, nil
}

// BindURIAndBody 同时绑定 URI 参数和请求体。
// 适用于 PUT /users/:id 带 JSON body 的路由。
func BindURIAndBody[T any](c *gin.Context) (*T, error) {
	return BindAll[T](c, BindConfig{URI: true, Body: true})
}

// BindURIAndQuery 同时绑定 URI 参数和查询参数。
// 适用于 GET /users/:id?include=posts 的路由。
func BindURIAndQuery[T any](c *gin.Context) (*T, error) {
	return BindAll[T](c, BindConfig{URI: true, Query: true})
}

// --- 常用命名约定的别名 ---

// BindPathAndJSON 是 BindURIAndBody 的别名。
// "Path" 是 URI 参数的常用替代名称。
func BindPathAndJSON[T any](c *gin.Context) (*T, error) {
	return BindURIAndBody[T](c)
}

// BindPathAndQuery 是 BindURIAndQuery 的别名。
func BindPathAndQuery[T any](c *gin.Context) (*T, error) {
	return BindURIAndQuery[T](c)
}

// BindPath 是 BindURI 的别名。
func BindPath[T any](c *gin.Context) (*T, error) {
	return BindURI[T](c)
}

// MustBindURI 绑定 URI，出错时 panic。
func MustBindURI[T any](c *gin.Context) *T {
	req, err := BindURI[T](c)
	if err != nil {
		panic(err)
	}
	return req
}

// MustBindPath 是 MustBindURI 的别名。
func MustBindPath[T any](c *gin.Context) *T {
	return MustBindURI[T](c)
}

// --- 基于 Result 的绑定 ---

// BindR 绑定并返回 Result。
func BindR[T any](c *gin.Context) gox.Result[*T] {
	return gox.Try(func() (*T, error) { return Bind[T](c) })
}

// BindJSONR 绑定 JSON 并返回 Result。
func BindJSONR[T any](c *gin.Context) gox.Result[*T] {
	return gox.Try(func() (*T, error) { return BindJSON[T](c) })
}

// BindQueryR 绑定查询参数并返回 Result。
func BindQueryR[T any](c *gin.Context) gox.Result[*T] {
	return gox.Try(func() (*T, error) { return BindQuery[T](c) })
}

// BindURIR 绑定 URI 并返回 Result。
func BindURIR[T any](c *gin.Context) gox.Result[*T] {
	return gox.Try(func() (*T, error) { return BindURI[T](c) })
}

// --- 基于 Optional 的绑定 ---

// BindO 绑定并返回 Optional（出错时为 None）。
func BindO[T any](c *gin.Context) gox.Optional[*T] {
	req, err := Bind[T](c)
	return gox.OFromErr(req, err)
}

// BindJSONO 绑定 JSON 并返回 Optional。
func BindJSONO[T any](c *gin.Context) gox.Optional[*T] {
	req, err := BindJSON[T](c)
	return gox.OFromErr(req, err)
}

// BindQueryO 绑定查询参数并返回 Optional。
func BindQueryO[T any](c *gin.Context) gox.Optional[*T] {
	req, err := BindQuery[T](c)
	return gox.OFromErr(req, err)
}

// BindURIO 绑定 URI 并返回 Optional。
func BindURIO[T any](c *gin.Context) gox.Optional[*T] {
	req, err := BindURI[T](c)
	return gox.OFromErr(req, err)
}
