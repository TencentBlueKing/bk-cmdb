package service

import (
	"encoding/json"

	"configcenter/src/common/metadata"
)

func parseMetadata(data string) (*metadata.Metadata, error) {
	meta := new(metadata.Metadata)
	if len(data) != 0 {
		if err := json.Unmarshal([]byte(data), meta); nil != err {
			return nil, err
		}
	}

	if meta.Label == nil || len(meta.Label) == 0 {
		meta = nil
	}

	return meta, nil
}
