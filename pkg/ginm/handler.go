package ginm

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandlerFunc 是泛型处理器函数类型。
// Req: 请求类型（如果不需要请求体可使用 struct{}）
// Resp: 响应数据类型
type HandlerFunc[Req, Resp any] func(c *gin.Context, req *Req) (Resp, error)

// HandlerFuncNoReq 是不需要请求绑定的处理器。
type HandlerFuncNoReq[Resp any] func(c *gin.Context) (Resp, error)

// Wrap 将泛型处理器转换为 gin.HandlerFunc，自动绑定请求。
// 根据 Content-Type 自动选择绑定方式（JSON、XML、Form 等）。
func Wrap[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := Bind[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapJSON 将泛型处理器转换为 gin.HandlerFunc，使用 JSON 绑定。
func WrapJSON[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindJSON[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapQuery 将泛型处理器转换为 gin.HandlerFunc，使用查询参数绑定。
func WrapQuery[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindQuery[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapURI 将泛型处理器转换为 gin.HandlerFunc，使用 URI 绑定。
func WrapURI[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindURI[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapNoReq 将不需要请求绑定的处理器转换为 gin.HandlerFunc。
func WrapNoReq[Resp any](handler HandlerFuncNoReq[Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := handler(c)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapPage 将分页处理器转换为 gin.HandlerFunc。
func WrapPage[Req any, Item any](handler func(c *gin.Context, req *Req) (PageResponse[Item], error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindQuery[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapURIAndJSON 将同时使用 URI 和 JSON 绑定的处理器转换为 gin.HandlerFunc。
// 适用于 PUT /users/:id 带 JSON body 的路由。
func WrapURIAndJSON[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindURIAndBody[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	}
}

// WrapWithStatus 将泛型处理器转换为 gin.HandlerFunc，使用自定义成功状态码。
func WrapWithStatus[Req, Resp any](handler HandlerFunc[Req, Resp], successStatus int) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := Bind[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(successStatus, OK(resp))
	}
}

// WrapCreated 包装返回 HTTP 201 Created 的处理器。
func WrapCreated[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return WrapWithStatus(handler, http.StatusCreated)
}

// WrapCreatedJSON 包装使用 JSON 绑定并返回 HTTP 201 Created 的处理器。
func WrapCreatedJSON[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindJSON[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := handler(c, req)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, OK(resp))
	}
}

// WrapAccepted 包装返回 HTTP 202 Accepted 的处理器。
func WrapAccepted[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return WrapWithStatus(handler, http.StatusAccepted)
}

// WrapNoContent 包装返回 HTTP 204 No Content 的处理器。
func WrapNoContent[Req any](handler func(c *gin.Context, req *Req) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := Bind[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		if err := handler(c, req); err != nil {
			handleError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// WrapNoContentJSON 包装使用 JSON 绑定并返回 HTTP 204 的处理器。
func WrapNoContentJSON[Req any](handler func(c *gin.Context, req *Req) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindJSON[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		if err := handler(c, req); err != nil {
			handleError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// --- 常用 HTTP 方法的便捷处理器 ---

// HandleGet 包装不需要请求绑定的 GET 处理器。
// 是 WrapNoReq 的简写。
func HandleGet[Resp any](handler func(c *gin.Context) (Resp, error)) gin.HandlerFunc {
	return WrapNoReq(handler)
}

// HandleGetWithQuery 包装带查询参数绑定的 GET 处理器。
// 是 WrapQuery 的简写。
func HandleGetWithQuery[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return WrapQuery(handler)
}

// HandlePost 包装使用 JSON 绑定的 POST 处理器。
// 是 WrapJSON 的简写。
func HandlePost[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return WrapJSON(handler)
}

// HandlePut 包装使用 URI 和 JSON 绑定的 PUT 处理器。
// 是 WrapURIAndJSON 的简写。
func HandlePut[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return WrapURIAndJSON(handler)
}

// HandleDelete 包装只返回错误的 DELETE 处理器。
func HandleDelete(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			handleError(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// DeleteWithURI 包装带 URI 绑定的 DELETE 处理器。
func DeleteWithURI[Req any](handler func(c *gin.Context, req *Req) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := BindURI[Req](c)
		if err != nil {
			handleError(c, err)
			return
		}

		if err := handler(c, req); err != nil {
			handleError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// HandlePatch 包装使用 URI 和 JSON 绑定的 PATCH 处理器。
func HandlePatch[Req, Resp any](handler HandlerFunc[Req, Resp]) gin.HandlerFunc {
	return WrapURIAndJSON(handler)
}

// handleError 处理错误并发送适当的 HTTP 响应。
func handleError(c *gin.Context, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		errStr := ""
		if apiErr.Err != nil && gin.Mode() != gin.ReleaseMode {
			errStr = apiErr.Err.Error()
		}
		c.JSON(apiErr.HTTPStatus, FailWithError[any](apiErr.Code, apiErr.Message, errStr))
		return
	}

	var bindErr *BindError
	if errors.As(err, &bindErr) {
		c.JSON(http.StatusBadRequest, Fail[any](http.StatusBadRequest, bindErr.Error()))
		return
	}

	var validationErrs *ValidationErrors
	if errors.As(err, &validationErrs) {
		c.JSON(http.StatusUnprocessableEntity, Response[*ValidationErrors]{
			Code:    http.StatusUnprocessableEntity,
			Message: "validation failed",
			Data:    validationErrs,
		})
		return
	}

	// 默认: 内部服务器错误
	errStr := ""
	if gin.Mode() != gin.ReleaseMode {
		errStr = err.Error()
	}
	c.JSON(http.StatusInternalServerError, FailWithError[any](
		http.StatusInternalServerError,
		"internal server error",
		errStr,
	))
}
