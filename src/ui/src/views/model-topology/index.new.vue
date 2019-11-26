<template>
    <div :class="['topo-wrapper', { hover: isTopoHover }]">
        <div class="toolbar">
            <cmdb-auth style="display: none;" ref="addBusinessLevel" :auth="$authResources({ type: $OPERATION.SYSTEM_TOPOLOGY })" @update-auth="handleReceiveAuth"></cmdb-auth>
            <template v-if="!topoEdit.isEdit">
                <cmdb-auth :auth="$authResources({ type: $OPERATION.SYSTEM_MODEL_GRAPHICS })">
                    <bk-button slot-scope="{ disabled }"
                        class="edit-button"
                        theme="primary"
                        :disabled="disabled"
                        @click="handleEditTopo">
                        {{$t('编辑拓扑')}}
                    </bk-button>
                </cmdb-auth>
            </template>
            <template v-else>
                <bk-button style="margin-top: -2px;" theme="primary" @click="handleExitEdit">
                    {{$t('返回')}}
                </bk-button>
                <p class="edit-cue">{{$t('所有更改已自动保存')}}</p>
            </template>
            <div class="vis-button-group">
                <i
                    :class="['bk-cc-icon', mainFullScreen ? 'icon-cc-fullscreen-outlined-reset' : 'icon-cc-fullscreen-outlined']"
                    @click="resizeFull"
                    v-bk-tooltips="$t(mainFullScreen ? '取消全屏' : '全屏')"
                >
                </i>
                <i class="bk-cc-icon icon-cc-fit" @click="resizeFit" v-bk-tooltips="$t('还原')"></i>
                <i class="bk-cc-icon icon-cc-zoom-out" @click="zoomOut" v-bk-tooltips="$t('缩小')"></i>
                <i class="bk-cc-icon icon-cc-zoom-in" @click="zoomIn" v-bk-tooltips="$t('放大')"></i>
                <div class="topo-legend">
                    <p class="legend-item built-in">
                        <i></i>
                        <span>{{$t('内置模型')}}</span>
                    </p>
                    <p class="legend-item custom">
                        <i></i>
                        <span>{{$t('自定义模型')}}</span>
                    </p>
                </div>
            </div>
        </div>

        <ul class="topo-nav">
            <li class="group-item">
                <div :class="['group-info', 'group-total', { 'selected': topoNav.selectedGroupId === -1 }]" @click="handleSelectGroup()">
                    <span class="group-name">全部模型</span>
                    <span class="model-count">{{localTopoModelList.length > 1000 ? '999+' : localTopoModelList.length}}</span>
                </div>
            </li>
            <li class="group-item" v-for="(group, groupIndex) in localClassifications" :key="groupIndex">
                <div
                    class="group-info"
                    :class="{
                        'active': topoNav.activeGroupId === group['bk_classification_id'],
                        'selected': topoNav.selectedGroupId === group['bk_classification_id'],
                        'invisible': topoNav.hideGroupIds.includes(group['bk_classification_id'])
                    }"
                    @click="handleSelectGroup(group)"
                >
                    <span class="toggle-arrow" @click.stop="handleSlideGroup(group)"><i class="bk-icon icon-angle-down"></i></span>
                    <span class="group-name">{{group['bk_classification_name']}}</span>
                    <span class="model-count">{{group['bk_objects'].length}}</span>
                    <i
                        class="bk-cc-icon icon-cc-hide"
                        @click.stop="handleToggleGroup(group)"
                    >
                    </i>
                </div>
                <cmdb-collapse-transition name="model-box">
                    <ul class="model-box" v-show="topoNav.activeGroupId === group['bk_classification_id']">
                        <li
                            v-for="(model, modelIndex) in group['bk_objects']"
                            :key="modelIndex"
                            class="model-item"
                            :class="{
                                'invisible': topoNav.hideNodeIds.includes(model['bk_obj_id']),
                                'selected': topoNav.selectedNodeId === model['bk_obj_id']
                            }"
                            @click="handleSelectNode(model)"
                        >
                            <i :class="[
                                'node-icon',
                                model['bk_obj_icon'],
                                {
                                    'is-public': model.ispre
                                }
                            ]">
                            </i>
                            <div class="info">
                                <p class="name" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</p>
                            </div>
                            <i
                                class="bk-cc-icon icon-cc-hide"
                                @click.stop="handleToggleNode(model, group)"
                            >
                            </i>
                        </li>
                    </ul>
                </cmdb-collapse-transition>
            </li>
        </ul>

        <bk-sideslider
            v-transfer-dom
            :width="slider.width"
            :is-show.sync="slider.isShow"
            :title="slider.title"
            @hidden="handleSliderCancel">
            <component
                class="model-slider-content"
                slot="content"
                v-if="slider.isShow"
                :is="slider.content"
                v-bind="slider.properties"
                @save="handleSliderSave"
                @cancel="handleSliderCancel"
            ></component>
        </bk-sideslider>

        <div class="global-model" ref="topo" v-bkloading="{ isLoading: loading }"></div>

        <div class="topology-node-tooltips" v-show="topoTooltip.hoverNode" ref="nodeTooltips">
            <div
                class="icon-box"
                ref="addEdgeIcon"
                @click="handleAddEdge"
            >
                <i class="icon-cc-line"></i>
            </div>
            <div
                class="icon-box"
                ref="hideNodeIcon"
                @click="handleHideNode"
            >
                <i class="icon-cc-hide"></i>
            </div>
        </div>

        <the-create-model
            :is-show.sync="addBusinessLevel.showDialog"
            :is-main-line="true"
            :title="$t('新建层级')"
            @confirm="handleCreateBusinessLevel"
        ></the-create-model>
    </div>
</template>

<script>
    import cytoscape from 'cytoscape'
    import edgehandles from 'cytoscape-edgehandles'
    import popper from 'cytoscape-popper'
    import theRelation from './children/create-relation'
    import theRelationDetail from './children/relation-detail'
    import theCreateModel from '@/components/model-manage/_create-model'
    import { generateObjIcon as GET_OBJ_ICON } from '@/utils/util'
    import { mapGetters, mapActions } from 'vuex'
    import memoize from 'lodash.memoize'
    import debounce from 'lodash.debounce'
    import throttle from 'lodash.throttle'

    // cytoscape实例，不能放到data中管理
    let cy = null
    // edge操作实例
    let eh = null

    const NODE_WIDTH = 55

    export default {
        components: {
            theRelation,
            theRelationDetail,
            theCreateModel
        },
        data () {
            return {
                specialModel: ['process', 'plat'],

                // 关联数据
                associationList: [],

                // 节点数据
                localTopoModelList: [],

                // 主线模型
                mainLineModelList: [],

                slider: {
                    width: 514,
                    isShow: false,
                    content: '',
                    properties: {},
                    title: this.$t('拓扑显示设置')
                },
                topoTooltip: {
                    hoverNode: null
                },
                topoEdit: {
                    isEdit: false
                },
                topoNav: {
                    activeGroupId: '',
                    // 选中的分组id，-1全部
                    selectedGroupId: -1,
                    hideGroupIds: [],
                    // 目前是偏平的结构，如果有查找的性能问题，可以考虑以groupId分组
                    hideNodeIds: [],
                    isSelectAll: true,
                    selectedNodeId: ''
                },
                addBusinessLevel: {
                    showDialog: false,
                    parent: null
                },
                loading: true,
                isTopoHover: false,
                createAuth: false
            }
        },
        computed: {
            ...mapGetters(['isAdminView', 'isBusinessSelected', 'supplierAccount', 'userName', 'featureTipsParams', 'mainFullScreen']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('objectModelClassify', [
                'classifications',
                'getModelById'
            ]),
            noPositionModels () {
                return this.localTopoModelList.filter(model => {
                    const position = model.position
                    const isMainNode = this.isMainNode(model)
                    return position.x === null && position.y === null && !isMainNode
                })
            },
            localClassifications () {
                return this.$tools.clone(this.classifications).map(classify => {
                    classify['bk_objects'] = classify['bk_objects'].filter(model => {
                        return !this.specialModel.includes(model['bk_obj_id']) && !model.bk_ispaused
                    })
                    return classify
                })
            },
            hideModelConfigKey () {
                return `${this.userName}_model_${this.isAdminView ? 'adminView' : this.bizId}_hide_models`
            },
            hideModels () {
                return this.usercustom[this.hideModelConfigKey] || {}
            }
        },
        watch: {
            'topoNav.hideNodeIds' (hideNodeIds) {
                this.saveHideModelConfig(hideNodeIds)
                this.toggleAddBusinessBtn()
            }
        },
        created () {
            if (typeof cytoscape('core', 'edgehandles') !== 'function') {
                cytoscape.use(edgehandles)
            }
            if (typeof cytoscape('core', 'popper') !== 'function') {
                cytoscape.use(popper)
            }

            // 已记录的隐藏节点信息
            const { hideNodeIds, hideGroupIds } = this.hideModels
            this.topoNav.hideNodeIds = hideNodeIds || []
            this.topoNav.hideGroupIds = hideGroupIds || []
        },
        mounted () {
            this.getMainLineModel()
            this.initNetwork()
        },
        destroyed () {
            cy = null
            eh = null

            // 取消全屏
            this.resizeFull(true)
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchAssociationType'
            ]),
            ...mapActions('objectMainLineModule', [
                'searchMainlineObject',
                'createMainlineObject'
            ]),
            ...mapActions('objectModelClassify', [
                'searchClassificationsObjects'
            ]),
            getMainLineModel () {
                return this.searchMainlineObject({})
            },
            getAssociationType () {
                return this.searchAssociationType({
                    params: {},
                    config: {
                        requestId: 'searchAssociationType'
                    }
                }).then(res => res.info)
            },
            initNetwork () {
                cy = window.cy = cytoscape({
                    container: this.$refs.topo,
                    autolock: true,
                    zoom: 1,
                    minZoom: 0.1,
                    maxZoom: 5,
                    wheelSensitivity: 0.05,
                    pixelRatio: 2,

                    // 元素定义，支持promise
                    elements: this.getTopoElements(),

                    layout: {
                        name: 'preset',
                        fit: true,
                        padding: 30,
                        ready: () => {
                            this.loadNodeImage()
                        },
                        stop: () => {
                            this.updateElementPostion()
                        }
                    },

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
                            selector: 'node.model',
                            style: {
                                'width': NODE_WIDTH,
                                'height': NODE_WIDTH,

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
                            selector: 'node.model.bg',
                            style: {
                                'background-image': 'data(bg.unselected)'
                            }
                        },
                        {
                            selector: 'node.model.hover, node.model:selected',
                            style: {
                                'background-image': 'data(bg.selected)',
                                'border-color': '#3a84ff',
                                'font-weight': 'bold'
                            }
                        },
                        {
                            selector: 'node.model.mask',
                            style: {
                                'opacity': 0.16
                            }
                        },

                        // 添加按钮
                        {
                            selector: 'node.add-business-btn',
                            style: {
                                'width': 20,
                                'height': 20,
                                'color': '#ffffff',
                                'text-valign': 'bottom',
                                'text-halign': 'center',
                                'font-size': '20px',
                                'text-margin-y': '-19px',
                                'font-family': 'arial',
                                'label': '+',
                                'shape': 'round-rectangle',
                                'background-color': '#3c96ff',
                                'display': 'none'
                            }
                        },

                        // edge样式配置
                        {
                            selector: 'edge.model',
                            style: {
                                'curve-style': 'bezier',
                                'label': 'data(label)',
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
                                'font-size': '14px',
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
                            selector: 'edge.model.hover',
                            style: {
                                'width': 3,
                                'line-color': '#3c96ff',
                                'source-arrow-color': '#3c96ff',
                                'target-arrow-color': '#3c96ff',
                                'font-weight': 'bold'
                            }
                        },
                        {
                            selector: 'edge.model.mask',
                            style: {
                                'opacity': 0.16
                            }
                        },

                        {
                            selector: '.edge-editing',
                            style: {
                                'curve-style': 'bezier',
                                'label': 'data(label)'
                            }
                        },

                        // edgehandle样式定义
                        {
                            selector: '.eh-handle',
                            style: {
                                // 不需要控制点
                                'display': 'none'
                            }
                        },
                        {
                            selector: '.eh-hover',
                            style: {
                                'background-color': '#ffb23a'
                            }
                        },
                        {
                            selector: '.eh-source',
                            style: {
                                'border-width': 2,
                                'border-color': '#ffb23a'
                            }
                        },
                        {
                            selector: '.eh-target',
                            style: {
                                'border-width': 2,
                                'border-color': '#ffb23a'
                            }
                        },
                        {
                            selector: '.eh-preview, .eh-ghost-edge, .edge-editing',
                            style: {
                                'curve-style': 'bezier',
                                'target-arrow-shape': 'triangle-backcurve',
                                'background-color': '#ffb23a',
                                'line-color': '#ffb23a',
                                'line-style': 'dashed',
                                'target-arrow-color': '#ffb23a',
                                'source-arrow-color': '#ffb23a'
                            }
                        },

                        {
                            selector: '.eh-ghost-edge.eh-preview-active',
                            style: {
                                'opacity': 0
                            }
                        }
                    ]
                })

                // 所有操作的事件绑定
                cy.on('ready', (event) => {
                    const cy = event.cy

                    // 初始化节点隐藏
                    cy.batch(() => {
                        this.topoNav.hideNodeIds.forEach((id) => {
                            cy.$(`node#${id}`).style('visibility', 'hidden').connectedEdges().style('visibility', 'hidden')
                        })
                    })
                }).on('layoutstop', (event) => {
                    this.fitMaxZoom(event.cy)
                }).on('resize', debounce((event) => {
                    const cy = event.cy
                    cy.fit()
                    this.fitMaxZoom(cy)
                }, 500)).on('mouseover', 'node.model', throttle((event) => {
                    const node = event.target
                    const nodeData = node.data()

                    // 添加hover样式
                    node.addClass('hover')
                    node.connectedEdges().addClass('hover')

                    // 显示tooltip
                    if (this.topoEdit.isEdit && !this.specialModel.includes(nodeData.id)) {
                        // 设置tooltip状态数据
                        this.topoTooltip.hoverNode = nodeData

                        // todo根据画布缩放值更新操作按钮大小

                        // 每次重新创建因content引用的内容只能移动一次无法反复使用
                        const popover = this.$bkPopover(node.popperRef(), {
                            content: this.$refs.nodeTooltips,
                            hideOnClick: true,
                            sticky: true,
                            placement: 'right',
                            interactive: true,
                            animateFill: false,
                            theme: 'node-tooltip',
                            boundary: this.$refs.topo,
                            trigger: 'manual',
                            distance: 6,
                            offset: 12
                        })

                        node.data('popover', popover)
                        popover.show()
                    }
                }, 60)).on('mouseout', 'node.model', throttle((event) => {
                    const node = event.target
                    node.removeClass('hover')
                    node.connectedEdges().removeClass('hover')

                    const popover = node.data('popover')
                    if (popover) {
                        popover.hide()
                    }

                    this.topoTooltip.hoverNode = null
                }), 60).on('dragfreeon', 'node.model', (event) => {
                    const node = event.target
                    const nodeData = node.data()
                    const position = node.position()
                    this.updateSingleNodePosition({
                        bk_obj_id: nodeData.id,
                        bk_inst_id: nodeData.instId,
                        node_type: nodeData.type,
                        position: {
                            x: Math.round(position.x),
                            y: Math.round(position.y)
                        }
                    })
                }).on('mouseover', 'edge', (event) => {
                    event.target.addClass('hover')
                    this.isTopoHover = true
                }).on('mouseout', 'edge', (event) => {
                    event.target.removeClass('hover')
                    this.isTopoHover = false
                }).on('click', 'edge', (event) => {
                    const edgeData = event.target.data()
                    this.slider.title = edgeData.label
                    this.slider.properties = {
                        objId: edgeData.source,
                        isEdit: this.topoEdit.isEdit,
                        asstId: edgeData.instId
                    }
                    this.showSlider('theRelationDetail')
                }).on('ehcomplete', (event, sourceNode, targetNode, addedEles) => {
                    this.slider.properties = {
                        fromObjId: sourceNode.data('id'),
                        toObjId: targetNode.data('id'),
                        topoModelList: this.localTopoModelList
                    }
                    this.slider.title = this.$t('新建关联')
                    this.showSlider('theRelation')
                }).on('ehhoverover', (event, sourceNode, targetNode) => {
                    targetNode.data('ehhoverover', true)
                }).on('ehhoverout', (event, sourceNode, targetNode) => {
                    targetNode.data('ehhoverover', false)
                }).on('click', 'node.add-business-btn', (event) => {
                    const node = event.target
                    this.handleAddBusinessLevel(node.data('model'))
                }).on('mouseover', 'node.add-business-btn', (event) => {
                    this.isTopoHover = true
                }).on('mouseout', 'node.add-business-btn', (event) => {
                    this.isTopoHover = false
                })
            },
            async updateNetwork () {
                // 全量更新画布元素，如存在性能问题则需要依赖数据返回做按需更新
                const elements = await this.getTopoElements()
                cy.json({ elements })

                this.loadNodeImage()
                this.updateElementPostion()
            },
            async getTopoElements () {
                const [asstData, mainLineData, nodeData] = await Promise.all([
                    this.getAssociationType(),
                    this.getMainLineModel(),
                    this.$store.dispatch('globalModels/searchModelAction', this.$injectMetadata())
                ])

                this.associationList = asstData
                this.localTopoModelList = nodeData
                this.mainLineModelList = mainLineData

                const elements = []

                this.loading = false

                // 包含分类属性的节点数据
                const nodeObjects = this.localClassifications.reduce((acc, cur) => acc.concat(cur['bk_objects']), [])

                this.localTopoModelList.forEach((nodeItem, i) => {
                    // nodes，模型节点
                    const nodeObjId = nodeItem.bk_obj_id
                    const isMainNode = this.isMainNode(nodeItem)
                    const nodeClasses = ['model']
                    if (isMainNode) {
                        nodeClasses.push('main')
                    }
                    if (nodeItem.ispre) {
                        nodeClasses.push('ispre')
                    }
                    elements.push({
                        data: {
                            id: nodeObjId,
                            name: nodeItem.node_name,
                            icon: nodeItem.bk_obj_icon,
                            groupId: (nodeObjects.find(item => item.bk_obj_id === nodeObjId) || {}).bk_classification_id,
                            instId: nodeItem.bk_inst_id,
                            type: nodeItem.node_type
                        },
                        position: {
                            x: nodeItem.position.x || 0,
                            y: nodeItem.position.y || 0
                        },
                        group: 'nodes',
                        locked: false,
                        selectable: false,
                        classes: nodeClasses.join(' ')
                    })

                    // edges
                    if (Array.isArray(nodeItem.assts) && nodeItem.assts.length) {
                        nodeItem.assts.forEach((asstItem, asstIndex) => {
                            // 关联关系源数据
                            const { direction, asstName, asstId } = this.getAsstDetail(asstItem['bk_asst_inst_id'])

                            // 所关联的节点必须存在
                            if (this.localTopoModelList.findIndex(({ bk_obj_id: objId }) => objId === asstItem.bk_obj_id) !== -1) {
                                elements.push({
                                    data: {
                                        id: asstItem['bk_inst_id'],
                                        label: asstName || asstId,
                                        source: nodeItem.bk_obj_id,
                                        target: asstItem.bk_obj_id,
                                        direction,
                                        instId: asstItem['bk_inst_id']
                                    },
                                    group: 'edges',
                                    selectable: true,
                                    classes: 'model'
                                })
                            }
                        })
                    }
                })

                // nodes，添加业务层级操作按钮
                this.mainLineModelList.forEach((model, i) => {
                    if (this.canAddBusinessLevel(model)) {
                        elements.push({
                            data: {
                                id: `add-business-btn-${model.bk_obj_id}`,
                                model
                            },
                            group: 'nodes',
                            locked: false,
                            classes: 'add-business-btn'
                        })
                    }
                })

                return elements
            },
            loadNodeImage () {
                // 缓存调用结果，减少相同icon的转换开销
                const makeSvg = memoize(this.makeSvg, data => data.icon)
                cy.nodes('.model').forEach(async (node, i) => {
                    const svg = await makeSvg(node.data())
                    node.data('bg', svg)
                    node.addClass('bg')
                })
            },
            updateElementPostion () {
                const extent = cy.extent()
                const isEdit = this.topoEdit.isEdit

                // 先给节点解锁
                cy.autolock(false)

                // 1. 设置主节点位置
                const centerPos = { x: (extent.x1 + extent.x2) / 2, y: (extent.y1 + extent.y2) / 2 }
                const startPosY = extent.y1 + NODE_WIDTH
                // const nodeSpace = extent.h * 0.8 / this.mainLineModelList.length
                const nodeSpace = 200

                // 坚排并lock
                this.mainLineModelList.forEach((model, i) => {
                    cy.nodes(`#${model['bk_obj_id']}`).position({
                        x: centerPos.x,
                        y: i * nodeSpace + startPosY
                    }).lock()
                })

                // 2. 摆放添加业务层级按钮节点
                cy.nodes('.add-business-btn').positions((node, i) => {
                    // 所属模型节点信息
                    const modelNodeId = node.data('model').bk_obj_id
                    const modelNode = cy.nodes(`#${modelNodeId}`)
                    const modelNodePos = modelNode.position()
                    const modelNodeHeight = modelNode.outerHeight() + 10

                    return {
                        x: modelNodePos.x,
                        y: modelNodePos.y + modelNodeHeight
                    }
                }).style('display', isEdit ? 'element' : 'none').lock()

                // 3. 摆放无位置节点
                const nodeCollection = cy.collection()
                this.noPositionModels.forEach((model, i) => {
                    const node = cy.nodes(`#${model['bk_obj_id']}`)
                    nodeCollection.merge(node)
                })
                const collectionBoundingBox = nodeCollection.boundingBox()
                const nodeTotal = nodeCollection.length
                const nodeGutter = 15
                // 设定一行最多5个
                const maxCountInOneRow = Math.min(nodeTotal, 5)
                const boundingBoxW = (collectionBoundingBox.w + nodeGutter) * maxCountInOneRow
                const rowTotal = Math.ceil(nodeTotal / maxCountInOneRow)
                const boundingBoxH = collectionBoundingBox.h * rowTotal
                nodeCollection.layout({
                    name: 'grid',
                    fit: false,
                    padding: 30,
                    rows: rowTotal,
                    boundingBox: { x1: extent.x2, y1: extent.y1, w: boundingBoxW, h: boundingBoxH },
                    stop: () => {
                        cy.fit()
                    }
                }).run()

                // 更新节点锁状态
                cy.autolock(!isEdit)
            },
            handleToggleGroup (group) {
                const groupId = group['bk_classification_id']
                const index = this.topoNav.hideGroupIds.indexOf(groupId)
                let display
                if (index !== -1) {
                    this.topoNav.hideGroupIds.splice(index, 1)
                    display = true
                } else {
                    this.topoNav.hideGroupIds.push(groupId)
                    display = false
                }
                this.toggleNodeByGroup(group, display)
            },
            handleToggleNode (model, group) {
                const nodeId = model['bk_obj_id']
                const groupId = group['bk_classification_id']

                // 当前节点在隐藏列表中的索引
                const index = this.topoNav.hideNodeIds.indexOf(nodeId)

                if (index !== -1) {
                    this.topoNav.hideNodeIds.splice(index, 1)

                    // 即时切换拓扑图中的节点显示状态
                    cy.$(`node#${nodeId}`).style('visibility', 'visible').connectedEdges().style('visibility', 'visible')
                } else {
                    this.topoNav.hideNodeIds.push(nodeId)
                    cy.$(`node#${nodeId}`).style('visibility', 'hidden').connectedEdges().style('visibility', 'hidden')
                }

                // 节点所关联的组中所有节点id
                const nodeIds = group['bk_objects'].map(model => model['bk_obj_id'])
                const nodeCount = nodeIds.length
                const hideNodeCount = this.topoNav.hideNodeIds.filter(id => nodeIds.includes(id)).length
                const hideGroupIndex = this.topoNav.hideGroupIds.indexOf(groupId)

                // 与group选择状态联动
                if (hideGroupIndex !== -1 && hideNodeCount !== nodeCount) {
                    this.topoNav.hideGroupIds.splice(hideGroupIndex, 1)
                }
                if (hideNodeCount === nodeCount) {
                    this.topoNav.hideGroupIds.push(groupId)
                }
            },
            handleSelectGroup (group) {
                if (group) {
                    const groupId = group['bk_classification_id']
                    const groupNodes = cy.$(`node[groupId='${groupId}']`)

                    // 通过样式降低其它节点透明度，使用batch降低开销
                    cy.startBatch()
                    cy.$('*').addClass('mask')
                    groupNodes.removeClass('mask').edgesWith(groupNodes).removeClass('mask')
                    cy.endBatch()

                    this.topoNav.selectedGroupId = group['bk_classification_id']
                } else {
                    // 选择全部
                    this.topoNav.selectedGroupId = -1
                    cy.$('*').removeClass('mask')
                }

                // 取消单个节点选择
                this.topoNav.selectedNodeId = ''
            },
            handleSelectNode (model) {
                const nodeId = model['bk_obj_id']
                this.topoNav.selectedNodeId = nodeId

                cy.startBatch()
                cy.$('*').addClass('mask')
                cy.$(`node#${nodeId}`).removeClass('mask')
                cy.endBatch()

                // 取消组选择
                this.topoNav.selectedGroupId = ''
            },
            toggleNodeByGroup (group, display) {
                const groupId = group['bk_classification_id']
                const nodeIds = group['bk_objects'].map(model => model['bk_obj_id'])

                if (display) {
                    // 显示则从隐藏记录中过滤掉
                    this.topoNav.hideNodeIds = this.topoNav.hideNodeIds.filter(id => !nodeIds.includes(id))
                } else {
                    this.topoNav.hideNodeIds = [...this.topoNav.hideNodeIds, ...nodeIds]
                }

                // 同时在拓扑图中显示/隐藏这组节点
                const visibility = display ? 'visible' : 'hidden'
                cy.$(`node[groupId='${groupId}']`).style('visibility', visibility).connectedEdges().style('visibility', visibility)
            },
            makeSvg (nodeData) {
                return new Promise((resolve, reject) => {
                    const image = new Image()
                    image.onload = () => {
                        const model = this.getModelById(nodeData.id)
                        const svg = {
                            unselected: `data:image/svg+xml;charset=utf-8,${encodeURIComponent(GET_OBJ_ICON(image, {
                                name: nodeData.name,
                                iconColor: model.ispre ? '#798aad' : '#3c96ff',
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
            handleRelationSave (params) {
                const fromNode = this.localTopoModelList.find(model => model['bk_obj_id'] === params['bk_obj_id'])
                if (!fromNode.hasOwnProperty('assts')) {
                    Object.assign(fromNode, { assts: [] })
                }
                fromNode.assts.push({
                    bk_asst_inst_id: this.associationList.find(asst => asst['bk_asst_id'] === params['bk_asst_id']).id,
                    bk_obj_id: params['bk_asst_obj_id'],
                    bk_inst_id: params.id,
                    asstInfo: params
                })
                // 完成edge添加
                this.completeEditingEdge(params)
            },
            handleRelationDetailSave (data) {
                if (data.type === 'delete') {
                    this.localTopoModelList.forEach(model => {
                        if (model.hasOwnProperty('assts')) {
                            const index = model.assts.findIndex(asst => {
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

                    // 删除edge
                    cy.edges(`[instId=${data.params.id}]`).remove()
                }
            },
            handleEditTopo () {
                this.topoEdit.isEdit = true

                // 解除锁定
                cy.autolock(false)

                // 启用或初始化edge编辑功能
                if (eh) {
                    eh.enable()
                } else {
                    eh = cy.edgehandles({
                        loopAllowed (node) {
                            return true
                        },
                        edgeParams (sourceNode, targetNode, i) {
                            return {
                                data: {
                                    label: ''
                                },
                                group: 'edges',
                                classes: 'edge-editing'
                            }
                        }
                    })
                }

                // 显示新建层级操作节点
                this.toggleAddBusinessBtn()
            },
            handleExitEdit () {
                this.topoEdit.isEdit = false
                cy.autolock(true)
                eh.disable()
                cy.nodes('.add-business-btn').style('display', 'none')
            },
            handleAddEdge () {
                const nodeId = this.topoTooltip.hoverNode.id
                const node = cy.$(`node#${nodeId}`)

                // 触发edge编辑，node为source
                eh.start(node)
            },
            handleHideNode () {
                const { id, groupId } = this.topoTooltip.hoverNode
                const group = this.localClassifications.find(item => item.bk_classification_id === groupId) || {}
                const model = this.localTopoModelList.find(item => item.bk_obj_id === id) || {}

                const node = cy.$(`node#${id}`)
                const popover = node.data('popover')
                if (popover) {
                    popover.hide()
                }

                this.handleToggleNode(model, group)
            },
            clearEditingEdge () {
                // 删除编辑中的edge
                cy.edges('.edge-editing').remove()
            },
            completeEditingEdge (params) {
                const asstInstId = this.associationList.find(asst => asst['bk_asst_id'] === params['bk_asst_id']).id
                const { direction, asstName, asstId } = this.getAsstDetail(asstInstId)
                const edge = cy.edges('.edge-editing')

                // style，使用model样式使其与初始化数据效果一致
                edge.removeClass('edge-editing').addClass('model')

                // update data，不可变属性需要使用move方法
                edge.move({
                    source: params.bk_obj_id,
                    target: params.bk_asst_obj_id
                })
                edge.data({
                    direction,
                    label: asstName || asstId,
                    instId: params.id
                })
            },
            handleSliderSave (params) {
                switch (this.slider.content) {
                    case 'theRelation':
                        this.handleRelationSave(params)
                        break
                    case 'theRelationDetail':
                        this.handleRelationDetailSave(params)
                        break
                    default:
                }
            },
            handleSliderCancel () {
                if (this.slider.content === 'theRelation') {
                    this.clearEditingEdge()
                }
                this.slider.isShow = false
            },
            canAddBusinessLevel (model) {
                return this.isAdminView && !['set', 'module', 'host'].includes(model['bk_obj_id'])
            },
            isMainNode (model) {
                const mainLineIds = this.mainLineModelList.map(model => model['bk_obj_id'])
                return mainLineIds.includes(model['bk_obj_id'])
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
            getAsstDetail (asstId) {
                const asst = this.associationList.find(asst => asst.id === asstId)
                return {
                    asstId: asst['bk_asst_id'],
                    asstName: asst['bk_asst_name'].length ? asst['bk_asst_name'] : asst['bk_asst_id'],
                    direction: asst.direction
                }
            },
            showSlider (content) {
                const {
                    slider
                } = this
                slider.content = content
                switch (content) {
                    case 'theRelation':
                    case 'theRelationDetail':
                    default:
                        slider.width = 514
                }
                slider.isShow = true
            },
            handleSlideGroup (group) {
                if (group['bk_classification_id'] !== this.topoNav.activeGroupId) {
                    this.topoNav.activeGroupId = group['bk_classification_id']
                } else {
                    this.topoNav.activeGroupId = ''
                }
            },
            resizeFit () {
                cy.fit()
            },
            resizeFull (reset) {
                const mainFullScreen = reset === true ? false : !this.mainFullScreen
                this.$store.commit('setLayoutStatus', { mainFullScreen })
            },
            zoomIn () {
                const zoom = cy.zoom()
                cy.zoom(zoom + 0.05)
            },
            zoomOut () {
                const zoom = cy.zoom()
                cy.zoom(zoom - 0.05)
            },
            updateSingleNodePosition (node) {
                this.$store.dispatch('globalModels/updateModelAction', {
                    params: {
                        origin: [node]
                    }
                })
            },
            handleAddBusinessLevel (model) {
                if (this.createAuth) {
                    this.addBusinessLevel.parent = model
                    this.addBusinessLevel.showDialog = true
                } else {
                    const addBusinessLevel = this.$refs.addBusinessLevel
                    if (addBusinessLevel) {
                        addBusinessLevel.$el && addBusinessLevel.$el.click()
                    }
                }
            },
            async handleCreateBusinessLevel (data) {
                try {
                    await this.createMainlineObject({
                        params: this.$injectMetadata({
                            'bk_asst_obj_id': this.addBusinessLevel.parent['bk_obj_id'],
                            'bk_classification_id': 'bk_biz_topo',
                            'bk_obj_icon': data['bk_obj_icon'],
                            'bk_obj_id': data['bk_obj_id'],
                            'bk_obj_name': data['bk_obj_name'],
                            'bk_supplier_account': this.supplierAccount,
                            'creator': this.userName
                        })
                    })

                    // 更新分组数据
                    await this.searchClassificationsObjects({
                        params: this.$injectMetadata({}),
                        config: {
                            clearCache: true,
                            requestId: 'post_searchClassificationsObjects'
                        }
                    })

                    // 更新拓扑图
                    this.updateNetwork()

                    this.cancelCreateBusinessLevel()
                } catch (e) {
                    console.log(e)
                }
            },
            cancelCreateBusinessLevel () {
                this.addBusinessLevel.parent = null
                this.addBusinessLevel.showDialog = false
            },
            saveHideModelConfig (hideNodeIds) {
                this.$store.dispatch('userCustom/saveUsercustom', {
                    [this.hideModelConfigKey]: { hideNodeIds, hideGroupIds: this.topoNav.hideGroupIds }
                })
            },
            toggleAddBusinessBtn () {
                if (this.topoEdit.isEdit) {
                    cy.nodes('.add-business-btn').forEach((node) => {
                        const model = node.data('model')
                        if (this.topoNav.hideNodeIds.includes(model.bk_obj_id)) {
                            node.style('display', 'none')
                        } else {
                            node.style('display', 'element')
                        }
                    })
                }
            },
            fitMaxZoom (cy) {
                const fitMaxZoom = 1
                if (cy.zoom() > fitMaxZoom) {
                    cy.zoom(fitMaxZoom)
                    cy.center()
                }
            },
            handleReceiveAuth (auth) {
                this.createAuth = auth
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topo-wrapper {
        position: relative;
        padding: 0;
        &.hover {
            cursor: pointer;
        }
    }
    .toolbar {
        padding: 9px 20px;
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
            >i {
                margin-left: 22px;
                font-size: 20px;
                cursor: pointer;
                outline: 0;
                color: #979ba5;
                padding: 6px;
                &:hover {
                    color: $cmdbBorderFocusColor;
                }
            }
        }
        .topo-legend {
            position: absolute;
            padding: 3px 10px;
            top: 57px;
            right: 8px;
            background: #fff;
            box-shadow: 0px 2px 1px 0px rgba(185, 203, 222, 0.5);
            font-size: 12px;
            z-index: 1;
            .legend-item {
                line-height: 30px;
                font-size: 0;
                &.custom i {
                    background: $cmdbBorderFocusColor;
                }
                &.built-in i {
                    background: #798aad;
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
        float: left;
        border: 1px solid $cmdbTableBorderColor;
        border-left: none;
        width: 210px;
        height: calc(100% - 50px);
        overflow: auto;
        background: #fff;
        @include scrollbar;
        .group-info {
            line-height: 42px;
            padding: 0 16px 0 5px;
            font-size: 14px;
            cursor: pointer;
            color: #63656e;
            position: relative;
            &.group-total {
                padding-left: 15px;
            }
            &:hover {
                background: #e1ecff;

                .icon-cc-hide {
                    display: inline-block;
                }
            }
            &:not(.group-total):hover {
                .model-count {
                    display: none;
                }
            }
            &.active {
                .icon-angle-down {
                    transform: rotate(180deg);
                }
            }
            &.selected {
                color: #3a84ff;
                background: #e1ecff;
                .model-count {
                    color: #fff;
                    background-color: #a2c5fd;
                }
            }
            &.invisible {
                .icon-cc-hide {
                    display: inline-block;
                }
                .group-name {
                    opacity: 0.5;
                }
                .model-count {
                    display: none;
                }
            }
            .model-count {
                position: absolute;
                right: 16px;
                top: 12px;
                padding: 0 5px;
                border-radius: 2px;
                font-size: 12px;
                color: #979ba5;
                background: #f0f1f5;
                height: 18px;
                line-height: 17px;
                text-align: center;
            }
            .toggle-arrow {
                padding: 0 8px 0 15px;
                margin-right: 2px;
            }
            .icon-angle-down {
                transition: all .2s;
                font-size: 12px;
                color: #979ba5;
                margin-top: -3px;
            }
            .icon-cc-hide {
                display: none;
                position: absolute;
                right: 16px;
                top: 12px;
                font-size: 18px;
                color: #979ba5;

                &:hover {
                    color: #3a84ff;
                }
            }
        }
        .model-box {
            padding: 8px 0;
        }
        .model-item {
            padding: 5px 16px 5px 20px;
            position: relative;
            cursor: pointer;
            &:hover {
                background: #ebf4ff;

                .icon-cc-hide {
                    display: inline-block;
                }
            }
            &.disabled {
                cursor: not-allowed;
                opacity: .6;
            }
            &.invisible {
                .icon-cc-hide {
                    display: inline-block;
                }
                .info,
                .node-icon {
                    opacity: 0.5;
                }
            }
            &.selected {
                background: #ebf4ff;
            }
            .node-icon {
                display: inline-block;
                margin-right: 5px;
                width: 36px;
                height: 36px;
                font-size: 20px;
                line-height: 34px;
                text-align: center;
                vertical-align: middle;
                color: $cmdbBorderFocusColor;
                border: 1px solid $cmdbTableBorderColor;
                border-radius: 50%;
                &.is-public {
                    color: #798aad;
                }
            }
            .info {
                display: inline-block;
                line-height: 18px;
                vertical-align: middle;
                font-size: 12px;
                .name {
                    @include ellipsis;
                    width: 100px;
                }
                .id {
                    color: $cmdbBorderColor;
                }
            }
            .icon-cc-hide {
                display: none;
                position: absolute;
                right: 16px;
                top: 14px;
                font-size: 18px;
                color: #979ba5;

                &:hover {
                    color: #3a84ff;
                }
            }
        }
    }
    .global-model {
        float: left;
        width: calc(100% - 210px);
        height: calc(100% - 50px);
        padding: 12px;
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
        font-size: 12px;
        line-height: 18px;
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
        .tooltips-option {
            padding: 0 5px;
        }
    }
    .topology-node-tooltips {
        color: #fff;
        .icon-box {
            display: block;
            height: 24px;
            width: 24px;
            line-height: 24px;
            font-size: 0px;
            border-radius: 12px;
            background: rgba(24, 24, 24, .8);
            text-align: center;
            cursor: pointer;
            white-space: nowrap;

            &+.icon-box {
                margin-top: 3px;
            }
        }
        [class^=icon-cc] {
            display: inline-block;
            vertical-align: middle;
            font-size: 12px;
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';

    .tippy-popper {
        transition: none!important;
    }

    .tippy-tooltip {
        &.node-tooltip-theme {
            background: none;
        }
    }
</style>
