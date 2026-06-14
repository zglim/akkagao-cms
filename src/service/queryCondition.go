package service

import (
	"bytes"
)

/**
conditionBuilder 统一拼接列表查询的 where 过滤条件。

各 service 的筛选条件都通过它来组织，新增筛选项时只需要追加一次 equal/appendRaw
调用，不需要每个 service 再各写一套风格不一致的字符串拼接。
*/
type conditionBuilder struct {
	buf bytes.Buffer
}

func newConditionBuilder() *conditionBuilder {
	return &conditionBuilder{}
}

/**
where 原样写入基础的 where 子句（例如 " where t.isdel = 1 "）。
*/
func (b *conditionBuilder) where(clause string) *conditionBuilder {
	b.buf.WriteString(clause)
	return b
}

/**
equal 在 value 非空时追加一个标准的等值过滤：" and column = 'value'"。
新增的普通筛选项都应当使用它，保证风格一致。
*/
func (b *conditionBuilder) equal(column, value string) *conditionBuilder {
	if value != "" {
		b.buf.WriteString(" and ")
		b.buf.WriteString(column)
		b.buf.WriteString(" = '")
		b.buf.WriteString(value)
		b.buf.WriteString("'")
	}
	return b
}

/**
appendRaw 在 guard 非空时原样追加 fragment。

用于少数历史遗留、写法不标准的过滤条件（保持其原有 SQL 不变）；
普通的等值筛选请使用 equal。
*/
func (b *conditionBuilder) appendRaw(guard, fragment string) *conditionBuilder {
	if guard != "" {
		b.buf.WriteString(fragment)
	}
	return b
}

func (b *conditionBuilder) String() string {
	return b.buf.String()
}
