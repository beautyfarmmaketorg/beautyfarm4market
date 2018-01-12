package proxy

import (
	"beautyfarm4market/entity"
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"crypto/tls"
	"crypto/x509"
)

func httpPostProxy(url string, req interface{}, res interface{}) entity.BaseResultEntity {
	baseRes := entity.GetBaseFailRes()
	jsonObj, err := json.Marshal(req)
	if err != nil {
		baseRes.IsSucess = false
		baseRes.Message = err.Error()
		return baseRes
	}
	jsonStr := string(jsonObj)
	fmt.Printf(jsonStr)
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
	fmt.Printf(s)
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

func GetClientWithCa(rootCa, rootKey string) *http.Client {
	var tr *http.Transport
	certs, err := tls.LoadX509KeyPair(rootCa, rootKey)
	if err != nil {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		ca, err := x509.ParseCertificate(certs.Certificate[0])
		if err != nil {
			return &http.Client{Transport: tr}
		}
		pool := x509.NewCertPool()
		pool.AddCert(ca)

		tr = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		}

	}
	return &http.Client{Transport: tr}
}
