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
	"beautyfarm4market/dal"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		sendMsgResult := entity.SendMsgResult{}
		mobileNo := r.FormValue("mobileNo")
		isSucess := false
		message := ""

		productIdStr:=r.FormValue("productId")
		if productIdStr == "" {
			productIdStr = "1"
		}
		productId, _ := strconv.ParseInt(productIdStr, 10, 64)
		productInfo:=dal.GetProductInfo(productId,true)
		productCode :=productInfo.Product_code
		//检查是否已经下过订单
		if hasOrdered := checkHasOrdered(mobileNo, productCode); hasOrdered {
			sendMsgResult.Code = "1" //
			json.NewEncoder(w).Encode(sendMsgResult)
			return
		}

		//vip用户
		if isVip := isVip(mobileNo); isVip {
			sendMsgResult.Code = "2" //
			json.NewEncoder(w).Encode(sendMsgResult)
			return
		}

		if enable, _ := checkCookieTime(mobileNo, w, r); !enable {
			isSucess = false
			message = "" // fmt.Sprintf("请稍后，剩余时间:%d 秒", leftSeconds)
		} else {
			code := getCode()
			if proxy.SendMsg(mobileNo, fmt.Sprintf(config.ConfigInfo.SmsOfVaild, code)) {
				isSucess = true
				message = "短信发送成功，请查看手机获取"
				setMobileCookie(w, mobileNo, 8600)
				setMobileCodeCookie(w, mobileNo, code, 8600)
			} else {
				message = "短信发送失败，请稍后重试"
			}
		}
		sendMsgResult.IsSucess = isSucess
		sendMsgResult.Message = message
		sendMsgResult.Mobile = mobileNo
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

func setMobileCodeCookie(w http.ResponseWriter, mobileNo string, code string, maxAge int) {
	cookie := http.Cookie{Name: fmt.Sprintf(config.ConfigInfo.CodeCookie, mobileNo),
		Value: code, Path: "/", Expires: time.Now().Add(time.Second * 30), MaxAge: maxAge}
	http.SetCookie(w, &cookie)
}

func getCode() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return timestamp[15:]
}
