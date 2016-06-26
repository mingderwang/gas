package controllers

import (
	"github.com/go-gas/gas"
	// "github.com/go-gas/gas/model"
	"github.com/go-gas/gas/Example/models"
	// "time"
)

func IndexPage(ctx *gas.Context) error {
	return ctx.Render("", "views/layout.html", "views/index.html")
}

func DefaultHi(ctx *gas.Context) error {
	a := map[string]string{
		"Name": ctx.GetParam("name"),
	}
	return ctx.Render(a, "views/layout2.html")

	// ctx.gas.Logger.Danger( time.Now().String() + " - YO~~" + ctx.GetParam("name") + "\n\n\n\n\n\n" )

	// ctx.JSON(200, ctx.GetParam("name"))

	// println("default hi")

	// return nil
}

func TestModel(ctx *gas.Context) error {
	m := ctx.GetModel()

	// u := &models.User{}
	// m.Save(u)
	// m.Builder().Select().From().Where()
	b, err := m.Builder().Where("id = ?", 1).Get(&models.TestUser{})
	if err != nil {
		println(err.Error())
	}
	for i := 0; i < len(b); i++ {
		for k, v := range b[i] {
			println(k, ": ", v)
		}
	}

	return ctx.STRING(200, "You can see testuser table schema and data print on console. If not please turn on your mysql server and modify the config file.")
}

func PostTest(ctx *gas.Context) error {
	// println("default hi")

	println(ctx.GetParam("Test"))

	// println(ctx.GetParam("name"))
	// println(ctx)

	// new model
	// u := &models.User{}
	// u.Builder = &model.MySQLBuilder{}
	// u.Builder.Select("ID", "Name").Get()
	// u.select("ID", "Name").get()
	// u.select("*").UserDetail().get()

	// a := map[string]string {
	//     "Name": ctx.GetParam("kkk"),
	// }
	// ctx.Render(a, "views/layout.html", "views/index.html")

	ctx.STRING(200, ctx.GetParam("kkk")+ctx.GetParam("Test"))

	return nil

}
