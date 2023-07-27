package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Json(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(fmt.Errorf("error encoding json: %w", err))
	}
	return nil
}

func JsonData(w http.ResponseWriter, data any) error {
	return Json(w, http.StatusOK, map[string]any{
		"code": http.StatusOK,
		"data": data,
	})
}

func JsonBusinessError(w http.ResponseWriter, msg string) error {
	return Json(w, http.StatusOK, map[string]any{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}

func JsonBadRequest(w http.ResponseWriter, msg string) error {
	return Json(w, http.StatusBadRequest, map[string]any{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}

func JsonServerError(w http.ResponseWriter, msg string) error {
	return Json(w, http.StatusInternalServerError, map[string]any{
		"code": http.StatusInternalServerError,
		"msg":  msg,
	})
}

func JsonUnauthorized(w http.ResponseWriter, msg string) error {
	return Json(w, http.StatusForbidden, map[string]any{
		"code": http.StatusForbidden,
		"msg":  msg,
	})
}

func JsonUnauthenticated(w http.ResponseWriter, msg string) error {
	return Json(w, http.StatusUnauthorized, map[string]any{
		"code": http.StatusUnauthorized,
		"msg":  msg,
	})
}
