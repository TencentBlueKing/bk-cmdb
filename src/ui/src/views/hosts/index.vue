<template>
    <div class="hosts-layout">
        <div class="hosts-topology"
            v-bkloading="{ isLoading: $loading(['getInstTopo', 'getInternalTopo']) }"
            :class="{ 'is-collapse': layout.topologyCollapse }">
            <bk-big-tree class="topology-tree"
                ref="tree"
                selectable
                :expand-on-click="false"
                :options="{
                    idKey: getNodeId,
                    nameKey: 'bk_inst_name',
                    childrenKey: 'child'
                }"
                @select-change="handleSelectChange">
                <template slot-scope="{ node, data }">
                    <i v-if="[1, 2].includes(data.default)"
                        :class="{
                            'internal-node-icon': true,
                            'icon-cc-host-free-pool': data.default === 1,
                            'icon-cc-host-breakdown': data.default === 2,
                            'is-selected': node.selected
                        }">
                    </i>
                    <i :class="['node-icon', { 'is-selected': node.selected }]" v-else>{{data.bk_obj_name[0]}}</i>
                    {{node.name}}
                </template>
            </bk-big-tree>
            <i class="topology-collapse-icon bk-icon icon-angle-left"
                @click="layout.topologyCollapse = !layout.topologyCollapse">
            </i>
        </div>
        <cmdb-hosts-table class="hosts-main" ref="hostsTable"
            delete-auth=""
            :show-collection="true"
            :edit-auth="$OPERATION.U_HOST"
            :save-auth="$OPERATION.U_HOST"
            :transfer-resource-auth="$OPERATION.HOST_TO_RESOURCE"
            :columns-config-key="columnsConfigKey"
            :columns-config-properties="columnsConfigProperties">
        </cmdb-hosts-table>
    </div>
</template>

<script>
    import { mapGetters, mapActions, mapState } from 'vuex'
    import cmdbHostsTable from '@/components/hosts/table'
    export default {
        components: {
            cmdbHostsTable
        },
        data () {
            return {
                properties: {
                    biz: [],
                    host: [],
                    set: [],
                    module: []
                },
                filter: {
                    selectedNode: null
                },
                layout: {
                    topologyCollapse: false
                },
                ready: false
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'isAdminView']),
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('hosts', ['filterParams']),
            columnsConfigKey () {
                return `${this.userName}_host_${this.isAdminView ? 'adminView' : this.bizId}_table_columns`
            },
            columnsConfigProperties () {
                const setProperties = this.properties.set.filter(property => ['bk_set_name'].includes(property['bk_property_id']))
                const moduleProperties = this.properties.module.filter(property => ['bk_module_name'].includes(property['bk_property_id']))
                const hostProperties = this.properties.host
                return [...setProperties, ...moduleProperties, ...hostProperties]
            }
        },
        watch: {
            filterParams () {
                if (this.ready) {
                    this.getHostList()
                }
            }
        },
        async created () {
            try {
                const [topologyInstance] = await Promise.all([
                    this.getBusinessTopology(),
                    this.getProperties()
                ])
                const businessNodeId = this.getNodeId(topologyInstance[0])
                this.$refs.tree.setData(topologyInstance)
                this.$refs.tree.setExpanded(businessNodeId)
                this.$refs.tree.setSelected(businessNodeId, {
                    emitEvent: true
                })
                this.ready = true
            } catch (e) {
                console.log(e)
            }
        },
        beforeDestroy () {
            this.ready = true
        },
        methods: {
            ...mapActions('objectModelProperty', ['batchSearchObjectAttribute']),
            getProperties () {
                return this.batchSearchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: { '$in': Object.keys(this.properties) },
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: 'getHostProperties'
                    }
                }).then(result => {
                    Object.keys(this.properties).forEach(objId => {
                        this.properties[objId] = result[objId]
                    })
                    return result
                })
            },
            async getBusinessTopology () {
                const [instance, internal] = await Promise.all([
                    this.getInstanceTopology(),
                    this.getInternalModules()
                ])
                const root = instance[0] || {}
                const children = root.child || []
                children.unshift(...(internal.module || []).map(module => {
                    return {
                        'default': ['空闲机', 'idle machine'].includes(module.bk_module_name) ? 1 : 2,
                        'bk_obj_id': 'module',
                        'bk_inst_id': module.bk_module_id,
                        'bk_inst_name': module.bk_module_name
                    }
                }))
                return instance
            },
            getInstanceTopology () {
                return this.$store.dispatch('objectMainLineModule/getInstTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: 'getInstTopo'
                    }
                })
            },
            getInternalModules () {
                return this.$store.dispatch('objectMainLineModule/getInternalTopo', {
                    bizId: this.bizId,
                    config: {
                        requestId: 'getInternalTopo'
                    }
                })
            },
            getNodeId (data) {
                return `${data.bk_obj_id}-${data.bk_inst_id}`
            },
            handleSelectChange (node) {
                this.filter.selectedNode = node
                this.getHostList()
            },
            getParams () {
                const defaultModel = ['biz', 'set', 'module', 'host', 'object']
                const modelInstKey = {
                    biz: 'bk_biz_id',
                    set: 'bk_set_id',
                    module: 'bk_module_id',
                    host: 'bk_host_id',
                    object: 'bk_inst_id'
                }
                const params = {
                    bk_biz_id: this.bizId,
                    ip: this.filterParams.ip,
                    condition: defaultModel.map(model => {
                        return {
                            bk_obj_id: model,
                            condition: this.filterParams[model] || [],
                            fields: []
                        }
                    })
                }
                const selectedNode = this.filter.selectedNode
                const selectedModel = defaultModel.includes(selectedNode.data.bk_obj_id) ? selectedNode.data.bk_obj_id : 'object'
                const selectedModelCondition = params.condition.find(model => model.bk_obj_id === selectedModel)
                selectedModelCondition.condition.push({
                    field: modelInstKey[selectedModel],
                    operator: '$eq',
                    value: selectedNode.data.bk_inst_id
                })
                return params
            },
            getHostList (resetPage = true) {
                const params = this.getParams()
                this.$refs.hostsTable.search(this.bizId, params, resetPage)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .hosts-layout{
        height: 100%;
        padding: 0;
        display: flex;
        .hosts-topology {
            position: relative;
            flex: 280px 0 0;
            height: 100%;
            border-right: 1px solid $cmdbLayoutBorderColor;
            &.is-collapse {
                width: 0;
                flex: 0 0 0;
                .topology-collapse-icon:before {
                    display: inline-block;
                    transform: rotate(180deg);
                }
            }
            .topology-collapse-icon {
                position: absolute;
                left: 100%;
                top: 50%;
                width: 16px;
                height: 100px;
                line-height: 100px;
                background: $cmdbLayoutBorderColor;
                border-radius: 0px 12px 12px 0px;
                transform: translateY(-50%);
                text-align: center;
                font-size: 12px;
                color: #fff;
                cursor: pointer;
                &:hover {
                    background: #699DF4;
                }
            }
        }
        .hosts-main{
            flex: 1;
            width: 0;
            height: 100%;
            padding: 20px;
        }
    }
    .topology-tree {
        max-height: 100%;
        padding: 10px 0;
        @include scrollbar-y;
        .node-icon {
            display: inline-block;
            width: 20px;
            height: 20px;
            vertical-align: middle;
            border-radius: 50%;
            background-color: #C4C6CC;
            line-height: 1.666667;
            text-align: center;
            font-size: 12px;
            font-style: normal;
            color: #FFF;
            &.is-selected {
                background-color: #3A84FF;
            }
        }
        .internal-node-icon.is-selected {
            color: #FFB400;
        }
    }
    .hosts-table{
        margin-top: 20px;
    }
</style>
