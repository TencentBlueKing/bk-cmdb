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

<script setup>
  import DropdownOptionButton from '../dropdown-option-button.vue'
  defineProps({
    commands: {
      type: [Object, Array],
      default: () => ([])
    }
  })
</script>

<template>
  <div class="more-action-menu">
    <slot name="append"></slot>
    <cmdb-dot-menu color="#979BA5">
      <ul v-if="commands.length">
        <li v-for="(cmd, cmdIndex) in commands" :key="cmdIndex">
          <template v-if="cmd.isShow">
            <!-- 带权限的按钮 -->
            <cmdb-auth
              v-if="cmd.auth"
              :auth="cmd.auth"
            >
              <dropdown-option-button
                class="more-action-menu-option"
                :disabled="cmd.disabled"
                @click="cmd.handler"
                v-bk-tooltips="$t(cmd.tips)"
              >
                {{ cmd.text }}
              </dropdown-option-button>
            </cmdb-auth>

            <!-- 无权限按钮 -->
            <dropdown-option-button
              v-else
              :disabled="cmd.disabled"
              @click="cmd.handler"
              v-bk-tooltips="$t(cmd.tips)"
            >
              {{ cmd.text }}
            </dropdown-option-button>
          </template>
        </li>
      </ul>
    </cmdb-dot-menu>
  </div>
</template>

<style  lang="scss" scoped>
.more-action-menu {
  @include space-between;
  justify-content: flex-start;

  .more-action-menu-option {
    font-weight: normal;
  }
}
.dropdown-option-button {
  cursor: pointer;
  padding: 0 15px !important;
}
</style>
