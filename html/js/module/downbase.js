// 云下载页面基类
var DownBase = {
    createNew: function (className) {
        var self = Page.createNew();
        // 本类名称
        self.className = className;
        // 标题头,***子类重写***
        self.header = "";
        // 是否初始化过(加载过账号列表)
        self.inited = false;
        // checkbox管理
        self.checkjar = Checkjar.createNew("files");
        // 文件树表格管理
        self.filesTable = FilesTable.createNew(self.className);
        // 账户管理
        self.am = AccountsManager.createNew();
        // 初始化
        self.init = function () {
            var str = downbase_template.substr(0);
            str = str.replace("{{header}}", self.header);
            str = str.replace(new RegExp("{{className}}", "gm"), self.className);
            str = str.replace("{{table}}", self.filesTable.getHtml());
            self.htmlContent = str;
        }
        // 是否激活
        self.toggleActivate = function () {
            if (self.activated) {
                // 是否初始化过
                if (!self.inited) {
                    self.getAccountList();
                } else {
                    Helper.callLater(self.fillHtml);
                }
                // 添加账户事件
                self.handleAddAccount();
            }
        }
        // 处理添加账户cookies文件
        // mount html到页面时调用
        self.handleAddAccount = function () {
            $('[data-toggle="tooltip"]').tooltip();
            $("#cookie_file").on("change", function () {
                var files = $(this).prop("files");
                var file = files[0];
                // 验证文件格式
                var prefix = "cookies_" + self.className;
                var suffix = ".json";
                var name = file.name;
                var valid = true;
                if (name.length < (prefix.length + suffix.length)) {
                    valid = false;
                }
                else if (name.substr(0, prefix.length) != prefix) {
                    valid = false;
                }
                else if (name.substr(name.length - suffix.length) != suffix) {
                    valid = false;
                }
                if (!valid) {
                    $.zui.messager.show('cookies文件名称格式不对', { type: 'danger', time: 3000 });
                    return;
                }
                // 读取
                var reader = new FileReader();
                reader.readAsText(file);
                reader.onload = function (e) {
                    var text = e.target.result;
                    C.getModule("net").send("cookies", "save", { filename: name, content: text, page: self.className });
                }
            })
        }
        // 获取账号列表
        self.getAccountList = function () {
            C.getModule("net").send(self.className, "getAccountList");
        }
        // 设置账户列表，只会调用一次，当后台账户列表发生变化时，则刷新整个页面
        self.setAccountList = function (list) {
            self.inited = true;
            if (!list || list.length == 0) {
                $.zui.messager.show('没有设置账号', { type: 'danger', time: 3000 });
                self.fillHtml();
            } else {
                for (var i = 0; i < list.length; i++) {
                    var name = list[i];
                    self.am.addAccount(name);
                }
                // 当前账户取第一个
                self.selectAccount(list[0]);
            }
        }
        // 选择账户
        self.selectAccount = function (accountName) {
            var am = self.am;
            am.selectAccount(accountName);
            if (am.curAccount != null && am.curAccount.curFile != null && !am.curAccount.curFile.loaded) {
                // 没有加载过的则加载
                self.loadData(am.curAccount.curFile.id);
            }
            else {
                self.fillHtml();
            }
        }
        // 进入文件夹
        self.enterDir = function (id) {
            var am = self.am;
            am.curAccount.selectFile(id);
            if (am.curAccount.curFile != null && !am.curAccount.curFile.loaded) {
                // 没有加载过的则加载
                self.loadData(am.curAccount.curFile.id);
            }
            else {
                self.fillHtml();
            }
        }
        // 排列数组
        // @param sortType "asc" "desc" ""
        self.sortTable = function (sortType) {
            // 先删除之前的点击时的tooltip
            $('.tooltip').remove();
            FilesTable.setSort(sortType);
            self.fillHtml();
        }
        // 加载数据,***由子类重写***
        self.loadData = function (id) {
        }
        // 设置数据，从服务器获得数据
        // {id,account,list:{id,title,size}}
        self.setData = function (data) {
            var id = data["id"];
            var accountName = data["account"];
            var list = data["list"];
            if (!list || !(list instanceof Array)) {
                // 返回数据错误
                $.zui.messager.show('返回列表数据错误', { type: 'danger', time: 3000 });
                return;
            }
            // 有可能数据回来的时候已经切换账户了
            var am = self.am;
            var account = am.getAccount(accountName);
            if (account) {
                // 返回数据为空且当前是根文件树，则可能是取不到数据
                if (list.length == 0 && account.curFile == account.rootFile) {
                    $.zui.messager.show('返回列表数据为空，可能是cookies已过期，请检查！', { type: 'warning', time: 3000 });
                }
                var filelist = self.transFilelist(list);
                account.setData(id, filelist);
                // 把当前文件树设为id
                account.selectFile(id);
                // console.log(am.curAccount.curFile);
                if (account == am.curAccount) {
                    // 渲染
                    self.fillHtml();
                }
            }
        }
        // 把服务器返回的数据转换成Fileinfo列表，***子类实现***
        self.transFilelist = function (list) {
        }
        // 刷新当前页面数据
        self.refresh = function () {
            self.loadData(self.am.curAccount.curFile.id);
        }
        // 回到主页面
        self.goBack = function () {
            var am = self.am;
            am.curAccount.curFile = am.curAccount.rootFile;
            self.fillHtml();
        }
        // 下载
        self.download = function () {
            var values = self.checkjar.values;
            if (values.length == 0) {
                $.zui.messager.show('请选择要下载的文件！', { type: 'important', time: 2500 });
                return;
            }
            var am = self.am;
            var curFile = am.curAccount.curFile;
            // 下载格式
            var list = self.getDownList(values);
            if (list.length > 0) {
                C.getModule("net").send(self.className, "download", { account: am.curAccount.name, list: list });
            }
        }
        // 由file.id列表获取下载的参数列表,***子类实现***
        self.getDownList = function (idList) {
        }
        // 根据参数填充页面
        self.fillHtml = function () {
            if (!self.activated) {
                return;
            }
            var am = self.am;
            if (am.curAccount == null) {
                // 没有账户
                $("#refresh").addClass("hidden");
                $("#accountList").addClass("hidden");
                $("#download").addClass("hidden");
                $("#back").addClass("hidden");
                $("#table").addClass("hidden");
                return;
            } else {
                $("#refresh").removeClass("hidden");
                $("#accountList").removeClass("hidden");
                $("#download").removeClass("hidden");
                $("#back").removeClass("hidden");
                $("#table").removeClass("hidden");
            }
            // 账号列表
            var accountList = am.accountList;
            if (accountList.length > 1) {
                var str = '<li class="nav-heading">账号列表</li>';
                for (var i = 0; i < accountList.length; i++) {
                    var a = accountList[i];
                    var active = "";
                    if (a == am.curAccount) {
                        active = 'class="active"';
                    }
                    str += '<li ' + active + ' account="' + a.name + '"><a href="javascript:C.getModule(\'' + self.className + '\').selectAccount(\'' + a.name + '\');">' + a.name + '</a></li>';
                }
                $("#accountList").html(str);
                $("#accountList").removeClass("hidden");
            }
            else {
                $("#accountList").addClass("hidden");
            }
            // 导航
            if (am.isRoot()) {
                $("#nav").addClass("hidden");
            } else {
                $("#nav").removeClass("hidden");
                var file = am.curAccount.curFile;
                var str = "";
                while (file != null) {
                    var a = "";
                    var content = file.title;
                    // 根目录
                    if (file.id == "") {
                        content = '<i class="icon icon-home"></i> 根目录';
                    }
                    if (file == am.curAccount.curFile) {
                        a = '<li class="active">' + content + '</li>';
                    }
                    else {
                        var url = "javascript:C.getModule(\'" + self.className + "\').enterDir('" + file.id + "')";
                        a = '<li><a href="' + url + '">' + content + '</a></li>';
                    }
                    str = a + str;
                    file = file.parent;
                }
                $("#nav").html(str);
            }
            // 文件列表
            self.filesTable.setData(am.getShowList());

            Helper.activateiCheck();
            $('[data-toggle="tooltip_sort"]').tooltip();
            // 不选择
            self.checkjar.values = [];
            self.checkjar.init();
        }

        return self;
    }
}

// 模板
var downbase_template = multiline(function () {/*
<h1 class="page-header">{{header}} &nbsp;&nbsp;&nbsp;
    <div class="btn btn-primary btn-file" data-toggle="tooltip" data-placement="right" data-container="body" title="格式: cookies_{{className}}_xxx.json">
        <i class="icon icon-folder-open"></i> 添加账户cookies文件
        <input id="cookie_file" type="file">
    </div>
</h1>
<ul id="accountList" class="nav nav-secondary hidden" style="padding-bottom:10px;">
</ul>
<div id="op_bar" style="padding-bottom:4px;">
    <a id="download" class="btn btn-primary hidden" href="javascript:C.getModule('{{className}}').download();" role="button">下载</a>
    &nbsp;&nbsp;&nbsp;
    <a id="refresh" class="btn btn-primary hidden" href="javascript:C.getModule('{{className}}').refresh();" role="button"><i class="icon-refresh"></i> 刷新</a>
    &nbsp;&nbsp;&nbsp;
</div>
<ol id="nav" class="breadcrumb hidden" style="margin-bottom:1px;padding:5px;">
</ol>
<div class="table-responsive">
    {{table}}
</div>
*/});
