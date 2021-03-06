package proxy

import (
	"beautyfarm4market/entity"
	"beautyfarm4market/config"
	"fmt"
	"beautyfarm4market/dal"
	"strconv"
)

func AddSoaOrder(mappingOrderNo string, accountNo string) (soaAddOrderResOut SoaAddOrderRes, baseResultEntity entity.BaseResultEntity) {
	var soaAddOrderRes SoaAddOrderRes
	soaAddOrderReq := getSoaAddOrderReq(mappingOrderNo, accountNo)
	url := fmt.Sprintf(config.ConfigInfo.OrderServiceUrl, config.ConfigInfo.AddOrUpdateOrderUrl)
	baseResultEntity = httpPostProxy(url, soaAddOrderReq, &soaAddOrderRes)
	soaAddOrderResOut = soaAddOrderRes
	return
}

func CancelSoaOrder(mappingOrderNo string,cancelmappingOrderNo string) (cancelSoaOrderResOut CancelSoaOrderRes, baseResultEntity entity.BaseResultEntity) {
	baseResultEntity = entity.GetBaseFailRes()
	c := config.ConfigInfo
	tempOrderInfo := dal.GetOrdersByMappingOrderNo(mappingOrderNo)
	if tempOrderInfo.MappingOrderNo == "" {
		baseResultEntity.Message = "订单不存在"
		return
	}
	cancelSoaOrderReq := CancelSoaOrderReq{
		AppId:          c.OrderServiceAppId,
		ModifyType:     "1",
		MappingOrderNo: cancelmappingOrderNo,
		AccountNo:      tempOrderInfo.AccountNo,
		Channel:        c.Channel,
		OrderType:      "3",
		OrderStatus:    "6",
		DetailList: []SoaProductDetail{
			{
				DetailListNo: "01",
				ProdCategory: "32",
				ProdNo:       tempOrderInfo.ProductCode,
				ProdName:     tempOrderInfo.ProductName,
				ProdUnit:     "件",
				OrderQty:     "-1",
				ProdPrice:    strconv.FormatFloat(tempOrderInfo.OrignalPrice, 'f', 2, 64),
				ProdAmt:      strconv.FormatFloat(-tempOrderInfo.OrignalPrice, 'f', 2, 64),
				OrderPrice:   strconv.FormatFloat(tempOrderInfo.TotalPrice, 'f', 2, 64),
				OrderAmt:     strconv.FormatFloat(-tempOrderInfo.TotalPrice, 'f', 2, 64),
				PayList: []SoaPayInfoDetail{
					{
						PayNo:       mappingOrderNo,
						PayCategory: "3",
						PayType:     "第三方支付",
						PayAmt:      strconv.FormatFloat(-tempOrderInfo.TotalPrice, 'f', 2, 64),
						PayTimes:    "-1",
					},
				},
				CardNo: tempOrderInfo.CardNo,
			},
		},
	}
	newDetailList := appendOtherProducts(tempOrderInfo.ProductCode, mappingOrderNo, tempOrderInfo.CardNo, "-1",-1,"-1")
	if len(newDetailList) > 0 {
		cancelSoaOrderReq.DetailList = newDetailList;
	}
	url := fmt.Sprintf(config.ConfigInfo.OrderServiceUrl, config.ConfigInfo.AddOrUpdateOrderUrl)
	var cancelSoaOrderRes CancelSoaOrderRes
	baseResultEntity = httpPostProxy(url, cancelSoaOrderReq, &cancelSoaOrderRes)
	cancelSoaOrderResOut = cancelSoaOrderRes
	return
}

type CancelSoaOrderReq struct {
	AppId          string             `json:"appId"`
	ModifyType     string             `json:"modifyType"`
	MappingOrderNo string             `json:"mappingOrderNo"`
	OrderNo        string             `json:"orderNo"`
	AccountNo      string             `json:"accountNo"`
	Channel        string             `json:"channel"`
	OrderType      string             `json:"orderType"`
	OrderStatus    string             `json:"orderStatus"`
	DetailList     []SoaProductDetail `json:"detailList"`
}

type CancelSoaOrderRes struct {
	SoaBaseRes4Orderservice
	OrderNo        string `json:"orderNo"`
	MappingOrderNo string `json:"mappingOrderNo"`
}

//获取下单请求
func getSoaAddOrderReq(mappingOrderNo string, accountNo string) SoaOrderDetai {
	temporder := dal.GetOrdersByMappingOrderNo(mappingOrderNo)
	soaAddOrderReq := SoaOrderDetai{
		AppId:          config.ConfigInfo.OrderServiceAppId,
		ModifyType:     "1",
		MappingOrderNo: temporder.MappingOrderNo,
		AccountNo:      accountNo,
		Channel:        config.ConfigInfo.Channel,
		OrderType:      "2",
		OrderStatus:    "6",
		DetailList: []SoaProductDetail{
			{
				DetailListNo: "01",
				ProdCategory: "32",
				ProdNo:       temporder.ProductCode,
				ProdName:     temporder.ProductName,
				ProdUnit:     "件",
				OrderQty:     "1",
				ProdPrice:    strconv.FormatFloat(temporder.OrignalPrice, 'f', 2, 64),
				ProdAmt:      strconv.FormatFloat(temporder.OrignalPrice, 'f', 2, 64),
				OrderPrice:   strconv.FormatFloat(temporder.TotalPrice, 'f', 2, 64),
				OrderAmt:     strconv.FormatFloat(temporder.TotalPrice, 'f', 2, 64),
				PayList: []SoaPayInfoDetail{
					{
						PayNo:       mappingOrderNo,
						PayCategory: "3",
						PayType:     "第三方支付",
						PayAmt:      strconv.FormatFloat(temporder.TotalPrice, 'f', 2, 64),
						PayTimes:    "1",
					},
				},
			},
		},
	}
	newDetailList := appendOtherProducts(temporder.ProductCode, mappingOrderNo, "", "1",1,"1")
	if len(newDetailList) > 0 {
		soaAddOrderReq.DetailList = newDetailList;
	}

	return soaAddOrderReq
}

func appendOtherProducts(productCode string, mappingOrderNo string, cardNo string, payTimes string,positive float64,orderQty string ) []SoaProductDetail {
	detailList := []SoaProductDetail{}
	if productCode == "1180100125" {
		detailList = []SoaProductDetail{
			{
				DetailListNo: "01",
				ProdCategory: "32",
				ProdNo:       "1180100125",
				ProdName:     "平衡修复水氧护理",
				ProdUnit:     "件",
				OrderQty:     orderQty,
				ProdPrice:    strconv.FormatFloat(980, 'f', 2, 64),
				ProdAmt:      strconv.FormatFloat(980*positive, 'f', 2, 64),
				OrderPrice:   strconv.FormatFloat(213, 'f', 2, 64),
				OrderAmt:     strconv.FormatFloat(213*positive, 'f', 2, 64),
				PayList: []SoaPayInfoDetail{
					{
						PayNo:       mappingOrderNo,
						PayCategory: "3",
						PayType:     "第三方支付",
						PayAmt:      strconv.FormatFloat(213*positive, 'f', 2, 64),
						PayTimes:    payTimes,
					},
				},
				CardNo: cardNo,
			},
			{
				DetailListNo: "02",
				ProdCategory: "32",
				ProdNo:       "1130600001",
				ProdName:     "芳香精油能量按摩",
				ProdUnit:     "件",
				OrderQty:     orderQty,
				ProdPrice:    strconv.FormatFloat(480, 'f', 2, 64),
				ProdAmt:      strconv.FormatFloat(480*positive, 'f', 2, 64),
				OrderPrice:   strconv.FormatFloat(105, 'f', 2, 64),
				OrderAmt:     strconv.FormatFloat(105*positive, 'f', 2, 64),
				PayList: []SoaPayInfoDetail{
					{
						PayNo:       mappingOrderNo,
						PayCategory: "3",
						PayType:     "第三方支付",
						PayAmt:      strconv.FormatFloat(105*positive, 'f', 2, 64),
						PayTimes:    payTimes,
					},
				},
				CardNo: cardNo,
			},
		}
	}
	return detailList
}

//订单服务共有请求字段
type SoaBaseReq4Orderservice struct {
	AppId string `json:"appId"`
}

//订单服务共有响应字段
type SoaBaseRes4Orderservice struct {
	Status       int    `json:"status"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

//下单接口请求实体
type SoaOrderDetai struct {
	AppId          string             `json:"appId"`
	ModifyType     string             `json:"modifyType"`
	MappingOrderNo string             `json:"mappingOrderNo"`
	AccountNo      string             `json:"accountNo"`
	Channel        string             `json:"channel"`
	OrderType      string             `json:"orderType"`
	OrderStatus    string             `json:"orderStatus"`
	DetailList     []SoaProductDetail `json:"detailList"`
}

//下单接口响应实体
type SoaAddOrderRes struct {
	SoaBaseRes4Orderservice
	OrderNo        string `json:"orderNo"`
	MappingOrderNo string `json:"mappingOrderNo"`
}

//产品明显
type SoaProductDetail struct {
	//院余号码
	CardNo       string             `json:"cardNo"`
	DetailListNo string             `json:"detailListNo"`
	ProdCategory string             `json:"prodCategory"`
	ProdNo       string             `json:"prodNo"`
	ProdName     string             `json:"prodName"`
	ProdUnit     string             `json:"prodUnit"`
	OrderQty     string             `json:"orderQty"`
	ProdPrice    string             `json:"prodPrice"`
	ProdAmt      string             `json:"prodAmt"`
	OrderPrice   string             `json:"orderPrice"`
	OrderAmt     string             `json:"orderAmt"`
	PayList      []SoaPayInfoDetail `json:"payList"`
}

//支付信息
type SoaPayInfoDetail struct {
	PayNo       string `json:"payNo"`
	PayCategory string `json:"payCategory"`
	PayType     string `json:"payType"`
	PayAmt      string `json:"payAmt"`
	PayTimes    string `json:"payTimes"`
}

//serverRes http请求结果
func GetSoaOrderDetail(orderNo string) (soaGetOrderDetailResOut SoaGetOrderDetailRes, serverRes entity.BaseResultEntity) {
	var soaGetOrderDetailRes SoaGetOrderDetailRes
	methodUrl := fmt.Sprintf(config.ConfigInfo.OrderServiceUrl, config.ConfigInfo.GetOrderDetailUrl)
	url := fmt.Sprintf(methodUrl, orderNo, config.ConfigInfo.OrderServiceAppId)
	serverRes = httpGetProxy(url, &soaGetOrderDetailRes)
	soaGetOrderDetailResOut = soaGetOrderDetailRes
	return
}

type OrderDetailInfo struct {
	OrderList []SoaProductDetail `json:"orderList"`
}

type SoaGetOrderDetailRes struct {
	SoaBaseRes4Orderservice
	OrderList []SoaOrderDetai `json:"orderList"`
}

//账户注册
func AddSoaAccount(accountRegisterReq AccountRegisterReq) (accountRegisterResOut AccountRegisterRes, baseResultEntity entity.BaseResultEntity) {
	var accountRegisterRes AccountRegisterRes
	accountRegisterReq.AppId = config.ConfigInfo.OrderServiceAppId
	accountRegisterReq.RegisterChannelType = config.ConfigInfo.RegisterChannelType
	url := fmt.Sprintf(config.ConfigInfo.OrderServiceUrl, config.ConfigInfo.AccountRegisterUrl)
	baseResultEntity = httpPostProxy(url, accountRegisterReq, &accountRegisterRes)
	accountRegisterResOut = accountRegisterRes
	return
}

type AccountRegisterReq struct {
	SoaBaseReq4Orderservice
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	RegisterChannelType string `json:"registerChannelType"`
	Gender              string `json:"gender"`
}

type AccountRegisterRes struct {
	SoaBaseRes4Orderservice
	AccountNo string `json:"accountNo"`
}

type SoaGetAccountInfoRes struct {
	SoaBaseRes4Orderservice
	AccountList []AccountInfo `json:"accountList"`
}

type AccountInfo struct {
	AccountNo string `json:"accountNo"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

//serverRes http请求结果
func GetSoaAccountInfo(mobileNo string) (soaGetAccountInfoResOut SoaGetAccountInfoRes, serverRes entity.BaseResultEntity) {
	var soaGetAccountInfoRes SoaGetAccountInfoRes
	methodUrl := fmt.Sprintf(config.ConfigInfo.OrderServiceUrl, config.ConfigInfo.GetAccountInfoUrl)
	url := fmt.Sprintf(methodUrl, mobileNo)
	serverRes = httpGetProxy(url, &soaGetAccountInfoRes)
	soaGetAccountInfoResOut = soaGetAccountInfoRes
	return
}
