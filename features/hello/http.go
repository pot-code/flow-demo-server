package hello

import (
	"gobit-demo/internal/api"
	"gobit-demo/internal/validate"
	"net/http"
)

func hello(r *http.Request, w http.ResponseWriter) error {
	return api.JsonData(w, "Hello World!")
}

func post(r *http.Request, w http.ResponseWriter) error {
	data := new(PostHelloDto)
	if err := api.DecodeFromRequestBody(data, r.Body); err != nil {
		return err
	}
	if err := validate.Validator.Struct(data); err != nil {
		return err
	}
	return nil
}
