// 360下载
var Yun360 = {
    createNew: function () {
        var self = DownBase.createNew("yun360");
        // 标题头
        self.header = "360云盘下载";
        // 加载数据
        self.loadData = function (id) {
            var file = self.am.curAccount.searchFile(id, null);
            var path = file.path;
            if (!path) {
                path = "";
            }
            C.getModule("net").send(self.className, "loadData", { account: self.am.curAccount.name, id: id, path: path }, true);
        }
        // 把服务器返回的数据转换成Fileinfo列表
        self.transFilelist = function (list) {
            var filelist = [];
            for (var i = 0; i < list.length; i++) {
                var obj = list[i];
                var file = Fileinfo.createNew(obj["id"], obj["title"], obj["size"]);
                file.path = obj["path"];
                file.isdir = obj["isdir"];
                filelist.push(file);
            }
            return filelist;
        }
        // 由file.id列表获取下载的参数列表
        self.getDownList = function (idList) {
            var curFile = self.am.curAccount.curFile;
            var list = [];
            for (var i = 0; i < idList.length; i++) {
                var id = idList[i];
                var f = self.am.curAccount.searchFile(id, curFile);
                if (f) {
                    list.push({ id: f.id, title: f.title, path: f.path });
                }
            }
            return list;
        }

        self.init();
        return self;
    }
}