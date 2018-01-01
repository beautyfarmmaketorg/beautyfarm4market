package handler

import (
	"net/http"
	"beautyfarm4market/dal"
	"beautyfarm4market/util"
)

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		viewLogReports := dal.GetViewLogReport()
		orderReports := dal.GetOrderReport()
		locals := make(map[string]interface{})
		locals["viewLogReports"] = viewLogReports
		locals["orderReports"] = orderReports
		util.RenderHtml(w, "report.html", locals)
		return
	}
}
