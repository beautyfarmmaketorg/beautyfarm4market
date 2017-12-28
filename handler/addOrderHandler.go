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
	"net/url"
)

func AddOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	result := AddOrderResponse{};
	result.IsSucess = true
	username := r.FormValue("username")
	mobileNo := r.FormValue("mobileNo")
	code := r.FormValue("code")
	productIdStr := r.FormValue("productId")
	channelcode := r.FormValue("channelcode")
	if productIdStr == "" {
		productIdStr = "1"
	}
	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	productInfo := dal.GetProductInfo(productId)
	messagecCodeCookieName := fmt.Sprintf(config.ConfigInfo.CodeCookie, mobileNo)
	cookieCode, err := r.Cookie(messagecCodeCookieName)
	if err == nil {
		if code != cookieCode.Value {
			result.Code = "-1"
			result.Message = "请输入正确的验证码"
			json.NewEncoder(w).Encode(result)
			return
		}
	} else {
		result.Code = "-2"
		result.Message = "请获取验证码"
		json.NewEncoder(w).Encode(result)
		return
	}
	productCode := productInfo.Product_code // r.FormValue("productCode")
	totalPrice := 1                         //1分钱
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
	mappingOrderNo, _ := addTempOrder(username, mobileNo, productId, accountNo, channelcode, clientIp) //正式订单
	userAgent := r.Header.Get("User-Agent")
	dal.AddLog(dal.LogInfo{Title: "User-Agent", Description: userAgent, Type: 1})
	if mappingOrderNo != "" {
		if !strings.Contains(userAgent, "MicroMessenger") {
			result.Code = "3" //成功下单跳转支付
			weChatUnifiedorderResponse := InvokeWeChatUnifiedorder(productCode, "product", mappingOrderNo,
				clientIp, totalPrice, r.Host, "MWEB", "")
			if weChatUnifiedorderResponse.ReturnCode == "SUCCESS" && weChatUnifiedorderResponse.MwebUrl != "" {
				dal.UpdateTempOrderPayStatus(mappingOrderNo, 1) //更新支付状态
				host := r.Host
				if ! strings.Contains(host, "http") {
					host = "http://" + host
				}
				redirect_url := url.QueryEscape(host + "/purchaseRes?mappingOrderNo=" + mappingOrderNo)
				result.Redirect = weChatUnifiedorderResponse.MwebUrl + "&redirect_url=" + redirect_url
				setMobileCodeCookie(w, mobileNo, "", -1)
			}
		} else {
			//微信环境 获取openid
			result.Code = "3" //成功下单跳转支付
			redirectURI := url.QueryEscape("http://" + r.Host + "/handlerWeChatLogin?mappingOrderNo=" + mappingOrderNo)
			weChatLoginUrl := fmt.Sprintf(config.ConfigInfo.WeChatLoginUrl, redirectURI)
			dal.AddLog(dal.LogInfo{Title: "weChatLoginUrl", Description: weChatLoginUrl, Type: 1})
			result.Redirect = weChatLoginUrl
		}
	}
	json.NewEncoder(w).Encode(result)
	return
}

type AddOrderResponse struct {
	entity.BaseResultEntity
	Redirect string `json:"redirect"` //微信外环境 表示拉起微信支付 微信内表示身份认证
}

//添加临时单
func addTempOrder(userName string, mobile string, productId int64, accountNo string, channel string, clientIp string) (mappingOrderNo string, res entity.BaseResultEntity) {
	res = entity.GetBaseSucessRes()
	mappingOrderNo = getMappingOrderNo()
	p := dal.GetProductInfo(productId)
	t := dal.TempOrder{
		MappingOrderNo: mappingOrderNo,
		UserName:       userName,
		MobileNo:       mobile,
		ProductCode:    p.Product_code,
		AccountNo:      accountNo,
		Channel:        channel,
		CreateDate:     time.Unix(time.Now().Unix(), 0).Format(config.ConfigInfo.TimeLayout),
		ModifyDate:     time.Unix(time.Now().Unix(), 0).Format(config.ConfigInfo.TimeLayout),
		TotalPrice:     p.Price,
		ProductName:    config.ConfigInfo.ProductName,
		OrderStatus:    1,
		ClientIp:       clientIp,
		OrignalPrice:   p.Orignal_price,
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
	accountRegisterReq.Gender = "2"
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
