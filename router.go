package gas

import (
	"net/http"
	"reflect"
	"strings"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	//"github.com/julienschmidt/httprouter"
	"github.com/julienschmidt/httprouter"
)

var supportRestProto = [7]string{"GET", "POST", "DELETE", "HEAD", "OPTIONS", "PUT", "PATCH"}

type (

	// Router class include httprouter and gas
	Router struct {
		g           *gas
		hr          fasthttprouter.Router
		middlewares []MiddlewareFunc
	}

	// MiddlewareFunc middlewarefunc define
	MiddlewareFunc func(CHandler) CHandler

	// CHandler is a function type for rout handler
	CHandler func(*Context) error

	// PanicHandler defined panic handler
	PanicHandler func(*Context, interface{}) error
)

// SetNotFoundHandler  set Notfound and Panic handler
func (r *Router) SetNotFoundHandler(h CHandler) {

	r.hr.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
		ctx.reset(w, req, r.g)

		// chain middleware functions
		var cpch CHandler // copy handle avoid repeat chain
		cpch = h

		for i := len(r.middlewares) - 1; i >= 0; i-- {
			cpch = r.middlewares[i](cpch)
		}

		if err := cpch(ctx); err != nil {

		}

		r.g.pool.Put(ctx)
	})
}

func (r *Router) SetPanicHandler(ph PanicHandler) {
	r.hr.PanicHandler = func(w http.ResponseWriter, req *http.Request, rcv interface{}) {
		// c := a.createContext(w, req)
		// a.panicFunc(c, rcv)
		// a.pool.Put(c)

		ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
		ctx.reset(w, req, r.g)

		if err := ph(ctx, rcv); err != nil {

		}

		r.g.pool.Put(ctx)
	}
}

func (r *Router) Use(m interface{}) {
	m = wrapMiddleware(m)

	r.middlewares = append(r.middlewares, m.(MiddlewareFunc))
}

// wrapMiddleware wraps middleware.
func wrapMiddleware(m interface{}) MiddlewareFunc {
	switch m := m.(type) {
	case MiddlewareFunc:
		return m
	case func(CHandler) CHandler:
		return m
	case CHandler:
		return wrapHandlerFuncToMiddlewareFunc(m)
	case func(c *Context) error:
		return wrapHandlerFuncToMiddlewareFunc(m)

	default:
		panic("unknown middleware")
	}
}

func wrapHandlerFuncToMiddlewareFunc(m CHandler) MiddlewareFunc {
	return func(h CHandler) CHandler {
		return func(c *Context) error {
			if err := m(c); err != nil {
				return err
			}

			return h(c)
		}
	}
}

func (r *Router) setRoute(method, path string, ch CHandler) {
	//r.hr.Handle(method, path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	r.hr.Handle(method, path, func(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {

		// ctx := Context{Rw: &rw, Req: req, ps: &ps, handlerFunc: ch}
		ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
		ctx.reset(rw, req, r.g)
		ctx.ps = &ps

		// chain middleware functions
		var cpch CHandler // copy handle avoid repeat chain
		cpch = ch

		for i := len(r.middlewares) - 1; i >= 0; i-- {
			cpch = r.middlewares[i](cpch)
		}

		if err := cpch(ctx); err != nil {
			// handle error
		}

		if ctx.isUseDB {
			defer ctx.CloseDB()
		}

		// ctx.handlerFunc = ch
		// ctx.Next()

		r.g.pool.Put(ctx)
	})
}

func checkHandler(h interface{}) CHandler {

	switch h := h.(type) {
	case CHandler:
		return h
	case func(*Context) error:
		return h
	default:
		panic("handler type error")
	}
}

func (r *Router) set(method, path string, ch CHandler) {
	r.setRoute(method, path, ch)
}

// Get REST funcs
func (r *Router) Get(path string, ch CHandler) {
	r.set("GET", path, ch)
}

// Post REST funcs
func (r *Router) Post(path string, ch CHandler) {
	r.set("POST", path, ch)
}

// Delete REST funcs
func (r *Router) Delete(path string, ch CHandler) {
	r.set("DELETE", path, ch)
}

// Head REST funcs
func (r *Router) Head(path string, ch CHandler) {
	r.set("HEAD", path, ch)
}

// Options REST funcs
func (r *Router) Options(path string, ch CHandler) {
	r.set("OPTIONS", path, ch)
}

// Put REST funcs
func (r *Router) Put(path string, ch CHandler) {
	r.set("PUT", path, ch)
}

// Patch REST funcs
func (r *Router) Patch(path string, ch CHandler) {
	r.set("PATCH", path, ch)
}

func (r *Router) StaticPath(dir string) {
	fileServer := http.FileServer(http.Dir(dir))

	r.hr.GET("/"+dir+"/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		r.URL.Path = p.ByName("filepath")
		fileServer.ServeHTTP(w, r)
	})
}

// REST for set all REST route
func (r *Router) REST(path string, c ControllerInterface) {
	// get all functions in controller
	refT := reflect.TypeOf(c)
	for i := 0; i < refT.NumMethod(); i++ {
		m := refT.Method(i)
		if checkSupportProto(m.Name) {
			revf := reflect.ValueOf(c)
			r.set(strings.ToUpper(m.Name), path, revf.MethodByName(m.Name).Interface().(func(*Context) error))
		}

	}
}

func checkSupportProto(proto string) bool {
	for _, v := range supportRestProto {
		if v == strings.ToUpper(proto) {
			return true
		}
	}

	return false
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.hr.ServeHTTP(w, req)
}
