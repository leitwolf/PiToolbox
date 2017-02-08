// 初始化
function init() {
    theme_init();
    FilesTable.init();
    var nav = Nav.createNew();
    C.registerModule("nav", nav);
    var net = Net.createNew();
    C.registerModule("net", net);
    var aria2 = Aria2.createNew();
    C.registerModule("aria2", aria2);
    var xunlei = Xunlei.createNew();
    C.registerModule("xunlei", xunlei);
    var yun360 = Yun360.createNew();
    C.registerModule("yun360", yun360);
    var xuanfeng = Xuanfeng.createNew();
    C.registerModule("xuanfeng", xuanfeng);

    var q = document.location.search;
    if (q == "?xunlei") {
        nav.loadItem("xunlei");
    }
    // else if (q == "?yun360") {
    //     nav.loadItem("yun360");
    // }
    else if (q == "?xuanfeng") {
        nav.loadItem("xuanfeng");
    }
    else {
        nav.loadItem("aria2");
    }
}