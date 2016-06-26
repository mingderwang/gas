package gas

import (
	"encoding/json"
	"errors"
	"github.com/go-gas/gas/model"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"html/template"
	"net/http"
)

type Context struct {
	context.Context
	RespWriter *ResponseWriter
	Req        *http.Request
	ps         *httprouter.Params
	// handlerFunc CHandler

	gas *gas //

	// DB
	isUseDB bool
	mobj    model.ModelInterface
}

// create context
func createContext(w *ResponseWriter, r *http.Request, g *gas) *Context {
	c := &Context{}
	c.RespWriter = w
	c.Req = r
	c.gas = g

	return c
}

// reset context when get it from buffer
func (ctx *Context) reset(w http.ResponseWriter, r *http.Request, g *gas) {
	ctx.Req = r
	ctx.RespWriter.reset(w)

	ctx.gas = g

	ctx.mobj = nil
	ctx.isUseDB = false
}

// func (ctx *Context) Next()  {
//     ctx.handlerFunc(ctx)
// }

// Get parameter from post or get value
func (ctx *Context) GetParam(name string) string {
	if ctx.Req.PostForm == nil || ctx.Req.Form == nil {
		ctx.Req.ParseForm()
	}

	if v := ctx.Req.FormValue(name); v != "" {
		return v
	}

	return ctx.ps.ByName(name)
}

// func (ctx *Context) GetAllParams()  {
//     res := make(map[string]string, 0)

//     for key, v := ctx.Req.Form {
//         res[key] = v[0]
//     }
// }

// Render function combined data and template to show
func (ctx *Context) Render(data interface{}, tplPath ...string) error {
	if len(tplPath) == 0 {
		return errors.New("File path can not be empty")
	}

	ctx.RespWriter.Header().Set(ContentType, TextHTMLCharsetUTF8)

	// tpls := strings.Join(tplPath, ", ")
	tmpl := template.New("gas")

	for _, tpath := range tplPath {
		tmpl = template.Must((tmpl.ParseFiles(tpath)))
	}

	err := tmpl.Execute(ctx.RespWriter, data)

	return err
	// if err != nil {
	//     // println(err)
	//     // panic(err)

	//     return err
	// }

	// return nil
}

// Set the response data-type to html
func (ctx *Context) HTML(code int, html string) error {

	ctx.RespWriter.Header().Set(ContentType, TextHTMLCharsetUTF8)
	ctx.RespWriter.WriteHeader(code)
	_, err := ctx.RespWriter.Write([]byte(html))

	return err
}

// Set the response data-type to plain text
func (ctx *Context) STRING(status int, data string) error {

	ctx.RespWriter.Header().Set(ContentType, TextPlainCharsetUTF8)
	ctx.RespWriter.WriteHeader(status)
	_, err := ctx.RespWriter.Write([]byte(data))

	return err
}

// Set response data-type to json and try to json encode
func (ctx *Context) JSON(status int, data interface{}) error {

	ctx.RespWriter.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	ctx.RespWriter.WriteHeader(status)

	// encode json string
	jsonstr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// fmt.Fprint(ctx.Writer, data)
	_, errr := ctx.RespWriter.Write(jsonstr)

	return errr
}

// Get model using context in controller
func (ctx *Context) GetModel() model.ModelInterface {
	m := ctx.gas.NewModel()

	if m != nil {
		ctx.isUseDB = true
		ctx.mobj = m

		return m
	}

	return nil
}

// Close db connection
func (ctx *Context) CloseDB() error {
	return ctx.mobj.Builder().GetDB().Close()
}
