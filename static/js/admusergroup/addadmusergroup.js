//初始化左边tree
$(document).ready(loadTree());
function loadTree() {
    initAdmUserGroupRoleTree("addadmgrouproletree", "/admusergroup/loadtreewithoutroot");
}

function submitAddAmdUserGroupForm() {
    var ids = collectAdmUserGroupRoleIds("addadmgrouproletree");
    if (ids == null) {
        return false;
    }

    var data = {
        ids: ids,
        groupname: $("input[name='admgroupusername']").val(),
        describe: $("input[name='admgroupuserdescribe']").val()
    };

    submitAdmUserGroupForm("/admusergroup/addadmusergroup", data, "addadmusergroup", "添加成功");
}

function clearAddAmdUserGroupForm() {
    $('#addamdusergroup').form('clear');
}
