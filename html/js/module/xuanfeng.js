// 旋风空间下载
var Xuanfeng = {
    createNew: function () {
        var self = DownBase.createNew("xuanfeng");
        // 标题头
        self.header = "旋风空间下载";
        // 加载数据
        self.loadData = function (id) {
            C.getModule("net").send(self.className, "loadData", { account: self.am.curAccount.name, id: id }, true);
        }
        // 把服务器返回的数据转换成Fileinfo列表
        self.transFilelist = function (list) {
            var filelist = [];
            for (var i = 0; i < list.length; i++) {
                var obj = list[i];
                var file = Fileinfo.createNew(obj["id"], obj["title"], obj["size"]);
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
                    list.push({ id: f.id, title: f.title });
                }
            }
            return list;
        }

        self.init();
        return self;
    }
}