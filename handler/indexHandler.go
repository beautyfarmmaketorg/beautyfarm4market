package handler

import (
	"net/http"
	"beautyfarm4market/util"
)

func IndexHandler(w http.ResponseWriter,r *http.Request)  {
	IsNewUser("18221647820","孙龙飞")
	if r.Method=="GET" {
		util.RenderHtml(w,"index.html",nil)
		return
	}
	return
}
