package types

import (
	"encoding/json"
)

// Merge merge second into self,if the key is the same then the new value replaces the old value.
func (cli *MapStr) Merge(second MapStr) {
	for key, val := range second {
		(*cli)[key] = val
	}
}

// ToJSON convert to json string
func (cli *MapStr) ToJSON() []byte {
	js, _ := json.Marshal(cli)
	return js
}
