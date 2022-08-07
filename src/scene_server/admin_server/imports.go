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

package main

import (

	// v3.9.x
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202002131522"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202008101530"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202008121631"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202008172134"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202010131456"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202010151455"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202010151650"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202010211805"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202010281615"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011021415"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011021501"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011171550"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011172152"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011192014"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011201146"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011241510"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011251014"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202011301723"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202012011450"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202101061721"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202102011055"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202102261105"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202103031533"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202103231621"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202104011012"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202104211151"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202105261459"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202106031151"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202106291420"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202106301910"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202107011154"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202107161611"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202107271940"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202107301510"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202107311844"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202112061431"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202112071130"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.9.202112071431"

	// v3.10.x
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202104221702"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202105251041"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202105261459"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202106031151"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202107011735"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202107021056"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202107161611"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202107271945"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202107301510"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202107311844"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202109131607"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202109181134"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202112071130"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202112071431"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202112171521"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202202181012"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202203011516"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202203021455"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202203031512"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202204181447"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202204271725"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202205182148"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202206081408"
	_ "configcenter/src/scene_server/admin_server/upgrader/y3.10.202207251408"
)
