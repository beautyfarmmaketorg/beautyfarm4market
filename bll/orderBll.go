package bll

import (
	"time"
	"strconv"
	"beautyfarm4market/dal"
	"beautyfarm4market/entity"
	"beautyfarm4market/config"
	"sort"
	"fmt"
	"strings"
	"beautyfarm4market/util"
	"encoding/xml"
	"net/http"
	"bytes"
	"io/ioutil"
)

//获取订单号待实现
func GetOrderNo(prefix string) string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return prefix + timestamp
}

func CancelOrder(orderNo string) entity.BaseResultEntity {
	res := entity.GetBaseFailRes()
	return res
}

func Refund(mappingOrderNo string, remark string) entity.BaseResultEntity {
	res := entity.GetBaseFailRes()
	tempOrderInfo := dal.GetOrdersByMappingOrderNo(mappingOrderNo)
	if tempOrderInfo.MappingOrderNo == "" {
		res.Message = "订单号不存在"
		return res
	}
	if tempOrderInfo.PayStatus != 2 {
		res.Message = "订单状态不正确"
		return res
	}
	c := config.ConfigInfo
	weChatRefundReq := WeChatRefundReq{
		Appid:          c.WeChatAppId,
		Mch_id:         c.WeChatMchId,
		Nonce_str:      strconv.FormatInt(time.Now().Unix(), 10),
		Out_refund_no:  GetOrderNo("Refund"),
		Refund_fee:     int(tempOrderInfo.TotalPrice * 100),
		Total_fee:      int(tempOrderInfo.TotalPrice * 100),
		Transaction_id: tempOrderInfo.WechatorderNo,
	}
	weChatRefundReq.Sign = getSign4WeChatRefund(weChatRefundReq)

	xmlStr, _ := xml.Marshal(weChatRefundReq);
	fmt.Printf(string(xmlStr))
	dal.AddLog(dal.LogInfo{Title: "RefundxmlReq" + mappingOrderNo, Description: string(xmlStr), Type: 1})
	if response, postErr := http.Post(config.ConfigInfo.WeRefundUrl, "text/plain", bytes.NewBuffer([]byte(xmlStr)));postErr==nil{
		defer response.Body.Close()
		check(postErr)
		weChatRefundRes := WeChatRefundRes{}
		if response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(response.Body)
			payResponseXml := string(body)
			weChatRefundRes = convenrtWeChatRefundRes(payResponseXml)
			if weChatRefundRes.Return_code == "SUCCESS" {
				res.IsSucess = true
				res.Message = "退款成功"
				dal.UpdateTempOrderPayStatus(mappingOrderNo, 3)
			} else {
				res.Code = weChatRefundRes.Return_code
				res.Message = weChatRefundRes.Return_msg
			}
			dal.AddLog(dal.LogInfo{Title: "weChatRefundRes_" + mappingOrderNo, Description: payResponseXml, Type: 1})
		} else {
			res.Code = string(response.StatusCode)
			res.Message = "其他错误"
		}
	}
	return res
}

func convenrtWeChatRefundRes(xmlStr string) WeChatRefundRes {
	v := WeChatRefundRes{}
	err := xml.Unmarshal([]byte(xmlStr), &v)
	check(err)
	return v
}

type WeChatRefundReq struct {
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce_str      string `xml:"nonce_str"`
	Out_refund_no  string `xml:"out_refund_no"`
	Refund_fee     int `xml:"refund_fee"`
	Refund_desc    string `xml:"refund_desc"`
	Total_fee      int `xml:"total_fee"`
	Transaction_id string `xml:"transaction_id"`
	Sign           string `xml:"sign"`
}

type WeChatRefundRes struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mch_id"`
	Nonce_str      string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Transaction_id string `xml:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no"`
	Out_refund_no  string `xml:"out_refund_no"`
	Refund_id      string `xml:"refund_id"`
	Refund_channel string `xml:"refund_channel"`
	Refund_fee     int    `xml:"refund_fee"`
}

func check(e error) {
	if e != nil {
		//panic(e)
	}
}

func getSign4WeChatRefund(e WeChatRefundReq) string {
	m := make(map[string]interface{}, 0)
	m["appid"] = e.Appid
	m["mch_id"] = e.Mch_id
	m["nonce_str"] = e.Nonce_str
	m["out_refund_no"] = e.Out_refund_no
	m["refund_fee"] = e.Refund_fee
	m["total_fee"] = e.Total_fee
	m["transaction_id"] = e.Transaction_id
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, m[k])
		value := fmt.Sprintf("%v", m[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	signStrings = signStrings + "key=" + config.ConfigInfo.WeChatKey
	return strings.ToUpper(util.GetMd5(signStrings))
}
