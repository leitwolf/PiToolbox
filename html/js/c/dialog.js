// 对话框
var Dialog = {
    // 当前动作
    curAction: null,
    // 当前对话框
    curDialog: null,
    // 删除确认
    // @param string action 动作
    deleteConfirm: function (action) {
        curAction = action;
        var str = dialog_template.substr(0);
        str = str.replace("{{content}}", "确定要删除吗？");
        str = str.replace("{{actionText}}", "删除");
        str = str.replace("{{action}}", "Dialog.doAction();");
        curDialog = new $.zui.ModalTrigger({ custom: str, showHeader: false, size: 'sm' });
        curDialog.show();
    },
    // 执行当前动作
    doAction: function () {
        if (curDialog) {
            curDialog.close();
            curDialog = null;
            if (curAction) {
                curAction.call(null);
            }
        }
    }
}

// 对话框模板
var dialog_template = multiline(function () {/*
<div class="modal-content">
    <div class="modal-body">
        <h3>{{content}}</h3>
    </div>
    <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
        <a class="btn btn-danger" href="javascript:{{action}}" role="button">{{actionText}}</a>
    </div>
</div>
*/});
