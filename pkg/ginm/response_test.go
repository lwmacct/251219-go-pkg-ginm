package ginm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOK_CreatesSuccessResponse(t *testing.T) {
	resp := OK("hello")
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "hello", resp.Data)
	assert.Empty(t, resp.Message)
	assert.Empty(t, resp.Error)
}

func TestOKWithMessage_IncludesMessage(t *testing.T) {
	resp := OKWithMessage("success", 42)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.Equal(t, 42, resp.Data)
}

func TestFail_CreatesErrorResponse(t *testing.T) {
	resp := Fail[string](400, "bad request")
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "bad request", resp.Message)
}

func TestFailWithError_IncludesErrorDetail(t *testing.T) {
	resp := FailWithError[string](500, "failed", "internal error")
	assert.Equal(t, 500, resp.Code)
	assert.Equal(t, "failed", resp.Message)
	assert.Equal(t, "internal error", resp.Error)
}

func TestNewPageResponse_CalculatesFields(t *testing.T) {
	items := []string{"a", "b", "c"}
	resp := NewPageResponse(items, 100, 2, 10)

	assert.Equal(t, items, resp.Items)
	assert.Equal(t, int64(100), resp.Total)
	assert.Equal(t, 2, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, 10, resp.TotalPages)
	assert.True(t, resp.HasMore)
}

func TestNewPageResponse_LastPage(t *testing.T) {
	resp := NewPageResponse([]int{1}, 10, 10, 1)
	assert.False(t, resp.HasMore)
}

func TestNewPageResponse_NilItemsBecomesEmptySlice(t *testing.T) {
	resp := NewPageResponse[int](nil, 0, 1, 10)
	assert.NotNil(t, resp.Items)
	assert.Empty(t, resp.Items)
}

func TestNewPageResponse_ZeroPageSize(t *testing.T) {
	resp := NewPageResponse([]int{1}, 10, 1, 0)
	assert.Equal(t, 0, resp.TotalPages)
}

func TestNewListResponse_CalculatesCount(t *testing.T) {
	items := []string{"a", "b", "c"}
	resp := NewListResponse(items)

	assert.Equal(t, items, resp.Items)
	assert.Equal(t, 3, resp.Count)
}

func TestNewListResponse_NilItemsBecomesEmptySlice(t *testing.T) {
	resp := NewListResponse[int](nil)
	assert.NotNil(t, resp.Items)
	assert.Empty(t, resp.Items)
	assert.Equal(t, 0, resp.Count)
}

func TestNewListResponse_EmptySlice(t *testing.T) {
	resp := NewListResponse([]int{})
	assert.NotNil(t, resp.Items)
	assert.Empty(t, resp.Items)
	assert.Equal(t, 0, resp.Count)
}
