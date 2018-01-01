$(function () {
    $(".productEditBtn").bind("click", function () {
        var productId = $(this).attr("productId")

        $.get("/prodcutdetail", {"productId": productId}, function (data) {
                initFrom(data);
                $("#idp").show();
                $("#pIsactive").show();
                $('#productDetail').modal({
                    keyboard: false
                })
            }
        );
    });

    function initFrom(data) {
        var props = Object.getOwnPropertyNames(data);
        for (var i = 0; i < props.length; i++) {
            var propName = props[i];
            var propValue = data[props[i]];
            var valueFiledObj = $(".valueFiled[name='" + propName + "']");
            if (valueFiledObj != undefined) {
                var attrName = valueFiledObj.prop("tagName");
                if (attrName == "INPUT" || attrName == "TEXTAREA") {
                    valueFiledObj.val(propValue);
                }
                else if (attrName == "SPAN") {
                    {
                        valueFiledObj.html(propValue);
                    }
                }
            }
        }
    }

    $("#submitBtn").bind("click", function () {
        var fileds = $(".valueFiled");
        var o = new Object();
        for (var i = 0; i < fileds.length; i++) {
            var filed = $(fileds[i]);
            var attrName = filed.prop("tagName");
            var filedName = filed.attr("name")
            var value = "";
            if (attrName == "INPUT" || attrName == "TEXTAREA") {
                value = filed.val();
            }
            else if (attrName == "SPAN") {
                value = filed.html();
            }
            if (filedName == "Product_id" || filedName == "Isactive") {
                value = parseInt(value);
            } else if (filedName == "Price" || filedName == "Orignal_price") {
                value = parseFloat(value);
            }
            o[filedName] = value;
        }
        $.post("/prodcutdetail", {"product": JSON.stringify(o)}, function (res) {
            if (res.isSucess) {
                alert("更新成功");
                $('#productDetail').modal('hide');
            }
        });
    });

    $(".uploadImageBtn").bind("click", function () {
        var productId = parseInt($("#productDetail [name='Product_id']").html());
        var imageObj = $(this).prev("[name='image']");
        var fileNameAttr = imageObj.attr("fileName");
        if (isNaN(productId)) {
            productId = Date.parse(new Date());
        }
        var fileName = productId + "_" + fileNameAttr;
        if (imageObj.val() == "") {
            return;
        }
        var fileData = imageObj[0].files[0];
        var fd = new FormData();
        fd.append('image', fileData);

        var fd = new FormData();
        fd.append("fileName", fileName);
        fd.append("image", imageObj[0].files[0]);
        $.ajax({
            url: "/upload",
            type: "POST",
            processData: false,
            contentType: false,
            data: fd,
            success: function (res) {
                if (res.isSucess) {
                    alert("上传成功")
                    $(".valueFiled[name='" + fileNameAttr + "']").html(res.message);
                }
            }
        });
    })

    $("#addProduct").bind("click", function () {
        $("#idp").hide();
        $("#pIsactive").hide();
        var obj = {
            "Prodcut_name": "",
            "Prodcut_desc": "",
            "Prodcut_rule": "",
            "Price": "",
            "Orignal_price": "",
            "Backgroud_image": "",
            "Rule_image": "",
            "PurhchaseBtn_image": "",
            "Product_code": "",
            "MaskImage": ""
        };
        initFrom(obj);
        $('#productDetail').modal({
            keyboard: false
        })
    })
})