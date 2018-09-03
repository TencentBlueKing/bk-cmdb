<template>
    <div class="relation-topology-layout">
        <div class="topology-container" ref="topologyContainer"></div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import Vis from 'vis'
    let NODE_ID = 1
    export default {
        data () {
            return {
                network: null,
                nodes: null,
                edges: null,
                ignore: ['plat']
            }
        },
        async mounted () {
            try {
                this.initNetwork()
                await this.getRelation(this.$parent.objId, this.$parent.instId)
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapActions('objectRelation', ['getInstRelation']),
            initNetwork () {
                this.nodes = new Vis.DataSet()
                this.edges = new Vis.DataSet()
                this.network = new Vis.Network(this.$refs.topologyContainer, {
                    nodes: this.nodes,
                    edges: this.edges
                }, {})
                this.nodes.on('click', () => {
                    console.log(arguments)
                })
            },
            getRelation (objId, instId, node = null) {
                return this.getInstRelation({
                    objId,
                    instId,
                    config: {
                        requestId: `get_${objId}_${instId}_relation`
                    }
                }).then(data => {
                    this.updateNetwork(data[0], node)
                    return data
                })
            },
            updateNetwork ({curr, next, prev}, node = null) {
                const currentId = node ? node.id : `${curr['bk_obj_id']}_${curr['bk_inst_id']}_${NODE_ID++}`
                const nextData = this.createRelationData(next, currentId, 'next')
                const prevData = this.createRelationData(prev, currentId, 'prev')
                let nodes = [
                    ...nextData.nodes,
                    ...prevData.nodes
                ]
                if (!node) {
                    nodes.push({
                        id: currentId,
                        label: curr['bk_inst_name']
                    })
                }
                this.nodes.add(nodes)
                this.edges.add([
                    ...nextData.edges,
                    ...prevData.edges
                ])
            },
            createRelationData (relation, currentId, type) {
                const nodes = []
                const edges = []
                relation.forEach(obj => {
                    if (!this.ignore.includes(obj['bk_obj_id']) && obj.count) {
                        obj.children.forEach(inst => {
                            const id = `${obj['bk_obj_id']}_${inst['bk_inst_id']}_${NODE_ID++}`
                            nodes.push({
                                id,
                                label: inst['bk_inst_name']
                            })
                            edges.push({
                                to: type === 'next' ? currentId : id,
                                from: type === 'next' ? id : currentId
                            })
                        })
                    }
                })
                return {nodes, edges}
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-topology-layout {
        height: calc(100% - 64px);
        background-color: #f9f9f9;
        .topology-container {
            height: 100%;
        }
    }
</style>