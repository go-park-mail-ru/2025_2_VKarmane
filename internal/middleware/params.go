package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetParam(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}
