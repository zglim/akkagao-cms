package service

import "strings"

// ConditionBuilder 统一查询条件拼接工具，
// 各 service 的 Gridlist 等查询方法共用此结构来组织 WHERE 子句，
// 新增筛选项时只需在对应 service 中追加 Add 调用即可。
type ConditionBuilder struct {
	base    string   // 基础 WHERE 子句，例如 "t.isdel = 1"
	clauses []string // 附加的 AND 条件
}

// NewConditionBuilder 创建 ConditionBuilder，base 为可选的基础 WHERE 条件。
// 如果 base 为空字符串，则生成的条件不带 WHERE 前缀，仅返回 " AND ..." 片段。
func NewConditionBuilder(base string) *ConditionBuilder {
	return &ConditionBuilder{base: base}
}

// Add 当 value 不为空时追加 "column = value" 形式的条件。
func (cb *ConditionBuilder) Add(column, value string) *ConditionBuilder {
	if value != "" {
		cb.clauses = append(cb.clauses, column+" = "+value)
	}
	return cb
}

// Build 拼接并返回完整的 WHERE 子句。
// 如果没有任何条件且 base 也为空，返回空字符串。
func (cb *ConditionBuilder) Build() string {
	var parts []string
	if cb.base != "" {
		parts = append(parts, cb.base)
	}
	parts = append(parts, cb.clauses...)
	if len(parts) == 0 {
		return ""
	}
	return " WHERE " + strings.Join(parts, " AND ")
}

// AppendTo 返回拼接在已有 WHERE 子句后面的 " AND ..." 片段，
// 适用于主查询已经自带 WHERE 条件（如 "where t.pid = ?"）的场景。
func (cb *ConditionBuilder) AppendTo() string {
	if len(cb.clauses) == 0 {
		return ""
	}
	return " AND " + strings.Join(cb.clauses, " AND ")
}
