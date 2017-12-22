package handler

import (
	"io/ioutil"
	"beautyfarm4market/util"
	"time"
	"strconv"
	"strings"
	"beautyfarm4market/config"
	"net/http"
	"bytes"
	"beautyfarm4market/dal"
	"io"
	"encoding/xml"
	"sort"
	"fmt"
)

func PayCallBackHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	weChatNotifyResponse := WeChatNotifyResponse{}
	xml.Unmarshal(body, &weChatNotifyResponse)
	mappingOrderNo := ""
	payCallBackRes := PayCallBackRes{ReturnCode: "FAIL"}
	if weChatNotifyResponse.Result_code == "SUCCESS" {
		mappingOrderNo = weChatNotifyResponse.Out_trade_no
		wechatOrderNo := weChatNotifyResponse.Transaction_id
		timeEnd := weChatNotifyResponse.Time_end
		processTempOrderRes := processTempOrder(mappingOrderNo, wechatOrderNo, timeEnd)
		if processTempOrderRes {
			payCallBackRes.ReturnCode = "SUCCESS"
			payCallBackRes.ReturnMsg = "OK"
		}
	}
	dal.AddLog(dal.LogInfo{Title: "payCallBackRes" + mappingOrderNo, Description: string(body), Type: 1})
	fmt.Println(string(body))
	resXml, _ := xml.Marshal(payCallBackRes)
	io.WriteString(w, string(resXml))
}

type PayCallBackRes struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

//添加正式单 发送短信 更新临时单
func processTempOrder(mappingOrderNo string, wechatOrderNo string, timeEnd string) bool {
	res := false
	tempOrder := dal.GetOrdersByMappingOrderNo(mappingOrderNo)
	if tempOrder.MappingOrderNo != "" {
		addRes := addFinalOrder(tempOrder.UserName, tempOrder.MobileNo, tempOrder.ProductCode, tempOrder.AccountNo, tempOrder.MappingOrderNo)
		if addRes.IsSucess {
			cardNo := addRes.CardNo
			res = dal.UpdateTempOrder(cardNo,addRes.OrderNo, mappingOrderNo, wechatOrderNo, timeEnd)
		}
	}
	return res
}

//调用微信统一下单接口
func InvokeWeChatUnifiedorder(productCode string, productName string,
	orderCode string, spbill_create_ip string, totalPrice int, host string) WeChatUnifiedorderResponse {
	payResponse := WeChatUnifiedorderResponse{}
	e := getwechatPayEntity(productCode, productName, orderCode, spbill_create_ip, totalPrice, host)
	sign := getSign(e)
	e.Sign = sign
	xmlStr, _ := xml.Marshal(e);
	dal.AddLog(dal.LogInfo{Title: "xmlReq_" + orderCode, Description: string(xmlStr), Type: 1})
	response, postErr := http.Post(config.ConfigInfo.UnifiedorderUrl, "text/plain", bytes.NewBuffer([]byte(xmlStr)))
	defer response.Body.Close()
	check(postErr)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		payResponseXml := string(body)
		payResponse = convenrtWeChatResponse(payResponseXml)
		dal.AddLog(dal.LogInfo{Title: "payResponse_" + orderCode, Description: payResponseXml, Type: 1})
	}
	return payResponse
}

func getSign(e wechatPayEntity) string {
	m := make(map[string]interface{}, 0)
	m["appid"] = e.Appid
	m["attach"] = e.Attach
	m["body"] = e.Body
	m["mch_id"] = e.Mch_id
	m["nonce_str"] = e.Nonce_str
	m["notify_url"] = e.Notify_url
	m["openid"] = e.Openid
	m["out_trade_no"] = e.Out_trade_no
	m["spbill_create_ip"] = e.Spbill_create_ip
	m["total_fee"] = e.Total_fee
	m["trade_type"] = e.Trade_type
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getwechatPayEntity(productCode string, productName string,
	orderCode string, spbill_create_ip string, totalPrice int, host string) wechatPayEntity {
	t := time.Now()
	nonce_str := strconv.FormatInt(t.UTC().UnixNano(), 10)
	notify_url := host + "/payCallBack" //异步回调
	entity := wechatPayEntity{
		Appid:            "wx7302aaa9857c055b",
		Attach:           productCode,
		Body:             productName,
		Mch_id:           "1301086301",
		Nonce_str:        nonce_str,
		Notify_url:       notify_url,
		Openid:           "",
		Out_trade_no:     orderCode,
		Spbill_create_ip: spbill_create_ip,
		Total_fee:        totalPrice,
		Trade_type:       "MWEB",
	}
	return entity
}

type wechatPayEntity struct {
	Appid            string `xml:"appid"`
	Attach           string `xml:"attach"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Openid           string `xml:"openid"`
	Out_trade_no     string `xml:"out_trade_no"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Trade_type       string `xml:"trade_type"`
	Sign             string `xml:"sign"`
}

//微信统一下单接口响应
type WeChatUnifiedorderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
	MwebUrl    string `xml:"mweb_url"`
}

func convenrtWeChatResponse(xmlStr string) WeChatUnifiedorderResponse {
	v := WeChatUnifiedorderResponse{}
	err := xml.Unmarshal([]byte(xmlStr), &v)
	check(err)
	return v
}

//https://pay.weixin.qq.com/wiki/doc/api/H5.php?chapter=9_7&index=7
type WeChatNotifyResponse struct {
	Appid          string `xml:"appid"`
	Attach         string `xml:"attach"`
	Bank_type      string `xml:"bank_type"`
	Fee_type       string `xml:"fee_type"`
	Is_subscribe   string `xml:"is_subscribe"`
	Mch_id         string `xml:"mch_id"`
	Nonce_str      string `xml:"nonce_str"`
	Openid         string `xml:"openid"`
	Out_trade_no   string `xml:"out_trade_no"`
	Result_code    string `xml:"result_code"`
	Sign           string `xml:"sign"`
	Sub_mch_id     string `xml:"sub_mch_id"`
	Time_end       string `xml:"time_end"`
	Total_fee      string `xml:"total_fee"`
	Coupon_fee     string `xml:"coupon_fee"`
	Coupon_count   string `xml:"coupon_count"`
	Coupon_type    string `xml:"coupon_type"`
	Coupon_id      string `xml:"coupon_id"`
	Trade_type     string `xml:"trade_type"`
	Transaction_id string `xml:"transaction_id"`
}
