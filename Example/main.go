package main

import (
	"github.com/gowebtw/gas"
	"github.com/gowebtw/gas/Example/routers"
	"github.com/gowebtw/gas/middleware"
)

func main() {
	g := gas.New()

	routers.RegistRout(g.Router)

	g.Router.Use(middleware.LogMiddleware)

	g.Run()
}
