// 一个要下载的文件或文件夹
// 根: file.parent=null||file.id=""
// 顶级的文件:file.parent.id=""
var Fileinfo = {
    createNew: function (id, title, size, url) {
        var self = {};
        // 标识
        self.id = id;
        // 文件名
        self.title = title;
        // 大小
        self.size = size;
        // 下载链接(一般url=""的话就是没有下层文件的)
        self.url = url;
        // 上一层文件，如果是root则是null
        self.parent = null;
        // 下层文件列表
        // [{file,file}]
        self.children = [];
        // 是否已加载下层文件
        self.loaded=false;
        // 添加子层文件
        self.addChild = function (file) {
            // 添加下一层文件
            self.children.push(file);
            file.parent = self;
            self.loaded=true;
        }

        return self;
    }
}