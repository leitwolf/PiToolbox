// 当前表格排列方式 asc desc 
var tableSort = "";
// 文件树表格类
var FilesTable = {
    // 初始化
    init: function () {
        var s = $.zui.store.get("sort");
        if (s) {
            tableSort = s;
        }
    },
    // 设置排序
    setSort: function (sort) {
        tableSort = sort;
        $.zui.store.set("sort", tableSort);
    },
    // 排序列表，不改变list
    sortList: function (list) {
        if (tableSort == "asc") {
            list.sort(function (a, b) {
                return a["title"].localeCompare(b["title"]);
            });
        }
        else if (tableSort == "desc") {
            list.sort(function (a, b) {
                return b["title"].localeCompare(a["title"]);
            });
        }
    },
    // @param string className 当前page的名称
    createNew: function (className) {
        var self = {};
        // 获取初始化时的表格html
        self.getHtml = function () {
            return filestable_template;
        }
        // 设置数据，更新页面
        // @param filelist 要显示的Fileinfo列表
        self.setData = function (filelist) {
            var str = "";
            var sortStr = "";
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
                FilesTable.sortList(files);
                FilesTable.sortList(bts);
                var newList = files.concat(bts);
                for (var i = 0; i < newList.length; i++) {
                    var file = newList[i];
                    str += '<tr>';
                    if (file.isdir) {
                        // 文件夹
                        var enterDirFunc = "C.getModule(\'" + className + "\').enterDir"
                        str += '<td><i class="icon icon-folder-close-alt icon-2x"></i></td>';
                        str += '<td><a href="javascript:' + enterDirFunc + '(\'' + file.id + '\');">' + file.title + '</a></td>';
                    } else {
                        str += '<td><input name="files" type="checkbox" value="' + file.id + '"></td>';
                        str += '<td>' + file.title + '</td>';
                    }
                    str += '<td>' + file.size + '</td>';
                    str += '</tr>';
                }
                // 排序图标
                var sortFunc = "C.getModule(\'" + className + "\').sortTable";
                var normal = '&nbsp;&nbsp;&nbsp;<a href="javascript:' + sortFunc + '(\'\');" data-toggle="tooltip_sort" data-placement="top" data-trigger="hover" data-container="body" title="原始排列"><i class="icon icon-align-justify"></i></a>';
                var asc = '&nbsp;&nbsp;&nbsp;<a href="javascript:' + sortFunc + '(\'asc\');" data-toggle="tooltip_sort" data-placement="top" data-trigger="hover" data-container="body" title="升序排列"><i class="icon icon-circle-arrow-up"></i></a>';
                var desc = '&nbsp;&nbsp;&nbsp;<a href="javascript:' + sortFunc + '(\'desc\');" data-toggle="tooltip_sort" data-placement="top" data-trigger="hover" data-container="body" title="降序排列"><i class="icon icon-circle-arrow-down"></i></a>';
                if (tableSort == "asc") {
                    asc = '&nbsp;&nbsp;&nbsp;<i class="icon icon-circle-arrow-up"></i>';
                }
                else if (tableSort == "desc") {
                    desc = '&nbsp;&nbsp;&nbsp;<i class="icon icon-circle-arrow-down"></i>';
                }
                else {
                    normal = '&nbsp;&nbsp;&nbsp;<i class="icon icon-align-justify"></i>';
                }
                var itemStr='&nbsp;&nbsp;&nbsp;<span class="label label-success">'+newList.length+"</span>";
                sortStr = normal + asc + desc+itemStr;
            }
            $("#filelist").html(str);
            $("#table_sort").html(sortStr);
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
            <th style="width:85%">文件名<div id="table_sort" style="display:inline-block;"></div></th>
            <th style="width:10%">文件大小</th>
        </tr>
    </thead>
    <tbody id="filelist">
    </tbody>
</table>
*/});
