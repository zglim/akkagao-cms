
//管理员组管理 - 公共工具函数

//设置tree的初始化参数
var admusergroupsetting = {
    check: {
        enable: true
    },
    data: {
        simpleData: {
            enable: true
        }
    }
};

/**
 * 获取zTree中已选节点的ID列表
 * @param {string} treeId zTree的DOM元素ID
 * @returns {object} { nodes: 选中节点数组, ids: 逗号分隔的ID字符串 }
 */
function getCheckedNodeIds(treeId) {
    var zTree = $.fn.zTree.getZTreeObj(treeId);
    var nodes = zTree.getCheckedNodes(true);
    var idArray = new Array(nodes.length);
    for (var i = 0; i < nodes.length; i++) {
        idArray[i] = nodes[i].id;
    }
    return {
        nodes: nodes,
        ids: idArray.join(",")
    };
}

/**
 * 校验是否选择了权限节点
 * @param {array} nodes 已选节点数组
 * @returns {boolean} 是否通过校验
 */
function validateCheckedNodes(nodes) {
    if (nodes.length == 0) {
        $.messager.alert('操作提示', "请至少选择一个权限", 'info');
        return false;
    }
    return true;
}

/**
 * 提交管理员组表单后的统一处理
 * @param {string} result 服务器返回结果
 * @param {string} windowId 窗口元素ID
 * @param {string} successMsg 成功提示文案
 */
function handleGroupSubmitResult(result, windowId, successMsg) {
    if (result == "success") {
        $(windowId).window("close");
        $.messager.alert('操作提示', successMsg, 'info');
        loadAdmUserGroupDatagrid();
    } else {
        $.messager.alert('操作提示', result, 'info');
    }
}
