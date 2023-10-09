/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import http from '@/api'

const create = (params, config) => http.post('createmany/quoted/instance', params, config)

const update = (params, config) => http.put('updatemany/quoted/instance', params, config)

const deletemany = (params, config) => http.delete('deletemany/quoted/instance', { ...config, data: params })

const find = (params, config) => http.post('findmany/quoted/instance', params, config)

export default {
  create,
  update,
  deletemany,
  find
}
