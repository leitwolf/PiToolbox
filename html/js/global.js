// 帮助类
var Helper = {
    // 延迟（下一帧）调用
    callLater: function (method, ...params) {
        setTimeout(function () {
            method.apply(null, params);
        }, 0);
    },
    // 挂载，替换target处的html内容
    mount: function (target, htmlContent) {
        $("#" + target).html(htmlContent);
    },
    // 激活iCheck组件，一般加载完页面之后
    activateiCheck: function () {
        $('input').iCheck({
            checkboxClass: 'icheckbox_square-blue',
            radioClass: 'iradio_square-blue',
            increaseArea: '10%'
        });
    },
    // 获取进度条
    // @param progress 0-100
    getProgressBar: function (progress) {
        var html = '<div class="progress">';
        html += '<div class="progress-bar" role="progressbar" aria-valuenow="' + progress + '" aria-valuemin="0" aria-valuemax="100" style="width: ' + progress + '%;">';
        html += progress + '%';
        html += '</div>';
        // html+='<div style="width:100%;position:relative;">'+progress + '%</div>';
        html += '</div>';
        return html;
    },
    // 获取状态标签
    // @param status [active,waiting,paused,completed,removed,error]
    getStatusLabel: function (status) {
        var html = '';
        if (status == "active") {
            html = '<span class="label label-info">下载中</span>';
        } else if (status == "waiting") {
            html = '<span class="label label-warning">等待</span>';
        } else if (status == "paused") {
            html = '<span class="label label-default">暂停</span>';
        } else if (status == "complete") {
            html = '<span class="label label-success">完成</span>';
        } else if (status == "removed") {
            html = '<span class="label label-danger">已删除</span>';
        } else if (status == "error") {
            html = '<span class="label label-danger">错误</span>';
        }
        return html;
    },
    // 把纯数字大小转成G M K的样式
    _G: 1024 * 1024 * 1024,
    _M: 1024 * 1024,
    _K: 1024,
    getReadableSize: function (size) {
        size = parseInt(size);
        var readable = "";
        if (size >= Helper._G) {
            var r = size / Helper._G;
            readable = r.toFixed(2) + "G";
        } else if (size >= Helper._M) {
            var r = size / Helper._M;
            readable = r.toFixed(2) + "M";
        } else if (size >= Helper._K) {
            var r = size / Helper._K;
            readable = r.toFixed(2) + "K";
        } else {
            readable = size + "";
        }
        return readable;
    },
    // 把纯数字时间少转成时，分，秒的形式
    _HOUR: 3600,
    _MINUTE: 60,
    getReadableTime: function (seconds) {
        seconds = parseInt(seconds);
        var readable = "";
        if (seconds >= Helper._HOUR) {
            var h = Math.floor(seconds / Helper._HOUR);
            readable += h + "h";
            seconds -= h * Helper._HOUR;
        }
        if (seconds >= Helper._MINUTE) {
            var m = Math.floor(seconds / Helper._MINUTE);
            readable += m + "m";
            seconds -= m * Helper._MINUTE;
        }
        readable += seconds + "s";
        return readable;
    },
}

// js中多行文本实现
// https://github.com/sindresorhus/multiline
var reCommentContents = /\/\*!?(?:\@preserve)?[ \t]*(?:\r\n|\n)([\s\S]*?)(?:\r\n|\n)[ \t]*\*\//;
var multiline = function (fn) {
    if (typeof fn !== 'function') {
        throw new TypeError('Expected a function');
    }
    var match = reCommentContents.exec(fn.toString());
    if (!match) {
        throw new TypeError('Multiline comment missing.');
    }
    return match[1];
};
