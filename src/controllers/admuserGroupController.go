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
	ids, groupname, describe, errMsg := this.parseGroupParams()
	if errMsg != "" {
		this.jsonResult(errMsg)
	}

	admusergroup := &model.Admusergroup{
		Groupname:  groupname,
		Des:        describe,
		Createtime: time.Now(),
		Updatetime: time.Now(),
		Isdel:      1}
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
	ids, groupname, describe, errMsg := this.parseGroupParams()
	if errMsg != "" {
		this.jsonResult(errMsg)
	}
	id, _ := this.GetInt64("id")

	admusergroup := &model.Admusergroup{
		Id:         id,
		Groupname:  groupname,
		Des:        describe,
		Createtime: time.Now(),
		Updatetime: time.Now(),
		Isdel:      1}
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
解析并校验管理员组参数（新增和修改共用）
*/
func (this *AdmUserGroupController) parseGroupParams() (ids, groupname, describe, errMsg string) {
	ids = this.GetString("ids")
	groupname = this.GetString("groupname")
	describe = this.GetString("describe")

	valid := validation.Validation{}
	valid.Required(groupname, "管理员组名称").Message("不能为空")
	valid.MaxSize(groupname, 20, "管理员组名称").Message("长度不能超过20个字符")
	valid.Required(describe, "描述信息").Message("不能为空")
	valid.MaxSize(describe, 50, "描述信息").Message("长度不能超过50个字符")
	valid.MinSize(ids, 1, "权限").Message("请至少选择一个")

	if valid.HasErrors() {
		errMsg = valid.Errors[0].Key + valid.Errors[0].Message
	}
	return
}

/**
展开一级目录节点
*/
func expandFirstLevelNodes(roles []model.RoleTree) {
	for i, role := range roles {
		if role.Pid == 0 {
			roles[i].Open = true
		}
	}
}

/**
加载权限树(用于添加管理员组的时候选择权限)
*/
func (this *AdmUserGroupController) Loadtreewithoutroot() {
	//查询树结构不加载root节点
	roles := service.RoleService.Listtree(false)
	//展开一级目录
	expandFirstLevelNodes(roles)
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
	expandFirstLevelNodes(roles)
	//根据已有权限打勾
	if roleIdMap != nil {
		for i, role := range roles {
			if _, ok := roleIdMap[role.Id]; ok {
				roles[i].Checked = true
			}
		}
	}
	this.jsonResult(roles)
}
