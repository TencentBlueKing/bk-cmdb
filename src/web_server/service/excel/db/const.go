package db

// 此文件存放公共的，需要暴露给其他文件的常量定义

// excel file const define
const (
	// HeaderLen excel header length
	HeaderLen = 6

	// HeaderTableLen excel table header length
	HeaderTableLen = 3

	// NameRowIdx excel name row index
	NameRowIdx = 0

	// TypeRowIdx excel type row index
	TypeRowIdx = 1

	// IDRowIdx excel id row index
	IDRowIdx = 2

	// TableNameRowIdx excel table name row index
	TableNameRowIdx = 3

	// TableTypeRowIdx excel table type row index
	TableTypeRowIdx = 4

	// TableIDRowIdx excel table id row index
	TableIDRowIdx = 5

	// InstRowIdx excel instance row index
	InstRowIdx = 6
)

// export instance const difine
const (
	// TopoObjID topo object id
	TopoObjID = "field_topo"

	// IDPrefix topo property id prefix
	IDPrefix = "bk_ext_"
)
