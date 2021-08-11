<template>
  <bk-select
    v-model="localValue"
    v-bind="$attrs"
    @toggle="handleToggle"
    @clear="() => $emit('clear')">
    <bk-option id="true" name="true"></bk-option>
    <bk-option id="false" name="false"></bk-option>
  </bk-select>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-bool',
    mixins: [activeMixin],
    props: {
      value: {
        type: [String, Boolean],
        default: ''
      }
    },
    computed: {
      localValue: {
        get() {
          return this.value.toString()
        },
        set(value) {
          const realValue = value && value === 'true'
          this.$emit('input', realValue)
          this.$emit('change', realValue)
        }
      }
    }
  }
</script>
