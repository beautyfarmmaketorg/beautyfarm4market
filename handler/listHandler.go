package handler

import (
	"net/http"
	"io/ioutil"
	"beautyfarm4market/util"
)

func ListHandler(w http.ResponseWriter,r *http.Request) {
	fileInfoArr,err:=ioutil.ReadDir(UPLOAD_DRI)
	if err!=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	locals:=make(map[string]interface{})
	imageArr:=[]string{}
	for _,fileInfo:=range fileInfoArr{
		imgid :=fileInfo.Name()
		imageArr = append(imageArr,imgid)
	}
	locals["images"] = imageArr
	if err:=util.RenderHtml(w,"list.html",locals);err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	return
}