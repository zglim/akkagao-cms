package service

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// batchInsertRelations 解析逗号分隔的 ID 字符串，逐条转为 int64 后调用 builder 构造关联对象并写入数据库。
// 返回 true 表示全部写入成功，返回 false 表示至少有一条解析或写入失败（部分失败）。
// 兼容"部分写入失败仍返回业务错误"的处理方式，由调用方根据返回值决定是否返回 BizError。
func batchInsertRelations(idsCSV string, builder func(childID int64) interface{}) (allOK bool) {
	allOK = true
	idArray := strings.Split(idsCSV, ",")
	for _, idStr := range idArray {
		childID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			beego.Debug("id 转换成数字异常，id：", idStr)
			allOK = false
			continue
		}
		rel := builder(childID)
		if _, err := o.Insert(rel); err != nil {
			beego.Warn("插入关联关系失败，rel：", rel, err.Error())
			allOK = false
		}
	}
	return
}
