package util

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"mahonia"
	"strconv"
)

func SendMsg(mobile string,content string)bool {
	isSucess:=false
	hex:=toHex(content)
	response,err:=http.Get("http://esms10.10690007.net/sms/mt?command=MT_REQUEST&spid=9508&sppassword=16r8XYC3&da=86"+mobile+"&dc=15&sm="+hex)
	if err==nil {
		defer response.Body.Close()
		if body,err:=ioutil.ReadAll(response.Body);err==nil {
			isSucess=checkSendMsgRes(string(body))
			fmt.Println(string(body));
		}
	}
	return isSucess
}

func checkSendMsgRes(responseBody string) bool{
	isSucess:=false
	arr:=strings.Split(responseBody,"&")
	for _,x:=range arr{
		innerArr :=strings.Split(x,"=")
		if len(innerArr)==2 {
			if innerArr[0]=="mtstat"&&innerArr[1]=="ACCEPTD" {
				isSucess = true
				break
			}
		}
	}
	return isSucess
}

func toHex(words string) string {
	enc:=mahonia.NewEncoder("GBK")
	hex:=""
	if output,ok:=enc.ConvertStringOK(words);ok{
		for _,c:=range []byte(output){
			hex+= strconv.FormatInt(int64(c), 16)
		}
	}
	return hex;
}