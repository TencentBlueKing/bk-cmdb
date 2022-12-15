<script>
  import { defineComponent, ref } from 'vue'

  export default defineComponent({
    setup() {
      const topHeight = ref(86)

      return {
        topHeight
      }
    }
  })
</script>

<template>
  <div class="details-general-layout" :style="{ '--topHeight': `${topHeight}px` }">
    <div class="details-top">
      <slot name="top"></slot>
    </div>
    <div class="details-main">
      <slot name="main"></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.details-general-layout {
  overflow: hidden;

  .details-top {
    max-height: 450px;
    padding: 12px 0 12px 24px;
    background:rgba(235, 244, 255, .6);
    border-bottom: 1px solid #dcdee5;
    @include scrollbar-y;
  }

  .details-main {
    height: calc(100% - var(--topHeight)) !important;
    min-height: 400px;
    .details-tab {
      height: 100%;
      ::v-deep {
        .bk-tab-header {
          padding: 0;
          margin: 0 20px;
        }
        .bk-tab-section {
          @include scrollbar-y;
          padding-bottom: 10px;
        }
      }
    }
  }
}
</style>
