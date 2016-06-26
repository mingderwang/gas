package main

import (
	"github.com/go-gas/gas"
	"github.com/go-gas/gas/Example/routers"
	"github.com/go-gas/gas/middleware"
)

func main() {
	g := gas.New()

	routers.RegistRout(g.Router)

	g.Router.Use(middleware.LogMiddleware)

	g.Run()
}
