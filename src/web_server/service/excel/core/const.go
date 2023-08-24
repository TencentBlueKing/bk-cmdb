package core

// 此文件存放公共的，需要暴露给其他文件的常量定义

// excel file const define
const (
	// HeaderLen excel表头所占行数
	HeaderLen = 6

	// HeaderTableLen excel表头里，表格相关所占用行数
	HeaderTableLen = 3

	// NameRowIdx excel表头「字段名」所在行位置
	NameRowIdx = 0

	// TypeRowIdx excel表头「字段类型」所在行位置
	TypeRowIdx = 1

	// IDRowIdx excel表头「字段标识」所在行位置
	IDRowIdx = 2

	// TableNameRowIdx excel表头「表格字段名」所在行位置
	TableNameRowIdx = 3

	// TableTypeRowIdx excel表头「表格字段类型」所在行位置
	TableTypeRowIdx = 4

	// TableIDRowIdx excel表头「表格字段标识」所在行位置
	TableIDRowIdx = 5

	// InstRowIdx excel 实例数据开始的位置
	InstRowIdx = 6
)

// export instance const define
const (
	// TopoObjID 导出主机实例时，「业务拓扑」这一字段的objID
	TopoObjID = "field_topo"

	// IDPrefix 导出主机实例时，字段标识的前缀
	IDPrefix = "bk_ext_"
)

const (
	// PropDefaultColIdx 属性所在列的默认值
	PropDefaultColIdx = 0
)

type HandleType string

const (
	// AddHost 添加主机
	AddHost HandleType = "addHost"

	// UpdateHost 更新主机
	UpdateHost HandleType = "updateHost"

	// AddInst 添加实例
	AddInst HandleType = "addInst"
)
