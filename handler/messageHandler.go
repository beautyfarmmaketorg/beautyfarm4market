package handler

import (
	"net/http"
	"log"
	"io"
)

func MessageHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method=="GET" {
		mobileNo:=r.FormValue("mobileNo")
		log.Printf("Mobile",mobileNo)
		io.WriteString(w,"{\"isOk\":true,\"message\":\"发送成功\",\"mobile\":\""+mobileNo+"\"}")
		return
	}
}
