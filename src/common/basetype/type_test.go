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
 
package basetype

import (
    "encoding/json"
    "fmt"
    "testing"
    "time"
)

type Data struct {
    // string
    Name *Type
    // int64
    Age *Type
    // bool
    Male *Type
    // time
    Birthday *Type
    Graduate *Type

    Classmate []map[string]*Type
}

func TestDemo(t *testing.T) {
    mapper := make(map[string]*Type)
    mapper["addr"]=NewMustType("shanghai")
    mapper["mailcode"]=NewMustType(2018)
    mapper["bool"]=NewMustType(true)

    data := Data{
        Name: NewMustType("breezeli"),
        Age: NewMustType(23),
        Male: NewMustType(true),
        Birthday: NewMustTimeType(time.Date(1989, time.January, 5, 0, 0, 0, 0, time.UTC).Unix()),
        Graduate: NewMustTimeType("2014-07-01T08:00:00+08:00"),
        Classmate: []map[string]*Type{mapper},
    }

    js, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        fmt.Printf("marshal failed, err: %v\n", err)
        return
    }

    fmt.Println(string(js))

    d := new(Data)
    err = json.Unmarshal(js, d)
    if err != nil {
        fmt.Println("unmarshal failed, err:", err)
        return
    }

    fmt.Println("Name: ", d.Name.String())
    fmt.Println("Age: ", d.Age.Int8())
    fmt.Println("Male: ", d.Male.Bool())
    fmt.Println("Birthday: ", d.Birthday.String())
    fmt.Println("Graduate: ", d.Graduate.String())
    fmt.Println("addr: ", data.Classmate[0]["addr"].String())
}
