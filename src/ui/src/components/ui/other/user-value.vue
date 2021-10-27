<template>
  <blueking-user-selector type="info"
    v-if="localValue.length"
    style="font-size: 12px;"
    :api="api"
    :value="localValue">
  </blueking-user-selector>
  <span v-else>--</span>
</template>

<script>
  import BluekingUserSelector from '@blueking/user-selector'
  export default {
    components: {
      BluekingUserSelector
    },
    props: {
      value: {
        type: String,
        default: ''
      }
    },
    data() {
      return {}
    },
    computed: {
      api() {
        const { userManage } = window.ESB
        if (userManage) {
          try {
            const url = new URL(userManage)
            return `${window.API_HOST}proxy/get/usermanage${url.pathname}`
          } catch (e) {
            console.error(e)
          }
        }
        return ''
      },
      localValue: {
        get() {
          if (this.value) {
            return this.value.split(',')
          }
          return []
        }
      }
    }
  }
</script>
