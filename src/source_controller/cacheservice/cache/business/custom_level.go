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

package business

import (
	"context"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/driver/mongodb"
)

func getMainlineTopology() ([]MainlineTopoAssociation, error) {
	relations := make([]MainlineTopoAssociation, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := mongodb.Client().Table(common.BKTableNameObjAsst).Find(filter).All(context.Background(), &relations)
	if err != nil {
		blog.Errorf("get mainline topology association failed, err: %v", err)
		return nil, err
	}
	return relations, nil
}

// rankMainlineTopology is to rank the biz topology to a array, start from biz to host
func rankMainlineTopology(relations []MainlineTopoAssociation) []string {
	rank := make([]string, 0)
	next := "biz"
	rank = append(rank, next)
	for _, relation := range relations {
		if relation.AssociateTo == next {
			rank = append(rank, relation.ObjectID)
			next = relation.ObjectID
			continue
		} else {
			for _, rel := range relations {
				if rel.AssociateTo == next {
					rank = append(rank, rel.ObjectID)
					next = rel.ObjectID
					break
				}
			}
		}
	}
	return rank
}
