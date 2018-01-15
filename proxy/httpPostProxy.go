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

var _tlsConfig *tls.Config

func getTLSConfig() (*tls.Config, error) {
	if _tlsConfig != nil {
		return _tlsConfig, nil
	}

	// load cert
	cert, err := tls.LoadX509KeyPair("c://ca//apiclient_cert.pem", "c://ca//apiclient_key.pem")
	if err != nil {
		return nil, err
	}

	// load root ca
	caData, err := ioutil.ReadFile("c://ca//rootca.pem")
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return _tlsConfig, nil
}

func GetClientWithCa() *http.Client {
	tlsConfig, err := getTLSConfig()
	if err != nil {

	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: tr}
}
