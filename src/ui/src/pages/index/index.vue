<template>
    <div class="index-wrapper">
        <v-search></v-search>
        <transition name="fade">
            <div class="classify-container" v-show="display.type === 'classify'">
                <ul class="classify-list">
                    <li class="classify-item clearfix" v-for="(classify, index) in classifies" :key="index" @click="showModels(classify)">
                        <div :class="['classify-navigator', {navigated: classify.navigated}]" @click.stop="toggleClassifyNavVisible(classify)">
                            <span class="icon-cc-thumbtack"></span>
                            <span>{{classify.navigated ? '取消导航' : '添加导航'}}</span>
                        </div>
                        <div class="classify-icon fl">
                            <i :class="classify['bk_classification_icon']"></i>
                        </div>
                        <div class="classify-name">
                            <h2  class="classify-name-text">{{classify['bk_classification_name']}}</h2>
                            <span class="classify-model-count">{{`模型：${classify['bk_objects'].length}个`}}</span>
                        </div>
                        <div class="classify-models">
                            <router-link exact class="classify-model-link"
                                v-for="(model, index) in classify['bk_objects']"
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
                <bk-button class="model-return" type="primary" @click="display.type = 'classify'">返回上级</bk-button>
                <div class="model-list">
                    <router-link exact class="model-item" v-for="(model, index) in displayClassify['bk_objects']"
                        :key="index"
                        :to="`organization/${model['bk_obj_id']}`">
                        <div :class="['classify-navigator', {navigated: model.navigated}]" @click.prevent.stop="toggleModelNavVisible(model)">
                            <span class="icon-cc-thumbtack"></span>
                            <span>{{model.navigated ? '取消导航' : '添加导航'}}</span>
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
                    type: 'classify',
                    displayClassifyId: '',
                    invisibleClassify: ['bk_host_manage', 'bk_biz_topo']
                }
            }
        },
        computed: {
            ...mapGetters(['allClassify', 'authority', 'adminAuthority']),
            ...mapGetters('usercustom', ['usercustom']),
            // 筛选出有权限的模型
            modelAuthority () {
                let authority = window.isAdmin === '1' ? this.adminAuthority : this.authority
                let modelConfig = authority['model_config']
                let classifyAuthority = Object.keys(modelConfig).filter(classifyId => {
                    return !this.display.invisibleClassify.includes(classifyId) && Object.keys(modelConfig[classifyId]).length
                })
                let modelAuthority = {}
                classifyAuthority.forEach(classifyId => {
                    modelAuthority[classifyId] = Object.keys(modelConfig[classifyId])
                })
                return modelAuthority
            },
            // 用户自定义导航内容
            navigation () {
                return this.usercustom.navigation || {}
            },
            // 根据用户模型权限，设置分类数据
            classifies () {
                let modelAuthority = this.modelAuthority
                // 1.删除没有权限的分类
                let allClassifyClone = JSON.parse(JSON.stringify(this.allClassify))
                allClassifyClone = allClassifyClone.filter(({bk_classification_id: bkClassficationId}) => {
                    return modelAuthority.hasOwnProperty(bkClassficationId)
                })
                // 2.删出有权限的分类下没有权限的模型
                allClassifyClone.forEach(({bk_classification_id: bkClassficationId, bk_objects: bkObjects}, index, originalClassify) => {
                    originalClassify[index]['bk_objects'] = bkObjects.filter(({bk_obj_id: bkObjId}) => {
                        return modelAuthority[bkClassficationId].includes(bkObjId)
                    })
                })
                // 3.检查是否被加入导航
                let navigation = this.navigation
                allClassifyClone.forEach(classify => {
                    let {bk_classification_id: bkClassficationId, bk_objects: bkObjects} = classify
                    bkObjects.forEach(bkObject => {
                        if (!navigation.hasOwnProperty(bkClassficationId)) {
                            bkObject.navigated = false
                        } else {
                            bkObject.navigated = navigation[bkClassficationId].some(navigationObjId => navigationObjId === bkObject['bk_obj_id'])
                        }
                    })
                    classify.navigated = !bkObjects.some(({navigated}) => !navigated)
                })
                return allClassifyClone.filter(({bk_objects: bkObjects}) => !!bkObjects.length)
            },
            // 当前模型视图所对应的分类
            displayClassify () {
                if (this.display.displayClassifyId && this.display.type === 'model') {
                    return this.classifies.find(({bk_classification_id: bkClassficationId}) => bkClassficationId === this.display.displayClassifyId)
                }
                return {'bk_objects': []}
            }
        },
        created () {
            // 获取用户自定义内容
            this.$store.dispatch('usercustom/getUserCustom')
        },
        methods: {
            // 显示视图
            showModels (classify) {
                this.display.type = 'model'
                this.display.displayClassifyId = classify['bk_classification_id']
            },
            // 分类视图添加/取消导航切换
            toggleClassifyNavVisible (classify) {
                let navigation = JSON.parse(JSON.stringify(this.navigation))
                let {
                    'bk_classification_id': bkClassficationId,
                    'bk_objects': bkObjects,
                    navigated
                } = classify
                navigation[bkClassficationId] = navigated ? [] : bkObjects.map(({bk_obj_id: bkObjId}) => bkObjId)
                this.$store.dispatch('usercustom/updateUserCustom', {navigation}, !navigated)
            },
            // 模型视图添加/取消导航切换
            toggleModelNavVisible (model) {
                let navigation = JSON.parse(JSON.stringify(this.navigation))
                let {
                    'bk_classification_id': bkClassficationId,
                    'bk_objects': bkObjects,
                    navigated
                } = this.displayClassify
                navigation[bkClassficationId] = navigation[bkClassficationId] || []
                if (model.navigated) {
                    navigation[bkClassficationId] = navigation[bkClassficationId].filter(bkObjId => bkObjId !== model['bk_obj_id'])
                } else {
                    navigation[bkClassficationId].push(model['bk_obj_id'])
                }
                this.$store.dispatch('usercustom/updateUserCustom', {navigation}, !model.navigated)
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