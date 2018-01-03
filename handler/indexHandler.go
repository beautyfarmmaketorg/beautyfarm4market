package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"strconv"
	"beautyfarm4market/dal"
	"html/template"
	"time"
	"sort"
	"fmt"
	"beautyfarm4market/proxy"
	"strings"
	"beautyfarm4market/config"
)

var ticketDic = make(map[string]TicketWithDate)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		channelcode := r.FormValue("channelcode")
		productIdStr := r.FormValue("productId")
		if productIdStr == "" {
			productIdStr = "1"
		}
		productId, err := strconv.ParseInt(productIdStr, 10, 64)
		clientIp := r.RemoteAddr
		pageUrl := "http://" + r.Host + r.RequestURI
		dal.AddViewLog(dal.ViewLog{Channel_code: channelcode, Pange_url: pageUrl, Client_ip: clientIp})
		if err == nil {
			p := dal.GetProductInfo(productId, true)
			if p.Product_id == 0 {
				util.RenderHtml(w, "notfound.html", nil)
				return
			}
			locals := make(map[string]interface{})
			shareImage := "http://" + r.Host + "/assets/images/liutao.jpg"
			pageInfo := PageInfo{Channelcode: channelcode, ProductId: productIdStr, Bg: p.Backgroud_image,
				Button: p.PurhchaseBtn_image, Rule: p.Rule_image, Mask: p.MaskImage, RuleDesc: template.HTML(p.Prodcut_rule),
				ProductName: p.Prodcut_name, ProductDesc: p.Prodcut_desc, PageUrl: pageUrl, ShareImage: shareImage}
			locals["pageInfo"] = pageInfo
			ticket := getTicketStr(w, r)
			locals["sign"] = getWeChatLoginParams(pageUrl, ticket)
			util.RenderHtml(w, "index.html", locals)
			return
		} else {
			util.RenderHtml(w, "notfound.html", nil)
		}
	}
	return
}

func getWeChatLoginParams(indexUrl string, ticket string) WeChatLoginParams {

	args := WeChatLoginParams{
		TimeStamp:   strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:    strconv.FormatInt(time.Now().Unix(), 10),
		JsapiTicket: ticket,
		IndexUrl:    indexUrl,
		AppId:       config.ConfigInfo.WeChatAppId,
	}
	sign := getSign4WeChat(args)
	args.PaySign = sign
	return args
}

func getTicketStr(w http.ResponseWriter, r *http.Request) string {
	ticketStr := ""
	ticketKey := "ticketKey"
	var ticketWithDate TicketWithDate
	if ticketDic != nil && len(ticketDic) > 0 {
		ticketWithDate = ticketDic[ticketKey]
	}

	if ticketWithDate.Ticket != "" {
		leftSecond := getCurrLeftSecond(ticketWithDate.DateTimeStr)
		if leftSecond > 100 {
			ticketStr = ticketWithDate.Ticket
		}
	}
	if ticketStr == "" {
		t, _ := proxy.GetTicket()
		ticketStr = t.Ticket
		setTicketDic(ticketKey, ticketStr)
	}
	return ticketStr
}

//返回相差的秒
func getCurrLeftSecond(cookTimeStr string) int {
	var leftSeconds float64 = -1
	loc, _ := time.LoadLocation("Local")
	cookTime, _ := time.ParseInLocation(config.ConfigInfo.TimeLayout, cookTimeStr, loc)
	leftSeconds = time.Now().Sub(cookTime).Seconds()
	return 7200 - int(leftSeconds)
}

func setTicketDic(ticketKey string, ticket string) {
	dataTimeStr := time.Unix(time.Now().Unix(), 0).Format(config.ConfigInfo.TimeLayout)
	ticketDic[ticketKey] = TicketWithDate{Ticket: ticket, DateTimeStr: dataTimeStr}
}

type TicketWithDate struct {
	Ticket      string
	DateTimeStr string
}

func getSign4WeChat(e WeChatLoginParams) string {
	m := make(map[string]interface{}, 0)
	m["timestamp"] = e.TimeStamp
	m["noncestr"] = e.NonceStr
	m["jsapi_ticket"] = e.JsapiTicket
	m["url"] = e.IndexUrl
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	var signStrings string
	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, m[k])
		value := fmt.Sprintf("%v", m[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}
	signStrings = strings.Trim(signStrings, "&")
	return util.GetSha1(signStrings)
}

type WeChatLoginParams struct {
	TimeStamp   string
	NonceStr    string
	PaySign     string
	IndexUrl    string
	JsapiTicket string
	AppId       string
}

type PageInfo struct {
	Channelcode string
	ProductId   string
	Bg          string
	Button      string
	Rule        string
	Mask        string
	RuleDesc    template.HTML
	ProductName string
	ProductDesc string
	PageUrl     string
	ShareImage  string
}
