$(function () {
    var payStatusObj = $("#payStatus")
    if (payStatusObj != undefined) {
        var paystatus = payStatusObj.val();
        if (paystatus == "1") {
            setInterval(refreshStatus, 5000);
        }
    }

    function refreshStatus() {
        var processMsgObj = $("#processMsg");
        processMsgObj.html("后台查询中...")
        var paystatus = $("#payStatus").val()
        if (paystatus == "2") {
            return;//已完成支付
        }
        $.post("/checkPurchaseRes", {"mappingOrderNo": $("#mappingOrderNo").val()}, function (res) {
            if (res.isSucess && res.payStatus == 2) {
                $("#alert-payIng").removeClass("show");
                $("#alert-success").addClass("show");
                $("#payStatus").val(res.payStatus)
                processMsgObj.html("")
            }
        });
    }

    $("#sucessBtn").click(function () {
        var url = $("#indexUrl").val();
        if (url.indexOf("http") == -1) {
            url = "http://" + url;
        }
        window.open(url);
    });
})
