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
    status: {
      type: String,
      default: ''
    },
    count: {
      type: Number,
      default: ''
    },
    text: {
      type: String,
      default: ''
    },
    tooltips: {
      type: String,
      default: ''
    }
  })
</script>

<template>
  <div :class="['diff-brand', status]">
    <i class="icon-box"></i>
    <div class="content">
      <span class="count">{{ count }}</span>
      <span :class="['text', { 'g-has-dashed-tooltips': tooltips.length > 0 }]"
        v-bk-tooltips="{ content: tooltips, disabled: !tooltips.length, allowHTML: true }">{{ text }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.diff-brand {
  --new-color: #2DCB56;
  --update-color: #FF9C01;
  --conflict-color: #EA3636;
  --unbound-color: #979BA5;
  --unchanged-color: #979BA5;

  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  .icon-box {
    width: 12px;
    height: 12px;
    background: #fff;
    border: 1px solid var(--new-color);
    border-radius: 2px;
  }
  .content {
    display: flex;
    align-items: center;
    gap: 2px;
  }
  .count {
    color: #63656E;
    font-weight: 700;
  }
  .text {
    color: #63656E;
  }

  &.new {
    .icon-box {
      background: #F2FFF4;
      border: 1px solid var(--new-color);
    }
    .count {
      color: var(--new-color);
    }
  }
  &.update {
    .icon-box {
      background: #FFF3E1;
      border: 1px solid var(--update-color);
    }
    .count {
      color: var(--update-color);
    }
  }
  &.conflict {
    .icon-box {
      background: #FFEEEE;
      border: 1px solid var(--conflict-color);
    }
    .count {
      color: var(--conflict-color);
    }
  }
  &.unbound {
    .icon-box {
      background: #DCDEE5;
      border: 1px solid var(--unbound-color);
    }
    .count {
      color: var(--unbound-color);
    }
  }
  &.unchanged {
    .icon-box {
      background: #fff;
      border: 1px solid var(--unchanged-color);
    }
    .count {
      color: var(--unchanged-color);
    }
  }
}
</style>
