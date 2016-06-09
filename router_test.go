package goslim

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPMethod(t *testing.T) {
	as := assert.New(t)

	// new goslim
	g := New()

	g.Router.Get("/test", func(c *Context) error {
		return c.STRING(200, "TEST")
	})

	g.Router.Post("/test", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
	})

	g.Router.Put("/test", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
	})

	g.Router.Patch("/", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
	})

	g.Router.Delete("/", func(c *Context) error {
		return c.STRING(200, "Deleted")
	})

	g.Router.Options("/", func(c *Context) error {
		return c.STRING(200, "Option")
	})

	g.Router.Head("/", func(c *Context) error {
		return c.STRING(200, "Head")
	})

	r, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
	as.Equal("TEST", w.Body.String())

	r, _ = http.NewRequest("POST", "/test", nil)
	r.ParseForm()
	r.Form.Add("Test", "POSTED")
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
	as.Equal("POSTED", w.Body.String())

	r, _ = http.NewRequest("PUT", "/test", nil)
	r.ParseForm()
	r.Form.Add("Test", "PUTED")
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
	as.Equal("PUTED", w.Body.String())

	r, _ = http.NewRequest("PATCH", "/", nil)
	r.ParseForm()
	r.Form.Add("Test", "PATCHED")
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
	as.Equal("PATCHED", w.Body.String())

	r, _ = http.NewRequest("DELETE", "/", nil)
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)

	r, _ = http.NewRequest("OPTIONS", "/", nil)
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)

	r, _ = http.NewRequest("HEAD", "/", nil)
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
}

type testController struct {
	ControllerInterface
}

func (cn *testController) Get(c *Context) error {
	return c.STRING(200, "Get Test")
}
func (cn *testController) Post(c *Context) error {
	return c.STRING(200, "Post Test"+c.GetParam("Test"))
}

func TestSetAllRESTFunc(t *testing.T) {
	as := assert.New(t)

	var c = &testController{}

	// new goslim
	g := New()

	g.Router.REST("/User", c)

	r, _ := http.NewRequest("GET", "/User", nil)
	w := httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
	as.Equal("Get Test", w.Body.String())

	r, _ = http.NewRequest("POST", "/User", nil)
	r.ParseForm()
	r.Form.Add("Test", "POSTED")
	w = httptest.NewRecorder()
	g.Router.ServeHTTP(w, r)
	as.Equal(200, w.Code)
	as.Equal("Post Test"+"POSTED", w.Body.String())

}
