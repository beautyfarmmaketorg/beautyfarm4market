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

		if enable, _ := checkCookieTime(mobileNo, w, r); !enable {
			isSucess = false
			message =""// fmt.Sprintf("请稍后，剩余时间:%d 秒", leftSeconds)
		} else {
			if proxy.SendMsg(mobileNo, fmt.Sprintf(config.ConfigInfo.SmsOfVaild, getCode())) {
				isSucess = true
				message = "短信发送成功，请查看手机获取"
				setMobileCookie(w, mobileNo, 8600)
			} else {
				message = "短信发送失败，请稍后重试"
			}
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

func checkCookieTime(mobileNo string, w http.ResponseWriter, r *http.Request) (enable bool, leftSecond int) {
	cookie, err := r.Cookie(fmt.Sprintf(config.ConfigInfo.MobileCookie, mobileNo))
	if err == nil {
		leftSecond = getLeftSecond(cookie.Value)
		if leftSecond <= 0 {
			setMobileCookie(w, mobileNo, -1)
			enable = true
		} else {
			enable = false
		}
		return
	}
	enable = true
	return
}

//返回相差的秒
func getLeftSecond(cookTimeStr string) int {
	var leftSeconds float64 = -1
	loc, _ := time.LoadLocation("Local")
	cookTime, _ := time.ParseInLocation(config.ConfigInfo.TimeLayout, cookTimeStr, loc)
	leftSeconds = time.Now().Sub(cookTime).Seconds()
	return 30 - int(leftSeconds)
}

func setMobileCookie(w http.ResponseWriter, mobileNo string, maxAge int) {
	dataTimeStr := time.Unix(time.Now().Unix(), 0).Format(config.ConfigInfo.TimeLayout)
	cookie := http.Cookie{Name: fmt.Sprintf(config.ConfigInfo.MobileCookie, mobileNo),
		Value: dataTimeStr, Path: "/", Expires: time.Now().Add(time.Second * 30), MaxAge: maxAge}
	http.SetCookie(w, &cookie)
}

func getCode() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return timestamp[15:]
}
