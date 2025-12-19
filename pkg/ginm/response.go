package ginm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是带泛型数据类型的标准 API 响应包装器。
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// OK 创建带数据的成功响应。
func OK[T any](data T) Response[T] {
	return Response[T]{
		Code: 0,
		Data: data,
	}
}

// OKWithMessage 创建带消息和数据的成功响应。
func OKWithMessage[T any](message string, data T) Response[T] {
	return Response[T]{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

// Fail 创建错误响应。
func Fail[T any](code int, message string) Response[T] {
	return Response[T]{
		Code:    code,
		Message: message,
	}
}

// FailWithError 创建带错误详情的错误响应。
func FailWithError[T any](code int, message string, err string) Response[T] {
	return Response[T]{
		Code:    code,
		Message: message,
		Error:   err,
	}
}

// PageResponse 表示带泛型元素类型的分页数据。
type PageResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	HasMore    bool  `json:"has_more"`
}

// NewPageResponse 创建新的分页响应。
func NewPageResponse[T any](items []T, total int64, page, pageSize int) PageResponse[T] {
	if items == nil {
		items = []T{}
	}
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	return PageResponse[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}
}

// ListResponse 用于非分页列表。
type ListResponse[T any] struct {
	Items []T `json:"items"`
	Count int `json:"count"`
}

// NewListResponse 创建列表响应。
func NewListResponse[T any](items []T) ListResponse[T] {
	if items == nil {
		items = []T{}
	}
	return ListResponse[T]{
		Items: items,
		Count: len(items),
	}
}

// JSON 发送带指定状态码的 JSON 响应。
func JSON[T any](c *gin.Context, status int, resp Response[T]) {
	c.JSON(status, resp)
}

// Success 发送 HTTP 200 的成功 JSON 响应。
func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, OK(data))
}

// SuccessWithMessage 发送带消息的成功 JSON 响应。
func SuccessWithMessage[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, OKWithMessage(message, data))
}

// SuccessPage 发送分页成功响应。
func SuccessPage[T any](c *gin.Context, items []T, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, OK(NewPageResponse(items, total, page, pageSize)))
}

// SuccessList 发送列表成功响应。
func SuccessList[T any](c *gin.Context, items []T) {
	c.JSON(http.StatusOK, OK(NewListResponse(items)))
}

// Error 发送错误 JSON 响应。
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Fail[any](code, message))
}

// ErrorWithDetail 发送带错误详情的错误 JSON 响应。
func ErrorWithDetail(c *gin.Context, httpStatus int, code int, message string, err error) {
	errStr := ""
	if err != nil && gin.Mode() != gin.ReleaseMode {
		errStr = err.Error()
	}
	c.JSON(httpStatus, FailWithError[any](code, message, errStr))
}

// --- 其他 HTTP 状态响应 ---

// NoContent 发送 HTTP 204 No Content 响应。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Created 发送带数据的 HTTP 201 Created。
func Created[T any](c *gin.Context, data T) {
	c.JSON(http.StatusCreated, OK(data))
}

// CreatedWithMessage 发送带消息和数据的 HTTP 201 Created。
func CreatedWithMessage[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusCreated, OKWithMessage(message, data))
}

// Accepted 发送带数据的 HTTP 202 Accepted。
func Accepted[T any](c *gin.Context, data T) {
	c.JSON(http.StatusAccepted, OK(data))
}

// Redirect 发送重定向响应。
func Redirect(c *gin.Context, code int, url string) {
	c.Redirect(code, url)
}

// RedirectPermanent 发送 HTTP 301 永久重定向。
func RedirectPermanent(c *gin.Context, url string) {
	c.Redirect(http.StatusMovedPermanently, url)
}

// RedirectTemporary 发送 HTTP 302 临时重定向。
func RedirectTemporary(c *gin.Context, url string) {
	c.Redirect(http.StatusFound, url)
}

// --- 文件响应 ---

// File 发送文件响应。
func File(c *gin.Context, filepath string) {
	c.File(filepath)
}

// FileAttachment 发送文件作为附件（下载）。
func FileAttachment(c *gin.Context, filepath, filename string) {
	c.FileAttachment(filepath, filename)
}

// DataAttachment 发送原始数据作为文件附件。
func DataAttachment(c *gin.Context, data []byte, filename string) {
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// --- 流式响应 ---

// Stream 通过 channel 发送数据，使用 JSON 行格式（NDJSON）。
// 序列化失败的元素会被跳过。
func Stream[T any](c *gin.Context, ch <-chan T) {
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		if item, ok := <-ch; ok {
			data, err := json.Marshal(item)
			if err != nil {
				// 跳过序列化失败的元素
				return true
			}
			w.Write(data)
			w.Write([]byte("\n"))
			return true
		}
		return false
	})
}

// StreamWithError 通过 channel 发送数据，遇到序列化错误时停止。
func StreamWithError[T any](c *gin.Context, ch <-chan T) error {
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Transfer-Encoding", "chunked")

	var marshalErr error
	c.Stream(func(w io.Writer) bool {
		if item, ok := <-ch; ok {
			data, err := json.Marshal(item)
			if err != nil {
				marshalErr = err
				return false
			}
			w.Write(data)
			w.Write([]byte("\n"))
			return true
		}
		return false
	})
	return marshalErr
}

// SSE 发送 Server-Sent Events。
// 序列化失败的元素会被跳过。
func SSE[T any](c *gin.Context, ch <-chan T) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		if item, ok := <-ch; ok {
			data, err := json.Marshal(item)
			if err != nil {
				// 跳过序列化失败的元素
				return true
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			return true
		}
		return false
	})
}

// SSEWithEvent 发送带自定义事件类型的 Server-Sent Events。
// 序列化失败的元素会被跳过。
func SSEWithEvent[T any](c *gin.Context, eventType string, ch <-chan T) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		if item, ok := <-ch; ok {
			data, err := json.Marshal(item)
			if err != nil {
				// 跳过序列化失败的元素
				return true
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
			return true
		}
		return false
	})
}

// --- 原始响应 ---

// Raw 发送带自定义 Content-Type 的原始数据。
func Raw(c *gin.Context, contentType string, data []byte) {
	c.Data(http.StatusOK, contentType, data)
}

// RawWithStatus 发送带自定义状态码和 Content-Type 的原始数据。
func RawWithStatus(c *gin.Context, status int, contentType string, data []byte) {
	c.Data(status, contentType, data)
}

// String 发送纯文本响应。
func String(c *gin.Context, format string, values ...any) {
	c.String(http.StatusOK, format, values...)
}

// HTML 发送 HTML 响应。
func HTML(c *gin.Context, html string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// XML 发送 XML 响应。
func XML[T any](c *gin.Context, data T) {
	c.XML(http.StatusOK, data)
}
