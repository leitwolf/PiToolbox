// 导航栏
var Nav = {
    createNew: function () {
        var self = {};
        // 当前加载的页面
        self.curPage = null;
        // 加载项，加载页面
        self.loadItem = function (name) {
            var items = $("li[nav-item]");
            // console.log(items);
            for (i = 0; i < items.length; i++) {
                item = items[i];
                if (item.id == name) {
                    $(item).addClass("active");
                } else {
                    $(item).removeClass();
                }
            }
            if (name == "aria2") {
                $("#nav_list").removeClass();
            }
            else {
                $("#nav_list").addClass("active");
            }
            // 加入页面
            if (self.curPage) {
                self.curPage.activated = false;
                self.curPage.toggleActivate();
                self.curPage = null;
            }
            var page = C.getModule(name);
            if (page) {
                page.mount();
                page.activated = true;
                page.toggleActivate();
                self.curPage = page;
            }
            // 下载界面钉住下载，刷新按钮
            // if (name != "aria2") {
            //     setTimeout(function() {
            //         $("#op_bar").stickUp({marginTop: '40px'});
            //         console.log("pin");
            //     }, 500);
            // }
        }

        return self;
    }
}