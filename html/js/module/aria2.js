// Aria2下载管理页面
var Aria2 = {
    createNew: function () {
        var self = Page.createNew();
        // 是否激活
        self.toggleActivate = function () {
            if (self.activated) {
                // 是否获取过版本号
                if (aria2Version != "") {
                    self.setVersion(aria2Version)
                } else {
                    C.getModule("net").send("aria2", "getVersion");
                }
            }
            else {
                clearTick();
                document.title = "PiToolbox";
            }
        }
        // 设置版本号
        self.setVersion = function (version) {
            aria2Version = version;
            if (version == "") {
                $("#version").html("");
                $.zui.messager.show('无法连接Aria2!', { type: 'danger', time: 2000 });
            }
            else {
                $("#version").html("Ver. " + version);
                requestData();
            }
        };
        // 设置数据
        // {speed:xxx,activeTasks:[],waitingTasks:[],stopedTasks=[]}
        self.setData = function (data) {
            // console.log("tick"+Math.floor(Math.random()*100));
            $("#speed").html('<i class="icon-download-alt"></i> ' + data.speed);
            var tasks;
            if (data.activeTasks && data.waitingTasks) {
                tasks = data.activeTasks.concat(data.waitingTasks);
            } else if (data.activeTasks) {
                tasks = data.activeTasks;
            } else {
                tasks = data.waitingTasks;
            }
            var html = activeTask.getHtml(tasks);
            $("#activeTasks").html(html);
            var html = stopedTask.getHtml(data.stopedTasks);
            $("#stopedTasks").html(html);
            Helper.activateiCheck();
            activeTask.checkjar.init();
            stopedTask.checkjar.init();

            if (self.activated) {
                // 修改页面标题
                document.title = data.speed;
            }
            doingRequestData = false;
            // 继续请求数据
            tickData();
        };
        // 刷新，当完成操作后执行
        self.refresh = function () {
            clearTick();
            requestData();
        }
        // 开始
        self.start = function () {
            if (!checkSelectValues("active")) {
                return;
            }
            console.log("values", activeTask.values);
            C.getModule("net").send("aria2", "start", activeTask.checkjar.values);
        };
        // 暂停
        self.pause = function () {
            if (!checkSelectValues("active")) {
                return;
            }
            C.getModule("net").send("aria2", "pause", activeTask.checkjar.values);
        };
        // 删除
        self.remove = function () {
            if (!checkSelectValues("active")) {
                return;
            }
            C.getModule("net").send("aria2", "remove", activeTask.checkjar.values);
        };
        // 开始所有
        self.startAll = function () {
            C.getModule("net").send("aria2", "startAll");
        };
        // 暂停所有
        self.pauseAll = function () {
            C.getModule("net").send("aria2", "pauseAll");
        };
        // 删除所有
        self.removeAll = function () {
            C.getModule("net").send("aria2", "remove", activeTask.checkjar.getAllValues());
        };
        // 删除已停止的
        self.removeStoped = function () {
            if (!checkSelectValues("stoped")) {
                return;
            }
            C.getModule("net").send("aria2", "removeStoped", stopedTask.checkjar.values);
        };
        // 删除所有停止的任务
        self.removeAllStoped = function () {
            C.getModule("net").send("aria2", "removeAllStoped");
        };
        // 版本号
        var aria2Version = "";
        // 任务类
        var activeTask = Aria2ActiveTask.createNew();
        var stopedTask = Aria2StopedTask.createNew();
        // 初始化
        var init = function () {
            var str = aria2_mainTemplate.substr(0);
            str = str.replace("{{activeTasks}}", activeTask.getHtml(null));
            str = str.replace("{{stopedTasks}}", stopedTask.getHtml(null));
            self.htmlContent = str;
            Helper.callLater(Helper.activateiCheck);
        }
        // 定时调用数据
        var tickTimeHandler = 0;
        var tickData = function () {
            tickTimeHandler = setTimeout(function () {
                requestData();
            }, 2000);
        }
        var clearTick = function () {
            if (tickTimeHandler) {
                clearTimeout(tickTimeHandler);
                tickTimeHandler = 0;
            }
        }
        // 请求数据
        // 是否在请求数据中，一次只能请求一次
        var doingRequestData = false;
        var requestData = function () {
            if (!self.activated || doingRequestData) {
                return;
            }
            doingRequestData = true;
            C.getModule("net").send("aria2", "getStat");
        }
        // 先检测是否已选择了任务
        // @param taskType [active stoped]
        var checkSelectValues = function (taskType) {
            var values;
            if (taskType == "active") {
                values = activeTask.checkjar.values;
            } else {
                values = stopedTask.checkjar.values;
            }
            if (values.length == 0) {
                // 没有选择项
                $.zui.messager.show('请选择任务！', { type: 'important', time: 2500 });
                return false;
            }
            return true;
        };

        init();
        return self;
    }
}

// 进行中的任务
var Aria2ActiveTask = {
    createNew: function () {
        var self = {};
        // checkbox管理
        var checkName = "active_tasks";
        self.checkjar = Checkjar.createNew(checkName);
        // 根据数据得出表结构
        // @param data [{gid,filename,size,speed,connections,progress,status}]
        self.getHtml = function (data) {
            var html = aria2_activeTemplate.substr(0);
            html = html.replace("{{checkName}}", checkName);
            // 项
            var body = "";
            if (!data || data.length == 0) {
                // 空表格
                body += '<tr>';
                body += '<td colspan="7">没有进行中的任务</td>';
                body += '</tr>';
            } else {
                for (var i = 0; i < data.length; i++) {
                    var item = data[i];
                    var size = item["size"];
                    var completedLength = item["completedLength"];
                    var speed = item["speed"];
                    // 剩余时间
                    var leftTimeStr = "---";
                    if (speed > 0) {
                        var leftTime = Math.floor((size - completedLength) / speed);
                        leftTimeStr = Helper.getReadableTime(leftTime);
                    }
                    var sizeStr = Helper.getReadableSize(size);
                    body += '<tr>';
                    body += '<td><input name="' + checkName + '" type="checkbox" value="' + item["gid"] + '"></td>';
                    body += '<td>' + item["filename"] + '</td>';
                    body += '<td>' + sizeStr + 'B</td>';
                    body += '<td>' + Helper.getReadableSize(item["speed"]) + 'B/s</td>';
                    body += '<td>' + leftTimeStr + '</td>';
                    body += '<td>' + Helper.getProgressBar(item["progress"]) + '</td>';
                    body += '<td>' + Helper.getStatusLabel(item["status"]) + '</td>';
                    body += '</tr>';
                }
            }
            html = html.replace("{{body}}", body);
            return html;
        }

        return self;
    }
}

// 已停止的任务
var Aria2StopedTask = {
    createNew: function () {
        var self = {};
        // checkbox管理
        var checkName = "stoped_tasks";
        self.checkjar = Checkjar.createNew(checkName);
        // 根据数据得出表结构
        // @param data [{gid,filename,size,progress,status}]
        self.getHtml = function (data) {
            var html = aria2_stopedTemplate.substr(0);
            html = html.replace("{{checkName}}", checkName);
            // 项
            var body = "";
            if (!data || data.length == 0) {
                // 空表格
                body += '<tr>';
                body += '<td colspan="5">没有已停止的任务</td>';
                body += '</tr>';
            } else {
                // 对列表进行反转，后加入的放在前面
                data.reverse();
                for (var i = 0; i < data.length; i++) {
                    var item = data[i];
                    body += '<tr>';
                    body += '<td><input name="' + checkName + '" type="checkbox" value="' + item["gid"] + '"></td>';
                    body += '<td>' + item["filename"] + '</td>';
                    body += '<td>' + Helper.getReadableSize(item["size"]) + 'B</td>';
                    body += '<td>' + Helper.getProgressBar(item["progress"]) + '</td>';
                    body += '<td>' + Helper.getStatusLabel(item["status"]) + '</td>';
                    body += '</tr>';
                }
            }
            html = html.replace("{{body}}", body);
            return html;
        }

        return self;
    }
}

// 主模板
var aria2_mainTemplate = multiline(function () {/*
<h1 class="page-header">Aria2下载管理 &nbsp;&nbsp;&nbsp;
<a class="btn btn-primary" href="#" role="button"><i class="icon-cog"></i> 设置</a>
&nbsp;&nbsp;&nbsp;<small id="version"></small>
&nbsp;&nbsp;&nbsp;<small id="speed"></small>
</h1>
<hr>
<div class="panel panel-primary">
    <div class="panel-heading">进行中的任务</div>
    <div class="panel-body">
        <div style="padding-bottom:4px;">
            <div class="btn-group" role="group">
                <a class="btn btn-primary" href="javascript:C.getModule('aria2').start();" role="button">开始</a>
                <a class="btn btn-primary" href="javascript:C.getModule('aria2').pause();" role="button">暂停</a>
                <a class="btn btn-danger" href="javascript:C.getModule('aria2').remove();" role="button">删除</a>
            </div>
            &nbsp;&nbsp;&nbsp;&nbsp;
            <a class="btn btn-primary" href="javascript:C.getModule('aria2').startAll();" role="button">开始所有</a>
            <a class="btn btn-primary" href="javascript:C.getModule('aria2').pauseAll();" role="button">暂停所有</a>
            <a class="btn btn-danger" href="javascript:C.getModule('aria2').remove();" role="button">删除所有</a>
        </div>
        <div id="activeTasks" class="table-responsive">
            {{activeTasks}}
        </div>
    </div>
</div>
<div class="panel panel-primary">
    <div class="panel-heading">已停止的任务</div>
    <div class="panel-body">
        <div style="padding-bottom:4px;">
            <a class="btn btn-danger" href="javascript:C.getModule('aria2').removeStoped();" role="button">删除</a>
            &nbsp;&nbsp;&nbsp;&nbsp;
            <a class="btn btn-danger" href="javascript:C.getModule('aria2').removeAllStoped();" role="button">删除所有</a>
        </div>
        <div id="stopedTasks" class="table-responsive">
            {{stopedTasks}}
        </div>
    </div>
</div>
*/});

// 进行中任务模板
var aria2_activeTemplate = multiline(function () {/*
<table class="table table-striped">
    <thead>
        <tr>
            <th width="3%">
                <input id="{{checkName}}" type="checkbox">
            </th>
            <th width="47%">文件名</th>
            <th width="8%">文件大小</th>
            <th width="5%">速度</th>
            <th width="10%">剩余时间</th>
            <th width="20%">进度</th>
            <th width="5%">状态</th>
        </tr>
    </thead>
    <tbody>
        {{body}}
    </tbody>
</table>
*/});

// 停止任务模板
var aria2_stopedTemplate = multiline(function () {/*
<table class="table table-striped">
    <thead>
        <tr>
            <th width="3%">
                <input id="{{checkName}}" type="checkbox">
            </th>
            <th width="50%">文件名</th>
            <th width="10%">文件大小</th>
            <th width="30%">进度</th>
            <th width="7%">状态</th>
        </tr>
    </thead>
    <tbody>
        {{body}}
    </tbody>
</table>
*/});
