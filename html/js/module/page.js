// 页面基类 aria2 xunlei等等
var Page = {
    createNew: function () {
        var self = {};
        // 生成的html
        self.htmlContent = "";
        // 当前是否是处理于激活状态
        self.activated = false;
        // 是否激活
        self.toggleActivate = function () {
        }
        // 加载替换到右边页面中
        self.mount = function () {
            $("#mainPage").html(self.htmlContent);
        }
        return self;
    }
}