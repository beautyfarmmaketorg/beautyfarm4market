package handler

import (
	"net/http"
	"beautyfarm4market/dal"
	"time"
	"beautyfarm4market/config"
	"strconv"
	"sort"
	"fmt"
	"strings"
	"beautyfarm4market/util"
	"html/template"
)

//记录openId 到cookie
func HandlerWeChatLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		dal.AddLog(dal.LogInfo{Title: "HandlerWeChatLoginHandler", Description: r.RequestURI, Type: 1})
		openId := r.FormValue("openid")
		mappingOrderNo := r.FormValue("mappingOrderNo")
		mappingOrderNoCookie := http.Cookie{Name: "mappingOrderNoCookie",
			Value: mappingOrderNo, Path: "/", Expires: time.Now().Add(time.Hour * 1), MaxAge: 8600}
		http.SetCookie(w, &mappingOrderNoCookie) //记录订单号
		tempOrderInfo := dal.GetOrdersByMappingOrderNo(mappingOrderNo);
		if tempOrderInfo.MappingOrderNo=="" {
			http.Redirect(w,r,r.Host,http.StatusFound);
			return
		}
		if tempOrderInfo.PayStatus == 2 {
			//已支付则跳转首页
			locals := make(map[string]interface{})
			p:=dal.GetProductInfo(tempOrderInfo.ProductId)
			pageInfo := PageInfo{Channelcode: tempOrderInfo.Channel, ProductId: strconv.Itoa(int(tempOrderInfo.ProductId)), Bg: p.Backgroud_image, Button: p.PurhchaseBtn_image, Rule: p.Rule_image, Mask: p.MaskImage,RuleDesc:template.HTML(p.Prodcut_rule)}
			locals["pageInfo"] = pageInfo
			util.RenderHtml(w, "index.html", locals)
			return
		}
		weChatUnifiedorderResponse := InvokeWeChatUnifiedorder(tempOrderInfo.ProductCode, tempOrderInfo.ProductName,
			mappingOrderNo,
			tempOrderInfo.ClientIp, tempOrderInfo.TotalPrice, r.Host, "JSAPI", openId)
		if weChatUnifiedorderResponse.ReturnCode == "SUCCESS" && weChatUnifiedorderResponse.PrepayId != "" {
			dal.UpdateTempOrderPayStatus(mappingOrderNo, 1) //更新支付状态
			weChatLoginAddOrderParams := getWeChatLoginAddOrderParams(weChatUnifiedorderResponse.PrepayId, r.Host+"/?productId="+strconv.FormatInt(tempOrderInfo.ProductId,10)+"&channelcode="+tempOrderInfo.Channel)
			locals := make(map[string]interface{})
			locals["weChatLoginAddOrderParams"] = weChatLoginAddOrderParams
			dal.AddJsonLog("weChatPayLocals", weChatLoginAddOrderParams)
			util.RenderHtml(w, "weChatPay.html", locals)
			return
		} else {
			locals := make(map[string]interface{})
			p:=dal.GetProductInfo(tempOrderInfo.ProductId)
			pageInfo := PageInfo{Channelcode: tempOrderInfo.Channel, ProductId: strconv.Itoa(int(tempOrderInfo.ProductId)), Bg: p.Backgroud_image, Button: p.PurhchaseBtn_image, Rule: p.Rule_image, Mask: p.MaskImage,RuleDesc:template.HTML(p.Prodcut_rule)}
			locals["pageInfo"] = pageInfo
			util.RenderHtml(w, "index.html", locals)
			return
		}
	}
}

func getWeChatLoginAddOrderParams(prepayId string, indexUrl string) WeChatLoginAddOrderParams {
	args := WeChatLoginAddOrderParams{
		AppId:     config.ConfigInfo.WeChatAppId,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  strconv.FormatInt(time.Now().Unix(), 10),
		Package:   "prepay_id=" + prepayId,
		SignType:  "MD5",
		IndexUrl:  indexUrl,
	}
	sign := getSign4WeChatPay(args)
	args.PaySign = sign
	return args
}

func getSign4WeChatPay(e WeChatLoginAddOrderParams) string {
	m := make(map[string]interface{}, 0)
	m["appId"] = e.AppId
	m["timeStamp"] = e.TimeStamp
	m["nonceStr"] = e.NonceStr
	m["package"] = e.Package
	m["signType"] = e.SignType
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, m[k])
		value := fmt.Sprintf("%v", m[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	signStrings = signStrings + "key=" + config.ConfigInfo.WeChatKey
	return strings.ToUpper(util.GetMd5(signStrings))
}

//微信公众号支付参数
type WeChatLoginAddOrderParams struct {
	AppId     string
	TimeStamp string
	NonceStr  string
	Package   string
	SignType  string
	PaySign   string
	IndexUrl  string
}
