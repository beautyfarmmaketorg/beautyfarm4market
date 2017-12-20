package handler

import (
	"io/ioutil"
	"fmt"
	"beautyfarm4market/util"
	"time"
	"strconv"
	"strings"
	"beautyfarm4market/config"
	"net/http"
	"bytes"
)

func GetPayUrl(productCode string, productName string,
	orderCode string, spbill_create_ip string, totalPrice string) string {
	payUrl := ""
	xmlFilePath := util.GetCurrentPath() + "/config/webchatpaytemplate.xml"
	template, err := ioutil.ReadFile(xmlFilePath)
	check(err)
	e := getwechatPayEntity(productCode, productName, orderCode, spbill_create_ip, totalPrice)
	sign := getSign(e)
	xmlStr := fmt.Sprintf(string(template), e.appid, e.attach, e.body, e.mch_id, e.nonce_str, e.notify_url,
		e.out_trade_no, e.spbill_create_ip, e.total_fee, e.trade_type, sign)

	response, postErr := http.Post(config.ConfigInfo.WeChatPayUrl, "text/plain", bytes.NewBuffer([]byte(xmlStr)))

	defer response.Body.Close()
	check(postErr)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		payUrl = string(body)
		//fmt.Printf(payUrl)
	}

	return payUrl
}

func getSign(e wechatPayEntity) string {
	template := "appid=%s&attach=%s&body=%s&mch_id=%s&nonce_str=%s&notify_url=%s&out_trade_no=%s&spbill_create_ip=%s&total_fee=%s&trade_type=%s&key=%s"
	str := fmt.Sprintf(template, e.appid, e.attach, e.body, e.mch_id, e.nonce_str, e.notify_url,
		e.out_trade_no, e.spbill_create_ip, e.total_fee, e.trade_type, config.ConfigInfo.WeChatKey)

	return strings.ToUpper(util.GetMd5(str))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getwechatPayEntity(productCode string, productName string,
	orderCode string, spbill_create_ip string, totalPrice string) wechatPayEntity {
	t := time.Now()
	nonce_str := strconv.FormatInt(t.UTC().UnixNano(), 10)
	notify_url := "http://www.baidu.com" //异步回调
	entity := wechatPayEntity{
		appid:            "wx7302aaa9857c055b",
		attach:           productCode,
		body:             productName,
		mch_id:           "1301086301",
		nonce_str:        nonce_str,
		notify_url:       notify_url,
		openid:           "",
		out_trade_no:     orderCode,
		spbill_create_ip: spbill_create_ip,
		total_fee:        totalPrice,
		trade_type:       "MWEB",
	}
	return entity
}

type wechatPayEntity struct {
	appid            string
	attach           string
	body             string
	mch_id           string
	nonce_str        string
	notify_url       string
	openid           string
	out_trade_no     string
	spbill_create_ip string
	total_fee        string
	trade_type       string
	sign             string
}
