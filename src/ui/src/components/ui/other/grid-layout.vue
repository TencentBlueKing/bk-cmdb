<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

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
      maxColumns: Number
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
      '--min-height': rowMinHeight,
      '--max-columns': maxColumns
    }">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.cmdb-grid-layout {
  display: grid;
  gap: var(--gap, 24px);
  grid-template-columns: repeat(var(--max-columns, auto-fill), minmax(var(--min-width, 200px), var(--max-width, 1fr)));
  grid-auto-rows: minmax(var(--min-height, 32px), auto);
}
</style>
