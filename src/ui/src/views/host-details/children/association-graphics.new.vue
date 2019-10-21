<template>
    <div class="relation-topology-layout" :class="{ 'full-screen': fullScreen }" ref="layout">
        <bk-button class="exit-full-screen icon-cc-resize-small" size="small" theme="default"
            v-show="fullScreen"
            @click="toggleFullScreen(false)">
            {{$t('退出')}}
        </bk-button>
        <div class="tolology-loading" v-bkloading="{ isLoading: $loading(getRelationRequestId) }">
            <div class="topology-container" ref="container">
            </div>
        </div>
        <ul class="topology-legend">
            <li class="legend-item"
                v-for="(legend, index) in legends[selectedNodeId]"
                :key="index"
                :class="{ inactive: !legend.active }"
                @click="handleToggleNodeVisibility(legend)">
                <i :class="legend.icon"></i>
                {{legend.name}} {{legend.count}}
            </li>
        </ul>
        <div
            class="topology-tooltips"
            ref="tooltips"
            v-show="hoverNodeData"
        >
            <a class="tooltips-option" href="javascript:void(0)"
                @click="handleShowDetails">
                {{$t('详情信息')}}
            </a>
        </div>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :is-show.sync="details.show"
            :title="details.title">
            <cmdb-details slot="content"
                v-if="details.show"
                :show-options="false"
                :inst="details.inst"
                :properties="details.properties"
                :property-groups="details.propertyGroups">
            </cmdb-details>
        </bk-sideslider>
    </div>
</template>

<script>
    import { mapActions, mapGetters, mapState } from 'vuex'
    import cytoscape from 'cytoscape'
    import popper from 'cytoscape-popper'
    import dagre from 'cytoscape-dagre'
    import { generateObjIcon as GET_OBJ_ICON } from '@/utils/util'
    import memoize from 'lodash.memoize'
    import debounce from 'lodash.debounce'
    import throttle from 'lodash.throttle'

    // cytoscape实例
    let cy = null

    let NODE_ID = 0

    export default {
        name: 'cmdb-host-association-graphics',
        data () {
            return {
                getRelationRequestId: null,
                fullScreen: false,
                // 关联数据
                associationList: [],

                // 实例列表
                instanceMap: {},

                layout: {
                    name: 'dagre',
                    rankDir: 'LR',
                    ranker: 'network-simplex',
                    fit: true,
                    padding: 60
                },

                legends: {},

                hoverNodeData: null,

                selectedNodeId: null,

                details: {
                    show: false,
                    title: '',
                    inst: {},
                    properties: [],
                    propertyGroups: []
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', [
                'getModelById'
            ]),
            ...mapState('hostDetails', ['info']),
            host () {
                return this.info.host || {}
            },
            id () {
                return parseInt(this.$route.params.id)
            }
        },
        created () {
            if (typeof cytoscape('core', 'popper') !== 'function') {
                cytoscape.use(popper)
            }
            if (typeof cytoscape('core', 'dagre') !== 'function') {
                cytoscape.use(dagre)
            }
        },
        mounted () {
            this.initNetwork()
        },
        methods: {
            ...mapActions('objectRelation', ['getInstRelationTopo', 'updateInstRelation']),
            ...mapActions('objectAssociation', [
                'searchAssociationType'
            ]),
            initNetwork () {
                cy = window.cy = cytoscape({
                    container: this.$refs.container,

                    minZoom: 0.5,
                    maxZoom: 5,
                    wheelSensitivity: 0.5,
                    pixelRatio: 2,

                    // 元素定义，支持promise
                    elements: this.initTopoElements(),

                    layout: this.layout,

                    style: [
                        {
                            selector: 'core',

                            // grabbed画布时
                            style: {
                                'active-bg-color': '#3c96ff',
                                'active-bg-size': '18px'
                            }
                        },

                        // 有关node样式配置
                        {
                            selector: 'node',
                            style: {
                                // 点击时不显示overlay
                                'overlay-opacity': 0
                            }
                        },
                        {
                            selector: 'node',
                            style: {
                                'width': 36,
                                'height': 36,

                                // 设置label文本
                                'label': 'data(name)',

                                // label
                                'color': '#868b97',
                                'text-valign': 'bottom',
                                'text-halign': 'center',
                                'font-size': '14px',
                                'text-margin-y': '9px',

                                // label换行
                                'text-wrap': 'wrap',
                                'text-max-width': '90px',
                                'text-overflow-wrap': 'anywhere',

                                // 背景图
                                'background-color': '#ffffff',
                                'background-fit': 'cover cover',
                                'border-width': 1,
                                'border-color': '#939393',
                                'border-opacity': 0.5
                            }
                        },
                        {
                            selector: 'node.root',
                            style: {
                                'width': 56,
                                'height': 56
                            }
                        },
                        {
                            selector: 'node.bg',
                            style: {
                                'background-image': 'data(bg.unselected)'
                            }
                        },
                        {
                            selector: 'node.hover, node:selected',
                            style: {
                                'background-image': 'data(bg.selected)',
                                'border-color': '#3a84ff',
                                'font-weight': 'bold'
                            }
                        },

                        // edge样式配置
                        {
                            selector: 'edge',
                            style: {
                                'curve-style': 'straight',
                                'target-arrow-shape': 'triangle-backcurve',
                                'opacity': 1,
                                'arrow-scale': 1.5,
                                'line-color': '#c3cdd7',
                                'target-arrow-color': '#c3cdd7',
                                'width': 2,

                                // 点击时overlay
                                'overlay-padding': '3px',

                                // label
                                'color': '#979ba5',
                                'font-size': '10px',
                                'text-background-opacity': 0.7,
                                'text-background-color': '#ffffff',
                                'text-background-shape': 'roundrectangle',
                                'text-background-padding': 2,
                                'text-border-opacity': 0.7,
                                'text-border-width': 1,
                                'text-border-style': 'solid',
                                'text-border-color': '#dcdee5',

                                'loop-direction': '45deg',
                                'loop-sweep': '90deg'
                            }
                        },
                        {
                            selector: 'edge[direction="none"]', // 无方向
                            style: {
                                'source-arrow-shape': 'none',
                                'target-arrow-shape': 'none'
                            }
                        },
                        {
                            selector: 'edge[direction="bidirectional"]', // 双向
                            style: {
                                'source-arrow-shape': 'triangle-backcurve',
                                'source-arrow-color': '#c3cdd7'
                            }
                        },
                        {
                            selector: 'edge.hover',
                            style: {
                                'width': 3,
                                'line-color': '#3c96ff',
                                'source-arrow-color': '#3c96ff',
                                'target-arrow-color': '#3c96ff',
                                'font-weight': 'bold'
                            }
                        }
                    ]
                })

                // 所有操作的事件绑定
                cy.on('layoutready', (event) => {
                    this.loadNodeImage()
                    event.cy.autolock(true)
                }).on('layoutstop', (event) => {
                    this.fitMaxZoom(event.cy)
                }).on('resize', debounce((event) => {
                    event.cy.fit()
                    this.fitMaxZoom(event.cy)
                }, 500)).on('mouseover', 'node', throttle((event) => {
                    const node = event.target
                    const nodeData = node.data()

                    // 添加hover样式
                    node.addClass('hover')

                    // 显示tooltip
                    this.hoverNodeData = nodeData
                    this.$refs.layout.style.cursor = 'pointer'

                    // 每次重新创建因content引用的内容只能移动一次无法反复使用
                    const popover = this.$bkPopover(node.popperRef(), {
                        content: this.$refs.tooltips,
                        hideOnClick: true,
                        sticky: true,
                        delay: [500, 1000],
                        placement: 'right',
                        interactive: true,
                        animateFill: false,
                        theme: 'node-tooltip',
                        trigger: 'manual',
                        distance: 2
                    })

                    node.data('popover', popover)
                    popover.show()
                }, 160)).on('mouseout', 'node', throttle((event) => {
                    const node = event.target
                    node.removeClass('hover')

                    const popover = node.data('popover')
                    if (popover) {
                        popover.hide()
                    }

                    this.hoverNodeData = null
                    this.$refs.layout.style.cursor = 'default'
                }, 160)).on('click', 'node', (event) => {
                    const node = event.target

                    if (node.data('loaded') !== true) {
                        this.handleSelectNode(node)
                    }

                    this.selectedNodeId = node.id()
                }).on('mouseover', 'edge', (event) => {
                    event.target.addClass('hover')
                }).on('mouseout', 'edge', (event) => {
                    event.target.removeClass('hover')
                })
            },
            async initTopoElements () {
                try {
                    const rootObjId = 'host'
                    const rootInstId = this.id
                    const rootInstName = this.host['bk_host_innerip']
                    const rootNodeId = `${rootObjId}_${rootInstId}_${NODE_ID++}`
                    const [asstData, relData] = await Promise.all([
                        this.getAssociationType(),
                        this.getRelationTopo(rootObjId, rootInstId)
                    ])

                    const topoData = relData.data
                    const { instance } = topoData
                    this.associationList = asstData
                    this.instanceMap = instance
                    this.selectedNodeId = rootNodeId

                    let elements = []

                    // 当前实例作为根节点
                    elements.push({
                        data: {
                            id: rootNodeId,
                            name: rootInstName,
                            icon: this.getModelById(rootObjId).bk_obj_icon,
                            objId: rootObjId,
                            instId: rootInstId,
                            modelName: this.getModelById(rootObjId).bk_obj_name,
                            loaded: true
                        },
                        group: 'nodes',
                        classes: 'root'
                    })

                    // 获取节点下的拓扑元素
                    const initElements = this.getTopoElements(topoData, rootNodeId)

                    // 组装成最终的拓扑
                    elements = [...elements, ...initElements]

                    return elements
                } catch (e) {
                    console.log(e)
                }
            },
            loadNodeImage () {
                // 缓存调用结果，减少相同icon的转换开销
                const makeSvg = memoize(this.makeSvg, data => data.icon)
                cy.nodes().forEach(async (node, i) => {
                    const svg = await makeSvg(node.data())
                    node.data('bg', svg)
                    node.addClass('bg')
                    if (node.id() === this.selectedNodeId) {
                        cy.$(`#${this.selectedNodeId}`).select()
                    }
                })
            },
            makeSvg (nodeData) {
                return new Promise((resolve, reject) => {
                    const image = new Image()
                    image.onload = () => {
                        const svg = {
                            unselected: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: nodeData.name,
                                iconColor: '#798aad',
                                backgroundColor: '#fff'
                            }))}`,
                            selected: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: nodeData.name,
                                iconColor: '#fff',
                                backgroundColor: '#3a84ff'
                            }))}`
                        }

                        resolve(svg)
                    }
                    image.src = `${window.location.origin}/static/svg/${nodeData.icon.substr(5)}.svg`
                })
            },
            getRelationTopo (objId, instId) {
                this.getRelationRequestId = `get_getInstRelationTopo_${objId}_${instId}`
                return this.getInstRelationTopo({
                    objId,
                    instId,
                    params: this.$injectMetadata(),
                    config: {
                        requestId: this.getRelationRequestId,
                        clearCache: true
                    }
                })
            },
            getAssociationType () {
                return this.searchAssociationType({
                    params: {},
                    config: {
                        requestId: 'searchAssociationType'
                    }
                }).then(res => res.info)
            },
            getAsstDetail (asstId) {
                const asst = this.associationList.find(asst => asst.bk_asst_id === asstId)
                return {
                    asstId: asst['bk_asst_id'],
                    asstName: asst['bk_asst_name'].length ? asst['bk_asst_name'] : asst['bk_asst_id'],
                    direction: asst.direction
                }
            },
            getInstDetail (objId, instId) {
                // 需要兼容不同实例的属性名称不一致
                const inst = this.instanceMap[objId].find(inst => inst.bk_inst_id === instId || inst.bk_host_id === instId)
                return {
                    id: instId,
                    name: inst['bk_inst_name'] ? inst['bk_inst_name'] : inst['bk_host_innerip']
                }
            },
            getTopoElements (topoData, rootNodeId) {
                const { association, instance } = topoData

                // 更新实例信息新的实例得以获取到正确的详情
                this.instanceMap = { ...this.instanceMap, ...instance }

                const asstSrc = association.src || []
                const asstDst = association.dst || []
                const elements = []

                // 拓扑图中的所有连线数据，用于查找连线关系
                const edges = cy.edges().map(egde => egde.data())

                // 所有以rootNodeId为目标的关联数据
                asstDst.forEach(item => {
                    const nodeIdPrefix = `${item.bk_obj_id}_${item.bk_inst_id}_`
                    const nodeId = `${nodeIdPrefix}${NODE_ID++}`

                    // 是否存在目标是当前根节点的连接
                    const exist = edges.find(edge => edge.target === rootNodeId)
                    // 不存在或者存在时来源实例不同
                    if (!exist || !exist.source.startsWith(nodeIdPrefix)) {
                        const nodeOptions = this.getNodeOptions({ nodeId, objId: item.bk_obj_id, instId: item.bk_inst_id })
                        const edgeOptions = this.getEdgeOptions({ source: nodeId, target: rootNodeId, asstId: item.bk_asst_id })
                        elements.push(nodeOptions)
                        elements.push(edgeOptions)

                        this.setLegends(rootNodeId, item.bk_obj_id, nodeId)
                    }
                })

                // 所有以rootNodeId为源的关联数据
                asstSrc.forEach(item => {
                    // 为源时，取目标实例id
                    const nodeIdPrefix = `${item.bk_asst_obj_id}_${item.bk_asst_inst_id}_`
                    const nodeId = `${nodeIdPrefix}${NODE_ID++}`
                    const exist = edges.find(item => item.source === rootNodeId)
                    if (!exist || !exist.target.startsWith(nodeIdPrefix)) {
                        const nodeOptions = this.getNodeOptions({ nodeId, objId: item.bk_asst_obj_id, instId: item.bk_asst_inst_id })
                        const edgeOptions = this.getEdgeOptions({ source: rootNodeId, target: nodeId, asstId: item.bk_asst_id })
                        elements.push(nodeOptions)
                        elements.push(edgeOptions)

                        this.setLegends(rootNodeId, item.bk_asst_obj_id, nodeId)
                    }
                })

                return elements
            },
            getNodeOptions ({ nodeId, objId, instId }) {
                const model = this.getModelById(objId)
                const options = {
                    data: {
                        id: nodeId,
                        objId,
                        instId,
                        name: this.getInstDetail(objId, instId).name,
                        icon: model.bk_obj_icon,
                        modelName: model.bk_obj_name
                    },
                    group: 'nodes',
                    classes: ''
                }

                return options
            },
            getEdgeOptions ({ source, target, asstId }) {
                const { direction } = this.getAsstDetail(asstId)
                const options = {
                    data: {
                        source,
                        target,
                        direction
                    },
                    group: 'edges',
                    classes: ''
                }

                return options
            },
            setLegends (rootNodeId, objId, nodeId) {
                const model = this.getModelById(objId)
                let nodelegends = this.legends[rootNodeId]
                const legendNew = {
                    id: objId,
                    name: model.bk_obj_name,
                    icon: model.bk_obj_icon,
                    active: true,
                    count: 1,
                    nodeIds: [nodeId]
                }
                if (nodelegends) {
                    const legend = nodelegends.find(legend => legend.id === objId)
                    if (legend) {
                        legend.count++
                        legend.nodeIds.push(nodeId)
                    } else {
                        nodelegends.push(legendNew)
                    }
                } else {
                    nodelegends = [legendNew]
                }

                this.legends = { ...this.legends, ...{ [rootNodeId]: nodelegends } }
            },
            async handleSelectNode (node) {
                const nodeData = node.data()
                const { objId, instId, id } = nodeData

                // 获取当前节点拓扑数据
                const relData = await this.getRelationTopo(objId, instId)
                const topoData = relData.data

                // 根据拓扑数据获取拓扑元素，根节点为当前节点
                const nodeElements = this.getTopoElements(topoData, id)

                // 将元素加入到拓扑图
                cy.add(nodeElements)

                // 重新应用layout
                cy.autolock(false)
                cy.layout(this.layout).run()

                // 标记为loaded
                node.data('loaded', true)
            },
            handleToggleNodeVisibility (legend) {
                legend.active = !legend.active

                cy.startBatch()
                legend.nodeIds.forEach(nodeId => cy.$(`#${nodeId}`).style('display', legend.active ? 'element' : 'none'))
                cy.endBatch()
            },
            async handleShowDetails () {
                const nodeData = this.hoverNodeData
                this.details.title = `${nodeData.modelName}-${nodeData.name}`
                try {
                    const [inst, properties, propertyGroups] = await Promise.all([
                        this.getInst(),
                        this.getProperties(),
                        this.getPropertyGroups()
                    ])
                    this.details.inst = this.$tools.flattenItem(properties, inst)
                    this.details.properties = properties
                    this.details.propertyGroups = propertyGroups
                    this.details.show = true
                } catch (e) {
                    this.details.inst = {}
                    this.details.properties = []
                    this.details.propertyGroups = []
                    this.details.show = false
                }

                if (this.hoverNodeData.popover) {
                    this.hoverNodeData.popover.hide()
                }
            },
            async getInst () {
                const modelId = this.hoverNodeData.objId
                if (modelId === 'host') {
                    return this.getHostDetails()
                } else if (modelId === 'biz') {
                    return this.getBusinessDetails()
                }
                return this.getInstDetails()
            },
            getHostDetails () {
                const hostId = this.hoverNodeData.instId
                return this.$store.dispatch('hostSearch/getHostBaseInfo', { hostId }).then(data => {
                    const inst = {}
                    data.forEach(field => {
                        inst[field['bk_property_id']] = field['bk_property_value']
                    })
                    return inst
                })
            },
            getBusinessDetails () {
                const bizId = this.hoverNodeData.instId
                return this.$store.dispatch('objectBiz/searchBusiness', {
                    params: {
                        condition: { 'bk_biz_id': bizId },
                        fields: [],
                        page: { start: 0, limit: 1 }
                    }
                }).then(({ info }) => info[0])
            },
            getInstDetails () {
                const modelId = this.hoverNodeData.objId
                const instId = this.hoverNodeData.instId
                const model = this.getModelById(modelId)
                return this.$store.dispatch('objectCommonInst/searchInst', {
                    objId: modelId,
                    params: this.$injectMetadata({
                        condition: {
                            [modelId]: [{
                                field: 'bk_inst_id',
                                operator: '$eq',
                                value: instId
                            }]
                        },
                        fields: {},
                        page: { start: 0, limit: 1 }
                    }, {
                        inject: !!this.$tools.getMetadataBiz(model)
                    })
                }).then(({ info }) => info[0])
            },
            getProperties () {
                const modelId = this.hoverNodeData.objId
                const model = this.getModelById(modelId)
                return this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                    params: this.$injectMetadata({
                        'bk_obj_id': modelId
                    }, {
                        inject: !!this.$tools.getMetadataBiz(model)
                    })
                })
            },
            getPropertyGroups () {
                const modelId = this.hoverNodeData.objId
                const model = this.getModelById(modelId)
                return this.$store.dispatch('objectModelFieldGroup/searchGroup', {
                    objId: modelId,
                    params: this.$injectMetadata({}, {
                        inject: !!this.$tools.getMetadataBiz(model)
                    })
                })
            },
            toggleFullScreen (fullScreen) {
                this.fullScreen = fullScreen
            },
            fitMaxZoom (cy) {
                const fitMaxZoom = 1
                if (cy.zoom() > fitMaxZoom) {
                    cy.zoom(fitMaxZoom)
                    cy.center()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-topology-layout {
        height: 100%;
        background-color: #f9f9f9;
        position: relative;
        &.full-screen {
            position: fixed;
            left: 0;
            right: 0;
            top: 0;
            bottom: 0;
            height: 100%;
            .exit-full-screen {
                position: absolute;
                top: 20px;
                right: 20px;
                z-index: 9999;
            }
        }
        .tolology-loading {
            height: 100%;
        }
        .topology-container {
            height: 100%;
        }
        .topology-legend {
            position: absolute;
            left: 20px;
            top: 20px;
            font-size: 14px;
            background-color: #f9f9f9;
            .legend-item {
                margin: 0 0 5px 0;
                cursor: pointer;
                &.inactive {
                    color: #c3cdd7;
                }
            }
        }
    }
    .topology-tooltips {
        position: relative;
        background-color: #fff;
        font-size: 12px;
        line-height: 22px;
        padding: 2px 8px;
        border-radius: 2px;
        border: 1px solid #b0becc;
        text-align: center;
        &::before {
            content: '';
            background: url('../../assets/images/tip.png');
            width: 6px;
            height: 10px;
            position: absolute;
            display: inline-block;
            left: -6px;
            top: 12px;
            transform: translate(0, -50%);
        }
        .tooltips-option {
            display: block;
            white-space: nowrap;
            &:hover {
                color: #3c96ff;
            }
            &.tooltips-option-delete {
                border-top: 1px solid #dde3e9;
                &:hover {
                    color: #ff5656;
                }
            }
        }
    }
</style>

<style lang="scss">
    .tippy-popper {
        transition: none!important;
    }

    .tippy-tooltip {
        &.node-tooltip-theme {
            background: none;
        }
    }
</style>
