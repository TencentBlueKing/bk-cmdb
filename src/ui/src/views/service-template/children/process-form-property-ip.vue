<template>
  <bk-select
    v-model="localValue"
    v-bind="$attrs"
    :popover-min-width="120"
    :clearable="false">
    <bk-option v-for="(name, id) in IPMap" :key="id" :id="id" :name="name"></bk-option>
  </bk-select>
</template>

<script>
  import { PROCESS_BIND_IP_ALL_MAP, PROCESS_BIND_IPV4_MAP, PROCESS_BIND_IPV6_MAP } from '@/dictionary/process-bind-ip.js'

  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      type: {
        type: String,
        default: ''
      }
    },
    computed: {
      IPMap() {
        if (this.type === 'v4') {
          return PROCESS_BIND_IPV4_MAP
        }
        if (this.type === 'v6') {
          return PROCESS_BIND_IPV6_MAP
        }
        return PROCESS_BIND_IP_ALL_MAP
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
