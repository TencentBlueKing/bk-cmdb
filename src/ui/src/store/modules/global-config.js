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
 * 全局配置数据模型，提供获取全局配置、更新全局配置的能力
 * 接口数据在这里做了适配 UI 的处理，后续服务接口有更新，直接在这里更新模型即可。
 */
import { getCurrentConfig, getDefaultConfig, updateConfig, updateIdleSet, createIdleModule, updateIdleModule, deleteIdleModule } from '@/service/global-config'
import to from 'await-to-js'
import { Base64 } from 'js-base64'
import { language } from '@/i18n'
import cloneDeep from 'lodash/cloneDeep'

const initialConfig = {
  backend: {
    maxBizTopoLevel: 0, // 最大拓扑层级数
  },
  site: {
    name: '蓝鲸配置平台', // 网站名
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
  }
}

// 备份的远程数据，在国际化处理时用来拼装出全量数据。
let currentConfigBackup = null

const state = () => ({
  auth: false, // 权限状态 true 为有权限，否则无
  updating: false, // 更新中状态
  loading: false, // 加载中状态
  language: language === 'zh_CN' ? 'cn' : language, // 后端保存的语言代码和前端的不一致，所以需要转换一下
  config: cloneDeep(initialConfig), // 用户自定义配置，
  defaultConfig: cloneDeep(initialConfig) // 默认配置，用于恢复初始化
})

/**
 * 反序列化远程数据，分离 UI 和服务层，便于后期维护
 * @param {Object} remoteData 后端数据
 * @param {string} lang 当前语种
 * @returns {Object}
 */
const unserializeConfig = (remoteData, lang) => {
  const newState = {
    backend: {
      maxBizTopoLevel: remoteData.backend.max_biz_topo_level,
      snapshotBizName: remoteData.backend.snapshot_biz_name
    },
    site: {
      name: remoteData.site.name.i18n[lang],
      separator: remoteData.site.separator
    },
    footer: {
      contact: remoteData.footer.contact.i18n[lang],
      copyright: remoteData.footer.copyright.i18n[lang]
    },
    validationRules: unserializeValidationRules(remoteData.validation_rules, lang),
    set: remoteData.set,
    idlePool: {
      idle: remoteData.idle_pool.idle,
      fault: remoteData.idle_pool.fault,
      recycle: remoteData.idle_pool.recycle,
      userModules: unserializeUserModules(remoteData.idle_pool.user_modules) || []
    }
  }

  return newState
}

/**
 * 序列化 state 为可传输给后端的数据
 * @param {Object} newConfig 前端 UI state
 * @param {string} lang 当前语种
 */
const serializeState = (newConfig, lang) => {
  const data = {
    backend: {
      max_biz_topo_level: newConfig.backend.maxBizTopoLevel,
      snapshot_biz_name: newConfig.backend.snapshotBizName,
    },
    site: {
      name: {
        ...currentConfigBackup.site.name,
        i18n: {
          ...currentConfigBackup.site.name.i18n,
          [lang]: newConfig.site.name
        }
      },
      separator: newConfig.site.separator
    },
    footer: {
      contact: {
        ...currentConfigBackup.footer.contact,
        i18n: {
          ...currentConfigBackup.footer.contact.i18n,
          [lang]: newConfig.footer.contact
        }
      },
      copyright: {
        ...currentConfigBackup.footer.copyright,
        i18n: {
          ...currentConfigBackup.footer.copyright.i18n,
          [lang]: newConfig.footer.copyright
        }
      }
    },
    validation_rules: serializeValidationRules(newConfig.validationRules, lang),
    set: newConfig.set,
    idle_pool: {
      idle: newConfig.idlePool.idle,
      fault: newConfig.idlePool.fault,
      recycle: newConfig.idlePool.recycle,
      user_modules: serializeUserModules(newConfig.idlePool.userModules, lang) || null
    }
  }

  return data
}

/**
 * 序列化验证规则，作 Base64 转换
 * @param {Object} validationRules 验证规则
 * @param {string} lang 当前语种
 * @returns {Object}
 */
const unserializeValidationRules = (validationRules, lang) => {
  const newRules = {}
  Object.keys(validationRules).forEach((key) => {
    newRules[key] = validationRules[key]
    try {
      newRules[key].value = Base64.decode(newRules[key].value)
    } catch (err) {
      console.log(err)
    }
    newRules[key].message = newRules[key].i18n[lang]
  })
  return newRules
}

/**
 * 序列化验证规则
 * @param {Object} rules 规则列表
 * @param {string} lang 当前语种
 * @returns {Object}
 */
const serializeValidationRules = (rules, lang) => {
  const newRules = {}
  Object.keys(rules).forEach((key) => {
    newRules[key] = rules[key]
    try {
      newRules[key].value = Base64.encode(newRules[key].value)
    } catch (err) {
      console.log(err)
    }
    newRules[key].i18n[lang] = newRules[key].message
    delete newRules[key].message
  })
  return newRules
}

/**
 * 反序列化用户自定义模块数据为前端 UI 可用数据
 * @param {Array} userModules 用户自定义模块数组
 */
const unserializeUserModules = (userModules = []) => userModules?.map(userModule => ({
  moduleKey: userModule.module_key,
  moduleName: userModule.module_name
}))

/**
 * 序列化用户自定义模块数组
 * @param {Array} userModules 用户自定义模块数组
 */
const serializeUserModules = (userModules = []) => userModules?.map(({ moduleKey, moduleName }) => ({
  module_key: moduleKey,
  module_name: moduleName
}))


const mutations = {
  setConfig(state, config) {
    state.config = config
  },
  setDefaultConfig(state, config) {
    state.defaultConfig = config
  },
  setAuth(state, auth) {
    state.auth = auth
  },
  setUpdating(state, updating) {
    state.updating = updating
  },
  setLoading(state, loading) {
    state.loading = loading
  },
}

const actions = {
  clearConfig({ commit }) {
    commit('setConfig', cloneDeep(initialConfig))
  },
  /**
   * 获取默认配置，用于恢复初始化操作
   * @returns {Promise}
   */
  fetchDefaultConfig({ commit, state }) {
    return getDefaultConfig()
      .then((config) => {
        commit('setDefaultConfig', unserializeConfig(config, state.language))
      })
      .catch((err) => {
        throw Error(`获取默认全局设置出现错误：${err.message}`)
      })
  },
  /**
   * 从后台获取配置，获取配置后会 set 配置到 state 中
   * @returns {Promise}
   */
  fetchConfig({ dispatch, commit, state }) {
    commit('setLoading', true)
    return getCurrentConfig()
      .then((config) => {
        currentConfigBackup = Object.freeze(cloneDeep(config))
        commit('setConfig', unserializeConfig(config, state.language))
      })
      .catch((err) => {
        dispatch('clearConfig')
        throw Error(`获取全局设置出现错误：${err.message}`)
      })
      .finally(() => {
        commit('setLoading', false)
      })
  },

  /**
   * 更新配置到后台，更新配置后会 fetchConfig
   * @param config 所有设置
   * @returns {Promise}
   */
  updateConfig({ state, dispatch, commit }, config) {
    return new Promise(async (resolve, reject) => {
      const newConfig = {
        ...cloneDeep(state.config),
        ...cloneDeep(config)
      }
      commit('setUpdating', true)
      const [updateErr] = await to(updateConfig(serializeState(newConfig, state.language)))

      if (updateErr) {
        reject(updateErr)
        commit('setUpdating', false)
        throw Error(`更新全局设置出现错误：${updateErr.message}`)
      }

      const [fetchErr] = await to(dispatch('fetchConfig'))

      if (fetchErr) {
        reject(fetchErr)
        commit('setUpdating', false)
        throw Error(fetchErr.message)
      }

      resolve()

      commit('setUpdating', false)
    })
  },

  /**
   * 更新空闲机池集群，更新后会 fetchConfig
   * @param {string} setKey 集群 Key
   * @param {string} setName 集群名称
   */
  updateIdleSet({ dispatch }, { setKey, setName }) {
    return new Promise(async (resolve, reject) => {
      const [updateErr] = await to(updateIdleSet({
        setKey,
        setName
      }))

      if (updateErr) {
        reject(updateErr)
        throw Error(`更新空闲机集群「${setKey}」出现错误：${updateErr.message}`)
      }

      const [fetchErr] = await to(dispatch('fetchConfig'))

      if (fetchErr) {
        reject(fetchErr)
        throw Error(fetchErr.message)
      }

      resolve()
    })
  },

  /**
   * 创建空闲机模块，创建后会 fetchConfig
   * @param {string} [moduleKey] 当增加新模块时，moduleKey 可以为空
   * @param {string} moduleName 新增模块的模块名称
   */
  createIdleModule({ dispatch }, { moduleKey, moduleName }) {
    return new Promise(async (resolve, reject) => {
      const [updateErr] = await to(createIdleModule({
        moduleKey,
        moduleName
      }))

      if (updateErr) {
        reject(updateErr)
        throw Error(`创建空闲机模块「${moduleKey}」出现错误：${updateErr.message}`)
      }

      const [fetchErr] = await to(dispatch('fetchConfig'))

      if (fetchErr) {
        reject(fetchErr)
        throw Error(fetchErr.message)
      }

      resolve()
    })
  },

  /**
   * 删除空闲模块，删除后会 fetchConfig
   * @param {string} moduleKey 需要删除的模块的 ID
   */
  deleteIdleModule({ dispatch }, { moduleKey, moduleName }) {
    return new Promise(async (resolve, reject) => {
      const [deleteErr] = await to(deleteIdleModule({
        moduleKey,
        moduleName
      }))

      if (deleteErr) {
        reject(deleteErr)
        throw Error(`删除空闲机模块「${moduleKey}」出现错误：${deleteErr.message}`)
      }

      const [fetchErr] = await to(dispatch('fetchConfig'))

      if (fetchErr) {
        reject(fetchErr)
        throw Error(fetchErr.message)
      }

      resolve()
    })
  },

  /**
   * 更新模块，更新后会 fetchConfig
   * @param {string} moduleKey 模块的 ID。
   * @param {string} moduleName 新增模块的模块名称
   */
  updateIdleModule({ dispatch }, { moduleKey, moduleName }) {
    return new Promise(async (resolve, reject) => {
      const [updateErr] = await to(updateIdleModule({
        moduleKey,
        moduleName
      }))

      if (updateErr) {
        reject(updateErr)
        throw Error(`更新空闲机模块「${moduleKey}」出现错误：${updateErr.message}`)
      }

      const [fetchErr] = await to(dispatch('fetchConfig'))

      if (fetchErr) {
        reject(fetchErr)
        throw Error(fetchErr.message)
      }

      resolve()
    })
  },
}

export default {
  namespaced: true,
  state,
  actions,
  mutations
}
