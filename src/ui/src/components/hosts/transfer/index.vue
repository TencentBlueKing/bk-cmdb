<template>
    <div class="transfer-layout clearfix" v-bkloading="{isLoading: loading}">
        <div class="columns-layout fl">
            <div class="business-layout">
                <label class="business-label">{{$t('Common[\'业务\']')}}</label>
                <cmdb-business-selector class="business-selector" v-model="businessId" :disabled="true">
                </cmdb-business-selector>
            </div>
            <div class="tree-layout">
                <cmdb-tree ref="topoTree" class="topo-tree"
                    children-key="child"
                    :selectable="false"
                    :id-generator="getTopoNodeId"
                    :tree="tree.data"
                    :before-click="beforeNodeSelect"
                    @on-click="handleNodeClick">
                    <div class="tree-node clearfix" slot-scope="{node, state}"
                        :class="{
                            'tree-node-selected': state.selected,
                            'tree-node-leaf-module': node['bk_obj_id'] === 'module'
                        }">
                        <cmdb-form-bool class="topo-node-checkbox"
                            v-if="node['bk_obj_id'] === 'module'"
                            :checked="selectedModuleStates.includes(state)"
                            :true-value="true"
                            :false-value="false"
                            @click.stop
                            @change="handleNodeCheck(...arguments, node, state)"
                            >
                        </cmdb-form-bool>
                        <template v-if="[1, 2].includes(node.default)">
                            <i class='topo-node-icon topo-node-icon-internal icon-cc-host-free-pool' v-if="node.default === 1"></i>
                            <i class='topo-node-icon topo-node-icon-internal icon-cc-host-breakdown' v-else></i>
                        </template>
                        <i class="topo-node-icon topo-node-icon-text" v-else>{{node['bk_obj_name'][0]}}</i>
                        <span class="topo-node-text">{{node['bk_inst_name']}}</span>
                    </div>
                </cmdb-tree>
            </div>
        </div>
        <div class="columns-layout fl">
            <div class="selected-layout">
                <label class="selected-label">{{$t('Hosts["已选中模块"]')}}</label>
            </div>
            <div class="modules-layout">
                <ul class="module-list">
                    <li class="module-item clearfix"
                        v-for="(state, index) in selectedModuleStates" :key="index">
                        <div class="module-info fl">
                            <span class="module-info-name">{{state.node['bk_inst_name']}}</span>
                            <span class="module-info-path">{{getModulePath(state)}}</span>
                        </div>
                        <i class="bk-icon icon-close fr" @click="removeSelectedModule(state, index)"></i>
                    </li>
                </ul>
            </div>
        </div>
        <div v-pre class="clearfix"></div>
        <div class="options-layout clearfix">
            <div class="increment-layout content-middle fl" v-if="showIncrementOption">
                <label class="cmdb-form-radio cmdb-radio-small" for="increment" :title="$t('Hosts[\'增量更新\']')">
                    <input id="increment" type="radio" v-model="increment" :value="true">
                    <span class="cmdb-radio-text">{{$t('Hosts["增量更新"]')}}</span>
                </label>
                <label class="cmdb-form-radio cmdb-radio-small" for="replacement" :title="$t('Hosts[\'完全替换\']')">
                    <input id="replacement" type="radio" v-model="increment" :value="false">
                    <span class="cmdb-radio-text">{{$t('Hosts["完全替换"]')}}</span>
                </label>
            </div>
            <div class="button-layout content-middle fr">
                <bk-button class= "transfer-button" type="primary"
                    :disabled="!selectedModuleStates.length"
                    @click="handleTransfer">
                    {{$t('Common[\'确认转移\']')}}
                </bk-button>
                <bk-button class= "transfer-button" type="default" @click="handleCancel">{{$t('Common[\'取消\']')}}</bk-button>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            selectedHosts: {
                type: Array,
                required: true,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                businessId: '',
                topoModel: [],
                tree: {
                    data: []
                },
                selectedModuleStates: [],
                increment: true
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['business']),
            currentBusiness () {
                return this.business.find(item => item['bk_biz_id'] === this.businessId)
            },
            hostIds () {
                return this.selectedHosts.map(host => host['host']['bk_host_id'])
            },
            showIncrementOption () {
                const isMoreThanOne = this.selectedHosts.length > 1
                const hasSpecialModule = this.selectedModuleStates.some(({node}) => node['bk_inst_id'] === 'source' || [1, 2].includes(node.default))
                return !!this.selectedModuleStates.length && isMoreThanOne && !hasSpecialModule
            },
            loading () {
                const requestIds = [
                    'get_searchMainlineObject',
                    `get_getInstTopo_${this.businessId}`,
                    `get_getInternalTopo_${this.businessId}`,
                    'post_transferHost'
                ]
                return this.$loading(requestIds)
            }
        },
        watch: {
            async businessId (businessId) {
                if (businessId) {
                    await this.getMainlineModel()
                    await this.getBusinessTopo()
                    if (this.selectedHosts.length === 1) {
                        this.setSelectedModuleStates()
                    }
                }
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
                    requestId: 'get_searchMainlineObject',
                    fromCache: true
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
                            requestId: `get_getInstTopo_${this.businessId}`,
                            fromCache: true
                        }
                    }),
                    this.getInternalTopo({
                        bizId: this.businessId,
                        config: {
                            requestId: `get_getInternalTopo_${this.businessId}`,
                            fromCache: true
                        }
                    })
                ]).then(([instTopo, internalTopo]) => {
                    const moduleModel = this.getModelByObjId('module')
                    const internalModule = internalTopo.module.map(module => {
                        return {
                            'default': ['空闲机', 'idle machine'].includes(module['bk_module_name']) ? 1 : 2,
                            'bk_obj_id': 'module',
                            'bk_obj_name': moduleModel['bk_obj_name'],
                            'bk_inst_id': module['bk_module_id'],
                            'bk_inst_name': module['bk_module_name']
                        }
                    })
                    this.tree.data = [{
                        'default': 0,
                        'bk_obj_id': 'module',
                        'bk_obj_name': this.$t('HostResourcePool["资源池"]'),
                        'bk_inst_id': 'source',
                        'bk_inst_name': this.$t('HostResourcePool["资源池"]'),
                        'child': []
                    }, {
                        expanded: true,
                        ...instTopo[0],
                        child: [...internalModule, ...instTopo[0].child]
                    }]
                })
            },
            setSelectedModuleStates () {
                const modules = this.selectedHosts[0]['module']
                const selectedStates = []
                modules.forEach(module => {
                    const nodeId = this.getTopoNodeId({
                        'bk_obj_id': 'module',
                        'bk_inst_id': module['bk_module_id']
                    })
                    const state = this.$refs.topoTree.getStateById(nodeId)
                    if (state) {
                        selectedStates.push(state)
                    }
                })
                this.selectedModuleStates = selectedStates
            },
            getModelByObjId (id) {
                return this.topoModel.find(model => model['bk_obj_id'] === id)
            },
            getTopoNodeId (node) {
                return `${node['bk_obj_id']}-${node['bk_inst_id']}`
            },
            beforeNodeSelect (node, state) {
                let confirmResolver
                let confirmRejecter
                const asyncConfirm = new Promise((resolve, reject) => {
                    confirmResolver = resolve
                    confirmRejecter = reject
                })
                if (node['bk_obj_id'] !== 'module') {
                    confirmResolver(true)
                } else {
                    const isSpecialNode = !!node.default || node['bk_inst_id'] === 'source'
                    const hasNormalNode = this.selectedModuleStates.some(({node}) => {
                        return !node.default && node['bk_inst_id'] !== 'source'
                    })
                    const hasSpecialNode = this.selectedModuleStates.some(({node}) => {
                        return node.default || node['bk_inst_id'] === 'source'
                    })
                    if (isSpecialNode && hasNormalNode) {
                        this.$bkInfo({
                            title: this.$t('Common[\'转移确认\']', {target: node['bk_inst_name']}),
                            confirmFn: () => {
                                this.selectedModuleStates = []
                                confirmResolver(true)
                            },
                            cancelFn: () => {
                                confirmResolver(false)
                            }
                        })
                    } else {
                        if (hasSpecialNode && !this.selectedModuleStates.includes(state)) {
                            this.selectedModuleStates = []
                        }
                        confirmResolver(true)
                    }
                }
                return asyncConfirm
            },
            async handleNodeCheck (checked, vNode, node, state) {
                if (!checked) {
                    this.selectedModuleStates = this.selectedModuleStates.filter(moduleState => moduleState !== state)
                } else {
                    const confirm = await this.beforeNodeSelect(node, state)
                    if (confirm) {
                        if (!this.selectedModuleStates.includes(state)) {
                            this.selectedModuleStates.push(state)
                        }
                    } else {
                        vNode.localChecked = false
                    }
                }
            },
            handleNodeClick (node, state) {
                const isModule = node['bk_obj_id'] === 'module'
                const isExist = this.selectedModuleStates.some(selectedState => selectedState === state)
                if (isModule) {
                    if (isExist) {
                        this.selectedModuleStates = this.selectedModuleStates.filter(selectedState => selectedState !== state)
                    } else {
                        this.selectedModuleStates.push(state)
                    }
                }
            },
            removeSelectedModule (state, index) {
                this.selectedModuleStates.splice(index, 1)
                if (state.selected) {
                    state.selected = false
                }
            },
            getModulePath (state) {
                if (state.node['bk_inst_id'] === 'source') {
                    return this.$t('Common["主机资源池"]')
                }
                const currentBusiness = this.currentBusiness
                if ([1, 2].includes(state.node.default)) {
                    return `${currentBusiness['bk_biz_name']}-${state.node['bk_inst_name']}`
                }
                return `${currentBusiness['bk_biz_name']}-${state.parent.node['bk_inst_name']}`
            },
            handleTransfer () {
                const toSource = this.selectedModuleStates.some(({node}) => node['bk_inst_id'] === 'source')
                const toIdle = this.selectedModuleStates.some(({node}) => node.default === 1)
                const toFault = this.selectedModuleStates.some(({node}) => node.default === 2)
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
                    this.$success(this.$t('Common[\'转移成功\']'))
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
                let hasSpecialNode = this.selectedModuleStates.some(({node}) => [1, 2].includes(node.default))
                let increment = this.increment
                if (this.hasSpecialNode || this.hostIds.length === 1) {
                    increment = false
                }
                return {
                    'bk_biz_id': this.businessId,
                    'bk_host_id': this.hostIds,
                    'bk_module_id': this.selectedModuleStates.map(({node}) => node['bk_inst_id']),
                    'is_increment': increment
                }
            },
            handleCancel () {
                this.$emit('on-cancel')
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
        .tree-node {
            font-size: 0;
            &:hover{
                .topo-node-icon.topo-node-icon-text{
                    background-color: #50abff;
                }
                .topo-node-icon.topo-node-icon-internal{
                    color: #50abff;
                }
            }
            &.tree-node-selected{
                .topo-node-icon.topo-node-icon-text{
                    background-color: #498fe0;
                }
                .topo-node-icon.topo-node-icon-internal{
                    color: #ffb400;
                }
            }
            &.tree-node-leaf-module {
                margin: 0 0 0 -2px !important;
                padding-left: 0;
            }
        }
        .topo-node-checkbox {
            position: relative;
            margin: 0 10px 0 0;
            transform: scale(0.888);
            background-color: #fff;
            z-index: 2;
        }
        .topo-node-icon{
            display: inline-block;
            vertical-align: middle;
            width: 18px;
            height: 18px;
            line-height: 16px;
            font-size: 12px;
            text-align: center;
            color: #fff;
            font-style: normal;
            background-color: #c3cdd7;
            border-radius: 50%;
            &.topo-node-icon-internal{
                font-size: 16px;
                color: $cmdbTextColor;
                background-color: transparent;
            }
        }
        .topo-node-text{
            display: inline-block;
            vertical-align: middle;
            padding: 0 0 0 8px;
            font-size: 14px;
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