$(function () {
    // === DataGrid 初始化 ===
    initRoleDatagrid();
    // === Tree 初始化 ===
    loadTree();
});

// ============ DataGrid 模块 ============

/**
 * 初始化权限列表 DataGrid
 */
function initRoleDatagrid() {
    $('#role_list').datagrid({
        url: 'role/gridlist',
        iconCls: 'icon-edit',
        width: 700,
        height: 'auto',
        nowrap: false,
        striped: true,
        border: true,
        collapsible: false,
        fit: true,
        remoteSort: false,
        idField: 'id',
        singleSelect: false,
        pagination: true,
        rownumbers: true,
        fitColumns: true,
        frozenColumns: [[
            { field: 'ck', checkbox: true }
        ]],
        toolbar: role_toolbar
    });
}

/**
 * 按权限ID重新加载 DataGrid 数据
 */
function loaddatagrid(id) {
    $('#role_list').datagrid('load', {
        roleid: id
    });
}

/**
 * 搜索功能
 */
var searRoleObj = {
    search: function () {
        $('#role_list').datagrid('load', {
            roleName: $('input[name="searchRoleName"]').val(),
            roleUrl: $('input[name="searchRoleUrl"]').val(),
            roleid: $("input[name='searchRolepid']").val()
        });
    }
};

/**
 * DataGrid 格式化：操作列
 */
function roleOpt(val, row, index) {
    return '<a href="#" onclick="openModifyRoleWin(' + row.id + ')">修改</a>';
}

/**
 * DataGrid 格式化：是否为菜单
 */
function roleIsMenu(val, row, index) {
    if (row.ismenu == 0) {
        return "是";
    } else if (row.ismenu == 1) {
        return "否";
    }
}

// ============ Tree 模块 ============

/**
 * zTree 配置
 */
var rolesetting = {
    data: {
        simpleData: {
            enable: true
        }
    },
    callback: {
        onClick: changeRoleList
    }
};

/**
 * 加载左侧权限树
 */
function loadTree(id) {
    var data = { id: id };
    $.post("/role/listtree", data, function (result) {
        $.fn.zTree.init($("#roletree"), rolesetting, result);
    });
}

/**
 * 点击树节点时，过滤右侧权限列表
 */
function changeRoleList(event, treeId, treeNode) {
    loaddatagrid(treeNode.id);
    $('#searchRolepid').val(treeNode.id);
}

// ============ 窗口操作模块 ============

/**
 * 打开添加权限窗口
 */
function openAddRoleWin() {
    $('#addRole').window({
        width: 400,
        height: 300,
        modal: true,
        maximizable: false,
        minimizable: false,
        collapsible: false,
        href: "/role/toadd"
    });
}

/**
 * 打开添加权限目录窗口
 */
function openAddRoleDirWin() {
    $('#addRoleDir').window({
        width: 400,
        height: 300,
        modal: true,
        maximizable: false,
        minimizable: false,
        collapsible: false,
        href: "/role/toadddir"
    });
}

/**
 * 打开修改权限窗口
 */
function openModifyRoleWin(roleid) {
    $("#roleid").attr("value", roleid);
    $('#modifyrole').window({
        width: 400,
        height: 300,
        modal: true,
        maximizable: false,
        minimizable: false,
        collapsible: false,
        href: "/role/tomodify?roleid=" + roleid
    });
}

// ============ 公共刷新方法 ============

/**
 * 刷新权限树和列表（新增/修改/删除后调用）
 */
function refreshRoleView(pid) {
    loadTree(pid);
    loaddatagrid(pid);
}

// ============ 删除操作 ============

/**
 * 删除选中的权限
 */
function deleteRole() {
    var selections = $('#role_list').datagrid('getSelections');
    if (selections.length == 0) {
        alert("请先选择要删除的列");
        return false;
    }

    if (!confirm("确定要删除选中的数据吗？")) {
        return false;
    }

    var idArray = new Array(selections.length);
    for (var i = 0; i < selections.length; i++) {
        idArray[i] = selections[i].id;
    }
    var ids = idArray.join(",");

    var pid = $("input[name='searchRolepid']").val();
    $.post("/role/deleterole", { ids: ids }, function (result) {
        refreshRoleView(pid);
        selections.length = 0;
        if (result == "success") {
            $.messager.alert('操作提示', "删除成功", 'info');
        } else {
            $.messager.alert('操作提示', result, 'warning');
        }
    });
}
