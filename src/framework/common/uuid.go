package common

import (
	"github.com/rs/xid"
)

// UUID a gloable id
func UUID() string {
	return xid.New().String()
}
