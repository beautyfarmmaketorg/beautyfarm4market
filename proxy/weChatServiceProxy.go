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
	Errcode    int `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Ticket     string `json:"ticket"`
	Expires_in int `json:"expires_in"`
}
