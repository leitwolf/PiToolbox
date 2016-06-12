// 文件树表格类
var FilesTable = {
    // @param string enterDirFunc 进入文件夹时要调用的函数，此函数接受一参数:file.id
    createNew: function (enterDirFunc) {
        var self = {};
        // 获取初始化时的表格html
        self.getHtml = function () {
            return filestable_template;
        }
        // 设置数据，更新页面
        // @param filelist 要显示的Fileinfo列表
        self.setData = function (filelist) {
            var str = "";
            if (filelist.length == 0) {
                str = '<tr><td colspan="3">没有文件</td></tr>';
            } else {
                // 把bt文件夹放后面
                var files = [];
                var bts = [];
                for (var i = 0; i < filelist.length; i++) {
                    var file = filelist[i];
                    if (file.isdir) {
                        bts.push(file);
                    } else {
                        files.push(file);
                    }
                }
                var newList = files.concat(bts);
                for (var i = 0; i < newList.length; i++) {
                    var file = newList[i];
                    str += '<tr>';
                    if (file.isdir) {
                        // 文件夹
                        str += '<td><i class="icon icon-folder-close-alt icon-2x"></i></td>';
                        str += '<td><a href="javascript:' + enterDirFunc + '(\'' + file.id + '\');">' + file.title + '</a></td>';
                    } else {
                        str += '<td><input name="files" type="checkbox" value="' + file.id + '"></td>';
                        str += '<td>' + file.title + '</td>';
                    }
                    str += '<td>' + file.size + '</td>';
                    str += '</tr>';
                }
            }
            $("#filelist").html(str);
        }

        return self;
    }
}
// 模板
var filestable_template = multiline(function () {/*
<table id="table" class="table table-hover hidden">
    <thead>
        <tr>
            <th style="width:5%"><input id="files" type="checkbox"></th>
            <th style="width:85%">文件名</th>
            <th style="width:10%">文件大小</th>
        </tr>
    </thead>
    <tbody id="filelist">
    </tbody>
</table>
*/});
