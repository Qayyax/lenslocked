package controllers

import (
	"net/http"

	"github.com/qayyax/lenslock/views"
)

func StaticController(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}
