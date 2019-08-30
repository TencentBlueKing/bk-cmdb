package common

const (
	CMDBINDEX       = "cmdb"
	INDICES         = "indices"
	TypeHost        = "host"
	TypeObject      = "object"
	TypeModel       = "model"
	TypeProcess     = "process"
	TypeApplication = "biz"

	TypeAggName  = "type_agg"
	TypeAggField = "_type"

	BkObjIdAggName    = "bk_obj_id_agg"
	BkObjIdAggField   = "bk_obj_id.keyword"
	BkBizMetaKey      = "metadata.label.bk_biz_id"
	BkSupplierAccount = "bk_supplier_account"
)

var (
	CmdbFindTypes = []string{BKTableNameBaseInst, BKTableNameBaseHost, BKTableNameObjDes}
	SpecialChar   = []string{"`", "~", "!", "@", "#", "$", "%", "^", "&", "*",
		"(", ")", "-", "_", "=", "+", "[", "{", "]", "}",
		"\\", "|", ";", ":", "'", "\"", ",", "<", ".", ">", "/", "?"}
)
