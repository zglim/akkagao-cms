
//初始化左边tree
$(document).ready(loadTree());
function loadTree() {
    url = "/admusergroup/loadtreewithoutroot"
    var data;
    $.post(url, data, function (result) {
        $.fn.zTree.init($("#addadmgrouproletree"), admusergroupsetting, result);
    });
}

function submitAddAmdUserGroupForm() {
    var result = getCheckedNodeIds("addadmgrouproletree");
    if (!validateCheckedNodes(result.nodes)) {
        return false;
    }

    url = "/admusergroup/addadmusergroup"

    var data = {
        ids: result.ids,
        groupname: $("input[name='admgroupusername']").val(),
        describe: $("input[name='admgroupuserdescribe']").val()
    };

    $.post(url, data, function (result) {
        handleGroupSubmitResult(result, '#addadmusergroup', "添加成功");
    });
}

function clearAddAmdUserGroupForm() {
    $('#addamdusergroup').form('clear');
}
