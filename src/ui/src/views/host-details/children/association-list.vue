<template>
    <div class="association-list" v-bkloading="{ isLoading: loading }">
        <template v-for="(item, itemIndex) in list">
            <cmdb-host-association-list-table
                v-for="(association, associationIndex) in item.associations"
                :key="association.id"
                :type="item.type"
                :id="item.id"
                :association-type="item.associationType"
                :instances="association.instances"
                :visible="!(itemIndex || associationIndex)">
            </cmdb-host-association-list-table>
        </template>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import cmdbHostAssociationListTable from './association-list-table.vue'
    export default {
        name: 'cmdb-host-association-list',
        components: {
            cmdbHostAssociationListTable
        },
        computed: {
            ...mapGetters('hostDetails', [
                'mainLine',
                'source',
                'sourceInstances',
                'target',
                'targetInstances',
                'associationTypes'
            ]),
            id () {
                return parseInt(this.$route.params.id)
            },
            hasAssociation () {
                return !!(this.source.length || this.target.length)
            },
            list () {
                try {
                    const list = []
                    const associations = [...this.source, ...this.target]
                    associations.forEach((association, index) => {
                        const isSource = index < this.source.length
                        const modelId = isSource ? association.bk_asst_obj_id : association.bk_obj_id
                        const topologyData = isSource ? this.targetInstances : this.sourceInstances
                        const modelTopologyData = topologyData.find(data => data.bk_obj_id === modelId) || {}
                        const associationData = {
                            ...association,
                            instances: modelTopologyData.children || []
                        }
                        const item = list.find(item => {
                            return isSource ? item.source === 'host' : item.target === 'host'
                        })
                        if (item) {
                            item.associations.push(associationData)
                        } else {
                            list.push({
                                type: isSource ? 'source' : 'target',
                                id: modelId,
                                associationType: this.associationTypes.find(target => target.bk_asst_id === association.bk_asst_id),
                                associations: [associationData]
                            })
                        }
                    })
                    return list
                } catch (e) {
                    console.log(e)
                }
                return false
            },
            loading () {
                return this.$loading([
                    'getAssociation',
                    'getMainLine',
                    'getAssociationType',
                    'getInstRelation'
                ])
            }
        },
        created () {
            this.getData()
        },
        methods: {
            async getData () {
                try {
                    const [source, target, mainLine, associationTypes, root] = await Promise.all([
                        this.getAssociation({ bk_obj_id: 'host' }),
                        this.getAssociation({ bk_asst_obj_id: 'host' }),
                        this.getMainLine(),
                        this.getAssociationType(),
                        this.getInstRelation()
                    ])
                    const availabelSource = this.getAvailableAssociation(source)
                    const availabelTarget = this.getAvailableAssociation(target, availabelSource)
                    this.setState({
                        source: availabelSource,
                        target: availabelTarget,
                        mainLine: mainLine.filter(model => !['biz', 'host'].includes(model.bk_obj_id)),
                        associationTypes: associationTypes,
                        root
                    })
                } catch (e) {
                    this.setState({
                        source: [],
                        target: [],
                        mainLine: [],
                        associationTypes: [],
                        root: {
                            prev: [],
                            next: []
                        }
                    })
                    console.error(e)
                }
            },
            setState ({ source, target, mainLine, associationTypes, root }) {
                this.$store.commit('hostDetails/setAssociation', { type: 'source', association: source })
                this.$store.commit('hostDetails/setAssociation', { type: 'target', association: target })
                this.$store.commit('hostDetails/setMainLine', mainLine)
                this.$store.commit('hostDetails/setInstances', { type: 'source', instances: root.prev })
                this.$store.commit('hostDetails/setInstances', { type: 'target', instances: root.next })
                this.$store.commit('hostDetails/setAssociationTypes', associationTypes)
            },
            getAssociation (condition) {
                return this.$store.dispatch('objectAssociation/searchObjectAssociation', {
                    params: this.$injectMetadata({ condition }),
                    config: {
                        requestId: 'getAssociation'
                    }
                })
            },
            getMainLine () {
                return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {}, {
                    config: {
                        requestId: 'getMainLine'
                    }
                })
            },
            async getAssociationType () {
                const { info } = await this.$store.dispatch('objectAssociation/searchAssociationType', {}, {
                    config: {
                        requestId: 'getAssociationType'
                    }
                })
                return Promise.resolve(info)
            },
            async getInstRelation () {
                const [root] = await this.$store.dispatch('objectRelation/getInstRelation', {
                    objId: 'host',
                    instId: this.id,
                    params: this.$injectMetadata(),
                    config: {
                        requestId: 'getInstRelation'
                    }
                })
                return Promise.resolve(root)
            },
            getAvailableAssociation (data, reference = []) {
                return data.filter(association => {
                    const sourceId = association.bk_obj_id
                    const targetId = association.bk_asst_obj_id
                    const isMainLine = this.mainLine.some(model => [sourceId, targetId].includes(model.bk_obj_id))
                    const isExist = reference.some(target => target.id === association.id)
                    return !isMainLine && !isExist
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association-list {
        height: 100%;
    }
</style>
