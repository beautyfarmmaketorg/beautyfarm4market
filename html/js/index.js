$(function () {
    //发送验证码
    $("#getValidCodeBtn").click(function () {
        var mobileNo = $("#phone").val();
        if(mobileNo==""){
            showErrMsg("请填写手机号码");
            return
        }
        $.get("/sendMsg",{"mobileNo":mobileNo},function(data){
            res = JSON.parse(data);
            if(res.isOk){
                showSucessMsg("短信发送成功，请查看手机获取");
            }else {
                showErrMsg("短信发送失败，请稍后重试");
            }
        });
   });
    
    $("#purchaseBtn").click(function () {
        if(!dataCheck())return;
        alert("购买成功！")
    })
    
    function dataCheck() {
        var username = $("#username").val();
        var mobileNo = $("#phone").val();
        var code = $("#code").val();
        if(username==""){
            showErrMsg("姓名不得为空");
            return false;
        }
        if(mobileNo==""){
            showErrMsg("手机号码不得为空");
            return false;
        }
        if(code==""){
            showErrMsg("验证码不得为空");
            return false;
        }
        return true;
    }

    $("input").focus(function () {
        recoverMsg();
    })
    
    function recoverMsg() {
        var warnMsg =$("#warnMsg");
        warnMsg.hide();
    }
    
    function showSucessMsg(msg) {
        var warnMsg =$("#warnMsg");
        warnMsg.removeClass();
        warnMsg.addClass("tips-sucess");
        warnMsg.addClass("tips-info");
        warnMsg.html(msg);
        warnMsg.show();
    }

    function showErrMsg(msg) {
        var warnMsg =$("#warnMsg");
        warnMsg.removeClass();
        warnMsg.addClass("tips-err");
        warnMsg.addClass("tips-info");
        warnMsg.html(msg);
        warnMsg.show();
    }
})