package ginm_test

import (
	"errors"
	"fmt"

	"github.com/lwmacct/251219-go-pkg-ginm/pkg/ginm"
)

func ExampleOK() {
	resp := ginm.OK("hello world")
	fmt.Println(resp.Code)
	fmt.Println(resp.Data)
	// Output:
	// 0
	// hello world
}

func ExampleOKWithMessage() {
	resp := ginm.OKWithMessage("created successfully", map[string]int{"id": 1})
	fmt.Println(resp.Code)
	fmt.Println(resp.Message)
	// Output:
	// 0
	// created successfully
}

func ExampleFail() {
	resp := ginm.Fail[any](400, "invalid request")
	fmt.Println(resp.Code)
	fmt.Println(resp.Message)
	// Output:
	// 400
	// invalid request
}

func ExampleNewPageResponse() {
	items := []string{"a", "b", "c"}
	resp := ginm.NewPageResponse(items, 100, 1, 10)
	fmt.Println("Items:", resp.Items)
	fmt.Println("Total:", resp.Total)
	fmt.Println("TotalPages:", resp.TotalPages)
	fmt.Println("HasMore:", resp.HasMore)
	// Output:
	// Items: [a b c]
	// Total: 100
	// TotalPages: 10
	// HasMore: true
}

func ExampleNewListResponse() {
	items := []int{1, 2, 3, 4, 5}
	resp := ginm.NewListResponse(items)
	fmt.Println("Count:", resp.Count)
	// Output:
	// Count: 5
}

func ExamplePaginateSlice() {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	resp := ginm.PaginateSlice(items, 2, 3)
	fmt.Println("Items:", resp.Items)
	fmt.Println("Page:", resp.Page)
	// Output:
	// Items: [4 5 6]
	// Page: 2
}

func ExampleNewPaginator() {
	p := ginm.NewPaginator[string](2, 25)
	fmt.Println("Offset:", p.Offset())
	fmt.Println("Limit:", p.Limit())
	// Output:
	// Offset: 25
	// Limit: 25
}

func ExamplePageQuery_Normalize() {
	q := &ginm.PageQuery{Page: 0, PageSize: 0}
	normalized := q.Normalize()
	fmt.Println("Page:", normalized.Page)
	fmt.Println("PageSize:", normalized.PageSize)
	fmt.Println("Order:", normalized.Order)
	// Output:
	// Page: 1
	// PageSize: 20
	// Order: desc
}

func ExampleErrBadRequest() {
	err := ginm.ErrBadRequest("missing required field")
	fmt.Println(err.HTTPStatus)
	fmt.Println(err.Message)
	// Output:
	// 400
	// missing required field
}

func ExampleErrNotFound() {
	err := ginm.ErrNotFound("user not found")
	fmt.Println(err.HTTPStatus)
	// Output:
	// 404
}

func ExampleNewBindError() {
	err := ginm.NewBindError("json", errors.New("invalid syntax"))
	fmt.Println(err.Source)
	// Output:
	// json
}

func ExampleValidationErrors() {
	ve := &ginm.ValidationErrors{}
	ve.Add("email", "invalid format")
	ve.Add("age", "must be positive")
	fmt.Println(ve.HasErrors())
	fmt.Println(len(ve.Errors))
	// Output:
	// true
	// 2
}
