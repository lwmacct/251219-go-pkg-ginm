package ginm

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IDParam 表示资源 ID 的 URI 参数。
type IDParam[ID any] struct {
	ID ID `binding:"required" uri:"id"`
}

// Resource 是 RESTful 资源接口，LQ 为列表查询参数类型
type Resource[T any, ID comparable, CI any, UI any, LQ any] interface {
	// List 返回分页的元素列表。
	List(c *gin.Context, query *LQ) (PageResponse[T], error)

	// Get 根据 ID 返回单个元素。
	Get(c *gin.Context, id ID) (*T, error)

	// Create 创建新元素。
	Create(c *gin.Context, input *CI) (*T, error)

	// Update 更新现有元素。
	Update(c *gin.Context, id ID, input *UI) (*T, error)

	// Delete 根据 ID 删除元素。
	Delete(c *gin.Context, id ID) error
}

// BaseResource 提供返回"未实现"错误的默认实现。
// 嵌入此结构体并仅覆盖你需要的方法。
type BaseResource[T any, ID comparable, CI any, UI any, LQ any] struct{}

func (r *BaseResource[T, ID, CI, UI, LQ]) List(c *gin.Context, query *LQ) (PageResponse[T], error) {
	return PageResponse[T]{}, ErrNotImplemented("List")
}

func (r *BaseResource[T, ID, CI, UI, LQ]) Get(c *gin.Context, id ID) (*T, error) {
	return nil, ErrNotImplemented("Get")
}

func (r *BaseResource[T, ID, CI, UI, LQ]) Create(c *gin.Context, input *CI) (*T, error) {
	return nil, ErrNotImplemented("Create")
}

func (r *BaseResource[T, ID, CI, UI, LQ]) Update(c *gin.Context, id ID, input *UI) (*T, error) {
	return nil, ErrNotImplemented("Update")
}

func (r *BaseResource[T, ID, CI, UI, LQ]) Delete(c *gin.Context, id ID) error {
	return ErrNotImplemented("Delete")
}

// ResourceConfig 包含资源注册的配置选项。
type ResourceConfig struct {
	// IDParam 是 URI 中 ID 参数的名称。默认值: "id"
	IDParam string
}

// ResourceOption 是资源注册的函数式选项。
type ResourceOption func(*ResourceConfig)

// WithIDParam 设置自定义的 ID 参数名称。
func WithIDParam(name string) ResourceOption {
	return func(cfg *ResourceConfig) {
		cfg.IDParam = name
	}
}

// RegisterResource 为资源注册所有 CRUD 路由。
// 创建的路由:
//   - GET    /           -> List
//   - GET    /:id        -> Get
//   - POST   /           -> Create
//   - PUT    /:id        -> Update
//   - DELETE /:id        -> Delete
func RegisterResource[T any, ID comparable, CI any, UI any, LQ any](
	group *gin.RouterGroup,
	resource Resource[T, ID, CI, UI, LQ],
	opts ...ResourceOption,
) {
	cfg := &ResourceConfig{IDParam: "id"}
	for _, opt := range opts {
		opt(cfg)
	}

	idPath := "/:" + cfg.IDParam

	// GET / - 列表
	group.GET("", func(c *gin.Context) {
		query, err := BindQuery[LQ](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := resource.List(c, query)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	})

	// GET /:id - 获取
	group.GET(idPath, func(c *gin.Context) {
		idParam, err := BindURI[IDParam[ID]](c)
		if err != nil {
			handleError(c, err)
			return
		}

		item, err := resource.Get(c, idParam.ID)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(item))
	})

	// POST / - 创建
	group.POST("", func(c *gin.Context) {
		input, err := BindJSON[CI](c)
		if err != nil {
			handleError(c, err)
			return
		}

		item, err := resource.Create(c, input)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, OK(item))
	})

	// PUT /:id - 更新
	group.PUT(idPath, func(c *gin.Context) {
		idParam, err := BindURI[IDParam[ID]](c)
		if err != nil {
			handleError(c, err)
			return
		}

		input, err := BindJSON[UI](c)
		if err != nil {
			handleError(c, err)
			return
		}

		item, err := resource.Update(c, idParam.ID, input)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(item))
	})

	// DELETE /:id - 删除
	group.DELETE(idPath, func(c *gin.Context) {
		idParam, err := BindURI[IDParam[ID]](c)
		if err != nil {
			handleError(c, err)
			return
		}

		if err := resource.Delete(c, idParam.ID); err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK[any](nil))
	})
}

// RegisterResourceReadOnly 仅注册只读路由（List 和 Get）。
func RegisterResourceReadOnly[T any, ID comparable, CI any, UI any, LQ any](
	group *gin.RouterGroup,
	resource Resource[T, ID, CI, UI, LQ],
	opts ...ResourceOption,
) {
	cfg := &ResourceConfig{IDParam: "id"}
	for _, opt := range opts {
		opt(cfg)
	}

	idPath := "/:" + cfg.IDParam

	// GET / - 列表
	group.GET("", func(c *gin.Context) {
		query, err := BindQuery[LQ](c)
		if err != nil {
			handleError(c, err)
			return
		}

		resp, err := resource.List(c, query)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(resp))
	})

	// GET /:id - 获取
	group.GET(idPath, func(c *gin.Context) {
		idParam, err := BindURI[IDParam[ID]](c)
		if err != nil {
			handleError(c, err)
			return
		}

		item, err := resource.Get(c, idParam.ID)
		if err != nil {
			handleError(c, err)
			return
		}

		c.JSON(http.StatusOK, OK(item))
	})
}
