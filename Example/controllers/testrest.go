package controllers

import (
	"github.com/go-gas/gas"
)

type RestController struct {
	gas.ControllerInterface
}

func (rc *RestController) Get(c *gas.Context) error {

	return c.STRING(200, "Test Get")
}
