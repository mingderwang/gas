package routers

import (
	"github.com/go-gas/gas"
	"github.com/go-gas/gas/Example/controllers"
)

func RegistRout(r *gas.Router) {

	// dc := &controllers.DefaultController{}
	r.Get("/", controllers.IndexPage)
	r.Get("/modeltest", controllers.TestModel)
	r.Post("/post/:kkk", controllers.PostTest)
	r.Get("/user/:name", controllers.DefaultHi)

	rc := &controllers.RestController{}
	r.REST("/User", rc)

}
