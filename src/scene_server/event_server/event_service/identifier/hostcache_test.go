package identifier

import (
	"configcenter/src/scene_server/event_server/types"
	"encoding/json"
	"fmt"
	"testing"
)

func TestIdentifier(t *testing.T) {
	i := HostIdentifier{
		HostID:          1,
		HostName:        "etcd-1",
		SupplierID:      0,
		SupplierAccount: "0",
		CloudID:         0,
		CloudName:       "default area",
		InnerIP:         "192.168.1.7",
		OuterIP:         "192.168.1.10",
		OSType:          "1",
		OSName:          "linux centos",
		Memory:          64131,
		CPU:             24,
		Disk:            271,
		Module: map[string]*Module{
			"module10": {
				BizID:      10,
				BizName:    "biz10",
				SetID:      10,
				SetName:    "set10",
				ModuleID:   10,
				ModuleName: "module10",
				SetStatus:  "1",
				SetEnv:     "1",
			},
		},
	}

	e := types.EventInst{}
	e.Data = nil
	e.EventType = types.EventTypeRelation
	e.ObjType = "hostidentifier"
	e.Action = types.EventActionUpdate

	d := types.EventData{CurData: i}
	e.Data = append(e.Data, d)

	out, _ := json.Marshal(e)
	fmt.Printf("%s\n", out)
}
