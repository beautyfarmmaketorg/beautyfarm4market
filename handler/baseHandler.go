package handler

import "net/http"

func SafeHandler(fn http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e,ok:=recover().(error);ok{
				http.Error(w,e.Error(),http.StatusInternalServerError)
			}
		}()
		fn(w,r)
	}
}