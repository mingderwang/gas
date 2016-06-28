package gas

import (
	"net/http"
	"testing"
	"github.com/gavv/httpexpect"
	"encoding/json"
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
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.Render(jsonMap, "testfiles/layout.html", "testfiles/index.html")
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client:   httpexpect.NewFastBinder(handler),
	})

	// run tests
	e.GET("/").
	Expect().
	Status(http.StatusOK).
	ContentType("text/html", "utf-8").
	Body().Equal(testHTML)

}

func TestHTML(t *testing.T) {
	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.HTML(200, testHTML)
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client:   httpexpect.NewFastBinder(handler),
	})

	// run tests
	e.GET("/").
	Expect().
	Status(http.StatusOK).
	ContentType("text/html", "utf-8").
	Body().Equal(testHTML)

}

func TestSTRINGResponse(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.STRING(200, tstr)
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client:   httpexpect.NewFastBinder(handler),
	})

	// run tests
	e.GET("/").
	Expect().
	Status(http.StatusOK).
	ContentType("text/plain", "utf-8").
	Body().Equal(tstr)


}

func TestJSONResponse(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	// set route
	g.Router.Get("/", func(ctx *Context) error {
		return ctx.JSON(200, jsonMap)
	})

	// create fasthttp.RequestHandler
	handler := g.Router.Handler

	// create httpexpect instance that will call fasthtpp.RequestHandler directly
	e := httpexpect.WithConfig(httpexpect.Config{
		Reporter: httpexpect.NewAssertReporter(t),
		Client:   httpexpect.NewFastBinder(handler),
	})

	js, _ := json.Marshal(jsonMap)

	// run tests
	e.GET("/").
	Expect().
	Status(http.StatusOK).
	ContentType("application/json", "utf-8").
	Body().Equal(string(js))

}
