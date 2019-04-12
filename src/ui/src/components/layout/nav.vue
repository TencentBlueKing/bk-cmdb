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
            <ul class="classify-list">
                <li class="classify-item"
                    v-for="(classify, index) in navigations"
                    :class="{
                        active: isClassifyActive(classify),
                        'is-open': openedClassify === classify.id,
                        'is-link': classify.hasOwnProperty('path')
                    }">
                    <h3 class="classify-info clearfix"
                        :class="{'classify-link': classify.hasOwnProperty('path')}"
                        @click="handleClassifyClick(classify)">
                        <i :class="['classify-icon', classify.icon]"></i>
                        <span class="classify-name">{{classify.i18n ? $t(classify.i18n) : classify.name}}</span>
                        <i class="toggle-icon bk-icon icon-angle-right"
                            v-if="classify.children && classify.children.length"
                            :class="{open: classify.id === openedClassify}">
                        </i>
                    </h3>
                    <div class="classify-models" 
                        v-if="classify.children && classify.children.length"
                        :style="getClassifyModelsStyle(classify)">
                        <router-link class="model-link" exact
                            v-for="(model, modelIndex) in classify.children"
                            :class="{
                                active: isRouterActive(model),
                                collection: classify.id === 'bk_collection'
                            }"
                            :key="modelIndex"
                            :to="model.path"
                            :title="model.i18n ? $t(model.i18n) : model.name">
                            {{model.i18n ? $t(model.i18n) : model.name}}
                            <i class="bk-icon icon-close"
                                v-if="classify.id === 'bk_collection'"
                                @click.stop.prevent="handleDeleteCollection(model)">
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
export default {
    data () {
        return {
            routerLinkHeight: 42,
            openedClassify: null,
            timer: null
        }
    },
    computed: {
        ...mapGetters(['navStick', 'navFold', 'admin', 'isAdminView']),
        ...mapGetters('objectModelClassify', ['classifications', 'authorizedNavigation', 'staticClassifyId']),
        ...mapGetters('userCustom', ['usercustom', 'classifyNavigationKey']),
        fixedClassifyId () {
            return [...this.staticClassifyId, 'bk_organization']
        },
        unfold () {
            return this.navStick || !this.navFold
        },
        // 当前导航对应的分类ID
        activeClassifyId () {
            return this.$classify.classificationId
        },
        navigations () {
            const navigations = this.$tools.clone(this.authorizedNavigation)
            if (this.admin) {
                if (this.isAdminView) {
                    return navigations.filter(classify => classify.classificationId !== 'bk_business_resource')
                }
                return navigations
            }
            navigations.forEach(classify => {
                classify.children = classify.children.filter(child => child.authorized)
            })
            return navigations.filter(classify => {
                return (classify.hasOwnProperty('path') && classify.authorized) || classify.children.length
            })
        },
        // 展开的分类子菜单高度
        openedClassifyHeight () {
            const openedClassify = this.navigations.find(classify => classify.id === this.openedClassify)
            if (openedClassify) {
                const modelsCount = openedClassify.children.length
                return modelsCount * this.routerLinkHeight
            }
            return 0
        }
    },
    watch: {
        activeClassifyId (id) {
            this.openedClassify = id
        }
    },
    methods: {
        isClassifyActive (classify) {
            const path = this.$route.meta.relative || this.$route.query.relative || this.$route.path
            return classify.path === path || this.activeClassifyId === classify.id
        },
        isActiveClosed (classify) {
            return this.activeClassifyId === classify.id &&
                (this.openedClassify === null || !this.unfold)
        },
        isAvailableClassify (classify) {
            if (classify.hasOwnProperty('path')) {
                return true
            }
            return classify.children.some(sub => sub.authorized)
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
        handleClassifyClick (classify) {
            this.checkPath(classify)
            this.toggleClassify(classify)
        },
        isRouterActive (model) {
            const path = this.$route.meta.relative || this.$route.query.relative || this.$route.path
            return model.path === path
        },
        getClassifyModelsStyle (classify) {
            return {
                height: (this.unfold && classify.id === this.openedClassify) ? this.openedClassifyHeight + 'px' : 0
            }
        },
        // 被点击的有对应的路由，则跳转
        checkPath (classify) {
            if (classify.hasOwnProperty('path')) {
                this.$router.push({path: classify.path})
            }
        },
        // 切换展开的分类
        toggleClassify (classify) {
            this.openedClassify = classify.id === this.openedClassify ? null : classify.id
        },
        // 切换导航展开固定
        toggleNavStick () {
            this.$store.commit('setNavStatus', {
                fold: !this.navFold,
                stick: !this.navStick
            })
        },
        handleDeleteCollection (model) {
            if (['biz', 'resource'].includes(model.id)) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [`is_${model.id}_collected`]: false
                })
            } else {
                const customNavigation = this.usercustom[this.classifyNavigationKey] || []
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.classifyNavigationKey]: customNavigation.filter(id => id !== model.id)
                })
            }
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

.classify-list {
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

    .classify-item {
        position: relative;
        transition: background-color $duration $cubicBezier;

        &.is-open {
            background-color: #202a3c;
        }
        &.active.is-link {
            background-color: #3a84ff;
            .classify-icon,
            .classify-name {
                color: #fff;
            }
        }
        .classify-info {
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

        .classify-icon {
            display: inline-block;
            vertical-align: top;
            margin: 13px 26px 13px 22px;
            font-size: 16px;
            color: rgba(255, 255, 255, .8);
        }

        .classify-name {
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

.classify-models {
    height: 0;
    line-height: 42px;
    font-size: 14px;
    overflow: hidden;
    transition: height $duration $cubicBezier;
    .model-link {
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
