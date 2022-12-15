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
  <bk-dialog
    v-model="isShow"
    :mask-close="false"
    :draggable="false"
    :width="730"
    :transfer="false"
    @after-leave="handleClosed">
    <div class="title" slot="tools">
      <span>{{$t('筛选条件')}}</span>
      <bk-input class="filter-input" v-model.trim="filter" clearable :placeholder="$t('请输入关键字搜索')"></bk-input>
    </div>
    <section class="property-selector">
      <div class="group"
        v-for="group in renderGroups"
        :key="group.id">
        <h2 class="group-title">
          {{group.name}}
          <span class="group-count">（{{group.children.length}}）</span>
        </h2>
        <ul class="property-list clearfix">
          <li class="property-item fl"
            v-for="property in group.children"
            :key="property.bk_property_id">
            <bk-checkbox class="property-checkbox"
              :checked="isChecked(property)"
              @change="handleToggleProperty(property, ...arguments)">
              {{property.bk_property_name}}
            </bk-checkbox>
          </li>
        </ul>
      </div>
    </section>
    <footer class="footer" slot="footer">
      <i18n class="selected-count"
        v-if="selected.length"
        path="已选择条数"
        tag="div">
        <template #count><span class="count">{{selected.length}}</span></template>
      </i18n>
      <div class="selected-options">
        <bk-button theme="primary" @click="confirm">{{$t('确定')}}</bk-button>
        <bk-button theme="default" @click="close">{{$t('取消')}}</bk-button>
      </div>
    </footer>
  </bk-dialog>
</template>

<script>
  import { mapGetters } from 'vuex'
  import FilterStore from './store'
  import throttle from 'lodash.throttle'
  export default {
    data() {
      return {
        filter: '',
        isShow: false,
        selected: [...FilterStore.selected],
        throttleFilter: throttle(this.handleFilter, 500, { leading: false }),
        renderGroups: []
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      propertyMap() {
        let modelPropertyMap = { ...FilterStore.modelPropertyMap }

        const ignoreHostProperties = ['bk_host_innerip', 'bk_host_outerip', '__bk_host_topology__']
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
      }
    },
    watch: {
      filter: {
        immediate: true,
        handler() {
          this.throttleFilter()
        }
      }
    },
    methods: {
      handleFilter() {
        if (!this.filter.length) {
          this.renderGroups = this.groups
        } else {
          const filteredGroups = []
          const filter = this.filter.toLowerCase()
          this.groups.forEach((group) => {
            const properties = group.children.filter((property) => {
              const name = property.bk_property_name.toLowerCase()
              return name.indexOf(filter) > -1
            })
            if (properties.length) {
              filteredGroups.push({
                ...group,
                children: properties
              })
            }
          })
          this.renderGroups = filteredGroups
        }
      },
      isChecked(property) {
        return this.selected.some(target => target.id === property.id)
      },
      handleToggleProperty(property, checked) {
        if (checked) {
          this.selected.push(property)
        } else {
          const index = this.selected.findIndex(target => target.id === property.id)
          index > -1 && this.selected.splice(index, 1)
        }
      },
      async confirm() {
        FilterStore.updateSelected(this.selected)
        FilterStore.updateUserBehavior(this.selected)
        this.close()
      },
      handleClosed() {
        this.$emit('closed')
      },
      open() {
        this.isShow = true
      },
      close() {
        this.isShow = false
      }
    }
  }
</script>

<style lang="scss" scoped>
    .title {
        display: flex;
        justify-content: space-between;
        align-items: center;
        vertical-align: middle;
        line-height: 31px;
        font-size: 24px;
        color: #444;
        padding: 15px 0 0 24px;
        .filter-input {
            width: 240px;
            margin-right: 45px;
        }
    }
    .property-selector {
        margin: 0 -24px -24px 0;
        height: 350px;
        @include scrollbar-y;
    }
    .group {
        margin-top: 15px;
        .group-title {
            position: relative;
            padding: 0 0 0 15px;
            line-height: 20px;
            font-size: 15px;
            font-weight: bold;
            color: #63656E;
            &:before {
                content: "";
                position: absolute;
                left: 0;
                top: 3px;
                width: 4px;
                height: 14px;
                background-color: #C4C6CC;
            }
            .group-count {
                color: #C4C6CC;
                font-weight: normal;
            }
        }
    }
    .property-list {
        padding: 10px 0 6px 0;
        .property-item {
            width: 33%;
        }
    }
    .property-checkbox {
        display: block;
        margin: 8px 20px 8px 0;
        @include ellipsis;
        /deep/ {
            .bk-checkbox-text {
                max-width: calc(100% - 25px);
                @include ellipsis;
            }
        }
    }
    .footer {
        display: flex;
        .selected-count {
            font-size: 14px;
            line-height: 32px;
            .count {
                color: #2DCB56;
                padding: 0 4px;
            }
        }
        .selected-options {
            margin-left: auto;
        }
    }
</style>
