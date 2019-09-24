<template>
    <div class="hosts-layout">
        <cmdb-resize-layout class="hosts-topology fl"
            direction="right"
            :handler-offset="3"
            :min="200"
            :max="480"
            v-bkloading="{ isLoading: $loading(['getInstTopo', 'getInternalTopo']) }"
            :class="{ 'is-collapse': layout.topologyCollapse }">
            <p class="topology-tips" v-if="showTopologyTips" v-show="!layout.topologyCollapse">
                <i class="icon icon-cc-exclamation-tips"></i>
                <i18n path="主机拓扑提示">
                    <a href="javascript:void(0)" place="link" @click="handleTopologyTipsClick">{{$t('服务拓扑')}}</a>
                </i18n>
                <i class="bk-icon icon-close" @click="handleCloseTopologyTips"></i>
            </p>
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
                    <i class="fl"
                        v-if="[1, 2].includes(data.default)"
                        :class="{
                            'internal-node-icon': true,
                            'icon-cc-host-free-pool': data.default === 1,
                            'icon-cc-host-breakdown': data.default === 2,
                            'is-selected': node.selected
                        }">
                    </i>
                    <i :class="['node-icon fl', { 'is-selected': node.selected }]" v-else>{{data.bk_obj_name[0]}}</i>
                    <span :class="['node-host-count fr', { 'is-selected': node.selected }]"
                        v-if="data.hasOwnProperty('host_count')">
                        {{data.host_count}}
                    </span>
                    <span class="node-name">{{node.name}}</span>
                </template>
            </bk-big-tree>
            <i class="topology-collapse-icon bk-icon icon-angle-left"
                @click="layout.topologyCollapse = !layout.topologyCollapse">
            </i>
        </cmdb-resize-layout>
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
    import { MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
    import cmdbHostsTable from '@/components/hosts/table'
    export default {
        components: {
            cmdbHostsTable
        },
        data () {
            const showTopologyTips = window.localStorage.getItem('showTopologyTips')
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
                showTopologyTips: showTopologyTips === null,
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
                this.getHostCount()
            } catch (e) {
                console.log(e)
            }
        },
        beforeDestroy () {
            this.ready = false
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
                        'bk_inst_name': module.bk_module_name,
                        'host_count': module.host_count
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
            async getHostCount () {
                try {
                    const data = await this.$store.dispatch('objectMainLineModule/getInstTopoInstanceNum', {
                        bizId: this.bizId
                    })
                    this.setHostCount(data)
                } catch (e) {
                    console.error(e)
                }
            },
            setHostCount (data) {
                data.forEach(datum => {
                    const id = this.getNodeId(datum)
                    const node = this.$refs.tree.getNodeById(id)
                    if (node) {
                        const count = datum.host_count
                        this.$set(node.data, 'host_count', count > 999 ? '999+' : count)
                    }
                    const child = datum.child
                    if (Array.isArray(child) && child.length) {
                        this.setHostCount(child)
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
            },
            handleTopologyTipsClick () {
                this.$router.push({
                    name: MENU_BUSINESS_SERVICE_TOPOLOGY
                })
            },
            handleCloseTopologyTips () {
                this.showTopologyTips = false
                window.localStorage.setItem('showTopologyTips', false)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .hosts-layout{
        border-top: 1px solid $cmdbLayoutBorderColor;
        padding: 0;
        .hosts-topology {
            position: relative;
            width: 280px;
            height: 100%;
            border-right: 1px solid $cmdbLayoutBorderColor;
            &.is-collapse {
                width: 0 !important;
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
            overflow: hidden;
            height: 100%;
            padding: 20px;
        }
    }
    .topology-tips {
        font-size: 12px;
        line-height: 16px;
        margin: 10px 0 0 0;
        padding: 2px 16px;
        .icon,
        .icon-close,
        span,
        a {
            display: inline-block;
            vertical-align: baseline;
        }
        a {
            color: #3A84FF;
        }
        .icon, .icon-close {
            cursor: pointer;
            vertical-align: -1px;
        }
        .icon-close {
            margin-left: 10px;
            &:hover {
                color: #3c96ff;
            }
        }
    }
    .topology-tree {
        width: 100%;
        max-height: 100%;
        padding: 10px 0;
        @include scrollbar-y;
        .node-icon {
            display: block;
            width: 20px;
            height: 20px;
            margin: 8px 4px 8px 0;
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
        .node-name {
            height: 36px;
            line-height: 36px;
            overflow: hidden;
            @include ellipsis;
        }
        .node-host-count {
            padding: 0 5px;
            margin: 9px 20px 9px 4px;
            height: 18px;
            line-height: 17px;
            border-radius: 2px;
            background-color: #f0f1f5;
            color: #979ba5;
            font-size: 12px;
            text-align: center;
            &.is-selected {
                background-color: #a2c5fd;
                color: #fff;
            }
        }
        .internal-node-icon{
            width: 20px;
            height: 20px;
            line-height: 20px;
            text-align: center;
            margin: 8px 4px 8px 0;
            &.is-selected {
                color: #FFB400;
            }
        }
    }
    .hosts-table{
        margin-top: 20px;
    }
</style>
