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
  <div class="property-selector" v-bkloading="{ isLoading: pending }">
    <div class="filter">
      <bk-input v-model.trim="keyword" :placeholder="$t('请输入字段名称')"></bk-input>
    </div>
    <div class="group-list"
      v-for="{ group, properties } in groupedPropertyies"
      v-show="properties.length"
      :key="group.id">
      <div class="group-header">
        <label class="group-label">{{group.bk_group_name}}</label>
        <bk-checkbox class="group-checkbox"
          :checked="isAllSelected(properties)"
          @change="setAllSelection(properties, ...arguments)">
          {{$t('全选')}}
        </bk-checkbox>
      </div>
      <ul class="property-list">
        <li class="property-item"
          v-for="property in properties"
          :key="property.id">
          <bk-checkbox class="property-checkbox"
            :title="property.bk_property_name"
            :checked="isSelected(property)"
            :disabled="isPreset(property)"
            @change="setSelection(property, ...arguments)">
            {{property.bk_property_name}}
          </bk-checkbox>
        </li>
      </ul>
    </div>
    <bk-exception type="search-empty" scene="part" v-show="!matchedProperties.length"></bk-exception>
  </div>
</template>

<script>
  import { ref, toRef, watch } from 'vue'
  import useFilter from '@/hooks/utils/filter'
  import useGroupProperty from '@/hooks/utils/group-property'
  import useProperty from '@/hooks/model/property'
  import useGroup from '@/hooks/model/group'
  import useState from './state'
  export default {
    name: 'export-property',
    setup() {
      const [exportState] = useState()
      // 加载属性与属性分组
      const [{ properties, pending }] = useProperty({
        bk_obj_id: exportState.bk_obj_id.value,
        bk_biz_id: exportState.bk_biz_id.value
      })
      const [{ groups }] = useGroup({
        bk_obj_id: exportState.bk_obj_id.value,
        bk_biz_id: exportState.bk_biz_id.value
      })

      // 设置筛选
      const keyword = ref('')
      const [matchedProperties] = useFilter({
        list: properties,
        keyword,
        target: 'bk_property_name'
      })
      const groupedPropertyies = useGroupProperty(groups, matchedProperties)

      // 设置预置属性
      const selection = toRef(exportState, 'fields')
      const presetProperties = ref([])
      watch(properties, (value) => {
        exportState.presetFields.value.forEach((field) => {
          const property = value.find(property => property.bk_property_id === field)
          property && presetProperties.value.push(property)
        })
        selection.value.push(...presetProperties.value)

        exportState.defaultSelectedFields.value.forEach((field) => {
          const property = value.find(property => property.bk_property_id === field)
          property && selection.value.indexOf(property) === -1 && selection.value.push(property)
        })
      })
      const isPreset = property => presetProperties.value.includes(property)

      // 用户勾选属性
      const setSelection = (item, selected) => {
        if (selected) {
          selection.value.push(item)
          return
        }
        const index = selection.value.indexOf(item)
        index > -1 && selection.value.splice(index, 1)
      }
      const setAllSelection = (properties, selected) => {
        if (selected) {
          selection.value = [...new Set([...selection.value, ...properties])]
        } else {
          selection.value = selection.value.filter(property => !properties.includes(property) || isPreset(property))
        }
      }
      const isSelected = property => selection.value.includes(property)
      const isAllSelected = properties => properties.every(property => selection.value.includes(property))

      return {
        keyword,
        matchedProperties,
        groupedPropertyies,
        selection,
        setSelection,
        isSelected,
        setAllSelection,
        isAllSelected,
        isPreset,
        pending
      }
    }
  }
</script>

<style lang="scss" scoped>
  .property-selector {
    padding: 20px 0 0 0;
    .filter {
      width: 340px;
    }
  }
  .group-list {
    margin: 25px 0 0 0;
    .group-header {
      display: flex;
      align-items: center;
      padding: 0 0 10px;
      .group-label {
        flex: 1;
        font-weight: 700;
        line-height: 20px;
        &:before {
          content: "";
          display: inline-block;
          vertical-align: top;
          width: 4px;
          height: 16px;
          background: #dcdee5;
          margin: 2px 9px 0 0;
        }
      }
      .group-checkbox {
        margin-left: auto;
      }
    }
  }
  .property-list {
    display: flex;
    flex-wrap: wrap;
    .property-item {
      width: 33%;
      line-height: 34px;
      .property-checkbox {
        display: inline-flex;
        vertical-align: middle;
        /deep/ {
          .bk-checkbox {
            width: 16px;
          }
          .bk-checkbox-text {
            max-width: 160px;
            @include ellipsis;
          }
        }
      }
    }
  }
</style>
