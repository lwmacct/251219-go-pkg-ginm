package ginm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePage_DefaultsForZeroValues(t *testing.T) {
	page, pageSize := normalizePage(0, 0)
	assert.Equal(t, DefaultPage, page)
	assert.Equal(t, DefaultPageSize, pageSize)
}

func TestNormalizePage_PreservesValidValues(t *testing.T) {
	page, pageSize := normalizePage(5, 50)
	assert.Equal(t, 5, page)
	assert.Equal(t, 50, pageSize)
}

func TestNormalizePage_CapsMaxPageSize(t *testing.T) {
	_, pageSize := normalizePage(1, 200)
	assert.Equal(t, MaxPageSize, pageSize)
}

func TestNormalizePage_HandlesNegativeValues(t *testing.T) {
	page, pageSize := normalizePage(-1, -1)
	assert.Equal(t, DefaultPage, page)
	assert.Equal(t, DefaultPageSize, pageSize)
}

func TestPageQuery_Normalize(t *testing.T) {
	q := &PageQuery{Page: 0, PageSize: 0, Order: ""}
	normalized := q.Normalize()

	assert.Equal(t, DefaultPage, normalized.Page)
	assert.Equal(t, DefaultPageSize, normalized.PageSize)
	assert.Equal(t, "desc", normalized.Order)
}

func TestPageQuery_Offset(t *testing.T) {
	tests := []struct {
		page     int
		pageSize int
		expected int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 10, 20},
		{0, 0, 0}, // normalized to page=1, so offset=0
	}

	for _, tt := range tests {
		q := &PageQuery{Page: tt.page, PageSize: tt.pageSize}
		assert.Equal(t, tt.expected, q.Offset(), "page=%d, pageSize=%d", tt.page, tt.pageSize)
	}
}

func TestPageQuery_Limit(t *testing.T) {
	q := &PageQuery{Page: 1, PageSize: 50}
	assert.Equal(t, 50, q.Limit())

	q = &PageQuery{Page: 1, PageSize: 0}
	assert.Equal(t, DefaultPageSize, q.Limit())
}

func TestNewPaginator(t *testing.T) {
	p := NewPaginator[string](2, 25)
	assert.Equal(t, 2, p.Page())
	assert.Equal(t, 25, p.PageSize())
	assert.Equal(t, 25, p.Offset())
	assert.Equal(t, 25, p.Limit())
}

func TestNewPaginator_NormalizesInvalidValues(t *testing.T) {
	p := NewPaginator[string](0, 0)
	assert.Equal(t, DefaultPage, p.Page())
	assert.Equal(t, DefaultPageSize, p.PageSize())
}

func TestPaginator_Paginate(t *testing.T) {
	p := NewPaginator[int](1, 10)
	items := []int{1, 2, 3, 4, 5}
	resp := p.Paginate(items, 100)

	assert.Equal(t, items, resp.Items)
	assert.Equal(t, int64(100), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, 10, resp.TotalPages)
	assert.True(t, resp.HasMore)
}

func TestPaginateSlice_FirstPage(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	resp := PaginateSlice(items, 1, 3)

	assert.Equal(t, []int{1, 2, 3}, resp.Items)
	assert.Equal(t, int64(10), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 3, resp.PageSize)
	assert.Equal(t, 4, resp.TotalPages)
	assert.True(t, resp.HasMore)
}

func TestPaginateSlice_MiddlePage(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	resp := PaginateSlice(items, 2, 3)

	assert.Equal(t, []int{4, 5, 6}, resp.Items)
	assert.Equal(t, 2, resp.Page)
}

func TestPaginateSlice_LastPage(t *testing.T) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	resp := PaginateSlice(items, 4, 3)

	assert.Equal(t, []int{10}, resp.Items)
	assert.False(t, resp.HasMore)
}

func TestPaginateSlice_BeyondRange(t *testing.T) {
	items := []int{1, 2, 3}
	resp := PaginateSlice(items, 10, 10)

	assert.Empty(t, resp.Items)
	assert.Equal(t, int64(3), resp.Total)
}

func TestPaginateSlice_EmptySlice(t *testing.T) {
	var items []int
	resp := PaginateSlice(items, 1, 10)

	assert.Empty(t, resp.Items)
	assert.Equal(t, int64(0), resp.Total)
	assert.False(t, resp.HasMore)
}
