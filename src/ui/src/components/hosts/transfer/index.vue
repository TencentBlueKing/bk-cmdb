<template>
    <div class="transfer-layout clearfix" v-bkloading="{ isLoading: loading }">
        <div class="columns-layout fl">
            <div class="business-layout">
                <label class="business-label">{{$t('业务')}}</label>
                <cmdb-business-selector class="business-selector" v-model="businessId" disabled>
                </cmdb-business-selector>
            </div>
            <cmdb-auth style="display: none;"
                :auth="transferResourceAuthData"
                @update-auth="handleReceiveAuth">
            </cmdb-auth>
            <div class="tree-layout">
                <bk-big-tree ref="topoTree" class="topo-tree"
                    v-cursor="{
                        active: !hasTransferResourceAuth,
                        auth: [transferResourceAuth],
                        selector: '.is-root.is-first-child'
                    }"
                    expand-icon="bk-icon icon-down-shape"
                    collapse-icon="bk-icon icon-right-shape"
                    :options="{
                        idKey: getTopoNodeId,
                        nameKey: 'bk_inst_name',
                        childrenKey: 'child'
                    }"
                    :selectable="false"
                    :expand-on-click="false"
                    :show-checkbox="shouldShowCheckbox"
                    :before-check="beforeNodeCheck"
                    :check-strictly="false"
                    @node-click="handleNodeClick"
                    @check-change="handleNodeCheck">
                    <div class="node-info clearfix" slot-scope="{ node, data }">
                        <i :class="['node-model-icon fl', { 'is-checked': node.checked }]">{{modelIconMap[data.bk_obj_id]}}</i>
                        <span class="node-name">{{data.bk_inst_name}}</span>
                    </div>
                </bk-big-tree>
            </div>
        </div>
        <div class="columns-layout fl">
            <div class="selected-layout">
                <label class="selected-label">{{$t('已选中模块')}}</label>
            </div>
            <div class="modules-layout">
                <ul class="module-list">
                    <li class="module-item clearfix"
                        v-for="(node, index) in selectedModules" :key="index">
                        <div class="module-info fl">
                            <span class="module-info-name">{{node.data.bk_inst_name}}</span>
                            <span class="module-info-path">{{getModulePath(node)}}</span>
                        </div>
                        <i class="bk-icon icon-close fr" @click="removeSelectedModule(node)"></i>
                    </li>
                </ul>
            </div>
        </div>
        <div v-pre class="clearfix"></div>
        <div class="options-layout clearfix">
            <div class="increment-layout content-middle fl" v-if="showIncrementOption">
                <label class="cmdb-form-radio cmdb-radio-small" for="increment" :title="$t('增量更新')">
                    <input id="increment" type="radio" v-model="increment" :value="true">
                    <span class="cmdb-radio-text">{{$t('增量更新')}}</span>
                </label>
                <label class="cmdb-form-radio cmdb-radio-small" for="replacement" :title="$t('完全替换')">
                    <input id="replacement" type="radio" v-model="increment" :value="false">
                    <span class="cmdb-radio-text">{{$t('完全替换')}}</span>
                </label>
            </div>
            <div class="button-layout content-middle fr">
                <bk-button class="transfer-button" theme="primary"
                    :disabled="!selectedModules.length"
                    @click="handleTransfer">
                    {{$t('确认转移')}}
                </bk-button>
                <bk-button class="transfer-button" theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            selectedHosts: {
                type: Array,
                required: true,
                default () {
                    return []
                }
            },
            transferResourceAuth: {
                type: [String, Array],
                default: ''
            }
        },
        data () {
            return {
                businessId: '',
                topoModel: [],
                tree: {
                    data: []
                },
                selectedModules: [],
                increment: true,
                hasTransferResourceAuth: false
            }
        },
        computed: {
            transferResourceAuthData () {
                const auth = this.transferResourceAuth
                if (!auth) return {}
                if (Array.isArray(auth) && !auth.length) return {}
                return this.$authResources({ type: auth })
            },
            hostIds () {
                return this.selectedHosts.map(host => host['host']['bk_host_id'])
            },
            hostModules () {
                const modules = []
                this.selectedHosts.forEach(host => {
                    host.module.forEach(module => {
                        modules.push(module)
                    })
                })
                return modules
            },
            showIncrementOption () {
                const hasSpecialModule = this.selectedModules.some(node => {
                    return node.data.bk_inst_id === 'resource' || [1, 2].includes(node.data.default)
                })
                return this.selectedModules.length && this.selectedHosts.length > 1 && !hasSpecialModule
            },
            loading () {
                const requestIds = [
                    'get_searchMainlineObject',
                    `get_getInstTopo_${this.businessId}`,
                    `get_getInternalTopo_${this.businessId}`,
                    'post_transferHost'
                ]
                return this.$loading(requestIds)
            },
            mainLineModels () {
                const models = this.$store.getters['objectModelClassify/models']
                return this.topoModel.map(data => models.find(model => model.bk_obj_id === data.bk_obj_id))
            },
            modelIconMap () {
                const map = {}
                this.mainLineModels.forEach(model => {
                    map[model.bk_obj_id] = model.bk_obj_name[0]
                })
                return map
            }
        },
        watch: {
            async businessId (businessId) {
                if (businessId) {
                    await this.getMainlineModel()
                    await this.getBusinessTopo()
                    if (this.selectedHosts.length === 1) {
                        this.setSelectedModules()
                    }
                }
            },
            showIncrementOption (show) {
                this.increment = show
            }
        },
        methods: {
            ...mapActions('objectMainLineModule', [
                'searchMainlineObject',
                'getInstTopo',
                'getInternalTopo'
            ]),
            ...mapActions('hostRelation', [
                'transferHostToResourceModule',
                'transferHostToIdleModule',
                'transferHostToFaultModule',
                'transferHostModule'
            ]),
            getMainlineModel () {
                return this.searchMainlineObject({
                    config: {
                        requestId: 'get_searchMainlineObject'
                    }
                }).then(topoModel => {
                    this.topoModel = topoModel
                    return topoModel
                })
            },
            getBusinessTopo () {
                return Promise.all([
                    this.getInstTopo({
                        bizId: this.businessId,
                        config: {
                            requestId: `get_getInstTopo_${this.businessId}`
                        }
                    }),
                    this.getInternalTopo({
                        bizId: this.businessId,
                        config: {
                            requestId: `get_getInternalTopo_${this.businessId}`
                        }
                    })
                ]).then(([instTopo, internalTopo]) => {
                    const internalModule = (internalTopo.module || []).map(module => {
                        return {
                            'default': ['空闲机', 'idle machine'].includes(module.bk_module_name) ? 1 : 2,
                            'bk_obj_id': 'module',
                            'bk_inst_id': module.bk_module_id,
                            'bk_inst_name': module.bk_module_name
                        }
                    })
                    const treeData = [{
                        'default': 0,
                        'bk_obj_id': 'module',
                        'bk_inst_id': 'resource',
                        'bk_inst_name': this.$t('资源池')
                    }, {
                        ...instTopo[0],
                        child: [...internalModule, ...instTopo[0].child]
                    }]
                    this.$refs.topoTree.setData(treeData)
                    if (!this.hasTransferResourceAuth) {
                        this.$nextTick(() => {
                            this.$refs.topoTree.setDisabled('module-resource')
                        })
                    }
                })
            },
            setSelectedModules () {
                this.$nextTick(() => {
                    const modules = this.selectedHosts[0]['module']
                    const moduleIds = modules.map(module => {
                        return this.getTopoNodeId({
                            'bk_obj_id': 'module',
                            'bk_inst_id': module.bk_module_id
                        })
                    })
                    this.$refs.topoTree.setChecked(moduleIds, {
                        checked: true,
                        emitEvent: true,
                        beforeCheck: false
                    })
                })
            },
            getTopoNodeId (node) {
                return `${node['bk_obj_id']}-${node['bk_inst_id']}`
            },
            shouldShowCheckbox (data) {
                return data.bk_obj_id === 'module'
            },
            handleNodeCheck (checked) {
                this.selectedModules = checked.map(id => this.$refs.topoTree.getNodeById(id))
            },
            beforeNodeCheck (node) {
                let confirmResolver
                const asyncConfirm = new Promise(resolve => {
                    confirmResolver = resolve
                })
                const data = node.data
                const isSpecialNode = !!data.default || data.bk_inst_id === 'resource'
                const hasNormalNode = this.selectedModules.some(selectedNode => {
                    const selectedNodeData = selectedNode.data
                    return !selectedNodeData.default && selectedNodeData.bk_inst_id !== 'resource'
                })
                if (isSpecialNode && hasNormalNode) {
                    this.$bkInfo({
                        title: this.$t('转移确认', { target: data.bk_inst_name }),
                        confirmFn: () => {
                            this.$refs.topoTree.removeChecked({ emitEvent: true })
                            confirmResolver(true)
                        },
                        cancelFn: () => {
                            confirmResolver(false)
                        },
                        zIndex: 2000
                    })
                } else {
                    const specialNodes = this.selectedModules.filter(selectedNode => {
                        const selectedNodeData = selectedNode.data
                        return selectedNodeData.default || selectedNodeData.bk_inst_id === 'resource'
                    })
                    if (specialNodes.length && !specialNodes.includes(node)) {
                        this.$refs.topoTree.removeChecked({ emitEvent: true })
                    }
                    confirmResolver(true)
                }
                return asyncConfirm
            },
            handleNodeClick (node) {
                const isModule = node.data.bk_obj_id === 'module'
                if (isModule) {
                    this.$refs.topoTree.setChecked(node.id, {
                        checked: !node.checked,
                        emitEvent: true,
                        beforeCheck: true
                    })
                } else {
                    this.$refs.topoTree.setExpanded(node.id, {
                        expanded: !node.expanded,
                        emitEvent: false
                    })
                }
            },
            removeSelectedModule (node) {
                this.$refs.topoTree.setChecked(node.id, {
                    checked: false,
                    emitEvent: true,
                    beforeCheck: false
                })
            },
            getModulePath (node) {
                const data = node.data
                if (data.bk_inst_id === 'resource') {
                    return this.$t('主机资源池')
                }
                return node.parents.map(parent => parent.data.bk_inst_name).join('-')
            },
            handleTransfer () {
                const toSource = this.selectedModules.some(node => node.data.bk_inst_id === 'resource')
                const toIdle = this.selectedModules.some(node => node.data.default === 1)
                const toFault = this.selectedModules.some(node => node.data.default === 2)
                const transferConfig = {
                    requestId: 'transferHost'
                }
                let transferPromise
                if (toSource) {
                    transferPromise = this.transferToSource(transferConfig)
                } else if (toIdle) {
                    transferPromise = this.transferToIdle(transferConfig)
                } else if (toFault) {
                    transferPromise = this.transferToFault(transferConfig)
                } else {
                    transferPromise = this.transerToModules(transferConfig)
                }
                transferPromise.then(() => {
                    this.$success(this.$t('转移成功'))
                    this.$emit('on-success')
                })
            },
            transferToSource (config) {
                return this.transferHostToResourceModule({
                    params: {
                        'bk_biz_id': this.businessId,
                        'bk_host_id': this.hostIds
                    },
                    config
                })
            },
            transferToIdle (config) {
                return this.transferHostToIdleModule({
                    params: this.getTransferParams(),
                    config
                })
            },
            transferToFault (config) {
                return this.transferHostToFaultModule({
                    params: this.getTransferParams(),
                    config
                })
            },
            transerToModules (config) {
                return this.transferHostModule({
                    params: this.getTransferParams(),
                    config
                })
            },
            getTransferParams () {
                const increment = this.hostIds.length === 1 ? false : this.increment
                return {
                    'bk_biz_id': this.businessId,
                    'bk_host_id': this.hostIds,
                    'bk_module_id': this.selectedModules.map(node => node.data.bk_inst_id),
                    'is_increment': increment
                }
            },
            handleCancel () {
                this.$emit('on-cancel')
            },
            handleReceiveAuth (auth) {
                this.hasTransferResourceAuth = auth
            }
        }
    }
</script>

<style lang="scss" scoped>
    .transfer-layout {
        height: 540px;
        width: 720px;
        .columns-layout{
            width: 50%;
            height: calc(100% - 61px);
        }
    }
    .business-layout {
        border-right: 1px solid $cmdbBorderColor;
        border-bottom: 1px solid $cmdbBorderColor;
        height: 65px;
        &:before {
            display: inline-block;
            width: 0;
            height: 100%;
            content: '';
            font-size: 0;
            vertical-align: middle;
        }
        .business-label {
            display: inline-block;
            vertical-align: middle;
            padding: 0 25px;
        }
        .business-selector {
            display: inline-block;
            vertical-align: middle;

            width: 245px;
        }
    }
    .tree-layout {
        height: 415px;
        border-right: 1px solid $cmdbBorderColor;
    }
    .topo-tree{
        padding: 0 0 0 20px;
        height: 100%;
        @include scrollbar-y;
    }
    .node-info {
        .node-model-icon {
            width: 22px;
            height: 22px;
            line-height: 21px;
            text-align: center;
            font-style: normal;
            font-size: 12px;
            margin: 9px 4px 0 6px;
            border-radius: 50%;
            background-color: #c4c6cc;
            color: #fff;
            &.is-checked {
                background-color: #3a84ff;
            }
        }
    }
    .selected-layout {
        height: 65px;
        line-height: 65px;
        border-bottom: 1px solid $cmdbBorderColor;
        .selected-label {
            padding: 0 0 0 25px;
        }
    }
    .modules-layout {
        height: calc(100% - 65px);
        @include scrollbar-y;
    }
    .module-list {
        .module-item {
            height: 44px;
            line-height: 16px;
            padding: 6px 0 6px 25px;
            &:hover:not(.disabled){
                background-color: #e2efff;
                .module-info-name{
                    color: #498fe0;
                }
                .module-info-path{
                    color: #93bff5;
                }
            }
            .module-info{
                width: 250px;
            }
            .module-info-name{
                display: block;
                font-size: 14px;
                color: $cmdbTextColor;
                @include ellipsis;
            }
            .module-info-path{
                display: block;
                font-size: 12px;
                color: #b9bdc1;
                @include ellipsis;
            }
            .icon-close{
                cursor: pointer;
                color: #9196a1;
                font-size: 15px;
                margin: 10px 25px 0 0;
                &:hover{
                    color: #3c96ff;
                }
            }
        }
    }
    .options-layout {
        height: 61px;
        border-top: 1px solid $cmdbBorderColor;
        .content-middle {
            height: 100%;
            &:before {
                display: inline-block;
                width: 0;
                height: 100%;
                content: '';
                font-size: 0;
                vertical-align: middle;
            }
        }
    }
    .increment-layout {
        width: 500px;
        white-space: nowrap;
        padding-left: 20px;
        .cmdb-form-radio {
            max-width: 230px;
            vertical-align: middle;
            @include ellipsis;
            margin-right: 10px;
        }
    }
    .button-layout {
        .transfer-button {
            margin: 0 15px 0 0;
        }
    }
</style>
