package handler

import (
	"net/http"
	"os"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DRI + imageId
	if exist := isExist(imagePath);!exist {
		http.NotFound(w,r)
		return
	}
	w.Header().Set("Content-Type","image")
	http.ServeFile(w,r,imagePath)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
