package routers

import (
	"github.com/gowebtw/gas"
	"github.com/gowebtw/gas/Example/controllers"
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
