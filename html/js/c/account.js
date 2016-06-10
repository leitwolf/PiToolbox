// 一个账户及文件管理
var Account = {
    createNew: function (name) {
        var self = {};
        // 账户名称
        self.name = name;
        // 根文件树 id=""
        self.rootFile = Fileinfo.createNew("", "", "", "");
        // 当前显示的文件树（实际显示的是file.children）
        self.curFile = self.rootFile;
        // 设置数据
        // @param id 所属的文件号，如果是""则是根目录
        // @param list 文件列表，只是服务器发来的数据
        self.setData = function (id, list) {
            var target = self.searchFile(id, null);
            // 没有找到对应的文件，则无效
            if (target == null) {
                console.log("no file id " + id);
                return;
            }
            // 要把之前的列表先清空
            target.children=[];
            for (var i = 0; i < list.length; i++) {
                var obj = list[i];
                var file = Fileinfo.createNew(obj["id"], obj["title"], obj["size"], obj["url"]);
                target.addChild(file);
            }
        }
        // 选择文件，显示其children
        self.selectFile = function (id) {
            self.curFile = self.searchFile(id, null);
        }
        // 通过id找到对应的文件
        // @param file 要查找的文件树 null则是从根目录开始查
        self.searchFile = function (id, file) {
            if (file == null) {
                file = self.rootFile;
            }
            if (id == file.id) {
                return file;
            }
            for (var i = 0; i < file.children.length; i++) {
                var f = file.children[i];
                // 在子列表中查找
                f1 = self.searchFile(id, f);
                if (f1 != null) {
                    return f1;
                }
            }
            return null;
        }

        return self;
    }
}