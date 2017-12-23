package proxy

import (
	"beautyfarm4market/config"
	"fmt"
	"beautyfarm4market/entity"
)

func GetAccessToken(code string) (accessTokenRes AccessToken, res entity.BaseResultEntity) {
	accessTokenRes = AccessToken{}
	c := config.ConfigInfo;
	url := fmt.Sprintf(c.WeChatAuthUrl, c.WeChatAppId, c.WeChatSecret, code)
	res = httpGetProxy(url, accessTokenRes)
	return

}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    string `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}
