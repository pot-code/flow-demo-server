package hello

import (
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"

	"github.com/labstack/echo/v4"
)

func hello(c echo.Context) error {
	return api.JsonData(c, "Hello World!")
}

func post(c echo.Context) error {
	data := new(PostHelloDto)
	if err := api.Bind(c, data); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}
	return nil
}
