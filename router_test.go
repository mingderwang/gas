package gas

import (
	"testing"
	"net/http"
)

func TestRouter_Get(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Get("/test", func(c *Context) error {
		return c.STRING(200, "TEST")
	})

	e := newHttpExpect(t, g.Router.Handler)
	e.GET("/test").Expect().Status(http.StatusOK).Body().Equal("TEST")
}

func TestRouter_Post(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Post("/test", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.POST("/test").WithFormField("Test", "POSTDATA").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("POSTDATA")
}

func TestRouter_Put(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Put("/test", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.PUT("/test").WithFormField("Test", "POSTDATA").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("POSTDATA")

	//r, _ := http.NewRequest("PUT", "/test", nil)
	//r.ParseForm()
	//r.Form.Add("Test", "PUTED")
	//w := httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
	//as.Equal("PUTED", w.Body.String())
}

func TestRouter_Patch(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Patch("/", func(c *Context) error {
		return c.STRING(200, c.GetParam("Test"))
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.PATCH("/").WithFormField("Test", "POSTDATA").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("POSTDATA")

	//r, _ := http.NewRequest("PATCH", "/", nil)
	//r.ParseForm()
	//r.Form.Add("Test", "PATCHED")
	//w := httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
	//as.Equal("PATCHED", w.Body.String())
}

func TestRouter_Delete(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Delete("/", func(c *Context) error {
		return c.STRING(200, "Deleted")
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.DELETE("/").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("Deleted")

	//r, _ := http.NewRequest("DELETE", "/", nil)
	//w := httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
}

func TestRouter_Options(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Options("/", func(c *Context) error {
		return c.STRING(200, "Option")
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.OPTIONS("/").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("Option")

	//r, _ := http.NewRequest("OPTIONS", "/", nil)
	//w := httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
}

func TestRouter_Head(t *testing.T) {
	//as := assert.New(t)

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.Head("/", func(c *Context) error {
		return c.STRING(200, "Head")
	})

	e := newHttpExpect(t, g.Router.Handler)
	ee := e.HEAD("/").Expect()
	ee.Status(http.StatusOK)
	ee.Body().Equal("Head")

	//r, _ := http.NewRequest("HEAD", "/", nil)
	//w := httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
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

func TestRouter_REST(t *testing.T) {
	//as := assert.New(t)

	var c = &testController{}

	// new gas
	g := New("testfiles/config_test.yaml")

	g.Router.REST("/User", c)

	e := newHttpExpect(t, g.Router.Handler)

	ee1 := e.GET("/User").Expect()
	ee1.Status(http.StatusOK)
	ee1.Body().Equal("Get Test")

	ee2 := e.POST("/User").WithFormField("Test", "POSTED").Expect()
	ee2.Status(http.StatusOK)
	ee2.Body().Equal("Post TestPOSTED")
	//r, _ := http.NewRequest("GET", "/User", nil)
	//w := httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
	//as.Equal("Get Test", w.Body.String())
	//
	//r, _ = http.NewRequest("POST", "/User", nil)
	//r.ParseForm()
	//r.Form.Add("Test", "POSTED")
	//w = httptest.NewRecorder()
	//g.Router.ServeHTTP(w, r)
	//as.Equal(200, w.Code)
	//as.Equal("Post Test"+"POSTED", w.Body.String())

}
