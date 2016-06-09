package controllers

import (
	"github.com/gowebtw/goslim"
)

type RestController struct {
	goslim.ControllerInterface
}

func (rc *RestController) Get(c *goslim.Context) error {

	return c.STRING(200, "Test Get")
}
