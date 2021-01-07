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

func TestClassification_Parse(t *testing.T) {
	type fields struct {
		ID                 int64
		ClassificationID   string
		ClassificationName string
		ClassificationType string
		ClassificationIcon string
		OwnerID            string
	}
	type args struct {
		data types.MapStr
	}

	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"id": 1, "bk_classification_id": "bk_classification_id",
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Classification
		wantErr bool
	}{
		{"", fields{}, args{arg1}, &Classification{ID: 1, ClassificationID: "bk_classification_id"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Classification{
				ID:                 tt.fields.ID,
				ClassificationID:   tt.fields.ClassificationID,
				ClassificationName: tt.fields.ClassificationName,
				ClassificationType: tt.fields.ClassificationType,
				ClassificationIcon: tt.fields.ClassificationIcon,
				OwnerID:            tt.fields.OwnerID,
			}
			got, err := cli.Parse(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Classification.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Classification.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassification_ToMapStr(t *testing.T) {
	arg1, _ := types.NewFromInterface(map[string]interface{}{
		"id": 1, "bk_classification_id": "bk_classification_id",
	})
	cli := &Classification{ID: 1}
	m := cli.ToMapStr()
	i, err1 := m.Int64("id")
	j, err2 := arg1.Int64("id")
	if i != j || err1 != nil || err2 != nil {
		t.Fail()

	}
}
