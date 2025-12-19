package ginm

// 默认分页常量。
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// normalizePage 将 page 和 pageSize 规范化为有效值。
// 这是分页默认值的唯一真相来源。
func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}

// PageQuery 是标准的分页查询结构体。
type PageQuery struct {
	Page     int    `form:"page" binding:"min=0"`
	PageSize int    `form:"page_size" binding:"min=0,max=100"`
	Sort     string `form:"sort"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// Normalize 返回应用默认值后的 PageQuery。
func (q *PageQuery) Normalize() PageQuery {
	page, pageSize := normalizePage(q.Page, q.PageSize)

	order := q.Order
	if order == "" {
		order = "desc"
	}

	return PageQuery{
		Page:     page,
		PageSize: pageSize,
		Sort:     q.Sort,
		Order:    order,
	}
}

// Offset 返回数据库偏移量。
func (q *PageQuery) Offset() int {
	page, pageSize := normalizePage(q.Page, q.PageSize)
	return (page - 1) * pageSize
}

// Limit 返回数据库限制数。
func (q *PageQuery) Limit() int {
	_, pageSize := normalizePage(q.Page, q.PageSize)
	return pageSize
}

// Paginator 处理特定类型的分页逻辑。
type Paginator[T any] struct {
	page     int
	pageSize int
}

// NewPaginator 创建新的分页器。
func NewPaginator[T any](page, pageSize int) *Paginator[T] {
	page, pageSize = normalizePage(page, pageSize)
	return &Paginator[T]{
		page:     page,
		pageSize: pageSize,
	}
}

// NewPaginatorFromQuery 从 PageQuery 创建分页器。
func NewPaginatorFromQuery[T any](q *PageQuery) *Paginator[T] {
	normalized := q.Normalize()
	return NewPaginator[T](normalized.Page, normalized.PageSize)
}

// Paginate 根据元素列表和总数创建 PageResponse。
func (p *Paginator[T]) Paginate(items []T, total int64) PageResponse[T] {
	return NewPageResponse(items, total, p.page, p.pageSize)
}

// Offset 返回数据库偏移量。
func (p *Paginator[T]) Offset() int {
	return (p.page - 1) * p.pageSize
}

// Limit 返回数据库限制数。
func (p *Paginator[T]) Limit() int {
	return p.pageSize
}

// Page 返回当前页码。
func (p *Paginator[T]) Page() int {
	return p.page
}

// PageSize 返回页面大小。
func (p *Paginator[T]) PageSize() int {
	return p.pageSize
}

// PaginateSlice 对内存中的切片进行分页。
func PaginateSlice[T any](items []T, page, pageSize int) PageResponse[T] {
	page, pageSize = normalizePage(page, pageSize)
	total := int64(len(items))

	start := (page - 1) * pageSize
	if start >= len(items) {
		return NewPageResponse([]T{}, total, page, pageSize)
	}

	end := min(start+pageSize, len(items))

	return NewPageResponse(items[start:end], total, page, pageSize)
}
