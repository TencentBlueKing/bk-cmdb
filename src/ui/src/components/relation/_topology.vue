<template>
    <div class="relation-topology-layout" :class="{ 'full-screen': fullScreen }">
        <bk-button class="exit-full-screen" size="small" theme="default"
            v-show="fullScreen"
            @click="toggleFullScreen(false)">
            <i class="icon-cc-resize-small"></i>
            {{$t('退出')}}
        </bk-button>
        <div class="tolology-loading" v-bkloading="{ isLoading: $loading(getRelationRequestId) }">
            <div class="topology-container" ref="container">
            </div>
        </div>
        <ul class="topology-legend">
            <li class="legend-item"
                v-for="(legend, index) in legends"
                :key="index"
                :class="{ inactive: !legend.active }"
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
                @click="handleShowDetails">
                {{$t('详情信息')}}
            </a>
        </div>
        <cmdb-topo-details
            v-if="details.show"
            :full-screen="fullScreen"
            :obj-id="details.objId"
            :inst-id="details.instId"
            :title="details.title"
            :show.sync="details.show">
        </cmdb-topo-details>
    </div>
</template>

<script>
    import { mapActions, mapGetters } from 'vuex'
    import cmdbTopoDetails from './_details.vue'
    import Vis from 'vis'
    let NODE_ID = 0
    export default {
        components: {
            cmdbTopoDetails
        },
        data () {
            return {
                getRelationRequestId: null,
                fullScreen: false,
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
                            size: 16,
                            vadjust: -5
                        },
                        scaling: {
                            min: 15,
                            max: 25
                        },
                        widthConstraint: {
                            maximum: 200
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
                    objId: null,
                    instId: null
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount'])
        },
        async mounted () {
            try {
                const [rootData] = await this.getRelation(this.$parent.objId, this.$parent.formatedInst['bk_inst_id'])
                const validRelation = rootData.next.filter(next => !this.ignore.includes(next['bk_obj_id']))
                if (validRelation.length) {
                    this.$emit('on-relation-loaded', validRelation)
                }
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('objectRelation', ['getInstRelation', 'updateInstRelation']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
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
                this.network.focus(node.id, { scale: 0.8 })
                this.network.selectNodes([node.id])
                this.legends = node.legends
            },
            getRelation (objId, instId, node = null) {
                this.getRelationRequestId = `get_getInstRelation_${objId}_${instId}`
                return this.getInstRelation({
                    objId,
                    instId,
                    params: this.$injectMetadata(),
                    config: {
                        requestId: this.getRelationRequestId,
                        clearCache: true
                    }
                }).then(async data => {
                    this.legends = []
                    await this.updateNetwork(data[0], node)
                    this.resetNetwork(node)
                    return data
                })
            },
            async updateNetwork ({ curr, next, prev }, node = null) {
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
                                from: type === 'next' ? node.id : currentNode.id,
                                label: node.data['bk_asst_name']
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
                return { nodes: relationNodes, edges: relationEdges, legends: relationLegends }
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
                    // const $tooltips = this.$refs.tooltips
                    const view = this.network.getViewPosition()
                    const scale = this.network.getScale()
                    const nodeBox = this.network.getBoundingBox(nodeId)
                    const containerBox = this.$refs.container.getBoundingClientRect()
                    // const tooltipsBox = $tooltips.getBoundingClientRect()
                    const left = containerBox.width / 2 + (nodeBox.right - view.x) * scale
                    const top = containerBox.height / 2 + (nodeBox.top - view.y) * scale
                    this.$refs.tooltips.style.left = left + 'px'
                    this.$refs.tooltips.style.top = top + 'px'
                })
            },
            handleToggleNodeVisibility (legend) {
                ['prev', 'next'].forEach(type => {
                    legend.node[type].forEach(node => {
                        const nodeLevel = legend.node.level
                        const level = nodeLevel === 0 ? [-1, 1] : [nodeLevel > 0 ? (nodeLevel + 1) : (nodeLevel - 1)]
                        if (level.includes(node.level) && node.data['bk_obj_id'] === legend.id) {
                            node.hidden = legend.active
                        }
                    })
                })
                legend.active = !legend.active
                this.resetNetwork(legend.node)
            },
            handleShowDetails () {
                const nodeData = this.hoverNode.data
                this.details.objId = nodeData['bk_obj_id']
                this.details.instId = nodeData['bk_inst_id']
                this.details.title = `${nodeData['bk_obj_name']}-${nodeData['bk_inst_name']}`
                this.details.show = true
            },
            toggleFullScreen (fullScreen) {
                this.fullScreen = fullScreen
                this.$nextTick(() => {
                    this.network.redraw()
                    this.network.moveTo({ scale: 0.8 })
                })
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
</style>
