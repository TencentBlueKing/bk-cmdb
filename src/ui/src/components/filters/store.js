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

import has from 'has'
import api from '@/api'
import Vue from 'vue'
import Utils from './utils'
import store from '@/store'
import i18n from '@/i18n'
import QS from 'qs'
import RouterQuery from '@/router/query'
import throttle from 'lodash.throttle'
import { CONTAINER_OBJECTS, TOPO_MODE_KEYS } from '@/dictionary/container.js'
import containerPropertyService from '@/service/container/property.js'

function getStorageHeader(type, key, properties) {
  if (!key) {
    return []
  }
  const data = store.state.userCustom[type][key] || []
  const header = []
  data.forEach((propertyId) => {
    const property = properties.find(property => property?.bk_property_id === propertyId)
    property && header.push(property)
  })
  return header
}

const FilterStore = new Vue({
  data() {
    return {
      config: {},
      components: {},
      properties: [],
      propertyGroups: [],
      selected: [],
      IP: Utils.getDefaultIP(),
      condition: {},
      header: [],
      collections: [],
      activeCollection: null,
      needResetPage: false,
      throttleSearch: throttle(this.dispatchSearch, 100, { leading: false }),
      topoMode: '',
      fixedPropertyIds: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id']
    }
  },
  computed: {
    bizId() {
      if (typeof this.config.bk_biz_id === 'function') {
        return this.config.bk_biz_id()
      }
      return this.config.bk_biz_id || void 0
    },
    userBehaviorKey() {
      return this.bizId ? 'topology_host_common_filter' : 'resource_host_common_filter'
    },
    collectable() {
      return !!this.bizId
    },
    request() {
      return {
        property: Symbol(this.bizId || 'property'),
        propertyGroup: Symbol(this.bizId || 'propertyGroup'),
        collections: Symbol(this.bizId || 'collections'),
        deleteCollection: id => `delete_collection_${id}`,
        updateCollection: id => `update_collection_${id}`,
        containerProperty: Symbol()
      }
    },
    isContainerTopo() {
      return this.topoMode === TOPO_MODE_KEYS.CONTAINER
    },
    searchHandler() {
      return this.config.searchHandler || (() => {})
    },
    globalHeader() {
      const key = this.config.header && this.config.header.global
      return getStorageHeader('globalUsercustom', key, this.modelPropertyMap.host)
    },
    customHeader() {
      let key
      let properties
      const hostProperties = this.getModelProperties('host')
      if (this.isContainerTopo) {
        // 容器拓扑中的自定义字段使用独立的key
        key = this.config.header && this.config.header.customContainer
        // 容器节点属性
        const nodeProperties = this.getModelProperties(CONTAINER_OBJECTS.NODE)

        properties = [...hostProperties, ...nodeProperties]
      } else {
        key = this.config.header && this.config.header.custom
        const moduleNameProperty = Utils.findPropertyByPropertyId('bk_module_name', this.properties, 'module')
        const setNameProperty = Utils.findPropertyByPropertyId('bk_set_name', this.properties, 'set')
        const bizNameProperty = Utils.findPropertyByPropertyId('bk_biz_name', this.properties, 'biz')
        properties = [...hostProperties, moduleNameProperty, setNameProperty, bizNameProperty]
      }

      return getStorageHeader('usercustom', key, properties)
    },
    presetHeader() {
      const hostProperties = this.getModelProperties('host')

      // 初始化属性为前6个
      const headerProperties = Utils.getInitialProperties(hostProperties)

      // 固定在前的几个属性
      const fixedProperties = this.fixedPropertyIds
        .map(propertyId => Utils.findPropertyByPropertyId(propertyId, hostProperties))

      // 资源-主机
      if (!this.bizId) {
        const topologyProperty = Utils.findPropertyByPropertyId('__bk_host_topology__', hostProperties)
        fixedProperties.push(topologyProperty)
      } else if (this.isContainerTopo) {
        // 容器节点
        const nodeProperties = this.getModelProperties(CONTAINER_OBJECTS.NODE)
        fixedProperties.push(...nodeProperties)
      } else {
        const moduleNameProperty = Utils.findPropertyByPropertyId('bk_module_name', this.properties, 'module')
        const setNameProperty = Utils.findPropertyByPropertyId('bk_set_name', this.properties, 'set')
        fixedProperties.push(moduleNameProperty, setNameProperty)
      }

      return Utils.getUniqueProperties(fixedProperties, headerProperties).slice(0, 6)
    },
    defaultHeader() {
      if (this.customHeader.length) return this.customHeader
      if (this.globalHeader.length) return this.globalHeader
      return this.presetHeader
    },
    modelPropertyMap() {
      const map = {}
      this.properties.forEach((property) => {
        const modelId = property.bk_obj_id
        const modelProperties = map[modelId] || []
        modelProperties.push(property)
        map[modelId] = modelProperties
      })
      return map
    },
    /**
     * 判断是否存在已生效的筛选条件
     * @returns {Boolean}
     */
    hasCondition() {
      const existedSelectedCondition = this.selected?.some((property) => {
        const { value } = this.condition[property.id]
        return value !== null && value !== undefined && !!value.toString().length
      })
      const existedIP = Utils.splitIP(this.IP.text)?.length > 0

      return existedSelectedCondition || existedIP
    },
    columnConfigProperties() {
      if (this.isContainerTopo) {
        const properties = FilterStore.properties.filter(property => property.bk_obj_id === 'host'
        || (property.bk_obj_id === CONTAINER_OBJECTS.NODE))

        return properties
      }

      const properties = FilterStore.properties.filter(property => property.bk_obj_id === 'host'
        || (property.bk_obj_id === 'module' && property.bk_property_id === 'bk_module_name')
        || (property.bk_obj_id === 'set' && property.bk_property_id === 'bk_set_name'))

      return properties
    }
  },
  watch: {
    selected() {
      this.initCondition()
    }
  },
  methods: {
    setupPropertyQuery() {
      const query = QS.parse(RouterQuery.get('filter'))
      const properties = []
      const condition = {}
      try {
        Object.keys(query).forEach((key) => {
          const [id, operator] = key.split('.')
          const property = Utils.findProperty(id, this.properties)
          const value = query[key].toString().split(',')
          if (property && operator && value.length) {
            properties.push(property)
            condition[property.id] = {
              operator: `$${operator}`,
              value: Utils.convertValue(value, `$${operator}`, property)
            }
          }
        })
      } catch (error) {
        Vue.prototype.$warn(i18n.t('解析查询链接出错提示'))
      }
      this.selected = properties
      this.condition = condition
    },
    setupNormalProperty() {
      const userBehavior = store.state.userCustom.usercustom[this.userBehaviorKey] || []
      const normal = userBehavior.length ? userBehavior : [
        ['bk_set_name', 'set'],
        ['bk_module_name', 'module'],
        ['operator', 'host'],
        ['bk_bak_operator', 'host'],
        ['bk_cloud_id', 'host']
      ]
      // eslint-disable-next-line max-len
      this.createOrUpdateCondition(normal.map(([field, model]) => ({ field, model })), { createOnly: true, useDefaultData: true })
    },
    setupIPQuery() {
      const query = QS.parse(RouterQuery.get('ip'))
      const { text = '', exact = 'true', inner = 'true', outer = 'true' } = query
      this.IP = {
        text: text.replace(/,/g, '\n'),
        exact: exact.toString() === 'true',
        inner: inner.toString() === 'true',
        outer: outer.toString() === 'true'
      }
    },
    initCondition() {
      const newConditon = {}
      this.selected.forEach((property) => {
        if (has(this.condition, property.id)) {
          newConditon[property.id] = this.condition[property.id]
        } else {
          newConditon[property.id] = Utils.getDefaultData(property)
        }
      })
      this.condition = newConditon
    },
    setTopoMode(mode) {
      this.topoMode = mode
    },
    setCondition(data = {}) {
      this.condition = data.condition || this.condition
      this.IP = data.IP || this.IP
      this.throttleSearch()
    },
    updateCondition(property, operator, value) {
      this.condition[property.id] = {
        operator,
        value
      }
      this.throttleSearch()
    },
    createOrUpdateCondition(data, options = {}) {
      const { createOnly = false, useDefaultData = false } = options
      data.forEach(({ field, model, operator, value }) => {
        // eslint-disable-next-line max-len
        const existProperty = this.selected.find(property => property.bk_property_id === field && property.bk_obj_id === model)
        if (!existProperty) {
          const property = Utils.findPropertyByPropertyId(field, this.getModelProperties(model))
          if (property) {
            const defaultData = Utils.getDefaultData(property)
            this.selected.push(property)
            this.$set(this.condition, property.id, {
              operator: useDefaultData ? defaultData.operator : operator,
              value: useDefaultData ? defaultData.value : value
            })
          }
        } else if (!createOnly) {
          const defaultData = Utils.getDefaultData(existProperty)
          const condition = this.condition[existProperty.id]
          condition.operator = useDefaultData ? defaultData.operator : operator
          condition.value = useDefaultData ? defaultData.value : value
        }
      })
      this.throttleSearch()
    },
    updateIP(data = {}) {
      Object.assign(this.IP, data)
      this.throttleSearch()
    },
    updateSelected(selected) {
      this.selected = selected
    },
    removeSelected(property) {
      const index = this.selected.findIndex(target => target.id.toString() === property.id.toString())
      index > -1 && this.selected.splice(index, 1)
      this.throttleSearch()
    },
    resetValue(property, silent = false) {
      const properties = Array.isArray(property) ? property : [property]
      properties.forEach((target) => {
        const { operator } = this.condition[target.id]
        const value = Utils.getOperatorSideEffect(target, operator, [])
        this.condition[target.id].value = value
      })
      !silent && this.throttleSearch()
    },
    resetAll(silent) {
      this.IP = Utils.getDefaultIP()
      this.resetValue(this.selected, silent)
    },
    resetIP() {
      this.IP = Utils.getDefaultIP()
      this.throttleSearch()
    },
    dispatchSearch() {
      this.setHeader()
      this.setQuery()
      this.searchHandler(this.condition)
      this.resetPage(false)
    },
    setQuery() {
      const query = {}
      Object.keys(this.condition).forEach((id) => {
        const { operator, value } = this.condition[id]
        if (String(value).length) {
          query[`${id}.${operator.replace('$', '')}`] = Array.isArray(value) ? value.join(',') : value
        }
      })

      const allQuery = {
        filter: QS.stringify(query, { encode: false }),
        ip: QS.stringify(this.IP.text.trim().length ? this.IP : {}, { encode: false }),
        _t: Date.now()
      }

      // 在触发搜索的场景中会设置needResetPage为true，同时需要满足当前业务存在分页的场景
      if (this.needResetPage && RouterQuery.get('page')) {
        allQuery.page = 1
      }

      RouterQuery.set(allQuery)
    },
    setHeader() {
      const suffixPropertyId = Object.keys(this.condition).filter(id => String(this.condition[id].value).trim().length)
      const suffixProperties = this.properties.filter(property => suffixPropertyId.includes(String(property.id)))

      // 默认的配置加上条件属性
      const header = [...this.defaultHeader, ...suffixProperties]

      // 固定显示的属性
      const presetProperty = this.fixedPropertyIds
        .map(propertyId => this.properties.find(property => property.bk_property_id === propertyId))

      this.header = Utils.getUniqueProperties(presetProperty, header)
    },
    setActiveCollection(collection, silent) {
      if (!collection) {
        this.activeCollection = null
        this.resetAll(silent)
        return
      }
      try {
        const IP = JSON.parse(collection.info)
        if (has(IP, 'ip_list')) { // 因老数据的操作符不可兼容，应用收藏条件时直接提示错误并返回
          this.$error(i18n.t('应用收藏条件失败提示'))
          return false
        }
        const queryParams = JSON.parse(collection.query_params)
        const condition = {}
        const selected = []
        let hasMissedProperty = false
        queryParams.forEach((query) => {
          // eslint-disable-next-line max-len
          const property = this.properties.find(property => property.bk_obj_id === query.bk_obj_id && property.bk_property_id === query.field)
          if (property) {
            selected.push(property)
            condition[property.id] = {
              operator: query.operator,
              value: query.value
            }
          } else {
            hasMissedProperty = true
          }
        })
        this.updateSelected(selected)
        this.setCondition({ IP, condition })
        this.activeCollection = collection
        hasMissedProperty && this.$warn(i18n.t('收藏条件未完全解析提示'))
      } catch (error) {
        console.error(error)
        this.$error(i18n.t('应用收藏条件失败提示'))
      }
    },
    getHeader() {
      // 取之前先设置为最新的值
      this.setHeader()

      return this.header
    },
    getSearchParams(topoMode) {
      const transformedIP = Utils.transformIP(this.IP.text)
      const flag = []
      this.IP.inner && flag.push('bk_host_innerip')
      this.IP.outer && flag.push('bk_host_outerip')
      const params = {
        bk_biz_id: this.bizId, // undefined会被忽略
        ip: {
          data: transformedIP.data,
          exact: this.IP.exact ? 1 : 0,
          flag: flag.join('|')
        }
      }

      if (topoMode === TOPO_MODE_KEYS.CONTAINER) {
        params.host_condition = Utils.transformContainerCondition(
          this.condition,
          this.selected,
          this.header.filter(property => !property?.isInject)
        )

        // 容器节点属性条件
        params.node_filter = Utils.transformContainerNodeCondition(
          this.condition,
          // TODO只取容器节点属性
          this.selected,
        )

        if (transformedIP.condition) {
          params.host_condition.condition.push(transformedIP.condition)
        }
      } else {
        params.condition = Utils.transformCondition(
          this.condition,
          this.selected,
          this.header.filter(property => !property?.isInject)
        )
        if (transformedIP.condition) {
          const { condition } = params.condition.find(condition => condition.bk_obj_id === 'host')
          condition.push(transformedIP.condition)
        }
      }

      return params
    },
    getModelProperties(modelId) {
      return [...this.modelPropertyMap[modelId] || []]
    },
    getProperty(propertyId, modelId) {
      return Utils.findPropertyByPropertyId(propertyId, this.getModelProperties(modelId))
    },
    setComponent(name, component) {
      this.components[name] = component
    },
    getComponent(name) {
      return this.components[name]
    },
    async getProperties() {
      const properties = await api.post('find/objectattr', {
        bk_biz_id: this.bizId,
        bk_obj_id: {
          $in: ['host', 'module', 'set', 'biz']
        }
      }, {
        requestId: this.request.property,
        fromCache: true
      })

      // Node的属性
      const nodeProperties = await containerPropertyService.getAll({
        objId: CONTAINER_OBJECTS.NODE
      }, {
        requestId: this.request.containerProperty,
        fromCache: true
      })

      const commonProperties = [...properties, ...nodeProperties]

      const hostIdProperty = Utils.defineProperty({
        id: 'bk_host_id',
        bk_obj_id: 'host',
        bk_property_id: 'bk_host_id',
        bk_property_name: 'ID',
        bk_property_index: -Infinity,
        bk_property_type: 'int'
      })
      const serviceTemplateProperty = Utils.defineProperty({
        id: 'service_template_id',
        bk_obj_id: 'module',
        bk_property_id: 'service_template_id',
        bk_property_name: i18n.t('服务模板'),
        bk_property_index: Infinity,
        bk_property_type: 'service-template'
      })

      if (this.bizId) {
        this.properties = [...commonProperties, hostIdProperty, serviceTemplateProperty]
      } else {
        const topologyProperty = Utils.defineProperty({
          id: Date.now() + 2,
          bk_obj_id: 'host',
          bk_property_id: '__bk_host_topology__',
          bk_property_index: Infinity,
          bk_property_name: i18n.t('业务拓扑'),
          bk_property_type: 'topology',
          required: true,
          isInject: true // 表示属性为前端注入，仅在视图中使用，不需要传递给后台。
        })
        this.properties = [...commonProperties, hostIdProperty, serviceTemplateProperty, topologyProperty]
      }

      return this.properties
    },
    async getPropertyGroups() {
      const groups = await api.post('find/objectattgroup/object/host', {
        bk_biz_id: this.bizId
      }, {
        requestId: this.request.propertyGroup
      })
      this.propertyGroups = groups
      return groups
    },
    async loadCollections() {
      const { info: collections } = await api.post('hosts/favorites/search', {
        condition: {
          bk_biz_id: this.bizId
        }
      }, {
        requestId: this.request.collections
      })
      this.collections = collections
      return this.collections
    },
    async createCollection(data) {
      const response = await api.post('hosts/favorites', data, {
        requestId: this.request.createCollection,
        globalError: false,
        transformData: false
      })
      if (response.result) {
        const collection = { id: response.data.id, ...data }
        this.collections.unshift(collection)
        this.activeCollection = collection
        this.dispatchSearch()
        return Promise.resolve()
      }
      return Promise.reject(response)
    },
    async removeCollection({ id }) {
      await api.delete(`hosts/favorites/${id}`, {
        requestId: this.request.deleteCollection(id)
      })
      const index = this.collections.findIndex(target => target.id === id)
      index > -1 && this.collections.splice(index, 1)
      return Promise.resolve()
    },
    async updateCollection(collection) {
      await api.put(`hosts/favorites/${collection.id}`, collection, {
        requestId: this.request.updateCollection(collection.id)
      })
      const raw = this.collections.find(raw => raw.id === collection.id)
      Object.assign(raw, collection)
      return Promise.resolve()
    },
    async updateUserBehavior(properties) {
      await store.dispatch('userCustom/saveUsercustom', {
        [this.userBehaviorKey]: properties.map(property => [property.bk_property_id, property.bk_obj_id])
      })
      return Promise.resolve()
    },
    resetPage(status = true) {
      this.needResetPage = status
    }
  }
})

/*
* config.bk_biz_id 业务id，仅业务拓扑需要
* config.searchHandler 触发搜索的方法
* config.header.custom 存储用户自定义的表头字段key
* config.header.global 存储全局默认表头字段key
*/

export async function setupFilterStore(config = {}) {
  FilterStore.config = config
  FilterStore.selected = []
  FilterStore.condition = {}
  FilterStore.components = {}
  FilterStore.activeCollection = null
  await Promise.all([
    FilterStore.getProperties(),
    FilterStore.getPropertyGroups()
  ])
  FilterStore.setupIPQuery()
  FilterStore.setupPropertyQuery()
  FilterStore.setupNormalProperty()
  FilterStore.setHeader()
  FilterStore.throttleSearch()
  return FilterStore
}

export default FilterStore
