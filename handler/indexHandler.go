package handler

import (
	"net/http"
	"beautyfarm4market/util"
)

func IndexHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method=="GET" {
		util.RenderHtml(w,"index.html",nil)
		return
	}
	return
}
