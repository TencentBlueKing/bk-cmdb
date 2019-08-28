<template>
    <nav class="nav-layout"
        :class="{ 'sticked': navStick }"
        @mouseenter="handleMouseEnter"
        @mouseleave="handleMouseLeave">
        <div class="nav-wrapper"
            :class="{ unfold: unfold, flexible: !navStick }">
            <div class="business-wrapper" v-if="showBusinessSelector">
                <cmdb-business-selector class="business-selector"></cmdb-business-selector>
            </div>
            <ul class="menu-list">
                <template v-for="(menu, index) in currentMenus">
                    <router-link class="menu-item is-link" tag="li" active-class="active"
                        v-if="menu.hasOwnProperty('route')"
                        :key="index"
                        :to="menu.route"
                        :title="$t(menu.i18n)">
                        <h3 class="menu-info clearfix">
                            <i :class="['menu-icon', menu.icon]"></i>
                            <span class="menu-name">{{$t(menu.i18n)}}</span>
                        </h3>
                    </router-link>
                    <li class="menu-item"
                        v-else
                        :key="index"
                        :class="{
                            'is-open': unfold && isMenuExpanded(menu)
                        }">
                        <h3 class="menu-info clearfix" @click="handleMenuClick(menu)">
                            <i :class="['menu-icon', menu.icon]"></i>
                            <span class="menu-name">{{$t(menu.i18n)}}</span>
                            <i class="toggle-icon bk-icon icon-angle-right"
                                v-if="menu.submenu && menu.submenu.length"
                                :class="{ open: unfold && isMenuExpanded(menu) }">
                            </i>
                        </h3>
                        <cmdb-collapse-transition>
                            <div class="menu-submenu"
                                v-if="menu.submenu && menu.submenu.length"
                                v-show="unfold && isMenuExpanded(menu)">
                                <router-link class="submenu-link"
                                    v-for="(submenu, submenuIndex) in menu.submenu"
                                    exact
                                    active-class="active"
                                    :class="{
                                        collection: menu.id === NAV_COLLECT
                                    }"
                                    :key="submenuIndex"
                                    :to="submenu.route"
                                    :title="$t(submenu.i18n)">
                                    {{$t(submenu.i18n)}}
                                    <i class="bk-icon icon-close"
                                        v-if="menu.id === NAV_COLLECT"
                                        @click.stop.prevent="handleDeleteCollection(submenu)">
                                    </i>
                                </router-link>
                            </div>
                        </cmdb-collapse-transition>
                    </li>
                </template>
            </ul>
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
    import { mapGetters } from 'vuex'
    import MENU_DICTIONARY from '@/dictionary/menu'
    import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                NAV_COLLECT: 'collect',
                routerLinkHeight: 42,
                timer: null,
                state: {},
                menus: []
            }
        },
        computed: {
            ...mapGetters(['navStick', 'navFold', 'admin']),
            ...mapGetters('menu', ['active']),
            ...mapGetters('userCustom', ['usercustom']),
            unfold () {
                return this.navStick || !this.navFold
            },
            owner () {
                return this.$route.matched[0].name
            },
            showBusinessSelector () {
                return this.owner === MENU_BUSINESS
            },
            currentMenus () {
                const target = MENU_DICTIONARY.find(menu => menu.id === this.owner)
                return (target && target.menu) || []
            }
        },
        created () {
            this.setDefaultExpand()
        },
        methods: {
            setDefaultExpand () {
                const expandedId = this.$tools.getValue(this.$route, 'meta.menu.parent')
                if (expandedId) {
                    this.$set(this.state, expandedId, { expanded: true })
                }
            },
            isMenuExpanded (menu) {
                if (this.state.hasOwnProperty(menu.id)) {
                    return this.state[menu.id].expanded
                }
                return false
            },
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
                const submenuCount = (menu.submenu || []).length
                return {
                    height: this.isMenuExpanded(menu) ? submenuCount * this.routerLinkHeight + 'px' : 0
                }
            },
            // 被点击的有对应的路由，则跳转
            checkPath (menu) {
                if (menu.path) {
                    this.$router.push({ path: menu.path })
                }
            },
            // 切换展开的分类
            toggleMenu (menu) {
                if (this.state.hasOwnProperty(menu.id)) {
                    this.state[menu.id].expanded = !this.state[menu.id].expanded
                } else {
                    this.$set(this.state, menu.id, { expanded: true })
                }
            },
            // 切换导航展开固定
            toggleNavStick () {
                this.$store.commit('setNavStatus', {
                    fold: !this.navFold,
                    stick: !this.navStick
                })
            },
            handleDeleteCollection (model) {
                const collectedModels = this.usercustom.collected_models
                this.$store.dispatch('userCustom/saveUsercustom', {
                    collected_models: collectedModels.filter(id => id !== model.id)
                })
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
    padding: 13px 0;
    height: 59px;
    border-bottom: 1px solid #DCDEE5;
    overflow: hidden;
    .business-selector {
        display: block;
        width: 240px;
        margin: 0 auto;
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
            background-color: #F0F1F5;
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
            background-color: #DCDEE5;
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
