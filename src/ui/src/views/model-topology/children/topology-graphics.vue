<template>
    <div class="graphics-layout" @mousemove="handleMouseMove"></div>
</template>

<script>
    import Vis from 'vis'
    import uuid from 'uuid/v4'
    import Graphics from './graphics.js'
    import {svgToImageUrl} from '@/utils/util'
    export default {
        name: 'cmdb-topology-graphics',
        data () {
            return {
                topologyData: [],
                associationList: []
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
                    this.associationList = associationData.info
                    this.topologyData = topologyData
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
                } catch (e) {
                    this.$error(e.message)
                }
            },
            async createNodes () {
                const images = await this.createNodeImages()
                return this.topologyData.map(model => {
                    const id = model['bk_obj_id']
                    const position = model.position || {}
                    const data = {
                        id,
                        label: model['node_name'],
                        fixed: typeof position.x === 'number',
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
                                from: model['bk_obj_id'],
                                to: asst['bk_obj_id']
                            })
                        })
                    }
                })
                return edges
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