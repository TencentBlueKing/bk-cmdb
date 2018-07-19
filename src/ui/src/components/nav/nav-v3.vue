<template>
    <div class="nav-wrapper"
        :class="{'sticked': navStick}"
        @mouseenter="handleMouseEnter"
        @mouseleave="handleMouseLeave">
        <div class="nav-layout" :class="{unfold: unfold, flexible: !navStick}">
            <div class="logo">
                <img src="@/common/svg/logo-cn.svg" alt="蓝鲸配置平台"
                    v-if="$i18n.locale === 'zh_CN'"
                    @click="$router.push('/index')">
                <img src="@/common/svg/logo-en.svg" alt="Configuration System"
                    v-else
                    @click="$router.push('/index')">
            </div>
            <ul class="classify-list">
                <li class="classify-item"
                    v-for="(classify, index) in [...staticClassify, ...customClassify]"
                    :class="{active: activeClassifyId === classify.id}">
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
                        <router-link exact class="model-link"
                            v-for="(model, modelIndex) in classify.children"
                            :key="modelIndex"
                            :to="model.path"
                            :title="model.i18n ? $t(model.i18n) : model.name">
                            {{model.i18n ? $t(model.i18n) : model.name}}
                        </router-link>
                    </div>
                    <i class="classify-corner"
                        v-show="!unfold"
                        v-if="classify.children && classify.children.length">
                    </i>
                </li>
            </ul>
            <i class="nav-stick icon-cc-nav-stick"
                v-show="unfold"
                :class="{'sticked': navStick}"
                :title="navStick ? $t('Index[\'收起导航\']') : $t('Index[\'固定导航\']')"
                @click="toggleNavStick">
            </i>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                staticClassifyId: ['bk_index', 'bk_host_manage', 'bk_organization', 'bk_back_config'],
                routerLinkHeight: 36,
                openedClassify: 'bk_index',
                timer: null
            }
        },
        computed: {
            ...mapGetters('navigation', ['fold', 'navStick', 'authorizedNavigation']),
            ...mapGetters('usercustom', ['usercustom', 'classifyNavigationKey']),
            unfold () {
                return this.navStick || !this.fold
            },
            customNavigation () {
                return this.usercustom[this.classifyNavigationKey] || []
            },
            allModels () {
                let models = []
                this.authorizedNavigation.forEach(classify => {
                    if (classify.children && classify.children.length) {
                        classify.children.forEach(model => {
                            models.push(model)
                        })
                    }
                })
                return models
            },
            // 固定到前面的分类
            staticClassify () {
                return this.authorizedNavigation.filter(classify => {
                    if (classify.id === 'bk_index') {
                        return true
                    }
                    return this.staticClassifyId.includes(classify.id) && classify.children.length
                })
            },
            // 用户自定义到导航的分类/模型
            customClassify () {
                let customClassify = []
                this.customNavigation.forEach(modelId => {
                    const classifyModel = this.allModels.find(model => model.id === modelId)
                    if (classifyModel) {
                        customClassify.push({
                            ...classifyModel,
                            children: []
                        })
                    }
                })
                return customClassify
            },
            // 当前导航对应的分类
            activeClassify () {
                const path = this.$route.fullPath
                return [...this.staticClassify, ...this.customClassify].find(classify => {
                    if (classify.hasOwnProperty('path')) {
                        return classify.path === path
                    } else if (classify.children && classify.children.length) {
                        return classify.children.some(model => model.path === path)
                    }
                })
            },
            // 当前导航对应的分类ID
            activeClassifyId () {
                return this.activeClassify ? this.activeClassify.id : this.$route.fullPath === '/index' ? 'bk_index' : null
            },
            // 当前导航对应的模型
            activeModel () {
                if (this.activeClassify) {
                    const path = this.$route.fullPath
                    return this.activeClassify.children.find(model => model.path === path)
                }
                return null
            },
            // 展开的分类子菜单高度
            openedClassifyHeight () {
                const openedClassify = this.authorizedNavigation.find(classify => classify.id === this.openedClassify)
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
            handleMouseEnter () {
                if (this.timer) {
                    clearTimeout(this.timer)
                }
                this.$store.commit('navigation/updateNavFold', false)
            },
            handleMouseLeave () {
                this.timer = setTimeout(() => {
                    this.$store.commit('navigation/updateNavFold', true)
                }, 300)
            },
            // 分类点击事件
            handleClassifyClick (classify) {
                this.checkPath(classify)
                this.toggleClassify(classify)
            },
            getClassifyModelsStyle (classify) {
                return {
                    height: (this.unfold && classify.id === this.openedClassify) ? this.openedClassifyHeight + 'px' : 0
                }
            },
            // 被点击的有对应的路由，则跳转
            checkPath (classify) {
                if (classify.hasOwnProperty('path')) {
                    this.$router.push(classify.path)
                }
            },
            // 切换展开的分类
            toggleClassify (classify) {
                this.openedClassify = classify.id === this.openedClassify ? null : classify.id
            },
            // 切换导航展开固定
            toggleNavStick () {
                this.$store.commit('navigation/updateNavFold', !this.fold)
                this.$store.commit('navigation/updateNavStick', !this.navStick)
            }
        }
    }
</script>

<style lang="scss" scoped>
    $cubicBezier: cubic-bezier(0.4, 0, 0.2, 1);
    $duration: 0.2s;
    .nav-wrapper{
        position: relative;
        width: 60px;
        height: 100%;
        transition: width $duration $cubicBezier;
        z-index: 1201;
        &.sticked{
            width: 240px;
        }
        .nav-layout{
            position: relative;
            width: 100%;
            height: 100%;
            background:rgba(46,50,59,1);
            transition: width $duration $cubicBezier;
            &.unfold{
                width: 240px;
            }
            &.unfold.flexible:after{
                content: "";
                position: absolute;
                width: 15px;
                height: 100%;
                left: 100%;
                top: 0;
            }
            .nav-stick{
                position: absolute;
                bottom: 9px;
                right: 13px;
                width: 32px;
                height: 32px;
                padding: 10px 9px;
                border-radius: 50%;
                transition: transform $duration $cubicBezier;
                transform: scale(0.8333) rotate(180deg);
                font-size: 12px;
                cursor: pointer;
                &.sticked{
                    transform: scale(0.8333);
                }
                &:hover{
                    background-color: #3a4156;
                }
            }
        }
    }
    .logo{
        height: 60px;
        padding: 12px 0 12px 15px;
        background-color: #4c84ff;
        color: #fff;
        overflow: hidden;
        cursor: pointer;
        img{
            height: 36px;
        }
    }
    .classify-corner{
        position: absolute;
        bottom: 0;
        right: -3px;
        width: 0;
        height: 0;
        border: 5px solid transparent;
        border-left-color: #737987;
        transform: rotate(45deg);
    }
    .classify-list{
        height: calc(100% - 120px); // 上下各去掉60px
        overflow-y: auto;
        overflow-x: hidden;
        white-space: nowrap;
        &::-webkit-scrollbar {
            width: 5px;
            height: 5px;
            &-thumb {
                border-radius: 20px;
                background: rgba(165, 165, 165, .3);
                box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
            }
        }
        .classify-item{
            position: relative;
            &:hover{
                background-color: rgba(21, 30, 42, .75);
            }
            &.active{
                background-color: #151e2a;
                .classify-icon{
                    color: #fff;
                }
            }
            .classify-info{
                margin: 0;
                padding: 0;
                height: 48px;
                line-height: 48px;
                font-weight: bold;
                white-space: nowrap;
                font-size: 0;
                cursor: pointer;
            }
            .classify-icon{
                display: inline-block;
                vertical-align: top;
                margin: 12px 18px;
                font-size: 24px;
                color: #737987;
            }
            .classify-name{
                display: inline-block;
                width: calc(100% - 120px);
                vertical-align: top;
                font-size: 14px;
                @include ellipsis;
            }
            .toggle-icon{
                display: inline-block;
                vertical-align: top;
                margin: 18px;
                font-size: 12px;
                transition: all $duration $cubicBezier;
                &.open{
                    transform: rotate(90deg);
                }
            }
        }
    }
    .classify-models{
        height: 0;
        padding: 0 0 0 63px;
        // padding: 4px 0 4px 63px;
        line-height: 36px;
        font-size: 12px;
        overflow: hidden;
        transition: height $duration $cubicBezier;
        .model-link{
            display: block;
            color: #c3cdd7;
            @include ellipsis;
            &:hover{
                color: #3c96ff;
            }
            &.active{
                color: #0082ff;
                font-weight: bold;
            }
        }
    }
</style>