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
import { transformHostSearchParams, localSort } from '@/utils/tools'

const getSearchUrl = (type) => {
  // 不同接口的权限不同
  const urls = {
    biz: 'findmany/hosts/search/with_biz',
    resource: 'findmany/hosts/search/resource',
    noauth: 'findmany/hosts/search/noauth'
  }
  return urls[type]
}

const getMany = async (type, { params, config }) => {
  const url = getSearchUrl(type)
  const data = http.post(url, transformHostSearchParams(params), config)

  if (data?.info) {
    data.info.forEach((host) => {
      localSort(host.module, 'bk_module_name')
      localSort(host.set, 'bk_set_name')
    })
  }

  return data
}

const getBizHosts = ({ params, config }) => {
  try {
    return getMany('biz', { params, config })
  } catch (error) {
    Promise.reject(error)
  }
}

const getResourceHosts = ({ params, config }) => {
  try {
    return getMany('resource', { params, config })
  } catch (error) {
    Promise.reject(error)
  }
}

const getHosts = ({ params, config }) => {
  try {
    return getMany('noauth', { params, config })
  } catch (error) {
    Promise.reject(error)
  }
}

// 在模型引用-根据名称(ip)查询主机使用
const find = async ({ bk_biz_id: bizId, params, config }) => {
  try {
    const url = getSearchUrl('noauth')
    const { count = 0, info: list = [] } = await http.post(url, {
      bk_biz_id: bizId || -1,
      ...params
    }, config)
    return { count, list }
  } catch (error) {
    console.error(error)
    return null
  }
}

// 在模型关联-查询主机详情中使用
const findOne = async ({ bk_host_id: hostId, bk_biz_id: bizId, config }) => {
  try {
    const url = getSearchUrl('noauth')
    const { info } = await http.post(url, {
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

// 在模型引用-根据ids查询主机使用
const findByIds = async ({ ids, bk_biz_id: bizId, config }) => {
  try {
    const url = getSearchUrl('noauth')
    const { count = 0, info: list = [] } = await http.post(url, {
      bk_biz_id: bizId || -1,
      condition: [
        { bk_obj_id: 'biz', condition: [], fields: [] },
        { bk_obj_id: 'set', condition: [], fields: [] },
        { bk_obj_id: 'module', condition: [], fields: [] },
        { bk_obj_id: 'host', condition: [{
          field: 'bk_host_id',
          operator: '$in',
          value: ids
        }], fields: [] }
      ]
    }, config)
    return { count, list }
  } catch (error) {
    console.error(error)
  }
}


const getTopoPath = (data, config) => http.post('find/host/topopath', data, config)

export default {
  find,
  findOne,
  getTopoPath,
  findByIds,
  getBizHosts,
  getResourceHosts,
  getHosts
}
