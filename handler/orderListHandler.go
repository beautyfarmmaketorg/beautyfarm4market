package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"beautyfarm4market/dal"
	"time"
	"strings"
)

func OrderListHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	mobile := r.FormValue("mobile")
	mappingOrderno := r.FormValue("mappingOrderno")
	wechatorderno := r.FormValue("wechatorderno")
	cardno := r.FormValue("cardno")
	args := Args{
		Name:           name,
		Mobile:         mobile,
		MappingOrderno: mappingOrderno,
		Wechatorderno:  wechatorderno,
		Cardno:         cardno,
	}

	allOrdersKey := "allOrdersKeys"
	cacheValue := util.GetCache(allOrdersKey)
	allOrders := []dal.TempOrder{}
	orders := []dal.TempOrder{}
	if cacheValue == nil {
		allOrders = dal.GetAllOrders()
		if len(allOrders) > 0 {
			util.SetCache(allOrdersKey, allOrders, 2*time.Minute)
		}
	} else {
		allOrders = cacheValue.([]dal.TempOrder)
	}
	for _, order := range allOrders {
		if (name == "" || strings.Index(order.UserName, name) > -1) &&
			(mobile == "" || strings.Index(order.MobileNo, mobile) > -1) &&
			(mappingOrderno == "" || strings.Index(order.MappingOrderNo, mappingOrderno) > -1) &&
			(wechatorderno == "" || strings.Index(order.WechatorderNo, wechatorderno) > -1) &&
			(cardno == "" || strings.Index(order.CardNo, cardno) > -1) {
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

func GetOrderExcel()  {
	
}

type Args struct {
	Name           string
	Mobile         string
	MappingOrderno string
	Wechatorderno  string
	Cardno         string
	Total          int
}
