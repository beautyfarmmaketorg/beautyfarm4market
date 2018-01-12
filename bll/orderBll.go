package bll

import (
	"time"
	"strconv"
	"beautyfarm4market/dal"
	"beautyfarm4market/entity"
	"beautyfarm4market/proxy"
	"beautyfarm4market/config"
	"sort"
	"fmt"
	"strings"
	"beautyfarm4market/util"
)

//获取订单号待实现
func GetOrderNo(prefix string) string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return prefix + timestamp
}

func CancelOrder(orderNo string)entity.BaseResultEntity {
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
	weChatRefundReq := proxy.WeChatRefundReq{
		Appid:          c.WeChatAppId,
		Mch_id:         c.WeChatMchId,
		Nonce_str:      strconv.FormatInt(time.Now().Unix(), 10),
		Out_refund_no:  GetOrderNo("Refund"),
		Refund_fee:     int(tempOrderInfo.TotalPrice * 100),
		Total_fee:      int(tempOrderInfo.TotalPrice * 100),
		Transaction_id: tempOrderInfo.WechatorderNo,
	}
	weChatRefundReq.Sign = getSign4WeChatRefund(weChatRefundReq)
	if refundRes, serverRes := proxy.WeChatRefund(weChatRefundReq); serverRes.IsSucess {
		if refundRes.Return_code == "SUCCESS" {
			res.IsSucess = true
			res.Message = "退款成功"
		} else {
			res.Code = refundRes.Return_code
			res.Message = refundRes.Return_msg
		}
	}
	return res
}

func getSign4WeChatRefund(e proxy.WeChatRefundReq) string {
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
