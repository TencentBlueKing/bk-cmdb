package common

import (
	"fmt"
	"strings"
)

const (
	CMDBINDEX       = "cmdb"
	TypeHost        = "host"
	TypeObject      = "object"
	TypeApplication = "biz"

	IndexAggName  = "index_agg"
	IndexAggField = "_index"

	BkBizMetaKey      = "metadata.label.bk_biz_id"
	BkSupplierAccount = "bk_supplier_account"
)

var (
	CmdbFindTypes = []string{BKTableNameBaseInst, BKTableNameBaseHost}
	SpecialChar   = []string{"`", "~", "!", "@", "#", "$", "%", "^", "&", "*",
		"(", ")", "-", "_", "=", "+", "[", "{", "]", "}",
		"\\", "|", ";", ":", "'", "\"", ",", "<", ".", ">", "/", "?"}
)

// GetIndexName get the index of es through mongo's db and collection
func GetIndexName(dbName string, collectionName string) string {
	dbName = strings.ToLower(dbName)
	collectionName = strings.ToLower(collectionName)
	return fmt.Sprintf("%s.%s", dbName, collectionName)
}
