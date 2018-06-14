package command

import (
	"bytes"
	"configcenter/src/source_controller/common/instdata"
	"encoding/json"
	"fmt"
)

type option struct {
	OwnerID  string
	position string
}

// Node topo node define
type Node struct {
	ObjID  string                 `json:"bk_obj_id,omitempty"`
	InstID int                    `json:"bk_inst_id,omitempty"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Childs []*Node                `json:"childs,omitempty"`
}

// type Data map[string]interface{}

func (d Node) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(d)
	fmt.Printf("%s\n")
	return buf.Bytes(), err
}

func (n *Node) getInstID() (int, error) {
	switch id := n.Data[instdata.GetIDNameByType(n.ObjID)].(type) {
	case int:
		return id, nil
	case int64:
		return int(id), nil
	default:
		return 0, fmt.Errorf("node has no instID: %+v", *n)
	}
}

func newNode(objID string) *Node {
	return &Node{ObjID: objID, Data: map[string]interface{}{}, Childs: []*Node{}}
}

// Topo define
type Topo struct {
	Mainline []string `json:"mainline"`
	Topo     *Node    `json:"topo"`
}
