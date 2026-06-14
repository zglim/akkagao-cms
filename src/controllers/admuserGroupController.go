package controllers

import (
	"cms/src/common"
	"cms/src/model"
	"cms/src/service"
	"time"

	"github.com/astaxie/beego/validation"
)

type AdmUserGroupController struct {
	BaseController
}

/**
进入管理员组管理页面
*/
func (this *AdmUserGroupController) List() {
	this.show("admusergroup/admUserGroupList.html")
}

/**
获取管理员组列表数据
*/
func (this *AdmUserGroupController) Gridlist() {
	groupName := this.GetString("groupName")
	pageNum, _ := this.GetInt("page")
	rowsNum, _ := this.GetInt("rows")
	p := common.NewPager(pageNum, rowsNum)

	count, admuserGroup := service.AdmUserGroupService.Gridlist(groupName, p)
	this.jsonResultPager(count, admuserGroup)
}

/**
进入添加页面
*/
func (this *AdmUserGroupController) Toadd() {
	this.show("admusergroup/addAdmusergroup.html")
}

/**
添加管理员组
*/
func (this *AdmUserGroupController) Addadmusergroup() {
	ids := this.GetString("ids")
	groupname := this.GetString("groupname")
	describe := this.GetString("describe")

	//参数校验(校验不通过会直接输出错误信息并终止本次请求)
	this.validateAdmUserGroupParam(groupname, describe, ids)

	admusergroup := buildAdmUserGroup(0, groupname, describe)
	if err := service.AdmUserGroupService.AddAdmUserGroup(admusergroup, ids); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/**
进入修改管理员组页面
*/
func (this *AdmUserGroupController) Tomodify() {
	id, _ := this.GetInt64("admusergroupid")
	admusergroup := service.AdmUserGroupService.GetAdmUserGroupById(id)
	this.Data["admusergroup"] = admusergroup
	this.show("admusergroup/modifyAdmusergroup.html")
}

/**
修改管理员组
*/
func (this *AdmUserGroupController) Modifyadmusergroup() {
	ids := this.GetString("ids")
	groupname := this.GetString("groupname")
	describe := this.GetString("describe")
	id, _ := this.GetInt64("id")

	//参数校验(校验不通过会直接输出错误信息并终止本次请求)
	this.validateAdmUserGroupParam(groupname, describe, ids)

	admusergroup := buildAdmUserGroup(id, groupname, describe)
	if err := service.AdmUserGroupService.Modifyadmusergroup(admusergroup, ids); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/**
删除管理员组
*/
func (this *AdmUserGroupController) Delete() {
	ids := this.GetString("ids")
	if err := service.AdmUserGroupService.Delete(ids); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/**
加载权限树(用于添加管理员组的时候选择权限)
*/
func (this *AdmUserGroupController) Loadtreewithoutroot() {
	//查询树结构不加载root节点
	roles := service.RoleService.Listtree(false)
	//展开一级目录
	expandRootRoleNodes(roles)
	this.jsonResult(roles)
}

/**
加载权限树(用于修改管理员组的时候选择权限-添加时选择的权限在修改的时候需要选中)
*/
func (this *AdmUserGroupController) Loadtreechecked() {
	admgroupuserid, _ := this.GetInt64("admgroupuserid")
	roleIdMap := service.AdmUserGroupService.GetAllRoleByGroupId(admgroupuserid)
	//查询树结构不加载root节点
	roles := service.RoleService.Listtree(false)
	//展开一级目录
	expandRootRoleNodes(roles)
	//根据已有权限打勾(roleIdMap 为 nil 时不会勾选任何节点)
	for i := range roles {
		if _, ok := roleIdMap[roles[i].Id]; ok {
			roles[i].Checked = true
		}
	}
	this.jsonResult(roles)
}

/**
校验管理员组的公共参数，校验不通过时直接输出第一条错误信息并终止本次请求
*/
func (this *AdmUserGroupController) validateAdmUserGroupParam(groupname, describe, ids string) {
	valid := validation.Validation{}
	valid.Required(groupname, "管理员组名称").Message("不能为空")
	valid.MaxSize(groupname, 20, "管理员组名称").Message("长度不能超过20个字符")
	valid.Required(describe, "描述信息").Message("不能为空")
	valid.MaxSize(describe, 50, "描述信息").Message("长度不能超过50个字符")
	valid.MinSize(ids, 1, "权限").Message("请至少选择一个")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过
		// 打印错误信息
		for _, err := range valid.Errors {
			this.jsonResult((err.Key + err.Message))
		}
	}
}

/**
根据表单参数组装管理员组对象，id 为 0 时表示新增
*/
func buildAdmUserGroup(id int64, groupname, describe string) *model.Admusergroup {
	return &model.Admusergroup{
		Id:         id,
		Groupname:  groupname,
		Des:        describe,
		Createtime: time.Now(),
		Updatetime: time.Now(),
		Isdel:      1}
}

/**
展开权限树的一级目录(pid 为 0 的节点)
*/
func expandRootRoleNodes(roles []model.RoleTree) {
	for i := range roles {
		if roles[i].Pid == 0 {
			roles[i].Open = true
		}
	}
}
