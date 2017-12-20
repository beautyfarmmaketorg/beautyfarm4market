package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"fmt"
)

func IndexHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method=="GET" {
		channelcode := r.FormValue("channelcode")
		fmt.Println(channelcode)
		util.RenderHtml(w,"index.html",nil)
		return
	}
	return
}
