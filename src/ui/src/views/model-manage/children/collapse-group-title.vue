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
    class="collapse-group-title"
    :class="{
      'is-dropdown-visible': isDropdownVisible,
      'is-collapse': collapse
    }">
    <div class="drag-icon" v-if="dragIcon"></div>
    <bk-icon class="group-collapse-icon" type="down-shape" />
    <span class="group-title-text">{{title}}</span>
    <bk-dropdown-menu
      @click.native.stop
      v-if="dropdownMenu"
      class="group-dropdown-menu"
      trigger="click"
      align="left"
    >
      <template #dropdown-trigger>
        <div class="more-operation-btn">
          <bk-icon type="more" />
        </div>
      </template>

      <template #dropdown-content>
        <ul class="bk-dropdown-list" v-if="visibleCommands.length">
          <li v-for="(cmd, cmdIndex) in visibleCommands" :key="cmdIndex">
            <!-- 带权限的按钮 -->
            <cmdb-auth
              v-if="cmd.auth"
              :auth="cmd.auth"
              @update-auth="cmd.onUpdateAuth"
            >
              <template #default="{ disabled: authDisabled }">
                <dropdown-option-button
                  v-bk-tooltips="{
                    content: cmd.disabledTooltips,
                    disabled: !cmd.disabled && cmd.disabledTooltips,
                    placement: 'right'
                  }"
                  :disabled="authDisabled || cmd.disabled"
                  @click="cmd.handler"
                >
                  {{ cmd.text }}
                </dropdown-option-button>
              </template>
            </cmdb-auth>

            <!-- 无权限按钮 -->
            <dropdown-option-button
              v-else
              :disabled="cmd.disabled"
              @click="cmd.handler"
            >
              {{ cmd.text }}
            </dropdown-option-button>
          </li>
        </ul>
      </template>
    </bk-dropdown-menu>
  </div>
</template>

<script>
  import { defineComponent, computed, ref, toRef } from '@vue/composition-api'
  import DropdownOptionButton from './dropdown-option-button.vue'
  import has from 'has'

  export default defineComponent({
    name: 'CollapseGroupTitle',
    components: {
      DropdownOptionButton
    },
    props: {
      // 标题内容
      title: {
        type: String,
        default: '',
        required: true
      },
      // 是否展示拖拽图标
      dragIcon: {
        type: Boolean,
        default: false
      },
      // 是否开启下拉菜单
      dropdownMenu: {
        type: Boolean,
        default: true
      },
      // 是否处于收缩状态
      collapse: {
        type: Boolean,
        default: false
      },
      /**
       * 下拉菜单指令组
       * @property {Array} commands
       * @property {String} commands[].text 指令按钮文案
       * @property {Object} commands[].auth 指令按钮权限，沿用 cmdb-auth 组件的配置
       * @property {Object} commands[].onUpdateAuth 指令按钮权限更新回调，沿用 cmdb-auth 组件的配置
       * @property {Boolean} commands[].visible 是否展示指令按钮
       * @property {Boolean} commands[].disabled 是否禁用指令按钮
       * @property {Boolean} commands[].disabledTooltips 禁用指令按钮提示文案
       * @property {Function} commands[].handler 指令按钮点击触发的事件
       */
      commands: {
        type: Array,
        default: () => []
      }
    },
    setup(props) {
      const isDropdownVisible = ref(false)
      const commandsRef = toRef(props, 'commands')
      const visibleCommands = computed(() => commandsRef.value.map((cmd) => {
        if (!has(cmd, 'visible')) {
          cmd.visible = true
        }
        if (!has(cmd, 'disabled')) {
          cmd.disabled = false
        }
        if (!has(cmd, 'auth')) {
          cmd.auth = null
        }
        if (!has(cmd, 'onUpdateAuth')) {
          cmd.onUpdateAuth = () => {}
        }

        return cmd
      })
        .filter(cmd => cmd.visible))

      return {
        visibleCommands,
        isDropdownVisible,
      }
    },
  })
</script>

<style lang="scss" scoped>
$titleHeight: 26px;

.collapse-group-title {
  display: flex;
  align-items: center;
  height: $titleHeight;
  border-radius: 2px;
  cursor: pointer;
  user-select: none;

  &:hover {
    background-color: #f0f1f5;

    .drag-icon {
      display: block;
    }
  }

  .drag-icon {
    margin-left: 5px;
    margin-right: 5px;
    flex-shrink: 0;
    display: none;
    @include dragIcon;
  }

  .group-collapse-icon {
    margin: 0 4px;
    color: #63656E;
    transition: transform 200ms ease;
  }

  .group-title-text {
    margin-right: 5px;
  }

  .group-dropdown-menu {
    .more-operation-btn {
      width: $titleHeight;
      height: $titleHeight;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 16px;
      color: #979BA5;
      cursor: pointer;
      user-select: none;

      &:hover {
        background-color: #eaebf0;
        color: #3a84ff;
      }
    }
  }

  &.is-dropdown-visible {
    .more-operation-btn {
      background-color: #eaebf0;
      color: #3a84ff;
    }
  }

  &.is-collapse {
    .group-collapse-icon{
      transform: rotate(-90deg);
    }
  }
}
</style>
