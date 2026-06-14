package controllers

import (
	"cms/src/common"
	"cms/src/model"
	"cms/src/service"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type RoleController struct {
	BaseController
}

// parseRoleFromRequest 从请求中统一读取权限参数（新增和修改共用）
func (this *RoleController) parseRoleFromRequest() *model.Role {
	id, _ := this.GetInt64("id")
	pid, _ := this.GetInt64("pid")
	name := this.GetString("name")
	roleurl := this.GetString("roleurl")
	ismenu, _ := this.GetInt8("ismenu")
	describe := this.GetString("describe")
	module := this.GetString("module")
	action := this.GetString("action")

	return &model.Role{
		Id:      id,
		Pid:     pid,
		Name:    name,
		Roleurl: roleurl,
		Ismenu:  ismenu,
		Des:     describe,
		Module:  module,
		Action:  action,
	}
}

// validateRoleFields 校验权限名称和描述，返回第一条错误信息，通过则返回空字符串
func (this *RoleController) validateRoleFields(name, describe string) string {
	valid := validation.Validation{}
	valid.Required(name, "权限名称").Message("不能为空")
	valid.MaxSize(name, 20, "权限名称").Message("长度不能超过20个字符")
	valid.Required(describe, "描述信息").Message("不能为空")
	valid.MaxSize(describe, 50, "描述信息").Message("长度不能超过50个字符")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return err.Key + err.Message
		}
	}
	return ""
}

/**
进入分页展示页面
*/
func (this *RoleController) List() {
	this.show("role/roleList.html")
}

/**
获取分页展示数据
*/
func (this *RoleController) Gridlist() {
	pageNum, _ := this.GetInt("page")
	rowsNum, _ := this.GetInt("rows")
	p := common.NewPager(pageNum, rowsNum)
	roleid, _ := this.GetInt("roleid")

	roleName := this.GetString("roleName")
	roleUrl := this.GetString("roleUrl")

	count, roles := service.RoleService.Gridlist(p, roleid, roleName, roleUrl)
	this.jsonResultPager(count, roles)
}

/**
加载权限树
*/
func (this *RoleController) Listtree() {
	id, _ := this.GetInt64("id")
	roles := service.RoleService.ListtreeForEdit(id)
	this.jsonResult(roles)
}

/**
进入添加权限页面
*/
func (this *RoleController) Toadd() {
	this.show("role/addRole.html")
}

/**
进入添加权限目录页面
*/
func (this *RoleController) Toadddir() {
	this.show("role/addRoleDir.html")
}

/**
添加权限
*/
func (this *RoleController) Addrole() {
	role := this.parseRoleFromRequest()

	if errMsg := this.validateRoleFields(role.Name, role.Des); errMsg != "" {
		this.jsonResult(errMsg)
	}

	beego.Debug("add role:", role)
	if err := service.RoleService.AddRole(role); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/**
进入修改页面，根据ID查询权限对象
*/
func (this *RoleController) Tomodify() {
	id, _ := this.GetInt64("roleid")
	role, err := service.RoleService.GetRoleById(id)
	if err != nil {
		this.jsonResult(err.Error())
	}
	//this.jsonResult(role)
	this.Data["role"] = role
	this.show("role/modifyRole.html")
}

/**
修改权限
*/
func (this *RoleController) Modify() {
	role := this.parseRoleFromRequest()

	if errMsg := this.validateRoleFields(role.Name, role.Des); errMsg != "" {
		this.jsonResult(errMsg)
	}

	beego.Debug(role)
	if err := service.RoleService.ModifyRole(role); err != nil {
		this.jsonResult("修改失败！")
	}
	this.jsonResult(SUCCESS)
}

/**
删除权限
*/
func (this *RoleController) Deleterole() {
	ids := this.GetStrings("ids")

	if err := service.RoleService.DeleteRole(ids); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}
