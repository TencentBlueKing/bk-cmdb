<template>
  <user-value :value="value" v-if="isUser"></user-value>
  <table-value
    :value="value"
    :show-on="showOn"
    :format-cell-value="formatCellValue"
    :property="property"
    v-else-if="isTable">
  </table-value>
  <service-template-value
    v-else-if="isServiceTemplate"
    :value="value"
    display-type="info">
  </service-template-value>
  <compmoent :is="tag" v-bind="attrs" v-else>{{displayValue}}</compmoent>
</template>

<script>
  import UserValue from './user-value'
  import TableValue from './table-value'
  import ServiceTemplateValue from '@/components/search/service-template'
  const ORG_CACHES = {}
  export default {
    name: 'cmdb-property-value',
    components: {
      UserValue,
      TableValue,
      ServiceTemplateValue
    },
    props: {
      value: {
        type: [String, Number, Array, Boolean],
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
      multiple: Boolean
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
        if (this.multiple && Array.isArray(value)) {
          displayQueue = value.map(subValue => this.getDisplayValue(subValue))
        } else {
          displayQueue = [this.getDisplayValue(value)]
        }
        const result = await Promise.all(displayQueue)
        this.displayValue = result.join(',')
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
        const cacheKey = (value || []).join('_')
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
      }
    }
  }
</script>

<style lang="scss" scoped>
    .value-primary-theme {
        color: $primaryColor;
        cursor: pointer;
    }
</style>
