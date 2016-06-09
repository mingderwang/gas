package main

import (
	"github.com/gowebtw/goslim"
	"github.com/gowebtw/goslim/Example/routers"
	"github.com/gowebtw/goslim/middleware"
)

func main() {
	g := goslim.New()

	routers.RegistRout(g.Router)

	g.Router.Use(middleware.LogMiddleware)

	g.Run()
}
