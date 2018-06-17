package command

import (
	"configcenter/src/common/blog"
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
	nodekey string
	mark    string
}

func (n *Node) getChildObjID() string {
	for _, child := range n.Childs {
		return child.ObjID
	}
	return ""
}

func (n *Node) getInstID() (int64, error) {
	id, err := getInt64(n.Data[instdata.GetIDNameByType(n.ObjID)])
	if nil != err {
		return 0, fmt.Errorf("node has no instID: %+v", *n)
	}
	return id, nil
}

func (n *Node) getChilDInstNames() (instnames []string) {
	instnamefield := n.getChilDInstNameField()
	for _, child := range n.Childs {
		name, ok := child.Data[instnamefield].(string)
		if !ok {
			blog.Errorf("child has no instname field %#v", child.Data[instnamefield])
			continue
		}
		instnames = append(instnames, name)
	}
	return
}
func (n *Node) getChilDInstNameField() string {
	for _, child := range n.Childs {
		return getInstnameField(child.ObjID)
	}
	return ""
}
func (n *Node) getInstNameField() string {
	return getInstnameField(n.ObjID)
}

func getInstnameField(obj string) string {
	switch obj {
	case "biz":
		return "bk_biz_name"
	case "set":
		return "bk_set_name"
	case "module":
		return "bk_module_name"
	default:
		return "bk_inst_name"
	}
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

// result: map[parentkey]map[childkey]node
func (n *Node) walk(walkfunc func(node *Node) error) error {
	if err := walkfunc(n); nil != err {
		return err
	}
	for _, child := range n.Childs {
		if err := child.walk(walkfunc); nil != err {
			return err
		}
	}
	return nil
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

// Topo define
type Topo struct {
	Mainline []string `json:"mainline"`
	Topo     *Node    `json:"topo"`
}
