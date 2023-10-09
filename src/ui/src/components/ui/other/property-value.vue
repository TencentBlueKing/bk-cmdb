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
  <user-value
    :value="value"
    ref="complexTypeComp"
    v-if="property.bk_property_type === PROPERTY_TYPES.OBJUSER">
  </user-value>
  <table-value
    ref="complexTypeComp"
    :value="value"
    :show-on="showOn"
    :format-cell-value="formatCellValue"
    :property="property"
    v-else-if="property.bk_property_type === PROPERTY_TYPES.TABLE">
  </table-value>
  <service-template-value
    v-else-if="property.bk_property_type === PROPERTY_TYPES.SERVICE_TEMPLATE"
    ref="complexTypeComp"
    :value="value"
    display-type="info">
  </service-template-value>
  <mapstring-value
    v-else-if="property.bk_property_type === PROPERTY_TYPES.MAP"
    ref="complexTypeComp"
    :value="value">
  </mapstring-value>
  <enumquote-value
    v-else-if="property.bk_property_type === PROPERTY_TYPES.ENUMQUOTE"
    ref="complexTypeComp"
    :value="value"
    :property="property">
  </enumquote-value>
  <org-value
    v-else-if="property.bk_property_type === PROPERTY_TYPES.ORGANIZATION"
    ref="complexTypeComp"
    :value="value"
    :property="property"
    :show-on="showOn"
    v-bind="$attrs">
  </org-value>
  <inner-table-value
    v-else-if="property.bk_property_type === PROPERTY_TYPES.INNER_TABLE"
    ref="complexTypeComp"
    :value="value"
    :property="property"
    :show-on="showOn"
    :instance="instance"
    v-bind="$attrs">
  </inner-table-value>
  <component
    class="value-container"
    :is="tag"
    v-bind="attrs"
    v-bk-overflow-tips
    v-else-if="isShowOverflowTips">
    {{displayValue}}
  </component>
  <component
    class="value-container"
    :is="tag"
    v-bind="attrs"
    v-else>
    {{displayValue}}
  </component>
</template>

<script>
  import UserValue from './user-value'
  import TableValue from './table-value'
  import ServiceTemplateValue from '@/components/search/service-template'
  import MapstringValue from './mapstring-value.vue'
  import EnumquoteValue from './enumquote-value.vue'
  import OrgValue from './org-value.vue'
  import InnerTableValue from './inner-table-value.vue'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { isUseComplexValueType } from '@/utils/tools'

  export default {
    name: 'cmdb-property-value',
    components: {
      UserValue,
      TableValue,
      ServiceTemplateValue,
      MapstringValue,
      EnumquoteValue,
      OrgValue,
      InnerTableValue
    },
    props: {
      value: {
        type: [String, Number, Array, Boolean, Object],
        default: ''
      },
      property: {
        type: [Object, String],
        default: () => ({})
      },
      options: {
        type: [Array, String, Object],
        default: () => ([])
      },
      showUnit: {
        type: Boolean,
        default: true
      },
      tag: {
        type: String,
        default: 'span'
      },
      className: {
        type: String,
        default: ''
      },
      theme: {
        type: String,
        default: 'default',
        validator(value) {
          return ['primary', 'default'].includes(value)
        }
      },
      showOn: {
        type: String,
        default: 'default',
        validator(value) {
          return ['default', 'cell'].includes(value)
        }
      },
      formatCellValue: Function,
      multiple: Boolean,
      isShowOverflowTips: Boolean,
      instance: {
        type: Object,
        default: () => ({})
      }
    },
    data() {
      return {
        displayValue: '',
        PROPERTY_TYPES
      }
    },
    computed: {
      attrs() {
        const attrs = {
          class: `value-${this.theme}-theme`
        }
        return attrs
      }
    },
    watch: {
      value(value) {
        this.setDisplayValue(value)
      }
    },
    created() {
      this.setDisplayValue(this.value)
    },
    methods: {
      async setDisplayValue(value) {
        if (isUseComplexValueType(this.property)) {
          return
        }

        let displayQueue
        if (this.multiple && Array.isArray(value)) {
          displayQueue = value.map(subValue => this.getDisplayValue(subValue))
        } else {
          displayQueue = [this.getDisplayValue(value)]
        }
        const result = await Promise.all(displayQueue)
        this.displayValue = result.join(', ')
      },
      async getDisplayValue(value) {
        const unit = this.property.unit || ''
        const displayValue = this.$options.filters.formatter(value, this.property, this.options)

        if ((this.showUnit && unit && displayValue !== '--')) {
          return `${displayValue}${unit}`
        }

        return String(displayValue).length ? displayValue : '--'
      },
      getCopyValue() {
        if (this.$refs?.complexTypeComp) {
          return this.$refs?.complexTypeComp?.getCopyValue?.()
        }
        return this.displayValue
      }
    }
  }
</script>

<style lang="scss" scoped>
  .value-container {
    display: block;
  }
  .value-primary-theme {
    color: $primaryColor;
    cursor: pointer;
    display: block;
  }
</style>
