package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"beautyfarm4market/entity"
	"encoding/json"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mobileNo := r.FormValue("mobileNo")
		isSucess := false
		message := ""
		if util.SendMsg(mobileNo,"123") {
			isSucess = true
			message ="短信发送成功，请查看手机获取"
		}else {
			message ="短信发送失败，请稍后重试"
		}
		sendMsgResult := entity.SendMsgResult{}
		sendMsgResult.IsSucess=isSucess
		sendMsgResult.Message = message
		sendMsgResult.Mobile=mobileNo

		w.Header().Set("Content-Type","application/json;charset=utf-8")
		json.NewEncoder(w).Encode(sendMsgResult)
		return
	}
}
