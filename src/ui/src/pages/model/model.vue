/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <div class="model-wrapper">
        <div class="left-tap-contain">
            <div class="list-tap" v-bkloading="{isLoading: isClassifyLoading}">
                <ul>
                    <li :class="{'active': topoView === 'models' && curClassify['bk_classification_id']===item['bk_classification_id']}" v-for="(item, index) in classifyList" @click="changeClassify(item)" @click.stop="topoView = 'models'">
                        <i :class="item['bk_classification_icon']"></i>
                        <span class="text">{{item['bk_classification_name']}}</span>
                    </li>
                </ul>
                <div class="bottom-btn-contain" @click="addClassify('add')">
                    <a class="bottom-btn" >
                        <span>{{$t('Common["新增"]')}}</span>
                    </a>
                </div>
            </div>
            <div :class="['global-models', {'active': topoView === 'GLOBAL_MODELS'}]">
                <button class="btn-global" @click="topoView = 'GLOBAL_MODELS'">
                    <i class="icon-cc-fullscreen"></i>
                    <span class="text">{{$t("ModelManagement['全局视图']")}}</span>
                </button>
            </div>
        </div>
        <div class="right-contain clearfix">
            <v-global-models v-if="topoView === 'GLOBAL_MODELS'"></v-global-models>
            <div class="model-box clearfix" v-bkloading="{isLoading: isTopoLoading}" v-else>
                <div class="model-diagram" v-if="curClassify['bk_classification_id']">
                    <div class="model-topo-box" v-show="topoList.length != 0 && !isCreateShow">
                        <template v-if="curClassify['bk_classification_id'] !== 'bk_biz_topo'">
                            <v-topo
                                :addModelAvailable="addModelAvailable"
                                :topo="topo"
                                :isTopoLoading.sync="isTopoLoading"
                                :modelNodes="topoList"
                                :curClassify="curClassify['bk_classification_id']"
                                :isChangeClassify="isChangeClassify"
                                :isPaused="curModelDiagram === 'disable'"
                                :disableModelList="disableModelList"
                                :isInnerType="curTempClassify['bk_classification_type'] === 'inner'"
                                @updateIsChangeClassify="updateIsChangeClassify"
                                @nodeClick="nodeClick"
                                @createModel="createModel"
                                @editModel="popShow('edit')"
                                @changeModelDiagram="changeModelDiagram"
                                @deleteClass="deleteClassify"
                            ></v-topo>
                        </template>
                        <template v-else>
                            <div class="model-content" v-if="topoList.length != 0">
                                <bk-button v-if="curTempClassify['bk_classification_id'] !== 'bk_biz_topo' && curTempClassify['bk_classification_id'] !== 'bk_host_manage'" type="primary" class="topo-btn edit" @click="popShow('edit')" :title="$t('Common[\'编辑\']')">
                                    <i class="icon icon-cc-edit"></i>
                                </bk-button>
                                <bk-button type="danger" class="topo-btn del" v-if="curTempClassify['bk_classification_type']!=='inner'" @click="deleteClassify">
                                    <i class="icon icon-cc-del"></i>
                                </bk-button>
                                <ul class="topo-wrapper clearfix" :class="'topo-wrapper-'+curClassify['bk_classification_id']">
                                    <li :class="{'line':curClassify['bk_classification_id'] ==='bk_biz_topo','locks-icon-content':item['ispre'], 'default':item['bk_obj_id'] === 'biz'}" v-for="item in topoList" @click="editModel(item, false)" v-if="!item['bk_ispaused']">
                                        <div class="content">
                                            <i v-if="curClassify['bk_classification_id']==='bk_biz_topo'&&item['bk_obj_id']!=='biz'&&item['bk_obj_id']!=='module'&&item['bk_obj_id']!=='host'" class="icon-add prev icon-cc-round-plus" @click.stop="addModel(item, 'prev')"></i>
                                            <div>
                                                <i class="topo-icon" :class="item['bk_obj_icon']"></i>
                                            </div>
                                            <div class="content-name">{{item['bk_obj_name']}}</div>
                                            <i class="icon-add next icon-cc-round-plus" @click.stop="addModel(item, 'next')" v-if="item['bk_obj_id']!=='set'&&item['bk_obj_id']!=='modules'&&curClassify['bk_classification_id']==='bk_biz_topo'&&item['bk_next_obj']===''&&item['bk_obj_id']!=='host'"></i>
                                        </div>
                                    </li>
                                </ul>
                            </div>
                        </template>
                    </div>
                    <div class="no-model-prompting tc" v-show="topoList.length == 0 || isCreateShow">
                        <bk-button v-if="curTempClassify['bk_classification_id'] !== 'bk_biz_topo' && curTempClassify['bk_classification_id'] !== 'bk_host_manage'" type="primary" class="topo-btn edit" @click="popShow('edit')" :title="$t('Common[\'编辑\']')">
                            <i class="icon icon-cc-edit"></i>
                        </bk-button>
                        <bk-button type="danger" class="topo-btn del" v-if="curTempClassify['bk_classification_type']!=='inner'" @click="deleteClassify">
                            <i class="icon icon-cc-del"></i>
                        </bk-button>
                        <button class="bk-button vis-button vis-enable" v-if="addModelAvailable" @click="isShowDisableList = true">
                            <i class="bk-icon icon-minus-circle-shape"></i>
                            <span class="vis-button-text">{{disableModelList.length}}</span>
                        </button>
                        <transition name="topo-disable-list">
                            <div class="topo-disable tl" v-show="isShowDisableList">
                                <label class="disable-title">
                                    <span>{{$t('ModelManagement["已停用模型"]')}}</span>
                                    <i class="bk-icon icon-arrows-right" @click="isShowDisableList = false"></i>
                                </label>
                                <ul class="disable-list" ref="disableList">
                                    <li class="disable-item" v-for="(model, index) in disableModelList" :key="index">
                                        <a class="disable-item-link" href="javascript:void(0)" @click="nodeClick(model, true)">{{model['bk_obj_name']}}</a>
                                    </li>
                                </ul>
                            </div>
                        </transition>
                        <img src="../../common/images/no_model_prompting.png" alt="">
                        <p v-if="isCreateShow">{{$t('ModelManagement["此分组下无已启用模型"]')}}</p>
                        <p v-else>{{$t('ModelManagement["此分组下无模型"]')}}</p>
                        <bk-button type="primary" class="create-btn" @click="showAddModel">{{$t('Common["立即创建"]')}}</bk-button>
                    </div>
                </div>
            </div>
        </div>
        <v-pop
            :isShow.sync="isPopShow"
            :type="category"
            :classification="curTempClassify"
            @confirm="saveClassify"
        ></v-pop>
        <v-sideslider :isShow.sync="slider.isBusinessShow"
        :title="sliderTitle"
        :hasCloseConfirm="true"
        :isCloseConfirmShow="slider.isCloseConfirmShow"
        @closeSlider="closeSliderConfirm">
            <div class="content slide-content clearfix" slot="content">
                <bk-tab :active-name="curTabName" @tab-changed="tabChanged">
                    <bk-tabpanel name="host" :title="$t('ModelManagement[\'模型配置\']')">
                        <v-field ref="field"
                        :bk_classification_id="curClassify['bk_classification_id']"
                        :type="curModel.type"
                        :id="curModel['id']"
                        :isShow="slider.isBusinessShow"
                        :objId="curModel['bk_obj_id']"
                        :isMainLine="curClassify['bk_classification_id']==='bk_biz_topo'"
                        :classificationId="curClassify['bk_classification_id']"
                        :associationId="curInsertInfo.preObj"
                        :isReadOnly="isModelDetailReadOnly"
                        :isCreateField="isCreateField"
                        :isSliderShow.sync="slider.isBusinessShow"
                        @getTopogical="getTopogical"
                        @cancel="cancel"
                        @baseInfoSuccess="baseInfoSuccess"
                        @newField="isNewField=!isNewField"
                        ></v-field>
                    </bk-tabpanel>
                    <bk-tabpanel name="layout" :title="$t('ModelManagement[\'字段分组\']')" :show="curModel.type==='change'">
                        <v-layout ref="layout"
                        :isShow="curTabName==='layout'"
                        :activeModel="curModel"
                        @cancel="cancel"
                        ></v-layout>
                    </bk-tabpanel>
                    <bk-tabpanel name="other" :title="$t('ModelManagement[\'其他操作\']')" :show="curModel.type==='change'">
                        <v-other
                            :parentClassificationId = "curClassify['bk_classification_id']"
                            :item="curModel"
                            :id="curModel['id']"
                            :isMainLine="curClassify['bk_classification_id']==='bk_biz_topo'"
                            @getTopogical="getTopogical"
                            :isReadOnly="isModelDetailReadOnly"
                            @deleteModel="cancel"
                            @closeSideSlider="cancel"
                            @stopModel="closeBusiness">
                        </v-other>
                    </bk-tabpanel>
                </bk-tab>
            </div>
        </v-sideslider>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    import vPop from './children/pop'
    import vSideslider from '@/components/slider/sideslider'
    import vBaseInfo from './children/baseInfo'
    import vField from './children/field'
    import vLayout from './children/layout'
    import vOther from './children/other'
    import vGlobalModels from './children/global-models'
    import vTopo from '@/components/topo/topo'
    import {mapGetters, mapActions} from 'vuex'
    const iconList = require('@/common/json/classIcon.json')
    export default {
        data () {
            return {
                topoView: 'models',
                isShowDisableList: false,
                isTopoLoading: false,           // 拓扑loading
                isClassifyLoading: false,       // 分组列表loading
                isNewField: false,
                topo: {                           // 普通拓扑
                    nodes: [],                      // 节点
                    edges: []                       // 关联关系
                },
                isModelDetailReadOnly: false,      // 模型详情是否只读
                category: '', // 分类
                isCreateField: true,       // 控制字段配置新增字段按钮的显示
                sliderTitle: {
                    text: '',
                    icon: 'icon-cc-model'
                },
                curTabName: 'host',
                slider: {
                    isBusinessShow: false,
                    isCloseConfirmShow: false
                },
                classifyList: [],           // 分类列表
                isEditClassify: false,
                curClassify: {},            // 当前类型
                curTempClassify: {},        // 当前类型修改内容 未保存
                curTempClassifyHoc: {},     // 新增分类临时定义变量
                topoList: [],               // 拓扑图 暂时为列表
                curModel: {
                    type: '',               // 新增还是修改 新增: new  修改: change
                    id: 0
                },
                isChoose: true,             // 判断编辑分类的时候是否选择了icon
                iconValue: 'icon-cc-business',               // 选择icon的值
                list: [],                      // icon 的值
                isChangeClassify: false,    // 是否点击的切换分组类型  true: 点击的切换分组 false: 点击的已启用/未启用
                curTopoStructure: [],       // 当前分类模型拓扑结构
                insertType: '',             // 插入模型  prev 向上 next 向下 mid 中间
                insertParams: {},           // 插入时相关参数
                curInsertInfo: {},          // 当前插入节点相关信息
                isPopShow: false,            // 新增编辑分类弹窗
                isInnerShow: false,           // 内置模型提示文本
                isIconDrop: false,            // 选择图标下拉框
                nowIndex: 0,                   // 选择图标下拉框当前index
                curModelDiagram: 'enable',      // 当前模块显示的图表类型,'enable'为已启用模型,'disable'为未启用模型
                defaultGroupId: -1             // 默认字段的分组id
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            addModelAvailable () {
                let notAvailable = ['bk_biz_topo', 'bk_host_manage', 'bk_organization']
                return notAvailable.indexOf(this.curClassify['bk_classification_id']) === -1
            },
            isAutoRouter () {
                return this.$route.query.hasOwnProperty('bk_classification_id')
            },
            disableModelList () {
                let disableModelList = []
                this.topoList.map(model => {
                    if (model['bk_ispaused']) {
                        disableModelList.push(model)
                    }
                })
                return disableModelList
            },
            isCreateShow () {
                return this.disableModelList.length === this.topoList.length && this.topoList.length
            }
        },
        watch: {
            'slider.isBusinessShow' (val) {
                if (!val) {
                    this.changeClassify()
                }
            }
        },
        methods: {
            ...mapActions(['getAllClassify']),
            closeSliderConfirm () {
                this.slider.isCloseConfirmShow = this.$refs.field.isCloseConfirmShow()
            },
            updateIsChangeClassify (val) {
                this.isChangeClassify = val
            },
            /*
                处理拓扑关系
                主要用户插入模型后把关联信息加上
            */
            setTopoRelation () {
                this.changeClassify(this.curClassify)
            },
            /*
                添加模型主关联
                item: 当前模块
                type: prev向上添加 next向下添加
            */
            addModel (item, type) {
                this.curTopoStructure.map(({bk_obj_id: objId, bk_pre_obj_id: preObj}) => {
                    if (objId === item['bk_obj_id']) {
                        this.curInsertInfo.preObj = preObj
                    }
                })
                // 调用新增模型
                this.showAddModel()
            },
            /*
                点击编辑弹窗
             */
            popShow (type) {
                if (type === 'edit') {
                    this.isPopShow = true
                    this.category = 'edit'
                }
                this.curTempClassify = JSON.parse(JSON.stringify(this.curClassify))
            },
            tabChanged (name) {
                // 防止重复点击
                if (this.curTabName === name) {
                    return
                }
                this.curTabName = name
                if (this.curTabName === 'base-info') {
                    this.$refs.baseInfo.getBaseInfo(this.curModel['bk_obj_id'])
                } else if (this.curTabName === 'host') {
                    this.$refs.field.init()
                }
            },
            closeBusiness () {
                this.slider.isBusinessShow = false
            },
            /*
                新增模型
            */
            createModel () {
                // 插入类型置空
                this.insertType = ''
                this.isModelDetailReadOnly = false
                this.showAddModel()
            },
            showAddModel () {
                this.sliderTitle.text = this.$t('ModelManagement["新增模型"]')
                this.sliderTitle.icon = 'icon-cc-model'
                this.curTabName = 'host'
                this.curModel = {}
                this.curModel.type = 'new'
                this.slider.isBusinessShow = true
            },
            /*
                编辑模型
                isModelDetailReadOnly: 模型详情是否可编辑
            */
            editModel (item, isModelDetailReadOnly) {
                if (item['bk_obj_id'] === 'biz' && this.curClassify['bk_classification_id'] === 'bk_biz_topo') {
                    return
                }
                this.isModelDetailReadOnly = isModelDetailReadOnly
                this.sliderTitle.text = item['bk_obj_name']
                this.sliderTitle.icon = 'icon-cc-host'
                this.curTabName = 'host'
                this.curModel = item
                this.curModel.type = 'change'
                let self = this
                setTimeout(function () {
                    self.$refs.field.init()
                }, 300)
                this.slider.isBusinessShow = true
                this.isEditClassify = false
            },
            /*
                重新启用模型确认弹框
            */
            restartModelConfirm (item) {
                this.isModelDetailReadOnly = true
                this.editModel(item)
            },
            /*
                重新启用模型
            */
            restartModel (item) {
                let params = {
                    bk_ispaused: false
                }
                this.$axios.put(`object/${item['id']}`, params).then(res => {
                    if (res.result) {
                        this.getTopogical()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                新增分类
            */
            addClassify (type) {
                if (type === 'add') {
                    this.isPopShow = true
                    this.category = 'add'
                }
                this.curTempClassify = {
                    bk_classification_icon: '',
                    bk_classification_id: '',
                    bk_classification_name: '',
                    bk_classification_type: ''
                }
            },
            /*
                保存分类
            */
            saveClassify (classification) {
                if (this.category === 'edit') {   // 更新分类信息
                    let params = {
                        bk_classification_name: classification['bk_classification_name'],
                        bk_classification_icon: classification['bk_classification_icon']
                    }
                    this.$axios.put(`object/classification/${this.curClassify['id']}`, params, {id: 'saveClassify'}).then(res => {
                        if (res.result) { // 保存时显示当前项
                            for (var i = 0; i < this.classifyList.length; i++) {
                                if (this.classifyList[i]['bk_classification_id'] === this.curClassify['bk_classification_id']) {
                                    this.getClassifyList(i)
                                }
                            }
                            this.isEditClassify = false
                            this.isPopShow = false
                            this.isChoose = true
                            this.curClassify['bk_classification_name'] = classification['bk_classification_name']
                            this.curClassify['bk_classification_icon'] = classification['bk_classification_icon']
                            this.curTempClassify = this.$deepClone(this.curClassify)
                            this.$store.commit('navigation/updateClassification', {
                                bk_classification_id: classification['bk_classification_id'],
                                bk_classification_name: classification['bk_classification_name'],
                                bk_classification_icon: classification['bk_classification_icon']
                            })
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else { // 新增分类
                    let createParams = {
                        bk_classification_id: classification['bk_classification_id'],
                        bk_classification_name: classification['bk_classification_name'],
                        bk_classification_icon: classification['bk_classification_icon']
                    }
                    this.$axios.post('object/classification', createParams, {id: 'saveClassify'}).then((res) => {
                        if (res.result) {
                            this.getClassifyList(this.classifyList.length)
                            this.isEditClassify = false
                            this.isPopShow = false
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            /*
                按照模型分类查询模型拓扑
                clsId: 模型分类ID
            */
            getTopoModelByClassify (clsId) {
                let p = new Promise((resolve, reject) => {
                    let url = ''
                    let method = ''
                    let data = {}
                    if (clsId === 'bk_biz_topo') {
                        url = `topo/model/${this.bkSupplierAccount}`
                        method = 'get'
                    } else {
                        url = 'objects/topo'
                        method = 'post'
                        data = {
                            bk_classification_id: clsId
                        }
                    }
                    this.$axios({
                        method: method,
                        url: url,
                        data: data
                    }).then(res => {
                        if (res.result) {
                            resolve(res.data)
                        } else {
                            reject(res['bk_error_msg'])
                        }
                    })
                })
                return p
            },
            /*
                获取业务 仅在业务拓扑分类下调用
            */
            getApp () {
                let p = new Promise((resolve, reject) => {
                    this.$axios.post('objects', {
                        bk_obj_id: 'biz'
                    }).then(res => {
                        if (res.result) {
                            resolve(res.data)
                        } else {
                            reject(res['bk_error_msg'])
                        }
                    })
                })
                return p
            },
            /*
                改变类型
            */
            changeClassify (item) {
                if (item === undefined) {
                    item = this.curClassify
                }
                this.isChangeClassify = true
                // this.isEditClassify = false
                this.curModelDiagram = 'enable'
                this.curClassify = item
                if (!this.curClassify['bk_classification_id']) {
                    this.classifyList.pop()
                }
                this.curTempClassify = {
                    id: item['id'],
                    bk_classification_name: item['bk_classification_name'],
                    bk_classification_type: item['bk_classification_type'],
                    bk_classification_id: item['bk_classification_id'],
                    bk_classification_icon: item['bk_classification_icon']
                }
                this.isTopoLoading = true
                // 切换时查询当前分类下的模型拓扑 及 当前分类下所有模型
                if (item['bk_classification_id'] === 'bk_biz_topo') {
                    Promise.all([
                        this.getTopoModelByClassify(item['bk_classification_id']),
                        this.getTopogical2(),
                        this.getApp()
                    ]).then(res => {
                        // 保存当前拓扑结构
                        this.curTopoStructure = res[0]
                        let topoList = res[1][0]['bk_objects']
                        topoList.push(res[2][0])
                        this.topoList = topoList
                        // 添加下一子节点相关属性
                        this.setModelAttr()
                        this.isTopoLoading = false
                    })
                } else {
                    Promise.all([
                        this.getTopoModelByClassify(item['bk_classification_id']),
                        this.getTopogical2()
                    ]).then(res => {
                        this.curTopoStructure = res[0]
                        // 当前分类下的所有模型
                        let curClsModel = res[1][0]['bk_objects']

                        // 与当前模型有关联的其他分类下的模型
                        let other = []
                        let asstObjects = res[1][0]['bk_asst_objects']
                        for (let key in asstObjects) {
                            let asstObject = asstObjects[key]
                            asstObject.map(obj => {
                                let isExist = false
                                for (let i = 0; i < curClsModel.length; i++) {
                                    if (obj['bk_obj_id'] === curClsModel[i]['bk_obj_id']) {
                                        isExist = true
                                        break
                                    }
                                }
                                if (!isExist) {
                                    let status = false
                                    // 防止插入重复的
                                    for (let i = 0; i < other.length; i++) {
                                        if (other[i]['bk_obj_id'] === obj['bk_obj_id']) {
                                            status = true
                                        }
                                    }
                                    if (!status) {
                                        other.push(obj)
                                    }
                                }
                            })
                        }
                        let modelIsPausedMap = {}
                        for (let i = 0; i < curClsModel.length; i++) {
                            modelIsPausedMap[curClsModel[i]['bk_obj_id']] = curClsModel[i]['bk_ispaused']
                        }
                        // 去掉关联对象被停用的
                        for (let i = other.length - 1; i >= 0; i--) {
                            let otherModel = other[i]
                            let isShow = false
                            for (let j = 0; j < this.curTopoStructure.length; j++) {
                                let structure = this.curTopoStructure[j]
                                let from = structure['from']
                                let to = structure['to']
                                if (otherModel['bk_obj_id'] === from['bk_obj_id']) {
                                    // 有处于启用状态的模型
                                    if (!modelIsPausedMap[to['bk_obj_id']]) {
                                        isShow = true
                                        break
                                    }
                                }
                                if (otherModel['bk_obj_id'] === to['bk_obj_id']) {
                                    // 有处于启用状态的模型
                                    if (!modelIsPausedMap[from['bk_obj_id']]) {
                                        isShow = true
                                        break
                                    }
                                }
                            }
                            if (!isShow) {
                                other.splice(i, 1)
                            }
                        }
                        this.topoList = res[1][0].bk_objects.concat(other)
                        let topoList = res[1][0].bk_objects.concat(other)
                        // topoList.concat(other)
                        this.topo = {
                            edges: res[0],
                            nodes: topoList
                        }
                    })
                }
            },
            /*
                添加下一子节点相关属性
            */
            setModelAttr () {
                // 当前分类为业务拓扑时调整顺序
                if (this.curClassify.bk_classification_id === 'bk_biz_topo') {
                    let tempList = []
                    this.curTopoStructure.map(val => {
                        this.topoList.map(li => {
                            if (val['bk_obj_id'] === li['bk_obj_id']) {
                                tempList.push(li)
                            }
                        })
                    })
                    this.topoList = tempList
                    this.topoList.splice()
                }
            },
            /*
                删除模型分类
            */
            deleteClassify () {
                var self = this
                this.$bkInfo({
                    title: this.$t('ModelManagement["确认要删除此分组？"]'),
                    confirmFn () {
                        self.deletes()
                    }
                })
            },
            deletes () {
                this.$axios.delete(`object/classification/${this.curClassify['id']}`).then(res => {
                    if (res.result) {
                        this.getClassifyList()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                保存基本信息成功
            */
            baseInfoSuccess (obj) {
                if (this.curModel.type === 'new') {
                    this.$store.dispatch('navigation/getClassifications', true)
                    this.curModel['id'] = obj['id']
                    this.sliderTitle.text = `${obj['bk_obj_name']}`
                    this.curModel['bk_obj_id'] = obj['bk_obj_id']
                    this.curModel.type = 'change'
                    this.curTabName = 'host'
                    // 处理拓扑关系
                    this.setTopoRelation()
                    let self = this
                    setTimeout(() => {
                        self.$refs.field.init()
                    }, 300)
                } else {
                    this.cancel()
                }
            },
            /*
                获取拓扑图信息
            */
            getTopogical () {
                this.changeClassify()
            },
            getTopogical2 () {
                let params = {
                    bk_classification_id: this.curClassify['bk_classification_id']
                }
                let p = new Promise((resolve, reject) => {
                    this.$axios.post(`object/classification/${this.bkSupplierAccount}/objects`, params).then(res => {
                        if (res.result) {
                            resolve(res.data)
                        } else {
                            reject(res['bk_error_msg'])
                        }
                    })
                })
                return p
            },
            /*
                取消按钮
            */
            cancel () {
                this.slider.isBusinessShow = false
            },
            /*
                更新topo图
            */
            updateTopo () {
                this.closeBusiness()
                this.getTopogical()
            },
            /*
                查询模型分类
                index: 查完后当前显示项
            */
            getClassifyList (index) {
                // 默认查全部
                let params = {}
                this.isClassifyLoading = true
                this.isTopoLoading = true
                this.$axios.post('object/classifications', params).then((res) => {
                    if (res.result) {
                        this.classifyList = res.data
                        if (!index) {
                            index = 0
                        }
                        this.curClassify = this.classifyList[index]
                        this.curTempClassify = {
                            id: this.classifyList[index]['id'],
                            bk_classification_type: this.classifyList[index]['bk_classification_type'],
                            bk_classification_name: this.classifyList[index]['bk_classification_name'],
                            bk_classification_id: this.classifyList[index]['bk_classification_id'],
                            bk_classification_icon: this.classifyList[index]['bk_classification_icon']
                        }
                        this.changeClassify(this.curClassify)
                    } else {
                        this.$bkInfo({
                            statusOpts: {
                                title: res['bk_error_msg'],
                                subtitle: false
                            },
                            type: 'error'
                        })
                    }
                    this.isClassifyLoading = false
                })
            },
            /*
                删除模型
                item: 当前项
            */
            deleteModel (item) {
                // 更新拓扑
                this.updateTopo()
            },
            /*
                切换当前显示的图表类型
                enable为已启用模型
                disable为未启用模型
            */
            changeModelDiagram (type) {
                this.isChangeClassify = false
                this.curModelDiagram = type
            },
            /* *****************模型拓扑相关回调 start********************* */
            /*
                节点点击回调
            */
            nodeClick (item, isDisable = false) {
                this.editModel(item, isDisable)
            },
            /* *****************模型拓扑相关回调 end********************* */
            init () {
                // 默认查全部
                let params = {}
                this.isClassifyLoading = true
                this.$axios.post('object/classifications', params).then((res) => {
                    if (res.result) {
                        this.classifyList = res.data
                        let curClassifyIndex = 0
                        if (this.isAutoRouter) {
                            res.data.forEach((classify, classifyIndex) => {
                                if (classify.bk_classification_id === this.$route.query.bk_classification_id) {
                                    curClassifyIndex = classifyIndex
                                }
                            })
                        }
                        this.curClassify = this.classifyList[curClassifyIndex]
                        this.curTempClassify = {
                            id: this.classifyList[curClassifyIndex]['id'],
                            bk_classification_type: this.classifyList[curClassifyIndex]['bk_classification_type'],
                            bk_classification_name: this.classifyList[curClassifyIndex]['bk_classification_name'],
                            bk_classification_id: this.classifyList[curClassifyIndex]['bk_classification_id'],
                            bk_classification_icon: this.classifyList[curClassifyIndex]['bk_classification_icon']
                        }
                        this.changeClassify(res.data[curClassifyIndex])
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.isClassifyLoading = false
                })
            }
        },
        mounted () {
            this.list = iconList
            this.init()
        },
        components: {
            vTopo,
            vSideslider,
            vBaseInfo,
            vField,
            vLayout,
            vOther,
            vPop,
            vGlobalModels
        }
    }
</script>

<style lang="scss" scoped>
    .global-models{
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 50px;
        line-height: 50px;
        text-align: center;
        background-color: #f7fafe;
        &.active{
            background-color: #e2efff;
            .btn-global{
                color: #3578da;
                border-color: currentColor;
            }
        }
        .btn-global{
            width: 120px;
            height: 32px;
            padding: 0;
            line-height: 30px;
            background-color: #ffffff;
            border-radius: 2px;
            border: solid 1px #d6d8df;
            font-size: 14px;
            color: $textColor;
            outline: 0;
            padding: 0;
            .icon-cc-fullscreen,
            .text{
                display: inline-block;
                vertical-align: middle;
            }
        }
    }
</style>
<style media="screen" lang="scss" scoped>
    $borderColor: #bec6de; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #f9f9f9; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; // 主要颜色
    $successColor: #30d878;
    $successColorHover: #27e97a;
    $successColorActive: #29d272;
    .model-wrapper{
        height: 100%;
        overflow: hidden;
        .slide-content{
            padding: 8px 20px 20px;
        }
        .cancel-btn-sider{
            color: $primaryHoverColor;
        }
        .del-btn{
            min-width: 124px;
            background: #ffffff;
            border: 1px solid #e7e9ef;
            display: inline-block;
            border-radius: 1px;
            height: 36px;
            line-height:32px;
            text-align:center;
            color:#bec6de;
            cursor:pointer;
            i{
                color: #bec6de;
                font-style: normal;
                font-size: 16px;
                vertical-align:middle;
                margin-right:10px;
            }
            span{
                vertical-align:middle;
            }
        }
        .left-tap-contain{
            width:188px;
            float:left;
            border-left: none;
            border-top: none;
            height: 100%;
            position: relative;
            border-right: 1px solid #dde4eb;
            .list-tap{
                height: calc(100% - 50px);
                overflow-y: auto;
                &::-webkit-scrollbar{
                    width: 6px;
                    height: 5px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 20px;
                    background: #a5a5a5;
                }
                ul{
                    >li{
                        height: 48px;
                        line-height: 48px;
                        padding: 0 30px 0 44px;
                        width: 100%;
                        cursor: pointer;
                        font-size: 14px;
                        color: #737987;
                        font-size: 14px;
                        position: relative;
                        white-space:nowrap;
                        text-overflow:ellipsis;
                        -o-text-overflow:ellipsis;
                        overflow: hidden;
                        i{
                            font-size: 16px;
                        }
                        .icon-left{
                            margin-left: -12px;
                        }
                        &:hover{
                            color: #3c96ff;
                            background: #f1f7ff;
                        }
                        .text{
                            padding:0 3px 0 5px;
                            min-width:64px;
                            vertical-align: top;
                        }
                        &.active{
                            color: #3c96ff;
                            background: #e2efff;
                        }
                    }
                }
            }
            .bottom-btn-contain{
                width: 148px;
                height:32px;
                background: #fff;
                cursor:pointer;
                font-size:0;
                margin: 10px auto;
                .bottom-btn{
                    display: block;
                    height: 32px;
                    line-height: 30px;
                    color: $primaryHoverColor;
                    border-radius: 2px;
                    color: #c3cdd7;
                    border: dashed 1px #c3cdd7;
                    text-align: center;
                    font-size: 14px;
                    &:hover{
                        border-color: #3c96ff;
                        color: #3c96ff;
                    }
                }
            }
        }
        .right-contain{
            width: calc(100% - 188px);
            float: left;
            height: 100%;
            .top-contain{
                height:56px;
                line-height:56px;
                width:100%;
                background: #fff;
                .del-type{
                    font-size: 0;
                    padding-right: 16px;
                    .icon-content{
                        width: 34px;
                        height: 27px;
                        line-height: 27px;
                        text-align: center;
                        display: inline-block;
                        cursor: pointer;
                        border-radius: 2px;
                    }
                    .icon-content-edit{
                        &:hover{
                            background: #498fe0;
                            i{
                                color: #fff;
                            }
                        }
                    }
                    .icon-content-del{
                        &:hover{
                            background: #ef4c4c;
                            i{
                                color: #fff;
                            }
                        }
                    }
                    .icon-cc-edit{
                        &:hover{
                            color: #042244;
                        }
                    }
                    .icon-cc-del{
                        &:hover{
                            color: #ef4c4c;
                        }
                    }
                }
            }
            .model-box{
                height: 100%;
                position: relative;
                z-index: 2;
                background-color: #f4f5f8;
                background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
                background-size: 10px 10px;
                .model-topo-box{
                    height: 100%;
                }
            }
            .model-diagram{
                overflow-y: auto;
                position: relative;
                z-index: 2;
                height: 100%;
                &::-webkit-scrollbar {
                    width: 6px;
                    height: 5px;
                    &-thumb {
                        border-radius: 20px;
                        background: #a5a5a5;
                        box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
                    }
                }
                .title{
                    color: #498fe0;
                    font-size:14px;
                    line-height:1;
                    text-align:left;
                    font-weight: bold;
                    .inner-model-title{
                        position: relative;
                        padding-left: 25px;
                        color: #6b7baa;
                        font-weight: normal;
                        font-size: 12px;
                        &:before{
                            content: '';
                            width: 5px;
                            height: 5px;
                            border: 1px solid #498fe0;
                            border-radius: 50%;
                            color: #498fe0;
                            font-weight: bold;
                            background: #498fe0;
                            top: 6px;
                            left: 35px;
                            position: absolute;
                        }
                    }
                }
                .new-model-contain{
                    position: absolute;
                    right: 20px;
                    &.new-model-change{
                        right: 252px;
                    }
                }
                .new-model-btn{
                    min-width: 124px;
                    display: inline-block;
                    border-radius: 1px;
                    height: 32px;
                    line-height:32px;
                    text-align:center;
                    color:#bec6de;
                    cursor:pointer;
                    font-size:14px;
                    color: #fff;
                }
            }
            .model-diagram{
                .topo-btn{
                    position: absolute;
                    width: 30px;
                    height: 30px;
                    line-height: 30px;
                    top: 9px;
                    padding: 0;
                    cursor: pointer;
                    border-radius: 50%;
                    box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.2);
                    border: none;
                    text-align: center;
                    font-size: 0;
                    background: #fff;
                    .icon{
                        font-size: 14px;
                        color: #737987;
                    }
                    &.edit{
                        left: 15px;
                        &:hover{
                            .icon{
                                color: #498fe0;
                            }
                        }
                    }
                    &.del{
                        right: 9px;
                        &:hover{
                            .icon{
                                color: #ef4c4c;
                            }
                        }
                    }
                }
                .model-content{
                    position: relative;
                    .topo-wrapper{
                        li{
                            // color: #498fe0!important;
                            /* overflow: hidden; */
                            text-overflow: ellipsis;
                            white-space: nowrap;
                            &.default{
                                border-style: solid;
                                border-width: 1px;
                                background: transparent;
                                background: #fff !important;
                                color: #d6d8df !important;
                                cursor: default;
                            }
                        }
                    }
                }
            }
            .model-diagram{
                .topo-btn{
                    position: absolute;
                    width: 30px;
                    height: 30px;
                    line-height: 30px;
                    top: 9px;
                    padding: 0;
                    cursor: pointer;
                    border-radius: 50%;
                    box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.2);
                    border: none;
                    text-align: center;
                    font-size: 0;
                    background: #fff;
                    .icon{
                        font-size: 14px;
                        color: #737987;
                    }
                    &.edit{
                        left: 15px;
                        &:hover{
                            .icon{
                                color: #498fe0;
                            }
                        }
                    }
                    &.del{
                        right: 9px;
                        &:hover{
                            .icon{
                                color: #ef4c4c;
                            }
                        }
                    }
                }
                .model-content{
                    position: relative;
                    .topo-wrapper{
                        li{
                            position: relative;
                            float: left;
                            width: 91px;
                            height: 91px;
                            border: 1px solid #d6d8df;
                            box-shadow: 0 0 10px transparent;
                            text-align: center;
                            margin-right: 20px;
                            margin-top: 30px;
                            border-radius: 50%;
                            padding: 0 5px;
                            cursor: pointer;
                            color: #498fe0;
                            background: #fff;
                            font-size: 12px;
                            font-weight: bold;
                            /*&.spacing-control{
                            }*/
                            &.line{
                                float: none;
                                margin-top: 60px;
                                &:first-child{
                                    margin-top: 0;
                                }
                                &::after{
                                    content: "";
                                    height: 60px;
                                    position: absolute;
                                    left: 50%;
                                    top: 89px;
                                    border: 1px dashed #d6d8df;
                                }
                                &:last-child{
                                    &::after{
                                        content: "";
                                        height: 100%;
                                        position: absolute;
                                        left: 56px;
                                        top: 50px;
                                        border:none;
                                    }
                                }
                            }
                            &:not(.default):hover{
                                border: 1px solid #d6d8df;
                                box-shadow: 0 2.8px 0 rgba(12, 34, 59, 0.05)
                            }
                            .content{
                                white-space: nowrap;
                                text-overflow: ellipsis;
                                overflow: hidden;
                                padding-top: 20px;
                                .topo-icon{
                                    font-size: 25px;
                                }
                                .content-name{
                                    margin-top: 10px;
                                    line-height: 1;
                                    white-space: nowrap;
                                    text-overflow: ellipsis;
                                    overflow: hidden;
                                }
                            }
                            &.locks-icon-content{
                                background: #6b7baa;
                                color: #fff;
                                .icon-add{
                                    color: #498fe0;
                                }
                            }
                            .icon-add{
                                position: absolute;
                                display: inline-block;
                                border-radius: 50%;
                                top: -41px;
                                padding: 1px 0;
                                font-size: 18px;
                                background: #eee;
                                &:hover{
                                    color: #50abff;
                                }
                                &.prev,&.next{
                                    left: 36px;
                                }
                            }
                        }
                    }
                }
                .no-model-prompting{
                    padding-top: 218px;
                    >img{
                        display: inline-block;
                        width: 200px;
                        margin-left: 10px;
                    }
                    .create-btn{
                        width: 208px;
                        // height: 38px;
                        // line-height: 36px;
                        // text-align: center;
                        // color: #ffffff;
                        // font-size: 14px;
                        // padding: 0 76px;
                        // border: none;
                        // border-radius: 2px;
                    }
                    p{
                        font-size: 14px;
                        color: #6b7baa;
                        margin: 0;
                        line-height: 14px;
                        margin-top: 23px;
                        margin-bottom: 20px;
                    }
                }
            }
        }
        .pop-content{ /*弹窗*/
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.6);
            z-index: 1300;
            .pop-detail{
                width: 610px;
                height: 400px;
                top: 50%;
                left: 50%;
                margin-left: -305px;
                margin-top: -200px;
                background: #fff;
                position: absolute;
                padding: 50px 0;
                .del-text{
                    position: absolute;
                    top: 0;
                    right: 0;
                    width: 27px;
                    height: 27px;
                    line-height: 26px;
                    border-radius: 50%;
                    text-align: center;
                    margin: 4px 4px 0 0;
                    border-radius: 50%;
                    background-repeat: no-repeat;
                    background-size: 11px 11px;
                    background-position: 50% 50%;
                    cursor: pointer;
                    display: inline-block;
                    &:hover{
                        background-color: #f3f3f3;
                    }
                    >i{
                        font-size: 10px;
                        color: #bec6de;
                    }
                }
                .form-wrapper{
                    color:　$primaryHoverColor;
                    font-size: 14px;
                    >h3{
                        font-size: 24px;
                        margin: 0;
                        padding: 0;
                        line-height: 24px;
                        font-weight: normal;
                        text-align: center;
                        margin-bottom: 25px;
                    }
                    .button-content{
                        margin-top: 20px;
                        margin-left: 10px;
                        button{
                            height: 36px;
                            line-height: 36px;
                            padding: 0 46px;
                            border-radius: 2px;
                            outline: none;
                            &.confirm-btn{
                                color: #Fff;
                                margin-right: 15px;
                            }
                        }
                    }
                    form{
                        text-align: left;
                        .from-label{
                            display: inline-block;
                            width: 174px;
                            padding: 0 12px 0 20px;
                            text-align: right;
                            font-size: 14px;
                            line-height: 36px;
                            vertical-align: top;
                            margin-top: 20px;
                        }
                        .form-contain{
                            display: inline-block;
                            margin-top: 20px;
                            .bk-select{
                                width: 100px!important;
                            }
                            .from-input{
                                width:270px;
                                height: 36px;
                                line-height: 36px;
                                border: 1px solid $borderColor;
                                outline: none;
                                padding: 0 15px;
                            }
                            .select-icon-content{ /* 图标选择下拉 */
                                display: inline-block;
                                position: relative;
                                .select-icon-show{
                                    /*width: 70px;*/
                                    height: 36px;
                                    line-height: 34px;
                                    border: 1px solid $borderColor;
                                    cursor: pointer;
                                    .icon-content{
                                        width: 47px;
                                        height: 100%;
                                        line-height: 36px;
                                        border-right: 1px solid $borderColor;
                                        text-align: center;
                                        float: left;
                                        >i{
                                            font-size: 24px;
                                            position: relative;
                                            top: -3px;
                                        }
                                    }
                                    .arrow{
                                        display: inline-block;
                                        width: 20px;
                                        height: 36px;
                                        text-align: center;
                                        line-height: 34px;
                                        vertical-align: bottom;
                                        i{
                                            color: #6b7baa;
                                            font-size: 12px;
                                        }
                                    }
                                }
                                .select-icon-list{
                                    padding: 10px;
                                    position: absolute;
                                    top: 44px;
                                    left: 0;
                                    width: 382px;
                                    min-height: 206px;
                                    border: 1px solid $borderColor;
                                    z-index: 500;
                                    background: #fff;
                                    ul{
                                        padding: 0;
                                        margin: 0;
                                        li{
                                            width: 60px;
                                            height: 46px;
                                            text-align: center;
                                            line-height: 46px;
                                            float: left;
                                            cursor: pointer;
                                            i{
                                                font-size: 24px;
                                                vertical-align: middle;
                                            }
                                            &:hover{
                                                background: #e2efff;
                                            }
                                            &:nth-child(6n){
                                                margin-right: 0;
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }

        }
        .content-button{
           background: #f9f9f9;
           height: 62px;
           padding: 14px 20px;
           font-size: 0;
           .btn{
               font-size: 14px;
               width: 110px;
               height: 34px;
               line-height: 32px;
               border-radius: 2px;
               margin-right: 10px;
           }
           .info{
               float: right;
               font-size: 14px;
               height: 34px;
               line-height: 34px;
               cursor: pointer;
               min-width:110px;
               background:#ffffff;
               display:inline-block;
               text-align:center;
               i{
                   font-style:normal;
               }
           }
        }
        .bk-select-input{
           height: 30px !important;
           line-height: 30px !important;
        }
        .bk-tab2{
           border: none;
           height: 100%;
        }
        .bk-tab2-pane{
           width:100%;
        }
        .is-disabled{
            cursor: not-allowed !important;
        }
    }


    .model-type-name{
        padding: 0 28px;
        display: inline-block;
        vertical-align: top;
        font-size: 14px;
        font-weight: bold;
        color: #6b7baa;
        border-left: 2px solid #6b7baa;
        height: 56px;
    }
    .model-type-more{
        position: relative;
        display: inline-block;
        vertical-align: top;
        width: 18px;
        height: 56px;
        margin-left: 40px;
        background: url(../../common/images/icon/icon-dot.png) no-repeat center center;
        &:hover{
            background-image: url(../../common/images/icon/icon-dot-hover.png);
            .model-type-more-content{
                display: block;
            }
        }
        .model-type-more-content{
            display: none;
            position: absolute;
            left: -30px;
            top: 50px;
            padding: 20px 30px;
            line-height: 24px;
            background-color: #fff;
            box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.1);
            z-index: 3;
            .more-content-list{
                white-space: nowrap;
                font-size: 12px;
                .more-content-list-label{
                    display: inline-block;
                    vertical-align: middle;
                    /* width: 48px; */
                    max-width: 48px; //临时使用，更多属性显示后恢复成width
                    text-align: right;
                    color: #bec6de;
                }
                .more-content-list-text{
                    display: inline-block;
                    vertical-align: middle;
                    color: #6b7baa;
                    margin-left: 8px;
                }
            }
        }
    }
    .topo-wrapper-bk_biz_topo{
        padding-top: 80px;
        li{
            margin-left: auto;
            margin-right: auto !important;
        }
    }
    .vis-button{
        position: absolute;
        width: 30px;
        height: 30px;
        line-height: 30px;
        top: 9px;
        padding: 0;
        cursor: pointer;
        border-radius: 50%;
        box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.2);
        border: none;
        text-align: center;
        font-size: 0;
        [class^="icon-cc-"],
        .vis-button-text{
            height: 100%;
            line-height: 30px;
            display: inline-block;
            vertical-align: middle;
            font-size: 12px;
            margin: 0 3px;
        }
    }
    .vis-enable{
        width: auto;
        right: 49px;
        padding: 0 10px;
        border-radius: 15px;
        font-weight: normal;
        text-align: left;
        font-size: 0;
        .bk-icon{
            font-weight: normal;
            vertical-align: middle;
            font-size: 14px;
        }
    }
    .topo-disable{
        position: absolute;
        right: 49px;
        top: 9px;
        width: 169px;
        background-color: #ffffff;
        box-shadow: 0px 1px 5px 0px 
            rgba(12, 34, 59, 0.2);
        border-radius: 2px;
        line-height: 40px;
        font-size: 14px;
        transform-origin: right top;
        transition: transform .3s ease-in-out;
        .disable-title{
            display: block;
            background-color: #fafbfd;
            padding: 0 0 0 22px;
            position: relative;
            .bk-icon{
                position: absolute;
                top: 12px;
                right: 12px;
                display: inline-block;
                font-size: 16px;
                transform: rotate(-45deg);
                cursor: pointer;
            }
        }
        .disable-list{
            padding: 6px;
            overflow: auto;
            min-height: 52px;
            @include scrollbar;
            .disable-item{
                .disable-item-link{
                    display: block;
                    padding: 0 0 0 16px;
                    color: inherit;
                    text-decoration: none;
                    transition: unset !important;
                    @include ellipsis;
                    &:hover{
                        background-color: #f1f7ff;
                        color: #3c96ff;
                    }
                    &:active{
                        background-color: #e2efff;
                        color: #3c96ff
                    }
                }
            }
        }
    }
    .topo-disable-list-enter,
    .topo-disable-list-leave-to{
        transform: scaleX(0) scaleY(0);
    }
</style>

<style media="screen" lang="scss">
    .model-wrapper{
        .bk-tab2-content{
            height: calc(100% - 58px);
            /*overflow-y: auto;*/
        }
    }
</style>
