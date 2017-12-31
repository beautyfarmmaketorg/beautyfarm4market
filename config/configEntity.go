package config

import (
	"fmt"
	"io/ioutil"
	"beautyfarm4market/util"
	"encoding/xml"
)

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
	CodeCookie          string
	TimeLayout          string
	RegisterChannelType string
	PosService          string
	PosServiceAppId     string
	AppSecret           string
	IsVipUrl            string
	SignTemplate        string
	ProductCode         string
	ProductName         string
	UnifiedorderUrl     string //微信统一下单接口
	WeChatKey           string
	WeChatLoginUrl      string //获取微信authCode
	WeChatAppId         string //微信公众号appid
	WeChatAuthUrl       string //获取openid
	WeChatSecret        string //微信秘钥
	ProxyDir            string
	Version             string
	Dbname              string
}

func init() {
	cfg := getConfig()
	ConfigInfo = configEntity{
		SmsUrl: "http://esms10.10690007.net/sms/mt",
		//下单接口配置
		OrderServiceUrl:    cfg.OrderServiceUrl,
		AddOrderUrl:        "order/orderInsertOrUpdate",
		GetOrderDetailUrl:  "order/orderDetailSelect?orderList=%s&appId=%s",
		AccountRegisterUrl: "account/accountRegister",
		GetAccountInfoUrl:  "account/accountListSelect?phone=%s",
		OrderServiceAppId:  "TEST",
		AccountNo:          "129147",
		Channel:            "XS0001",
		//下单接口配置END
		SmsOfOrderSucess:    "您已成功购买%s，凭证号为：%s。请直接致电门店预约后到店护理或拨打4008206142预约咨询。",
		SmsOfVaild:          "%s（美丽田园手机验证码，请完成验证）， 如非本人操作，请忽略本短信。",
		MobileCookie:        "code%s",
		CodeCookie:          "messagecode%s",
		TimeLayout:          "2006-01-02 15:04:05",
		RegisterChannelType: "3003",
		PosService:          cfg.PosService,
		IsVipUrl:            "customer/isVipMember?org_no=beautyfarm&mobile=%s&appid=%s&timestamp=%d&sign=%s",
		SignTemplate:        "appid=%s&secretkey=%s&timestamp=%d",
		PosServiceAppId:     "bf_market",
		AppSecret:           "Vit+HmAG8a+7JCyIEPmR5A==",
		ProductCode:         "1110300002",
		ProductName:         "纯新胶原精华护理",
		UnifiedorderUrl:     "https://api.mch.weixin.qq.com/pay/unifiedorder",
		WeChatKey:           "meilitianyuan2016isgood2016igood",
		WeChatLoginUrl:      "http://bfwechat.beautyfarm.com.cn/Url.aspx?url=%s&scope=snsapi_base,snsapi_login",
		WeChatAppId:         "wx7302aaa9857c055b",
		WeChatAuthUrl:       "https://api.weixin.qq.com/sns/oauth2/access_token?appid=?&secret=?&code=?&grant_type=authorization_code",
		WeChatSecret:        "077cb94b72154e7b3d7db95ba40a83cb",
		ProxyDir:            "beautyfarm4market",
		Version:             "v1",
		Dbname:              cfg.Dbname,
	}
	fmt.Printf("init Config")
}

type config struct {
	Dbname          string `xml:"dbname"`
	PosService      string `xml:"posService"`
	OrderServiceUrl string `xml:"orderServiceUrl"`
}

func getConfig() config {
	cfg := config{}
	absoluteViewsDir := util.GetCurrentPath() + "/config.xml"
	data, err := ioutil.ReadFile(absoluteViewsDir)
	if err == nil {
		xml.Unmarshal(data, &cfg)
	}
	return cfg
}
