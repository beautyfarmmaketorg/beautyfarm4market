package handler

import (
	"net/http"
	"io"
	"os"
)

func UploadHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == "GET" {
		io.WriteString(w,"<html><head><title>上传照片</title></head><body> <form method=\"post\" action=\"/upload\""+
			" enctype =\"multipart/form-data\">"+
				" Chose an image to upload :<input name=\"image\" type=\"file\" />"+
					"<input type=\"submit\" value=\"upload\"/></form></body></html>")
		return
	}

	if r.Method=="POST" {
		f,h,err:=r.FormFile("image")
		if err!=nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		filename:=h.Filename
		defer f.Close()

		t,err :=os.Create(UPLOAD_DRI+filename)
		if err!=nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _,err:=io.Copy(t,f);err!=nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		http.Redirect(w,r,"/view?id="+filename,http.StatusFound)
	}
}
