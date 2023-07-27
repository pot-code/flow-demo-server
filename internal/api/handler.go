package api

import (
	"net/http"

	"github.com/pot-code/gobit/pkg/validate"
	"github.com/rs/zerolog/log"
)

type Handler func(r *http.Request, w http.ResponseWriter) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(r, w)
	if err != nil {
		switch e := err.(type) {
		case validate.ValidationError:
			JsonBusinessError(w, e[0].Error())
		case *DecoderError:
			JsonBadRequest(w, e.Error())
		default:
			log.Err(err).Msg("")
			JsonServerError(w, e.Error())
		}
	}
}
