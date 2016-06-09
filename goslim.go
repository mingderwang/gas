package goslim

import (
	"fmt"
	"github.com/gowebtw/goslim/logger"
	"github.com/gowebtw/goslim/model"
	"net/http"
	"strings"
	"sync"
)

const (
	//-------------
	// Media types
	//-------------

	ApplicationJSON                  = "application/json"
	ApplicationJSONCharsetUTF8       = ApplicationJSON + "; " + CharsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + CharsetUTF8
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + CharsetUTF8
	ApplicationForm                  = "application/x-www-form-urlencoded"
	ApplicationProtobuf              = "application/protobuf"
	ApplicationMsgpack               = "application/msgpack"
	TextHTML                         = "text/html"
	TextHTMLCharsetUTF8              = TextHTML + "; " + CharsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + CharsetUTF8
	MultipartForm                    = "multipart/form-data"

	//---------
	// Charset
	//---------

	CharsetUTF8 = "charset=utf-8"

	//---------
	// Headers
	//---------

	AcceptEncoding     = "Accept-Encoding"
	Authorization      = "Authorization"
	ContentDisposition = "Content-Disposition"
	ContentEncoding    = "Content-Encoding"
	ContentLength      = "Content-Length"
	ContentType        = "Content-Type"
	Location           = "Location"
	Upgrade            = "Upgrade"
	Vary               = "Vary"
	WWWAuthenticate    = "WWW-Authenticate"
	XForwardedFor      = "X-Forwarded-For"
	XRealIP            = "X-Real-IP"
)

type (
	Goslim struct {
		Router *Router
		Config *config
		Model  *goslimModel
		pool   sync.Pool
		Logger *logger.Logger
	}

	goslimModel struct {
		model.Model
	}
)

// New goslim Object
func New() *Goslim {
	g := &Goslim{}

	// init logger
	g.Logger = logger.New("log/system.log")

	// init pool
	g.pool.New = func() interface{} {
		// c := &Context{}
		// c.Writer = &c.writercache
		c := createContext(new(ResponseWriter), nil, g)

		return c
	}

	// load config
	g.Config = &config{
		Mode:       "DEV",
		ListenAddr: "localhost",
		ListenPort: "8080",
		PubDir:     "public",
	}
	g.Config.loadConfig("config/default.yaml")

	// set router
	g.Router = &Router{g: g}

	// set default not found handler
	g.Router.SetNotFoundHandler(defaultNotFoundHandler)

	// set default panic handler
	g.Router.SetPanicHandler(defaultPanicHandler)

	// set static file path
	g.Router.StaticPath(g.Config.PubDir)
	// fileServer := http.FileServer(http.Dir(g.Config.PubDir))
	// g.Router.Get("/"+ g.Config.PubDir +"/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//     w.Header().Set("Vary", "Accept-Encoding")
	//     w.Header().Set("Cache-Control", "public, max-age=7776000")
	//     r.URL.Path = p.ByName("filepath")
	//     fileServer.ServeHTTP(w, r)
	// })

	// init model
	// g.Model = &goslimModel{}

	// if Config.Db.SQLDriver != "" && Config.Db.Username != "" && Config.Db.Dbname != "" {
	//     Model.Conn(Config.Db.Username, Config.Db.Password, Config.Db.Dbname)
	// }

	// add Log middleware
	// g.Router.Use(middleware.LogMiddleware)

	return g
}

func defaultNotFoundHandler(c *Context) error {
	return c.STRING(404, "Page Not Found.")
}

func defaultPanicHandler(c *Context, rcv interface{}) error {
	logStr := fmt.Sprintf("Panic occurred...rcv: %v", rcv)
	c.Goslim.Logger.Error(logStr)

	var output string
	if c.Goslim.Config.Mode == "DEV" {
		output = logStr
	} else {
		output = "Sorry...some error occurred..."
	}

	return c.STRING(500, output)
}

func (g *Goslim) LoadConfig(configPath string) {
	g.Config.loadConfig(configPath)
}

// Run framework
func (g *Goslim) Run() {
	// controller.Hi()
	// print("go slim go go go")

	fmt.Println("Server is Listen on: " + g.Config.ListenAddr + ":" + g.Config.ListenPort)
	if err := http.ListenAndServe(g.Config.ListenAddr+":"+g.Config.ListenPort, g.Router); err != nil {
		panic(err)
	}
}

func (g *Goslim) NewDb() model.SlimDbInterface {
	c := g.Config

	var d model.SlimDbInterface

	switch strings.ToLower(c.Db.SQLDriver) {
	case "mysql":
		d = new(model.MysqlDb)
	default:
		panic("Unknow Database Driver: " + g.Config.Db.SQLDriver)

	}

	d.ConnectWithConfig(g.Config.Db)

	return d

	// err := m.Connect(c.Db.Protocal, c.Db.Hostname, c.Db.Port, c.Db.Username, c.Db.Password, c.Db.Dbname, "charset=" + c.Db.Charset)
	// if err != nil {
	//     panic("Connection error: " + err.Error())
	// }

	// m.TestConn()

	// return m
}

func (g *Goslim) NewModel() model.ModelInterface {
	// get db
	// db := g.NewDb()
	c := g.Config

	var db model.SlimDbInterface
	var m model.ModelInterface
	var builder model.BuilderInterface

	switch strings.ToLower(c.Db.SQLDriver) {
	case "mysql":
		db = new(model.MysqlDb)
		m = new(model.MySQLModel)
		builder = new(model.MySQLBuilder)
	default:
		panic("Unknow Database Driver: " + g.Config.Db.SQLDriver)

	}

	err := db.ConnectWithConfig(g.Config.Db)
	if err != nil {
		panic(err.Error())
	}
	m.SetDB(db)
	builder.SetDB(db)
	m.SetBuilder(builder)

	return m
}
