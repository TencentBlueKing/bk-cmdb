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
                :visible="!(itemIndex || associationIndex)"
                @delete-association="handleDeleteAssociation(...arguments, item)">
            </cmdb-host-association-list-table>
        </template>
    </div>
</template>

<script>
    import cmdbHostAssociationListTable from './association-list-table.vue'
    export default {
        name: 'cmdb-host-association-list',
        components: {
            cmdbHostAssociationListTable
        },
        data () {
            return {
                mainLine: [],
                source: [],
                sourceInstances: [],
                target: [],
                targetInstances: [],
                associationTypes: []
            }
        },
        computed: {
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
                    this.mainLine = mainLine.filter(model => !['biz', 'host'].includes(model.bk_obj_id))
                    this.source = this.getAvailableAssociation(source)
                    this.target = this.getAvailableAssociation(target, this.source)
                    this.associationTypes = associationTypes
                    this.sourceInstances = root.prev
                    this.targetInstances = root.next
                } catch (e) {
                    this.mainLine = []
                    this.source = []
                    this.target = []
                    this.associationTypes = []
                    this.sourceInstances = []
                    this.targetInstances = []
                    console.error(e)
                }
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
            },
            handleDeleteAssociation (associationInstance, item) {
                const type = item.type
                const instances = {
                    source: this.targetInstances,
                    target: this.sourceInstances
                }
                const associations = instances[type].find(data => data.bk_obj_id === item.id)
                associations.children = associations.children.filter(association => association.asso_id !== associationInstance.asso_id)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association-list {
        height: 100%;
    }
</style>
