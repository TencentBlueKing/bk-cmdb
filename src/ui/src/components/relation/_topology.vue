<template>
    <div class="relation-topology-layout">
        <div class="tolology-loading" v-bkloading="{isLoading: $loading(getRelationRequestId)}">
            <div class="topology-container" ref="container">
            </div>
        </div>
        <ul class="topology-legend">
            <li class="legend-item" 
                v-for="(legend, index) in legends"
                :key="index"
                :class="{inactive: !legend.active}"
                @click="handleToggleNodeVisibility(legend)">
                <i :class="legend.icon"></i>
                {{legend.name}} {{legend.count}}
            </li>
        </ul>
        <div class="topology-tooltips" ref="tooltips"
            v-if="hoverNode"
            @mouseover="handleTooltipsOver"
            @mouseleave="handleTooltipsLeave">
            <a class="tooltips-option" href="javascript:void(0)"
                @click="handleShowDetails">{{$t('Common["详情信息"]')}}</a>
            <a class="tooltips-option tooltips-option-delete" href="javascript:void(0)"
                v-if="hoverNode.level === 1"
                @click="handleDeleteRelation">
                {{$t('Common["删除关联"]')}}
            </a>
        </div>
        <div class="topology-details-layout" v-if="details.show">
            <div class="details-container" ref="detailsContainer" v-bkloading="{isLoading: $loading()}" v-click-outside="handleHideDetails">
                <div class="details-title" ref="detailsTitle">
                    {{details.title}}
                    <i class="bk-icon icon-close" @click="handleHideDetails"></i>
                </div>
                <cmdb-details ref="detailsPopup" class="details-popup"
                    :showOptions="false"
                    :inst="details.inst"
                    :properties="details.properties"
                    :property-groups="details.propertyGroups">
                </cmdb-details>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import Vis from 'vis'
    let NODE_ID = 0
    export default {
        data () {
            return {
                getRelationRequestId: null,
                network: null,
                nodes: [],
                edges: [],
                legends: [],
                hoverNode: null,
                hoverTimer: null,
                rootNode: null,
                options: {
                    physics: false,
                    interaction: {
                        dragNodes: false,
                        hover: true
                    },
                    edges: {
                        color: {
                            color: '#c3cdd7',
                            highlight: '#c3cdd7',
                            hover: '#c3cdd7'
                        },
                        smooth: {
                            enabled: false
                        },
                        arrows: 'middle'
                    },
                    nodes: {
                        shape: 'image',
                        font: {
                            color: '#737987',
                            size: 12,
                            vadjust: -5
                        },
                        scaling: {
                            min: 15,
                            max: 25
                        },
                        widthConstraint: {
                            maximum: 50
                        }
                    },
                    layout: {
                        hierarchical: {
                            direction: 'LR',
                            nodeSpacing: 90
                        }
                    }
                },
                ignore: ['plat'],
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
            ...mapGetters(['supplierAccount'])
        },
        async mounted () {
            try {
                await this.getRelation(this.$parent.objId, this.$parent.instId)
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('objectRelation', ['getInstRelation', 'updateInstRelation']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectCommonInst', ['searchInst']),
            ...mapActions('objectBiz', ['searchBusiness']),
            ...mapActions('hostSearch', ['getHostBaseInfo']),
            resetNetwork (node = null) {
                this.network && this.network.destroy()
                this.network = new Vis.Network(this.$refs.container, {
                    nodes: new Vis.DataSet(this.nodes),
                    edges: new Vis.DataSet(this.edges)
                }, this.options)
                this.network.on('selectNode', data => {
                    this.handleSelectNode(data.nodes[0])
                })
                this.network.on('hoverNode', data => {
                    this.handleHoverNode(data)
                })
                this.network.on('blurNode', data => {
                    this.handleBlurNode(data)
                })
                this.network.on('dragStart', data => {
                    this.handleDragStart(data)
                })
                node = node || this.nodes[0]
                this.network.focus(node.id, {scale: 0.8})
                this.network.selectNodes([node.id])
                this.legends = node.legends
            },
            getRelation (objId, instId, node = null) {
                this.getRelationRequestId = `get_${objId}_${instId}_relation`
                return this.getInstRelation({
                    objId,
                    instId,
                    config: {
                        requestId: this.getRelationRequestId
                    }
                }).then(async data => {
                    this.legends = []
                    await this.updateNetwork(data[0], node)
                    this.resetNetwork(node)
                    return data
                })
            },
            async updateNetwork ({curr, next, prev}, node = null) {
                node = node || await this.createRootNode(curr)
                node.next = node.next || []
                node.prev = node.prev || []
                const [nextData, prevData] = await Promise.all([
                    this.createRelationData(next, node, 'next'),
                    this.createRelationData(prev, node, 'prev')
                ])
                node.next = [...node.next, ...nextData.nodes]
                node.prev = [...node.prev, ...prevData.nodes]
                const allLegends = [...nextData.legends, ...prevData.legends]
                const uniqueLegends = []
                allLegends.forEach(legend => {
                    const uniqueLegend = uniqueLegends.find(uniqueLegend => uniqueLegend.id === legend.id)
                    if (uniqueLegend) {
                        uniqueLegend.count += legend.count
                    } else {
                        uniqueLegends.push(legend)
                    }
                })
                node.legends = uniqueLegends
            },
            async createRootNode (root) {
                const node = {
                    id: `${root['bk_obj_id']}_${root['bk_inst_id']}_${NODE_ID++}`,
                    label: root['bk_inst_name'],
                    data: root,
                    loaded: true,
                    children: [],
                    level: 0,
                    value: 25
                }
                root['_id'] = node.id
                try {
                    const image = await this.createNodeImages(root)
                    node.image = image
                } catch (e) {
                    node.shape = 'dot'
                }
                this.rootNode = node
                this.nodes.push(node)
                return node
            },
            async createRelationData (relation, currentNode, type) {
                const relationNodes = []
                const relationEdges = []
                const relationLegends = []
                for (let i = 0; i < relation.length; i++) {
                    const obj = relation[i]
                    if (this.ignore.includes(obj['bk_obj_id']) || !obj.count) continue
                    const children = obj.children
                    for (let j = 0; j < children.length; j++) {
                        const inst = children[j]
                        inst['bk_obj_id'] = obj['bk_obj_id']
                        if (!this.exist(currentNode, inst, type)) {
                            const node = {
                                id: `${inst['bk_obj_id']}_${inst['bk_inst_id']}_${NODE_ID++}`,
                                label: inst['bk_inst_name'],
                                data: inst,
                                loaded: false,
                                level: this.getRelationNodeLevel(currentNode, type),
                                children: [],
                                [type === 'next' ? 'prev' : 'next']: [currentNode],
                                [type]: [],
                                legends: [],
                                value: 15
                            }
                            inst['_id'] = node.id
                            currentNode.children.push(node)
                            const edge = {
                                to: type === 'next' ? currentNode.id : node.id,
                                from: type === 'next' ? node.id : currentNode.id
                            }
                            const legend = relationLegends.find(legend => legend.id === obj['bk_obj_id'])
                            if (legend) {
                                legend.count++
                            } else {
                                relationLegends.push({
                                    id: obj['bk_obj_id'],
                                    name: obj['bk_obj_name'],
                                    icon: obj['bk_obj_icon'],
                                    node: currentNode,
                                    active: true,
                                    count: 1
                                })
                            }
                            try {
                                const instImages = await this.createNodeImages(inst)
                                node.image = instImages
                            } catch (e) {
                                node.shape = 'dot'
                            }
                            relationNodes.push(node)
                            relationEdges.push(edge)
                        }
                    }
                }
                this.nodes = [...this.nodes, ...relationNodes]
                this.edges = [...this.edges, ...relationEdges]
                return {nodes: relationNodes, edges: relationEdges, legends: relationLegends}
            },
            exist (currentNode, newInst, type) {
                return currentNode[type].some(node => {
                    return node.data['bk_obj_id'] === newInst['bk_obj_id'] && node.data['bk_inst_id'] === newInst['bk_inst_id']
                })
            },
            getRelationNodeLevel (currentNode, type) {
                if (currentNode.level === 0) {
                    return type === 'next' ? currentNode.level + 1 : currentNode.level - 1
                } else if (currentNode.level < 0) {
                    return currentNode.level - 1
                } else {
                    return currentNode.level + 1
                }
            },
            async createNodeImages (data) {
                const image = await this.loadNodeImage(data)
                const base64 = {
                    selected: this.createBase64Image(image, [60, 150, 255]),
                    unselected: this.createBase64Image(image, [98, 104, 127])
                }
                return Promise.resolve(base64)
            },
            loadNodeImage (data) {
                return new Promise((resolve, reject) => {
                    let useDefaultWhenError = true
                    const image = new Image()
                    image.onload = () => {
                        resolve(image)
                    }
                    image.onerror = () => {
                        if (useDefaultWhenError) {
                            useDefaultWhenError = false
                            image.src = `./static/svg/cc-default.svg`
                        } else {
                            reject(new Error(`Can not load object icon, object id: ${data['bk_obj_id']}, object icon: ${data['bk_obj_icon']}`))
                        }
                    }
                    image.src = `./static/svg/${data['bk_obj_icon'].substr(5)}.svg`
                })
            },
            createBase64Image (image, rgb) {
                let canvas = document.createElement('canvas')
                const ctx = canvas.getContext('2d')
                ctx.clearRect(0, 0, canvas.width, canvas.height)
                canvas.width = image.width
                canvas.height = image.height
                ctx.drawImage(image, 0, 0, image.width, image.height)
                const imageData = ctx.getImageData(0, 0, image.width, image.height)
                for (let i = 0; i < imageData.data.length; i += 4) {
                    imageData.data[i] = rgb[0]
                    imageData.data[i + 1] = rgb[1]
                    imageData.data[i + 2] = rgb[2]
                }
                ctx.putImageData(imageData, 0, 0)
                const svg = `<svg xmlns="http://www.w3.org/2000/svg" stroke="" xmlns:xlink="http://www.w3.org/1999/xlink" width="50" height="50">
                    <rect x="" height="50" width="50" style="fill: #f9f9f9; border: none"/>
                    <image width="100%" xlink:href="${canvas.toDataURL('image/svg')}"></image>
                </svg>`
                canvas = null
                return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
            },
            async handleSelectNode (id) {
                const node = this.nodes.find(node => node.id === id)
                if (!node.loaded) {
                    this.hoverNode = null
                    const data = node.data
                    await this.getRelation(data['bk_obj_id'], data['bk_inst_id'], node)
                    node.loaded = true
                } else {
                    this.legends = node.legends
                }
            },
            handleHoverNode (data) {
                this.$refs.container.style.cursor = 'pointer'
                clearTimeout(this.hoverTimer)
                this.popupTooltips(data)
            },
            handleBlurNode (data) {
                this.$refs.container.style.cursor = 'default'
                this.hoverTimer = setTimeout(() => {
                    this.hoverNode = null
                }, 300)
            },
            handleDragStart (data) {
                this.hoverNode = null
            },
            handleTooltipsOver () {
                clearTimeout(this.hoverTimer)
            },
            handleTooltipsLeave () {
                this.hoverTimer = setTimeout(() => {
                    this.hoverNode = null
                }, 300)
            },
            popupTooltips (data) {
                const nodeId = data.node
                this.hoverNode = this.nodes.find(node => node.id === nodeId)
                this.$nextTick(() => {
                    const $tooltips = this.$refs.tooltips
                    const view = this.network.getViewPosition()
                    const scale = this.network.getScale()
                    const nodeBox = this.network.getBoundingBox(nodeId)
                    const containerBox = this.$refs.container.getBoundingClientRect()
                    const tooltipsBox = $tooltips.getBoundingClientRect()
                    const left = containerBox.width / 2 + (nodeBox.right - view.x) * scale
                    const top = containerBox.height / 2 + (nodeBox.top - view.y) * scale
                    this.$refs.tooltips.style.left = left + 'px'
                    this.$refs.tooltips.style.top = top + 'px'
                })
            },
            handleToggleNodeVisibility (legend) {
                ['prev', 'next'].forEach(type => {
                    legend.node[type].forEach(node => {
                        const level = legend.node.level === 0 ? [-1, 1] : [legend.node.level]
                        if (level.includes(node.level) && node.data['bk_obj_id'] === legend.id) {
                            node.hidden = legend.active
                        }
                    })
                })
                legend.active = !legend.active
                this.resetNetwork(legend.node)
            },
            async handleDeleteRelation () {
                try {
                    const hoverNodeData = this.$tools.clone(this.hoverNode.data)
                    const rootNodeData = this.rootNode.data
                    const properties = await this.getObjectProperties(rootNodeData['bk_obj_id'])
                    const relationProperty = properties.find(property => property['bk_asst_obj_id'] === hoverNodeData['bk_obj_id'])
                    const relationAfterDeleted = this.getRelationAfterDeleted(this.rootNode, hoverNodeData)
                    const data = {
                        updateType: 'remove',
                        objId: rootNodeData['bk_obj_id'],
                        relation: relationAfterDeleted,
                        id: relationProperty['bk_property_id'],
                        multiple: relationProperty['bk_property_type'] === 'multiasst',
                        value: hoverNodeData['bk_inst_id'],
                        params: {}
                    }
                    if (rootNodeData['bk_obj_id'] === 'host') {
                        data.params['bk_host_id'] = rootNodeData['bk_inst_id'].toString()
                    } else {
                        const keyMap = {
                            host: 'bk_host_id',
                            biz: 'bk_biz_id'
                        }
                        data[keyMap[rootNodeData['bk_obj_id']] || 'bk_inst_id'] = rootNodeData['bk_inst_id']
                    }
                    await this.updateInstRelation({
                        params: data,
                        config: {
                            requestId: `update_${rootNodeData['bk_obj_id']}_relation`
                        }
                    })
                    this.updateDeletedRelation(hoverNodeData)
                } catch (e) {
                    console.error(e)
                }
            },
            updateDeletedRelation (deleteNodeData) {
                const deleteNode = this.nodes.find(node => node.id === deleteNodeData['_id'])
                const allDeleteNodes = [deleteNode, ...this.getAllDeleteNodes(deleteNode)]
                this.nodes = this.nodes.filter(node => !allDeleteNodes.includes(node))
                this.resetNetwork()
            },
            getAllDeleteNodes (deleteNode) {
                let allDeleteNodes = []
                deleteNode.children.forEach(node => {
                    allDeleteNodes = [...allDeleteNodes, node, ...this.getAllDeleteNodes(node)]
                })
                return allDeleteNodes
            },
            getObjectProperties (objId) {
                return this.searchObjectAttribute({
                    params: {
                        'bk_supplier_account': this.supplierAccount,
                        'bk_obj_id': objId
                    },
                    config: {
                        requestId: `get_${objId}_attribute`,
                        fromCache: true
                    }
                })
            },
            getRelationAfterDeleted (rootNode, hoverNodeData) {
                const relation = []
                rootNode.next.forEach(node => {
                    const nodeData = node.data
                    if (nodeData['bk_obj_id'] === hoverNodeData['bk_obj_id']) {
                        relation.push(nodeData['bk_inst_id'])
                    }
                })
                return relation
            },
            async handleShowDetails () {
                const node = this.hoverNode
                const objId = node.data['bk_obj_id']
                this.details.show = true
                this.details.title = `${node.data['bk_obj_name']}-${node.data['bk_inst_name']}`
                const [properties, propertyGroups] = await Promise.all([
                    this.getObjectProperties(objId),
                    this.searchGroup({
                        objId,
                        config: {
                            requestId: `get_${objId}_property_groups`
                        }
                    })
                ])
                this.details.properties = properties
                this.details.propertyGroups = propertyGroups
                const inst = await this.getNodeDetails(node)
                this.details.inst = this.$tools.flatternList(properties, [inst])[0]
                this.$nextTick(() => {
                    const detailsContainerHeight = this.$refs.detailsContainer.getBoundingClientRect().height
                    const detailsTitleHeight = this.$refs.detailsTitle.getBoundingClientRect().height
                    this.$refs.detailsPopup.$el.style.height = detailsContainerHeight - detailsTitleHeight + 'px'
                })
            },
            handleHideDetails () {
                this.details = {
                    show: false,
                    title: '',
                    inst: {},
                    properties: [],
                    propertyGroups: []
                }
            },
            getNodeDetails (node) {
                const objId = node.data['bk_obj_id']
                const instId = node.data['bk_inst_id']
                let promise
                if (objId === 'host') {
                    promise = this.getHostDetails(instId)
                } else if (objId === 'biz') {
                    promise = this.getBusinessDetails(instId)
                } else {
                    promise = this.getInstDetails(objId, instId)
                }
                return promise
            },
            getHostDetails (hostId) {
                return this.getHostBaseInfo({hostId}).then(data => {
                    const inst = {}
                    data.forEach(field => {
                        inst[field['bk_property_id']] = field['bk_property_value']
                    })
                    return inst
                })
            },
            getBusinessDetails (businessId) {
                return this.searchBusiness({
                    params: {
                        condition: {'bk_biz_id': businessId},
                        fields: [],
                        page: {start: 0, limit: 1}
                    }
                }).then(({info}) => info[0])
            },
            getInstDetails (objId, instId) {
                return this.searchInst({
                    objId,
                    params: {
                        condition: {
                            [objId]: [{
                                field: 'bk_inst_id',
                                operator: '$eq',
                                value: instId
                            }]
                        },
                        fields: {},
                        page: {start: 0, limit: 1}
                    }
                }).then(({info}) => info[0])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-topology-layout {
        height: calc(100% - 64px);
        background-color: #f9f9f9;
        position: relative;
        .tolology-loading {
            height: 100%;
        }
        .topology-container {
            height: 100%;
        }
        .topology-legend {
            position: absolute;
            left: 10px;
            top: 10px;
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
        .topology-tooltips {
            position: absolute;
            top: 0;
            left: 0;
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
    }
    .topology-details-layout {
        position: fixed;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        text-align: right;
        &:before {
            content: "";
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
        }
        .details-container {
            position: relative;
            display: inline-block;
            width: 710px;
            min-height: 250px;
            max-height: 80%;
            margin: 0 45px 0 0;
            vertical-align: middle;
            text-align: left;
            background-color: #fff;
            box-shadow: 0px 2px 9.6px 0.4px rgba(0, 0, 0, 0.4);
            z-index: 100;
            .details-title {
                position: relative;
                height: 49px;
                padding: 0 0 0 16px;
                border-bottom: 1px solid $cmdbBorderColor;
                line-height: 48px;
                color: #333948;
                background-color: #f7f7f7;
                .icon-close {
                    position: absolute;
                    right: 6px;
                    top: 12px;
                    padding: 6px;
                    font-size: 12px;
                    cursor: pointer;
                    color: #333948;
                    border-radius: 50%;
                    &:hover {
                        background-color: #e5e5e5;
                    }
                }
            }
            .details-popup {
                padding-bottom: 20px;
                @include scrollbar-y;
            }
        }
    }
</style>