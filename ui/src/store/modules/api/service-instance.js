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

/* eslint-disable no-unused-vars */

import $http from '@/api'

const actions = {
  getHostServiceInstances(context, { params, config }) {
    return $http.post('findmany/proc/web/service_instance/with_host', params, config)
  },
  getModuleServiceInstances(context, { params, config }) {
    return $http.post('findmany/proc/web/service_instance', params, config)
  },
  createProcServiceInstanceWithRaw(context, { params, config }) {
    return $http.post('create/proc/service_instance', params, config)
  },
  createProcServiceInstanceByTemplate(context, { params, config }) {
    return $http.post('create/proc/service_instance', params, config)
  },
  createProcServiceInstancePreview(context, { params, config }) {
    return $http.post('create/proc/service_instance/preview', params, config)
  },
  deleteServiceInstance(context, { config }) {
    return $http.delete('deletemany/proc/service_instance', config)
  },
  removeServiceTemplate(context, { config }) {
    return $http.delete('delete/proc/template_binding_on_module', config)
  },
  getInstanceIpByHost(context, { hostId, config }) {
    return $http.get(`${window.API_HOST}hosts/${hostId}/listen_ip_options`, config)
  },
  previewDeleteServiceInstances(context, { params, config }) {
    return $http.post('deletemany/proc/service_instance/preview', params, config)
  },
  getMoudleProcessList(context, { params, config }) {
    return $http.post('findmany/proc/process_instance/name_ids', params, config)
  },
  getProcessListById(context, { params, config }) {
    return $http.post('findmany/proc/process_instance/detail/by_ids', params, config)
  },
  batchUpdateProcess(context, { params, config }) {
    return $http.put('update/proc/process_instance/by_ids', params, config)
  },
  updateServiceInstance(context, { bizId, params, config }) {
    return $http.put(`updatemany/proc/service_instance/biz/${bizId}`, params, config)
  },

  // 获取主机是否有服务实例
  getNoServiceInstanceHost(context, { params, config }) {
    return $http.post('findmany/proc/host/with_no_service_instance', params, config)
  }
}

export default {
  namespaced: true,
  actions
}
