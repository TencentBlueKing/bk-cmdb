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

<template>
  <div
    v-show="showWrapper"
    v-transfer-dom
    class="dialog-wrapper"
    :style="{
      zIndex,
    }"
    @click.stop.prevent="closeDialog">
    <div ref="resizeTrigger">
      <transition name="dialog-fade">
        <div
          v-if="showBody"
          ref="body"
          class="dialog-body"
          :class="{ 'is-scrollable': bodyScroll }"
          :style="bodyStyle">
          <div v-if="showHeader" ref="header" class="dialog-header">
            <slot name="header"></slot>
          </div>
          <div ref="dialogContent" class="dialog-content">
            <slot></slot>
          </div>
          <div v-if="showFooter" ref="footer" class="dialog-footer">
            <slot name="footer"></slot>
          </div>
          <i
            v-if="showCloseIcon"
            class="bk-icon icon-close"
            @click="handleCloseDialog"></i>
        </div>
      </transition>
    </div>
  </div>
</template>

<script>
import {
  addResizeListener,
  removeResizeListener,
} from '@/utils/resize-events.js'
export default {
  name: 'cmdb-dialog',
  props: {
    value: Boolean,
    showHeader: {
      type: Boolean,
      default: true,
    },
    showFooter: {
      type: Boolean,
      default: true,
    },
    showCloseIcon: {
      type: Boolean,
      default: true,
    },
    width: {
      type: Number,
      default: 720,
    },
    height: Number,
    minHeight: Number,
    bodyScroll: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      timer: null,
      bodyHeight: 0,
      showWrapper: false,
      showBody: false,
      zIndex: 2000,
    }
  },
  computed: {
    autoResize() {
      return typeof this.height !== 'number'
    },
    bodyStyle() {
      const style = {
        width: `${this.width}px`,
        '--height': `${this.autoResize ? this.bodyHeight : this.height}px`,
      }
      if (!this.autoResize) {
        style.height = `${this.height}px`
        style.maxHeight = 'initial'
      }
      if (this.minHeight) {
        style.minHeight = `${this.minHeight}px`
      }
      return style
    },
  },
  watch: {
    value: {
      immediate: true,
      handler(value) {
        if (value) {
          this.showWrapper = true
          // eslint-disable-next-line no-underscore-dangle
          this.zIndex = window.__bk_zIndex_manager.nextZIndex()
          this.$nextTick(() => {
            this.showBody = true
          })
        } else {
          this.showBody = false
          this.timer && clearTimeout(this.timer)
          this.timer = setTimeout(() => {
            this.showWrapper = false
          }, 300)
        }
      },
    },
  },
  mounted() {
    addResizeListener(this.$refs.resizeTrigger, this.resizeHandler)
  },
  beforeDestroy() {
    removeResizeListener(this.$refs.resizeTrigger, this.resizeHandler)
  },
  methods: {
    resizeHandler() {
      this.$nextTick(() => {
        this.bodyHeight = Math.min(
          this.$APP.height,
          this.$refs.resizeTrigger.offsetHeight
        )
      })
    },
    handleCloseDialog() {
      this.$emit('close')
      this.$emit('input', false)
    },
    closeDialog(event) {
      const dom = this.$refs.dialogContent
      if (dom) {
        if (!dom.contains(event.target)) {
          this.handleCloseDialog()
        }
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.dialog-wrapper {
  position: fixed;
  inset: 0;
  background-color: rgb(0 0 0 / 60%);
  z-index: 2000;

  @include scrollbar;

  .dialog-body {
    position: relative;
    margin: 0 auto;
    margin-top: calc((100vh - var(--height)) / 3);
    max-height: calc(100vh - 225px);
    min-height: 100px;
    border-radius: 2px;
    background-color: #fff;
    box-shadow: 0 4px 12px 0 rgb(0 0 0 / 20%);
    overflow: hidden;

    &.is-scrollable {
      @include scrollbar;
    }

    .icon-close {
      position: absolute;
      top: 6px;
      right: 6px;
      width: 32px;
      height: 32px;
      line-height: 32px;
      font-size: 22px;
      font-weight: 700;
      text-align: center;
      color: #d8d8d8;
      cursor: pointer;

      &:hover {
        color: #979ba5;
      }
    }
  }
}
</style>

<style lang="scss">
.dialog-fade-enter-active,
.dialog-fade-leave-active {
  transition: opacity 0.3s ease;
}

.dialog-fade-enter,
.dialog-fade-leave-to {
  opacity: 0;
}
</style>
