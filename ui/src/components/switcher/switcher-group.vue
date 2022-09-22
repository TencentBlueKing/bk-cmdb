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
  <bk-popover ref="popover" v-bind="bindProps">
    <span class="switcher-group">
      <slot></slot>
    </span>
    <div slot="content" class="switcher-group-tips">
      <span>{{tips}}</span>
      <i class="bk-icon icon-close" @click="handleCloseTips"></i>
    </div>
  </bk-popover>
</template>

<script>
  export default {
    name: 'cmdb-switcher-group',
    props: {
      tipsKey: {
        type: String,
        default: ''
      },
      tips: {
        type: String,
        default: ''
      },
      value: {
        type: [String, Number]
      }
    },
    data() {
      return {
        isSwitcherGroup: true
      }
    },
    computed: {
      bindProps() {
        let showOnInit = false
        if (this.tipsKey) {
          showOnInit = window.localStorage.getItem(this.tipsKey) === null
        }
        return {
          theme: 'switcher-group-tips',
          placement: 'bottom',
          trigger: 'manual',
          showOnInit,
          disabled: !this.tips,
          tippyOptions: {
            hideOnClick: false
          },
          ...this.$attrs
        }
      },
      active: {
        get() {
          return this.value
        },
        set(active) {
          this.$emit('input', active)
          this.$emit('change', active)
        }
      }
    },
    methods: {
      setActive(name) {
        this.active = name
        this.handleCloseTips()
      },
      handleCloseTips() {
        this.$refs.popover.instance.hide()
        window.localStorage.setItem(this.tipsKey, false)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .switcher-group {
        display: inline-flex;
        justify-content: center;
        align-items: center;
    }
    .switcher-group-tips {
        .bk-icon {
            font-size: 14px;
            margin-left: 20px;
        }
    }
</style>

<style lang="scss">
    .tippy-tooltip.switcher-group-tips-theme {
        background-color: #699DF4;
        .tippy-arrow {
            border-bottom-color: #699DF4;
        }
    }
</style>
