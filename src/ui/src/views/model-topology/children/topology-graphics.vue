<template>
    <div class="graphics-layout"
        @dragover.prevent
        @drop="handleDrop($event)"
        @mousemove="handleMouseMove">
    </div>
</template>

<script>
    import Vis from 'vis'
    import uuid from 'uuid/v4'
    import Graphics from './graphics.js'
    import { svgToImageUrl } from '@/utils/util'
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-topology-graphics',
        data () {
            return {}
        },
        computed: {
            ...mapGetters('globalModels', [
                'isEditMode',
                'topologyData',
                'topologyMap',
                'options',
                'edgeOptions'
            ]),
            ...mapGetters('objectAssociation', ['associationList'])
        },
        watch: {
            isEditMode (val) {
                this.instance.changeMode(val)
            },
            options (options) {
                this.instance.updateOptions(options)
            },
            edgeOptions (edgeOptions) {
                this.instance.updateEdges(edgeOptions)
            }
        },
        mounted () {
            this.getGraphicsData()
        },
        methods: {
            async getGraphicsData () {
                try {
                    const [
                        associationData,
                        topologyData
                    ] = await Promise.all([
                        this.getAssociationData(),
                        this.getTopologyData()
                    ])
                    this.$store.commit('globalModels/setTopologyData', topologyData)
                    this.$store.commit('objectAssociation/setAssociationList', associationData.info)
                    this.initGraphics()
                } catch (e) {
                    this.associationList = []
                    this.topologyData = []
                    console.log(e)
                }
            },
            getAssociationData () {
                return this.$store.dispatch('objectAssociation/searchAssociationType', {
                    params: {},
                    config: {
                        requestId: 'searchAssociationType'
                    }
                })
            },
            getTopologyData () {
                return this.$store.dispatch('globalModels/searchModelAction', this.$injectMetadata())
            },
            async initGraphics () {
                try {
                    const nodes = await this.createNodes()
                    const edges = this.createEdges()
                    // this.instance特意设置为非响应式，因为Vue的数据劫持会影响Vis内部自定义的数据劫持
                    this.instance = new Graphics(this.$el, {
                        nodes: nodes,
                        edges: edges
                    })
                    this.instance.on('deleteNode', this.handleDeleteNode)
                    this.instance.on('dragNode', this.handleDragNode)
                    this.instance.on('addEdge', this.handleAddEdge)
                    this.instance.on('edgeClick', this.handleEdgeClick)
                    this.instance.on('stabilized', this.handleStabilized)
                } catch (e) {
                    this.$error(e.message)
                }
            },
            async createNodes () {
                const images = await this.createNodeImages()
                return this.topologyData.map(model => {
                    const id = model['bk_obj_id']
                    const position = model.position || {}
                    const fixed = typeof position.x === 'number'
                    const data = {
                        id,
                        label: model['node_name'],
                        fixed: fixed,
                        hidden: !fixed,
                        ...position
                    }

                    if (images[id]) {
                        data.image = images[id]
                    } else {
                        data.shape = 'dot'
                    }

                    return data
                })
            },
            createEdges () {
                const edges = []
                this.topologyData.forEach(model => {
                    if (Array.isArray(model.assts)) {
                        model.assts.forEach(asst => {
                            edges.push({
                                id: asst['bk_inst_id'],
                                from: model['bk_obj_id'],
                                to: asst['bk_obj_id'],
                                label: this.getEdgeLable(asst),
                                arrows: this.getEdgeArrows(asst),
                                data: asst
                            })
                        })
                    }
                })
                return edges
            },
            getEdgeLable (association) {
                const data = this.associationList.find(data => data.id === association['bk_asst_inst_id']) || {}
                return data['bk_asst_name'] || data['bk_asst_id']
            },
            createNodeImages () {
                const data = this.topologyData
                const images = {}
                const total = data.length
                let counter = 0
                let loadedResolver
                const checkCounter = () => {
                    if (counter === total) {
                        loadedResolver(images)
                    }
                }
                data.forEach(model => {
                    const image = new Image()
                    image.onload = () => {
                        images[model['bk_obj_id']] = {
                            unselected: svgToImageUrl(image, {
                                name: model['node_name'],
                                iconColor: this.$tools.getMetadataBiz(model) ? '#3c96ff' : '#868b97',
                                backgroundColor: '#fff'
                            }),
                            selected: svgToImageUrl(image, {
                                name: model['node_name'],
                                iconColor: '#fff',
                                backgroundColor: '#3a84ff'
                            })
                        }
                        counter++
                        checkCounter()
                    }
                    image.onerror = () => {
                        images[model['bk_obj_id']] = false
                        counter++
                        checkCounter()
                    }
                    image.src = `${window.location.origin}/static/svg/${model['bk_obj_icon'].substr(5)}.svg`
                })
                return new Promise((resolve, reject) => {
                    loadedResolver = resolve
                })
            },
            handleMouseMove (event) {
                if (this.instance) {
                    this.instance.shadowNodeFollowMouse(event)
                }
            },
            handleDeleteNode (nodeId, edges) {
                const edgeCount = edges.length
                if (edgeCount) {
                    this.$bkInfo({
                        title: this.$t('移除失败'),
                        subTitle: this.$tc('移除失败提示', edgeCount, { asstNum: edgeCount })
                    })
                    return false
                } else {
                    let resolver = null
                    const promise = new Promise(resolve => {
                        resolver = resolve
                    })
                    this.$bkInfo({
                        title: this.$t('确定移除模型?'),
                        subTitle: this.$t('移除模型提示'),
                        confirmFn: () => {
                            const data = this.topologyMap[nodeId]
                            this.updateSavedPosition([{
                                'bk_inst_id': data['bk_inst_id'],
                                'bk_obj_id': data['bk_obj_id'],
                                'node_type': data['node_type'],
                                'position': { x: null, y: null }
                            }])
                            resolver(true)
                        },
                        cancelFn: () => {
                            resolver(false)
                        }
                    })
                    return promise
                }
            },
            handleDragNode (nodeId, position) {
                const data = this.topologyMap[nodeId]
                this.updateSavedPosition([{
                    'bk_inst_id': data['bk_inst_id'],
                    'bk_obj_id': data['bk_obj_id'],
                    'node_type': data['node_type'],
                    'position': position
                }])
            },
            handleAddEdge (edge) {
                this.$store.commit('globalModels/setAssociation', {
                    show: true,
                    edge
                })
                return new Promise((resolve, reject) => {
                    const commitMethod = 'globalModels/setAddEdgePromise'
                    this.$store.commit(commitMethod, {
                        resolve: data => {
                            resolve({
                                id: data['bk_inst_id'],
                                label: this.getEdgeLable(data),
                                data: data,
                                arrows: this.getEdgeArrows(data)
                            })
                            this.$store.commit(commitMethod, { resolve: null, reject: null })
                        },
                        reject: result => {
                            this.$store.commit(commitMethod, { resolve: null, reject: null })
                            reject(result)
                        }
                    })
                })
            },
            getEdgeArrows (data) {
                const associationId = data['bk_asst_inst_id']
                const association = this.associationList.find(association => association.id === associationId)
                return {
                    from: association.direction === 'bidirectional',
                    to: association.direction !== 'none'
                }
            },
            handleEdgeClick (edge) {
                this.$store.commit('globalModels/setAssociation', {
                    show: true,
                    edge
                })
            },
            handleStabilized (positions) {
                const updateQueue = []
                Object.keys(positions).forEach(modelId => {
                    const data = this.topologyMap[modelId]
                    const newPosition = positions[modelId]
                    const oldPosition = data.position
                    if (newPosition.x !== oldPosition.x || newPosition.y !== oldPosition.y) {
                        updateQueue.push({
                            'bk_inst_id': data['bk_inst_id'],
                            'bk_obj_id': data['bk_obj_id'],
                            'node_type': data['node_type'],
                            'position': newPosition
                        })
                    }
                })
                this.updateSavedPosition(updateQueue)
            },
            handleDrop (event) {
                const modelId = event.dataTransfer.getData('modelId')
                if (modelId) {
                    const data = this.topologyMap[modelId]
                    const position = this.instance.convertNodePosition(event)
                    this.instance.showNode(modelId, position)
                    this.updateSavedPosition([{
                        'bk_inst_id': data['bk_inst_id'],
                        'bk_obj_id': data['bk_obj_id'],
                        'node_type': data['node_type'],
                        'position': position
                    }])
                }
            },
            // 后端对位置的存储为int64,需取整
            convertPosition (updateQueue) {
                updateQueue.forEach(queue => {
                    const position = queue.position
                    if (typeof position.x === 'number') {
                        position.x = Math.floor(position.x)
                        position.y = Math.floor(position.y)
                    }
                })
            },
            updateSavedPosition (updateQueue) {
                if (!updateQueue.length) {
                    return false
                }
                this.convertPosition(updateQueue)
                this.$store.commit('globalModels/updateTopologyData', updateQueue)
                this.$store.dispatch('globalModels/updateModelAction', {
                    params: this.$injectMetadata({
                        origin: updateQueue
                    })
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .graphics-layout {
        background-color: #f4f5f8;
        background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
        background-size: 10px 10px;
    }
</style>
