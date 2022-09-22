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
  <div class="breadcrumbs-layout clearfix">
    <template v-if="customize.enable">
      <i class="icon icon-cc-arrow fl" v-if="customize.backward" @click="customize.backward"></i>
      <h1 class="current fl">{{customize.title}}</h1>
    </template>
    <template v-else>
      <i class="icon icon-cc-arrow fl" v-if="(from && current)" @click="handleClick"></i>
      <h1 class="current fl">{{current}}</h1>
    </template>
  </div>
</template>

<script>
  import has from 'has'
  import { mapGetters } from 'vuex'
  import { Base64 } from 'js-base64'
  import Bus from '@/utils/bus'
  import { changeDocumentTitle } from '@/utils/change-document-title'

  export default {
    data() {
      return {
        /**
         * 自定义面包屑，如果启用了自定义面包屑，在模拟路由改变时需要手动关掉还原成自动生成的面包屑，否则会造成未知问题。
         * 自定义面包屑应该在没有路由变更时使用，一旦有路由变更会自动关闭。
         */
        customize: {
          enable: false, // 是否开启自定义面包屑
          title: '', // 自定义面包屑标题
          backward: null // 自定义面包屑返回函数
        },
      }
    },
    computed: {
      ...mapGetters(['title']),
      current() {
        const menuI18n = this.$route.meta.menu.i18n && this.$t(this.$route.meta.menu.i18n)
        return this.title || this.$route.meta.title || menuI18n
      },
      hideBack() {
        return (this.$route.meta.layout || {}).breadcrumbs?.back === false
      },
      defaultFrom() {
        const menu = this.$route.meta.menu || {}
        if (menu.relative) {
          return { name: Array.isArray(menu.relative) ? menu.relative[0] : menu.relative }
        }
        return null
      },
      latest() {
        let latest
        if (has(this.$route.query, '_f')) {
          try {
            const historyList = JSON.parse(window.sessionStorage.getItem('history'))
            latest = historyList.pop()
          } catch (e) {
            // ignore
          }
        }
        return latest
      },
      from() {
        if (this.latest) {
          try {
            return JSON.parse(Base64.decode(this.latest))
          } catch (error) {
            return this.defaultFrom
          }
        }
        return this.defaultFrom
      }
    },
    watch: {
      // 如果路由有所变更，则关闭手动设置面包屑
      current() {
        this.disableCustomize()
      }
    },
    created() {
      Bus.$on('enable-customize-breadcrumbs', this.enableCustomize)
      Bus.$on('disable-customize-breadcrumbs', this.disableCustomize)
    },
    methods: {
      async handleClick() {
        this.$routerActions.redirect({ ...this.from, back: true })
      },
      enableCustomize({ title, backward }) {
        if (!title) return

        this.customize.enable = true
        this.customize.title = title
        if (backward && typeof backward === 'function') {
          this.customize.backward = backward
        }
        changeDocumentTitle(title)
      },
      disableCustomize() {
        this.customize.enable = false
        this.customize.title = ''
        this.customize.backward = null
        changeDocumentTitle()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .breadcrumbs-layout {
        display: flex;
        align-items: center;
        padding: 14px 20px;
        height: 53px;
        background: #fff;
        box-shadow: 0px 2px 4px 0px rgba(0, 0, 0, 0.06);
        .icon-cc-arrow {
            display: block;
            width: 24px;
            height: 24px;
            line-height: 24px;
            font-size: 14px;
            text-align: center;
            margin-right: 3px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                color: #699df4;
            }
        }
        .current {
            font-size: 16px;
            line-height: 24px;
            color: #313238;
            font-weight: normal;
            @include ellipsis;
        }
    }
</style>
