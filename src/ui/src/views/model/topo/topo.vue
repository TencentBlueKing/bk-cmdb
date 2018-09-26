/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */
<template>
    <div class="topo-box model" :class="{'no-edit': bkClassificationId === 'bk_host_manage'}" v-bkloading="{isLoading: isLoading}">
        <div ref="topo" class="topo" v-show="modelList.length !== 0"></div>
        <button v-if="bkClassificationId !== 'bk_host_manage'" class="bk-button vis-button vis-setup" @click="editClassify" v-tooltip="$t('Common[\'编辑\']')">
            <i class="icon icon-cc-edit"></i>
        </button>
        <button class="bk-button vis-button vis-zoomExtends bk-icon icon-full-screen" @click="resizeFull" v-tooltip="$t('ModelManagement[\'还原\']')">
        </button>
        <button class="bk-button vis-button vis-zoomIn bk-icon icon-plus" @click="zoomIn" v-tooltip="$t('ModelManagement[\'放大\']')">
        </button>
        <button class="bk-button vis-button vis-zoomOut bk-icon icon-minus" @click="zoomOut" v-tooltip="$t('ModelManagement[\'缩小\']')">
        </button>
        <div class="has-model" v-if="modelList.length !== 0">
            <bk-button type="primary" v-if="addModelAvailable" class="bk-button vis-button vis-create" @click="createModel">
                <span class="vis-button-text">{{$t('Common["新增"]')}}</span>
            </bk-button>
        </div>
        <div class="no-model-prompting tc" v-if="!modelList.length">
            <img src="../../../assets/images/no_model_prompting.png">
            <p v-if="modelList.length === 0 && disableModelList.length">{{$t('ModelManagement["此分组下无已启用模型"]')}}</p>
            <p v-else>{{$t('ModelManagement["此分组下无模型"]')}}</p>
            <bk-button type="primary" class="create-btn" @click="createModel">{{$t('Common["立即创建"]')}}</bk-button>
        </div>
        <bk-button class="bk-button vis-button vis-enable" v-if="addModelAvailable && disableModelList.length" @click="isShowDisableList = true">
            <i class="bk-icon icon-minus-circle-shape"></i>
            <span class="vis-button-text">{{disableModelList.length}}</span>
        </bk-button>
        <transition name="topo-disable-list">
            <div class="topo-disable" v-show="isShowDisableList">
                <label class="disable-title">
                    <span>{{$t('ModelManagement["已停用模型"]')}}</span>
                    <i class="bk-icon icon-arrows-right" @click="isShowDisableList = false"></i>
                </label>
                <ul class="disable-list" ref="disableList">
                    <li class="disable-item" v-for="(model, index) in disableModelList" :key="index">
                        <a class="disable-item-link" href="javascript:void(0)" @click="editModel(model)">{{model['bk_obj_name']}}</a>
                    </li>
                </ul>
            </div>
        </transition>
        <bk-button type="default" class="vis-button vis-del" v-tooltip="$t('ModelManagement[\'删除\']')" v-if="!isInnerType" @click="deleteClassify">
            <i class="icon icon-cc-del"></i>
        </bk-button>
    </div>
</template>

<script>
    import Vis from 'vis'
    import { mapGetters, mapActions, mapMutations } from 'vuex'
    import { generateObjIcon as GET_OBJ_ICON } from '@/utils/util'
    export default {
        data () {
            return {
                isLoading: false,
                networkInstance: null,
                networkDataSet: {
                    nodes: null,
                    edges: null
                },
                network: {
                    nodes: null,
                    edges: null,
                    options: {
                        nodes: {
                            shape: 'image',
                            size: 45,
                            widthConstraint: 55,
                            shadow: {
                                enabled: true,
                                color: 'rgba(0,0,0,0.1)',
                                x: 0,
                                y: 2,
                                size: 4
                            }
                        },
                        edges: {
                            color: {
                                color: '#c3cdd7',
                                highlight: '#3c96ff'
                            },
                            smooth: {
                                type: 'curvedCW',
                                roundness: 0
                            },
                            arrows: {
                                to: {
                                    scaleFactor: 0.6
                                },
                                from: {
                                    scaleFactor: 0.6
                                }
                            }
                        },
                        physics: {
                            enabled: true,
                            barnesHut: {
                                avoidOverlap: 0.5,
                                springLength: 150
                            }
                        }
                    }
                },
                topoStructure: [],
                disableModelList: [],
                isShowDisableList: false
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            activeClassify () {
                let activeClassify = this.classifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === this.bkClassificationId)
                return activeClassify
            },
            isInnerType () {
                return this.$classifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === this.bkClassificationId)['bk_classification_type'] === 'inner'
            },
            addModelAvailable () {
                return !['bk_biz_topo', 'bk_host_manage', 'bk_organization'].includes(this.bkClassificationId)
            },
            bkClassificationId () {
                return this.$route.params.classifyId
            },
            modelList () {
                let currentClassify = this.$classifications.find(({bk_classification_id: bkClassificationId}) => this.bkClassificationId === bkClassificationId)
                let modelList = currentClassify['bk_objects'].filter(({bk_ispaused: bkIspaused}) => !bkIspaused)
                this.disableModelList = currentClassify['bk_objects'].filter(({bk_ispaused: bkIspaused}) => bkIspaused)
                let asstList = []
                for (let key in currentClassify['bk_asst_objects']) {
                    let object = currentClassify['bk_asst_objects'][key]
                    object.map(asstModel => {
                        let model = asstList.find(({bk_classification_id: bkClassificationId}) => {
                            return bkClassificationId === asstModel['bk_classification_id']
                        })
                        
                        this.topoStructure.map(structure => {
                            let isExist = -1
                            if (asstModel['bk_obj_id'] === structure.to['bk_obj_id']) {
                                isExist = modelList.findIndex(model => {
                                    return model['bk_obj_id'] === structure.from['bk_obj_id']
                                })
                            } else if (asstModel['bk_obj_id'] === structure.from['bk_obj_id']) {
                                isExist = modelList.findIndex(model => {
                                    return model['bk_obj_id'] === structure.to['bk_obj_id']
                                })
                            }
                            if (isExist !== -1 && asstList.findIndex(({bk_obj_id: objId}) => objId === asstModel['bk_obj_id']) === -1 && modelList.findIndex(({bk_obj_id: objId}) => objId === asstModel['bk_obj_id']) === -1) {
                                asstList.push(asstModel)
                            }
                        })
                    })
                }
                return [...modelList, ...asstList]
            },
            noPositionNodes () {
                return this.network.nodes.filter(node => {
                    const position = node.data.position
                    return position.x === null && position.y === null
                })
            }
        },
        watch: {
            '$route.params.classifyId' () {
                this.initTopo()
            }
        },
        methods: {
            ...mapActions('objectModelClassify', [
                'deleteClassification'
            ]),
            ...mapMutations('objectModelClassify', [
                'deleteClassify'
            ]),
            async deleteClassify () {
                this.$bkInfo({
                    title: this.$t('ModelManagement["确认要删除此分组？"]'),
                    confirmFn: async () => {
                        await this.deleteClassification({
                            id: this.activeClassify['id']
                        })
                        this.$router.push('/model/bk_host_manage')
                        this.$store.commit('objectModelClassify/deleteClassify', this.bkClassificationId)
                    }
                })
            },
            editClassify () {
                this.$emit('editClassify')
            },
            createModel () {
                this.$emit('createModel')
            },
            editModel (model) {
                this.$emit('editModel', model)
            },
            resizeFull () {
                this.networkInstance.moveTo({scale: 1})
            },
            zoomIn () {
                let scale = this.networkInstance.getScale()
                scale += 0.05
                this.networkInstance.moveTo({scale: scale})
            },
            zoomOut () {
                let scale = this.networkInstance.getScale()
                if (scale > 0.05) {
                    scale -= 0.05
                }
                this.networkInstance.moveTo({scale: scale})
            },
            async initTopo () {
                this.isLoading = true
                await this.setEdges()
                this.setNodes()
                this.networkInstance = new Vis.Network(this.$refs.topo, {
                    nodes: this.networkDataSet.nodes,
                    edges: this.networkDataSet.edges
                }, this.network.options)
                this.addListener()
                this.isLoading = false
            },
            async getTopoStructure () {
                const res = await this.$store.dispatch('objectModel/searchObjectTopo', {params: {bk_classification_id: this.bkClassificationId}})
                let topoStructure = []
                res.map(structure => {
                    let index = topoStructure.findIndex(item => {
                        if ((item.to['bk_obj_id'] === structure.to['bk_obj_id'] && item.from['bk_obj_id'] === structure.from['bk_obj_id']) || (item.to['bk_obj_id'] === structure.from['bk_obj_id'] && item.from['bk_obj_id'] === structure.to['bk_obj_id'])) {
                            return true
                        }
                    })
                    if (index === -1) {
                        topoStructure.push(structure)
                    } else {
                        topoStructure[index]['label_name'] += `,${structure['label_name']}`
                    }
                })
                this.topoStructure = topoStructure
            },
            setNodes () {
                let nodes = []
                this.network.nodes = this.modelList.map(nodeData => {
                    let fontColor = nodeData['ispre'] ? '#6894c8' : '#868b97'
                    if (nodeData['bk_classification_id'] !== this.bkClassificationId || nodeData['bk_obj_id'] === 'plat') {
                        fontColor = '#c3cdd7'
                    }
                    let node = {
                        id: nodeData['bk_obj_id'],
                        image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON({
                            name: nodeData['bk_obj_name'],
                            backgroundColor: '#fff',
                            fontColor
                        }))}`,
                        data: nodeData
                    }
                    if (nodeData.hasOwnProperty('position') && nodeData['position'] !== '') {
                        let position = JSON.parse(nodeData['position'])
                        if (position.hasOwnProperty(this.bkClassificationId)) {
                            node.physics = false
                            node.x = position[this.bkClassificationId].x
                            node.y = position[this.bkClassificationId].y
                        }
                    }
                    return node
                })
                this.networkDataSet.nodes = new Vis.DataSet(this.network.nodes)
            },
            async setEdges () {
                await this.getTopoStructure()
                let edges = []
                this.topoStructure.map((edge, index) => {
                    edges.push({
                        dashes: !(edge['from']['bk_classification_id'] === this.bkClassificationId && edge['to']['bk_classification_id'] === this.bkClassificationId),
                        label: edge['label_name'],
                        arrows: edge['arrows'],
                        to: edge['to']['bk_obj_id'],
                        from: edge['from']['bk_obj_id']
                    })
                })
                this.network.edges = edges
                this.networkDataSet.edges = new Vis.DataSet(this.network.edges)
            },
            // 加载节点icon并更新
            loadNodeImage () {
                this.network.nodes.forEach(node => {
                    let image = new Image()
                    image.onload = () => {
                        let fontColor = node.data['ispre'] ? '#6894c8' : '#868b97'
                        let iconColor = node.data['ispre'] ? '#6894c8' : '#868b97'
                        if (node.data['bk_classification_id'] !== this.bkClassificationId || node.data['bk_obj_id'] === 'plat') {
                            fontColor = '#c3cdd7'
                            iconColor = '#c3cdd7'
                        }
                        this.networkDataSet.nodes.update({
                            id: node.id,
                            image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: node.data['bk_obj_name'],
                                fontColor,
                                iconColor,
                                backgroundColor: '#fff'
                            }))}`
                        })
                    }
                    image.src = `${window.location.origin}/static/svg/${node['data']['bk_obj_icon'].substr(5)}.svg`
                })
            },
            async updateNodePosition (nodeId) {
                let model = this.modelList.find(({bk_obj_id: bkObjId}) => {
                    return bkObjId === nodeId
                })
                if (model) {
                    let params = {
                        bk_supplier_account: this.supplierAccount,
                        bk_classification_id: this.bkClassificationId,
                        position: ''
                    }
                    let pos = {}
                    const nodePositions = this.networkInstance.getPositions(nodeId)[nodeId]
                    if (model.hasOwnProperty('position') && model.position !== '') {
                        pos = JSON.parse(model.position)
                    }
                    pos[this.bkClassificationId] = nodePositions
                    params.position = JSON.stringify(pos)
                    await this.$store.dispatch('objectModel/updateObject', {id: model.id, params})
                    let updateModel = {
                        bk_classification_id: this.bkClassificationId,
                        bk_obj_id: nodeId,
                        position: params.position
                    }
                    this.$store.commit('objectModelClassify/updateModel', updateModel)
                }
            },
            // 拓扑稳定后执行事件
            // 1.取消物理模拟
            // 2.配置拖拽结束监听，更新位置信息
            // 3.设置无位置信息的单节点位置
            // 4.加载节点图标
            listenerCallback () {
                this.networkInstance.setOptions({
                    physics: {
                        enabled: false
                    }
                })
                this.networkInstance.on('dragEnd', (params) => {
                    if (params.nodes.length) {
                        this.updateNodePosition(params.nodes[0])
                    }
                    this.networkInstance.unselectAll()
                })
                this.networkInstance.on('hoverNode', () => {
                    this.$refs.topo.style.cursor = 'pointer'
                })
                this.networkInstance.on('blurNode', () => {
                    this.$refs.topo.style.cursor = 'default'
                })
                this.networkInstance.on('click', (params) => {
                    // 点击了具体某个节点
                    if (params.nodes.length) {
                        let id = params.nodes[0]
                        let model = this.modelList.find(({bk_obj_id: bkObjId}) => bkObjId === id)
                        if (model && model['bk_classification_id'] === this.bkClassificationId && id !== 'plat') {
                            this.$emit('editModel', model)
                        }
                    }
                })
                this.loadNodeImage()
                this.networkInstance.fit()
                this.loading = false
            },
            addListener () {
                const networkInstance = this.networkInstance
                networkInstance.once('stabilized', this.listenerCallback)
                if (!this.noPositionNodes.length) {
                    this.listenerCallback()
                }
            }
        },
        mounted () {
            this.initTopo()
        }
    }
</script>

<style lang="scss" scoped>
    .topo-box{
        position: relative;
        width: 100%;
        height: 100%;
        .topo{
            height: 100%;
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
        &:not(.vis-create){
            background: #fff;
        }
        [class^="icon-cc-"],
        .vis-button-text{
            height: 100%;
            line-height: 30px;
            display: inline-block;
            vertical-align: middle;
            font-size: 12px;
            margin: 0 3px;
        }
        [class^="icon-cc-"],
        &.bk-icon{
            font-weight: bold;
            font-size: 14px;
        }
        &.vis-zoomIn,
        &.vis-zoomOut,
        &.vis-zoomExtends,
        &.vis-setup,
        &.vis-enable{
            &:hover{
                color: #6eb1ff !important;
            }
            &:active{
                color: #3188ed;
            }
        }
    }
    .model.no-edit{
        .vis-zoomIn{
            left: 54px;
        }
        .vis-zoomOut{
            left: 92px;
        }
        .vis-zoomExtends{
            left: 15px;
        }
    }
    .vis-zoomIn{
        left: 92px;
    }
    .vis-zoomOut{
        left: 132px;
    }
    .vis-zoomExtends{
        left: 54px;
    }
    .vis-up,
    .vis-down,
    .vis-left,
    .vis-right{
        display: none;
    }
    .vis-create{
        left: 170px;
        font-size: 18px;
        border-radius: 14px;
        width: 60px;
        font-size: 0;
    }
    .vis-setup{
        left: 15px;
        .icon{
            position: relative;
            top: -1px;
            font-weight: normal;
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
    .vis-del{
        right: 9px;
        color: #737987;
        .icon{
            font-weight: normal;
        }
        &:hover .icon{
            color: #ef4c4c;
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
    .no-model-prompting{
        padding-top: 218px;
        >img{
            display: inline-block;
            width: 200px;
            margin-left: 10px;
        }
        .create-btn{
            width: 208px;
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
</style>
