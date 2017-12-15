package handler

import (
	"net/http"
	"io"
)

func IndexHandler(w http.ResponseWriter,r *http.Request)  {
	io.WriteString(w,"hello,world!")
}
