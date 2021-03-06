package handler

import (
	"net/http"
	"beautyfarm4market/dal"
	"beautyfarm4market/util"
	"strconv"
	"encoding/json"
	"beautyfarm4market/entity"
)

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		allProducts := dal.GetAllProductInfos()
		locals := make(map[string]interface{})
		locals["allProducts"] = allProducts
		util.RenderHtml(w, "productList.html", locals)
		return
	}
}

func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	res := entity.GetBaseSucessRes();
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if r.Method == "GET" {
		productIdStr := r.FormValue("productId")
		productId, _ := strconv.ParseInt(productIdStr, 10, 64)
		product := dal.GetProductInfo(productId, false)
		json.NewEncoder(w).Encode(product)
		return
	} else if r.Method == "POST" {
		productIdJson := r.FormValue("product")
		productInfo := dal.ProductInfo{}
		json.Unmarshal([]byte(productIdJson), &productInfo)
		if productInfo.Product_id > 0 {
			updateRes := dal.UpdateProductInfo(productInfo)
			res.IsSucess = updateRes
		} else {
			addRes := dal.AddProductInfo(productInfo);
			res.IsSucess=addRes
		}
		json.NewEncoder(w).Encode(res);
		return
	}
}
