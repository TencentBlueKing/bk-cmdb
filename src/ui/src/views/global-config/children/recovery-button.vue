<template>
  <!-- <bk-popconfirm trigger="click" :disabled="loading" @confirm="$listeners.confirm" :title="$t('确定恢复初始化？')"
    :content="$t('恢复初始化将会把所有选项内容恢复至初始化状态，不可恢复，请谨慎操作。')"> -->
  <bk-button size="small" :disabled="loading" @click="$listeners.confirm" text class="recovery-button">
    {{$t('恢复初始化')}}
  </bk-button>
  <!-- </bk-popconfirm> -->
</template>

<script>
  import { defineComponent } from '@vue/composition-api'
  export default defineComponent({
    props: {
      loading: {
        type: Boolean,
        default: false,
      },
    }
  })
</script>
<style lang="scss" scoped>
.recovery-button {
  padding: 0;
}
</style>
