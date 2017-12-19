package handler

import (
	"net/http"
	"beautyfarm4market/util"
	"beautyfarm4market/dal"
)

func OrderListHandler(w http.ResponseWriter, r *http.Request) {
	orders := dal.GetAllOrders()
	locals := make(map[string]interface{})
	locals["orders"] = orders
	if err := util.RenderHtml(w, "orderlist.html", locals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
