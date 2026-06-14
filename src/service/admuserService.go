package service

import (
	"cms/src/common"
	"cms/src/model"
	"strings"
	"time"

	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type admUserService struct{}

/**
分页查询管理员列表
*/
func (this *admUserService) Gridlist(pager *common.Pager, admuserid, admusermail, admusername, admuserphone, accout string) (count int, admusers []model.Admuser) {
	countsql := "select count(1) from t_admuser t "
	condition := genAdmUserCondition(admuserid, admusermail, admusername, admuserphone, accout)
	if err := o.Raw(countsql + condition).QueryRow(&count); err != nil || count < 1 {
		beego.Debug("select admuser count err or result is null.")
		return
	}

	listsql := "select id,accout,mail,name,phone,department,password,createtime,updatetime,isdel from t_admuser t "
	if _, err := o.Raw(listsql+condition+common.LIMIT, pager.GetBegin(), pager.GetLen()).QueryRows(&admusers); err != nil {
		beego.Warn("select admuserList from db error.")
		return
	}
	return
}

/**
按照参数拼接sql查询条件
*/
func genAdmUserCondition(admuserid, admusermail, admusername, admuserphone, accout string) (condition string) {
	condition = " where t.isdel = 1 "
	if !strings.EqualFold(admuserid, "") {
		condition += " and t.id = " + admuserid + "'"
	}
	if !strings.EqualFold(admusermail, "") {
		condition += " and t.mail = '" + admusermail + "'"
	}
	if !strings.EqualFold(admusername, "") {
		condition += " and t.name =  '" + admusername + "'"
	}
	if !strings.EqualFold(admuserphone, "") {
		condition += " and t.phone =  '" + admuserphone + "'"
	}
	if !strings.EqualFold(accout, "") {
		condition += " and t.accout =  '" + accout + "'"
	}
	beego.Debug("condition is : ", condition)
	return
}

/**
添加管理员
*/
func (this *admUserService) AddAdmUser(admUser *model.Admuser, groupIds string) error {
	admUserId, err := o.Insert(admUser)
	if err != nil {
		beego.Warn("insert admUser fail, admUser:", admUser, err.Error())
		return &common.BizError{"添加失败,账号已经存在"}
	}
	if rebuildRelations(groupIds, func(id int64) interface{} {
		return &model.UserGroupRel{Userid: admUserId, Groupid: id, Isdel: 1}
	}) {
		return &common.BizError{"出现异常，部分权限添加失败，请补充添加权限。"}
	}
	return nil
}

/**
修改管理员
*/
func (this *admUserService) ModifyAdmUser(admUser *model.Admuser, groupIds string) error {
	updateSql := "UPDATE t_admuser SET "

	set := updateSet(admUser)
	condition := " where id = ? "

	id := admUser.Id
	if _, err := o.Raw(updateSql+set+condition, id).Exec(); err != nil {
		beego.Warn("update admUser fail, admUser:", admUser, err.Error())
		return &common.BizError{"修改失败"}
	}

	//逻辑删除所有用户和组关联关系UserGroupRel
	delRelSql := "update t_user_group_rel set isdel = 0 where userid = ?"
	if _, err := o.Raw(delRelSql, admUser.Id).Exec(); err != nil {
		return &common.BizError{"修改失败"}
	}

	//重新添加关联关系
	if rebuildRelations(groupIds, func(id int64) interface{} {
		return &model.UserGroupRel{Userid: admUser.Id, Groupid: id, Isdel: 1}
	}) {
		return &common.BizError{"出现异常，部分权限修改失败，请补充添加权限。"}
	}

	return nil
}

func updateSet(admUser *model.Admuser) string {
	set := ""
	if !strings.EqualFold(admUser.Password, "") {
		set += " password = '" + admUser.Password + "',"
	}
	if !strings.EqualFold(admUser.Accout, "") {
		set += " accout = '" + admUser.Accout + "',"
	}
	if !strings.EqualFold(admUser.Mail, "") {
		set += " mail = '" + admUser.Mail + "',"
	}
	if !strings.EqualFold(admUser.Name, "") {
		set += " name = '" + admUser.Name + "',"
	}
	if !strings.EqualFold(admUser.Phone, "") {
		set += " phone = '" + admUser.Phone + "',"
	}
	if !strings.EqualFold(admUser.Department, "") {
		set += " department = '" + admUser.Department + "',"
	}
	set += " updatetime = '" + time.Now().Format("2006-01-02 15:04:05") + "'"

	return set
}

/**
删除管理员基本信息
*/
func (this *admUserService) Delete(userids string) error {
	delUserSql := "update t_admuser set isdel = 0 where id in (" + userids + ")"
	if _, err := o.Raw(delUserSql).Exec(); err != nil {
		return &common.BizError{"删除管理员基本信息失败"}
	}
	delRelSql := "update t_user_group_rel set isdel = 0 where userid in (" + userids + ")"
	if _, err := o.Raw(delRelSql).Exec(); err != nil {
		return &common.BizError{"删除管理员和组关系失败"}
	}
	return nil
}

/**
登陆鉴权
*/
func (this *admUserService) Authentication(accout, encodePwd string) (admuser *model.Admuser, err error) {
	selectSql := "select id,password from t_admuser t where t.accout = '" + accout + "' and isdel =1"
	if err := o.Raw(selectSql).QueryRow(&admuser); err != nil {
		if err == orm.ErrNoRows {
			return nil, &common.BizError{"账号不存在"}
		}
		return nil, &common.BizError{"登陆失败，请稍后重试"}
	}
	if !strings.EqualFold(encodePwd, admuser.Password) {
		return nil, &common.BizError{"密码错误"}
	}
	return admuser, nil
}

/**
根据ID查询管理员
*/
func (this *admUserService) GetUserById(id int64) (admuser *model.Admuser, err error) {
	admuser = &model.Admuser{Id: id}
	if err := o.Read(admuser); err != nil {
		if err == orm.ErrNoRows {
			err = &common.BizError{"账号不存在"}
			return nil, err
		}
		err = &common.BizError{"系统错误"}
		return nil, err
	}
	return admuser, nil
}

func (this *admUserService) GetAllCheckGroup(id int64) map[int64]bool {
	var list orm.ParamsList
	num, err := o.Raw("SELECT groupid from t_user_group_rel t where isdel=1 and t.userid = ?", id).ValuesFlat(&list)
	if err != nil || num < 1 {
		return nil
	}
	roleIdMap := make(map[int64]bool, len(list))
	for i := 0; i < len(list); i++ {
		idStr := list[i].(string)
		id, _ := strconv.ParseInt(idStr, 10, 64)
		roleIdMap[id] = true
	}
	return roleIdMap
}
