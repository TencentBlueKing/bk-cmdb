<template>
    <div class="topo-wrapper" :class="{'has-nav': topoEdit.isEdit}">
        <div class="toolbar">
            <bk-button class="edit-button" type="primary" @click="editTopo">
                {{$t('ModelManagement["编辑拓扑"]')}}
            </bk-button>
            <div class="vis-button-group">
                <bk-button class="vis-button vis-zoomExtends bk-icon icon-full-screen" @click="resizeFull" v-tooltip="$t('ModelManagement[\'还原\']')"></bk-button>
                <bk-button class="vis-button vis-zoomIn bk-icon icon-plus" @click="zoomIn" v-tooltip="$t('ModelManagement[\'放大\']')"></bk-button>
                <bk-button class="vis-button vis-zoomOut bk-icon icon-minus" @click="zoomOut" v-tooltip="$t('ModelManagement[\'缩小\']')"></bk-button>
                <bk-button class="vis-button vis-setting icon-cc-setting" @click="showSlider('theDisplay')" v-tooltip="$t('ModelManagement[\'拓扑显示设置\']')"></bk-button>
                <bk-button class="vis-button vis-example" @click="toggleExample">
                    <span class="vis-button-text">{{$t('ModelManagement["图例"]')}}</span>
                    <i class="bk-icon icon-angle-down" :class="{'rotate': isShowExample}"></i>
                </bk-button>
                <cmdb-collapse-transition name="topo-example-list">
                    <div class="topo-example" v-show="isShowExample">
                        <p class="example-item">
                            <i></i>
                            <span>{{$t('ModelManagement["自定义模型"]')}}</span>
                        </p>
                        <p class="example-item">
                            <i></i>
                            <span>{{$t('ModelManagement["内置模型"]')}}</span>
                        </p>
                    </div>
                </cmdb-collapse-transition>
            </div>
        </div>
        <template v-if="topoEdit.isEdit">
            <div class="topo-save-title">
                <bk-button type="primary" @click="saveTopo">
                    {{$t('Common["保存"]')}}
                </bk-button>
                <bk-button type="primary" @click="exitEdit">
                    {{$t('Common["返回"]')}}
                </bk-button>
            </div>
            <ul class="topo-nav">
                <li class="group-item" v-for="(group, groupIndex) in classifications" :key="groupIndex">
                    <div class="group-info"
                        :class="{'active': topoNav.activeGroup === group['bk_classification_id']}"
                        @click="toggleGroup(group)">
                        <span class="group-name">{{group['bk_classification_name']}}</span>
                        <span class="model-count">{{group['bk_objects'].length}}</span>
                        <i class="bk-icon icon-angle-down"></i>
                    </div>
                    <cmdb-collapse-transition name="model-box">
                        <ul class="model-box" v-show="topoNav.activeGroup === group['bk_classification_id']">
                            <li class="model-item"
                            v-for="(model, modelIndex) in group['bk_objects']"
                            :key="modelIndex"
                            :class="{'disabled': isModelInTopo(model)}"
                            :draggable="!isModelInTopo(model)"
                            @dragstart="handleDragstart(model, $event)">
                                <div v-if="!isModelInTopo(model)">
                                    <i :class="model['bk_obj_icon']"></i>
                                    <div class="info">
                                        <p class="name">{{model['bk_obj_name']}}</p>
                                        <p class="id">{{model['bk_obj_id']}}</p>
                                    </div>
                                </div>
                                <div v-else>
                                    <i :class="model['bk_obj_icon']"></i>
                                    <div class="info">
                                        <p class="name">{{model['bk_obj_name']}}</p>
                                        <p class="id">{{model['bk_obj_id']}}</p>
                                    </div>
                                </div>
                            </li>
                        </ul>
                    </cmdb-collapse-transition>
                </li>
            </ul>
        </template>
        
        <cmdb-slider
            :width="slider.width"
            :isShow.sync="slider.isShow"
            :title="slider.title"
            @close="handleSliderCancel">
            <component 
                class="slider-content"
                slot="content"
                :is="slider.content"
                v-bind="slider.properties"
                @save="handleSliderSave"
                @cancel="handleSliderCancel"
            ></component>
        </cmdb-slider>
        <div class="global-model" @dragover.prevent="" @drop="handleDrop" @mousemove="handleMouseMove" ref="topo" v-bkloading="{isLoading: loading}"></div>
        <svg class="topo-line" v-if="topoEdit.line.x1 && topoEdit.line.x2">
            <line :x1="topoEdit.line.x1" :y1="topoEdit.line.y1" :x2="topoEdit.line.x2" :y2="topoEdit.line.y2" stroke="#c3cdd7" stroke-width="1"></line>
        </svg>
        <ul class="topology-edge-tooltips" ref="edgeTooltips"
            @mouseover="handleEdgeTooltipsOver"
            @mouseleave="handleEdgeTooltipsLeave"
            v-if="topoTooltip.hoverEdge && topoTooltip.hoverEdge.labelList.length > 1">
            <li class="tooltips-option" 
                :key="labelIndex"
                v-for="(labelInfo, labelIndex) in topoTooltip.hoverEdge.labelList"
                @click="handleShowDetails(labelInfo)">
                {{labelInfo.text}}
            </li>
        </ul>
        <div class="topology-node-tooltips"
            ref="nodeTooltips"
            v-if="topoTooltip.hoverNode"
            @mouseover="handleNodeTooltipsOver"
            @mouseleave="handleNodeTooltipsLeave"
            @click="deleteNode">
            <i class="bk-icon icon-close"></i>
        </div>
    </div>
</template>

<script>
    import Vis from 'vis'
    import theDisplay from './topo-detail/display'
    import theRelation from './topo-detail/relation'
    import theRelationDetail from './topo-detail/relation-detail'
    import { generateObjIcon as GET_OBJ_ICON } from '@/utils/util'
    import { mapGetters, mapActions } from 'vuex'
    import throttle from 'lodash.throttle'
    const NAV_WIDTH = 200
    export default {
        components: {
            theDisplay,
            theRelation,
            theRelationDetail
        },
        data () {
            return {
                associationList: [],
                slider: {
                    width: 514,
                    isShow: false,
                    content: '',
                    properties: {},
                    title: this.$t('ModelManagement["拓扑显示设置"]')
                },
                displayConfig: {
                    isShowModelName: true,
                    isShowModelAsst: true
                },
                topoTooltip: {
                    hoverNode: null,
                    hoverNodeTimer: null,
                    hoverEdge: null,
                    hoverEdgeTimer: null
                },
                topoEdit: {
                    isEdit: false,
                    activeEdge: {
                        from: '',
                        to: ''
                    },
                    edges: [],
                    nodes: [],
                    line: {
                        center: {
                            x: 0,
                            y: 0
                        },
                        dragStart: {
                            x1: 0,
                            y1: 0,
                            x2: 0,
                            y2: 0
                        },
                        x1: 0,
                        y1: 0,
                        x2: 0,
                        y2: 0
                    }
                },
                topoNav: {
                    activeGroup: ''
                },
                topoModelList: [],
                localTopoModelList: [],
                isShowExample: false,
                loading: true,
                networkInstance: null,
                networkDataSet: {
                    nodes: null,
                    edges: null
                },
                handleMouseMove: Function,
                network: {
                    nodes: null,
                    edges: null,
                    options: {
                        interaction: {
                            hover: true
                        },
                        manipulation: {
                            enabled: true,
                            addEdge: (data, callback) => {
                                this.topoEdit.activeEdge = data
                                callback(data)
                                this.handleEdgeCreate(data)
                            }
                        },
                        nodes: {
                            shape: 'image',
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
                                highlight: '#3c96ff',
                                hover: '#3c96ff'
                            },
                            font: {
                                background: '#fff'
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
                }
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            noPositionNodes () {
                return this.network.nodes.filter(node => {
                    const position = node.data.position
                    return position.x === null && position.y === null
                })
            }
        },
        watch: {
            'topoEdit.activeEdge.from' (objId) {
                if (objId === '') {
                    this.topoEdit.line.x1 = 0
                    this.topoEdit.line.y1 = 0
                }
            }
        },
        mounted () {
            this.initNetwork()
            this.initMoveFunction()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType',
                'createObjectAssociation',
                'updateObjectAssociation',
                'deleteObjectAssociation'
            ]),
            ...mapActions('objectModel', [
                'deleteObject'
            ]),
            isModelInTopo (model) {
                return this.network.nodes.findIndex(node => node.id === model['bk_obj_id']) > -1
            },
            createAsst (params) {
                return this.createObjectAssociation({
                    params
                })
            },
            updateAsst (params) {
                return this.updateObjectAssociation({
                    params: {
                        id: params.id,
                        bk_obj_asst_name: params['bk_obj_asst_name']
                    }
                })
            },
            deleteAsst (params) {
                return this.deleteObjectAssociation({
                    id: params.id
                })
            },
            clearEditData () {
                this.topoEdit.isEdit = false
                this.topoEdit.activeEdge.from = ''
                this.topoEdit.activeEdge.to = ''
                this.topoEdit.edges = []
                this.topoEdit.nodes = []
            },
            async exitEdit () {
                this.localTopoModelList = this.$tools.clone(this.topoModelList)
                this.topoEdit.isEdit = false
                this.updateNetwork()
                await this.$nextTick()
                this.clearEditData()
            },
            updatePositions () {
                const nodeIds = this.network.nodes.map(({id}) => id)
                const positions = this.networkInstance.getPositions(nodeIds)
                let nodes = []
                this.network.nodes.forEach(({id, x, y}) => {
                    if (positions[id].x !== x || positions[id].y !== y) {
                        nodes.push(id)
                    }
                })
                this.updateNodePosition(this.networkDataSet.nodes.get(nodes))
            },
            async saveTopo () {
                let createAsstArray = []
                let updateAsstArray = []
                let deleteAsstArray = []
                let deleteObjectArray = []
                this.topoEdit.edges.filter(({type}) => type === 'create').forEach(data => {
                    createAsstArray.push(this.createAsst(data.params))
                })
                this.topoEdit.edges.filter(({type}) => type === 'update').forEach(data => {
                    updateAsstArray.push(this.updateAsst(data.params))
                })
                this.topoEdit.edges.filter(({type}) => type === 'delete').forEach(data => {
                    deleteAsstArray.push(this.deleteAsst(data.params))
                })
                this.topoEdit.nodes.filter(({type}) => type === 'delete').forEach(data => {
                    let id = this.$allModels.find(model => model['bk_obj_id'] === data.params.objId).id
                    deleteObjectArray.push(this.deleteObject({id}))
                })
                this.updatePositions()
                await Promise.all(createAsstArray)
                await Promise.all(updateAsstArray)
                await Promise.all(deleteAsstArray)
                await Promise.all(deleteObjectArray)
                this.topoEdit.isEdit = false
                this.topoModelList = this.$tools.clone(this.localTopoModelList)
            },
            handleDisplaySave (displayConfig) {
                this.displayConfig.isShowModelName = displayConfig.isShowModelName
                this.displayConfig.isShowModelAsst = displayConfig.isShowModelAsst
                this.localTopoModelList = displayConfig.topoModelList
                this.updateNetwork()
            },
            handleRelationSave (params) {
                this.topoEdit.edges.push({
                    type: 'create',
                    params
                })
                let fromNode = this.localTopoModelList.find(model => model['bk_obj_id'] === this.topoEdit.activeEdge.from)
                if (!fromNode.hasOwnProperty('assts')) {
                    Object.assign(fromNode, {assts: []})
                }
                fromNode.assts.push({
                    bk_asst_inst_id: this.associationList.find(asst => asst['bk_asst_id'] === params['bk_asst_id']).id,
                    bk_obj_id: params['bk_asst_obj_id'],
                    bk_inst_id: '',
                    checked: true,
                    asstInfo: params
                })
                this.updateNetwork()
            },
            handleRelationDetailSave (data) {
                if (data.type === 'delete') {
                    this.localTopoModelList.forEach(model => {
                        if (model.hasOwnProperty('assts')) {
                            let index = model.assts.findIndex(asst => {
                                if (asst['bk_inst_id'] !== '') {
                                    return asst['bk_inst_id'] === data.params.id
                                } else {
                                    return asst.asstInfo['bk_obj_id'] === data.params['bk_obj_id'] && asst.asstInfo['bk_asst_id'] === data.params['bk_asst_id'] && asst.asstInfo['bk_asst_obj_id'] === data.params['bk_asst_obj_id']
                                }
                            })
                            if (index > -1) {
                                model.assts.splice(index, 1)
                            }
                        }
                    })
                    if (data.params.id === '') {
                        let edgeIndex = this.topoEdit.edges.findIndex(({params}) => {
                            let isExist = true
                            for (let key in params) {
                                if (key !== 'bk_obj_asst_name' && params[key] !== data.params[key]) {
                                    isExist = false
                                }
                            }
                            return isExist
                        })
                        if (edgeIndex > -1) {
                            this.topoEdit.edges.splice(edgeIndex, 1)
                        }
                    } else {
                        this.topoEdit.edges.push(data)
                    }
                } else { // update
                    let edge = this.topoEdit.edges.find(({params}) => {
                        let isExist = true
                        for (let key in params) {
                            if (key !== 'bk_obj_asst_name' && params[key] !== data.params[key]) {
                                isExist = false
                            }
                        }
                        return isExist
                    })
                    if (edge) {
                        edge.params['bk_obj_asst_name'] = data.params['bk_obj_asst_name']
                    } else {
                        this.topoEdit.edges.push(data)
                    }
                }
                this.updateNetwork()
            },
            editTopo () {
                this.topoEdit.isEdit = true
            },
            deleteNode () {
                let {
                    hoverNode
                } = this.topoTooltip
                this.$bkInfo({
                    title: this.$tc('ModelManagement["确定删除模型？"]', hoverNode.label, {name: hoverNode.label}),
                    confirmFn: () => {
                        this.localTopoModelList = this.localTopoModelList.filter(model => model['bk_obj_id'] !== hoverNode.id)
                        this.topoEdit.nodes.push({
                            type: 'delete',
                            params: {
                                objId: hoverNode.id
                            }
                        })
                        this.topoEdit.edges = this.topoEdit.edges.filter(edge => edge.params['bk_obj_id'] !== hoverNode.id && edge.params['bk_asst_obj_id'] !== hoverNode.id)
                        this.topoTooltip.hoverNode = null
                        this.topoTooltip.hoverNodeTimer = null
                        this.updateNetwork()
                    }
                })
            },
            handleDragstart (model, event) {
                event.dataTransfer.setData('objId', model['bk_obj_id'])
            },
            handleDrop (event) {
                let objId = event.dataTransfer.getData('objId')
                let node = this.localTopoModelList.find(model => model['bk_obj_id'] === objId)
                let originPosition = this.networkInstance.getViewPosition()
                let container = this.$refs.topo.getBoundingClientRect()
                let scale = this.networkInstance.getScale()
                node.position.x = originPosition.x - ((container.left + container.right) / 2 - event.clientX) / scale
                node.position.y = originPosition.y - ((container.top + container.bottom) / 2 - event.clientY) / scale
                node.draged = true
                this.updateNetwork()
            },
            clearActiveEdge () {
                this.topoEdit.activeEdge = {
                    from: '',
                    to: ''
                }
            },
            handleSliderSave (params) {
                switch (this.slider.content) {
                    case 'theDisplay':
                        this.handleDisplaySave(params)
                        break
                    case 'theRelation':
                        this.handleRelationSave(params)
                        this.clearActiveEdge()
                        break
                    case 'theRelationDetail':
                        this.handleRelationDetailSave(params)
                        break
                    default:
                }
                this.clearHoverTooltip()
            },
            handleSliderCancel () {
                if (this.slider.content === 'theRelation') {
                    this.clearActiveEdge()
                    this.updateNetwork()
                }
                this.clearHoverTooltip()
                this.slider.isShow = false
            },
            handleEdgeCreate (data) {
                this.slider.properties = {
                    fromObjId: data.from,
                    toObjId: data.to,
                    topoModelList: this.localTopoModelList,
                    edges: this.topoEdit.edges
                }
                this.showSlider('theRelation')
            },
            handleEdgeClick (edgeId) {
                let edge = this.network.edges.find(({id}) => id === edgeId)
                if (edge.labelList.length === 1) {
                    this.handleShowDetails(edge.labelList[0])
                }
            },
            initMoveFunction () {
                this.handleMouseMove = throttle(event => {
                    this.topoEdit.line.x2 = event.layerX
                    this.topoEdit.line.y2 = event.layerY
                }, 50)
            },
            handleNodeClick (data) {
                if (!this.topoEdit.isEdit) {
                    return
                }
                if (this.topoEdit.activeEdge.from === '') {
                    this.topoEdit.activeEdge.from = data['nodes'][0]
                    this.topoEdit.line.x1 = data.pointer.DOM.x
                    this.topoEdit.line.y1 = data.pointer.DOM.y
                } else if (this.topoEdit.activeEdge.to === '') {
                    this.topoEdit.activeEdge.to = data['nodes'][0]
                    this.updateNetwork()
                    this.slider.properties = {
                        fromObjId: this.topoEdit.activeEdge.from,
                        toObjId: this.topoEdit.activeEdge.to,
                        topoModelList: this.localTopoModelList,
                        edges: this.topoEdit.edges
                    }
                    this.showSlider('theRelation')
                }
            },
            popupEdgeTooltips (data) {
                const edgeId = data.edge
                this.topoTooltip.hoverEdge = this.network.edges.find(edge => edge.id === edgeId)
                if (this.topoTooltip.hoverEdge.labelList.length > 1) {
                    this.$nextTick(() => {
                        const view = this.networkInstance.getViewPosition()
                        const scale = this.networkInstance.getScale()
                        const nodes = this.networkInstance.getConnectedNodes(edgeId)
                        const nodePositions = this.networkInstance.getPositions(nodes)
                        const edgeLeft = (nodePositions[nodes[0]].x + nodePositions[nodes[1]].x) / 2
                        const edgeTop = (nodePositions[nodes[0]].y + nodePositions[nodes[1]].y) / 2
                        const containerBox = this.$refs.topo.getBoundingClientRect()
                        let left = containerBox.width / 2 + (edgeLeft - view.x) * scale + 18
                        if (this.topoEdit.isEdit) {
                            left += NAV_WIDTH
                        }
                        const top = containerBox.height / 2 + (edgeTop - view.y) * scale - 18
                        this.$refs.edgeTooltips.style.left = left + 'px'
                        this.$refs.edgeTooltips.style.top = top + 'px'
                    })
                }
            },
            popupNodeTooltips (data) {
                const nodeId = data.node
                this.topoTooltip.hoverNode = this.network.nodes.find(node => node.id === nodeId)
                this.$nextTick(() => {
                    const view = this.networkInstance.getViewPosition()
                    const scale = this.networkInstance.getScale()
                    const nodeBox = this.networkInstance.getBoundingBox(nodeId)
                    const containerBox = this.$refs.topo.getBoundingClientRect()
                    const left = containerBox.width / 2 + (nodeBox.right - view.x - 18) * scale + NAV_WIDTH
                    const top = containerBox.height / 2 + (nodeBox.top - view.y) * scale
                    this.$refs.nodeTooltips.style.left = left + 'px'
                    this.$refs.nodeTooltips.style.top = top + 'px'
                })
            },
            clearHoverTooltip () {
                this.topoTooltip.hoverNode = null
                this.topoTooltip.hoverNodeTimer = null
            },
            handleHoverEdge (data) {
                this.$refs.topo.style.cursor = 'pointer'
                clearTimeout(this.topoTooltip.hoverEdgeTimer)
                this.popupEdgeTooltips(data)
            },
            handleHoverNode (data) {
                this.$refs.topo.style.cursor = 'pointer'
                clearTimeout(this.topoTooltip.hoverNodeTimer)
                this.popupNodeTooltips(data)
            },
            handleBlurEdge (data) {
                this.$refs.topo.style.cursor = 'default'
                this.topoTooltip.hoverEdgeTimer = setTimeout(() => {
                    this.topoTooltip.hoverEdge = null
                }, 300)
            },
            handleBlurNode (data) {
                this.$refs.topo.style.cursor = 'default'
                this.topoTooltip.hoverNodeTimer = setTimeout(() => {
                    this.topoTooltip.hoverNode = null
                }, 300)
            },
            handleEdgeTooltipsOver () {
                clearTimeout(this.topoTooltip.hoverEdgeTimer)
            },
            handleEdgeTooltipsLeave () {
                this.topoTooltip.hoverEdgeTimer = setTimeout(() => {
                    this.topoTooltip.hoverEdgeTimer = null
                }, 300)
            },
            handleNodeTooltipsOver () {
                clearTimeout(this.topoTooltip.hoverNodeTimer)
            },
            handleNodeTooltipsLeave () {
                this.topoTooltip.hoverNodeTimer = setTimeout(() => {
                    this.topoTooltip.hoverNodeTimer = null
                }, 300)
            },
            handleShowDetails (labelInfo) {
                this.slider.title = labelInfo.text
                this.slider.properties = {
                    objId: labelInfo.objId,
                    isEdit: this.topoEdit.isEdit,
                    asstId: labelInfo.asst['bk_inst_id'],
                    asstInfo: labelInfo.asst.asstInfo || {}
                }
                this.showSlider('theRelationDetail')
            },
            getAssociationName (asstId) {
                let asst = this.associationList.find(asst => asst.id === asstId)
                if (asst) {
                    if (asst['bk_asst_name'].length) {
                        return asst['bk_asst_name']
                    }
                    return asst['bk_asst_id']
                }
            },
            getAssociationType () {
                return this.searchAssociationType({
                    params: {},
                    config: {
                        requestId: 'searchAssociationType'
                    }
                }).then(res => {
                    this.associationList = res.info
                })
            },
            showSlider (content) {
                let {
                    slider
                } = this
                slider.content = content
                switch (content) {
                    case 'theDisplay':
                        slider.properties = {
                            topoModelList: this.localTopoModelList,
                            associationList: this.associationList,
                            isShowModelName: this.displayConfig.isShowModelName,
                            isShowModelAsst: this.displayConfig.isShowModelAsst
                        }
                        slider.width = 600
                        break
                    case 'theRelation':
                    case 'theRelationDetail':
                    default:
                        slider.width = 514
                }
                slider.isShow = true
            },
            toggleGroup (group) {
                if (group['bk_classification_id'] !== this.topoNav.activeGroup) {
                    this.topoNav.activeGroup = group['bk_classification_id']
                } else {
                    this.topoNav.activeGroup = ''
                }
            },
            toggleExample () {
                this.isShowExample = !this.isShowExample
            },
            resizeFull () {
                this.networkInstance.moveTo({scale: 1})
                this.networkInstance.fit()
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
            updateNetwork () {
                this.setNodes(this.localTopoModelList)
                this.setEdges(this.localTopoModelList)
                let scale = this.networkInstance.getScale()
                let origin = this.networkInstance.getViewPosition()
                this.networkInstance = new Vis.Network(this.$refs.topo, {
                    nodes: this.networkDataSet.nodes,
                    edges: this.networkDataSet.edges
                }, this.network.options)
                this.addListener()
                this.networkInstance.moveTo({
                    position: origin,
                    scale
                })
            },
            async initNetwork () {
                await this.getAssociationType()
                const response = await this.$store.dispatch('globalModels/searchModelAction')
                this.localTopoModelList = response
                this.localTopoModelList.forEach(model => {
                    this.$set(model, 'draged', false)
                    if (model.hasOwnProperty('assts') && model.assts.length) {
                        model.assts.forEach(asst => {
                            this.$set(asst, 'checked', true)
                        })
                    }
                })
                this.topoModelList = this.$tools.clone(response)
                this.setNodes(this.localTopoModelList)
                this.setEdges(this.localTopoModelList)
                this.networkInstance = new Vis.Network(this.$refs.topo, {
                    nodes: this.networkDataSet.nodes,
                    edges: this.networkDataSet.edges
                }, this.network.options)
                this.initPosition()
                this.addListener()
            },
            // 设置节点数据
            setNodes (data) {
                let nodes = []
                let asstList = []
                data.forEach(nodeData => {
                    if (nodeData.hasOwnProperty('assts')) {
                        asstList = [...asstList, ...nodeData.assts]
                    }
                })
                data.forEach(nodeData => {
                    if ((nodeData.hasOwnProperty('assts') || asstList.findIndex(({bk_obj_id: objId}) => objId === nodeData['bk_obj_id']) > -1) || nodeData.draged) {
                        const node = {
                            id: nodeData['bk_obj_id'],
                            image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON({
                                name: nodeData['node_name'],
                                backgroundColor: '#fff'
                            }))}`,
                            data: nodeData
                        }
                        if (this.displayConfig.isShowModelName) {
                            node.label = nodeData['node_name']
                        }
                        if (nodeData['position']['x'] !== null && nodeData['position']['y'] !== null) {
                            node.physics = false
                            node.x = nodeData['position']['x']
                            node.y = nodeData['position']['y']
                        }
                        nodes.push(node)
                    }
                })
                this.network.nodes = nodes
                this.networkDataSet.nodes = new Vis.DataSet(this.network.nodes)
            },
            // 设置连线数据
            setEdges (data) {
                let edges = []
                data.forEach(node => {
                    if (Array.isArray(node.assts) && node.assts.length) {
                        node.assts.forEach(asst => {
                            if (this.displayConfig.isShowModelAsst && asst.checked) {
                                const twoWayAsst = this.getTwoWayAsst(node, asst, edges)
                                // 存在则不重复添加
                                let edge = edges.find(edge => edge.to === asst['bk_obj_id'] && edge.from === node['bk_obj_id'])
                                if (edge) {
                                    edge.labelList.push({
                                        text: this.getAssociationName(asst['bk_asst_inst_id']),
                                        arrows: 'to',
                                        objId: node['bk_obj_id'],
                                        asst
                                    })
                                    edge.label = String(edge.labelList.length)
                                } else if (twoWayAsst) { // 双向关联，将已存在的线改为双向
                                    twoWayAsst.arrows = 'to,from'
                                    twoWayAsst.labelList.push({
                                        text: this.getAssociationName(asst['bk_asst_inst_id']),
                                        arrows: 'from',
                                        objId: node['bk_obj_id'],
                                        asst
                                    })
                                    twoWayAsst.label = String(twoWayAsst.labelList.length)
                                } else {
                                    edges.push({
                                        from: node['bk_obj_id'],
                                        to: asst['bk_obj_id'],
                                        arrows: 'to',
                                        label: this.getAssociationName(asst['bk_asst_inst_id']),
                                        labelList: [{
                                            text: this.getAssociationName(asst['bk_asst_inst_id']),
                                            arrows: 'to',
                                            objId: node['bk_obj_id'],
                                            asst
                                        }]
                                    })
                                }
                            }
                        })
                    }
                })
                this.network.edges = edges
                this.networkDataSet.edges = new Vis.DataSet(this.network.edges)
            },
            // 设置单节点位置并更新其节点位置信息
            setSingleNodePosition () {
                const edges = this.network.edges
                const noPositionSingleNode = this.network.nodes.filter(node => !edges.some(edge => edge.from === node.id || edge.to === node.id))
                const fixedNodeY = this.getSingleNodePositionY()
                const fixedDistance = 70
                const compensateDistance = noPositionSingleNode.length % 2 === 0 ? fixedDistance / 2 : 0
                const middleIndex = Math.floor(noPositionSingleNode.length / 2)
                noPositionSingleNode.forEach((node, index) => {
                    node.x = (index - middleIndex) * fixedDistance - compensateDistance
                    node.y = fixedNodeY
                })
                if (noPositionSingleNode.length) {
                    this.networkDataSet.nodes.update(noPositionSingleNode)
                }
            },
            // 获取单节点摆放的初始化Y轴坐标
            getSingleNodePositionY () {
                const asstNodes = this.network.nodes.filter(node => this.network.edges.some(edge => edge.to === node.id || edge.from === node.id))
                if (asstNodes.length) {
                    const asstNodePositions = this.networkInstance.getPositions(asstNodes.map(node => node.id))
                    return Math.min(...asstNodes.map(node => asstNodePositions[node.id]['y'])) - 70
                }
                return 0
            },
            // 判断是否是双向关联A.to = B.from && A.from = B.to
            getTwoWayAsst (node, asst, edges) {
                return edges.find(edge => edge.from === asst['bk_obj_id'] && edge.to === node['bk_obj_id'])
            },
            // 加载节点icon并更新
            loadNodeImage () {
                this.network.nodes.forEach(node => {
                    let image = new Image()
                    image.onload = () => {
                        this.networkDataSet.nodes.update({
                            id: node.id,
                            image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: node.data['node_name'],
                                iconColor: node.data.ispre ? '#868b97' : '#3c96ff',
                                backgroundColor: '#fff'
                            }))}`
                        })
                    }
                    image.src = `${window.location.origin}/static/svg/${node['data']['bk_obj_icon'].substr(5)}.svg`
                })
            },
            initPosition () {
                let nodesId = []
                this.topoModelList.forEach(model => {
                    if (model.hasOwnProperty('assts') && model.assts.length) {
                        model.assts.forEach(asst => {
                            nodesId.push(asst['bk_obj_id'])
                        })
                    }
                })
                nodesId = nodesId.filter(id => {
                    return this.topoModelList.some(({bk_obj_id: objId, position}) => objId === id && position.x === null && position.y === null)
                })
                if (nodesId.length) {
                    this.updateNodePosition(this.networkDataSet.nodes.get(nodesId))
                }
            },
            // 批量更新节点位置信息
            async updateNodePosition (updateNodes) {
                if (!updateNodes.length) return
                const nodePositions = this.networkInstance.getPositions(updateNodes.map(node => node.id))
                const params = updateNodes.map(node => {
                    const nodeData = node.data
                    return {
                        'bk_obj_id': node.id,
                        'bk_inst_id': nodeData['bk_inst_id'],
                        'node_type': nodeData['node_type'],
                        'position': {
                            x: nodePositions[node.id]['x'],
                            y: nodePositions[node.id]['y']
                        }
                    }
                })
                await this.$store.dispatch('globalModels/updateModelAction', {params})
                updateNodes.forEach(node => {
                    let model = this.localTopoModelList.find(({bk_obj_id: objId}) => objId === node.id)
                    model.position.x = nodePositions[node.id]['x']
                    model.position.y = nodePositions[node.id]['y']
                })
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
                    this.networkInstance.unselectAll()
                })
                // this.setSingleNodePosition()
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
                networkInstance.on('dragStart', data => {
                    this.topoTooltip.hoverNode = null
                    this.topoEdit.line.center = data.event.center
                    this.topoEdit.line.dragStart.x1 = this.topoEdit.line.x1
                    this.topoEdit.line.dragStart.y1 = this.topoEdit.line.y1
                    this.topoEdit.line.dragStart.x2 = this.topoEdit.line.x2
                    this.topoEdit.line.dragStart.y2 = this.topoEdit.line.y2
                })
                networkInstance.on('hoverEdge', data => {
                    this.handleHoverEdge(data)
                })
                networkInstance.on('hoverNode', data => {
                    if (this.topoEdit.isEdit) {
                        this.handleHoverNode(data)
                    }
                })
                networkInstance.on('blurEdge', data => {
                    this.handleBlurEdge(data)
                })
                networkInstance.on('blurNode', data => {
                    if (this.topoEdit.isEdit) {
                        this.handleBlurNode(data)
                    }
                })
                networkInstance.on('click', data => {
                    if (data['edges'].length === 1 && data['nodes'].length === 0) {
                        this.handleEdgeClick(data['edges'][0])
                    }
                    if (data['nodes'].length === 1) {
                        this.handleNodeClick(data)
                    } else {
                        this.topoEdit.activeEdge = {
                            from: '',
                            to: ''
                        }
                    }
                })
                networkInstance.on('zoom', data => {
                    this.clearActiveEdge()
                })
                networkInstance.on('dragging', data => {
                    if (this.topoEdit.activeEdge.from) {
                        let {
                            line
                        } = this.topoEdit
                        let offsetX = (line.center.x - data.event.center.x)
                        let offsetY = (line.center.y - data.event.center.y)
                        line.x1 = line.dragStart.x1 - offsetX
                        line.y1 = line.dragStart.y1 - offsetY
                        line.x2 = line.dragStart.x2 - offsetX
                        line.y2 = line.dragStart.y2 - offsetY
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topo-wrapper {
        position: relative;
        padding: 0;
        height: 100%;
        &.has-nav {
            .edit-button {
                display: none;
            }
            .topo-nav {
                display: block;
            }
            .global-model {
                float: left;
                width: calc(100% - 200px);
            }
        }
    }
    .toolbar {
        .edit-button {
            position: absolute;
            padding: 0 10px;
            border-radius: 18px;
            z-index: 1;
            top: 10px;
            left: 20px;
        }
        .vis-button-group {
            position: absolute;
            top: 10px;
            right: 20px;
            z-index: 1;
            font-size: 0;
        }
        .vis-button {
            margin-left: 10px;
            width: 36px;
            height: 36px;
            line-height: 36px;
            padding: 0;
            cursor: pointer;
            border-radius: 50%;
            box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.2);
            border: none;
            text-align: center;
            z-index: 1;
            &.vis-example {
                width: auto;
                padding: 0 15px;
                border-radius: 18px;
                font-size: 0;
                .vis-button-text {
                    font-size: 14px;
                    vertical-align: middle;
                }
                .icon-angle-down {
                    font-size: 12px;
                    vertical-align: middle;
                    transition: all .2s;
                    &.rotate {
                        transform: rotate(180deg);
                    }
                }
            }
        }
        .topo-example {
            position: absolute;
            padding: 3px 10px;
            top: 46px;
            right: 0;
            width: 100px;
            height: 66px;
            background: #fff;
            box-shadow: 0px 2px 1px 0px rgba(185, 203, 222, 0.5);
            font-size: 12px;
            z-index: 1;
            &:before {
                position: absolute;
                top: -10px;
                right: 18px;
                content: "";
                border: 5px solid transparent;
                border-bottom-color: #fff;
            }
            .example-item {
                line-height: 30px;
                font-size: 0;
                &:first-child i{
                    background: $cmdbBorderFocusColor;
                }
                &:last-child i{
                    background: #868b97;
                }
                i {
                    display: inline-block;
                    margin-right: 6px;
                    width: 12px;
                    height: 12px;
                    border-radius: 2px;
                    vertical-align: middle;
                }
                span {
                    font-size: 12px;
                    vertical-align: middle;
                }
            }
        }
    }
    .topo-save-title {
        position: absolute;
        padding: 11px;
        top: 0;
        left: 0;
        width: 100%;
        height: 58px;
        background: #fff;
        font-size: 0;
        z-index: 1;
        .bk-button {
            margin-right: 10px;
        }
    }
    .topo-nav {
        display: none;
        float: left;
        border: 1px solid $cmdbTableBorderColor;
        border-left: none;
        width: 200px;
        height: calc(100% - 58px);
        margin-top: 58px;
        overflow: auto;
        @include scrollbar;
        .group-info {
            line-height: 42px;
            padding: 0 20px 0 15px;
            font-size: 14px;
            cursor: pointer;
            &:hover,
            &.active {
                background: $cmdbBorderFocusColor;
                color: #fff;
                .model-count {
                    background: #fff;
                }
                .icon-angle-down {
                    color: #fff;
                }
            }
            &.active {
                opacity: .65;
                .icon-angle-down {
                    transform: rotate(180deg);
                }
            }
            .model-count {
                padding: 0 5px;
                border-radius: 4px;
                font-size: 12px;
                color: $cmdbBorderFocusColor;
                background: #ebf4ff;
            }
            .icon-angle-down {
                transition: all .2s;
                float: right;
                margin-top: 15px;
                font-size: 12px;
                color: $cmdbBorderColor;
            }
        }
        .model-box {
            padding: 5px 0;
        }
        .model-item {
            padding: 7px 12px;
            cursor: move;
            &:hover {
                background: #ebf4ff;
            }
            &.disabled {
                cursor: not-allowed;
                opacity: .6;
            }
            i {
                display: inline-block;
                margin-right: 5px;
                padding-top: 7px;
                width: 36px;
                height: 36px;
                font-size: 20px;
                line-height: 1;
                text-align: center;
                vertical-align: middle;
                color: $cmdbBorderFocusColor;
                border: 1px solid $cmdbTableBorderColor;
                border-radius: 50%;
            }
            .info {
                display: inline-block;
                line-height: 18px;
                vertical-align: middle;
                font-size: 12px;
                .id {
                    color: $cmdbBorderColor;
                }
            }
        }
    }
    .global-model {
        width: 100%;
        height: 100%;
        background-color: #f4f5f8;
        background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
        background-size: 10px 10px;
    }
    .topo-line {
        position: absolute;
        top: 0;
        left: 200px;
        width: calc(100% - 200px);
        height: 100%;
        z-index: 9;
        pointer-events: none;
    }
    .topology-edge-tooltips {
        position: absolute;
        padding: 5px;
        top: 0;
        left: 0;
        min-width: 100px;
        font-size: 12px;
        color: #868b97;
        background: #fff;
        box-shadow:0px 2px 4px 0px rgba(147,147,147,0.5);
        border-radius:2px;
        cursor: pointer;
        :hover {
            color: $cmdbBorderFocusColor;
            background: #ebf4ff;
        }
        &:after {
            position: absolute;
            content: '';
            border: 6px solid transparent;
            border-right-color: #fff;
            top: 16px;
            left: -12px;
            z-index: 1;
        }
        &:before {
            position: absolute;
            content: '';
            border: 6px solid transparent;
            border-right-color: $cmdbTableBorderColor;
            top: 16px;
            left: -13px;
            z-index: 1;
        }
    }
    .topology-node-tooltips {
        position: absolute;
        padding-top: 1px;
        top: 0;
        left: 0;
        display: inline-block;
        width: 14px;
        height: 14px;
        text-align: center;
        border-radius: 50%;
        font-size: 12px;
        color: #fff;
        background: $cmdbDangerColor;
        .bk-icon {
            transform: scale(.5);
            font-weight: bold;
            vertical-align: top;
            cursor: pointer;
        }
    }
</style>