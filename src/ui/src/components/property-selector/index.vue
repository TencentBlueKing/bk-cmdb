<template>
  <bk-select
    v-model="localValue"
    :searchable="searchable"
    font-size="medium"
    :clearable="false"
    v-bind="$attrs">
    <bk-option v-for="option in options"
      :key="option.bk_property_id"
      :id="option.bk_property_id"
      :name="option.bk_property_name">
    </bk-option>
  </bk-select>
</template>

<script>
  export default {
    name: 'cmdb-property-selector',
    props: {
      properties: {
        type: Array,
        default: () => ([])
      },
      value: {
        type: [String, Number],
        default: ''
      },
      searchable: {
        type: Boolean,
        default: true
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      },
      options() {
        return this.properties.filter(property => !!property.id)
      }
    }
  }
</script>
