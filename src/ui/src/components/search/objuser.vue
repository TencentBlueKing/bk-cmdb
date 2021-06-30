<template>
  <cmdb-form-objuser
    v-model="localValue"
    v-bind="$attrs"
    @clear="() => $emit('clear')"
    @focus="handleToggle(true)"
    @blur="handleToggle(false)">
  </cmdb-form-objuser>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-objuser',
    mixins: [activeMixin],
    props: {
      value: {
        type: Array,
        default: () => ([])
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value.join(',')
        },
        set(value) {
          const values = value.split(',')
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    }
  }
</script>
