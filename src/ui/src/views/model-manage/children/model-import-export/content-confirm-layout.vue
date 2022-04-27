<template>
  <div class="content-confirm-layout">
    <div class="content-confirm-aside">
      <slot name="aside"></slot>
    </div>
    <div class="content-confirm-main">
      <slot name="main"></slot>
    </div>
  </div>
</template>

<script>
  import { defineComponent } from '@vue/composition-api'

  export default defineComponent({
    name: 'ContentConfirmLayout'
  })
</script>

<style lang="scss" scoped>
$borderColor: #dcdee5;
.content-confirm {
  &-layout {
    display: flex;
    height: 100%;
  }

  &-aside {
    position: relative;
    flex: 0 0 360px;
    width: 360px;
    overflow-y: auto;
    overflow-x: visible;
    background: #fafbfd;

    &::after {
      content:"";
      display: block;
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      width: 1px;
      background-color: $borderColor;
    }
  }

  &-main {
    flex: 1 1 50%;
    margin-left: 25px;
    height: 100%;
    overflow-y: auto;
  }
}
</style>
