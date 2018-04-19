package model

// model the metadata structure definition of the model
type model struct {
	ObjCls      string `json:"bk_classification_id"`
	ObjIcon     string `json:"bk_obj_icon"`
	ObjectID    string `json:"bk_obj_id"`
	ObjectName  string `json:"bk_obj_name"`
	IsPre       bool   `json:"ispre"`
	IsPaused    bool   `json:"bk_ispaused"`
	Position    string `json:"position"`
	OwnerID     string `json:"bk_supplier_account"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Modifier    string `json:"modifier"`
}
