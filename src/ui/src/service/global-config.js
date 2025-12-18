/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import $http from '@/api'
import { getPlatformConfig, titleSeparator, setShortcutIcon } from '@blueking/platform-config'

export const initialConfig = {
  backend: {
    maxBizTopoLevel: 0, // 最大拓扑层级数
  },
  site: {
    name: '', // 网站名
    separator: '|' // 网站名称路由分隔符
  },
  footer: {
    contact: '', // 联系方式
    copyright: '' // 脚部版权
  },
  validationRules: [], // 用户自定义验证规则
  set: '', // 集群名称
  idlePool: {
    idle: '', // 空闲机
    fault: '', // 故障机
    recycle: '', // 待回收
    userModules: [] // 用户自定义模块
  },
  publicConfig: {
    name: '配置平台',
    nameEn: 'CMDB',
    productName: '蓝鲸配置平台',
    productNameEn: 'BK CMDB',
    brandName: '蓝鲸智云',
    brandNameEn: 'Tencent BlueKing',
    favicon: `${window.Site.publicPath}static/favicon.ico`,
    version: window.Site.buildVersion
  }
}

/**
 * 更新全局设置-业务通用
 * @param {string} type  可选值:backend(更新拓扑层级配置)
 * @param {Object} config 业务通用参数
 * @returns {Promise}
 */
export const updateGlobalConfig = (type, config) => $http.put(`admin/update/config/global_config/${type}`, config)

/**
 * 更新平台配置
 * @param {string} type  可选值: [id_generator]查询id_generator配置内容
 * @param {Object} config 设置参数
 * @returns {Promise}
 */
export const updatePlatformConfig = (type, config) => $http.put(`admin/update/config/platform_config/${type}`, config)

/**
 * 获取当前的平台设置
 * @param {string} type id_generator(id生成器相关配置)
 * @returns {Object}
 */
export const getCurrentPlatformConfig = async type => await $http.get(`admin/find/config/platform_config/${type}`)

/**
 * 获取当前用户的全局设置
 * @returns {Object}
 */
export const getCurrentConfig = async () => {
  const { sharedResUrl } = window.Site

  let publicConfigPromise
  if (sharedResUrl) {
    const resUrl = sharedResUrl.endsWith('/') ? sharedResUrl : `${sharedResUrl}/`
    publicConfigPromise = getPlatformConfig(`${resUrl}bk_cmdb/base.js`, initialConfig.publicConfig)
  } else {
    publicConfigPromise = getPlatformConfig(initialConfig.publicConfig)
  }
  const [config, publicConfig] = await Promise.all([
    $http.get('admin/find/config/global_config'), // 获取业务通用，业务空闲池配置
    publicConfigPromise
  ])

  setShortcutIcon(publicConfig.favicon)

  return {
    ...config,
    site: {
      name: publicConfig.i18n.name,
      separator: titleSeparator
    },
    footer: {
      contact: publicConfig.i18n.footerInfoHTML,
      copyright: publicConfig.footerCopyrightContent
    },
    publicConfig
  }
}

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
}) => $http.post('topo/delete/biz/extra_module', {
  module_key: moduleKey,
  module_name: moduleName
})
