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
  <div class="bk-button-group">
    <bk-button
      :class="{ 'is-selected': active === 'instance' }"
      @click="handleSwitch('instance')">
      {{$t('实例')}}
    </bk-button>
    <bk-button ref="tipsReference"
      :class="{ 'is-selected': active === 'process' }"
      v-bk-tooltips="{
        content: '#tipsContent',
        allowHtml: true,
        placement: 'bottom-end',
        disabled: tipsDisabled || !showTips,
        showOnInit: !tipsDisabled,
        hideOnClick: false,
        trigger: 'manual',
        theme: 'view-switer-tips',
        zIndex: getZIndex()
      }"
      @click="handleSwitch('process')">
      {{$t('进程')}}
    </bk-button>
    <span v-if="showTips" class="tips-content" id="tipsContent">
      {{$t('切换进程视角提示语')}}
      <i class="bk-icon icon-close" @click="handleCloseTips"></i>
    </span>
  </div>
</template>

<script>
  import RouterQuery from '@/router/query'
  export default {
    props: {
      showTips: {
        type: Boolean,
        default: true
      }
    },
    data() {
      return {
        tipsDisabled: !!window.localStorage.getItem('service_instance_view_switcher'),
        active: RouterQuery.get('view', 'instance')
      }
    },
    watch: {
      active(active) {
        if (active === 'process') {
          this.hideTips()
        }
      }
    },
    mounted() {
      if (this.active === 'process') {
        this.hideTips()
      }
    },
    methods: {
      getZIndex() {
        // eslint-disable-next-line no-underscore-dangle
        return window.__bk_zIndex_manager.nextZIndex()
      },
      handleSwitch(active) {
        RouterQuery.set({ view: active })
      },
      hideTips() {
        this.tipsDisabled = true
        const { tippyInstance } = this.$refs.tipsReference.$el
        tippyInstance && tippyInstance.hide()
      },
      handleCloseTips() {
        window.localStorage.setItem('service_instance_view_switcher', 'closed')
        this.hideTips()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .bk-button-group {
        .tips-content {
            display: none;
        }
    }
    .tips-content {
        display: flex;
        align-items: center;
        position: relative;
        .bk-icon {
            margin-left: 20px;
            cursor: pointer;
            font-size: 16px;
        }
    }
</style>

<style lang="scss">
    .tippy-tooltip.view-switer-tips-theme {
        background-color: #699df4;
        color: #fff;
        .tippy-arrow {
            border-bottom-color: #699df4;
        }
    }
</style>
