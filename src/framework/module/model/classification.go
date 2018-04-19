package model

// classification the model classification definition
type classification struct {
	ClassificationID   string `json:"bk_classification_id"`
	ClassificationName string `json:"bk_classification_name"`
	ClassificationType string `json:"bk_classification_type"`
	ClassificationIcon string `json:"bk_classification_icon"`
}
