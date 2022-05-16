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
  <div class="cmdb-main-inject"
    :style="{
      width: scrollbar ? 'calc(100% - 17px)' : '100%'
    }">
    <slot></slot>
  </div>
</template>

<script>
  export default {
    name: 'cmdb-main-inject',
    props: {
      injectType: {
        type: String,
        default: 'append',
        validator(val) {
          return ['append', 'prepend'].includes(val)
        }
      }
    },
    computed: {
      scrollbar() {
        return this.$store.state.scrollerState.scrollbar
      }
    },
    mounted() {
      const $mainLayout = document.querySelector('.main-layout')
      if (this.injectType === 'append') {
        $mainLayout.append(this.$el)
      } else {
        $mainLayout.insertBefore(this.$el, $mainLayout.firstChild)
      }
    },
    beforeDestroy() {
      this.$el.remove()
    }
  }
</script>
