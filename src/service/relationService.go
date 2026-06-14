package service

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

/**
rebuildRelations 解析逗号分隔的 id 列表，逐条转换为数字后通过 build 构造关联对象并写入。
只要任意一条 id 解析失败或写入失败，就返回 true，沿用服务层“部分写入失败仍返回业务错误”的处理方式，
由调用方据此返回对应的业务错误文案。后续新增其它关联关系时，只需提供对应的 build 即可复用本方法。
*/
func rebuildRelations(ids string, build func(id int64) interface{}) (partialFailure bool) {
	for _, raw := range strings.Split(ids, ",") {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			beego.Warn("id 转换成数字异常，id：", raw)
			partialFailure = true
			continue
		}
		rel := build(id)
		if _, err := o.Insert(rel); err != nil {
			beego.Warn("添加关联关系失败", rel, err.Error())
			partialFailure = true
			continue
		}
	}
	return partialFailure
}
