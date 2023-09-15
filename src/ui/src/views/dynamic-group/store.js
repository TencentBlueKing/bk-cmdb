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

import Vue from 'vue'
import api from '@/api'
import Utils from '@/components/filters/utils'
import store from '@/store'

function getStorageHeader(type, key, properties) {
  if (!key) {
    return []
  }
  const data = store.state.userCustom[type][key] || []
  const header = []
  data.forEach((propertyId) => {
    const property = properties?.find(property => property?.bk_property_id === propertyId)
    property && header.push(property)
  })
  return header
}

const FilterStore = new Vue({
  data() {
    return {
      config: {},
      fixedPropertyIds: ['bk_host_id', 'bk_host_innerip', 'bk_cloud_id'],
      IP: Utils.getDefaultIP(),
      properties: [],
      propertyGroups: [],
      request() {
        return {
          property: Symbol('property')
        }
      },
    }
  },
  computed: {
    columnConfigProperties() {
      if (this.isDynamicGroupSet) {
        return this.getModelProperties('set')
      }
      const properties = this.properties.filter((property) => {
        const { bk_obj_id: objId, bk_property_id: propId } = property
        const isHost = objId === 'host'
        const isModuleName = objId === 'module' && propId === 'bk_module_name'
        const isSetName = objId === 'set' && propId === 'bk_set_name'
        const isBizName = objId === 'biz' && propId  === 'bk_biz_name'
        return isHost || isModuleName || isSetName || isBizName
      })
      return properties
    },
    isDynamicGroupSet() {
      const { mode } = this.config
      return mode === 'set'
    },
    customHeader() {
      let key
      let properties
      const hostProperties = this.getModelProperties('host')
      const setProperties = this.getModelProperties('set')
      if (this.isDynamicGroupSet) {
        // 集群字段配置
        key = this.config.header && this.config.header.cluster
        properties = [...setProperties]
      }  else {
        key = this.config.header && this.config.header.custom
        const moduleNameProperty = Utils.findPropertyByPropertyId('bk_module_name', this.properties, 'module')
        const setNameProperty = Utils.findPropertyByPropertyId('bk_set_name', this.properties, 'set')
        const bizNameProperty = Utils.findPropertyByPropertyId('bk_biz_name', this.properties, 'biz')
        properties = [...hostProperties, moduleNameProperty, setNameProperty, bizNameProperty]
      }
      return getStorageHeader('usercustom', key, properties)
    },
    presetHeader() {
      const model = this.isDynamicGroupSet ? 'set' : 'host'
      const hostProperties = this.getModelProperties(model)
      // 初始化属性为前6个
      return Utils.getInitialProperties(hostProperties)
    },
    defaultHeader() {
      if (this.customHeader.length) return this.customHeader
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
  },
  methods: {
    // 设置动态分组查询对象选中值
    setDynamicGroupModel(mode) {
      this.config.mode = mode
    },
    getModelProperties(modelId) {
      return [...this.modelPropertyMap[modelId] || []]
    },
    setHeader() {
      // 默认的配置加上条件属性
      const header = [...this.defaultHeader]
      // 固定显示的属性
      const presetProperty = this.isDynamicGroupSet ? [] : this.fixedPropertyIds
        .map(propertyId => this.properties.find(property => property.bk_property_id === propertyId))

      this.header = Utils.getUniqueProperties(presetProperty, header)
    },
    getHeader() {
      // 取之前先设置为最新的值
      this.setHeader()
      // 由于属性数据异步加载，可能会存在无效的数据，过滤后返回
      return this.header.filter(header => header)
    },
    getSearchParams() {
      const header = this.getHeader()
      const dynamicDefaultIp = {
        data: [],
        exact: 1,
        flag: 'bk_host_innerip|bk_host_outerip'
      }
      const params = {
        ip: dynamicDefaultIp,
        ipv6: dynamicDefaultIp
      }

      params.condition = Utils.transformCondition(
        this.condition,
        this.selected,
        header.filter(property => !property?.isInject)
      )

      return params
    },
    updateSelected(selected) {
      this.selected = selected
    },
    setCondition(data = {}) {
      this.condition = data.condition || this.condition
    },
    setDynamicCollection(data) {
      if (!data) {
        return
      }
      const condition = {}
      const selected = []
      Object.keys(data).forEach((key) => {
        const item = data[key]
        selected.push(item?.property)
        condition[item?.property?.id] = {
          operator: item?.operator,
          value: item?.value
        }
      })
      this.updateSelected(selected)
      this.setCondition({ condition })
    },
    async getProperties() {
      const properties = await api.post('find/objectattr/web', {
        bk_obj_id: {
          $in: ['host', 'module', 'set', 'biz']
        }
      }, {
        requestId: this.request.property
      })

      const hostIdProperty = Utils.defineProperty({
        id: 'bk_host_id',
        bk_obj_id: 'host',
        bk_property_id: 'bk_host_id',
        bk_property_name: 'ID',
        bk_property_index: -Infinity,
        bk_property_type: 'int'
      })

      this.properties = [...properties, hostIdProperty]
      return this.properties
    }
  }
})

export async function setupFilterStore(config = {}) {
  FilterStore.config = config
  FilterStore.selected = []
  FilterStore.condition = {}
  await FilterStore.getProperties()

  FilterStore.setHeader()
  return FilterStore
}

export default FilterStore
