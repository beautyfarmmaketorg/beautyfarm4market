package handler

import (
	"net/http"
	"encoding/json"
	"beautyfarm4market/proxy"
	"beautyfarm4market/bll"
	"beautyfarm4market/entity"
	"beautyfarm4market/dal"
)

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	res := entity.GetBaseFailRes()
	mappingOrderNo := r.FormValue("mappingOrderNo")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	cancelOrderNo := bll.GetOrderNo("CNC")
	cancelRes, serviceRes := proxy.CancelSoaOrder(mappingOrderNo, cancelOrderNo)
	if serviceRes.IsSucess && cancelRes.ErrorCode == "200" {
		res.IsSucess = true
		res.Message = "取消院余成功"
		dal.UpdateTempOrderStatus(mappingOrderNo, 3)
	} else {
		res.IsSucess = false
		res.Message = cancelRes.ErrorMessage
	}
	json.NewEncoder(w).Encode(res)
	return
}

func RefundOrderHandler(w http.ResponseWriter, r *http.Request) {
	res := entity.GetBaseFailRes()
	mappingOrderNo := r.FormValue("mappingOrderNo")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	refundRes := bll.Refund(mappingOrderNo, "后台取消")
	if refundRes.IsSucess {
		res.IsSucess = true
		res.Message = "退款成功"
		dal.UpdateTempOrderStatus(mappingOrderNo, 3)
	} else {
		res.IsSucess = false
		res.Message = "退款失败"
	}
	json.NewEncoder(w).Encode(res)
	return
}
