// 迅雷下载
var Xunlei = {
    createNew: function () {
        var self = Page.createNew();
        // 是否激活
        self.toggleActivate = function () {
            if (self.activated) {
                // 是否初始化过
                if (!inited) {
                    self.getAccountList();
                } else {
                    Helper.callLater(fillHtml);
                }
            }
        }
        // 获取账号列表
        self.getAccountList = function () {
            C.getModule("net").send("xunlei", "getAccountList");
        }
        // 设置账户列表，只会调用一次，当后台账户列表发生变化时，则刷新整个页面
        self.setAccountList = function (list) {
            inited = true;
            if (!list || list.length == 0) {
                $.zui.messager.show('没有设置账号', { type: 'danger', time: 3000 });
                fillHtml();
            } else {
                for (var i = 0; i < list.length; i++) {
                    var name = list[i];
                    am.addAccount(name);
                }
                // 当前账户取第一个
                self.selectAccount(list[0]);
            }
        }
        // 选择账户
        self.selectAccount = function (account) {
            var items = $("li[account]");
            for (var i = 0; i < items.length; i++) {
                var item = items[i];
                if (item.account == account) {
                    $(item).addClass("active");
                }
                else {
                    $(item).removeClass("active");
                }
            }
            am.selectAccount(account);
            if (am.curAccount != null && am.curAccount.curFile != null && !am.curAccount.curFile.loaded) {
                // 没有加载过的则加载
                self.loadData(am.curAccount.curFile.id);
            }
            else {
                fillHtml();
            }
        }
        // 加载数据
        self.loadData = function (id) {
            C.getModule("net").send("xunlei", "loadData", { account: am.curAccount.name, id: id });
        }
        // 设置数据，从服务器获得数据
        // {id,account,list}
        self.setData = function (data) {
            var id = data["id"];
            var list = data["list"];
            if (!list || !(list instanceof Array)) {
                // 返回数据错误
                $.zui.messager.show('返回列表数据错误', { type: 'danger', time: 3000 });
                return;
            }
            if (list.length == 0) {
                $.zui.messager.show('返回列表数据为空，可能是cookies已过期，请检查！', { type: 'warning', time: 3000 });
            }
            am.curAccount.setData(id, list);
            // 把当前文件树设为id
            am.curAccount.selectFile(id);
            // 渲染
            fillHtml();
        }
        // 刷新当前页面数据
        self.refresh = function () {
            self.loadData(am.curAccount.curFile.id);
        }
        // 回到主页面
        self.goBack = function () {
            am.curAccount.curFile = am.curAccount.rootFile;
            fillHtml();
        }
        // 下载
        self.download = function () {
            var values = checkjar.values;
            if (values.length == 0) {
                $.zui.messager.show('请选择要下载的文件！', { type: 'important', time: 2500 });
                return;
            }
            var curFile = am.curAccount.curFile;
            // 下载格式 [{title,url}]
            var list = [];
            for (var i = 0; i < values.length; i++) {
                var id = values[i];
                var f = am.curAccount.searchFile(id, curFile);
                if (f) {
                    list.push({ title: f.title, url: f.url });
                }
            }
            if (list.length > 0) {
                C.getModule("net").send("xunlei", "download", { account: am.curAccount.name, list: list });
            }
        }

        var checkjar = Checkjar.createNew("files");
        // 是否已初始化过
        var inited = false;
        // 账户管理
        var am = AccountsManager.createNew();
        // 初始化
        var init = function () {
            self.htmlContent = xunlei_template.substr(0);
        }
        // 根据参数填充页面
        var fillHtml = function () {
            if (!self.activated) {
                return;
            }
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
                    str += '<li ' + active + ' account="' + a.name + '"><a href="javascript:C.getModule(\'xunlei\').selectAccount(\'' + a.name + '\');">' + a.name + '</a></li>';
                }
                $("#accountList").html(str);
                $("#accountList").removeClass("hidden");
            }
            else {
                $("#accountList").addClass("hidden");
            }
            // 是否是主页面
            if (am.isRoot()) {
                $("#back").addClass("hidden");
            } else {
                $("#back").removeClass("hidden");
            }
            // 文件列表
            var str = '';
            var curList = am.getShowList();
            if (curList.length == 0) {
                str = '<tr><td colspan="3">没有文件</td></tr>';
            } else {
                // 把bt文件夹放后面
                var files = [];
                var bts = [];
                for (var i = 0; i < curList.length; i++) {
                    var file = curList[i];
                    if (file["url"] == "") {
                        bts.push(file);
                    } else {
                        files.push(file);
                    }
                }
                var newList = files.concat(bts);
                for (var i = 0; i < newList.length; i++) {
                    var file = newList[i];
                    str += '<tr>';
                    if (file["url"] == "") {
                        // bt文件夹
                        str += '<td><img src="bt.png" /></td>';
                        str += '<td><a href="javascript:C.getModule(\'xunlei\').loadData(\'' + file.id + '\');">' + file.title + '</a></td>';
                    } else {
                        str += '<td><input name="files" type="checkbox" value="' + file.id + '"></td>';
                        str += '<td>' + file.title + '</td>';
                    }
                    str += '<td>' + file.size + '</td>';
                    str += '</tr>';
                }
            }
            $("#filelist").html(str);

            Helper.activateiCheck();
            // 不选择
            checkjar.values = [];
            checkjar.init();
        }

        init();
        return self;
    }
}

// 模板
var xunlei_template = multiline(function () {/*
<h1 class="page-header">迅雷离线下载</h1>
<hr/>
<ul id="accountList" class="nav nav-secondary hidden" style="padding-bottom:10px;">
</ul>
<div style="padding-bottom:4px;">
    <a id="download" class="btn btn-primary hidden" href="javascript:C.getModule('xunlei').download();" role="button">下载</a>
    &nbsp;&nbsp;&nbsp;
    <a id="refresh" class="btn btn-primary hidden" href="javascript:C.getModule('xunlei').refresh();" role="button"><i class="icon-refresh"></i> 刷新</a>
    &nbsp;&nbsp;&nbsp;
    <a id="back" class="btn btn-primary hidden" href="javascript:C.getModule('xunlei').goBack();" role="button"><i class="icon-home"></i> 返回主页面</a>
</div>
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
