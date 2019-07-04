layui.use(['form', 'laydate', 'layer'], function () {
    var form = layui.form;
    var layer = layui.layer;
});

function showAllContent(o, data) {
    layer.open({
        type: 1,
        area: ['600px', '360px'],
        shadeClose: true, //点击遮罩关闭
        content: '\<\div style="padding:20px;display:block;word-break: break-all;word-wrap: break-word;line-height:22px">' + data + '\<\/div>'
    });
}

//url 提交地址
//data 表单数据
//display 提醒方式 msg/alert
//jumpType  跳转还是刷新 reload/herf
function formSubmit(url, data, display,jumpType) {
    $.ajax({
        url: url,
        data: data,
        type: "post",
        dataType: "json",
        success: function (data) {
            var messge = "网络繁忙...";
            if (data.msg) {
                messge = data.msg;
            }
            if (display === "msg"){
                layer.msg(data.msg,{icon:1,time:1000},function () {
                    if(data.code === 0){
                        if(jumpType === "reload"){
                            closeCurrentIframe();
                            window.parent.location.reload();
                        } else {
                            closeCurrentIframe();
                            window.parent.location.href = data.data
                        }
                    }else {
                        return false
                    }
                });
            } else {
                layer.alert(messge, {icon: 6, time: 5000}, function () {
                    if(data.code === 0){
                        if(jumpType === "reload"){
                            closeCurrentIframe();
                            window.parent.location.reload();
                        } else {
                            closeCurrentIframe();
                            console.log(data.data);
                            window.parent.location.href = data.data
                        }
                    }else {
                        return false
                    }
                });
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            layer.alert(messge, {icon: 6}, function () {
                // 获得frame索引
                var index = parent.layer.getFrameIndex(window.name);
                //关闭当前frame
                parent.layer.close(index);
            });
        },
        beforeSend: function () {
        },
        complete: function () {
        }
    });
}

function closeCurrentIframe() {
    if (top.location != self.location)
    {
        // 获得frame索引
        var index = parent.layer.getFrameIndex(window.name);
        console.log(index);
        //关闭当前frame
        parent.layer.close(index);
    }
}