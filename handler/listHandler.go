package handler

import (
	"net/http"
	"io/ioutil"
	"io"
)

func ListHandler(w http.ResponseWriter,r *http.Request) {
	fileInfoArr,err:=ioutil.ReadDir(UPLOAD_DRI)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	var listHtml string
	for _,fileInfo:=range fileInfoArr{
		imgid :=fileInfo.Name()
		listHtml+="<li><a href=\"/view?id="+imgid+"\">"+imgid+"</a></li>"
	}
	io.WriteString(w,"<html><body><ol>"+listHtml+"</ol></body></html>")
}