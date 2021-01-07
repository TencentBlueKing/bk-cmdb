package service

import (
	"encoding/json"
)

func parseModelBizID(data string) (int64, error) {
	model := &struct {
		BizID int64 `json:"bk_biz_id"`
	}{}

	if len(data) != 0 {
		if err := json.Unmarshal([]byte(data), model); nil != err {
			return 0, err
		}
	}

	return model.BizID, nil
}
