<template>
    <div class="topo-wrapper" :class="{'has-nav': topoEdit.isEdit}">
        <div class="toolbar">
            <template v-if="!topoEdit.isEdit">
                <bk-button class="edit-button" type="primary"
                    :disabled="!authority.includes('update')"
                    @click="editTopo">
                    {{$t('ModelManagement["编辑拓扑"]')}}
                </bk-button>
            </template>
            <template v-else>
                <bk-button type="primary" @click="exitEdit">
                    {{$t('Common["返回"]')}}
                </bk-button>
                <p class="edit-cue">{{$t('ModelManagement["所有更改已自动保存"]')}}</p>
            </template>
            <div class="vis-button-group">
                <i class="bk-icon icon-full-screen" @click="resizeFull" v-tooltip="$t('ModelManagement[\'还原\']')"></i>
                <i class="bk-icon icon-plus" @click="zoomIn" v-tooltip="$t('ModelManagement[\'放大\']')"></i>
                <i class="bk-icon icon-minus" @click="zoomOut" v-tooltip="$t('ModelManagement[\'缩小\']')"></i>
                <i class="icon-cc-setting"
                    v-if="authority.includes('update')"
                    v-tooltip="$t('ModelManagement[\'拓扑显示设置\']')"
                    @click="showSlider('theDisplay')">
                </i>
                <div class="topo-example">
                    <p class="example-item">
                        <i></i>
                        <span>{{$t('ModelManagement["自定义模型"]')}}</span>
                    </p>
                    <p class="example-item">
                        <i></i>
                        <span>{{$t('ModelManagement["内置模型"]')}}</span>
                    </p>
                </div>
            </div>
        </div>
        <template v-if="topoEdit.isEdit">
            <ul class="topo-nav">
                <li class="group-item" v-for="(group, groupIndex) in localClassifications" :key="groupIndex">
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
            <defs>
                <marker id="arrow" viewBox="0 0 10 10"
                    refX="1" refY="5" 
                    markerUnits="strokeWidth"
                    markerWidth="5"
                    markerHeight="5"
                    orient="auto">
                    <path d="M 0 0 L 10 5 L 0 10 z" fill="#ffb23a"/>
                </marker>
            </defs>
            <line :x1="topoEdit.line.x1" :y1="topoEdit.line.y1" :x2="topoEdit.line.x2" :y2="topoEdit.line.y2" stroke="#ffb23a" stroke-width="2" marker-end="url(#arrow)" stroke-dasharray="5,2"></line>
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
            @mouseleave="handleNodeTooltipsLeave">
            <span class="icon-box is-line" @click="addEdge">
                <i class="icon-cc-line"></i>
            </span>
            <span class="icon-box is-del" @click="deleteNode">
                <i class="icon-cc-del"></i>
            </span>
        </div>
    </div>
</template>

<script>
    import Vis from 'vis'
    import theDisplay from './children/display-config'
    import theRelation from './children/create-relation'
    import theRelationDetail from './children/relation-detail'
    import { generateObjIcon as GET_OBJ_ICON } from '@/utils/util'
    import { mapGetters, mapActions } from 'vuex'
    import throttle from 'lodash.throttle'
    const NAV_WIDTH = 200
    const TOOLBAR_HEIHGT = 50
    export default {
        components: {
            theDisplay,
            theRelation,
            theRelationDetail
        },
        data () {
            return {
                specialModel: ['process', 'plat'],
                associationList: [],
                slider: {
                    width: 514,
                    isShow: false,
                    content: '',
                    properties: {},
                    title: this.$t('ModelManagement["拓扑显示设置"]')
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
            ...mapGetters('globalModels', ['modelConfig']),
            noPositionNodes () {
                return this.network.nodes.filter(node => {
                    const position = node.data.position
                    return position.x === null && position.y === null
                })
            },
            localClassifications () {
                return this.$tools.clone(this.classifications).map(classify => {
                    classify['bk_objects'] = classify['bk_objects'].filter(model => {
                        return !this.isModelInTopo(model) &&
                            !this.specialModel.includes(model['bk_obj_id']) &&
                            !model.bk_ispaused
                    })
                    return classify
                })
            },
            authority () {
                return this.$store.getters.admin ? ['search', 'update', 'delete'] : []
            },
            displayConfig () {
                return {
                    isShowModelName: this.modelConfig.isShowModelName === undefined ? true : this.modelConfig.isShowModelName,
                    isShowModelAsst: this.modelConfig.isShowModelAsst === undefined ? true : this.modelConfig.isShowModelAsst
                }
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
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["模型拓扑"]'))
        },
        mounted () {
            this.initNetwork()
            this.initMoveFunction()
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType'
            ]),
            addEdge () {
                if (this.topoEdit.activeEdge.from === '') {
                    const nodeId = this.topoTooltip.hoverNode.id
                    const view = this.networkInstance.getViewPosition()
                    const positions = this.networkInstance.getPositions([nodeId])
                    const containerBox = this.$refs.topo.getBoundingClientRect()
                    const scale = this.networkInstance.getScale()
                    this.topoEdit.activeEdge.from = nodeId
                    this.topoEdit.line.x1 = (containerBox.left + containerBox.right) / 2 - (view.x - positions[nodeId].x) * scale - containerBox.x
                    this.topoEdit.line.y1 = (containerBox.top + containerBox.bottom) / 2 - (view.y - positions[nodeId].y) * scale - containerBox.y + TOOLBAR_HEIHGT
                }
            },
            isModelInTopo (model) {
                return this.network.nodes.findIndex(node => node.id === model['bk_obj_id']) > -1
            },
            clearEditData () {
                this.topoEdit.isEdit = false
                this.topoEdit.activeEdge.from = ''
                this.topoEdit.activeEdge.to = ''
            },
            async exitEdit () {
                this.topoEdit.isEdit = false
                // this.updateNetwork()
                this.networkInstance.setOptions({nodes: {fixed: true}})
                await this.$nextTick()
                this.networkInstance.redraw()
                this.clearEditData()
            },
            handleDisplaySave (displayConfig) {
                this.displayConfig.isShowModelName = displayConfig.isShowModelName
                this.displayConfig.isShowModelAsst = displayConfig.isShowModelAsst
                this.localTopoModelList = displayConfig.topoModelList
                this.updateNetwork()
            },
            handleRelationSave (params) {
                let fromNode = this.localTopoModelList.find(model => model['bk_obj_id'] === params['bk_obj_id'])
                if (!fromNode.hasOwnProperty('assts')) {
                    Object.assign(fromNode, {assts: []})
                }
                fromNode.assts.push({
                    bk_asst_inst_id: this.associationList.find(asst => asst['bk_asst_id'] === params['bk_asst_id']).id,
                    bk_obj_id: params['bk_asst_obj_id'],
                    bk_inst_id: params.id,
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
                }
                this.updateNetwork()
            },
            editTopo () {
                this.networkInstance.setOptions({
                    nodes: {fixed: false}
                })
                this.topoEdit.isEdit = true
                this.$nextTick(() => {
                    this.networkInstance.redraw()
                })
            },
            checkNodeAsst (node) {
                let asstNum = 0
                this.localTopoModelList.forEach(model => {
                    if (model.hasOwnProperty('assts') && model.assts.length) {
                        if (model['bk_obj_id'] === node.id) {
                            asstNum += model.assts.length
                        } else {
                            model.assts.forEach(asst => {
                                if (asst['bk_obj_id'] === node.id) {
                                    asstNum++
                                }
                            })
                        }
                    }
                })
                if (asstNum) {
                    this.$bkInfo({
                        title: this.$t('ModelManagement["移除失败"]'),
                        content: this.$tc('ModelManagement["移除失败提示"]', asstNum, {asstNum})
                    })
                }
                return !!asstNum
            },
            deleteNode () {
                let {
                    hoverNode
                } = this.topoTooltip
                if (this.checkNodeAsst(hoverNode)) {
                    return
                }
                this.$bkInfo({
                    title: this.$t('ModelManagement["确定移除模型?"]'),
                    content: this.$t('ModelManagement["移除模型提示"]'),
                    confirmFn: () => {
                        let node = this.localTopoModelList.find(model => model['bk_obj_id'] === hoverNode.id)
                        node.position = {x: null, y: null}

                        this.updateSingleNodePosition({
                            bk_obj_id: node['bk_obj_id'],
                            bk_inst_id: node['bk_inst_id'],
                            node_type: node['node_type'],
                            position: {
                                x: node.position.x,
                                y: node.position.y
                            }
                        })
                        
                        this.topoTooltip.hoverNode = null
                        this.topoTooltip.hoverNodeTimer = null
                        this.updateNetwork()
                    },
                    cancelFn: () => {
                        this.topoTooltip.hoverNode = null
                        this.topoTooltip.hoverNodeTimer = null
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
                node.position.x = Math.floor(originPosition.x - ((container.left + container.right) / 2 - event.clientX) / scale)
                node.position.y = Math.floor(originPosition.y - ((container.top + container.bottom) / 2 - event.clientY) / scale)
                this.updateNetwork()
                this.updateSingleNodePosition({
                    bk_obj_id: node['bk_obj_id'],
                    bk_inst_id: node['bk_inst_id'],
                    node_type: node['node_type'],
                    position: {
                        x: node.position.x,
                        y: node.position.y
                    }
                })
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
            handleEdgeClick (edgeId) {
                let edge = this.network.edges.find(({id}) => id === edgeId)
                if (edge.labelList.length === 1) {
                    this.handleShowDetails(edge.labelList[0])
                }
            },
            initMoveFunction () {
                this.handleMouseMove = throttle(event => {
                    this.topoEdit.line.x2 = event.layerX
                    this.topoEdit.line.y2 = event.layerY + TOOLBAR_HEIHGT
                }, 50)
            },
            handleNodeClick (data) {
                if (!this.topoEdit.isEdit) {
                    return
                }
                if (this.topoEdit.activeEdge.from && this.topoEdit.activeEdge.to === '') {
                    this.topoEdit.activeEdge.to = data['nodes'][0]
                    this.slider.properties = {
                        fromObjId: this.topoEdit.activeEdge.from,
                        toObjId: this.topoEdit.activeEdge.to,
                        topoModelList: this.localTopoModelList
                    }
                    this.slider.title = this.$t('ModelManagement["新建关联"]')
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
                        const top = containerBox.height / 2 + (edgeTop - view.y) * scale - 18 + TOOLBAR_HEIHGT
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
                    const left = containerBox.width / 2 + (nodeBox.right - view.x) * scale + NAV_WIDTH - 8
                    const top = containerBox.height / 2 + (nodeBox.top - view.y) * scale - 8 + TOOLBAR_HEIHGT
                    this.$refs.nodeTooltips.style.left = left + 'px'
                    this.$refs.nodeTooltips.style.top = top + 'px'
                })
            },
            clearHoverTooltip () {
                this.topoTooltip.hoverNode = null
                this.topoTooltip.hoverNodeTimer = null
                this.topoTooltip.hoverEdge = null
                this.topoTooltip.hoverEdgeTimer = null
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
                        this.slider.title = this.$t('ModelManagement["拓扑显示设置"]')
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
                this.networkInstance.setOptions({nodes: {fixed: true}})
                if (this.authority.includes('update')) {
                    this.initPosition()
                }
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
                    if (((nodeData.hasOwnProperty('assts') && nodeData.assts.length) || asstList.findIndex(({bk_obj_id: objId}) => objId === nodeData['bk_obj_id']) > -1) || (nodeData.position.x !== null && nodeData.position.y !== null)) {
                        const node = {
                            id: nodeData['bk_obj_id'],
                            image: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON({
                                name: nodeData['node_name'],
                                backgroundColor: '#fff'
                            }))}`,
                            data: nodeData,
                            hidden: !this.modelConfig[nodeData.bk_obj_id]
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
                            if (asst.checked) {
                                const twoWayAsst = this.getTwoWayAsst(node, asst, edges)
                                // 存在则不重复添加
                                let edge = edges.find(edge => edge.to === asst['bk_obj_id'] && edge.from === node['bk_obj_id'])
                                if (edge) { // 已存在主动关联
                                    this.updateEdgeArrows(edge, this.getEdgeArrows(asst))
                                    edge.labelList.push({
                                        text: this.getAssociationName(asst['bk_asst_inst_id']),
                                        objId: node['bk_obj_id'],
                                        asst
                                    })
                                    if (this.displayConfig.isShowModelAsst) {
                                        edge.label = String(edge.labelList.length)
                                    }
                                } else if (twoWayAsst) { // 被关联
                                    this.updateEdgeArrows(twoWayAsst, this.getEdgeArrows(asst))
                                    twoWayAsst.labelList.push({
                                        text: this.getAssociationName(asst['bk_asst_inst_id']),
                                        objId: node['bk_obj_id'],
                                        asst
                                    })
                                    if (this.displayConfig.isShowModelAsst) {
                                        twoWayAsst.label = String(twoWayAsst.labelList.length)
                                    }
                                } else { // 无关联关系
                                    const edge = {
                                        from: node['bk_obj_id'],
                                        to: asst['bk_obj_id'],
                                        arrows: this.getEdgeArrows(asst),
                                        labelList: [{
                                            text: this.getAssociationName(asst['bk_asst_inst_id']),
                                            objId: node['bk_obj_id'],
                                            asst
                                        }]
                                    }
                                    if (this.displayConfig.isShowModelAsst) {
                                        edge.label = this.getAssociationName(asst['bk_asst_inst_id'])
                                    }
                                    edges.push(edge)
                                }
                            }
                        })
                    }
                })
                this.network.edges = edges
                this.networkDataSet.edges = new Vis.DataSet(this.network.edges)
            },
            updateEdgeArrows (edge, arrows) {
                if (edge.arrows === '') {
                    edge.arrows = arrows
                } else if (edge.arrows === 'to' && arrows === 'to,from') {
                    edge.arrows = arrows
                }
            },
            getEdgeArrows (asst) {
                const asstType = this.associationList.find(({id}) => id === asst['bk_asst_inst_id'])['direction']
                let arrows = ''
                switch (asstType) {
                    case 'bidirectional':
                        arrows = 'to,from'
                        break
                    case 'src_to_dest':
                        arrows = 'to'
                        break
                    case 'none':
                    default:
                        break
                }
                return arrows
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
                        nodesId.push(model['bk_obj_id'])
                        model.assts.forEach(asst => {
                            nodesId.push(asst['bk_obj_id'])
                        })
                    }
                })
                nodesId = [...new Set(nodesId)]
                nodesId = nodesId.filter(id => {
                    return this.topoModelList.some(({bk_obj_id: objId, position}) => objId === id && position.x === null && position.y === null)
                })
                if (nodesId.length) {
                    this.updateNodePosition(this.networkDataSet.nodes.get(nodesId))
                }
            },
            updateSingleNodePosition (node) {
                this.$store.dispatch('globalModels/updateModelAction', {params: [node]})
            },
            // 批量更新节点位置信息
            async updateNodePosition (updateNodes, removeNodes = []) {
                if (!updateNodes.length && !removeNodes.length) return
                let nodePositions = []
                let params = []
                if (updateNodes.length) {
                    nodePositions = this.networkInstance.getPositions(updateNodes.map(node => node.id))
                    params = updateNodes.map(node => {
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
                }
                if (removeNodes.length) {
                    removeNodes.forEach(node => {
                        params.push({
                            'bk_obj_id': node['bk_obj_id'],
                            'bk_inst_id': node['bk_inst_id'],
                            'node_type': node['node_type'],
                            'position': {
                                x: null,
                                y: null
                            }
                        })
                    })
                }

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
                this.networkInstance.on('dragEnd', data => {
                    if (this.topoEdit.isEdit && data.nodes.length === 1) {
                        const nodeId = data.nodes[0]
                        const position = this.networkInstance.getPositions([nodeId])
                        const model = this.localTopoModelList.find(({bk_obj_id: objId}) => objId === nodeId)
                        model.position.x = position[nodeId].x
                        model.position.y = position[nodeId].y
                        this.updateSingleNodePosition({
                            bk_obj_id: model['bk_obj_id'],
                            bk_inst_id: model['bk_inst_id'],
                            node_type: model['node_type'],
                            position: {
                                x: model.position.x,
                                y: model.position.y
                            }
                        })
                    }
                    this.networkInstance.unselectAll()
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
                    if (this.topoEdit.isEdit && !this.specialModel.includes(data.node)) {
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
                    if (data['nodes'].length === 1 && !this.specialModel.includes(data['nodes'][0])) {
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
                    this.clearHoverTooltip()
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
        padding: 7px 20px;
        width: 100%;
        height: 50px;
        background: #fff;
        font-size: 0;
        .bk-button {
            margin-right: 10px;
        }
        .edit-cue {
            display: inline-block;
            font-size: 14px;
            color: #a4aab3;
            line-height: 36px;
            vertical-align: middle;
        }
        i {
            font-size: 14px;
        }
        .vis-button-group {
            float: right;
            padding-top: 11px;
            >i {
                margin-left: 32px;
                font-size: 14px;
                font-weight: bold;
                cursor: pointer;
                &:hover {
                    color: $cmdbBorderFocusColor;
                }
            }
        }
        .topo-example {
            position: absolute;
            padding: 3px 10px;
            top: 57px;
            right: 8px;
            background: #fff;
            box-shadow: 0px 2px 1px 0px rgba(185, 203, 222, 0.5);
            font-size: 12px;
            z-index: 1;
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
    .topo-nav {
        display: none;
        float: left;
        border: 1px solid $cmdbTableBorderColor;
        border-left: none;
        width: 200px;
        height: calc(100% - 50px);
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
                color: #868b97;
            }
        }
        .model-item {
            padding: 7px 12px;
            cursor: move;
            &:first-child {
                padding-top: 12px;
            }
            &:last-child {
                padding-bottom: 12px;
            }
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
        height: calc(100% - 50px);
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
        top: 0;
        left: 0;
        font-size: 12px;
        color: #fff;
        .icon-box {
            position: absolute;
            display: inline-block;
            padding-top: 3px;
            width: 20px;
            height: 20px;
            border-radius: 50%;
            background: #181818;
            text-align: center;
            visibility: top;
            line-height: 1;
            transform: scale(.8);
            cursor: pointer;
            &.is-del {
                top: 18px;
                left: 8px;
            }
        }
        .bk-icon {
            transform: scale(.5);
            font-weight: bold;
            vertical-align: top;
            cursor: pointer;
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
