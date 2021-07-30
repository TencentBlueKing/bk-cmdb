<template>
  <bk-date-picker
    type="daterange"
    transfer
    :value="localValue"
    v-bind="$attrs"
    format="yyyy-MM-dd"
    @change="handleChange"
    @open-change="handleToggle"
    @clear="() => $emit('clear')">
  </bk-date-picker>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-date',
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
