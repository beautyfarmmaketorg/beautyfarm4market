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
	result := AddOrderResponse{};
	result.IsSucess = true
	username := r.FormValue("username")
	mobileNo := r.FormValue("mobileNo")
	//code := r.FormValue("code")
	productCode := config.ConfigInfo.ProductCode // r.FormValue("productCode")
	totalPrice := 1                              //1分钱
	clientIp := r.Header.Get("Remote_addr")
	if (clientIp == "") {
		clientIp = r.RemoteAddr
	}
	if strings.Index(clientIp, "[::1]") > -1 {
		clientIp = "223.104.210.86:42425"
	}
	if clientIp != "" {
		arr := strings.Split(clientIp, ":")
		if len(arr) == 2 {
			clientIp = arr[0]
		}
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
	mappingOrderNo, _ := addTempOrder(username, mobileNo, config.ConfigInfo.ProductCode,
		config.ConfigInfo.ProductName, accountNo, 1) //正式订单
	if mappingOrderNo != "" {
		result.Code = "3" //成功下单跳转支付
		weChatUnifiedorderResponse := InvokeWeChatUnifiedorder(productCode, "product", mappingOrderNo,
			clientIp, totalPrice, r.Host)
		if weChatUnifiedorderResponse.ReturnCode == "SUCCESS" && weChatUnifiedorderResponse.MwebUrl != "" {
			dal.UpdateTempOrderPayStatus(mappingOrderNo, 1) //更新支付状态
			host := r.Host
			if ! strings.Contains(host, "http") {
				host = "http://" + host
			}
			redirect_url := host + "/purchaseRes?mappingOrderNo=" + mappingOrderNo
			result.PayUrl = weChatUnifiedorderResponse.MwebUrl + "&redirect_url=" +redirect_url
		}
	}
	json.NewEncoder(w).Encode(result)
	return
}

type AddOrderResponse struct {
	entity.BaseResultEntity
	PayUrl string
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
		OrderStatus:    1,
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

type AddFinalOrderRes struct {
	entity.BaseResultEntity
	CardNo  string
	OrderNo string
}

//下正式订单
func addFinalOrder(userName string, mobile string, productCode string, accountNo string, mappingOrderNo string) AddFinalOrderRes {
	result := AddFinalOrderRes{}
	soaAddOrderRes, serverRes := proxy.AddSoaOrder(mappingOrderNo, accountNo)
	if serverRes.IsSucess && soaAddOrderRes.ErrorCode == "200" {
		result.IsSucess = true
		orderNo := soaAddOrderRes.OrderNo
		cardNo, productName := getCardNo(orderNo)
		result.CardNo = cardNo
		result.OrderNo = soaAddOrderRes.OrderNo
		sendOrderSucessMessage(cardNo, productName, mobile) //发送成功下单短信
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

func sendOrderSucessMessage(cardNo, productName, mobileNo string) bool {
	smsContent := fmt.Sprintf(config.ConfigInfo.SmsOfOrderSucess, productName, cardNo)
	proxy.SendMsg(mobileNo, smsContent)
	return true
}

func getCardNo(orderNo string) (cardNo, productName string) {
	//获取院余号
	soaGetOrderDetailResOut, serverRes := proxy.GetSoaOrderDetail(orderNo)
	if !serverRes.IsSucess {
		return
	}
	if soaGetOrderDetailResOut.ErrorCode == "200" && len(soaGetOrderDetailResOut.OrderList) > 0 {
		orderDetail := soaGetOrderDetailResOut.OrderList[0]
		if len(orderDetail.DetailList) > 0 {
			cardNo = orderDetail.DetailList[0].CardNo
			productName = orderDetail.DetailList[0].ProdName
		}
	}
	return
}
