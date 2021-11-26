<template>
  <bk-popconfirm trigger="click" :disabled="loading" :title="$t('确认重置所有选项？')" @confirm="handleConfirm"
    :content="$t('重置将会把所有选项内容恢复至初始化状态，不可恢复，请谨慎操作。')">
    <bk-button class="action-button" :loading="loading">{{$t('重置')}}</bk-button>
  </bk-popconfirm>
</template>

<script>
  export default {
    props: {
      loading: {
        type: Boolean,
        default: false
      }
    },
    setup(props, { emit }) {
      const handleConfirm = () => emit('reset')

      return { handleConfirm }
    }
  }
</script>

<style lang="scss" scoped>
@import url("../style.scss");
</style>
