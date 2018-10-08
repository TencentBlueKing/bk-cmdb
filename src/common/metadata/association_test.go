/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package metadata

import (
	types "configcenter/src/common/mapstr"
	"reflect"
	"testing"
)

func TestAssociation_Parse(t *testing.T) {
	type fields struct {
		ID               int64
		ObjectID         string
		OwnerID          string
		AsstForward      string
		AsstObjID        string
		AsstName         string
		ObjectAttID      string
		ClassificationID string
		ObjectIcon       string
		ObjectName       string
	}
	type args struct {
		data types.MapStr
	}

	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"id": 1, "bk_supplier_account": "owner",
	})
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Association
		wantErr bool
	}{
		{"", fields{}, args{arg1}, &Association{ID: 1, OwnerID: "owner"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Association{
				ID:               tt.fields.ID,
				ObjectID:         tt.fields.ObjectID,
				OwnerID:          tt.fields.OwnerID,
				AsstForward:      tt.fields.AsstForward,
				AsstObjID:        tt.fields.AsstObjID,
				AsstName:         tt.fields.AsstName,
				ObjectAttID:      tt.fields.ObjectAttID,
				ClassificationID: tt.fields.ClassificationID,
				ObjectIcon:       tt.fields.ObjectIcon,
				ObjectName:       tt.fields.ObjectName,
			}
			got, err := cli.Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Association.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Association.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssociation_ToMapStr(t *testing.T) {
	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"id": 1, "bk_supplier_account": "owner",
	})
	cli := &Association{ID: 1}
	m := cli.ToMapStr()
	i, err1 := m.Int64("id")
	j, err2 := arg1.Int64("id")
	if i != j || err1 != nil || err2 != nil {
		t.Fail()

	}

}

func TestInstAsst_Parse(t *testing.T) {
	type fields struct {
		ID           int64
		InstID       int64
		ObjectID     string
		AsstInstID   int64
		AsstObjectID string
	}
	type args struct {
		data types.MapStr
	}

	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"id": 1, "bk_obj_id": "obj_id",
	})
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *InstAsst
		wantErr bool
	}{
		{"", fields{}, args{arg1}, &InstAsst{ID: 1, ObjectID: "obj_id"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &InstAsst{
				ID:           tt.fields.ID,
				InstID:       tt.fields.InstID,
				ObjectID:     tt.fields.ObjectID,
				AsstInstID:   tt.fields.AsstInstID,
				AsstObjectID: tt.fields.AsstObjectID,
			}
			got, err := cli.Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstAsst.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstAsst.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstAsst_ToMapStr(t *testing.T) {
	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"id": 1,
	})
	cli := &InstAsst{ID: 1}
	m := cli.ToMapStr()
	i, err1 := m.Int64("id")
	j, err2 := arg1.Int64("id")
	if i != j || err1 != nil || err2 != nil {
		t.Fail()

	}
}

func TestMainlineObjectTopo_Parse(t *testing.T) {
	type fields struct {
		ObjID      string
		ObjName    string
		OwnerID    string
		NextObj    string
		NextName   string
		PreObjID   string
		PreObjName string
	}
	type args struct {
		data types.MapStr
	}
	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"bk_obj_id": "obj_id",
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *MainlineObjectTopo
		wantErr bool
	}{
		{"", fields{}, args{arg1}, &MainlineObjectTopo{ObjID: "obj_id"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &MainlineObjectTopo{
				ObjID:      tt.fields.ObjID,
				ObjName:    tt.fields.ObjName,
				OwnerID:    tt.fields.OwnerID,
				NextObj:    tt.fields.NextObj,
				NextName:   tt.fields.NextName,
				PreObjID:   tt.fields.PreObjID,
				PreObjName: tt.fields.PreObjName,
			}
			got, err := cli.Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MainlineObjectTopo.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MainlineObjectTopo.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMainlineObjectTopo_ToMapStr(t *testing.T) {
	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"bk_obj_id": "obj_id",
	})
	cli := &MainlineObjectTopo{ObjID: "obj_id"}
	m := cli.ToMapStr()
	s1, err1 := m.String("bk_obj_id")
	s2, err2 := arg1.String("bk_obj_id")
	if s1 != s2 || err1 != nil || err2 != nil {
		t.Fail()

	}
}
