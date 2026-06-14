
//初始化tree
$(document).ready(loadTree());
function loadTree() {
    var admgroupuserid = $("input[name='admgroupuserid']").val()
    url = "/admusergroup/loadtreechecked"
    var data = {
        admgroupuserid: admgroupuserid,
    };
    $.post(url, data, function (result) {
        $.fn.zTree.init($("#modifyadmgrouproletree"), admusergroupsetting, result);
    });
}

/**
 * 修改管理员组
 */
function submitModifyAmdUserGroupForm() {
    var result = getCheckedNodeIds("modifyadmgrouproletree");
    if (!validateCheckedNodes(result.nodes)) {
        return false;
    }

    url = "/admusergroup/modifyadmusergroup"

    var data = {
        ids: result.ids,
        id: $("input[name='admgroupuserid']").val(),
        groupname: $("input[name='ag_m_name']").val(),
        describe: $("input[name='ag_m_describe']").val()
    };

    $.post(url, data, function (result) {
        handleGroupSubmitResult(result, '#modifyadmusergroup', "修改成功");
    });
}

function clearModifyAmdUserGroupForm() {
    $('#modifyadmusergroup').form('clear');
}
