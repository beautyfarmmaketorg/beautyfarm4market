package util

import (
	"net/http"
	"html/template"
	"io/ioutil"
	"path"
	"log"
)

const(
	VIEW_DRI = "/beautyfarm4market/html/"
)

var tempaltes =make(map[string]*template.Template)

func init()  {
	absoluteViewsDir:=GetCurrentPath()+VIEW_DRI
	fileInfoArr,err:=ioutil.ReadDir(absoluteViewsDir)
	if err!=nil {
		panic(err)
		return
	}
	var templateName,templatePath string
	for _,fileInfo:=range fileInfoArr{
		templateName = fileInfo.Name()
		if ext:=path.Ext(templateName);ext!=".html" {
			continue
		}
		templatePath=absoluteViewsDir+templateName
		log.Printf("loading path",templatePath)
		t:=template.Must(template.ParseFiles(templatePath))
		tempaltes[templateName] = t
	}
}


func RenderHtml(w http.ResponseWriter,tmpl string,locals map[string]interface{}) (err error) {
	err =tempaltes[tmpl].Execute(w,locals)
	return
}
