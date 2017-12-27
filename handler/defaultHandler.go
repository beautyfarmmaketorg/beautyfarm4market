package handler

import (
	"net/http"
	"io"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok")
}
