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
