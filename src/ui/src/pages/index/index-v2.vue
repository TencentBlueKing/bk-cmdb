<template>
    <div class="index-wrapper">
        <v-search></v-search>
        <div class="classify-layout">
            <ul class="classify-list" ref="classifyList">
                <li class="classify-item"
                    v-for="(classify, classifyIndex) in sortedClassifications"
                    :index="classifyIndex"
                    @mouseleave="toggleClassify(classify, classifyIndex, 'close')">
                    <div ref="classifyItemLayout" class="item-layout">
                        <i class="classify-pin icon-cc-pin" v-if="!staticClassify.includes(classify['bk_classification_id'])"
                            :class="{highlight: classifyNavigation.includes(classify['bk_classification_id'])}"
                            :title="classifyNavigation.includes(classify['bk_classification_id']) ? $t('Index[\'取消导航\']') : $t('Index[\'固定到导航\']')"
                            @click.stop.prevent="pinClassify(classify)">
                        </i>
                        <div class="classify-info">
                            <i :class="['classify-icon', classify['bk_classification_icon']]"></i>
                            <p class="classify-name">{{classify['bk_classification_name']}}</p>
                        </div>
                        <transition-group tag="div" name="fade" class="classify-models" :duration="300" @afterLeave="removeToggleAttribute(classifyIndex)">
                            <ul :class="['model-list', $i18n.locale]"
                                v-for="col in Math.ceil(classify['bk_objects'].length / modelsPerCol)"
                                v-show="col === 1 || openedClassify === classify['bk_classification_id']"
                                :key="col">
                                <li class="model-item clearfix"
                                    v-for="modelColIndex in modelsPerCol"
                                    v-if="getModelIndex(col, modelColIndex) < classify['bk_objects'].length">
                                    <router-link exact class="model-name fl"
                                        :to="getModelLink(classify, col, modelColIndex)"
                                        :title="getModelObject(classify, col, modelColIndex)['bk_obj_name']">
                                        {{getModelObject(classify, col, modelColIndex)['bk_obj_name']}}
                                    </router-link>
                                    <i class="model-stick icon-cc-stick"
                                        v-if="!['bk_host_manage'].includes(classify['bk_classification_id']) && getModelIndex(col, modelColIndex)"
                                        @click.stop.prevent="stickModel(classify, col, modelColIndex)">
                                    </i>
                                </li>
                            </ul>
                        </transition-group>
                        <a href="javascript:void(0)" class="model-more"
                            v-if="openedClassify !== classify['bk_classification_id'] && classify['bk_objects'].length > 5"
                            :title="$t('Index[\'更多\']')"
                            @click.stop.prevent="toggleClassify(classify, classifyIndex, 'open')">{{$t("Index['更多']")}}</a>
                    </div>
                </li>
            </ul>
        </div>
        <p class="copyright">
            Copyright © 2012-{{year}} Tencent BlueKing. All Rights Reserved. 腾讯蓝鲸 版权所有
        </p>
    </div>
</template>
 <script>
    import vSearch from './children/search'
    import bus from '@/eventbus/bus'
    import throttle from 'lodash.throttle'
    import { bk_host_manage as bkHostManage } from '@/common/json/static_navigation.json'
    import { mapGetters } from 'vuex'
    import {addResizeListener, removeResizeListener} from '@/utils/resize-event.js'
    export default {
        components: {
            vSearch
        },
        data () {
            const year = (new Date()).getFullYear()
            const hostManageClassification = {
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
            return {
                year,
                hostManageClassification,
                modelsPerCol: 5,
                openedClassify: null,
                staticClassify: ['bk_host_manage', 'bk_organization'],
                notModelClassify: ['bk_host_manage', 'bk_back_config'],
                throttleLayout: null
            }
        },
        computed: {
            ...mapGetters('navigation', ['authorizedClassifications']),
            ...mapGetters('usercustom', ['usercustom', 'classifyNavigationKey', 'classifyModelSequenceKey']),
            classifyNavigation () {
                return this.usercustom[this.classifyNavigationKey] || []
            },
            classifyModelSequence () {
                return this.usercustom[this.classifyModelSequenceKey] || {}
            },
            sortedClassifications () {
                const classifications = this.$deepClone(this.authorizedClassifications)
                classifications.forEach((classify, index) => {
                    if (this.classifyModelSequence.hasOwnProperty(classify['bk_classification_id'])) {
                        classify['bk_objects'].sort((modelA, modelB) => {
                            return this.getModelSequence(classify, modelA) - this.getModelSequence(classify, modelB)
                        })
                    }
                })
                return [this.hostManageClassification, ...classifications]
            }
        },
        watch: {
            sortedClassifications () {
                this.throttleLayout()
            }
        },
        mounted () {
            this.initResizeListener()
        },
        beforeDestory () {
            if (this.throttleLayout) {
                removeResizeListener(this.$el, this.throttleLayout)
            }
        },
        methods: {
            initResizeListener () {
                this.throttleLayout = throttle(() => {
                    this.doLayout()
                }, 50, {leading: false})
                addResizeListener(this.$el, this.throttleLayout)
            },
            doLayout () {
                const classifyItemSpace = 230
                const wrapperWidth = this.$el.getBoundingClientRect().width
                const maxItemInRow = Math.floor(wrapperWidth / classifyItemSpace)
                const $classifyList = this.$refs.classifyList
                if (this.sortedClassifications.length > maxItemInRow) {
                    $classifyList.style.width = maxItemInRow * classifyItemSpace + 'px'
                } else {
                    $classifyList.style.width = 'auto'
                }
            },
            getModelSequence (classify, model) {
                if (this.classifyModelSequence.hasOwnProperty(classify['bk_classification_id'])) {
                    const sequence = this.classifyModelSequence[classify['bk_classification_id']]
                    const modelSequence = sequence.indexOf(model['bk_obj_id'])
                    return modelSequence === -1 ? classify['bk_objects'].length : modelSequence
                }
                return classify['bk_objects'].length
            },
            getModelIndex (col, colIndex) {
                return (col - 1) * this.modelsPerCol + colIndex - 1
            },
            getModelObject (classify, col, colIndex) {
                return classify['bk_objects'][this.getModelIndex(col, colIndex)]
            },
            getModelLink (classify, col, colIndex) {
                const model = this.getModelObject(classify, col, colIndex)
                if (this.notModelClassify.includes(classify['bk_classification_id'])) {
                    return model.path
                }
                return `/organization/${model['bk_obj_id']}`
            },
            stickModel (classify, col, colIndex) {
                const model = this.getModelObject(classify, col, colIndex)
                const modelId = model['bk_obj_id']
                const classifyId = classify['bk_classification_id']
                let newClassifyModelSequence = {...this.classifyModelSequence}
                let modelSequence = newClassifyModelSequence[classifyId] || []
                const index = modelSequence.indexOf(modelId)
                if (index !== -1) {
                    modelSequence.splice(index, 1)
                }
                modelSequence.unshift(modelId)
                newClassifyModelSequence[classifyId] = modelSequence
                const upateParams = {}
                upateParams[this.classifyModelSequenceKey] = newClassifyModelSequence
                this.$store.dispatch('usercustom/updateUserCustom', upateParams)
            },
            async pinClassify (classify) {
                let isPin = false
                const classifyId = classify['bk_classification_id']
                const pinParams = {}
                let classifyNavigation = [...this.classifyNavigation]
                if (classifyNavigation.includes(classifyId)) {
                    classifyNavigation.splice(classifyNavigation.indexOf(classifyId), 1)
                } else {
                    classifyNavigation.push(classifyId)
                    isPin = true
                }
                pinParams[this.classifyNavigationKey] = classifyNavigation
                const pinResult = await this.$store.dispatch('usercustom/updateUserCustom', pinParams)
                if (isPin && pinResult.result) {
                    bus.$emit('handlePinClassify', classify)
                }
                if (pinResult.result) {
                    this.$alertMsg(isPin ? this.$t('Index["添加导航成功"]') : this.$t('Index["取消导航成功"]'), 'success')
                }
            },
            toggleClassify (classify, index, type = 'open') {
                const classifyItemLayout = this.$refs.classifyItemLayout[index]
                if (type === 'open') {
                    const wrapperRect = this.$el.getBoundingClientRect()
                    const classifyItemLayoutRect = classifyItemLayout.getBoundingClientRect()
                    const columns = Math.ceil(classify['bk_objects'].length / this.modelsPerCol)
                    const openWidth = columns * (this.$i18n.locale === 'en' ? 150 : 120) + 60
                    const openDirection = classifyItemLayoutRect.left + openWidth < wrapperRect.width ? 'item-layout-open-right' : 'item-layout-open-left'
                    classifyItemLayout.classList.add('item-layout-open')
                    classifyItemLayout.classList.add(openDirection)
                    classifyItemLayout.setAttribute('data-open-direction', openDirection)
                    classifyItemLayout.style.width = openWidth + 'px'
                    this.openedClassify = classify['bk_classification_id']
                } else {
                    classifyItemLayout.style.width = '100%'
                    this.openedClassify = null
                }
            },
            removeToggleAttribute (index) {
                const classifyItemLayout = this.$refs.classifyItemLayout[index]
                classifyItemLayout.classList.remove('item-layout-open')
                classifyItemLayout.classList.remove(classifyItemLayout.getAttribute('data-open-direction'))
            }
        }
    }
 </script>

 <style lang="scss" scoped>
    $itemBorder: #cbdef6;
    $modelColor: #3c96ff;
    .index-wrapper{
        height: 100%;
        padding: 60px 0 0 0;
        position: relative;
    }
    .classify-layout{
        margin: 15px 0 0;
        height: calc(100% - 120px);
        overflow-y: auto;
        overflow-x: hidden;
        @include scrollbar;
    }
    .classify-list{
        font-size: 0;
        height: 100%;
        max-width: 1610px;
        margin: 0 auto;
        .classify-item{
            display: inline-block;
            font-size: 14px;
            width: 180px;
            height: 286px;
            margin: 40px 25px 0;
            position: relative;
        }
    }
    .item-layout{
        position: absolute;
        top: 0;
        width: 100%;
        height: 100%;
        border-radius: 2px;
        border: solid 1px $itemBorder;
        background-color: #fff;
        transition: box-shadow,width .3s ease-in-out;
        &-open{
            z-index: 1;
            &-right{
                left: 0;
            }
            &-left{
                right: 0;
            }
        }
        &:hover{
            box-shadow: 0 3px 8px 0px rgba(37, 81, 140, .15);
        }
        .classify-pin{
            position: absolute;
            top: 10px;
            right: 8px;
            color: rgba(110, 169, 249, 0.4);
            cursor: pointer;
            &.highlight,
            &:hover{
                color: #3c80dc;
            }
        }
        .classify-info{
            height: 110px;
            padding: 18px 20px 0;
            border-bottom: 1px dashed #dfecfc;
            .classify-icon{
                display: block;
                width: 50px;
                height: 50px;
                line-height: 50px;
                text-align: center;
                color: #fff;
                font-size: 32px;
                border-radius: 50%;
                background-image: linear-gradient(180deg, rgba(73, 136, 223, 0.8) 0%, rgba(32, 96, 216, 0.8) 100%), 
                                  linear-gradient(rgba(110, 169, 249, 0.8), rgba(110, 169, 249, 0.8));
                background-blend-mode: normal, normal;
                opacity: 0.8;
            }
            .classify-name{
                margin: 6px 0 0 0;
                line-height: 24px;
                color: #479cf9;
                font-size: 18px;
                @include ellipsis;
            }
        }
        .model-more{
            position: absolute;
            bottom: 11px;
            right: 8px;
            font-size: 12px;
            color: rgba(71, 156, 249, 0.5);
            &:hover{
                color: rgb(71, 156, 249);
            }
        }
    }
    .classify-models{
        position: relative;
        height: 176px;
        white-space: nowrap;
        font-size: 0;
        padding: 0 0 0 10px;
    }
    .model-list{
        display: inline-block;
        vertical-align: top;
        height: 170px;
        padding: 10px 0;
        font-size: 14px;
        .model-item{
            width: 120px;
            height: 30px;
            .model-name{
                border-radius: 10px;
                padding: 0 10px;
                max-width: 103px;
                height: 22px;
                line-height: 22px;
                cursor: pointer;
                color: $modelColor;
                transition: none !important;
                @include ellipsis;
            }
            .model-stick{
                margin-left: 3px;
                cursor: pointer;
                color: rgba(110, 169, 249, 0.4);
                display: none;
                &:hover{
                   color: $modelColor;
                }
            }
            &:hover{
                .model-name{
                    background-color: #f1f7ff;
                }
                .model-stick{
                    display: inline-block;
                }
            }
        }
    }
    .model-list.en{
        .model-item{
            width: 150px;
            .model-name{
                max-width: 133px;
            }
        }
    }
    .copyright{
        text-align: center;
        position: absolute;
        width: 100%;
        bottom: 0;
        left: auto;
        padding: 0 0 10px 0;
        margin: 0;
        font-size: 14px;
        color: rgba(116, 120, 131, 0.5);
    }
 </style>