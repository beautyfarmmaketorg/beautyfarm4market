package main

import (
	"beautyfarm4market/handler"
	"net/http"
	"log"
)



func main()  {
	http.HandleFunc("/",handler.IndexHandler)
	http.HandleFunc("/upload",handler.UploadHandler)
	http.HandleFunc("/view",handler.ViewHandler)
	http.HandleFunc("/list",handler.ListHandler)
	err:=http.ListenAndServe(":8099",nil)
	if err!=nil {
		log.Fatal("ListenAndServe: ",err.Error())
	}
}

