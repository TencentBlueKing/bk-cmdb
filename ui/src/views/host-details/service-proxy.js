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

/**
 * 主机详情服务代理，代理业务和业务集的接口请求
 */
import router from '@/router'
import store from '@/store'
import { MENU_BUSINESS_SET } from '@/dictionary/menu-symbol.js'
import { findProcessByServiceInstance } from '@/service/business-set/process-instance'
import { findAggregationLabels, findServiceInstanceWithHost } from '@/service/business-set/service-instance'
import { HostService } from '@/service/business-set/host'
import { findTopoPath } from '@/service/business-set/topology'

const inBusinessSet = () => router.currentRoute.matched[0].name === MENU_BUSINESS_SET

export const serviceInstanceProcessesProxy = (params, config) => {
  if (inBusinessSet()) {
    const { bizSetId } = store.state.bizSet
    return findProcessByServiceInstance(bizSetId, params, config)
  }

  return store.dispatch('processInstance/getServiceInstanceProcesses', {
    params,
    config
  })
}

export const historyLabelProxy = (params, config) => {
  if (inBusinessSet()) {
    const { bizSetId } = store.state.bizSet
    return findAggregationLabels(bizSetId, params, config)
  }

  return store.dispatch('instanceLabel/getHistoryLabel', {
    params,
    config
  })
}

export const hostInfoProxy = (params, config) => {
  if (inBusinessSet()) {
    const { bizSetId } = store.state.bizSet
    return HostService.findOne(bizSetId, params, config)
  }

  return store.dispatch('hostSearch/searchHost', { params, config })
}

export const topoPathProxy = (bizId, params, config) => {
  if (inBusinessSet()) {
    const { bizSetId } = store.state.bizSet
    return findTopoPath({ bizSetId, bizId }, params, config)
  }

  return store.dispatch('objectMainLineModule/getTopoPath', { bizId, params, config })
}

export const hostServiceInstancesProxy = (params, config) => {
  if (inBusinessSet()) {
    const { bizSetId } = store.state.bizSet
    return findServiceInstanceWithHost(bizSetId, params, config)
  }
  return store.dispatch('serviceInstance/getHostServiceInstances', { params, config })
}
