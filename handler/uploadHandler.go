package handler

import (
	"net/http"
	"io"
	"os"
	"beautyfarm4market/util"
	"path"
	"beautyfarm4market/entity"
	"encoding/json"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	res := entity.GetBaseFailRes();
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		f, h, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filename := r.FormValue("fileName") + path.Ext(h.Filename)
		if filename == "" {
			filename = h.Filename
		}
		defer f.Close()
		path := util.GetCurrentPath() + "/html/images/" + filename;
		if isExist, _ := PathExists(path); isExist {
			os.Remove(path)
		}
		t, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.IsSucess = true
		res.Message = filename
	}
	json.NewEncoder(w).Encode(res)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
