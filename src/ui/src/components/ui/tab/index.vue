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
  defineProps({
    tabs: {
      type: Array,
      default: () => []
    },
    active: {
      type: [String, Symbol]
    }
  })

  const emit = defineEmits(['change'])

  const handleChange = (tab, index) => {
    emit('change', tab, index)
  }
</script>
<script>
  export default {
    name: 'cmdb-tab'
  }
</script>

<template>
  <div class="tab">
    <div :class="['tab-item', { active: active === tab.id }]"
      v-for="(tab, index) in tabs"
      :key="index"
      @click="handleChange(tab, index)">
      <span class="tab-text">{{ tab.text }}</span>
      <em class="tab-count" v-if="tab.count >= 0">{{ tab.count }}</em>
      <slot name="append"></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.tab {
  display: flex;
  gap: 8px;

  .tab-item {
    display: flex;
    align-items: center;
    gap: 8px;
    height: 42px;
    padding: 0 24px;
    justify-content: center;
    background: #EAEBF0;
    border-radius: 4px 4px 0 0;
    cursor: pointer;

    .tab-text {
      font-size: 14px;
      color: #63656E;
    }
    .tab-count {
      display: flex;
      align-items: center;
      justify-content: center;
      min-width: 30px;
      height: 18px;
      font-size: 14px;
      color: #63656E;
      font-style: normal;
      background: #DCDEE5;
      border-radius: 10px;
      margin-top: -1px;
    }

    &:hover,
    &.active {
      .tab-text {
        color: $primaryColor;
      }
    }

    &.active {
      background: #fff;
      .tab-count {
        color: $primaryColor;
        background: #F0F5FF;
      }
    }
  }
}
</style>
