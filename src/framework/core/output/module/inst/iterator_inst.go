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

package inst

type iterator struct {
}

func (cli *iterator) Next() (Inst, error) {
	// TODO:当迭代至最后一条数据后，从数据库中读取下一组实例数据，每组固定条数
	return nil, nil
}

func (cli *iterator) ForEach(callbackItem func(item Inst)) error {
	return nil
}
