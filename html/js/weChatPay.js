$(function () {
    function onBridgeReady() {
        WeixinJSBridge.invoke(
            'getBrandWCPayRequest', {
                "appId": $("#appId").val(),     //公众号名称，由商户传入
                "timeStamp": $("#timeStamp").val(),         //时间戳，自1970年以来的秒数
                "nonceStr": $("#nonceStr").val(), //随机串
                "package": $("#package").val(),
                "signType": $("#signType").val(),         //微信签名方式：
                "paySign": $("#paySign").val() //微信签名
            },
            function (res) {
                if (res.err_msg == "get_brand_wcpay_request:ok") {
                    $("#alert-success").addClass("show");
                    $("#alert-default").removeClass("show");
                } else {
                    window.location.href = "http://bfwechat.beautyfarm.com.cn:8009";
                }
            }
        );
    }

    if (typeof WeixinJSBridge == "undefined") {
        if (document.addEventListener) {
            document.addEventListener('WeixinJSBridgeReady', onBridgeReady, false);
        } else if (document.attachEvent) {
            document.attachEvent('WeixinJSBridgeReady', onBridgeReady);
            document.attachEvent('onWeixinJSBridgeReady', onBridgeReady);
        }
    } else {
        onBridgeReady();
    }

    $("#sucessBtn").click(function () {
        var url = $("#indexUrl").val();
        if (url.indexOf("http") == -1) {
            url = "http://" + url;
        }
        window.location.href = "http://bfwechat.beautyfarm.com.cn:8009";
    });

})