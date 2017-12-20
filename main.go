package main

import (
	"beautyfarm4market/handler"
	"net/http"
	"log"
	_ "beautyfarm4market/config"
	"beautyfarm4market/config"
)

func main() {
	handler.GetPayUrl(config.ConfigInfo.ProductCode, config.ConfigInfo.ProductName, "123456",
		"14.23.150.211", 1)

	mux := http.NewServeMux()
	handler.StaticDirHandler(mux, "/assets/", 0)
	mux.HandleFunc("/", handler.SafeHandler(handler.IndexHandler))
	mux.HandleFunc("/upload", handler.SafeHandler(handler.UploadHandler))
	mux.HandleFunc("/view", handler.SafeHandler(handler.ViewHandler))
	mux.HandleFunc("/list", handler.SafeHandler(handler.ListHandler))
	mux.HandleFunc("/sendMsg", handler.SafeHandler(handler.MessageHandler))
	mux.HandleFunc("/addOrder", handler.SafeHandler(handler.AddOrderHandler))
	mux.HandleFunc("/orderList", handler.SafeHandler(handler.OrderListHandler))
	err := http.ListenAndServe(":8009", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
