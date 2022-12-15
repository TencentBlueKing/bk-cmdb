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
    v-if="isUser">
  </user-value>
  <table-value
    ref="complexTypeComp"
    :value="value"
    :show-on="showOn"
    :format-cell-value="formatCellValue"
    :property="property"
    v-else-if="isTable">
  </table-value>
  <service-template-value
    v-else-if="isServiceTemplate"
    ref="complexTypeComp"
    :value="value"
    display-type="info">
  </service-template-value>
  <mapstring-value
    v-else-if="isMapstring"
    ref="complexTypeComp"
    :value="value">
  </mapstring-value>
  <component
    :is="tag"
    v-bind="attrs"
    v-bk-overflow-tips
    v-else-if="isShowOverflowTips">
    {{displayValue}}
  </component>
  <component
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
  const ORG_CACHES = {}
  export default {
    name: 'cmdb-property-value',
    components: {
      UserValue,
      TableValue,
      ServiceTemplateValue,
      MapstringValue
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
      isShowOverflowTips: Boolean
    },
    data() {
      return {
        displayValue: ''
      }
    },
    computed: {
      attrs() {
        const attrs = {
          class: `value-${this.theme}-theme`
        }
        return attrs
      },
      isUser() {
        const type = typeof this.property === 'object' ? this.property.bk_property_type : this.property
        return type === 'objuser'
      },
      isTable() {
        return this.property.bk_property_type === 'table'
      },
      isServiceTemplate() {
        return this.property.bk_property_type === 'service-template'
      },
      isOrg() {
        return this.property.bk_property_type === 'organization'
      },
      isMapstring() {
        return this.property.bk_property_type === 'map'
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
        if (this.isUser || this.isTable) return
        let displayQueue
        if (this.multiple && Array.isArray(value) && !this.isOrg) {
          displayQueue = value.map(subValue => this.getDisplayValue(subValue))
        } else {
          displayQueue = [this.getDisplayValue(value)]
        }
        const result = await Promise.all(displayQueue)
        this.displayValue = result.join(', ')
      },
      async getDisplayValue(value) {
        let displayValue
        const isPropertyObject = Object.prototype.toString.call(this.property) === '[object Object]'
        const type = isPropertyObject ? this.property.bk_property_type : this.property
        const unit = isPropertyObject ? this.property.unit : ''
        if (type === 'organization') {
          displayValue = await this.getOrganization(value)
        } else {
          displayValue = this.$options.filters.formatter(value, this.property, this.options)
        }
        // eslint-disable-next-line no-nested-ternary
        return (this.showUnit && unit && displayValue !== '--')
          ? `${displayValue}${unit}`
          : String(displayValue).length
            ? displayValue
            : '--'
      },
      async getOrganization(value) {
        let displayValue
        const cacheKey = Array.isArray(value) ? value.join('_') : String(value)
        if (ORG_CACHES[cacheKey]) {
          return ORG_CACHES[cacheKey]
        }

        if (!value || !value.length) {
          displayValue = '--'
        } else {
          const res = await this.$store.dispatch('organization/getDepartment', {
            params: {
              lookup_field: 'id',
              exact_lookups: value.join(',')
            },
            config: {
              fromCache: true,
              requestId: `get_department_id_${cacheKey}`
            }
          })
          const names = (res.results || []).map(item => item.full_name)
          displayValue = names.join('; ') || '--'
        }

        ORG_CACHES[cacheKey] = displayValue
        return displayValue
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
  .value-primary-theme {
    color: $primaryColor;
    cursor: pointer;
    display: block;
  }
</style>
