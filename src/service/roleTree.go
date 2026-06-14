package service

import (
	"bytes"
	"cms/src/model"

	"github.com/astaxie/beego"
)

/**
权限树 / 菜单树相关的数据准备。

把"查询节点 + 展开 + 选中 + 点击字段组装"这些零碎步骤集中在 service 层，
控制器只负责调用对应入口，不再各自重复一遍树节点的整理逻辑。
*/

/**
listtree 查询权限树的原始节点。
@param needRoot:查询的数据集中是否需要包含root节点
*/
func (this *roleService) listtree(needRoot bool) []model.RoleTree {
	var buf bytes.Buffer
	buf.WriteString("SELECT id, pid, name, roleurl, ismenu, des from t_role t ")
	if !needRoot {
		buf.WriteString(" where t.id != 0")
	}
	var roles []model.RoleTree
	beego.Debug("查询权限树sql：", buf.String())
	_, err := o.Raw(buf.String()).QueryRows(&roles)
	if err != nil {
		beego.Error("查询权限树的role列表异常，error message：", err.Error())
	}
	beego.Debug("生成权限树的数据：", roles)
	return roles
}

/**
ListtreeForEdit 用于添加/修改权限页面的权限树：包含 root 节点，
展开一级目录以及当前正在操作的节点（让新增/选中的节点及时展示出来）。
*/
func (this *roleService) ListtreeForEdit(id int64) []model.RoleTree {
	roles := this.listtree(true)
	expandTopDirs(roles)
	expandNode(roles, id)
	return roles
}

/**
ListtreeForSelect 用于选择权限时的权限树（不含 root 节点），并展开一级目录。
*/
func (this *roleService) ListtreeForSelect() []model.RoleTree {
	roles := this.listtree(false)
	expandTopDirs(roles)
	return roles
}

/**
ListtreeChecked 用于修改管理员组时选择权限：在 ListtreeForSelect 的基础上，
把 checked 集合中的节点设置为选中状态（checked 为 nil 时只展开一级目录）。
*/
func (this *roleService) ListtreeChecked(checked map[int64]bool) []model.RoleTree {
	roles := this.listtree(false)
	expandTopDirs(roles)
	markChecked(roles, checked)
	return roles
}

/**
expandTopDirs 展开所有一级目录（pid 为 0 的根节点）。
*/
func expandTopDirs(roles []model.RoleTree) {
	for i := range roles {
		if roles[i].Pid == 0 {
			roles[i].Open = true
		}
	}
}

/**
expandNode 展开 id 匹配的节点（例如刚刚新增的权限节点）。
*/
func expandNode(roles []model.RoleTree, id int64) {
	for i := range roles {
		if roles[i].Id == id {
			roles[i].Open = true
		}
	}
}

/**
markChecked 把 id 命中 checked 集合的节点设置为选中。
*/
func markChecked(roles []model.RoleTree, checked map[int64]bool) {
	for i := range roles {
		if checked[roles[i].Id] {
			roles[i].Checked = true
		}
	}
}

/**
decorateMenu 整理左侧菜单树：展开所有父节点，并给带 url 的叶子节点组装点击字段。
*/
func decorateMenu(roles []model.RoleTree) {
	pidMap := make(map[int64]bool, len(roles))
	for _, role := range roles {
		pidMap[role.Pid] = true
	}
	for i := range roles {
		//展开所有父节点
		if pidMap[roles[i].Id] {
			roles[i].Open = true
			continue
		}
		if roles[i].Roleurl != "" {
			roles[i].Click = "click: addTab('" + roles[i].Name + "','" + roles[i].Roleurl + "')"
		}
	}
}
