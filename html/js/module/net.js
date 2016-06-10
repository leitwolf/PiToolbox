// 网络通信模块
var Net = {
    createNew: function () {
        var slef = {};
        // 发送地址
        var url = "action";
        // 发送
        self.send = function (module, action, data) {
            if (data == undefined) {
                data = null;
            }
            var obj = { module: module, action: action, data: data };
            $.ajax({
                type: "POST",
                url: url,
                dataType: "json",
                cache: false,
                data: JSON.stringify(obj),
                success: function (result) {
                    console.log(result);
                    handler(result);
                },
                error: function (req, textStatus, errorThrown) {
                    $.zui.messager.show('服务器错误: ' + textStatus + errorThrown, { type: 'danger', time: 3000 });
                }
            });
        }
        // 返回数据的处理，分发处理
        var handler = function (data1) {
            var module = data1.module;
            var action = data1.action;
            var data = data1.data;
            var err = data1.err;
            if (err != "") {
                $.zui.messager.show('Error: ' + err, { type: 'danger', time: 2000 });
                return;
            }
            if (module == "aria2") {
                var aria2 = C.getModule("aria2");
                if (action == "getVersion") {
                    aria2.setVersion(data);
                }
                else if (action == "getStat") {
                    aria2.setData(data);
                }
                else {
                    aria2.refresh();
                }
            } else if (module == "xunlei") {
                var xunlei = C.getModule("xunlei");
                if (action == "getAccountList") {
                    xunlei.setAccountList(data);
                } else if (action == "loadData") {
                    xunlei.setData(data);
                } else if (action == "download") {
                    $.zui.messager.show('添加下载成功', { type: 'success', time: 2000 });
                }
            }
        }

        return self;
    }
}