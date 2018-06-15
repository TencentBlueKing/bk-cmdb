package command

import (
	"configcenter/src/source_controller/common/instdata"
	"fmt"
	"strings"
)

type option struct {
	OwnerID  string
	position string
}

// Node topo node define
type Node struct {
	ObjID   string                 `json:"bk_obj_id,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Childs  []*Node                `json:"childs,omitempty"`
	instID  int
	nodekey string
}

func (n *Node) getInstID() (int64, error) {
	id, err := getInt64(n.Data[instdata.GetIDNameByType(n.ObjID)])
	if nil != err {
		return 0, fmt.Errorf("node has no instID: %+v", *n)
	}
	return id, nil
}

func getInt64(v interface{}) (int64, error) {
	switch id := v.(type) {
	case int:
		return int64(id), nil
	case int64:
		return int64(id), nil
	case float32:
		return int64(id), nil
	case float64:
		return int64(id), nil
	default:
		return 0, fmt.Errorf("v is not number : %+v", v)
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

// result: map[parentkey]map[childkey]node
func (n *Node) expand(parent string, modelKeys map[string][]string, result map[string]map[string]*Node) {
	if result[parent] == nil {
		result[parent] = map[string]*Node{}
	}
	nodeKey := n.getNodekey(parent, modelKeys[n.ObjID])
	n.nodekey = nodeKey
	result[parent][nodeKey] = n
	for _, child := range n.Childs {
		child.expand(nodeKey, modelKeys, result)
	}
}

// nodekey: model2[key1:value,key2:value]-model2[key1:value,key2:value]
func (n *Node) getNodekey(parentKey string, keys []string) (nodekey string) {
	nodekey = n.ObjID + "["
	if "" != parentKey {
		nodekey = parentKey + "-" + nodekey
	}
	kv := []string{}
	for _, key := range keys {
		kv = append(kv, key+":"+fmt.Sprint(n.Data[key]))
	}
	nodekey = strings.Join(kv, ",") + "]"
	return nodekey
}
