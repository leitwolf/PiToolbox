// 多账户管理
var AccountsManager = {
    createNew: function () {
        var self = {};
        // 账户列表
        self.accountList = [];
        // 当前使用的账户 account
        self.curAccount = null;
        // 添加一个账户
        self.addAccount = function (name) {
            var account = Account.createNew(name);
            self.accountList.push(account);
        }
        // 选择账户
        self.selectAccount = function (name) {
            self.curAccount = self.getAccount(name);
            if (self.curAccount == null) {
                console.log("no account name " + name);
            }
        }
        // 获取指定的账户
        self.getAccount = function (name) {
            for (var i = 0; i < self.accountList.length; i++) {
                var account = self.accountList[i];
                if (account.name == name) {
                    return account;
                }
            }
            return null;
        }
        // 是否是根目录，没有当前账户或没有当前文件树的都算
        self.isRoot = function () {
            return self.curAccount == null || self.curAccount.curFile == null || self.curAccount.curFile.id == "";
        }
        // 获取当前要显示的文件列表
        self.getShowList = function () {
            var list = [];
            if (self.curAccount && self.curAccount.curFile) {
                list = self.curAccount.curFile.children;
            }
            return list;
        }

        return self;
    }
}