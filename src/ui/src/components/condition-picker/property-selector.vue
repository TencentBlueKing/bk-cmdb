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
  <div class="property-selector-content" :style="{
    height: `${height}px`
  }">
    <div class="property-selector-options">
      <bk-input class="options-filter"
        v-model.trim="filter"
        right-icon="icon-search"
        :placeholder="$t('请输入字段名称或唯一标识')"
        clearable
        v-autofocus>
      </bk-input>
    </div>
    <div class="property-selector-container">
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
          ref="checkboxRef"
          :disabled="getCheckDisabled(model.bk_obj_id)"
          :indeterminate="indeterminate[model.bk_obj_id]"
          :checked="allChecked[model.bk_obj_id]"
          @change="handleChangeAllCheck(model.bk_obj_id, ...arguments)"
          class="all-check"
        >{{$t('全选')}}</bk-checkbox>
        <div class="group-property-list">
          <bk-checkbox
            :class="['group-property-item',
                     { 'is-checked': isChecked(property),
                       'is-checked-diabled': isDisabled(model, property) }]"
            v-for="property in matchedPropertyMap[model.bk_obj_id]"
            v-show="isShowProperty(property)"
            :key="property.id"
            :title="property.bk_property_name"
            :checked="isChecked(property)"
            :disabled="isDisabled(model, property)"
            @change="handleChange(property, ...arguments)">
            <div style="width: calc(100% - 30px);"
              v-bk-tooltips.top-start="{
                disabled: !isDisabled(model, property),
                content: getDisabledTip(property)
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
  import { computed, ref, watch, reactive } from 'vue'
  import { t } from '@/i18n'
  import debounce from 'lodash.debounce'
  import { DYNAMIC_GROUP_COND_NAMES, DYNAMIC_GROUP_COND_TYPES } from '@/dictionary/dynamic-group'

  const emit = defineEmits(['change'])

  const props = defineProps({
    height: {
      type: Number,
      default: 490
    },
    selected: {
      type: Array,
      default: () => ([])
    },
    disabledPropertyMap: {
      type: Object,
      default: () => ({})
    },
    models: {
      type: Array,
      default: () => ([])
    },
    propertyMap: {
      type: [Object, Array],
      default: () => ({})
    },
    conditionType: {
      type: String,
      default: DYNAMIC_GROUP_COND_TYPES.IMMUTABLE // condition: 锁定条件 varCondition：可变条件
    },
  })
  const checkboxRef = ref('')
  const indeterminate = reactive({})
  const allChecked = reactive({})
  const dataEmpty = reactive({
    type: 'empty',
    payload: {
      defaultText: t('暂无数据')
    }
  })

  const disabledPropertyCounts = reactive({})
  const matchedPropertyMap = ref(props.propertyMap)
  const localSelected = ref([...props.selected])
  const filter = ref('')

  const propertyMap = computed(() => props.propertyMap)
  const isShowEmpty = computed(() => {
    let isNoData = true
    Object.values(matchedPropertyMap.value)?.forEach((value) => {
      if (value?.length > 0) isNoData = false
    })
    return isNoData
  })

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

  const isDisabled = (model, property) => props.disabledPropertyMap[model.bk_obj_id].includes(property.bk_property_id)

  const getLength = (bkObjId) => {
    const length = matchedPropertyMap.value[bkObjId]?.length || 0
    const disabledLength = disabledPropertyCounts[bkObjId] || 0
    return { length, disabledLength }
  }

  const getCheckDisabled = (bkObjId) => {
    const { length, disabledLength } = getLength(bkObjId)
    if (length === disabledLength) return true
    return false
  }

  const getDisabledTip = (property) => {
    const type = property?.conditionType ?? props.conditionType
    if (type !== props.conditionType) {
      return t('条件已被添加', {
        condition: t(DYNAMIC_GROUP_COND_NAMES[type])
      })
    }
    return t('该字段不支持配置')
  }

  const updateLocalSelected = (property, checked) => {
    property.conditionType = props.conditionType
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
    emit('change')
  }

  const handleChangeAllCheck = (bkObjId, checked) => {
    indeterminate[bkObjId] = false
    allChecked[bkObjId] = checked
    matchedPropertyMap.value?.[bkObjId]?.forEach((target) => {
      const isDisabled = props.disabledPropertyMap[bkObjId].includes(target.bk_property_id)
      if (!isDisabled) {
        updateLocalSelected(target, checked)
      }
    })
    emit('change')
  }

  // 判断相应的全选/半选状态
  const allCheckState = ({ bk_obj_id: bkObjId }) => {
    const { length, disabledLength } = getLength(bkObjId)
    if (length === 0) return
    const matchedPropertyMapIdSet = new Set()
    matchedPropertyMap.value[bkObjId]?.forEach(property => matchedPropertyMapIdSet.add(property?.id))
    const currentCheckedCount = localSelected.value.filter(target => target.bk_obj_id === bkObjId
      && matchedPropertyMapIdSet.has(target.id)
      && (target?.conditionType ?? DYNAMIC_GROUP_COND_TYPES.IMMUTABLE) === props.conditionType)
      ?.length || 0
    // 默认是一个都没选的状态
    let isIndeterminate = false // 半选
    let isChecked = false

    if (currentCheckedCount > 0) {
      if (currentCheckedCount === length - disabledLength) {
        isChecked = true
      } else {
        isIndeterminate = true
      }
    }
    indeterminate[bkObjId] = isIndeterminate
    allChecked[bkObjId] = isChecked
  }

  const handleClearFilter = () => filter.value = ''

  const initChecked = () => {
    props.models.forEach((model) => {
      const objId = model?.bk_obj_id
      if (objId) {
        allCheckState({ bk_obj_id: objId })
      }
    })
  }

  const initDisabledProperty = () => {
    Object.keys(matchedPropertyMap.value)?.forEach((bkObjId) => {
      let length = 0
      matchedPropertyMap.value[bkObjId]?.forEach((target) => {
        const isDisabled = props.disabledPropertyMap[bkObjId].includes(target.bk_property_id)
        if (isDisabled) {
          length += 1
        }
      })
      disabledPropertyCounts[bkObjId] = length
    })
  }

  initDisabledProperty()
  initChecked()

  watch(() => filter.value, (filter) => {
    handleFilter(filter)
    dataEmpty.type = filter ? 'search' : 'empty'
  }, {
    immediate: true
  })

  defineExpose({
    localSelected
  })

</script>

<style lang="scss" scoped>
.property-selector-content {
  width: 400px;
  max-height: 500px;
  padding: 10px 14px;
  margin: -.3rem -.6rem;
}
.property-selector-container {
  max-height: calc(100% - 32px);
  margin-right: -14px;
  margin-left: -14px;
  padding: 0 14px;
  @include scrollbar-y;
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

  .all-check {
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

      &.is-checked-diabled {
        background: #f9fafd;
        :deep(.bk-checkbox-text),
        .icon-cc-selected {
          color: #dcdee5;
        }
      }

      :deep {
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
