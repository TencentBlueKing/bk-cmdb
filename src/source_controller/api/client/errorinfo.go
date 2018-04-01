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
 
package client

import (
	"errors"
)

var (
	Err_Not_Found_Anything = errors.New("found nothing")
	Err_Not_Set_Input      = errors.New("input nothing")
	Err_Not_Set_ObjID      = errors.New("not set objid")
	Err_Request_Object     = errors.New("http request failed")
	Err_Decode_Json        = errors.New("decode json failed")
	Err_Creaate_Object     = errors.New("create object failed")
	Err_Request            = errors.New("request failed, error")
)
