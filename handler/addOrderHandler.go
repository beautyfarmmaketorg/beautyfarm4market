package handler

import (
	"net/http"
	"beautyfarm4market/entity"
	"encoding/json"
	"beautyfarm4market/proxy"
	"fmt"
	"beautyfarm4market/config"
	"time"
	"strconv"
)

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	mobileNo := r.FormValue("mobileNo")
	code := r.FormValue("code")
	result := addFinalOrder(username, mobileNo, code) //正式订单
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(result)
	return
}

//是否是新用户
func isNewUser(mobile string) bool {
	return false
}

//是否购买过某个产品
func hasOrdered(mobile string, productCode string) bool {
	return false
}

//下正式订单
func addFinalOrder(userName string, mobile string, productCode string) entity.BaseResultEntity {
	result := entity.GetBaseFailRes()
	soaAddOrderRes, serverRes := proxy.AddSoaOrder(getMappingOrderNo())
	if serverRes.IsSucess && soaAddOrderRes.ErrorCode == "200" {
		result.IsSucess = true
		orderNo := soaAddOrderRes.OrderNo
		sendOrderSucessMessage(orderNo, mobile) //发送成功下单短信
		result.Message = "响应成功"
	} else {
		result.IsSucess = false
		result.Message = soaAddOrderRes.ErrorMessage
		result.Code = soaAddOrderRes.ErrorCode
	}
	return result
}

//获取订单号待实现
func getMappingOrderNo() string {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	return "P" + timestamp
}

func sendOrderSucessMessage(orderNo string, mobileNo string) bool {
	//获取院余号
	soaGetOrderDetailResOut, serverRes := proxy.GetSoaOrderDetail(orderNo)
	if !serverRes.IsSucess {
		return false;
	}
	if soaGetOrderDetailResOut.ErrorCode == "200" && len(soaGetOrderDetailResOut.OrderList) > 0 {
		orderDetail := soaGetOrderDetailResOut.OrderList[0]
		if len(orderDetail.DetailList) > 0 {
			yuanYuNo := orderDetail.DetailList[0].CardNo
			productName:=orderDetail.DetailList[0].ProdName
			smsContent := fmt.Sprintf(config.ConfigInfo.SmsOfOrderSucess,productName,yuanYuNo)
			proxy.SendMsg(mobileNo, smsContent)
			return true
		}
	}
	return false
}
