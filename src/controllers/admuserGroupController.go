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

	//参数校验
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
	ids := this.GetString("ids")
	groupname := this.GetString("groupname")
	describe := this.GetString("describe")
	id, _ := this.GetInt64("id")

	//参数校验
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
加载权限树(用于添加管理员组的时候选择权限)
*/
func (this *AdmUserGroupController) Loadtreewithoutroot() {
	//查询树结构不加载root节点，展开一级目录
	roles := service.RoleService.ListtreeForSelect()
	this.jsonResult(roles)
}

/**
加载权限树(用于修改管理员组的时候选择权限-添加时选择的权限在修改的时候需要选中)
*/
func (this *AdmUserGroupController) Loadtreechecked() {
	admgroupuserid, _ := this.GetInt64("admgroupuserid")
	roleIdMap := service.AdmUserGroupService.GetAllRoleByGroupId(admgroupuserid)
	//查询树结构不加载root节点，展开一级目录并选中已分配的权限
	roles := service.RoleService.ListtreeChecked(roleIdMap)
	this.jsonResult(roles)
}
