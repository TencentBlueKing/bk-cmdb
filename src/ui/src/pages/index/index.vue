<template>
    <div class="index-wrapper">
        <v-search></v-search>
        <transition name="fade">
            <div class="classify-container" v-show="display.type === 'classification'">
                <ul class="classify-list">
                    <li class="classify-item clearfix" v-for="(classification, index) in authorizedClassifications" :key="index" @click="showModels(classification)">
                        <div 
                            :class="['classify-navigator', {navigated: isClassificationNavigated(classification)}]"
                            @click.stop="toggleClassifyNavVisible(classification)">
                            <span class="icon-cc-thumbtack"></span>
                            <span>{{isClassificationNavigated(classification) ? $t("Index['取消导航']") : $t("Index['添加导航']")}}</span>
                        </div>
                        <div class="classify-icon fl">
                            <i :class="classification['bk_classification_icon']"></i>
                        </div>
                        <div class="classify-name">
                            <h2  class="classify-name-text">{{classification['bk_classification_name']}}</h2>
                            <span class="classify-model-count">
                                {{$tc("Index['模型数量']",  classification['bk_objects'].length, {count: classification['bk_objects'].length})}}
                            </span>
                        </div>
                        <div class="classify-models">
                            <router-link exact class="classify-model-link"
                                v-for="(model, index) in classification['bk_objects']"
                                v-if="index < 4"
                                :key="index"
                                :title="model['bk_obj_name']"
                                :to="`/organization/${model['bk_obj_id']}`"
                                 @click.stop>
                                {{model['bk_obj_name']}}
                            </router-link>
                        </div>
                    </li>
                </ul>
            </div>
        </transition>
        <transition name="fade">
            <div class="model-container" v-show="display.type === 'model'">
                <bk-button class="model-return" type="primary" @click="display.type = 'classification'">{{$t("Index['返回上级']")}}</bk-button>
                <div class="model-list">
                    <router-link exact class="model-item" v-for="(model, index) in displayClassification['bk_objects']"
                        :key="index"
                        :to="`/organization/${model['bk_obj_id']}`">
                        <div :class="['classify-navigator', {navigated: isModelNavigated(model)}]" @click.prevent.stop="toggleModelNavVisible(model)">
                            <span class="icon-cc-thumbtack"></span>
                            <span>{{isModelNavigated(model) ? $t("Index['取消导航']") : $t("Index['添加导航']")}}</span>
                        </div>
                        <div class="model-icon fl">
                            <i :class="model['bk_obj_icon']"></i>
                        </div>
                        <div class="model-name">
                            <h3 class="model-name-text">{{model['bk_obj_name']}}</h3>
                            <span class="model-name-id">{{model['bk_obj_id']}}</span>
                        </div>
                    </router-link>
                </div>
            </div>
        </transition>
    </div>
</template>
<script>
    import vSearch from './children/search'
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                search: {
                    keyword: ''
                },
                display: {
                    type: 'classification',
                    displayClassificationId: null
                }
            }
        },
        computed: {
            ...mapGetters('usercustom', ['usercustom']),
            ...mapGetters('navigation', ['authorizedClassifications']),
            // 用户自定义导航内容
            userCustomNavigation () {
                return this.usercustom.navigation || {}
            },
            // 当前模型视图所对应的分类
            displayClassification () {
                if (this.display.displayClassificationId && this.display.type === 'model') {
                    return this.authorizedClassifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === this.display.displayClassificationId)
                }
                return {'bk_objects': []}
            }
        },
        methods: {
            // 显示视图
            showModels (classification) {
                this.display.type = 'model'
                this.display.displayClassificationId = classification['bk_classification_id']
            },
            // 检测分类是否已添加到导航
            isClassificationNavigated (classification) {
                let userCustomNavigation = this.userCustomNavigation
                if (!userCustomNavigation.hasOwnProperty(classification['bk_classification_id'])) {
                    return false
                }
                return !classification['bk_objects'].some(model => {
                    return !userCustomNavigation[classification['bk_classification_id']].includes(model['bk_obj_id'])
                })
            },
            isModelNavigated (model) {
                let userCustomNavigation = this.userCustomNavigation
                if (!userCustomNavigation.hasOwnProperty(model['bk_classification_id'])) {
                    return false
                }
                return userCustomNavigation[model['bk_classification_id']].includes(model['bk_obj_id'])
            },
            // 分类视图添加/取消导航切换
            toggleClassifyNavVisible (classification) {
                let navigation = JSON.parse(JSON.stringify(this.userCustomNavigation))
                let isFormerNavigated = this.isClassificationNavigated(classification)
                navigation[classification['bk_classification_id']] = isFormerNavigated ? [] : classification['bk_objects'].map(({bk_obj_id: bkObjId}) => bkObjId)
                this.$store.dispatch('usercustom/updateUserCustom', {navigation}).then(res => {
                    if (res.result) {
                        this.$alertMsg(isFormerNavigated ? this.$t("Index['取消导航成功']") : this.$t("Index['添加导航成功']"), 'success')
                    }
                })
            },
            // 模型视图添加/取消导航切换
            toggleModelNavVisible (model) {
                let navigation = JSON.parse(JSON.stringify(this.userCustomNavigation))
                let isFormerNavigated = this.isModelNavigated(model)
                let classificationId = model['bk_classification_id']
                navigation[classificationId] = navigation[classificationId] || []
                if (isFormerNavigated) {
                    navigation[classificationId] = navigation[classificationId].filter(bkObjId => bkObjId !== model['bk_obj_id'])
                } else {
                    navigation[classificationId].push(model['bk_obj_id'])
                }
                this.$store.dispatch('usercustom/updateUserCustom', {navigation}).then(res => {
                    if (res.result) {
                        this.$alertMsg(isFormerNavigated ? this.$t("Index['取消导航成功']") : this.$t("Index['添加导航成功']"), 'success')
                    }
                })
            }
        },
        components: {
            vSearch
        }
    }
</script>
<style lang="scss" scoped>
    $hoverColor: #3c96ff;
    .index-wrapper{
        height: 100%;
        padding: 20px 0;
        background-color: #e5eaef;
        overflow: hidden;
    }
    .classify-navigator{
        position: absolute;
        top: 10px;
        right: 10px;
        font-size: 12px;
        color: $textColor;
        cursor: pointer;
        .icon-cc-thumbtack{
            margin: 2px 0;
            display: inline-block;
        }
        &.navigated{
            color: $hoverColor;
        }
    }
    .classify-container{
        margin: 20px 0;
        padding: 0 10px;
        height: calc(100% - 154px);
        overflow-y: auto;
        overflow-x: hidden;
        @include scrollbar;
    }
    .classify-list{
        font-size: 0;
        margin: 0 auto;
        .classify-item{
            font-size: 14px;
            position: relative;
            display: inline-block;
            width: 390px;
            height: 200px;
            margin: 10px;
            background-color: #fff;
            box-shadow: 0px 1px 3px 0px rgba(0, 0, 0, 0.13);
            border-radius: 2px;
            transition: box-shadow .1s linear;
            &:hover{
                box-shadow: 0px 8px 12px 0px rgba(17, 28, 62, 0.2);
                .classify-icon{
                    color: $hoverColor;
                }
            }
            .classify-icon{
                width: 140px;
                height: 140px;
                line-height: 140px;
                font-size: 60px;
                text-align: center;
                transition: color .1s linear;
            }
            .classify-name{
                height: 140px;
                overflow: hidden;
                padding: 50px 10px 0 0;
            }
            .classify-name-text{
                font-size: 24px;
                font-weight: normal;
                color:#333948;
                margin: 0;
                @include ellipsis;
            }
            .classify-model-count{
                line-height: 24px;
                font-size: 14px;
                color:$textColor;
            }
            .classify-models{
                height: 60px;
                padding: 0 18px;
                background-color: #fafbfd;
                &:before{
                    @include verticalHack;
                }
            }
            .classify-model-link{
                display: inline-block;
                vertical-align: middle;
                max-width: 70px;
                margin: 0 18px;
                color: $textColor;
                @include ellipsis;
                &:hover{
                    color: $hoverColor;
                    text-decoration: underline;
                }
            }
        }
    }
    .model-container{
        position: relative;
        height: calc(100% - 114px);
        .model-return{
            margin: 0 20px 4px;
        }
    }
    .model-list{
        height: calc(100% - 40px);
        overflow-x: hidden;
        overflow-y: auto;
        font-size: 0;
        @include scrollbar;
        .model-item{
            position: relative;
            display: inline-block;
            width: 316px;
            height: 130px;
            margin: 10px 0 10px 20px;
            font-size: 14px;
            background-color: #ffffff;
            box-shadow: 0px 1px 3px 0px rgba(0, 0, 0, 0.13);
            border-radius: 2px;
            color: $textColor;
            transition: box-shadow .1s linear !important;
            &:hover{
                    box-shadow: 0px 8px 12px 0px rgba(17, 28, 62, 0.2);
            }
            .model-icon{
                width: 127px;
                height: 100%;
                font-size: 50px;
                text-align: center;
                background-color: #fafbfd;
                &:before{
                    @include verticalHack;
                }
            }
            .model-name{
                height: 100%;
                overflow: hidden;
                padding: 44px 10px 0 19px;
            }
            .model-name-text{
                margin: 0;
                font-size: 18px;
                font-weight: normal;
                color: #333948;
                @include ellipsis;
            }
            .model-name-id{
                line-height: 24px;
                color: $textColor;
                @include ellipsis;
            }
        }
    }
</style>
<style lang="scss">
    /* 导航栏收起时的九宫格宽度 */
    @media (max-width: 1679px){
        .content-control{
            .index-wrapper{
                .classify-list{
                    width: 1220px;
                    .classify-item{
                        width: 380px;
                    }
                }
            }
        }
    }
    @media (min-width: 1680px) and (max-width: 1759px) {
        .content-control{
            .index-wrapper{
                .classify-list{
                    width: 1620px;
                    .classify-item{
                        width: 380px;
                    }
                }
            }
        }
    }
    @media (min-width: 1760px) {
        .content-control{
            .index-wrapper{
                .classify-list{
                    width: 1700px;
                }
            }
        }
    }
    /* 导航展示时的九宫格宽度 */
    @media (max-width: 1439px) {
        .content-wrapper{
            .index-wrapper{
                .classify-list{
                    width: 860px;
                }
            }
        }
    }
    @media (min-width: 1440px) and (max-width: 1499px) {
        .content-wrapper{
            .index-wrapper{
                .classify-list{
                    width: 1220px;
                    .classify-item{
                        width: 380px;
                    }
                }
            }
        }
    }
    @media (min-width: 1500px) and (max-width: 1919px) {
        .content-wrapper{
            .index-wrapper{
                .classify-list{
                    width: 1280px;
                }
            }
        }
    }
</style>