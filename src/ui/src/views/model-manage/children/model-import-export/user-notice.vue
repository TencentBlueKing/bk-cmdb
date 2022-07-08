<template>
  <div class="user-notice">
    <div class="notice-icon">
      <img alt="Notice" src="@/assets/images/notice-pane.png" />
    </div>
    <h2 class="notice-title">
      {{ t("用户须知") }}
    </h2>
    <div class="notice-content">
      <slot></slot>
    </div>
  </div>
</template>

<script>
  import { defineComponent } from '@vue/composition-api'
  import { t } from '@/i18n'

  export default defineComponent({
    name: 'UserNotice',
    setup() {
      return {
        t
      }
    }
  })
</script>

<style lang="scss" scoped>
.notice {
  &-icon {
    margin: 30px auto 0;
    width: 158px;
    height: 106px;

    img {
      width: 100%;
      height: 100%;
    }
  }

  &-title {
    text-align: center;
    margin-top: 20px;
    font-weight: 400;
    font-size: 24px;
  }

  &-content {
    width: 680px;
    margin: 20px auto 0;
    padding: 16px 16px 26px;
    font-size: 12px;
    line-height: 20px;
    background: #fafbfd;
    border: 1px solid #dcdee5;
    border-radius: 2px;
  }
}
</style>
