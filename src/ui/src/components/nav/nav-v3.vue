<template>
    <div class="nav-wrapper" :class="{'sticked': navStick}">
        <div class="nav-layout">
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
                    v-for="(classify, index) in [...staticClassify, ...customClassify]" @click="checkPath(classify)"
                    :class="{active: activeClassifyId === classify.id}">
                    <h3 class="classify-info clearfix" @click="toggleClassify(classify)">
                        <i :class="['classify-icon', classify.icon]"></i>
                        <span class="classify-name">{{classify.i18n ? $t(classify.i18n) : classify.name}}</span>
                        <i class="toggle-icon bk-icon icon-angle-right"
                            v-if="classify.children && classify.children.length"
                            :class="{open: classify.id === openedClassify}">
                        </i>
                    </h3>
                    <div :class="['classify-models', {'open': classify.id === openedClassify}]"
                        :style="{height: classify.id === openedClassify ? openedClassifyHeight + 'px' : 0}">
                        <router-link exact class="model-link"
                            v-for="(model, modelIndex) in classify.children"
                            :key="modelIndex"
                            :to="model.path"
                            :title="model.i18n ? $t(model.i18n) : model.name">
                            {{model.i18n ? $t(model.i18n) : model.name}}
                        </router-link>
                    </div>
                    <i class="classify-corner" v-if="classify.children && classify.children.length"></i>
                </li>
            </ul>
            <i class="nav-stick bk-icon icon-back"
                :class="{'sticked': navStick}"
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
                routerLinkHeight: 24,
                openedClassify: 'bk_index'
            }
        },
        computed: {
            ...mapGetters('navigation', ['navStick', 'authorizedNavigation']),
            ...mapGetters('usercustom', ['usercustom', 'classifyNavigationKey', 'classifyModelSequenceKey']),
            customNavigation () {
                return this.usercustom[this.classifyNavigationKey] || []
            },
            classifyModelSequence () {
                return this.usercustom[this.classifyModelSequenceKey] || {}
            },
            staticClassify () {
                const classifies = this.$deepClone(this.authorizedNavigation.filter(classify => this.staticClassifyId.includes(classify.id)))
                classifies.forEach((classify, index) => {
                    if (this.classifyModelSequence.hasOwnProperty(classify.id) && !['bk_host_manage'].includes(classify.id)) {
                        classify['children'].sort((modelA, modelB) => {
                            return this.getModelSequence(classify, modelA) - this.getModelSequence(classify, modelB)
                        })
                    }
                })
                return classifies
            },
            customClassify () {
                const classifies = this.$deepClone(this.authorizedNavigation.filter(classify => this.customNavigation.includes(classify.id)))
                classifies.forEach((classify, index) => {
                    if (this.classifyModelSequence.hasOwnProperty(classify.id)) {
                        classify['children'].sort((modelA, modelB) => {
                            return this.getModelSequence(classify, modelA) - this.getModelSequence(classify, modelB)
                        })
                    }
                })
                return classifies
            },
            activeClassify () {
                const path = this.$route.fullPath
                return this.authorizedNavigation.find(classify => classify.children.some(model => model.path === path))
            },
            activeClassifyId () {
                return this.activeClassify ? this.activeClassify.id : this.$route.fullPath === '/index' ? 'bk_index' : null
            },
            activeModel () {
                if (this.activeClassify) {
                    const path = this.$route.fullPath
                    return this.activeClassify.children.find(model => model.path === path)
                }
                return null
            },
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
            getModelSequence (classify, model) {
                if (this.classifyModelSequence.hasOwnProperty(classify.id)) {
                    const sequence = this.classifyModelSequence[classify.id]
                    const modelSequence = sequence.indexOf(model.id)
                    return modelSequence === -1 ? classify.children.length : modelSequence
                }
                return classify.children.length
            },
            checkPath (classify) {
                if (classify.hasOwnProperty('path')) {
                    this.$router.push(classify.path)
                }
            },
            toggleNavStick () {
                this.$store.commit('navigation/updateNavStick', !this.navStick)
            },
            toggleClassify (classify) {
                this.openedClassify = classify.id
            }
        }
    }
</script>

<style lang="scss" scoped>
    .nav-wrapper{
        position: relative;
        width: 60px;
        height: 100%;
        transition: width .1s ease-in;
        z-index: 1201;
        &.sticked{
            width: 240px;
            .classify-corner{
                display: none;
            }
            .classify-models{
                display: block;
            }
        }
        &:hover{
            .nav-layout{
                width: 240px;
                .nav-stick{
                    opacity: 1;
                }
            }
            .classify-corner{
                display: none;
            }
            .classify-models{
                display: block;
            }
        }
        .nav-layout{
            position: relative;
            width: 100%;
            height: 100%;
            background:rgba(46,50,59,1);
            transition: width .1s ease-in;
            .nav-stick{
                position: absolute;
                bottom: 17px;
                right: 21px;
                opacity: 0;
                transition: opacity .1s ease-in;
            }
        }
    }
    .logo{
        height: 60px;
        padding: 12px 0 12px 15px;
        background-color: #4c84ff;
        color: #fff;
        overflow: hidden;
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
        @include scrollbar;
        .classify-item{
            position: relative;
            &.active{
                .classify-info{
                    color: #fff;
                    background-color: #3a4156;
                }
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
                transition: all .1s ease-in;
                &.open{
                    transform: rotate(90deg);
                }
            }
        }
    }
    .classify-models{
        display: none;
        height: 0;
        padding: 0 0 0 63px;
        // padding: 4px 0 4px 63px;
        line-height: 24px;
        font-size: 12px;
        overflow: hidden;
        transition: height .1s ease-in;
        .model-link{
            display: block;
            color: #c3cdd7;
            @include ellipsis;
            &.active{
                color: #fff;
                font-weight: bold;
            }
            &:hover{
                color: #fff;
            }
        }
    }
</style>