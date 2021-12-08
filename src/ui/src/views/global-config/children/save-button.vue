<template>
  <bk-popconfirm trigger="click" :disabled="loading" :title="$t('确认保存配置？')"
    :content="$t('保存后将会立即生效，请谨慎操作！')" @confirm="handleConfirm">
    <bk-button theme="primary" class="action-button" :loading="loading">{{$t('保存')}}</bk-button>
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
      const handleConfirm = () => emit('save')

      return { handleConfirm }
    }
  }
</script>

<style lang="scss" scoped>
@import url("../style.scss");
</style>
