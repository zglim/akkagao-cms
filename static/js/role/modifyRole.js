function submitModifyRoleForm() {
    var pid = $("input[name='rolepid']").val();
    var data = {
        pid: pid,
        id: $("input[name='m_roleid']").val(),
        name: $("input[name='m_rolename']").val(),
        roleurl: $("input[name='m_roleurl']").val(),
        module: $("input[name='m_module']").val(),
        action: $("input[name='m_action']").val(),
        ismenu: $("#m_roleismenu").combobox('getValue'),
        describe: $("input[name='m_roledescribe']").val()
    };

    $.post("/role/modify", data, function (result) {
        if (result == "success") {
            clearModifyRoleForm();
            refreshRoleView(pid);
            $('#modifyrole').window('close');
            $.messager.alert('操作提示', "修改成功", 'info');
        }
    });
}

function clearModifyRoleForm() {
    $('#modifyrole').form('clear');
}
