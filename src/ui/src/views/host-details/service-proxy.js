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
