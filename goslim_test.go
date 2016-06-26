package gas

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	indexString = "indexpage"
	testStaticString = "This is a static file"
)

func Testgas(t *testing.T) {
	assert := assert.New(t)

	// new gas
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

func Testgas_Static(t *testing.T) {
	assert := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	//g.Router.Get("/", indexPage)

	// create request
	req, _ := http.NewRequest("GET", "/testfiles/static.txt", nil)
	rec := httptest.NewRecorder()
	g.Router.ServeHTTP(rec, req)

	assert.Equal(200, rec.Code)
	assert.Equal(testStaticString, rec.Body.String())
}

func indexPage(ctx *Context) error {
	return ctx.STRING(200, indexString)
}
