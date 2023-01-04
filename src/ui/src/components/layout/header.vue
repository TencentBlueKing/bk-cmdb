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
  <header class="header-layout" v-test-id.global="'header'">
    <div class="logo">
      <router-link class="logo-link" to="/index">
        {{$t('蓝鲸配置平台')}}
      </router-link>
    </div>
    <nav class="header-nav" v-test-id.global="'headerNav'">
      <router-link class="header-link"
        v-for="nav in visibleMenu"
        :to="getHeaderLink(nav)"
        :key="nav.id"
        :class="{
          active: isLinkActive(nav)
        }">
        {{$t(nav.i18n)}}
      </router-link>
    </nav>
    <section class="header-info">
      <bk-popover class="info-item"
        theme="light header-info-popover"
        :arrow="false"
        :tippy-options="{
          animateFill: false,
          hideOnClick: false
        }">
        <i :class="`bk-icon icon-${currentSysLang.icon} lang-icon`"></i>
        <div slot="content">
          <div
            v-for="(lang, index) in sysLangs"
            :key="index"
            :class="['link-item', { active: $i18n.locale === lang.id }]"
            @click="handleToggleLang(lang.id)">
            <i :class="`bk-icon icon-${lang.icon} lang-icon`"></i>
            <span>{{lang.name}}</span>
          </div>
        </div>
      </bk-popover>
      <bk-popover class="info-item"
        theme="light header-info-popover"
        animation="fade"
        placement="bottom-end"
        :arrow="false"
        :tippy-options="{
          animateFill: false,
          hideOnClick: false
        }">
        <i class="question-icon icon-cc-default"></i>
        <template slot="content">
          <a class="link-item" target="_blank" :href="helpDocUrl">{{$t('产品文档')}}</a>
          <a class="link-item" target="_blank" href="https://bk.tencent.com/s-mart/community">{{$t('问题反馈')}}</a>
          <a class="link-item" target="_blank" href="https://github.com/Tencent/bk-cmdb">{{$t('开源社区')}}</a>
        </template>
      </bk-popover>
      <bk-popover class="info-item"
        theme="light header-info-popover"
        animation="fade"
        placement="bottom-end"
        :arrow="false"
        :tippy-options="{
          animateFill: false,
          hideOnClick: false
        }">
        <span class="info-user">
          <span class="user-name">{{userName}}</span>
          <i class="user-icon bk-icon icon-angle-down"></i>
        </span>
        <template slot="content">
          <a class="link-item" href="javascript:void(0)"
            @click="handleLogout">
            <i class="icon-cc-logout"></i>
            {{$t('注销')}}
          </a>
        </template>
      </bk-popover>
    </section>
  </header>
</template>

<script>
  import has from 'has'
  import menu from '@/dictionary/menu'
  import {
    MENU_BUSINESS,
    MENU_BUSINESS_SET,
    MENU_BUSINESS_SET_TOPOLOGY,
    MENU_BUSINESS_HOST_AND_SERVICE
  } from '@/dictionary/menu-symbol'
  import { mapGetters } from 'vuex'
  import {
    getBizSetIdFromStorage,
    getBizSetRecentlyUsed
  } from '@/utils/business-set-helper.js'
  import { changeLocale } from '@/i18n'
  import { LANG_SET } from '@/i18n/constants'

  export default {
    data() {
      return {
        sysLangs: LANG_SET
      }
    },
    computed: {
      ...mapGetters(['userName']),
      ...mapGetters('objectBiz', ['bizId']),
      helpDocUrl() {
        return this.$Site.helpDocUrl || 'http://docs.bk.tencent.com/product_white_paper/cmdb/'
      },
      visibleMenu() {
        return menu.filter((menuItem) => {
          if (!has(menuItem, 'visibility')) {
            return true
          }

          if (typeof menuItem.visibility === 'function') {
            return menuItem.visibility(this)
          }
          return menuItem.visibility
        })
      },
      isBusinessView() {
        const { matched: [topRoute] } = this.$route
        return topRoute?.name === MENU_BUSINESS_SET || topRoute?.name === MENU_BUSINESS
      },
      currentSysLang() {
        return this.sysLangs.find(lang => lang.id === this.$i18n.locale) || {}
      }
    },
    methods: {
      isLinkActive(nav) {
        const { matched: [topRoute] } = this.$route
        if (!topRoute) {
          return false
        }
        return topRoute.name === nav.id
          || (topRoute.name === MENU_BUSINESS_SET && nav.id === MENU_BUSINESS) // 业务集时高亮业务菜单
      },
      getHeaderLink(nav) {
        const link = { name: nav.id }
        if (nav.id === MENU_BUSINESS) {
          const isBizSetUsed = getBizSetRecentlyUsed()
          const id = isBizSetUsed ? getBizSetIdFromStorage() : this.bizId
          const paramName = isBizSetUsed ? 'bizSetId' : 'bizId'
          link.name = isBizSetUsed ? MENU_BUSINESS_SET_TOPOLOGY : MENU_BUSINESS_HOST_AND_SERVICE
          link.params = {
            [paramName]: id
          }
        }
        return link
      },
      handleToggleLang(locale) {
        changeLocale(locale)
      },
      handleLogout() {
        this.$http.post(`${window.API_HOST}logout`, {
          http_scheme: window.location.protocol.replace(':', '')
        }).then((data) => {
          window.location.href = data.url
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
    .header-layout {
        position: relative;
        display: flex;
        height: 58px;
        background-color: #182132;
        z-index: 1002;
    }
    .logo {
        flex: 292px 0 0;
        font-size: 0;
        .logo-link {
            display: inline-block;
            vertical-align: middle;
            height: 58px;
            line-height: 58px;
            margin-left: 23px;
            padding-left: 38px;
            color: #fff;
            font-size: 16px;
            background: url("../../assets/images/logo.svg") no-repeat 0 center;
        }
    }
    .header-nav {
        flex: 3;
        font-size: 0;
        white-space: nowrap;
        .header-link {
            display: inline-block;
            vertical-align: middle;
            height: 58px;
            line-height: 58px;
            padding: 0 25px;
            color: #979BA5;
            font-size: 14px;
            &:hover {
                background-color: rgba(49, 64, 94, .5);
                color: #fff;
            }
            &.router-link-active,
            &.active {
                background-color: rgba(49, 64, 94, 1);
                color: #fff;
            }
        }
    }
    .header-info {
        flex: 1;
        text-align: right;
        white-space: nowrap;
        @include middleBlockHack;
    }
    .info-item {
        @include inlineBlock;
        margin: 0 20px 0 0;
        text-align: left;
        font-size: 0;
        cursor: pointer;
        &:last-child {
            margin-right: 24px;
        }
        .tippy-active {
            .bk-icon {
                color: #fff;
            }
            .user-icon {
                transform: rotate(-180deg);
            }
        }
        .question-icon,
        .lang-icon {
            font-size: 16px;
            color: #DCDEE5;
            &:hover {
                color: #fff;
            }
        }
        .info-user {
            font-size: 14px;
            font-weight: bold;
            color: #fff;
            .user-name {
                max-width: 150px;
                @include inlineBlock;
                @include ellipsis;
            }
            .user-icon {
                margin-left: -4px;
                transition: transform .2s linear;
                font-size: 20px;
                color: #fff;
            }
        }
        .lang-icon {
            font-size: 20px;
        }
    }

    .header-info-popover-theme {
      .link-item {
          display: flex;
          padding: 0 20px;
          height: 40px;
          font-size: 14px;
          white-space: nowrap;
          align-items: center;
          cursor: pointer;

          &:hover,
          &.active {
              background-color: #f1f7ff;
              color: #3a84ff;
          }

          .lang-icon {
              font-size: 16px;
              margin-right: 4px;
          }
      }
    }
</style>

<style>
    .tippy-tooltip.header-info-popover-theme {
        padding-left: 0 !important;
        padding-right: 0 !important;
        overflow: hidden;
        border-radius: 2px !important;
    }
</style>
