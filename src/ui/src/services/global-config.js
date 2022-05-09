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

import $http from '@/api'

/**
 * 更新全局设置
 * @param {Object} globalConfig 全局设置全量集合
 * @returns {Promise}
 */
export const updateConfig = globalConfig => $http.put('admin/update/system_config/platform_setting', globalConfig)

/**
 * 获取当前用户的全局设置
 * @returns {Promise}
 */
export const getCurrentConfig = () => $http.get('admin/find/system_config/platform_setting/current')

/**
 * 获取默认的全局设置，用来恢复为默认值
 * @returns {Promise}
 */
export const getDefaultConfig = () => $http.get('admin/find/system_config/platform_setting/initial')

/**
 * 更新空闲机集群
 * @param {string} setKey 集群 ID
 * @param {string} setName 集群名称
 * @returns {Promise}
 */
export const updateIdleSet = ({
  setKey,
  setName,
}) => $http.post('topo/update/biz/idle_set', {
  type: 'set',
  set: {
    set_key: setKey,
    set_name: setName,
  }
})

/**
 * 创建空闲机模块
 * @param {string} moduleKey 模块 ID
 * @param {string} moduleName 模块名称
 * @returns {Promise}
 */
export const createIdleModule = ({
  moduleKey,
  moduleName,
}) => $http.post('topo/update/biz/idle_set', {
  type: 'module',
  module: {
    module_key: moduleKey,
    module_name: moduleName,
  }
})

/**
 * 更新空闲机模块，同一个接口，前端分离职责
 * @param {string} moduleKey 模块 ID
 * @param {string} moduleName 模块名称
 * @returns {Promise}
 */
export const updateIdleModule = ({
  moduleKey,
  moduleName,
}) => $http.post('topo/update/biz/idle_set', {
  type: 'module',
  module: {
    module_key: moduleKey,
    module_name: moduleName,
  }
})

/**
 * 删除空闲机模块
 * @param {string} moduleKey 需要删除的模块的 ID，其中 idle、fault、recycle 为内置模块 ID，不可删除。
 * @param {string} moduleName 模块名称
 * @returns {Promise}
 */
export const deleteIdleModule = ({
  moduleKey,
  moduleName
}) => $http.post('topo/delete/biz/extra_moudle', {
  module_key: moduleKey,
  module_name: moduleName
})
