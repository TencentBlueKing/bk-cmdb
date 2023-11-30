<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <addCondition
    ref="addConditionComp"
    :selected="selected"
    :disabled-property-map="disabledProperties"
    :models="models"
    :property-map="propertyMap">
  </addCondition>
</template>

<script>
  import { mapGetters } from 'vuex'
  import FilterStore from './store'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import addCondition from '@/components/add-condition'

  export default {
    components: {
      addCondition
    },
    data() {
      return {
        selected: [...FilterStore.selected]
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      propertyMap() {
        let modelPropertyMap = { ...FilterStore.modelPropertyMap }

        const ignoreHostProperties = ['bk_host_innerip', 'bk_host_outerip', '__bk_host_topology__', 'bk_host_innerip_v6', 'bk_host_outerip_v6']
        modelPropertyMap.host = modelPropertyMap.host
          .filter(property => !ignoreHostProperties.includes(property.bk_property_id))

        // 暂时不支持node对象map类型的字段
        modelPropertyMap.node = modelPropertyMap.node
          ?.filter(property => !['map'].includes(property.bk_property_type))

        const getPropertyMapExcludeBy = (exclude = []) => {
          const excludes = !Array.isArray(exclude) ? [exclude] : exclude
          const propertyMap = []
          for (const [key, value] of Object.entries(modelPropertyMap)) {
            if (!excludes.includes(key)) {
              propertyMap[key] = value
            }
          }
          return propertyMap
        }

        // 资源-主机视图
        if (!FilterStore.bizId) {
          // 非已分配
          if (!FilterStore.isResourceAssigned) {
            return getPropertyMapExcludeBy('node')
          }
          return modelPropertyMap
        }

        // 当前处于业务节点，使用除业务外全量的字段(包括node)
        if (FilterStore.isBizNode) {
          return getPropertyMapExcludeBy('biz')
        }

        // 容器拓扑
        if (FilterStore.isContainerTopo) {
          return {
            host: modelPropertyMap.host || [],
            node: modelPropertyMap.node || [],
          }
        }

        // 业务拓扑主机，不需要业务和Node模型字段
        modelPropertyMap = {
          host: modelPropertyMap.host || [],
          module: modelPropertyMap.module || [],
          set: modelPropertyMap.set || []
        }
        return modelPropertyMap
      },
      groups() {
        const sequence = ['host', 'module', 'set', 'node', 'biz']
        return Object.keys(this.propertyMap).map((modelId) => {
          const model = this.getModelById(modelId) || {}
          return {
            id: modelId,
            name: model.bk_obj_name,
            children: this.propertyMap[modelId]
          }
        })
          .sort((groupA, groupB) => sequence.indexOf(groupA.id) - sequence.indexOf(groupB.id))
      },
      models() {
        return this.groups.map(group => ({
          id: group.id,
          bk_obj_name: group.name,
          bk_obj_id: group.id
        }))
      },
      disabledProperties() {
        const disabledPropertyMap = {}
        this.groups.forEach((group) => {
          disabledPropertyMap[group.id] = group.children
            .filter(item => item.bk_property_type === PROPERTY_TYPES.INNER_TABLE)
            .map(item => item.bk_property_id)
        })
        return disabledPropertyMap
      }
    },
    methods: {
      async confirm() {
        const selected = this.$refs?.addConditionComp?.localSelected ?? this.selected
        FilterStore.updateSelected(selected)
        FilterStore.updateUserBehavior(selected)
        this.close()
      },
      close() {
        this.$emit('closed')
      },
      handleClosed() {
        this.$emit('closed')
      },
    }
  }
</script>
