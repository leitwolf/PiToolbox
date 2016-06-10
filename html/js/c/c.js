// 全局容器
var C = {
    // 模块列表，结构{name,module}
    moduleList: [],
    // 注册页面模块
    registerModule: function (name, module) {
        C.moduleList.push({ name: name, module: module });
    },
    // 获取已注册的页面模块
    getModule: function (name) {
        for (var i = 0; i < C.moduleList.length; i++) {
            var m = C.moduleList[i];
            if (m.name == name) {
                return m.module;
            }
        }
        console.log("No module: " + name + "!");
        return null;
    }
}