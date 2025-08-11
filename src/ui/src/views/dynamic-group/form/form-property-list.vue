<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="form-property-list" ref="propertyList">
    <bk-form form-type="vertical" :label-width="400">
      <bk-form-item
        v-for="property in properties"
        :key="property.id"
        :label="getPropertyLabel(property)">
        <div class="form-property-item">
          <form-operator-selector class="item-operator"
            v-if="!withoutOperator.includes(property.bk_property_type)"
            :property="property"
            :custom-type-map="customOperatorTypeMap"
            :symbol-map="operatorSymbolMap"
            :desc-map="operatorDescMap"
            v-model="condition[property.id].operator"
            :disabled="disabled"
            @selected="handleOperatorChange(property, ...arguments)">
          </form-operator-selector>
          <div class="item-value">
            <component
              class="form-element"
              :is="getComponentType(property)"
              :placeholder="getPlaceholder(property)"
              :property="property"
              :data-vv-name="property.bk_property_id"
              :data-vv-as="property.bk_property_name"
              v-bind="getBindProps(property)"
              v-model="condition[property.id].value"
              display-tag
              :disabled="disabled"
              v-validate="'required'"
              @change="handleChange"
              @click.native="handleClick"
              :is-paste-split="isIPField(property.bk_property_id)"
              :popover-options="{
                duration: 0,
                onShown: handleShow,
                onHidden: handleHidden
              }"
              v-bk-tooltips.top="{
                disabled: !property.placeholder,
                theme: 'light',
                trigger: 'click',
                content: property.placeholder
              }"
              @inputchange="hanleInputChange">
            </component>
          </div>
          <i class="item-remove bk-icon icon-close" v-if="!disabled" @click="handleRemove(property)"></i>
          <i class="item-toggle bk-icon icon-sort" v-if="!disabled" @click="handleToggle(property)" v-bk-tooltips.top="{
            content: $t('切换为', {
              condition: conditionName
            })
          }"></i>
        </div>
        <p class="form-error" v-if="errors.has(property.bk_property_id)">{{errors.first(property.bk_property_id)}}</p>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  import FormOperatorSelector from '@/components/filters/operator-selector.vue'
  import has from 'has'
  import {
    QUERY_OPERATOR,
    QUERY_OPERATOR_OTHER_SYMBOL,
    QUERY_OPERATOR_HOST_SYMBOL,
    QUERY_OPERATOR_OTHER_DESC,
    QUERY_OPERATOR_HOST_DESC
  } from '@/utils/query-builder-operator'
  import { DYNAMIC_GROUP_COND_TYPES, DYNAMIC_GROUP_COND_NAMES } from '@/dictionary/dynamic-group'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'
  const { IMMUTABLE, VARIABLE } = DYNAMIC_GROUP_COND_TYPES

  export default {
    components: {
      FormOperatorSelector
    },
    inject: ['dynamicGroupForm'],
    props: {
      disabled: {
        type: Boolean,
        default: false
      },
      conditionType: {
        type: String,
        default: IMMUTABLE // condition: 锁定条件 varCondition：可变条件
      }
    },
    data() {
      return {
        condition: {},
        withoutOperator: ['date', 'time', 'bool', 'service-template']
      }
    },
    computed: {
      conditionName() {
        const exchangeType = {
          [VARIABLE]: IMMUTABLE,
          [IMMUTABLE]: VARIABLE
        }
        return this.$t(DYNAMIC_GROUP_COND_NAMES[exchangeType[this.conditionType]])
      },
      bizId() {
        return this.dynamicGroupForm.bizId
      },
      properties() {
        return this.dynamicGroupForm.selectedProperties
          .filter(property => property.conditionType === this.conditionType)
      },
      availableModels() {
        return this.dynamicGroupForm.availableModels
      },
      modelMap() {
        const modelMap = {}
        this.availableModels.forEach((model) => {
          modelMap[model.bk_obj_id] = model
        })
        return modelMap
      },
      details() {
        return this.dynamicGroupForm.details
      },
      searchObjId() {
        return this.dynamicGroupForm.formData.bk_obj_id
      },
      customOperatorTypeMap() {
        const { EQ, NE, GTE, LTE, RANGE, IN, NIN, LIKE, CONTAINS, CONTAINS_CS } = QUERY_OPERATOR
        const operatorTypeMap = {
          float: [EQ, NE, GTE, LTE, RANGE],
          int: [EQ, NE, GTE, LTE, RANGE]
        }
        if (this.searchObjId === BUILTIN_MODELS.HOST) {
          return {
            ...operatorTypeMap,
            longchar: [IN, NIN, CONTAINS, LIKE],
            singlechar: [IN, NIN, CONTAINS, LIKE],
            array: [IN, NIN, CONTAINS, LIKE],
            object: [IN, NIN, CONTAINS, LIKE]
          }
        }
        return {
          ...operatorTypeMap,
          longchar: [IN, NIN, CONTAINS, CONTAINS_CS],
          singlechar: [IN, NIN, CONTAINS, CONTAINS_CS],
          array: [IN, NIN, CONTAINS, CONTAINS_CS],
          object: [IN, NIN, CONTAINS, CONTAINS_CS]
        }
      },
      operatorSymbolMap() {
        if (this.searchObjId === BUILTIN_MODELS.HOST) {
          return QUERY_OPERATOR_HOST_SYMBOL
        }
        return QUERY_OPERATOR_OTHER_SYMBOL
      },
      operatorDescMap() {
        if (this.searchObjId === BUILTIN_MODELS.HOST) {
          return QUERY_OPERATOR_HOST_DESC
        }
        return QUERY_OPERATOR_OTHER_DESC
      }
    },
    watch: {
      properties: {
        immediate: true,
        handler() {
          this.updateCondition()
        }
      },
      condition: {
        handler(condition) {
          const { length } = Object.keys(condition)
          if (length) {
            Object.assign(this.dynamicGroupForm.storageCondition, condition)
          }
        },
        deep: true
      }
    },
    methods: {
      isIPField(id) {
        const IPField = ['bk_host_outerip_v6', 'bk_host_innerip_v6', 'bk_host_innerip', 'bk_host_outerip']
        return IPField.includes(id)
      },
      handleClick(e) {
        if (~e?.target?.className.indexOf('is-focus')) {
          // select专属
          return
        }
        this.calcPosition('click')
      },
      handleChange() {
        this.calcPosition()
      },
      hanleInputChange() {
        this.calcPosition()
      },
      handleShow() {
        this.calcPosition()
      },
      handleHidden() {
        this.$refs.propertyList.classList.remove('over-height')
      },
      calcPosition(type = 'change') {
        if (type === 'click') this.$refs.propertyList.classList.remove('over-height')
        if (!this.target) return

        this.$nextTick(() => {
          const parent = document.querySelector('.dynamic-group-form')
          const { scrollHeight } = parent
          const { height } = parent.getClientRects()[0]
          if (scrollHeight > Math.ceil(height)) {
            this.$refs.propertyList.classList.add('over-height')
          }
        })
      },
      getDefaultData(property) {
        const defaultMap = {
          bool: {
            operator: '$eq',
            value: ''
          },
          date: {
            operator: '$range',
            value: []
          },
          float: {
            operator: '$eq',
            value: ''
          },
          int: {
            operator: '$eq',
            value: ''
          },
          time: {
            operator: '$range',
            value: []
          },
          'service-template': {
            operator: '$in',
            value: []
          }
        }
        return {
          operator: '$in',
          value: [],
          ...defaultMap[property.bk_property_type]
        }
      },
      setDetailsCondition() {
        const { conditionType } = this
        Object.values(this.condition).forEach((condition) => {
          const modelId = condition.property.bk_obj_id
          const propertyId = condition.property.bk_property_id
          // eslint-disable-next-line max-len
          const detailsCondition = this.details.info[conditionType].find(detailsCondition => detailsCondition.bk_obj_id === modelId)
          const detailsFieldData = detailsCondition.condition.find(data => data.field === propertyId
            && data.conditionType === conditionType)
          condition.operator = detailsFieldData.operator
          condition.value = detailsFieldData.value
        })
      },
      updateCondition() {
        const newCondition = {}
        const { storageCondition } = this.dynamicGroupForm
        this.properties.forEach((property) => {
          if (has(storageCondition, property.id)) {
            // 修改conditionType为property.conditionType，保证conditionType正确
            storageCondition[property.id].property.conditionType = property.conditionType
            newCondition[property.id] = storageCondition[property.id]
          } else {
            newCondition[property.id] = {
              property,
              ...this.getDefaultData(property)
            }
          }
        })
        this.condition = newCondition
      },
      handleOperatorChange(property, operator) {
        if (operator === QUERY_OPERATOR.RANGE) {
          this.condition[property.id].value = []
        } else if ([QUERY_OPERATOR.LIKE, QUERY_OPERATOR.CONTAINS, QUERY_OPERATOR.CONTAINS_CS].includes(operator)) {
          const currentValue = this.condition[property.id].value
          this.condition[property.id].value = Array.isArray(currentValue) ? (currentValue[0] || '') : currentValue
        } else {
          const defaultValue = this.getDefaultData(property).value
          const currentValue = this.condition[property.id].value
          const isTypeChanged = (Array.isArray(defaultValue)) !== (Array.isArray(currentValue))
          this.condition[property.id].value = isTypeChanged ? defaultValue : currentValue
        }
      },
      handleRemove(property) {
        if (this.disabled) return
        this.$emit('remove', property)
      },
      handleToggle(property) {
        if (this.disabled) return
        this.$emit('toggle', property, this.conditionType)
      },
      getComponentType(property) {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = property
        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'

        if ((isSetName || isModuleName)
          && ![QUERY_OPERATOR.CONTAINS, QUERY_OPERATOR.LIKE].includes(this.condition[property.id].operator)) {
          return `cmdb-search-${modelId}`
        }

        return `cmdb-search-${propertyType}`
      },
      getBindProps(property) {
        const props = this.getNormalProps(property)
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId
        } = property
        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'
        if (isSetName || isModuleName) {
          return Object.assign(props, { bizId: this.bizId })
        }
        return props
      },
      getNormalProps(property) {
        const type = property.bk_property_type
        if (['list', 'enum', 'enumquote', 'enummulti'].includes(type)) {
          return {
            options: property.option || []
          }
        }
        if (type === 'objuser') {
          return {
            fastSelect: true
          }
        }
        return {}
      },
      getPropertyLabel(property) {
        const modelId = property.bk_obj_id
        const propertyName = property.bk_property_name
        const modelName = this.modelMap[modelId].bk_obj_name
        return `${modelName} - ${propertyName}`
      },
      getPlaceholder(property) {
        const selectTypes = ['list', 'enum', 'timezone', 'organization', 'date', 'time', 'enumquote', 'enummulti', 'foreignkey']
        const name = this.$t(this.conditionType === IMMUTABLE ? '条件值' : '默认值')
        if (selectTypes.includes(property.bk_property_type)) {
          return this.$t('请选择xx', { name })
        }
        return this.$t('请输入xx', { name })
      }
    }
  }
</script>

<style lang="scss" scoped>
.over-height {
  .g-expand {
    bottom: -32px;
  }
}
.form-property-list {
  /deep/ .bk-form-item {
    padding: 8px;
    margin: -8px;
    margin-bottom: 4px !important;

    .bk-label {
      cursor: pointer;
      .bk-label-text {
        width: calc(100% - 20px);
        @include ellipsis;
      }
    }

    &:hover {
      background: #F0F1F5;

      .item-remove, .item-toggle {
        visibility: visible;
      }
    }
  }
  :deep(.bk-select-tag-container.is-focus) {
    max-height: 200px;
    @include scrollbar;
  }
  .form-property-item {
    display: flex;
    align-items: flex-start;
    .item-operator {
      flex: 128px 0 0;
      margin-right: 10px;
    }
    .item-value {
      flex: 1;
      display: flex;
      align-items: self-start;
      position: relative;

      .form-element {
        width: 100%;
      }
    }
    .item-remove {
      font-size: 20px;
      visibility: hidden;
      cursor: pointer;
      position: absolute;
      right: 0;
      top: -32px;
      &:hover {
        color: #EA3636;
      }
    }
    .item-toggle {
      font-size: 14px;
      visibility: hidden;
      cursor: pointer;
      position: absolute;
      right: 20px;
      top: -29px;
      transform: rotate(90deg);
      &:hover {
        color: #EA3636;
      }
    }
  }
  .form-error {
    font-size: 12px;
    line-height: 14px;
    color: $dangerColor;
  }
}
</style>
