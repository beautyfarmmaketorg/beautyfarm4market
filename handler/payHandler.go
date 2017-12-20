package handler

import (
	"io/ioutil"
	"fmt"
	"beautyfarm4market/util"
	"time"
	"strconv"
	"strings"
)

func GetPayUrl(productCode string, productName string,
	orderCode string, spbill_create_ip string, totalPrice int) {
	xmlFilePath := util.GetCurrentPath() + "/config/webchatpaytemplate.xml"
	template, err := ioutil.ReadFile(xmlFilePath)
	check(err)
	e := getwechatPayEntity(productCode, productName, orderCode, spbill_create_ip, totalPrice)
	sign := getSign(e)
	xmStr := fmt.Sprintf(string(template), e.appid, e.attach, e.body, e.mch_id, e.nonce_str, e.notify_url,
		e.openid, e.out_trade_no, e.spbill_create_ip, e.total_fee, e.trade_type, sign)
	fmt.Printf(xmStr)
}

func getSign(e wechatPayEntity) string {
	//template := "appid=%s&attach=%s&body=%s&mch_id=%s&nonce_str=%s&notify_url=%s&openid=%s&out_trade_no=%s&spbill_create_ip=%s&total_fee=%s&trade_type=%s&key=%s"
	//str := fmt.Sprintf(template, e.appid, e.attach, e.body, e.mch_id, e.nonce_str, e.notify_url, e.openid, e.out_trade_no, e.spbill_create_ip, e.total_fee, e.trade_type, "meilitianyuan2016isgood2016igood")
	var str2 string
	str2 += "appid=" + e.appid
	str2 += "&attach=" + e.attach
	str2 += "&body=" + e.body
	str2 += "&mch_id=" + e.mch_id
	str2 += "&nonce_str=" + e.nonce_str
	var str3 string= "&notify_url=" + e.notify_url
	str3 += "&openid=" + e.openid
	str3 += "&out_trade_no=" + e.out_trade_no
	str3 += "&spbill_create_ip=" + e.spbill_create_ip
	str3 += "&total_fee=" + string(e.total_fee)
	str3 += "&trade_type=" + e.trade_type
	str3 += "&key=meilitianyuan2016isgood2016igood"

	str4:="纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理纯新胶原精华护理"
	return strings.ToUpper(util.GetMd5(str4))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getwechatPayEntity(productCode string, productName string,
	orderCode string, spbill_create_ip string, totalPrice int) wechatPayEntity {
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
		openid:           "oUpF8uMuAJO_M2pxb1Q9zNjWeS6o",
		out_trade_no:     orderCode,
		spbill_create_ip: spbill_create_ip,
		total_fee:        totalPrice,
		trade_type:       "NATIVE",
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
	total_fee        int
	trade_type       string
	sign             string
}
