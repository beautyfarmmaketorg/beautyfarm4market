package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"strconv"
	"beautyfarm4market/dal"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		channelcode := r.FormValue("channelcode")
		productIdStr := r.FormValue("productId")
		if productIdStr == "" {
			productIdStr = "1"
		}
		productId, err := strconv.ParseInt(productIdStr, 10, 64)
		clientIp := r.RemoteAddr
		pageUrl := r.Host + "/index"
		dal.AddViewLog(dal.ViewLog{Channel_code: channelcode, Pange_url: pageUrl, Client_ip: clientIp})
		if err == nil {
			p := dal.GetProductInfo(productId)
			if p.Product_id == 0 {
				util.RenderHtml(w, "notfound.html", nil)
				return
			}
			locals := make(map[string]interface{})
			pageInfo := PageInfo{Channelcode: channelcode, ProductId: productIdStr}
			locals["pageInfo"] = pageInfo
			util.RenderHtml(w, "index.html", locals)
			return
		} else {
			util.RenderHtml(w, "notfound.html", nil)
		}
	}
	return
}

type PageInfo struct {
	Channelcode string
	ProductId   string
}
