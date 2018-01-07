package proxy

import (
	"beautyfarm4market/entity"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

func httpPostProxy(url string, req interface{}, res interface{}) entity.BaseResultEntity {
	baseRes := entity.GetBaseFailRes()
	jsonObj, err := json.Marshal(req)
	if err != nil {
		baseRes.IsSucess = false
		baseRes.Message = err.Error()
		return baseRes
	}
	//jsonStr := string(jsonObj)
	//fmt.Printf(jsonStr)
	response, postErr := http.Post(url, "application/json", bytes.NewBuffer([]byte(jsonObj)))
	defer response.Body.Close()
	if postErr != nil {
		baseRes.IsSucess = false
		baseRes.Message = postErr.Error()
		return baseRes
	}
	if response.StatusCode == 200 {
		baseRes = entity.GetBaseSucessRes()
	}
	body, _ := ioutil.ReadAll(response.Body)
	s := string(body)
	json.Unmarshal([]byte(s), res);
	return baseRes
}

func httpGetProxy(url string, res interface{}) entity.BaseResultEntity {
	baseRes := entity.GetBaseFailRes()
	response, getErr := http.Get(url)
	defer response.Body.Close()
	if getErr != nil {
		baseRes.IsSucess = false
		baseRes.Message = getErr.Error()
		return baseRes
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		s := string(body)
		json.Unmarshal([]byte(s), res);
		return entity.GetBaseSucessRes()
	}
	return baseRes
}
