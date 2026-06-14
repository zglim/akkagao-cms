//初始化tree
$(document).ready(loadTree());
function loadTree() {
    var admgroupuserid = $("input[name='admgroupuserid']").val();
    initAdmUserGroupRoleTree("modifyadmgrouproletree", "/admusergroup/loadtreechecked", {
        admgroupuserid: admgroupuserid
    });
}

/**
 * 修改管理员组
 */
function submitModifyAmdUserGroupForm() {
    var ids = collectAdmUserGroupRoleIds("modifyadmgrouproletree");
    if (ids == null) {
        return false;
    }

    var data = {
        ids: ids,
        id: $("input[name='admgroupuserid']").val(),
        groupname: $("input[name='ag_m_name']").val(),
        describe: $("input[name='ag_m_describe']").val()
    };

    submitAdmUserGroupForm("/admusergroup/modifyadmusergroup", data, "modifyadmusergroup", "修改成功");
}

function clearModifyAmdUserGroupForm() {
    $('#modifyadmusergroup').form('clear');
}
