package controllers

import (
	"cms/src/common"
	"cms/src/model"
	"cms/src/service"
	"time"

	"github.com/astaxie/beego/validation"
)

type AdmUserController struct {
	BaseController
}

/*
*
进入管理员列表页面
*/
func (this *AdmUserController) List() {
	this.show("admUser/admUserList.html")
}

/*
*
获取分页展示数据
*/
func (this *AdmUserController) Gridlist() {
	pageNum, _ := this.GetInt("page")
	rowsNum, _ := this.GetInt("rows")
	admusermail := this.GetString("admusermail")
	admuserphone := this.GetString("admuserphone")
	admusername := this.GetString("admusername")
	admuserid := this.GetString("admuserid")
	account := this.GetString("admaccout")
	p := common.NewPager(pageNum, rowsNum)
	count, admuser := service.AdmUserService.Gridlist(p, admuserid, admusermail, admusername, admuserphone, account)
	this.jsonResultPager(count, admuser)
}

/*
*
进入添加页面
*/
func (this *AdmUserController) Toaddadmuser() {
	this.show("admUser/addAdmUser.html")
}

// parseAndValidateAdmuser 从请求中提取管理员参数并校验，返回组装好的 model.Admuser 和 groupIds。
// isModify 为 true 时表示修改流程：密码可选、额外读取 userId；
// isModify 为 false 时表示新增流程：密码必填、无需 userId。
// 两个流程的前端提交参数名 groupIds 不一致（新增用 ids，修改用 groupids），在此统一读取。
func (this *AdmUserController) parseAndValidateAdmuser(isModify bool) (*model.Admuser, string, bool) {
	account := this.GetString("account")
	mail := this.GetString("mail")
	name := this.GetString("name")
	phone := this.GetString("phone")
	department := this.GetString("department")
	password := this.GetString("password")

	// 新增和修改前端提交的组 ID 参数名不同，分别读取
	var groupIds string
	if isModify {
		groupIds = this.GetString("groupids")
	} else {
		groupIds = this.GetString("ids")
	}

	valid := validation.Validation{}
	valid.Required(account, "账号").Message("不能为空")
	valid.MaxSize(account, 20, "账号").Message("长度不能超过20个字符")
	valid.Required(mail, "邮箱").Message("不能为空")
	valid.MaxSize(mail, 50, "邮箱").Message("长度不能超过50个字符")
	valid.Email(mail, "邮箱").Message("格式错误")
	valid.Required(name, "姓名").Message("不能为空")
	valid.MaxSize(name, 20, "姓名").Message("长度不能超过20个字符")
	valid.Required(phone, "手机号码").Message("不能为空")
	valid.MaxSize(phone, 15, "手机号码").Message("长度不能超过15个字符")
	valid.Required(department, "部门").Message("不能为空")
	valid.MaxSize(department, 20, "部门").Message("长度不能超过20个字符")

	if isModify {
		// 修改时密码可选，仅在填写时校验长度
		if len(password) > 0 {
			valid.MaxSize(password, 20, "密码").Message("长度不能超过20个字符")
		}
	} else {
		// 新增时密码必填
		valid.Required(password, "密码").Message("不能为空")
		valid.MaxSize(password, 20, "密码").Message("长度不能超过20个字符")
	}

	valid.MinSize(groupIds, 1, "组信息").Message("请至少选择一个")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			this.jsonResult(err.Key + err.Message)
		}
		return nil, "", false
	}

	// 密码加密：新增必填直接加密；修改时仅在填写了密码才加密
	if !isModify || len(password) > 0 {
		password = common.EncodeMessageMd5(password)
	}

	admuser := &model.Admuser{
		Accout:     account,
		Name:       name,
		Mail:       mail,
		Phone:      phone,
		Department: department,
		Password:   password,
		Createtime: time.Now(),
		Updatetime: time.Now(),
		Isdel:      1,
	}

	if isModify {
		userId, _ := this.GetInt64("userId")
		admuser.Id = userId
	}

	return admuser, groupIds, true
}

/*
*
添加管理员
*/
func (this *AdmUserController) Addadmuser() {
	admuser, groupIds, ok := this.parseAndValidateAdmuser(false)
	if !ok {
		return
	}
	if err := service.AdmUserService.AddAdmUser(admuser, groupIds); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/*
*
进入修改页面
*/
func (this *AdmUserController) Tomodifyadmuser() {
	admUserId, _ := this.GetInt64("admUserId")
	admUser, _ := service.AdmUserService.GetUserById(admUserId)
	this.Data["admuser"] = admUser
	this.show("admUser/modifyAdmUser.html")
}

/*
*
修改管理员
*/
func (this *AdmUserController) Modifyyadmuser() {
	admuser, groupIds, ok := this.parseAndValidateAdmuser(true)
	if !ok {
		return
	}
	if err := service.AdmUserService.ModifyAdmUser(admuser, groupIds); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/*
*
删除
*/
func (this *AdmUserController) Delete() {
	userids := this.GetString("userids")
	if err := service.AdmUserService.Delete(userids); err != nil {
		this.jsonResult(err.Error())
	}
	this.jsonResult(SUCCESS)
}

/*
*
获取管理员组列表数据
修改管理员的时候需要加载管理员组列表，并且设置已经选择的权限为选中状态
*/
func (this *AdmUserController) Gridgrouplist() {
	admUserId, _ := this.GetInt64("admUserId")
	groupName := this.GetString("groupName")
	pageNum, _ := this.GetInt("page")
	rowsNum, _ := this.GetInt("rows")
	p := common.NewPager(pageNum, rowsNum)

	count, admuserGroup := service.AdmUserGroupService.Gridlist(groupName, p)
	checkedGroupId := service.AdmUserService.GetAllCheckGroup(admUserId)

	admUserCheckGroup := make([]model.Admusergroupcheck, len(admuserGroup))

	for index, admuser := range admuserGroup {
		admUserCheck := model.Admusergroupcheck{
			Id:         admuser.Id,
			Groupname:  admuser.Groupname,
			Des:        admuser.Des,
			Createtime: admuser.Createtime,
			Updatetime: admuser.Updatetime,
			Isdel:      admuser.Isdel,
			Check:      checkedGroupId[admuser.Id]}
		admUserCheckGroup[index] = admUserCheck
	}

	this.jsonResultPager(count, admUserCheckGroup)
}
