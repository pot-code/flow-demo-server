package api

import (
	"github.com/labstack/echo/v4"
	"github.com/pot-code/gobit/pkg/validate"
	"github.com/rs/zerolog/log"
)

func ErrorHandler(err error, c echo.Context) {
	switch e := err.(type) {
	case validate.ValidationError:
		JsonBusinessError(c, e[0].Error())
	case *validate.ValidationResult:
		JsonBadRequest(c, e.Error())
	case *BindError:
		JsonBadRequest(c, e.Error())
	default:
		log.Err(err).Msg("")
		JsonServerError(c, e.Error())
	}
}
