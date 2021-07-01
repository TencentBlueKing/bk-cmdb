<template>
  <bk-select
    searchable
    v-model="localValue"
    v-bind="$attrs"
    :multiple="multiple"
    @clear="() => $emit('clear')"
    @toggle="handleToggle">
    <bk-option v-for="option in options"
      :key="option.id"
      :id="option.id"
      :name="option.name">
    </bk-option>
  </bk-select>
</template>

<script>
  import activeMixin from './mixins/active'
  export default {
    name: 'cmdb-search-enum',
    mixins: [activeMixin],
    props: {
      value: {
        type: [String, Array],
        default: () => ([])
      },
      options: {
        type: Array,
        default: () => ([])
      }
    },
    computed: {
      multiple() {
        return Array.isArray(this.value)
      },
      localValue: {
        get() {
          return this.value
        },
        set(value) {
          this.$emit('input', value)
          this.$emit('change', value)
        }
      }
    }
  }
</script>
