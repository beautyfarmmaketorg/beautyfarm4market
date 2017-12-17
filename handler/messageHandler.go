package handler

import (
	"net/http"
	"beautyfarm4market/proxy"
	"beautyfarm4market/entity"
	"encoding/json"
	"fmt"
	"beautyfarm4market/config"
	"time"
	"strconv"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		mobileNo := r.FormValue("mobileNo")
		isSucess := false
		message := ""
		if proxy.SendMsg(mobileNo, fmt.Sprintf(config.ConfigInfo.SmsOfVaild, getCode())) {
			isSucess = true
			message = "短信发送成功，请查看手机获取"
		} else {
			message = "短信发送失败，请稍后重试"
		}
		sendMsgResult := entity.SendMsgResult{}
		sendMsgResult.IsSucess = isSucess
		sendMsgResult.Message = message
		sendMsgResult.Mobile = mobileNo

		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(sendMsgResult)
		return
	}
}

func getCode() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return timestamp[15:]
}
