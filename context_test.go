package goslim

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	jsonMap = map[string]string{
		"Test": "index page",
	}

	tstr = "Test String"

	testHTML = `<html>
    <head>
        <title>index page</title>
    </head>
    
    <body>
        <b>This is index page</b>
    </body>
</html>`
)

func TestRender(t *testing.T) {
	as := assert.New(t)

	// new goslim
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.Render(jsonMap, "testfiles/layout.html", "testfiles/index.html")
	})

	// create request
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	g.Router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal(testHTML, rec.Body.String())
	as.Equal(TextHTMLCharsetUTF8, rec.Header().Get(ContentType))
}

func TestHTML(t *testing.T) {
	as := assert.New(t)

	// new goslim
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.HTML(200, testHTML)
	})

	// create request
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	g.Router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal(testHTML, rec.Body.String())
	as.Equal(TextHTMLCharsetUTF8, rec.Header().Get(ContentType))
}

func TestSTRINGResponse(t *testing.T) {
	as := assert.New(t)

	// new goslim
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.STRING(200, tstr)
	})

	// create request
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	g.Router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	as.Equal(tstr, rec.Body.String())
	as.Equal(TextPlainCharsetUTF8, rec.Header().Get(ContentType))
}

func TestJSONResponse(t *testing.T) {
	as := assert.New(t)

	// new goslim
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.JSON(200, jsonMap)
	})

	// create request
	req, _ := http.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	g.Router.ServeHTTP(rec, req)

	as.Equal(200, rec.Code)
	js, _ := json.Marshal(jsonMap)
	as.Equal(string(js), rec.Body.String())
	as.Equal(ApplicationJSONCharsetUTF8, rec.Header().Get(ContentType))
}
