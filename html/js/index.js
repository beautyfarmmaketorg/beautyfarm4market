$(function () {
    //发送验证码
    $("#getValidCodeBtn").click(function () {
        var mobileNo = $("#phone").val();
        if(mobileNo==""){
            showErrMsg("请填写手机号码");
            return
        }
        $.get("/sendMsg",{"mobileNo":mobileNo},function(res){
            if(res.isSucess){
                showSucessMsg(res.message);
            }else {
                showErrMsg(res.message);
            }
        });
   });
    
    $("#purchaseBtn").click(function () {
        if(!dataCheck())return;
        var username = $("#username").val();
        var mobileNo = $("#phone").val();
        var code = $("#code").val();
        $.post("/addOrder",{"username":username,"mobileNo":mobileNo,"code":code},function(res){
            if(res.isSucess){
                var sucessAlert= $("#alert-success");
                sucessAlert.addClass("show");
            }else {
                showErrMsg(res.message);
            }
        });
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