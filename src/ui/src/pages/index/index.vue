<template>
    <div class="index-wrapper">
        <v-search></v-search>
        <transition name="fade">
            <div class="classify-container" v-show="display.type === 'classification'">
                <ul class="classify-list">
                    <li class="classify-item" v-for="(classification, index) in sortedClassifications" :key="index">
                        <i v-if="!isHostManage(classification)"
                            :class="['icon-cc-star', {navigated: isClassificationNavigated(classification)}]"
                            @click.stop="toggleClassifyNavVisible(classification)">
                        </i>
                        <div :class="['classify-info-layout', `classify-info-layout-${classification['bk_classification_id']}`]" @click="showModels(classification)">
                            <span class="classify-info">
                                <i :class="['classify-info-icon', classification['bk_classification_icon']]"></i>
                                <span class="classify-info-name">{{classification['bk_classification_name']}}</span>
                            </span>
                        </div>
                        <ul :class="['classify-models-list', {'has-more': classification['bk_objects'].length > 3}]">
                            <li class="classify-model-item" 
                                v-for="(model, index) in classification['bk_objects'].filter((model, index) => index <= 4)">
                                <router-link exact
                                    :class="['classify-model-link', `classify-model-link-${model['bk_classification_id']}`]"
                                    :key="index"
                                    :title="model['bk_obj_name']"
                                    :to="getRouterLink(model)"
                                     @click.stop>
                                    {{model['bk_obj_name']}}
                                </router-link>
                            </li>
                            <li class="classify-model-item" v-if="classification['bk_objects'].length > 5">
                                <a href="javascript:void(0)" class="classify-model-link more" @click="showModels(classification)"></a>
                            </li>
                        </ul>
                    </li>
                </ul>
            </div>
        </transition>
        <transition name="fade">
            <div class="model-container" v-show="display.type === 'model'">
                <bk-button class="model-return" @click="display.type = 'classification'">{{$t("Index['返回上级']")}}</bk-button>
                <div class="model-list">
                    <router-link exact class="model-item" v-for="(model, index) in sortedDisplayModels"
                        :key="index"
                        :to="`/organization/${model['bk_obj_id']}`">
                        <i :class="['icon-cc-star', {navigated: isModelNavigated(model)}]" @click.stop.prevent="toggleModelNavVisible(model)"></i>
                        <div class="model-name-layout fl">
                            <span class="model-name">
                                <h3 class="model-name-text">{{model['bk_obj_name']}}</h3>
                                <span class="model-name-id">{{model['bk_obj_id']}}</span>
                            </span>
                        </div>
                        <div class="model-icon-layout">
                            <i :class="['model-icon', model['bk_obj_icon']]"></i>
                        </div>
                    </router-link>
                </div>
            </div>
        </transition>
    </div>
</template>
<script>
    import vSearch from './children/search'
    import { bk_host_manage as bkHostManage } from '@/common/json/static_navigation.json'
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                display: {
                    type: 'classification',
                    displayClassificationId: null
                },
                hostManageClassification: {
                    'bk_classification_icon': bkHostManage.icon,
                    'bk_classification_id': bkHostManage.id,
                    'bk_classification_name': this.$t(bkHostManage.i18n),
                    'bk_classification_type': 'inner',
                    'bk_objects': bkHostManage.children.map(nav => {
                        return {
                            'bk_obj_name': this.$t(nav.i18n),
                            'bk_obj_id': nav.id,
                            'path': nav.path,
                            'bk_classification_id': bkHostManage.id
                        }
                    })
                }
            }
        },
        computed: {
            ...mapGetters('usercustom', ['usercustom']),
            ...mapGetters('navigation', ['authorizedClassifications']),
            sortedClassifications () {
                let sortedClassifications = [this.hostManageClassification, ...this.authorizedClassifications]
                return sortedClassifications
                // 已添加到导航的排到前面
                /* const navigatedValue = {
                    true: 1,
                    false: 0
                }
                return sortedClassifications.sort((classificationA, classificationB) => {
                    return navigatedValue[this.isClassificationNavigated(classificationB)] - navigatedValue[this.isClassificationNavigated(classificationA)]
                }) */
            },
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
            },
            sortedDisplayModels () {
                let models = [...this.displayClassification['bk_objects']]
                return models
                // 已添加到导航的排到前面
                /* const navigatedValue = {
                    true: 1,
                    false: 0
                }
                return models.sort((modelA, modelB) => {
                    return navigatedValue[this.isModelNavigated(modelB)] - navigatedValue[this.isModelNavigated(modelA)]
                }) */
            }
        },
        methods: {
            isHostManage (obj) {
                return obj['bk_classification_id'] === bkHostManage.id
            },
            getRouterLink (model) {
                if (this.isHostManage(model)) {
                    return model.path
                }
                return `/organization/${model['bk_obj_id']}`
            },
            // 显示视图
            showModels (classification) {
                const classificationId = classification['bk_classification_id']
                if (this.isHostManage(classification)) return
                this.display.type = 'model'
                this.display.displayClassificationId = classificationId
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
        overflow: hidden;
        background: #fff url(../../common/images/index-bg.png) center bottom no-repeat;
    }
    .icon-cc-star{
        position: absolute;
        top: 7px;
        right: 13px;
        font-size: 18px;
        color: #5f95de;
        cursor: pointer;
        &.navigated{
            color: #f5ef90;
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
        width: 1372px;
        .classify-item{
            font-size: 14px;
            position: relative;
            display: inline-block;
            width: 323px;
            height: 230px;
            overflow: hidden;
            margin: 0 10px 20px;
            background-color: #fff;
            transition: box-shadow .2s ease-out;
            &:hover{
                box-shadow: 0px 5px 11px 0px rgba(0, 0, 0, 0.2);
                .classify-info-layout{
                    height: 155px;
                }
                .classify-models-list{
                    height: 75px;
                    &:after{
                        opacity: 0;
                    }
                }
                .classify-model-item:nth-child(n+4){
                    opacity: 1;
                }
            }
        }
    }
    .classify-info-layout {
        height: 165px;
        text-align: center;
        color: #fff;
        cursor: pointer;
        border-top-left-radius: 2px;
        border-top-right-radius: 2px;
        transition: height .2s ease-out;
        background-image: linear-gradient(#6ea9f9, #6ea9f9),linear-gradient(10deg, #42c0ff 0%, #3c96ff 100%), linear-gradient(#6ea9f9, #6ea9f9);
        &-bk_host_manage {
            cursor: default;
        }
        &:before{
            content: "";
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
        }
        .classify-info{
            display: inline-block;
            vertical-align: middle;
            margin-top: 6px;
            .classify-info-icon{
                display: block;
                font-size: 48px;
            }
            .classify-info-name{
                display: block;
                font-size: 20px;
                margin-top: 6px;
            }
        }
    }
    .classify-models-list{
        height: 65px;
        padding: 7px 0;
        position: relative;
        font-size: 0;
        border: solid 1px #c3cdd7;
        border-top: none;
        border-bottom-left-radius: 2px;
        border-bottom-right-radius: 2px;
        transition: height .2s ease-out;
        &.has-more:after{
            content: '';
            position: absolute;
            bottom: 18px;
            left: 50%;
            width: 4px;
            height: 4px;
            margin-left: -2px;
            border-radius: 50%;
            background-color: rgba(195, 205, 215, .4);
            box-shadow: -7px 0 rgba(195, 205, 215, .4), 7px 0 rgba(195, 205, 215, .4);
            transition: opacity .2s ease-out;
            pointer-events: none;
        }
        .classify-model-item {
            display: inline-block;
            vertical-align: middle;
            width: 107px;
            text-align: center;
            padding: 4px 4px;
            font-size: 14px;
            transition: all .2s ease-out;
            &:nth-child(n+4) {
                opacity: 0;
            }
            &:last-child{
                .classify-model-link{
                    max-width: 150px;
                }
            }
            .classify-model-link{
                display: inline-block;
                height: 20px;
                line-height: 20px;
                text-align: center;
                padding: 0 8px;
                color: #3c96ff;
                border-radius: 10px;
                max-width: 100%;
                @include ellipsis;
                &:hover{
                    background-color: #ebf0f7;
                }
                &.more{
                    border-radius: 0;
                    background-color: transparent;
                    position: relative;
                    overflow: visible;
                    &:after{
                        content: '';
                        position: absolute;
                        top: 50%;
                        left: 50%;
                        width: 4px;
                        height: 4px;
                        margin: -2px 0 0 -2px;
                        border-radius: 50%;
                        background-color: #c3cdd7;
                        box-shadow: -7px 0 #c3cdd7, 7px 0 #c3cdd7;
                        pointer-events: none;
                    }
                }
            }
        }
    }
    .model-container{
        position: relative;
        height: calc(100% - 114px);
        margin: 0 auto;
        width: 1660px;
        .model-return{
            width: 102px;
            margin: 0 20px 14px;
        }
    }
    .model-list{
        height: calc(100% - 50px);
        overflow-x: hidden;
        overflow-y: auto;
        font-size: 0;
        @include scrollbar;
        .model-item{
            position: relative;
            display: inline-block;
            width: 306px;
            height: 112px;
            margin: 0 0 20px 20px;
            font-size: 14px;
            background-color: #f5faff;
            border-radius: 2px;
            border: solid 1px #c6d4e3;
            transition: box-shadow .2s ease-out !important;
            &:hover{
                background-color: #ffffff;
                box-shadow: 0px 5px 11px 0px rgba(0, 0, 0, 0.2);
                border: solid 1px #499dff;
                .model-name-text,
                .model-name-id{
                    color: #3c96ff;
                }
                .model-icon-layout{
                    .model-icon{
                        color: #6fb1ff;
                    }
                }
            }
            .icon-cc-star{
                color: #dde4eb;
                top: 6px;
                right: 12px;
                &.navigated{
                    color: #f9cd6e;
                }
            }
            .model-icon-layout{
                height: 100%;
                overflow: hidden;
                text-align: center;
                font-size: 0;
                &:before{
                    content: "";
                    width: 0;
                    height: 100%;
                    display: inline-block;
                    vertical-align: middle;
                }
                .model-icon{
                    font-size: 80px;
                    display: inline-block;
                    vertical-align: middle;
                    color:#c1d9f5;
                }
            }
            .model-name-layout{
                height: 100%;
                width: 165px;
                padding: 0px 0px 0 30px;
                &:before{
                    content: "";
                    width: 0;
                    height: 100%;
                    display: inline-block;
                    vertical-align: middle;
                }
                .model-name{
                    display: inline-block;
                    vertical-align: middle;
                    width: 100%;
                }
            }
            .model-name-text{
                margin: 0;
                font-size: 18px;
                font-weight: normal;
                color: #333948;
                @include ellipsis;
            }
            .model-name-id{
                display: block;
                line-height: 24px;
                color: $textColor;
                @include ellipsis;
            }
        }
    }
</style>
<style lang="scss">
    /* 导航栏收起时的分类列表宽度 */
    /* screen <= 1279 放3个分类 */
    @media (max-width: 1452px){
        .content-wrapper.fold{
            .index-wrapper{
                .search-box{
                    width: 700px;
                }
                .classify-list{
                    width: 1029px;
                }
            }
        }
    }
    @media (min-width: 1453px){
        .content-wrapper.fold{
            .index-wrapper{
                .classify-list{
                    width: 1372px;
                }
            }
        }
    }
    /* 导航展开时的分类列表宽度 */
    /* screen <= 1419 放2个分类 */
    @media (max-width: 1612px) {
        .content-wrapper{
            .index-wrapper{
                .classify-list{
                    width: 1029px;
                }
            }
        }
    }

    /* 导航栏收起时的模型列表宽度 */
    @media (max-width: 1389px) {
        .content-wrapper.fold{
            .index-wrapper{
                .model-container {
                    width: 1000px;
                }
            }
        }
    }
    @media (min-width: 1390px) and (max-width: 1709px) {
        .content-wrapper.fold{
            .index-wrapper{
                .model-container{
                    width: 1330px;
                }
            }
        }
    }
    /* 导航栏展开时的模型列表宽度 */
    @media (max-width: 1544px) {
        .content-wrapper{
            .index-wrapper{
                .model-container{
                    width: 1000px;
                }
            }
        }
    }
    @media (min-width: 1545px) and (max-width: 1864px) {
        .content-wrapper{
            .index-wrapper{
                .model-container{
                    width: 1330px;
                }
            }
        }
    }
</style>