package gas

import (
	//"net/http"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"reflect"
	"strings"
	//"github.com/julienschmidt/httprouter"
	//"github.com/julienschmidt/httprouter"
)

var supportRestProto = [7]string{"GET", "POST", "DELETE", "HEAD", "OPTIONS", "PUT", "PATCH"}

type (

	// Router class include httprouter and gas
	Router struct {
		*fasthttprouter.Router
		g           *gas
		middlewares []MiddlewareFunc
	}

	// MiddlewareFunc middlewarefunc define
	MiddlewareFunc func(GasHandler) GasHandler

	// CHandler is a function type for rout handler
	GasHandler func(*Context) error

	// PanicHandler defined panic handler
	PanicHandler func(*Context, interface{}) error
)

func newRouter(g *gas) *Router {
	fastR := fasthttprouter.New()
	r := &Router{}
	r.Router = fastR
	r.g = g

	return r
}

//func (r *Router) wrapGasHandlerToFasthttpRequestHandler(h GasHandler) fasthttp.RequestHandler {
//	// type RequestHandler func(ctx *RequestCtx)
//	return func(ctx *fasthttp.RequestCtx) {
//		gasCtx := r.g.pool.Get().(*Context)
//		gasCtx.reset(ctx, nil, r.g)
//
//		// chain middleware functions
//		var cpch GasHandler // copy handle avoid repeat chain
//		cpch = h
//
//		for i := len(r.middlewares) - 1; i >= 0; i-- {
//			cpch = r.middlewares[i](cpch)
//		}
//
//		if err := cpch(ctx); err != nil {
//			// handle error
//		}
//
//		if gasCtx.isUseDB {
//			defer gasCtx.CloseDB()
//		}
//
//		// ctx.handlerFunc = ch
//		// ctx.Next()
//
//		r.g.pool.Put(gasCtx)
//	}
//}

func (r *Router) wrapGasHandlerToFasthttpRouterHandler(h GasHandler) fasthttprouter.Handle {
	// type Handle func(*fasthttp.RequestCtx, Params)
	return func(ctx *fasthttp.RequestCtx, ps fasthttprouter.Params) {
		gasCtx := r.g.pool.Get().(*Context)
		gasCtx.reset(ctx, &ps, r.g)

		// chain middleware functions
		var cpch GasHandler // copy handle avoid repeat chain
		cpch = h

		for i := len(r.middlewares) - 1; i >= 0; i-- {
			cpch = r.middlewares[i](cpch)
		}

		if err := cpch(gasCtx); err != nil {
			// handle error
		}

		if gasCtx.isUseDB {
			defer gasCtx.CloseDB()
		}

		// ctx.handlerFunc = ch
		// ctx.Next()

		r.g.pool.Put(gasCtx)
	}
}

// SetNotFoundHandler  set Notfound and Panic handler
func (r *Router) SetNotFoundHandler(h GasHandler) {
	r.NotFound = func(fctx *fasthttp.RequestCtx) {
		ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
		ctx.reset(fctx, nil, r.g)

		// chain middleware functions
		var cpch GasHandler // copy handle avoid repeat chain
		cpch = h

		for i := len(r.middlewares) - 1; i >= 0; i-- {
			cpch = r.middlewares[i](cpch)
		}

		if err := cpch(ctx); err != nil {

		}

		r.g.pool.Put(ctx)
	}
	//r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	//	ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
	//	ctx.reset(w, req, r.g)
	//
	//	// chain middleware functions
	//	var cpch GasHandler // copy handle avoid repeat chain
	//	cpch = h
	//
	//	for i := len(r.middlewares) - 1; i >= 0; i-- {
	//		cpch = r.middlewares[i](cpch)
	//	}
	//
	//	if err := cpch(ctx); err != nil {
	//
	//	}
	//
	//	r.g.pool.Put(ctx)
	//})
}

func (r *Router) SetPanicHandler(ph PanicHandler) {
	r.PanicHandler = func(fctx *fasthttp.RequestCtx, rcv interface{}) {
		ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
		ctx.reset(fctx, nil, r.g)

		if err := ph(ctx, rcv); err != nil {

		}

		r.g.pool.Put(ctx)
	}
	//r.hr.PanicHandler = func(w http.ResponseWriter, req *http.Request, rcv interface{}) {
	//	// c := a.createContext(w, req)
	//	// a.panicFunc(c, rcv)
	//	// a.pool.Put(c)
	//
	//	ctx := r.g.pool.Get().(*Context) //createContext(rw, req)
	//	ctx.reset(w, req, r.g)
	//
	//	if err := ph(ctx, rcv); err != nil {
	//
	//	}
	//
	//	r.g.pool.Put(ctx)
	//}
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
	case func(GasHandler) GasHandler:
		return m
	case GasHandler:
		return wrapHandlerFuncToMiddlewareFunc(m)
	case func(c *Context) error:
		return wrapHandlerFuncToMiddlewareFunc(m)

	default:
		panic("unknown middleware")
	}
}

func wrapHandlerFuncToMiddlewareFunc(m GasHandler) MiddlewareFunc {
	return func(h GasHandler) GasHandler {
		return func(c *Context) error {
			if err := m(c); err != nil {
				return err
			}

			return h(c)
		}
	}
}

func (r *Router) setRoute(method, path string, ch GasHandler) {
	//r.hr.Handle(method, path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	r.Handle(method, path, r.wrapGasHandlerToFasthttpRouterHandler(ch))
}

//func checkHandler(h interface{}) GasHandler {
//
//	switch h := h.(type) {
//	case GasHandler:
//		return h
//	case func(*Context) error:
//		return h
//	default:
//		panic("handler type error")
//	}
//}

func (r *Router) set(method, path string, ch GasHandler) {
	r.setRoute(method, path, ch)
}

// Get REST funcs
func (r *Router) Get(path string, ch GasHandler) {
	r.set("GET", path, ch)
}

// Post REST funcs
func (r *Router) Post(path string, ch GasHandler) {
	r.set("POST", path, ch)
}

// Delete REST funcs
func (r *Router) Delete(path string, ch GasHandler) {
	r.set("DELETE", path, ch)
}

// Head REST funcs
func (r *Router) Head(path string, ch GasHandler) {
	r.set("HEAD", path, ch)
}

// Options REST funcs
func (r *Router) Options(path string, ch GasHandler) {
	r.set("OPTIONS", path, ch)
}

// Put REST funcs
func (r *Router) Put(path string, ch GasHandler) {
	r.set("PUT", path, ch)
}

// Patch REST funcs
func (r *Router) Patch(path string, ch GasHandler) {
	r.set("PATCH", path, ch)
}

func (r *Router) StaticPath(dir string) {

	//fileServer := http.FileServer(http.Dir(dir))

	//r.hr.GET("/"+dir+"/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//	w.Header().Set("Vary", "Accept-Encoding")
	//	w.Header().Set("Cache-Control", "public, max-age=7776000")
	//	r.URL.Path = p.ByName("filepath")
	//	fileServer.ServeHTTP(w, r)
	//})
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

//func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	r.ServeHTTP(w, req)
//}
