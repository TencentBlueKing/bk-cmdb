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
 
package mgoclient

import (
	"testing"
	"time"
)

func TestEscapeHtml(t *testing.T) {
	var (
		testingScript1 = `<script>alert(0)</script>`
		testingScript2 = `<script>alert(1)</script>`
		testingScript3 = `<script>alert(2)</script>`
		expect1        = `&lt;script&gt;alert(0)&lt;/script&gt;`
		expect2        = `&lt;script&gt;alert(1)&lt;/script&gt;`

		testingScript4 = "<>&"
		expect4        = `&lt;&gt;&#39;&amp;`
	)

	type Foo struct {
		Description string
		LastTime    time.Time
		Name        *string
		F           *Foo
	}

	type Embed struct {
		Address string
		Foo
	}

	origin1 := map[string]interface{}{
		"last_time":   time.Now(),
		"description": testingScript1,
		"name":        &testingScript2,
	}
	EscapeHtml(origin1)
	if origin1["description"] != expect1 {
		t.Errorf("expect %s, got %s", expect1, origin1["description"])
	}

	origin2 := Foo{
		LastTime:    time.Now(),
		Description: testingScript1,
		Name:        &testingScript2,
	}
	EscapeHtml(&origin2)
	if origin2.Description != expect1 {
		t.Errorf("expect %s, got %s", expect1, origin2.Description)
	}
	if *origin2.Name != expect2 {
		t.Errorf("expect %s, got %s", expect2, *origin2.Name)
	}

	origin3 :=
		Embed{
			Address: testingScript3,
			Foo: Foo{
				LastTime:    time.Now(),
				Description: testingScript1,
				Name:        &testingScript2,
			},
		}
	EscapeHtml(&origin3)
	if origin3.Description != expect1 {
		t.Errorf("expect %s, got %s", expect1, origin3.Description)
	}
	if *origin3.Name != expect2 {
		t.Errorf("expect %s, got %s", expect2, *origin3.Name)
	}

	origin4 := map[string]interface{}{
		"last_time":   time.Now(),
		"description": testingScript4,
	}
	EscapeHtml(origin4)
	if origin1["description"] != expect4 {
		t.Errorf("expect ---%s---, got ---%s---", expect4, origin4["description"])
	}
}
