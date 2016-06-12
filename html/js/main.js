// 初始化
function init() {
    theme_init();
    var sidebar = Sidebar.createNew();
    C.registerModule("sidebar", sidebar);
    var net = Net.createNew();
    C.registerModule("net", net);
    var aria2 = Aria2.createNew();
    C.registerModule("aria2", aria2);
    var xunlei = Xunlei.createNew();
    C.registerModule("xunlei", xunlei);
    var yun360 = Yun360.createNew();
    C.registerModule("yun360", yun360);

    sidebar.loadItem("aria2");
}