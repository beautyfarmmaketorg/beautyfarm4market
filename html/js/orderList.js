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

    $.datepicker.setDefaults({
        dateFormat: "yy-mm-dd",
    });

    $("#beginDate").datepicker();
    $("#endDate").datepicker();

    var payStatus = $("#payStatusStrHidden").val();
    if(payStatus!=""){
        $("#payStatusSel").find("option[value='"+payStatus+"']").attr("selected",true);
    }

    var channel = $("#channelStrHidden").val();
    if(channel!=""){
        $("#ChannelSel").find("option[value='"+channel+"']").attr("selected",true);
    }

})