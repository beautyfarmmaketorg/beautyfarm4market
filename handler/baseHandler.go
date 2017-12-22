package handler

import (
	"net/http"
	"beautyfarm4market/dal"
	"reflect"
)

func SafeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		funName:=reflect.TypeOf(fn).Name()
		defer func() {
			if e, ok := recover().(error); ok {
				dal.AddLog(dal.LogInfo{Title: funName, Description: e.Error(), Type: 3})
				http.Error(w, e.Error(), http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}
