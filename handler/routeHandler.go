package handler

import (
	"net/http"
)

func RouteHandler(w http.ResponseWriter, r *http.Request) {
	channelcode := r.FormValue("channelcode")
	http.Redirect(w, r, "http://promotion.beautyfarm.com.cn:8009/?channelcode="+
		channelcode, http.StatusFound)
}
