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
import { normalizationTopo } from '@/service/container/transition.js'
import { rollReqUseCount } from '../utils.js'

export const requestIds = {
  getTopology: Symbol('getTopology')
}

const getWithStat = async (bizId, config = {}) => {
  try {
    const res = await http.post(`find/topoinst_with_statistics/biz/${bizId}`, {}, {
      requestId: requestIds.getTopology,
      ...config
    })
    return res
  } catch (error) {
    console.error(error)
  }
}

const getContainerTopo = async ({ bizId, params }, config) => {
  try {
    const topoList = await rollReqUseCount(`find/kube/topo_path/bk_biz_id/${bizId}`, params, { limit: 100 }, config)
    return normalizationTopo(topoList, params.bk_reference_id)
  } catch (error) {
    console.error(error)
    return Promise.reject(error)
  }
}

const getContainerTopoNodeStats = async ({ bizId, params }, config) => {
  const [hostStats, podStats] = await Promise.all([
    http.post(`find/kube/${bizId}/topo_node/host/count`, params, config),
    http.post(`find/kube/${bizId}/topo_node/pod/count`, params, config)
  ])

  return {
    hostStats,
    podStats
  }
}

const geFulltWithStat = async (bizId, config = {}) => {
  try {
    const res = await getWithStat(bizId, {
      ...config,
      params: {
        with_default: 1
      }
    })
    return res
  } catch (error) {
    console.error(error)
  }
}

export default {
  getWithStat,
  geFulltWithStat,
  getContainerTopo,
  getContainerTopoNodeStats
}
