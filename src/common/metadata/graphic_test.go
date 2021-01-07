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
	"reflect"
	"testing"
)

func TestTopoGraphics_FillBlank(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}

	bizid := 0
	account := "0"
	ext := map[string]interface{}{}
	pos := &Position{}

	tests := []struct {
		name   string
		fields fields
		want   *TopoGraphics
	}{
		{"", fields{}, &TopoGraphics{BizID: &bizid, SupplierAccount: &account, Ext: ext, Position: pos}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			if got := tg.FillBlank(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TopoGraphics.FillBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopoGraphics_SetNodeType(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	nodetype := "nodetype"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{NodeType: &nodetype}, args{"nodetype"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetNodeType(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetObjID(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	objid := "objid"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{ObjID: &objid}, args{"objid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetObjID(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetInstID(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val int
	}
	instid := 1
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{InstID: &instid}, args{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetInstID(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetNodeName(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	nodename := "nodename"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{NodeName: &nodename}, args{"nodename"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetNodeName(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetIsPre(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val bool
	}
	ispre := true
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{IsPre: &ispre}, args{true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetIsPre(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetPosition(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val *Position
	}
	pos := &Position{}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{Position: pos}, args{&Position{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetPosition(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetExt(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{Ext: map[string]interface{}{"k": "v"}}, args{map[string]interface{}{"k": "v"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetExt(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetIcon(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	icon := "icon"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{Icon: &icon}, args{"icon"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetIcon(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetScopeType(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	ScopeType := "ScopeType"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{ScopeType: &ScopeType}, args{"ScopeType"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetScopeType(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetScopeID(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	ScopeID := "ScopeID"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{ScopeID: &ScopeID}, args{"ScopeID"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetScopeID(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetBizID(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val int
	}
	BizID := 1
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{BizID: &BizID}, args{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetBizID(tt.args.val)
		})
	}
}

func TestTopoGraphics_SetSupplierAccount(t *testing.T) {
	type fields struct {
		ScopeType       *string
		ScopeID         *string
		NodeType        *string
		ObjID           *string
		IsPre           *bool
		InstID          *int
		NodeName        *string
		Position        *Position
		Ext             map[string]interface{}
		Icon            *string
		BizID           *int
		SupplierAccount *string
		Assts           []GraphAsst
	}
	type args struct {
		val string
	}
	SupplierAccount := "SupplierAccount"
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"", fields{SupplierAccount: &SupplierAccount}, args{"SupplierAccount"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := &TopoGraphics{
				ScopeType:       tt.fields.ScopeType,
				ScopeID:         tt.fields.ScopeID,
				NodeType:        tt.fields.NodeType,
				ObjID:           tt.fields.ObjID,
				IsPre:           tt.fields.IsPre,
				InstID:          tt.fields.InstID,
				NodeName:        tt.fields.NodeName,
				Position:        tt.fields.Position,
				Ext:             tt.fields.Ext,
				Icon:            tt.fields.Icon,
				BizID:           tt.fields.BizID,
				SupplierAccount: tt.fields.SupplierAccount,
				Assts:           tt.fields.Assts,
			}
			tg.SetSupplierAccount(tt.args.val)
		})
	}
}
