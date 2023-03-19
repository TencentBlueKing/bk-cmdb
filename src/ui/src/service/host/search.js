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
import { BUILTIN_MODELS } from '@/dictionary/model-constants'
import { transformHostSearchParams, localSort } from '@/utils/tools'

const find = async ({ params, config }) => {
  try {
    const bizCond = params?.condition?.find(cond => cond.bk_obj_id === BUILTIN_MODELS.BUSINESS)
    let url = 'findmany/hosts/search/without_biz'
    if (bizCond?.condition?.length || params.bk_biz_id > 0) {
      url = 'findmany/hosts/search/with_biz'
    }

    const data = http.post(url, transformHostSearchParams(params), config)

    if (data?.info) {
      data.info.forEach((host) => {
        localSort(host.module, 'bk_module_name')
        localSort(host.set, 'bk_set_name')
      })
    }

    return data
  } catch (error) {
    Promise.reject(error)
  }
}

const findOne = async ({ bk_host_id: hostId, bk_biz_id: bizId, config }) => {
  try {
    const { info } = await http.post('hosts/search', {
      bk_biz_id: bizId || -1,
      condition: [
        { bk_obj_id: 'biz', condition: [], fields: [] },
        { bk_obj_id: 'set', condition: [], fields: [] },
        { bk_obj_id: 'module', condition: [], fields: [] },
        { bk_obj_id: 'host', condition: [{
          field: 'bk_host_id',
          operator: '$eq',
          value: hostId
        }], fields: [] }
      ],
      id: { flag: 'bk_host_innerip', exact: 1, data: [] }
    }, config)
    const [instance] = info
    return instance ? instance.host : null
  } catch (error) {
    console.error(error)
    return null
  }
}

const getTopoPath = (data, config) => http.post('find/host/topopath', data, config)

export default {
  find,
  findOne,
  getTopoPath
}
