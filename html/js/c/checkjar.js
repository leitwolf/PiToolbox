// 
// 管理checkbox
// 
var Checkjar = {
    createNew: function (name) {
        var self = {};
        // name 所有子checkbox的name,主checkbox的id
        self.name = name;
        // 已选择的值列表，要保存下来
        self.values = [];
        // 初始化，刷新一次列表的时候执行
        // 操作内容：
        // 1、把之前选择的value对应checkbox选上
        // 2、给新的checkbox加上事件监听
        self.init = function () {
            var list = $("input[name='" + name + "']");
            for (var i = 0; i < list.length; i++) {
                var check = list[i];
                var v = $(check).val();
                $(check).iCheck('uncheck');
                for (var j = 0; j < self.values.length; j++) {
                    if (v == self.values[j]) {
                        $(check).iCheck('check');
                    }
                }
                // 监听事件
                $(check).on("ifClicked", function () {
                    Helper.callLater(self.update);
                });
            }
            // 总开关的事件
            $("#" + self.name).on("ifClicked", function (event) {
                Helper.callLater(function () {
                    var checked = event.target.checked;
                    // console.log(checked);
                    var values = [];
                    for (var i = 0; i < list.length; i++) {
                        var check = list[i];
                        var v = $(check).val();
                        if (checked) {
                            $(check).iCheck('check');
                            values.push(v);
                        } else {
                            $(check).iCheck('uncheck');
                        }
                    }
                    self.values = values;
                })
            });

            self.update();
        }
        // 更新
        // 1、选中值列表 2、主checkbox是否选中
        self.update = function () {
            var list = $("input:checked[name='" + name + "']");
            var values = [];
            for (var i = 0; i < list.length; i++) {
                var check = list[i];
                var v = $(check).val();
                values.push(v);
            }
            self.values = values;
            // 已选中的项目数和全部的项目数一样，则说明已全选
            var list2 = $("input[name='" + name + "']");
            if (list.length == list2.length && list.length > 0) {
                $("#" + self.name).iCheck('check');
            } else {
                $("#" + self.name).iCheck('uncheck');
            }
        }
        // 获取所有的值
        self.getAllValues = function () {
            var list = $("input[name='" + name + "']");
            var values = [];
            for (var i = 0; i < list.length; i++) {
                var check = list[i];
                var v = $(check).val();
                values.push(v);
            }
            return values;
        }
        // 记得返回
        return self;
    }
}