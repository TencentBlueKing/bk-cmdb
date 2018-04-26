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
 
package model

var _ Iterator = (*iterator)(nil)

type iterator struct {
}

func newModelIterator() (Iterator, error) {
	// TODO: 在实例化的时候默认查询一定数量的Model，每次调用Next返回数据中的一个，当读取到缓存的数据的最后一条后开始重新组织数据，直到将数据库里的数据全部读取完毕
	return nil, nil
}

func (cli *iterator) Next() (Model, error) {

	return nil, nil
}
