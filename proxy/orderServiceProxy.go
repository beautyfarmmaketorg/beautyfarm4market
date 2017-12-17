package proxy

import (
	"beautyfarm4market/entity"
	"beautyfarm4market/config"
	"fmt"
)

func AddSoaOrder(mappingOrderNo string) (soaAddOrderResOut SoaAddOrderRes, baseResultEntity entity.BaseResultEntity) {
	var soaAddOrderRes SoaAddOrderRes
	soaAddOrderReq := getSoaAddOrderReq(mappingOrderNo)
	url:=fmt.Sprintf(config.ConfigInfo.OrderServiceUrl,config.ConfigInfo.AddOrderUrl)
	baseResultEntity = httpPostProxy(url, soaAddOrderReq, &soaAddOrderRes)
	soaAddOrderResOut = soaAddOrderRes
	return
}

//获取下单请求
func getSoaAddOrderReq(mappingOrderNo string) SoaOrderDetai {
	soaAddOrderReq := SoaOrderDetai{
		AppId:          config.ConfigInfo.OrderServiceAppId,
		ModifyType:     "1",
		MappingOrderNo: mappingOrderNo,
		AccountNo:      config.ConfigInfo.AccountNo,
		Channel:        config.ConfigInfo.Channel,
		OrderType:      "2",
		OrderStatus:    "6",
		DetailList: []SoaProductDetail{
			{
				DetailListNo: "01",
				ProdCategory: "32",
				ProdNo:       "1110300002",
				ProdName:     "纯新胶原精华护理",
				ProdUnit:     "件",
				OrderQty:     "1",
				ProdPrice:    "1",
				ProdAmt:      "1",
				OrderPrice:   "1",
				OrderAmt:     "1",
				PayList: []SoaPayInfoDetail{
					{
						PayNo:       "P000004",
						PayCategory: "3",
						PayType:     "第三方支付",
						PayAmt:      "1",
						PayTimes:    "1",
					},
				},
			},
		},
	}
	return soaAddOrderReq
}

//订单服务共有请求字段
type SoaBaseReq4Orderservice struct {

}

//订单服务共有响应字段
type SoaBaseRes4Orderservice struct {
	ErrorCode      string `json:"errorCode"`
	ErrorMessage   string `json:"errorMessage"`
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
func GetSoaOrderDetail(orderNo string)(soaGetOrderDetailResOut SoaGetOrderDetailRes, serverRes entity.BaseResultEntity)  {
	var soaGetOrderDetailRes SoaGetOrderDetailRes
	methodUrl:=fmt.Sprintf(config.ConfigInfo.OrderServiceUrl,config.ConfigInfo.GetOrderDetailUrl)
	url:=fmt.Sprintf(methodUrl,orderNo,config.ConfigInfo.OrderServiceAppId)
	serverRes= httpGetProxy(url, &soaGetOrderDetailRes)
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