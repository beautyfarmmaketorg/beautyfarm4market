package config

import "fmt"

var ConfigInfo configEntity

type configEntity struct {
	SmsUrl              string
	OrderServiceUrl     string
	AddOrderUrl         string
	AccountRegisterUrl  string
	OrderServiceAppId   string
	AccountNo           string
	Channel             string
	SmsOfOrderSucess    string
	GetOrderDetailUrl   string
	GetAccountInfoUrl   string
	SmsOfVaild          string
	MobileCookie        string
	TimeLayout          string
	RegisterChannelType string
	PosService          string
	PosServiceAppId     string
	AppSecret           string
	IsVipUrl            string
	SignTemplate        string
	ProductCode         string
	ProductName         string
	WeChatPayUrl        string
	WeChatKey string
}

func init() {
	ConfigInfo = configEntity{
		SmsUrl: "http://esms10.10690007.net/sms/mt",
		//下单接口配置
		OrderServiceUrl:    "http://180.169.107.155:8866/api/v1/bf-cam-adapter/spec/%s",
		AddOrderUrl:        "order/orderInsertOrUpdate",
		GetOrderDetailUrl:  "order/orderDetailSelect?orderList=%s&appId=%s",
		AccountRegisterUrl: "account/accountRegister",
		GetAccountInfoUrl:  "account/accountListSelect?phone=%s",
		OrderServiceAppId:  "TEST",
		AccountNo:          "129147",
		Channel:            "XS0001",
		//下单接口配置END
		SmsOfOrderSucess:    "您已成功购买产品%s，您的院余号为%s",
		SmsOfVaild:          "%s（美丽田园手机验证码，请完成验证）， 如非本人操作，请忽略本短信。",
		MobileCookie:        "code%s",
		TimeLayout:          "2006-01-02 15:04:05",
		RegisterChannelType: "3003",
		PosService:          "http://pos.beautyfarm.com.cn:8070/%s",
		IsVipUrl:            "customer/isVipMember?org_no=beautyfarm&mobile=%s&appid=%s&timestamp=%d&sign=%s",
		SignTemplate:        "appid=%s&secretkey=%s&timestamp=%d",
		PosServiceAppId:     "bf_market",
		AppSecret:           "Vit+HmAG8a+7JCyIEPmR5A==",
		ProductCode:         "1110300002",
		ProductName:         "纯新胶原精华护理",
		WeChatPayUrl:        "https://api.mch.weixin.qq.com/pay/unifiedorder",
		WeChatKey:"meilitianyuan2016isgood2016igood",
	}
	fmt.Printf("init Config")
}
