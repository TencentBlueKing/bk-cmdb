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
  <div class="clearfix">
    <dynamic-navigation class="main-navigation" v-show="!isEntry"></dynamic-navigation>
    <dynamic-breadcrumbs class="main-breadcrumbs" ref="breadcrumbs" v-if="showBreadcrumbs"></dynamic-breadcrumbs>
    <div class="main-layout">
      <div class="main-scroller" ref="scroller">
        <router-view class="main-views" :name="view" ref="view"></router-view>
      </div>
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import dynamicNavigation from './dynamic-navigation'
  import dynamicBreadcrumbs from './dynamic-breadcrumbs'
  import {
    addResizeListener,
    removeResizeListener
  } from '@/utils/resize-events'
  import { MENU_ENTRY, MENU_ADMIN } from '@/dictionary/menu-symbol'
  import throttle from 'lodash.throttle'
  export default {
    components: {
      dynamicNavigation,
      dynamicBreadcrumbs
    },
    data() {
      return {
        refreshKey: Date.now(),
        meta: this.$route.meta,
        scrollerObserver: null,
        scrollerObserverHandler: null
      }
    },
    computed: {
      ...mapGetters(['globalLoading']),
      view() {
        return this.meta.view
      },
      isEntry() {
        const [topRoute] = this.$route.matched
        return topRoute && [MENU_ENTRY, MENU_ADMIN].includes(topRoute.name)
      },
      showBreadcrumbs() {
        return this.$route.meta.layout && this.$route.meta.layout.breadcrumbs
      }
    },
    watch: {
      $route() {
        this.meta = this.$route.meta
      }
    },
    created() {
      this.scrollerObserverHandler = throttle(() => {
        const { scroller } = this.$refs
        if (scroller) {
          const gutter = scroller.offsetHeight - scroller.clientHeight
          this.$store.commit('setAppHeight', this.$root.$el.offsetHeight - gutter)
          this.$store.commit('setScrollerState', {
            scrollbar: scroller.scrollHeight > scroller.offsetHeight
          })
        }
      }, 300, { leading: false, trailing: true })
    },
    mounted() {
      addResizeListener(this.$refs.scroller, this.scrollerObserverHandler)
      this.addScrollerObserver()
    },
    beforeDestory() {
      removeResizeListener(this.$refs.scroller, this.scrollerObserverHandler)
      this.scrollerObserver && this.scrollerObserver.disconnect()
    },
    methods: {
      addScrollerObserver() {
        this.scrollerObserver = new MutationObserver(this.scrollerObserverHandler)
        this.scrollerObserver.observe(this.$refs.scroller, {
          attributes: true,
          childList: true,
          subtree: true
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .main-navigation {
        float: left;
    }
    .main-breadcrumbs {
        overflow: hidden;
        position: relative;
        background-color: #fff;
        z-index: 100;
    }
    .main-layout {
        position: relative;
        overflow: hidden;
        height: calc(100% - 53px);
        z-index: 99;
    }
    .main-scroller {
        height: 100%;
        overflow: auto;
    }
    .main-views {
        position: relative;
        height: 100%;
        min-width: 1089px;
    }
</style>
