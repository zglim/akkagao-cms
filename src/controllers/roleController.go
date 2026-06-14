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
	//节点展开规则统一交给 service 处理（展开一级目录和当前编辑的节点）
	this.jsonResult(service.RoleService.ListtreeForEdit(id))
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
读取新增/修改权限时提交的表单参数，构造 Role 对象（不含主键 Id）
*/
func (this *RoleController) parseRoleForm() *model.Role {
	pid, _ := this.GetInt64("pid")
	ismenu, _ := this.GetInt8("ismenu")
	return &model.Role{
		Pid:     pid,
		Name:    this.GetString("name"),
		Roleurl: this.GetString("roleurl"),
		Ismenu:  ismenu,
		Des:     this.GetString("describe"),
		Module:  this.GetString("module"),
		Action:  this.GetString("action")}
}

/**
校验权限表单参数，校验不通过时输出错误信息并终止请求
*/
func (this *RoleController) validateRoleForm(role *model.Role) {
	valid := validation.Validation{}
	valid.Required(role.Name, "权限名称").Message("不能为空")
	valid.MaxSize(role.Name, 20, "权限名称").Message("长度不能超过20个字符")
	valid.Required(role.Des, "描述信息").Message("不能为空")
	valid.MaxSize(role.Des, 50, "描述信息").Message("长度不能超过50个字符")

	if valid.HasErrors() {
		// 如果有错误信息，证明验证没通过，打印错误信息
		for _, err := range valid.Errors {
			this.jsonResult((err.Key + err.Message))
		}
	}
}

/**
添加权限
*/
func (this *RoleController) Addrole() {
	role := this.parseRoleForm()
	this.validateRoleForm(role)

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
	role := this.parseRoleForm()
	role.Id, _ = this.GetInt64("id")
	this.validateRoleForm(role)

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
