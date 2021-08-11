<template>
  <bk-date-picker
    type="datetimerange"
    transfer
    :value="localValue"
    v-bind="$attrs"
    format="yyyy-MM-dd HH:mm:ss"
    @clear="() => $emit('clear')"
    @change="handleChange"
    @open-change="handleToggle">
  </bk-date-picker>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-time',
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
          return [...this.value]
        },
        set(values) {
          this.$emit('input', values)
          this.$emit('change', values)
        }
      }
    },
    methods: {
      handleChange(values) {
        if (values.toString() === this.value.toString()) return
        this.localValue = values.filter(value => !!value)
      }
    }
  }
</script>
