<template>
    <div class="association-list" v-bkloading="{ isLoading: loading }">
        <div class="association-empty" v-if="!hasRelation">
            <div class="empty-content">
                <i class="bk-icon icon-empty">
                </i>
                <span>{{$t('暂无关联关系')}}</span>
            </div>
        </div>
        <template v-else>
            <template v-for="item in list">
                <cmdb-relation-list-table
                    ref="associationListTable"
                    v-for="association in item.associations"
                    :key="association.id"
                    :type="item.type"
                    :id="item.id"
                    :source="sourceInstances"
                    :target="targetInstances"
                    :association-type="item.associationType">
                </cmdb-relation-list-table>
            </template>
        </template>
    </div>
</template>

<script>
    import bus from '@/utils/bus.js'
    import { mapGetters, mapActions } from 'vuex'
    import cmdbRelationListTable from './list-table.vue'
    export default {
        name: 'cmdb-relation-list',
        components: {
            cmdbRelationListTable
        },
        props: {
            associationObject: {
                type: Array,
                required: true
            },
            associationTypes: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                sourceInstances: [],
                targetInstances: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', ['models']),
            objId () {
                return this.$parent.objId
            },
            instId () {
                return this.$parent.formatedInst['bk_inst_id']
            },
            hasRelation () {
                return this.$parent.hasRelation
            },
            list () {
                try {
                    const list = []
                    this.associationObject.forEach((association, index) => {
                        const isSource = association['bk_obj_id'] === this.objId
                        const modelId = isSource ? association.bk_asst_obj_id : association.bk_obj_id
                        list.push({
                            type: isSource ? 'source' : 'target',
                            id: modelId,
                            associationType: this.associationTypes.find(target => target.bk_asst_id === association.bk_asst_id) || {},
                            associations: [association]
                        })
                    })
                    return list
                } catch (e) {
                    console.log(e)
                }
                return []
            },
            loading () {
                return this.$loading([
                    'getInstRelation'
                ])
            }
        },
        watch: {
            list () {
                this.expandFirstListTable()
            }
        },
        created () {
            this.getInstRelation()

            bus.$on('association-change', async () => {
                await this.getInstRelation()
            })

            this.expandFirstListTable()
        },
        beforeDestroy () {
            bus.$off('association-change')
        },
        methods: {
            ...mapActions('objectAssociation', [
                'searchObjectAssociation'
            ]),
            async getInstRelation () {
                const res = await this.$store.dispatch('objectRelation/getInstRelation', {
                    objId: this.objId,
                    instId: this.instId,
                    params: {},
                    config: {
                        requestId: 'getInstRelation'
                    }
                })

                if (res) {
                    const data = res[0] || { prev: [], next: [] }
                    this.sourceInstances = data.prev
                    this.targetInstances = data.next
                }
            },
            expandFirstListTable () {
                this.$nextTick(() => {
                    if (this.$refs.associationListTable) {
                        const [firstAssociationListTable] = this.$refs.associationListTable
                        firstAssociationListTable && (firstAssociationListTable.expanded = true)
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association-list {
        height: 100%;
    }
    .association-empty {
        height: 100%;
        text-align: center;
        font-size: 14px;
        &:before {
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
            content: "";
        }
        .empty-content {
            display: inline-block;
            vertical-align: middle;
            .bk-icon {
                display: inline-block;
                margin: 0 0 10px 0;
                font-size: 65px;
                color: #c3cdd7;
            }
            span {
                display: inline-block;
                width: 100%;
            }
        }
    }
</style>
