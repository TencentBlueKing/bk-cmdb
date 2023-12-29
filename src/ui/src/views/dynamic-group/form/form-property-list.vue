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
  <div class="form-property-list">
    <bk-form form-type="vertical" :label-width="400">
      <bk-form-item
        v-for="property in properties"
        :key="property.id"
        :label="getPropertyLabel(property)">
        <div class="form-property-item">
          <form-operator-selector class="item-operator"
            v-if="!withoutOperator.includes(property.bk_property_type)"
            :property="property"
            :custom-type-map="customTypeMap"
            v-model="condition[property.id].operator"
            :disabled="disabled"
            @selected="handleOperatorChange(property, ...arguments)">
          </form-operator-selector>
          <div class="item-value">
            <component
              class="form-element"
              :is="getComponentType(property)"
              :placeholder="getPlaceholder(property)"
              :data-vv-name="property.bk_property_id"
              :data-vv-as="property.bk_property_name"
              v-bind="getBindProps(property)"
              v-model="condition[property.id].value"
              display-tag
              :disabled="disabled"
              v-validate="'required'"
              v-bk-tooltips.top="{
                disabled: !property.placeholder,
                theme: 'light',
                trigger: 'click',
                content: property.placeholder
              }">
            </component>
          </div>
          <i class="item-remove bk-icon icon-close" v-if="!disabled" @click="handleRemove(property)"></i>
        </div>
        <p class="form-error" v-if="errors.has(property.bk_property_id)">{{errors.first(property.bk_property_id)}}</p>
      </bk-form-item>
    </bk-form>
  </div>
</template>

<script>
  import FormOperatorSelector from '@/components/filters/operator-selector.vue'
  import has from 'has'
  import { QUERY_OPERATOR } from '@/utils/query-builder-operator'

  export default {
    components: {
      FormOperatorSelector
    },
    inject: ['dynamicGroupForm'],
    props: {
      disabled: {
        type: Boolean,
        value: false
      }
    },
    data() {
      const { EQ, NE, GTE, LTE, RANGE } = QUERY_OPERATOR
      return {
        condition: {},
        withoutOperator: ['date', 'time', 'bool', 'service-template'],
        customTypeMap: {
          float: [EQ, NE, GTE, LTE, RANGE],
          int: [EQ, NE, GTE, LTE, RANGE]
        }
      }
    },
    computed: {
      bizId() {
        return this.dynamicGroupForm.bizId
      },
      properties() {
        return this.dynamicGroupForm.selectedProperties
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
      }
    },
    watch: {
      properties: {
        immediate: true,
        handler() {
          this.updateCondition()
        }
      }
    },
    methods: {
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
        Object.values(this.condition).forEach((condition) => {
          const modelId = condition.property.bk_obj_id
          const propertyId = condition.property.bk_property_id
          // eslint-disable-next-line max-len
          const detailsCondition = this.details.info.condition.find(detailsCondition => detailsCondition.bk_obj_id === modelId)
          const detailsFieldData = detailsCondition.condition.find(data => data.field === propertyId)
          condition.operator = detailsFieldData.operator
          condition.value = detailsFieldData.value
        })
      },
      updateCondition() {
        const newConditon = {}
        this.properties.forEach((property) => {
          if (has(this.condition, property.id)) {
            newConditon[property.id] = this.condition[property.id]
          } else {
            newConditon[property.id] = {
              property,
              ...this.getDefaultData(property)
            }
          }
        })
        this.condition = newConditon
      },
      handleOperatorChange(property, operator) {
        if (operator === '$range') {
          this.condition[property.id].value = []
        } else if (operator === '$regex') {
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
      getComponentType(property) {
        const {
          bk_obj_id: modelId,
          bk_property_id: propertyId,
          bk_property_type: propertyType
        } = property
        const isSetName = modelId === 'set' && propertyId === 'bk_set_name'
        const isModuleName = modelId === 'module' && propertyId === 'bk_module_name'

        if ((isSetName || isModuleName) && this.condition[property.id].operator !== '$regex') {
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
        const selectTypes = ['list', 'enum', 'timezone', 'organization', 'date', 'time', 'enumquote', 'enummulti']
        if (selectTypes.includes(property.bk_property_type)) {
          return this.$t('请选择xx', { name: property.bk_property_name })
        }
        return this.$t('请输入xx', { name: property.bk_property_name })
      }
    }
  }
</script>

<style lang="scss" scoped>
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

      .item-remove {
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
    align-items: center;
    .item-operator {
      flex: 110px 0 0;
      margin-right: 10px;
    }
    .item-value {
      flex: 1;
      margin: 0 10px 0 0;
      width: calc(100% - 150px);
      display: flex;
      align-items: center;

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
  }
  .form-error {
    font-size: 12px;
    line-height: 14px;
    color: $dangerColor;
  }
}
</style>
