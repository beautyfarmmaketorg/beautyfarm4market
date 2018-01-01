package handler

import (
	"net/http"
	"beautyfarm4market/util"
)

func BackyardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		util.RenderHtml(w, "backyard.html", nil)
		return
	}
}
