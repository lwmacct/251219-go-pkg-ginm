package ginm

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testRequest struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Email string `json:"email" form:"email"`
}

type testQueryParams struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Search   string `form:"search"`
}

func createTestContext(method, path string, body []byte, contentType string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func TestBindJSON_Success(t *testing.T) {
	body := []byte(`{"name": "John", "email": "john@example.com"}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	result, err := BindJSON[testRequest](c)
	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, "john@example.com", result.Email)
}

func TestBindJSON_ValidationError(t *testing.T) {
	body := []byte(`{"email": "john@example.com"}`) // missing required name
	c, _ := createTestContext("POST", "/", body, "application/json")

	_, err := BindJSON[testRequest](c)
	require.Error(t, err)

	var bindErr *BindError
	assert.ErrorAs(t, err, &bindErr)
	assert.Equal(t, "json", bindErr.Source)
}

func TestBindJSON_InvalidJSON(t *testing.T) {
	body := []byte(`{invalid json}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	_, err := BindJSON[testRequest](c)
	require.Error(t, err)
}

func TestBindQuery_Success(t *testing.T) {
	c, _ := createTestContext("GET", "/?page=2&page_size=10&search=test", nil, "")

	result, err := BindQuery[testQueryParams](c)
	require.NoError(t, err)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, "test", result.Search)
}

func TestBindQuery_EmptyParams(t *testing.T) {
	c, _ := createTestContext("GET", "/", nil, "")

	result, err := BindQuery[testQueryParams](c)
	require.NoError(t, err)
	assert.Equal(t, 0, result.Page)
	assert.Equal(t, 0, result.PageSize)
}

func TestBindForm_Success(t *testing.T) {
	body := []byte("name=John&email=john@example.com")
	c, _ := createTestContext("POST", "/", body, "application/x-www-form-urlencoded")

	result, err := BindForm[testRequest](c)
	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
}

func TestMustBindJSON_Panics(t *testing.T) {
	body := []byte(`{"email": "john@example.com"}`) // missing required name
	c, _ := createTestContext("POST", "/", body, "application/json")

	assert.Panics(t, func() {
		MustBindJSON[testRequest](c)
	})
}

func TestMustBindJSON_Success(t *testing.T) {
	body := []byte(`{"name": "John"}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	assert.NotPanics(t, func() {
		result := MustBindJSON[testRequest](c)
		assert.Equal(t, "John", result.Name)
	})
}

func TestBindJSONR_ReturnsResult(t *testing.T) {
	body := []byte(`{"name": "John"}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	result := BindJSONR[testRequest](c)
	assert.True(t, result.IsOk())
	assert.Equal(t, "John", result.Unwrap().Name)
}

func TestBindJSONR_ReturnsErrorResult(t *testing.T) {
	body := []byte(`{}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	result := BindJSONR[testRequest](c)
	assert.True(t, result.IsErr())
}

func TestBindJSONO_ReturnsOptional(t *testing.T) {
	body := []byte(`{"name": "John"}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	opt := BindJSONO[testRequest](c)
	assert.True(t, opt.IsSome())
	assert.Equal(t, "John", opt.MustGet().Name)
}

func TestBindJSONO_ReturnsNone(t *testing.T) {
	body := []byte(`{}`)
	c, _ := createTestContext("POST", "/", body, "application/json")

	opt := BindJSONO[testRequest](c)
	assert.True(t, opt.IsNone())
}

func TestBindConfig(t *testing.T) {
	cfg := BindConfig{
		URI:   true,
		Query: true,
		Body:  false,
	}
	assert.True(t, cfg.URI)
	assert.True(t, cfg.Query)
	assert.False(t, cfg.Body)
}
