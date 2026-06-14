//管理员组添加/修改窗口的公共逻辑

//tree的初始化参数(添加、修改共用)
var admusergroupsetting = {
    check: {
        enable: true
        // chkboxType: { "Y": "", "N": "" }
    },
    data: {
        simpleData: {
            enable: true
        }
    }
};

//初始化权限树
function initAdmUserGroupRoleTree(treeId, url, data) {
    $.post(url, data, function (result) {
        $.fn.zTree.init($("#" + treeId), admusergroupsetting, result);
    });
}

//收集已选中的权限节点ID，未选中时提示并返回 null
function collectAdmUserGroupRoleIds(treeId) {
    var zTree = $.fn.zTree.getZTreeObj(treeId);
    var nodes = zTree.getCheckedNodes(true);
    //判断选中的节点数，如果没有选中节点则提示操作错误
    if (nodes.length == 0) {
        $.messager.alert('操作提示', "请至少选择一个权限", 'info');
        return null;
    }
    //获取所有选中的节点ID
    var idArray = new Array(nodes.length);
    for (var i = 0; i < nodes.length; i++) {
        idArray[i] = nodes[i].id;
    }
    return idArray.join(",");
}

//提交管理员组表单，成功后关闭窗口、提示并刷新列表
function submitAdmUserGroupForm(url, data, windowId, successMsg) {
    $.post(url, data, function (result) {
        if (result == "success") {
            $('#' + windowId).window("close");
            $.messager.alert('操作提示', successMsg, 'info');
            loadAdmUserGroupDatagrid();
        } else {
            $.messager.alert('操作提示', result, 'info');
        }
    });
}
