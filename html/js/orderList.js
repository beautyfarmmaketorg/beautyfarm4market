$(function () {
    $(".cancelBtn").bind("click", function () {
        var orderNo = $(this).attr("data");
        alert(orderNo);

        $.post("/cancelOrder", {
            "mappingOrderNo": orderNo
        }, function (res) {
            alert(res.message)
        });
    })

    $(".refundBtn").bind("click", function () {
        var orderNo = $(this).attr("data");
        alert(orderNo);
        $.post("/refundOrder", {
            "mappingOrderNo": orderNo
        }, function (res) {
            alert(res.message)
        });
    })

})