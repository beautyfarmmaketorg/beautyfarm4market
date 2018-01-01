$(function () {
    $(".productEditBtn").bind("click", function () {
        var productId = $(this).attr("productId")

        $.get("/prodcutdetail", {"productId": productId}, function (data) {
                var formObj = $("#productDetail form");
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
                $('#productDetail').modal({
                    keyboard: false
                })
            }
        );
    });

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
            }
        });
    });

    $(".uploadImageBtn").bind("click", function () {
        var productId = parseInt($("#productDetail [name='Product_id']").html());
        var imageObj = $(this).prev("[name='image']");
        var fileNameAttr = imageObj.attr("fileName");
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
})