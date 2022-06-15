<script lang="ts">
  import { computed, defineComponent } from '@vue/composition-api'

  export default defineComponent({
    props: {
      mode: {
        type: String,
        default: 'detail'
      },
      gap: {
        type: [Number, String]
      },
      minWidth: {
        type: [Number, String]
      },
      maxWidth: {
        type: [Number, String]
      },
      minHeight: {
        type: [Number, String]
      },
    },
    provide() {
      return {
        mode: this.mode
      }
    },
    setup(props) {
      const gridGap = computed(() => (isNaN(Number(props.gap)) ? props.gap : `${props.gap}px`))

      const colMinWidth = computed(() => (isNaN(Number(props.minWidth)) ? props.minWidth : `${props.minWidth}px`))
      const colMaxWidth = computed(() => (isNaN(Number(props.maxWidth)) ? props.maxWidth : `${props.maxWidth}px`))

      const rowMinHeight = computed(() => (isNaN(Number(props.minHeight)) ? props.minHeight : `${props.minHeight}px`))

      return {
        gridGap,
        colMinWidth,
        colMaxWidth,
        rowMinHeight
      }
    }
  })
</script>

<template>
  <div
    class="cmdb-grid-layout"
    :style="{
      '--gap': gridGap,
      '--min-width': colMinWidth,
      '--max-width': colMaxWidth,
      '--min-height': rowMinHeight
    }">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
  .cmdb-grid-layout {
    display: grid;
    gap: var(--gap, 24px);
    grid-template-columns: repeat(auto-fill, minmax(var(--min-width, 200px), var(--max-width, 1fr)));
    grid-auto-rows: minmax(var(--min-height, 32px), auto);
  }
</style>
