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
  <div v-if="showTips" class="cmdb-tips" :style="tipsStyle">
    <i v-if="icon" :class="icon" :style="iconStyle"></i>
    <p
      class="tips-content"
      :class="{ ellipsis: localEllipsis, pr20: showClose }">
      <slot></slot>
      <a v-if="moreLink" class="more" :href="moreLink" target="_blank"
        >{{ $t('更多详情') }} &gt;&gt;</a
      >
    </p>
    <i v-if="showClose" class="icon-cc-tips-close" @click="handleClose"></i>
  </div>
</template>

<script>
export default {
  name: 'cmdb-tips',
  props: {
    tipsKey: {
      type: String,
      default: '',
    },
    tipsStyle: {
      type: Object,
      default: () => ({}),
    },
    icon: {
      type: [String, Boolean],
      default: 'icon icon-cc-exclamation-tips',
    },
    iconStyle: {
      type: Object,
      default: () => ({}),
    },
    ellipsis: {
      type: Boolean,
      default: false,
    },
    moreLink: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      localEllipsis: false,
      showClose: false,
      showTips: true,
    }
  },
  watch: {
    moreLink(link) {
      this.localEllipsis = link ? false : this.ellipsis
    },
  },
  created() {
    this.setStatus()
  },
  methods: {
    setStatus() {
      let value = !this.tipsKey
      if (this.tipsKey) {
        const localValue = window.localStorage.getItem(this.tipsKey)
        value = localValue === null ? true : localValue === 'true'
      }
      this.$emit('input', value)
      this.showTips = value
      this.showClose = !!this.tipsKey
    },
    handleClose() {
      window.localStorage.setItem(this.tipsKey, false)
      this.$emit('input', false)
      this.showTips = false
    },
  },
}
</script>

<style lang="scss" scoped>
.cmdb-tips {
  display: flex;
  min-height: 30px;
  font-size: 12px;
  background: #f0f8ff;
  border-radius: 2px;
  border: 1px solid #a3c5fd;
  padding: 0 16px;
  align-items: center;

  .icon {
    flex: 16px 0 0;
    text-align: center;
    line-height: 16px;
    font-size: 16px;
    color: #3a84ff;
    margin-right: 5px;
  }

  .tips-content {
    flex: 1;

    &.ellipsis {
      @include ellipsis;
    }

    .more {
      display: inline-block;
      margin-left: 20px;
      color: #3a84ff;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .icon-cc-tips-close {
    align-self: center;
    width: 12px;
    height: 12px;
    font-size: 12px;
    color: #979ba5;
    cursor: pointer;
  }
}
</style>
