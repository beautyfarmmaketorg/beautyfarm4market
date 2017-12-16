package handler

import (
	"net/http"
	"io"
	"os"
	"beautyfarm4market/util"
)

func UploadHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == "GET" {
		util.RenderHtml(w,"upload.html",nil)
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
