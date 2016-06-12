// 一个要下载的文件或文件夹
// 根: file.parent=null||file.id=""
// 顶级的列表文件:file.parent.id=""
var Fileinfo = {
    createNew: function (id, title, size) {
        var self = {};
        // 标识
        self.id = id;
        // 文件名
        self.title = title;
        // 大小
        self.size = size;
        // 下载链接(有些并没有，要实时的)
        self.url = "";
        // 是否是文件夹
        self.isdir = false;
        // 上一层文件，如果是root则是null
        self.parent = null;
        // 下层文件列表
        // [{file,file}]
        self.children = [];
        // 是否已加载下层文件(当isdir=true时使用)
        self.loaded = false;
        // 添加子层文件
        self.addChild = function (file) {
            // 添加下一层文件
            self.children.push(file);
            file.parent = self;
            self.loaded = true;
        }

        return self;
    }
}