// 管理主题选择
// 当前使用的主题
var themeCurrent = "default";
function theme_init() {
    $('[data-toggle="tooltip"]').tooltip();
    var t = $.zui.store.get("theme");
    if (t) {
        themeCurrent = t;
    }
    $(".theme-tile[data-theme='" + themeCurrent + "']").addClass("active");
    theme_loadCSS();

    // 监听点击事件
    $(".theme-tile").on("click", function () {
        if (!$(this).hasClass("active")) {
            // 把之前的主题不选择
            $(".theme-tile[data-theme='" + themeCurrent + "']").removeClass("active");
            $(this).addClass("active");
            themeCurrent = $(this).attr("data-theme");
            // console.log($(this));
            // 保存当前的主题名
            $.zui.store.set("theme", themeCurrent);
            theme_loadCSS();
        }
    })
}
// 加载css文件
function theme_loadCSS() {
    var url = "lib/zui/css/zui-theme";
    if (themeCurrent != "default") {
        url += "-" + themeCurrent;
    }
    url += ".css";
    $('head').append('<link href="' + url + '" rel="stylesheet" type="text/css" />');
}