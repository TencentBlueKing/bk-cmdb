<template>
    <div :class="['nav-wrapper', {fold: fold}]">
        <div :class="['nav-logo', language]" @click="turnToIndex"></div>
        <ul class="nav-list">
            <li v-for="(classification, index) in customNavigation"
                :key="index"
                :class="['nav-item', {open: classification['id'] === unfoldClassificationId}]">
                <router-link exact class="nav-classification index" v-if="classification['id'] === 'bk_index'"
                    :to="classification['path']"
                    :title="classification['i18n'] ? $t(classification['i18n']) : classification['name']">
                    <div @click="toggleClassification(classification)">
                        <i :class="['nav-classification-icon', classification['icon']]"></i>
                        <span class="nav-classification-name" :title="classification['i18n'] ? $t(classification['i18n']) : classification['name']">
                            {{classification['i18n'] ? $t(classification['i18n']) : classification['name']}}
                        </span>
                    </div>
                </router-link>
                <div class="nav-classification" v-else @click="toggleClassification(classification)">
                    <i :class="['nav-classification-icon', classification['icon']]"></i>
                    <span class="nav-classification-name" :title="classification['i18n'] ? $t(classification['i18n']) : classification['name']">
                        {{classification['i18n'] ? $t(classification['i18n']) : classification['name']}}
                    </span>
                    <i class="bk-icon icon-angle-down" v-show="classification['children'].length"></i>
                </div>
                <div class="nav-classification-model" v-if="classification['children'].length"
                    :style="{height: classification['id'] === unfoldClassificationId ? `${36 * classification['children'].length}px` : 0}">
                    <router-link class="nav-classification-link" exact v-for="(model, index) in classification['children']"
                        :key="index"
                        :to="model['path']"
                        :title="model['i18n'] ? $t(model['i18n']) : model['name']">
                        {{model['i18n'] ? $t(model['i18n']) : model['name']}}
                    </router-link>
                </div>
            </li>
        </ul>
    </div>
</template>
<script>
    import {mapGetters} from 'vuex'
    export default {
        data () {
            return {
                unfoldClassificationId: 'bk_index'
            }
        },
        computed: {
            ...mapGetters('navigation', ['fold', 'customNavigation']),
            language () {
                return this.$i18n.locale
            }
        },
        watch: {
            '$route.path' (newPath, oldPath) {
                this.setUnfoldClassificationId()
            }
        },
        methods: {
            turnToIndex () {
                this.$router.push('/')
            },
            setUnfoldClassificationId () {
                let path = this.$route.path
                if (path === '/') {
                    this.unfoldClassificationId = 'bk_index'
                    return
                }
                let activeClassification = this.customNavigation.find(classification => {
                    return classification.children.some(model => model['path'] === path)
                })
                this.unfoldClassificationId = activeClassification ? activeClassification['id'] : null
            },
            toggleClassification (classification) {
                if (classification['id'] === 'bk_index') {
                    this.unfoldClassificationId = classification['id']
                } else {
                    this.unfoldClassificationId = this.unfoldClassificationId === classification['id'] ? null : classification['id']
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    $navTextColor: #c9d0e6;
    $navActiveColor: #283556;
    .nav-wrapper{
        height: 100%;
        width: 220px;
        background-color: #334162;
        transition: all .5s;
        &.fold{
            width: 60px;
        }
    }
    .nav-logo{
        display: block;
        height: 62px;
        background: transparent center center no-repeat;
        background-size: 173px 31px;
        &.en{
            background-image: url(../../common/images/nav_title.png);
        }
        &.zh_CN{
            background-image: url(../../common/images/nav-title-zh.png);
        }
    }

    .nav-list{
        height: calc(100% - 62px - 86px - 30px); /* 62px : logo高度; 86px: 版权高度; 30px: 底部留白间距 */
        overflow: auto;
        color: $navTextColor;
        @include scrollbar;
        .nav-item{
            cursor: pointer;
            &.open{
                background-color: #2f3c5d;
                .icon-angle-down{
                    transform: rotate(0deg);
                }
            }
        }
    }
    .nav-classification{
        font-size: 14px;
        height: 48px;
        line-height: 48px;
        color: $navTextColor;
        font-size: 0;
        &:hover{
            background-color: $navActiveColor;
        }
        &.active{
            color: #fff;
            background-color: $navActiveColor;
        }
        .nav-classification-icon{
            display: inline-block;
            vertical-align: middle;
            margin: 0 4px 0 38px;
            font-size: 16px;
        }
        .nav-classification-name{
            display: inline-block;
            vertical-align: middle;
            font-weight: 700;
            width: 130px;
            padding: 0 10px 0 15px;
            font-size: 14px;
            @include ellipsis;
        }
        .icon-angle-down{
            display: inline-block;
            vertical-align: middle;
            font-size: 12px;
            transition: transform .5s cubic-bezier(.23, 1, .23, 1);
            transform: rotate(90deg);
        }
    }
    .nav-classification-model{
        overflow: hidden;
        transition: height .5s cubic-bezier(.23, 1, .23, 1);
        .nav-classification-link{
            display: block;
            height: 36px;
            line-height: 36px;
            padding: 0 10px 0 73px;
            color: $navTextColor;
            font-size: 14px;
            @include ellipsis;
            &.router-link-exact-active,
            &:hover{
                color: #fff;
                background-color: $navActiveColor;
            }
        }
    }
</style>