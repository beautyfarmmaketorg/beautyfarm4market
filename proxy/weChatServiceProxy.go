package proxy

import (
	"beautyfarm4market/config"
	"fmt"
	"beautyfarm4market/entity"
)

func GetAccessToken() (accessTokenRes AccessToken, res entity.BaseResultEntity) {
	accessTokenRes = AccessToken{}
	c := config.ConfigInfo;
	url := fmt.Sprintf(c.WeChatTokenUrl, c.WeChatAppId, c.WeChatSecret)
	res = httpGetProxy(url, &accessTokenRes)
	return
}

func GetTicket() (ticket Ticket, res entity.BaseResultEntity) {

	ticket = Ticket{}
	c := config.ConfigInfo;
	tokenInfo, serverRes := GetAccessToken()
	accessToken := ""
	if serverRes.IsSucess {
		accessToken = tokenInfo.AccessToken
	}
	url := fmt.Sprintf(c.WeChatTicketUrl, accessToken)
	res = httpGetProxy(url, &ticket)
	return
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Ticket struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Ticket     string `json:"ticket"`
	Expires_in int    `json:"expires_in"`
}

func WeChatRefund(weChatRefundReq WeChatRefundReq) (weChatRefundRes WeChatRefundRes, res entity.BaseResultEntity) {
	weChatRefundRes = WeChatRefundRes{}
	c := config.ConfigInfo
	url := c.WeRefundUrl
	res = httpPostProxy(url, weChatRefundReq, &weChatRefundRes)
	return
}

type WeChatRefundReq struct {
	Appid          string
	Mch_id         string
	Nonce_str      string
	Out_refund_no  string
	Refund_fee     int
	Refund_desc    string
	Total_fee      int
	Transaction_id string
	Sign           string
}

type WeChatRefundRes struct {
	Return_code    string `json:"return_code"`
	Return_msg     string `json:"return_msg"`
	Appid          string `json:"appid"`
	Mch_id         string `json:"mch_id"`
	Nonce_str      string `json:"nonce_str"`
	Sign           string `json:"sign"`
	Result_code    string `json:"result_code"`
	Transaction_id string `json:"transaction_id"`
	Out_trade_no   string `json:"out_trade_no"`
	Out_refund_no  string `json:"out_refund_no"`
	Refund_id      string `json:"refund_id"`
	Refund_channel string `json:"refund_channel"`
	Refund_fee     int    `json:"refund_fee"`
}
