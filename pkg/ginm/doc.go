// Package ginm 提供基于 Go 泛型的 Gin Web 框架类型安全工具。
//
// 核心功能：
//   - 泛型响应类型: Response[T], PageResponse[T]
//   - 类型安全请求绑定: Bind[T], BindJSON[T], BindQuery[T]
//   - 泛型处理器包装: Wrap, WrapJSON, WrapQuery
//   - 类型安全上下文工具: ContextKey[T], Get[T], Set[T]
//   - RESTful 资源接口，支持完整 CRUD
//
// 使用示例：
//
//	type CreateUserReq struct {
//	    Name  string `json:"name" binding:"required"`
//	    Email string `json:"email" binding:"required,email"`
//	}
//
//	func CreateUser(c *gin.Context, req *CreateUserReq) (*User, error) {
//	    return &User{ID: 1, Name: req.Name, Email: req.Email}, nil
//	}
//
//	r.POST("/users", ginm.WrapJSON(CreateUser))
package ginm
