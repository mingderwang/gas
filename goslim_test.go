package goslim

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	indexString = "indexpage"
)

func TestGoslim(t *testing.T) {
	assert := assert.New(t)

	// new goslim
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", indexPage)

	// create request
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	g.Router.ServeHTTP(rec, req)

	assert.Equal(200, rec.Code)
	assert.Equal(indexString, rec.Body.String())

}

func indexPage(ctx *Context) error {
	return ctx.STRING(200, indexString)
}
