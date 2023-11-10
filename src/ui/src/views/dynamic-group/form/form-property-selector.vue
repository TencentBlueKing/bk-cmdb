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
  <div class="property-selector-content" slot="content" ref="propertySelector">
    <div class="property-selector-options">
      <bk-input class="options-filter"
        v-model.trim="filter"
        right-icon="icon-search"
        :placeholder="$t('请输入字段名称或唯一标识')"
        clearable
        v-autofocus>
      </bk-input>
    </div>
    <div class="all-property-selector">
      <div class="property-selector-group clearfix"
        v-for="model in models"
        v-show="isShowGroup(model)"
        :key="model.id">
        <label class="group-label">
          {{model.bk_obj_name}}
          <span class="count">
            （{{matchedPropertyMap[model.bk_obj_id].length}}）
          </span>
        </label>
        <bk-checkbox
          :indeterminate="indeterminate[model.bk_obj_id]"
          :checked="allChecked[model.bk_obj_id]"
          @change="handleChangeAllCheck(model.bk_obj_id, ...arguments)"
          class="allCheck"
        >{{$t('全选')}}</bk-checkbox>
        <div class="group-property-list">
          <bk-checkbox
            :class="['group-property-item', { 'is-checked': isChecked(property) }]"
            v-for="property in matchedPropertyMap[model.bk_obj_id]"
            v-show="isShowProperty(property)"
            :key="property.id"
            :title="property.bk_property_name"
            :checked="isChecked(property)"
            :disabled="disabledPropertyMap[model.bk_obj_id].includes(property.bk_property_id)"
            @change="handleChange(property, ...arguments)">
            <div style="width: calc(100% - 30px);"
              v-bk-tooltips.top-start="{
                disabled: !disabledPropertyMap[model.bk_obj_id].includes(property.bk_property_id),
                content: $t('该字段不支持配置')
              }">
              <div class="group-property-name" v-bk-overflow-tips>{{property.bk_property_name}}</div>
            </div>
            <i class="icon-cc-selected"></i>
          </bk-checkbox>
        </div>
      </div>
    </div>

    <cmdb-data-empty v-if="isShowEmpty" slot="empty"
      :stuff="dataEmpty"
      @clear="handleClearFilter"></cmdb-data-empty>
  </div>
</template>

<script setup>
  import { computed, ref, watch, inject, reactive, nextTick } from 'vue'
  import { t } from '@/i18n'
  import debounce from 'lodash.debounce'
  import store from '@/store'

  const props = defineProps({
    selected: {
      type: Array,
      default: () => ([])
    },
    handler: Function
  })
  const dynamicGroupForm = inject('dynamicGroupForm')
  const propertySelector = ref('')

  watch(() => propertySelector.value, (val) => {
    const height  = store.state.appHeight
    nextTick(() => {
      const { bottom = 0 } = val?.getClientRects()?.[0]
      const dis = bottom - height
      if (dis > -10) {
        // 改变气泡框的高度
        val.getElementsByClassName('all-property-selector')[0].style.height = `${420 - Math.abs(dis)}px`
      }
    })
  })

  const indeterminate = reactive({
    host: false,
    set: false,
    module: false
  })
  const allChecked = reactive({
    host: false,
    set: false,
    module: false
  })
  const disabledPropertyMap = reactive(dynamicGroupForm.disabledPropertyMap)
  const dataEmpty = reactive({
    type: 'empty',
    payload: {
      defaultText: t('暂无数据')
    }
  })

  const matchedPropertyMap = ref(dynamicGroupForm.propertyMap)
  const localSelected = ref([...props.selected])
  const filter = ref('')

  const target = computed(() => dynamicGroupForm.formData.bk_obj_id)
  const propertyMap = computed(() => dynamicGroupForm.propertyMap)
  const models = computed(() => {
    if (target.value === 'host') {
      return dynamicGroupForm.availableModels
    }
    return dynamicGroupForm.availableModels.filter(model => model.bk_obj_id === target.value)
  })
  const isShowEmpty = computed(() => matchedPropertyMap.value.host.length === 0
    && matchedPropertyMap.value.module.length === 0
    && matchedPropertyMap.value.set.length === 0)

  const handleFilter = debounce((filter) => {
    if (!filter.length) {
      matchedPropertyMap.value = propertyMap.value
    } else {
      const matchedPropertyMapOther = {}
      const lowerCaseFilter = filter.toLowerCase()
      Object.keys(propertyMap.value).forEach((modelId) => {
        matchedPropertyMapOther[modelId] = propertyMap.value[modelId].filter((property) => {
          const lowerCaseName = property.bk_property_name.toLowerCase()
          const lowerPropertyId = property.bk_property_id.toLowerCase()
          return lowerCaseName.indexOf(lowerCaseFilter) > -1 || lowerPropertyId.indexOf(lowerCaseFilter) > -1
        })
      })
      matchedPropertyMap.value = matchedPropertyMapOther
    }
    Object.keys(matchedPropertyMap.value)?.forEach(property => allCheckState({ bk_obj_id: property }))
  }, 300)

  const isShowGroup = model => !!matchedPropertyMap.value?.[model.bk_obj_id]?.length

  const isShowProperty = (property) => {
    const modelId = property.bk_obj_id
    return matchedPropertyMap.value?.[modelId]?.some(target => target === property)
  }

  const isChecked = property => localSelected.value.some(target => target.id === property.id)

  const updateLocalSelected = (property, checked) => {
    const index = localSelected.value.findIndex(target => target.id === property.id)
    // 如果checked为true 并且 index === -1 则push
    if (checked && index === -1) {
      localSelected.value.push(property)
    }
    // 如果checked为false 并且 index > -1 则splice
    if (!checked && index > -1) {
      localSelected.value.splice(index, 1)
    }
  }

  const handleChange = (property, checked) => {
    updateLocalSelected(property, checked)
    allCheckState(property)
  }

  const handleChangeAllCheck = (bkObjId, checked) => {
    indeterminate[bkObjId] = false
    allChecked[bkObjId] = checked
    matchedPropertyMap.value?.[bkObjId]?.forEach((target) => {
      const isDisabled = disabledPropertyMap[bkObjId].includes(target.bk_property_id)
      if (!isDisabled) {
        updateLocalSelected(target, checked)
      }
    })
  }

  // 判断相应的全选/半选状态
  const allCheckState = ({ bk_obj_id: bkObjId }) => {
    const length = matchedPropertyMap.value[bkObjId]?.length || 0
    if (length === 0) return
    const matchedPropertyMapIdSet = new Set()
    matchedPropertyMap.value[bkObjId]?.forEach(property => matchedPropertyMapIdSet.add(property?.id))
    const nowChecked = localSelected.value.filter(target => target.bk_obj_id === bkObjId
      && matchedPropertyMapIdSet.has(target.id))?.length || 0
    // 默认是一个都没选的状态
    let isIndeterminate = false // 半选
    let isChecked = false

    if (nowChecked > 0) {
      if (nowChecked === length) {
        isChecked = true
      } else {
        isIndeterminate = true
      }
    }

    indeterminate[bkObjId] = isIndeterminate
    allChecked[bkObjId] = isChecked
  }

  const handleClearFilter = () => filter.value = ''

  defineExpose({
    confirm: () => {
      props.handler && props.handler([...localSelected.value])
    }
  })

  watch(() => filter.value, (filter) => {
    handleFilter(filter)
    dataEmpty.type = filter ? 'search' : 'empty'
  }, {
    immediate: true
  })

</script>

<style lang="scss" scoped>
.property-selector-content {
  width: 400px;
  max-height: 500px;
  padding: 10px 14px;
  margin: -.3rem -.6rem;
}
.all-property-selector {
  max-height: 440px;
  margin-right: -14px;
  margin-left: -14px;
  padding: 0 14px;
  @include scrollbar-y(6px, transparent);
}
.property-selector-group {
  margin-top: 15px;

  .group-label {
    display: block;
    font-weight: bold;
    font-size: 12px;
    color: #313237;
    float: left;

    .count {
      font-size: 12px;
      color: #63656E;
      font-weight: normal;
    }
  }

  .allCheck {
    float: right;
    :deep(.bk-checkbox-text) {
      font-size: 12px;
    }
  }

  .group-property-list {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    margin-top: 4px;
    gap: 3px 14px;
    float: left;
    width: 100%;

    .group-property-item {
      display: inline-flex;
      align-items: center;
      flex: calc(50% - 4px);
      line-height: 32px;
      padding-left: 6px;
      margin-left: -6px;

      .group-property-name {
        display: block;
        width: 100%;
        @include ellipsis;
      }

      .icon-cc-selected {
        font-size: 24px;
        color: #3A84FF;
        opacity: 0;
      }

      &.is-checked,
      &:hover {
        background: #F5F7FA;
        border-radius: 2px;
      }

      &.is-checked {
        :deep(.bk-checkbox-text) {
          color: #3A84FF;
        }
        .icon-cc-selected {
          opacity: 1;
        }
      }

      /deep/ {
        .bk-checkbox {
          flex: 16px 0 0;
          opacity: 0;
          position: absolute;
        }
        .bk-checkbox-text {
          font-size: 12px;
          padding-right: 10px;
          margin: 0;
          width: 100%;
          @include space-between;
          @include ellipsis;
        }
      }
    }
  }
}
</style>
