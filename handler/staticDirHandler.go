package handler

import (
	"net/http"
	"beautyfarm4market/util"
)

const (
	ListDir = 0x0001
    VIEW_DRI = "/html/"
)

func StaticDirHandler(mux *http.ServeMux,prefix string,flags int )  {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file:=util.GetCurrentPath()+VIEW_DRI+r.URL.Path[len(prefix)-1:]
		if (flags&ListDir==0) {
			if exists:=isExist(file);!exists {
				http.NotFound(w,r)
				return
			}
		}
		http.ServeFile(w,r,file)
	})
}