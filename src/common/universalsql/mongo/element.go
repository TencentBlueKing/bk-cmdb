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

package mongo

import (
    "configcenter/src/common/mapstr"
    "configcenter/src/common/universalsql"
)

type element struct {
    Key string
    Val interface{}
}

type Eq element

var _ universalsql.ConditionElement = (*Eq)(nil)

func (e *Eq) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        e.Key: mapstr.MapStr{
            universalsql.EQ: e.Val,
        },
    }
}

type Neq element

var _ universalsql.ConditionElement = (*Neq)(nil)

func (n *Neq) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        n.Key: mapstr.MapStr{
            universalsql.NEQ: n.Val,
        },
    }
}

type Gt element

var _ universalsql.ConditionElement = (*Gt)(nil)

func (g *Gt) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        g.Key: mapstr.MapStr{
            universalsql.GT: g.Val,
        },
    }
}

type Lt element

var _ universalsql.ConditionElement = (*Lt)(nil)

func (l *Lt) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        l.Key: mapstr.MapStr{
            universalsql.LT: l.Val,
        },
    }
}

type Gte elemen

var _ universalsql.ConditionElement = (*Gte)(nil)

func (g *Gte) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        g.Key: mapstr.MapStr{
            universalsql.GTE: g.Val,
        },
    }
}

type Lte element

var _ universalsql.ConditionElement = (*Lte)(nil)

func (l *Lte) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        l.Key: mapstr.MapStr{
            universalsql.LTE: l.Val,
        },
    }
}

type In element

var _ universalsql.ConditionElement = (*In)(nil)

func (i *In) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        i.Key: mapstr.MapStr{
            universalsql.IN: i.Val,
        },
    }
}

type Nin element

var _ universalsql.ConditionElement = (*Nin)(nil)

func (n *Nin) ToMapStr() mapstr.MapStr {
    return mapstr.MapStr{
        n.Key: mapstr.MapStr{
            universalsql.NIN: n.Val,
        },
    }
}
