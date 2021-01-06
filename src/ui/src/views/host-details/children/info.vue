<template>
    <div class="info">
        <div class="info-basic">
            <i :class="['info-icon', model.bk_obj_icon]"></i>
            <span class="info-ip">{{hostIp}}</span>
            <span class="info-area">{{cloudArea}}</span>
        </div>
        <div class="info-topology">
            <div class="topology-label">
                <span>{{$t('所属拓扑')}}</span>
                <span v-if="topologyList.length > 1" v-bk-tooltips="{
                    content: $t(isSingleColumn ? '切换双列显示' : '切换单列显示'),
                    interactive: false
                }">
                    <i class="topology-toggle icon-cc-single-column" v-if="isSingleColumn" @click="toggleDisplayType"></i>
                    <i class="topology-toggle icon-cc-double-column" v-else @click="toggleDisplayType"></i>
                </span>
                <span v-pre style="padding: 0 5px;">:</span>
            </div>
            <ul class="topology-list"
                :class="{ 'is-single-column': isSingleColumn }"
                :style="getListStyle(topologyList)">
                <li :class="['topology-item', { 'is-internal': item.isInternal }]"
                    v-for="(item, index) in topologyList"
                    :key="index">
                    <span class="topology-path" v-bk-overflow-tips @click="handlePathClick(item)">{{item.path}}</span>
                    <cmdb-auth :auth="[
                        { type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] },
                        { type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] },
                        { type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }
                    ]">
                        <i class="topology-remove-trigger icon-cc-tips-close"
                            v-if="!item.isInternal"
                            v-bk-tooltips="{ content: $t('从该模块移除'), interactive: false }"
                            @click="handleRemove(item.id)">
                        </i>
                    </cmdb-auth>
                </li>
            </ul>
            <a class="action-btn view-all"
                href="javascript:void(0)"
                v-if="showMore"
                @click="viewAll">
                {{$t('更多信息')}}
                <i class="bk-icon icon-angle-down" :class="{ 'is-all-show': showAll }"></i>
            </a>
            <a class="action-btn change-topology" v-if="isBusinessHost"
                href="javascript:void(0)"
                @click="handleEditTopo">
                {{$t('修改')}}
                <i class="icon icon-cc-edit-shape"></i>
            </a>
        </div>
        <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height">
            <component
                :is="dialog.component"
                v-bind="dialog.componentProps"
                :confirm-loading="confirmLoading"
                @cancel="handleDialogCancel"
                @confirm="handleDialogConfirm">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import { mapState, mapGetters } from 'vuex'
    import {
        MENU_BUSINESS_TRANSFER_HOST,
        MENU_BUSINESS_HOST_AND_SERVICE,
        MENU_RESOURCE_BUSINESS_HOST_DETAILS,
        MENU_RESOURCE_HOST,
        MENU_BUSINESS_HOST_DETAILS
    } from '@/dictionary/menu-symbol'
    import ModuleSelectorWithTab from '@/views/business-topology/host/module-selector-with-tab.vue'
    export default {
        name: 'cmdb-host-info',
        components: {
            [ModuleSelectorWithTab.name]: ModuleSelectorWithTab
        },
        data () {
            return {
                displayType: window.localStorage.getItem('host_topology_display_type') || 'double',
                showAll: false,
                topoNodesPath: [],
                dialog: {
                    show: false,
                    component: null,
                    componentProps: {},
                    width: 828,
                    height: 600
                },
                request: {
                    moveToIdleModule: Symbol('moveToIdleModule')
                }
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            ...mapGetters('hostDetails', ['isBusinessHost']),
            business () {
                const biz = this.info.biz || []
                return biz[0]
            },
            bizId () {
                return this.business.bk_biz_id
            },
            isSingleColumn () {
                return this.displayType === 'single'
            },
            host () {
                return this.info.host || {}
            },
            modules () {
                return this.info.module || []
            },
            hostIp () {
                if (Object.keys(this.host).length) {
                    const hostList = this.host.bk_host_innerip.split(',')
                    const host = hostList.length > 1 ? `${hostList[0]}...` : hostList[0]
                    return host
                } else {
                    return ''
                }
            },
            cloudArea () {
                return (this.host.bk_cloud_id || []).map(cloud => {
                    return `${this.$t('云区域')}：${cloud.bk_inst_name} (ID：${cloud.bk_inst_id})`
                }).join('\n')
            },
            topologyList () {
                const modules = this.info.module || []
                return this.topoNodesPath.map(item => {
                    const instId = item.topo_node.bk_inst_id
                    const module = modules.find(module => module.bk_module_id === instId)
                    return {
                        id: instId,
                        path: item.topo_path.reverse().map(node => node.bk_inst_name).join(' / '),
                        isInternal: module && module.default !== 0
                    }
                }).sort((itemA, itemB) => {
                    return itemA.path.localeCompare(itemB.path, 'zh-Hans-CN', { sensitivity: 'accent' })
                })
            },
            showMore () {
                if (this.isSingleColumn) {
                    return this.topologyList.length > 1
                }
                return this.topologyList.length > 2
            },
            model () {
                return this.$store.getters['objectModelClassify/getModelById']('host')
            },
            confirmLoading () {
                return this.$loading(Object.values(this.request))
            }
        },
        watch: {
            async info () {
                await this.getModulePathInfo()
            }
        },
        methods: {
            async getModulePathInfo () {
                try {
                    const modules = this.info.module || []
                    const biz = this.info.biz || []
                    const result = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
                        bizId: biz[0].bk_biz_id,
                        params: {
                            topo_nodes: modules.map(module => ({ bk_obj_id: 'module', bk_inst_id: module.bk_module_id }))
                        }
                    })
                    this.topoNodesPath = result.nodes || []
                } catch (e) {
                    console.error(e)
                    this.topoNodesPath = []
                }
            },
            viewAll () {
                this.showAll = !this.showAll
                this.$emit('info-toggle')
            },
            getListStyle (items) {
                const itemHeight = 21
                const itemMargin = 9
                const length = this.isSingleColumn ? items.length : Math.ceil(items.length / 2)
                return {
                    height: (this.showAll ? length : 1) * (itemHeight + itemMargin) + 'px',
                    flex: (!this.isSingleColumn && items.length === 1) ? 'none' : ''
                }
            },
            toggleDisplayType () {
                this.displayType = this.displayType === 'single' ? 'double' : 'single'
                this.$emit('info-toggle')
                window.localStorage.setItem('host_topology_display_type', this.displayType)
            },
            handlePathClick (item) {
                if (this.isBusinessHost) {
                    this.$routerActions.open({
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        params: {
                            bizId: this.bizId
                        },
                        query: {
                            node: `module-${item.id}`
                        }
                    })
                } else {
                    const modules = this.info.module || []
                    this.$routerActions.open({
                        name: MENU_RESOURCE_HOST,
                        params: {
                            bizId: this.bizId
                        },
                        query: {
                            scope: '1',
                            directory: modules[0].bk_module_id
                        }
                    })
                }
            },
            handleRemove (moduleId) {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        bizId: this.bizId,
                        type: 'remove',
                        module: moduleId
                    },
                    query: {
                        sourceModel: 'module',
                        sourceId: moduleId,
                        resources: this.$route.params.id
                    },
                    history: true
                })
            },
            handleEditTopo () {
                this.dialog.component = ModuleSelectorWithTab.name
                this.dialog.componentProps.modules = this.modules
                this.dialog.componentProps.business = this.business
                this.dialog.show = true
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDialogConfirm () {
                const [tab, ...params] = [...arguments]
                const { tabName, moduleType } = tab
                if (tabName !== 'acrossBusiness') {
                    if (moduleType === 'idle') {
                        const isAllIdleSetHost = this.modules.every(module => module.default !== 0)
                        if (isAllIdleSetHost) {
                            this.transferDirectly(...params)
                        } else {
                            this.gotoTransferPage(...params, moduleType)
                        }
                    } else {
                        this.gotoTransferPage(...params, moduleType)
                    }
                } else {
                    this.moveHostToOtherBusiness(...params)
                }
            },
            async transferDirectly (modules) {
                try {
                    const internalModule = modules[0]
                    await this.$http.post(
                        `host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}`, {
                            bk_host_ids: [this.host.bk_host_id],
                            default_internal_module: internalModule.data.bk_inst_id,
                            remove_from_node: {
                                bk_inst_id: this.bizId,
                                bk_obj_id: 'biz'
                            }
                        }, {
                            requestId: this.request.moveToIdleModule
                        }
                    )
                    this.dialog.show = false
                    this.$success('转移成功')
                    this.$emit('change')
                } catch (e) {
                    console.error(e)
                }
            },
            gotoTransferPage (modules, moduleType) {
                const query = {
                    sourceModel: 'biz',
                    sourceId: this.bizId,
                    targetModules: modules.map(node => node.data.bk_inst_id).join(','),
                    resources: [this.host.bk_host_id].join(','),
                    single: 1
                }
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        bizId: this.bizId,
                        type: moduleType
                    },
                    query: query,
                    history: true
                })
            },
            async moveHostToOtherBusiness (modules, targetBizId) {
                try {
                    const [targetModule] = modules
                    await this.$http.post('hosts/modules/across/biz', {
                        src_bk_biz_id: this.bizId,
                        dst_bk_biz_id: targetBizId,
                        bk_host_id: [this.host.bk_host_id],
                        bk_module_id: targetModule.data.bk_inst_id
                    })

                    this.dialog.show = false
                    this.$success('转移成功')

                    // 跳转刷新
                    const routeParams = {
                        id: this.host.bk_host_id
                    }
                    const routeName = this.$route.name
                    if (routeName === MENU_RESOURCE_BUSINESS_HOST_DETAILS) {
                        routeParams.business = targetBizId
                    } else if (routeName === MENU_BUSINESS_HOST_DETAILS) {
                        routeParams.bizId = targetBizId
                    }
                    this.$routerActions.redirect({
                        name: routeName,
                        params: {
                            ...this.$route.params,
                            ...routeParams
                        },
                        reload: false
                    })
                } catch (error) {
                    console.error(error)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info {
        max-height: 450px;
        padding: 11px 0 2px 24px;
        background:rgba(235, 244, 255, .6);
        border-bottom: 1px solid #dcdee5;
        @include scrollbar-y;
    }
    .info-basic {
        font-size: 0;
        .info-icon {
            display: inline-block;
            width: 38px;
            height: 38px;
            margin: 0 11px 0 0;
            border: 1px solid #dde4eb;
            border-radius: 50%;
            background-color: #fff;
            vertical-align: middle;
            line-height: 38px;
            text-align: center;
            font-size: 18px;
            color: #3a84ff;
        }
        .info-ip {
            display: inline-block;
            vertical-align: middle;
            line-height: 38px;
            font-size: 16px;
            font-weight: bold;
            color: #333948;
        }
        .info-area {
            display: inline-block;
            vertical-align: middle;
            height: 18px;
            margin-left: 10px;
            padding: 0 5px;
            line-height: 16px;
            font-size: 12px;
            color: #979BA5;
            border: 1px solid #C4C6CC;
            border-radius: 2px;
            background-color: #fff;
        }
    }
    .info-topology {
        line-height: 19px;
        display: flex;
        .topology-label {
            display: flex;
            align-items: center;
            align-self: baseline;
            padding: 0 0 0 50px;
            font-size: 14px;
            font-weight: bold;
            line-height: 20px;
            .topology-toggle {
                font-size: 16px;
                margin: 0 0 0 5px;
                cursor: pointer;
                &:hover {
                    opacity: .75;
                }
            }
        }
        .topology-list {
            flex: 1;
        }
        .action-btn {
            align-self: flex-start;
            margin: 0 14px;
            font-size: 12px;
            color: #007eff;
        }
        .view-all {
            .bk-icon {
                display: inline-block;
                vertical-align: -1px;
                font-size: 20px;
                margin-left: -4px;
                transition: transform .2s linear;
                &.is-all-show {
                    transform: rotate(-180deg);
                }
            }
        }
        .change-topology {
            .icon-cc-edit-shape {
                font-size: 14px;
            }
        }
    }
    .topology-list {
        display: flex;
        flex-wrap: wrap;
        overflow: hidden;
        color: #63656e;
        will-change: height;
        transition: height .2s ease-in;
        max-width: 700px;
        &.is-single-column {
            max-width: 850px;
            display: inline-block;
            flex: none;
            .topology-item {
                width: auto;
            }
        }
        .topology-item {
            flex: 0 1 50%;
            width: 50%;
            height: 20px;
            font-size: 0px;
            margin: 0 0 9px 0;
            padding: 0 15px 0 0;
            line-height: 20px;
            &:only-child {
                flex: 1 1 50%;
            }
            &:hover {
                .topology-remove-trigger {
                    opacity: 1;
                }
            }
            .topology-path {
                display: inline-block;
                vertical-align: middle;
                font-size: 14px;
                max-width: calc(100% - 30px);
                cursor: pointer;
                @include ellipsis;
                &:hover {
                    color: $primaryColor;
                }
            }
            .topology-remove-trigger {
                opacity: 0;
                font-size: 20px;
                cursor: pointer;
                margin: 0 0 0 10px;
                color: $textColor;
                transform: scale(.5);
                &:hover {
                    color: $primaryColor;
                }
            }
            &.is-internal {
                .topology-path {
                    max-width: 100%;
                }
            }
        }
    }
    .topology-list.right-list {
        margin: 0 0 0 90px;
    }
</style>
