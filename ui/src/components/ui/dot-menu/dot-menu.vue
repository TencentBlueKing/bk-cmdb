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
  <bk-popover
    ref="popover"
    trigger="click"
    placement="bottom-start"
    :sticky="true"
    :arrow="false"
    theme="light dot-menu-popover"
    :class="['dot-menu', {
      'is-open': open
    }]"
    :always="open"
    :on-show="show"
    :on-hide="hide">
    <i class="menu-trigger"
      :style="{
        '--color': color,
        '--hoverColor': hoverColor
      }">
    </i>
    <div class="menu-content" slot="content" @click.capture="handleContentClick">
      <slot></slot>
    </div>
  </bk-popover>
</template>

<script>
  export default {
    name: 'cmdb-dot-menu',
    props: {
      color: {
        type: String,
        default: '#979BA5'
      },
      hoverColor: {
        type: String,
        default: '#3A84FF'
      },
      closeWhenMenuClick: {
        type: Boolean,
        default: true
      }
    },
    data() {
      return {
        open: false
      }
    },
    methods: {
      show() {
        this.open = true
      },
      hide() {
        this.open = false
      },
      handleContentClick() {
        if (this.closeWhenMenuClick) {
          // eslint-disable-next-line no-underscore-dangle
          this.$refs.popover && this.$refs.popover.$refs.reference._tippy.hide()
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .dot-menu {
        width: 25px;
        height: 20px;
        line-height: 20px;
        text-align: center;
        font-size: 0;
        @include inlineBlock;
        &:hover,
        &.is-open {
            display: inline-block !important;
            .menu-trigger:before{
                background-color: var(--hoverColor);
                box-shadow: 0 -5px 0 0 var(--hoverColor), 0 5px 0 0 var(--hoverColor);
            }
        }
        /deep/ .bk-tooltip-ref {
            width: 100%;
            outline: none;
        }
        .menu-trigger {
            @include inlineBlock;
            width: 100%;
            cursor: pointer;
            &:before {
                @include inlineBlock;
                content: "";
                width: 3px;
                height: 3px;
                border-radius: 50%;
                background-color: var(--color);
                box-shadow: 0 -5px 0 0 var(--color), 0 5px 0 0 var(--color);
            }
        }
    }
</style>
<style lang="scss">
    .dot-menu-popover-theme {
        top: -6px;
        left: 10px;
        padding: 0 !important;
        .menu-content {
            font-size: 14px !important;
            background-color: #ffffff;
            button {
                font-size: 14px !important;
            }
        }
    }
</style>
