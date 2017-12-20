package handler

import (
	"net/http"
	"beautyfarm4market/entity"
	"encoding/json"
	"beautyfarm4market/proxy"
	"fmt"
	"beautyfarm4market/config"
	"strconv"
	"beautyfarm4market/dal"
	"time"
	"strings"
)

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	result := entity.GetBaseSucessRes();
	username := r.FormValue("username")
	mobileNo := r.FormValue("mobileNo")
	//code := r.FormValue("code")
	productCode := config.ConfigInfo.ProductCode // r.FormValue("productCode")
	totalPrice := 1
	clientIp := r.Header.Get("Remote_addr")
	if (clientIp == "") {
		clientIp = r.RemoteAddr
	}
	if strings.Index(clientIp, "[::1]") > -1 {
		clientIp = "192.168.1.1"
	}
	fmt.Printf("clientIp:", clientIp);
	//检查是否已经下过订单
	if hasOrdered := checkHasOrdered(mobileNo, productCode); hasOrdered {
		result.Code = "1" //
		json.NewEncoder(w).Encode(result)
		return
	}

	//vip用户
	if isVip := isVip(mobileNo); isVip {
		result.Code = "2" //vip
		json.NewEncoder(w).Encode(result)
		return
	}

	//闪客且没有下过订单
	accountNo, _ := getAccountNo(mobileNo, username)
	mappingOrderNo, result := addTempOrder(username, mobileNo, config.ConfigInfo.ProductCode,
		config.ConfigInfo.ProductName, accountNo, 1) //正式订单
	if mappingOrderNo != "" {
		result.Code = "3" //成功下单跳转支付
		xmlStr := GetPayUrl(productCode, "product", mappingOrderNo,
			clientIp, strconv.Itoa(totalPrice))
		fmt.Printf(xmlStr)
	}
	//result = addFinalOrder(username, mobileNo, code, accountNo, mappingOrderNo) //正式订单

	json.NewEncoder(w).Encode(result)
	return
}

//添加临时单
func addTempOrder(userName string, mobile string, productCode string, productName string, accountNo string, channel int) (mappingOrderNo string, res entity.BaseResultEntity) {
	res = entity.GetBaseSucessRes()
	mappingOrderNo = getMappingOrderNo()
	t := dal.TempOrder{
		MappingOrderNo: mappingOrderNo,
		UserName:       userName,
		MobileNo:       mobile,
		ProductCode:    productCode,
		AccountNo:      accountNo,
		Channel:        channel,
		CreateDate:     time.Unix(time.Now().Unix(), 0).Format(config.ConfigInfo.TimeLayout),
		ModifyDate:     time.Unix(time.Now().Unix(), 0).Format(config.ConfigInfo.TimeLayout),
		TotalPrice:     1,
		ProductName:    config.ConfigInfo.ProductName,
	}
	isSucess := dal.AddTempOrder(t)
	res.IsSucess = isSucess
	return
}

//检查是否已经下过订单
func checkHasOrdered(mobileNo string, productCode string) bool {
	orders := dal.GetOrdersByMobile(mobileNo, productCode)
	return len(orders) > 0
}

//是否是新用户 是的话则注册并且 返回accountNo用于下单
func getAccountNo(mobile string, userName string) (accountNo string, isNewCreate bool) {
	accountNo = ""
	accountRegisterReq := proxy.AccountRegisterReq{}
	accountRegisterReq.Phone = mobile
	accountRegisterReq.Name = userName
	accountRegisterRes, serverRes := proxy.AddSoaAccount(accountRegisterReq)
	if serverRes.IsSucess {
		if accountRegisterRes.ErrorCode == "200" {
			accountNo = accountRegisterRes.AccountNo
			isNewCreate = true
		}
	} else if accountRegisterRes.Status == 400 {
		//账户已存在则调用查询账户信息接口 获取AccountNo
		getAccountInfoRes, serverRes4Account := proxy.GetSoaAccountInfo(mobile)
		if serverRes4Account.IsSucess && len(getAccountInfoRes.AccountList) > 0 {
			accountNo = getAccountInfoRes.AccountList[0].AccountNo
		}
	}
	return
}

//是否是新用户
func isVip(mobile string) bool {
	isVip := false
	soaIsVipResOut, serverRes := proxy.IsVip(mobile)
	if serverRes.IsSucess && soaIsVipResOut.Status == 1 {
		isVip = soaIsVipResOut.Data.IsVip || soaIsVipResOut.Data.IsMarketVip
	}
	return isVip
}

//是否购买过某个产品
func hasOrdered(mobile string, productCode string) bool {
	return false
}

//下正式订单
func addFinalOrder(userName string, mobile string, productCode string, accountNo string, mappingOrderNo string) entity.BaseResultEntity {
	result := entity.GetBaseFailRes()
	soaAddOrderRes, serverRes := proxy.AddSoaOrder(mappingOrderNo, accountNo)
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
			productName := orderDetail.DetailList[0].ProdName
			smsContent := fmt.Sprintf(config.ConfigInfo.SmsOfOrderSucess, productName, yuanYuNo)
			proxy.SendMsg(mobileNo, smsContent)
			return true
		}
	}
	return false
}
