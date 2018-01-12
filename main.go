package main

import (
	"beautyfarm4market/handler"
	"net/http"
	"log"
	_ "beautyfarm4market/config"
	"beautyfarm4market/config"
	"beautyfarm4market/proxy"
	"beautyfarm4market/bll"
	"fmt"
)

func main() {
	cancelRes, serviceRes := proxy.CancelSoaOrder("BZ1515723284147993500")
	if serviceRes.IsSucess && cancelRes.ErrorCode == "200" {
		fmt.Printf("取消院余成功")
	} else {
		fmt.Printf(cancelRes.ErrorMessage)
	}
	refundRes:=bll.Refund("BZ1515723284147993500","测试取消")
	fmt.Printf(refundRes.Message)
	fmt.Printf(cancelRes.OrderNo, serviceRes.Message)
	cfg := config.GetConfigFromXml()
	mux := http.NewServeMux()
	handler.StaticDirHandler(mux, "/assets/", 0)
	mux.HandleFunc("/", handler.SafeHandler(handler.IndexHandler))
	mux.HandleFunc("/upload", handler.SafeHandler(handler.UploadHandler))
	mux.HandleFunc("/view", handler.SafeHandler(handler.ViewHandler))
	mux.HandleFunc("/list", handler.SafeHandler(handler.ListHandler))
	mux.HandleFunc("/sendMsg", handler.SafeHandler(handler.MessageHandler))
	mux.HandleFunc("/addOrder", handler.SafeHandler(handler.AddOrderHandler))
	mux.HandleFunc("/orderList", handler.SafeHandler(handler.OrderListHandler))
	mux.HandleFunc("/promotion", handler.SafeHandler(handler.RouteHandler))
	mux.HandleFunc("/payCallBack", handler.SafeHandler(handler.PayCallBackHandler))
	mux.HandleFunc("/purchaseRes", handler.SafeHandler(handler.PurchaseResHandler))
	mux.HandleFunc("/checkPurchaseRes", handler.SafeHandler(handler.PurchaseResHandler))
	mux.HandleFunc("/handlerWeChatLogin", handler.SafeHandler(handler.HandlerWeChatLoginHandler))
	mux.HandleFunc("/favicon.ico", handler.SafeHandler(handler.DefaultHandler))
	mux.HandleFunc("/prodcut", handler.SafeHandler(handler.ProductHandler))
	mux.HandleFunc("/prodcutdetail", handler.SafeHandler(handler.ProductDetailHandler))
	mux.HandleFunc("/report", handler.SafeHandler(handler.ReportHandler))
	mux.HandleFunc("/backyard", handler.SafeHandler(handler.BackyardHandler))
	err := http.ListenAndServe(cfg.Port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
