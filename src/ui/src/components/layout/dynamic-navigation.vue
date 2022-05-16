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
  <nav class="nav-layout" v-test-id.dynamic="'asideNav'"
    :class="{ 'sticked': navStick }"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave">
    <div class="nav-wrapper"
      :class="{ unfold: unfold, flexible: !navStick }">
      <div class="business-wrapper" v-if="isBusinessView">
        <transition name="fade">
          <cmdb-business-mix-selector class="business-selector" v-test-id.dynamic
            v-show="unfold"
            :value="currentBusinessViewPayload.selectedId"
            :popover-options="{
              appendTo: () => this.$el
            }"
            @select="handleToggleBusiness">
          </cmdb-business-mix-selector>
        </transition>
        <transition name="fade">
          <i class="business-flag bk-icon icon-angle-down" v-show="!unfold"></i>
        </transition>
      </div>
      <div class="menu-list">
        <template v-for="(menu, index) in currentMenus">
          <router-link class="menu-item is-link" tag="a" active-class="active" style="display: block;"
            v-if="menu.hasOwnProperty('route')"
            ref="menuLink"
            :class="{
              'is-relative-active': isRelativeActive(menu)
            }"
            :key="index"
            :to="getMenuLink(menu)"
            :title="$t(menu.i18n)">
            <h3 class="menu-info clearfix">
              <i :class="['menu-icon', menu.icon]"></i>
              <span class="menu-name">{{$t(menu.i18n)}}</span>
            </h3>
          </router-link>
        </template>
      </div>
      <div class="nav-option">
        <i class="nav-stick icon icon-cc-nav-toggle"
          :class="{
            sticked: navStick
          }"
          :title="navStick ? $t('收起导航') : $t('固定导航')"
          @click="toggleNavStick">
        </i>
      </div>
    </div>
  </nav>
</template>
<script>
  import has from 'has'
  import { mapGetters } from 'vuex'
  import MENU_DICTIONARY from '@/dictionary/menu'
  import {
    MENU_BUSINESS,
    MENU_BUSINESS_SET,
    MENU_BUSINESS_SET_TOPOLOGY,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_RESOURCE,
    MENU_RESOURCE_INSTANCE,
  } from '@/dictionary/menu-symbol'
  import { BUILTIN_MODEL_RESOURCE_MENUS } from '@/dictionary/model-constants.js'

  export default {
    data() {
      return {
        routerLinkHeight: 42,
        timer: null,
        state: {},
        hasExactActive: false
      }
    },
    computed: {
      ...mapGetters(['navStick', 'navFold', 'admin']),
      ...mapGetters('objectModelClassify', ['classifications', 'models']),
      ...mapGetters('userCustom', ['resourceCollection']),
      isBusinessNav() {
        const { matched: [topRoute] } = this.$route
        return topRoute?.name === MENU_BUSINESS
      },
      isBusinessSetNav() {
        const { matched: [topRoute] } = this.$route
        return topRoute?.name === MENU_BUSINESS_SET
      },
      // 业务或业务集统称为业务视图
      isBusinessView() {
        return this.isBusinessNav || this.isBusinessSetNav
      },
      // 当前业务视图附属数据
      currentBusinessViewPayload() {
        if (!this.isBusinessView) {
          return {}
        }
        const paramName = this.isBusinessNav ? 'bizId' : 'bizSetId'
        const id = Number(this.$route.params[paramName])
        const selectedId = `${id}-${this.isBusinessNav ? 'biz' : 'bizset'}`

        // 路由参数中的业务(集)id,参数名称,业务选择器的id值
        return { id, paramName, selectedId }
      },
      unfold() {
        return this.navStick || !this.navFold
      },
      owner() {
        return this.$route.matched[0].name
      },
      collection() {
        if (this.owner === MENU_RESOURCE) {
          return this.resourceCollection
        }
        return []
      },
      collectionMenus() {
        return this.collection.map((id) => {
          const model = this.models.find(model => model.bk_obj_id === id)
          return {
            i18n: model.bk_obj_name,
            icon: model.bk_obj_icon,
            id: `collection_${id}`,
            route: this.getCollectionRoute(model)
          }
        })
      },
      currentMenus() {
        const target = MENU_DICTIONARY.find(menu => menu.id === this.owner)
        const menus = [...((target && target.menu) || [])]
        if (this.owner === MENU_RESOURCE) {
          menus.splice(1, 0, ...this.collectionMenus)
        }
        return menus
      },
      relativeActiveName() {
        const relative = this.$tools.getValue(this.$route, 'meta.menu.relative')
        if (relative && !this.hasExactActive) {
          const names = Array.isArray(relative) ? relative : [relative]
          let relativeActiveName = null
          for (let index = 0; index < names.length; index++) {
            const name = names[index]
            const isActive = this.currentMenus.some((menu) => {
              if (has(menu, 'route')) {
                return menu.route.name === name
              }
              if (menu.submenu && menu.submenu.length) {
                return menu.submenu.some(submenu => submenu.route.name === name)
              }
              return false
            })
            if (isActive) {
              relativeActiveName = name
              break
            }
          }
          return relativeActiveName
        }
        return null
      }
    },
    watch: {
      $route: {
        immediate: true,
        handler() {
          this.setDefaultExpand()
          this.checkExactActive()
        }
      },
      isBusinessNav(isBusinessNav) {
        isBusinessNav && this.refreshAuthorizedList()
      }
    },
    methods: {
      setDefaultExpand() {
        const expandedId = this.$route.meta?.menu?.parent
        if (expandedId) {
          this.$set(this.state, expandedId, { expanded: true })
        } else if (this.relativeActiveName) {
          const parent = this.currentMenus.find((menu) => {
            if (has(menu, 'route')) {
              return menu.route.name === this.relativeActiveName
            }
            return menu.submenu.some(submenu => submenu.route.name === this.relativeActiveName)
          })
          if (parent) {
            this.$set(this.state, parent.id, { expanded: true })
          }
        }
      },
      checkExactActive() {
        if (!this.$refs.menuLink) {
          return
        }
        this.$nextTick(() => {
          this.hasExactActive = this.$refs.menuLink.some(link => link.$el.classList.contains('active'))
        })
      },
      isRelativeActive(menu) {
        return menu.route.name === this.relativeActiveName
      },
      getCollectionRoute(model) {
        if (has(BUILTIN_MODEL_RESOURCE_MENUS, model.bk_obj_id)) {
          return {
            name: BUILTIN_MODEL_RESOURCE_MENUS[model.bk_obj_id]
          }
        }
        return {
          name: MENU_RESOURCE_INSTANCE,
          params: {
            objId: model.bk_obj_id
          }
        }
      },
      getMenuLink(menu) {
        if (this.isBusinessNav) {
          return {
            name: menu.route.name,
            params: {
              bizId: this.$store.getters['objectBiz/bizId']
            }
          }
        }
        return menu.route
      },
      handleMouseEnter() {
        if (this.timer) {
          clearTimeout(this.timer)
        }
        this.$store.commit('setNavStatus', { fold: false })
      },
      handleMouseLeave() {
        this.timer = setTimeout(() => {
          this.$store.commit('setNavStatus', { fold: true })
        }, 300)
      },
      // 切换导航展开固定
      toggleNavStick() {
        this.$store.commit('setNavStatus', {
          fold: !this.navFold,
          stick: !this.navStick
        })
      },
      handleToggleBusiness(value, newId, isBizSet) {
        // 切换前的业务(集)id值
        const { selectedId } = this.currentBusinessViewPayload

        // 非业务视图或切换的前后值相同则取消执行
        if (!this.isBusinessView || value === selectedId) {
          return false
        }

        const paramName = isBizSet ? 'bizSetId' : 'bizId'
        this.$routerActions.redirect({
          name: isBizSet ? MENU_BUSINESS_SET_TOPOLOGY : MENU_BUSINESS_HOST_AND_SERVICE,
          params: {
            ...this.$route.params,
            [paramName]: newId
          },
          reload: true
        })
      },
      async refreshAuthorizedList() {
        try {
          const { info } = await this.$store.dispatch('objectBiz/getAuthorizedBusiness')
          this.$store.commit('objectBiz/setAuthorizedBusiness', Object.freeze(info))
          const { bizId } = this.$route.params
          const exist = info.some(biz => biz.bk_biz_id.toString() === bizId.toString())
          if (!exist) {
            this.$route.matched[0].meta.view = 'permission'
          }
        } catch (error) {
          console.error(error)
        }
      }
    }
  }
</script>
<style lang="scss" scoped>
$cubicBezier: cubic-bezier(0.4, 0, 0.2, 1);
$duration: 0.2s;
$color: #63656E;

.nav-layout {
    position: relative;
    width: 60px;
    height: 100%;
    transition: width $duration $cubicBezier;
    z-index: 1000;
    &.sticked {
        width: 260px;
    }

    .nav-wrapper {
        position: relative;
        width: 100%;
        height: 100%;
        border-right: 1px solid #DCDEE5;
        background: #fff;
        transition: width $duration $cubicBezier;

        &.unfold {
            width: 260px;
        }

        &.unfold.flexible:after {
            content: "";
            position: absolute;
            width: 15px;
            height: 100%;
            left: 100%;
            top: 0;
        }
    }
}

.business-wrapper {
    position: relative;
    padding: 10px 0;
    height: 53px;
    border-bottom: 1px solid #DCDEE5;
    overflow: hidden;
    .business-selector {
        display: block;
        width: 240px;
        margin: 0 auto;
    }
    .business-flag {
        position: absolute;
        left: 10px;
        top: 9px;
        width: 32px;
        height: 32px;
        line-height: 30px;
        text-align: center;
        font-size: 20px;
        border: 1px solid #C4C6CC;
        border-radius: 2px;
    }
}

.menu-list {
    padding: 10px 0;
    height: calc(100% - 120px);
    overflow-y: auto;
    overflow-x: hidden;
    white-space: nowrap;

    &::-webkit-scrollbar {
        width: 5px;
        height: 5px;

        &-thumb {
            border-radius: 20px;
            background: rgba(165, 165, 165, .3);
            box-shadow: inset 0 0 6px hsla(0, 0%, 80%, .3);
        }
    }

    .menu-item {
        position: relative;
        &.is-link {
            .icon-close {
                display: none;
                position: absolute;
                right: 20px;
                top: 10px;
                width: 25px;
                height: 25px;
                line-height: 25px;
                font-size: 16px;
                text-align: center;
                color: $color;
                &:hover {
                    opacity: .8;
                }
                &:before {
                    display: block;
                    transform: scale(.65);
                }
            }
            &:hover {
                .icon-close {
                    display: block;
                }
            }
        }
        &:hover {
            background-color: #F6F6F9;
        }
        &.is-relative-active.is-link,
        &.active.is-link {
            background-color: #E1ECFF;
            .menu-icon,
            .menu-name,
            .icon-close {
                color: $primaryColor;
            }
        }
        .menu-info {
            margin: 0;
            padding: 0;
            height: 42px;
            line-height: 42px;
            white-space: nowrap;
            font-size: 0;
            font-weight: normal;
            color: $color;
            cursor: pointer;
        }

        .menu-icon {
            display: inline-block;
            vertical-align: top;
            margin: 13px 26px 13px 22px;
            font-size: 16px;
            color: #979BA5;
        }

        .menu-name {
            display: inline-block;
            width: calc(100% - 120px);
            vertical-align: top;
            font-size: 14px;
            @include ellipsis;
        }

        .toggle-icon {
            display: inline-block;
            vertical-align: top;
            margin: 14px;
            font-size: 20px;
            transition: all $duration $cubicBezier;

            &.open {
                transform: rotate(90deg);
            }
        }
    }
}

.menu-submenu {
    line-height: 42px;
    font-size: 14px;
    overflow: hidden;
    transition: height $duration $cubicBezier;
    .submenu-link {
        position: relative;
        display: block;
        padding: 0 0 0 64px;
        color: $color;
        @include ellipsis;
        &:hover {
            background-color: #E8E9EF;
        }
        &.is-relative-active,
        &.active {
            color: $primaryColor;
            background-color: #E1ECFF;
            &::before {
                background-color: #fff;
            }
        }
        &:before {
            content: "";
            position: absolute;
            left: 29px;
            top: 19px;
            width: 4px;
            height: 4px;
            border-radius: 50%;
            background-color: #c4c6cc;
        }
    }
}

.copyright {
    margin: 17px 0 17px 22px;
    font-size: 20px;
    line-height: 28px;
    color: $color;
    transform-origin: left;
    transform: scale(.5);
    white-space: nowrap;
}

.nav-option {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 50px;
    line-height: 49px;
    border-top: 1px solid #DCDEE5;
    font-size: 0;
    color: #63656E;
    &:before {
        content: "";
        display: inline-block;
        height: 100%;
        width: 0;
        vertical-align: middle;
    }
    .nav-stick {
        display: inline-block;
        vertical-align: middle;
        width: 32px;
        height: 32px;
        margin: 0 0 0 13px;
        line-height: 32px;
        text-align: center;
        font-size: 14px;
        cursor: pointer;
        transition: transform $duration $cubicBezier;
        &:hover {
            opacity: .8;
        }
        &.sticked {
            transform: rotate(180deg);
        }
    }
}
</style>
