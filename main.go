package main

import (
	"beautyfarm4market/handler"
	"net/http"
	"log"
)


func main()  {
	mux:=http.NewServeMux()
	handler.StaticDirHandler(mux,"/assets/",0)
	mux.HandleFunc("/",handler.SafeHandler(handler.IndexHandler))
	mux.HandleFunc("/upload",handler.SafeHandler(handler.UploadHandler))
	mux.HandleFunc("/view",handler.SafeHandler(handler.ViewHandler))
	mux.HandleFunc("/list",handler.SafeHandler(handler.ListHandler))
	mux.HandleFunc("/sendMsg",handler.SafeHandler(handler.MessageHandler))
	err:=http.ListenAndServe(":8099",mux)
	if err!=nil {
		log.Fatal("ListenAndServe: ",err.Error())
	}
}

