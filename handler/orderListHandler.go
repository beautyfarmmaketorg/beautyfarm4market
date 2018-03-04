package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"beautyfarm4market/dal"
	"time"
	"strings"
	"beautyfarm4market/config"
	"strconv"
)

func OrderListHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	mobile := r.FormValue("mobile")
	mappingOrderno := r.FormValue("mappingOrderno")
	wechatorderno := r.FormValue("wechatorderno")
	cardno := r.FormValue("cardno")
	beginDateStr := r.FormValue("beginDateStr")
	endDateStr := r.FormValue("endDateStr")
	payStatusStr := r.FormValue("payStatusStr")
	channel := r.FormValue("channel")
	if channel == "" {
		channel = "-1"
	}
	payStatus := -1
	if payStatusStr != "" {
		payStatus, _ = strconv.Atoi(payStatusStr)
	}

	if beginDateStr == "" {
		d, _ := time.ParseDuration("-24h")
		beginDateStr = time.Now().Add(d).Format(config.ConfigInfo.DateLayout)
	}

	loc, _ := time.LoadLocation("Local")
	beginDate, _ := time.ParseInLocation(config.ConfigInfo.DateLayout, beginDateStr, loc)
	endDate, _ := time.ParseInLocation(config.ConfigInfo.DateLayout, endDateStr, loc)

	args := Args{
		Name:           name,
		Mobile:         mobile,
		MappingOrderno: mappingOrderno,
		Wechatorderno:  wechatorderno,
		Cardno:         cardno,
		BeginDateStr:   beginDateStr,
		EndDateStr:     endDateStr,
		PayStatusStr:   payStatusStr,
		Channel:        channel,
		Channels:       dal.GetAllChannel(),
	}

	allOrdersKey := "allOrdersKeys"
	cacheValue := util.GetCache(allOrdersKey)
	allOrders := []dal.TempOrder{}
	orders := []dal.TempOrder{}
	if cacheValue == nil {
		allOrders = dal.GetAllOrders()
		if len(allOrders) > 0 {
			//util.SetCache(allOrdersKey, allOrders, 2*time.Minute)
		}
	} else {
		allOrders = cacheValue.([]dal.TempOrder)
	}
	for _, order := range allOrders {
		createDate, _ := time.ParseInLocation(config.ConfigInfo.TimeLayout, order.CreateDate, loc)
		if (name == "" || strings.Index(order.UserName, name) > -1) &&
			(mobile == "" || strings.Index(order.MobileNo, mobile) > -1) &&
			(mappingOrderno == "" || strings.Index(order.MappingOrderNo, mappingOrderno) > -1) &&
			(wechatorderno == "" || strings.Index(order.WechatorderNo, wechatorderno) > -1) &&
			(cardno == "" || strings.Index(order.CardNo, cardno) > -1) &&
			(beginDateStr == "" || createDate.After(beginDate)) &&
			(endDateStr == "" || createDate.Before(endDate)) &&
			(payStatus == -1 || order.PayStatus == payStatus) &&
			(channel == "-1" || order.Channel == channel) {
			orders = append(orders, order)
		}
	}
	args.Total = len(orders)
	locals := make(map[string]interface{})
	locals["orders"] = orders
	locals["args"] = args
	if err := util.RenderHtml(w, "orderlist.html", locals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

func GetOrderExcel() {

}

type Args struct {
	Name           string
	Mobile         string
	MappingOrderno string
	Wechatorderno  string
	Cardno         string
	Total          int
	BeginDateStr   string
	EndDateStr     string
	PayStatusStr   string
	Channel        string
	Channels       []string
}
