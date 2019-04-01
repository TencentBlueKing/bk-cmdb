<template>
    <nav class="nav-layout"
        :class="{'sticked': navStick, 'admin-view': isAdminView}"
        @mouseenter="handleMouseEnter"
        @mouseleave="handleMouseLeave">
        <div class="nav-wrapper"
            :class="{unfold: unfold, flexible: !navStick}">
            <div class="logo" @click="$router.push({name: 'index'})">
                <span class="logo-text">
                    {{$t('Nav["蓝鲸配置平台"]')}}
                </span>
                <span class="logo-tag" v-if="isAdminView" :title="$t('Nav[\'后台管理标题\']')">
                    {{$t('Nav["后台管理"]')}}
                </span>
            </div>
            <ul class="menu-list">
                <li class="menu-item"
                    v-for="(menu, index) in menus"
                    :class="{
                        active: activeMenu === menu.id,
                        'is-open': openedMenu === menu.id,
                        'is-link': menu.path
                    }">
                    <h3 class="menu-info clearfix"
                        :class="{'menu-link': menu.path}"
                        @click="handleMenuClick(menu)">
                        <i :class="['menu-icon', menu.icon]"></i>
                        <span class="menu-name">{{menu.i18n ? $t(menu.i18n) : menu.name}}</span>
                        <i class="toggle-icon bk-icon icon-angle-right"
                            v-if="menu.submenu && menu.submenu.length"
                            :class="{open: menu.id === openedMenu}">
                        </i>
                    </h3>
                    <div class="menu-submenu" 
                        v-if="menu.submenu && menu.submenu.length"
                        :style="getMenuModelsStyle(menu)">
                        <router-link class="submenu-link" exact
                            v-for="(submenu, submenuIndex) in menu.submenu"
                            :class="{
                                active: activeMenu === submenu.id,
                                collection: menu.id === NAV_COLLECT
                            }"
                            :key="submenuIndex"
                            :to="submenu.path"
                            :title="submenu.i18n ? $t(submenu.i18n) : submenu.name">
                            {{submenu.i18n ? $t(submenu.i18n) : submenu.name}}
                            <i class="bk-icon icon-close"
                                v-if="menu.id === NAV_COLLECT"
                                @click.stop.prevent="handleDeleteCollection(submenu)">
                            </i>
                        </router-link>
                    </div>
                </li>
            </ul>
            <div class="nav-option">
                <i class="nav-stick icon icon-cc-nav-toggle"
                    :class="{
                        sticked: navStick
                    }"
                    :title="navStick ? $t('Index[\'收起导航\']') : $t('Index[\'固定导航\']')"
                    @click="toggleNavStick">
                </i>
            </div>
        </div>
    </nav>
</template>
<script>
import { mapGetters } from 'vuex'
import MENU, { NAV_COLLECT } from '@/dictionary/menu'
import MODEL_ROUTER_CONFIG, { GET_MODEL_PATH } from '@/views/general-model/router.config'
export default {
    data () {
        return {
            NAV_COLLECT,
            routerLinkHeight: 42,
            openedMenu: null,
            activeMenu: null,
            timer: null
        }
    },
    computed: {
        ...mapGetters(['navStick', 'navFold', 'admin', 'isAdminView']),
        ...mapGetters('userCustom', ['usercustom', 'classifyNavigationKey']),
        collectMenus () {
            const collectMenus = []
            const collectedModelIds = this.usercustom[this.classifyNavigationKey] || []
            const getModelById = this.$store.getters['objectModelClassify/getModelById']
            collectedModelIds.forEach((modelId, index) => {
                const model = getModelById(modelId)
                if (model) {
                    collectMenus.push({
                        id: model.bk_obj_id,
                        name: model.bk_obj_name,
                        path: GET_MODEL_PATH(modelId),
                        order: index
                    })
                }
            })
            return collectMenus
        },
        menus () {
            const menus = this.$tools.clone(MENU)
            const routes = this.$router.options.routes
            const isAuthorized = this.$store.getters['auth/isAuthorized']
            routes.forEach(route => {
                const meta = (route.meta || {})
                const auth = meta.auth || {}
                const menu = meta.menu
                if (menu) {
                    const authorized = auth.view ? isAuthorized(...auth.view.split('.'), true) : true
                    if (authorized) {
                        if (menu.parent) {
                            const parent = menus.find(parent => parent.id === menu.parent) || {}
                            const submenu = parent.submenu || []
                            submenu.push(menu)
                        } else {
                            const parent = menus.find(parent => parent.id === menu.id) || {}
                            Object.assign(parent, menu)
                        }
                    }
                }
            })
            const collectMenu = menus.find(menu => menu.id === NAV_COLLECT) || {}
            const collectSubmenu = collectMenu.submenu || []
            Array.prototype.push.apply(collectSubmenu, this.collectMenus)
            const availableMenus = menus.filter(menu => {
                return menu.path ||
                    (Array.isArray(menu.submenu) && menu.submenu.length)
            })
            availableMenus.forEach(menu => {
                if (Array.isArray(menu.submenu)) {
                    menu.submenu.sort((prev, next) => prev.order - next.order)
                }
            })
            return availableMenus
        },
        unfold () {
            return this.navStick || !this.navFold
        },
        // 展开的分类子菜单高度
        openedMenuHeight () {
            const openedMenu = this.menus.find(menu => menu.id === this.openedMenu)
            if (openedMenu) {
                const submenuCount = (openedMenu.submenu || []).length
                return submenuCount * this.routerLinkHeight
            }
            return 0
        }
    },
    methods: {
        handleMouseEnter () {
            if (this.timer) {
                clearTimeout(this.timer)
            }
            this.$store.commit('setNavStatus', { fold: false })
        },
        handleMouseLeave () {
            this.timer = setTimeout(() => {
                this.$store.commit('setNavStatus', { fold: true })
            }, 300)
        },
        // 分类点击事件
        handleMenuClick (menu) {
            this.checkPath(menu)
            this.toggleMenu(menu)
        },
        getMenuModelsStyle (menu) {
            return {
                height: (this.unfold && menu.id === this.openedMenu) ? this.openedMenuHeight + 'px' : 0
            }
        },
        // 被点击的有对应的路由，则跳转
        checkPath (menu) {
            if (menu.path) {
                this.$router.push({path: menu.path})
            }
        },
        // 切换展开的分类
        toggleMenu (menu) {
            this.openedMenu = menu.id === this.openedMenu ? null : menu.id
        },
        // 切换导航展开固定
        toggleNavStick () {
            this.$store.commit('setNavStatus', {
                fold: !this.navFold,
                stick: !this.navStick
            })
        },
        handleDeleteCollection (model) {
            const customNavigation = this.usercustom[this.classifyNavigationKey] || []
            this.$store.dispatch('userCustom/saveUsercustom', {
                [this.classifyNavigationKey]: customNavigation.filter(id => id !== model.id)
            })
        }
    }
}
</script>
<style lang="scss" scoped>
$cubicBezier: cubic-bezier(0.4, 0, 0.2, 1);
$duration: 0.2s;
$color: #979ba5;

.nav-layout {
    position: relative;
    width: 60px;
    height: 100%;
    transition: width $duration $cubicBezier;

    &.sticked {
        width: 260px;
    }

    .nav-wrapper {
        position: relative;
        width: 100%;
        height: 100%;
        background: #182132;
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

.logo {
    height: 60px;
    padding: 0 0 0 64px;
    border-bottom: 1px solid rgba(255, 255, 255, .05);
    background-color: #182132;
    line-height: 59px;
    color: #fff;
    font-size: 0;
    font-weight: bold;
    white-space: nowrap;
    overflow: hidden;
    cursor: pointer;
    background: url('../../assets/images/logo.svg') no-repeat;
    background-position: 16px 14px;
    .logo-text {
        display: inline-block;
        vertical-align: middle;
        font-size: 18px;
    }
    .logo-tag {
        display: inline-block;
        padding: 0 8px;
        margin: 0 0 0 4px;
        vertical-align: middle;
        border-radius: 2px;
        color: #282b41;
        font-size: 20px;
        font-weight: normal;
        line-height: 32px;
        background: #18b48a;
        transform: scale(0.5);
        transform-origin: left center;
    }
}

.menu-list {
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
        transition: background-color $duration $cubicBezier;

        &.is-open {
            background-color: #202a3c;
        }
        &.active.is-link {
            background-color: #3a84ff;
            .menu-icon,
            .menu-name {
                color: #fff;
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
            color: rgba(255, 255, 255, .8);
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
            margin: 18px;
            font-size: 12px;
            transition: all $duration $cubicBezier;

            &.open {
                transform: rotate(90deg);
            }
        }
    }
}

.menu-submenu {
    height: 0;
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
        &.collection {
            padding-right: 50px;
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
                    color: #fff;
                }
                &:before {
                    display: block;
                    transform: scale(.5);
                }
            }
            &:hover {
                .icon-close {
                    display: block;
                }
            }
        }
        &:hover {
            color: #fff;
            background-color: #303c4c;
        }
        &.active {
            color: #fff;
            background-color: #3a84ff;
        }
        &:before {
            content: "";
            position: absolute;
            left: 29px;
            top: 19px;
            width: 4px;
            height: 4px;
            border-radius: 50%;
            background-color: currentColor;
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
    height: 55px;
    line-height: 54px;
    border-top: 1px solid rgba(255, 255, 255, .05);
    font-size: 0;
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
            color: #fff;
        }
        &.sticked {
            transform: rotate(180deg);
        }
    }
}
</style>
