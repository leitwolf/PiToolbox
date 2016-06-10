// 侧边栏
var Sidebar = {
    createNew: function () {
        var self = {};
        // 当前加载的页面
        self.curPage=null;
        // 加载项，加载页面
        self.loadItem = function (name) {
            var items = $("li[sidebar-item]");
            // console.log(items);
            for (i = 0; i < items.length; i++) {
                item = items[i];
                if (item.id == name) {
                    $(item).addClass("active");
                } else {
                    $(item).removeClass();
                }
            }
            // 加入页面
            if (self.curPage) {
                self.curPage.activated=false;
                self.curPage.toggleActivate();
                self.curPage=null;
            }
            var page=C.getModule(name);
            if (page) {
                page.mount();
                page.activated=true;
                page.toggleActivate();
                self.curPage=page;
            }
        }
        
        return self;
    }
}