<template>
    <div class="transfer-menu">
        <bk-dropdown-menu
            trigger="click"
            font-size="medium"
            :disabled="disabled"
            @show="handleMenuToggle(true)"
            @hide="handleMenuToggle(false)">
            <div class="dropdown-trigger-btn" style="padding-left: 19px;" slot="dropdown-trigger">
                <span>{{$t('转移到')}}</span>
                <i :class="['bk-icon icon-angle-down', { 'icon-flip': isShow }]"></i>
            </div>
            <ul class="bk-dropdown-list" slot="dropdown-content">
                <li><a href="javascript:;" @click="transferToIdleModule">{{$t('空闲模块')}}</a></li>
                <li><a href="javascript:;" @click="transferToBizModule">{{$t('业务模块')}}</a></li>
                <li><a href="javascript:;" @click="transferToResourcePool">{{$t('主机池')}}</a></li>
            </ul>
        </bk-dropdown-menu>
        <cmdb-dialog v-model="dialog.show" :width="dialog.width" :height="dialog.height">
            <component
                :is="dialog.component"
                v-bind="dialog.props"
                @cancel="handleDialogCancel"
                @confirm="handleDialogConfirm">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import HostStore from './host-store'
    import ModuleSelector from '../../business-topology/host/module-selector'
    import MoveToResourceConfirm from '../../business-topology/host/move-to-resource-confirm'
    import { MENU_BUSINESS_TRANSFER_HOST } from '@/dictionary/menu-symbol'
    import RouterQuery from '@/router/query'
    export default {
        components: {
            [ModuleSelector.name]: ModuleSelector,
            [MoveToResourceConfirm.name]: MoveToResourceConfirm
        },
        data () {
            return {
                isShow: false,
                dialog: {
                    width: 830,
                    height: 600,
                    show: false,
                    component: null,
                    props: {}
                },
                request: {
                    moveToIdleModule: Symbol('moveToIdleModule'),
                    moveToResource: Symbol('moveToResource')
                }
            }
        },
        computed: {
            disabled () {
                return !HostStore.isSelected
            }
        },
        methods: {
            handleMenuToggle (isShow) {
                this.isShow = isShow
            },
            validateSameBiz () {
                if (!HostStore.isSameBiz) {
                    this.$error(this.$t('仅支持对相同业务下的主机进行操作'))
                    return false
                }
                return true
            },
            transferToIdleModule () {
                const valid = this.validateSameBiz()
                if (!valid) {
                    return false
                }
                if (HostStore.isAllResourceHost) {
                    this.$error(this.$t('仅支持对业务下的主机进行操作'))
                    return false
                }
                const props = {
                    moduleType: 'idle',
                    business: HostStore.uniqueBusiness,
                    title: this.$t('转移主机到空闲模块')
                }
                this.dialog.props = props
                this.dialog.width = 830
                this.dialog.height = 600
                this.dialog.component = ModuleSelector.name
                this.dialog.show = true
            },
            transferToBizModule () {
                const valid = this.validateSameBiz()
                if (!valid) {
                    return false
                }
                if (HostStore.isAllResourceHost) {
                    this.$error(this.$t('仅支持对业务下的主机进行操作'))
                    return false
                }
                const props = {
                    moduleType: 'business',
                    business: HostStore.uniqueBusiness,
                    title: this.$t('转移主机到业务模块')
                }
                const selection = HostStore.getSelected()
                const firstSelectionModules = selection[0].module.map(module => module.bk_module_id).sort()
                const firstSelectionModulesStr = firstSelectionModules.join(',')
                const allSame = selection.slice(1).every(item => {
                    const modules = item.module.map(module => module.bk_module_id).sort().join(',')
                    return modules === firstSelectionModulesStr
                })
                if (allSame) {
                    props.previousModules = firstSelectionModules
                }
                this.dialog.props = props
                this.dialog.width = 830
                this.dialog.height = 600
                this.dialog.component = ModuleSelector.name
                this.dialog.show = true
            },
            transferToResourcePool () {
                const isSameBiz = this.validateSameBiz()
                if (!isSameBiz) {
                    return false
                }
                if (HostStore.isAllResourceHost) {
                    this.$error('所选主机已在主机池中')
                    return false
                }
                const isAllIdleModule = HostStore.isAllIdleModule
                if (!isAllIdleModule) {
                    this.$error(this.$t('仅支持对空闲机模块下的主机进行操作'))
                    return false
                }
                const [bizId] = HostStore.bizSet
                this.dialog.props = {
                    count: HostStore.getSelected().length,
                    bizId: bizId
                }
                this.dialog.width = 400
                this.dialog.height = 231
                this.dialog.component = MoveToResourceConfirm.name
                this.dialog.show = true
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDialogConfirm () {
                this.dialog.show = false
                if (this.dialog.component === ModuleSelector.name) {
                    if (this.dialog.props.moduleType === 'idle') {
                        if (HostStore.isAllIdleSet) {
                            this.transferDirectly(...arguments)
                        } else {
                            this.gotoTransferPage(...arguments)
                        }
                    } else {
                        this.gotoTransferPage(...arguments)
                    }
                } else if (this.dialog.component === MoveToResourceConfirm.name) {
                    this.moveHostToResource(...arguments)
                }
            },
            async transferDirectly (modules) {
                try {
                    const bizId = HostStore.uniqueBusiness.bk_biz_id
                    const internalModule = modules[0]
                    await this.$http.post(
                        `host/transfer_with_auto_clear_service_instance/bk_biz_id/${bizId}`, {
                            bk_host_ids: HostStore.getSelected().map(data => data.host.bk_host_id),
                            default_internal_module: internalModule.data.bk_inst_id,
                            remove_from_node: {
                                bk_inst_id: bizId,
                                bk_obj_id: 'biz'
                            }
                        }, {
                            requestId: this.request.moveToIdleModule
                        }
                    )
                    HostStore.clear()
                    this.$success('转移成功')
                    RouterQuery.set({
                        _t: Date.now(),
                        page: 1
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            gotoTransferPage (modules) {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_TRANSFER_HOST,
                    params: {
                        bizId: HostStore.uniqueBusiness.bk_biz_id,
                        type: this.dialog.props.moduleType
                    },
                    query: {
                        targetModules: modules.map(node => node.data.bk_inst_id).join(','),
                        resources: HostStore.getSelected().map(item => item.host.bk_host_id).join(',')
                    },
                    history: true
                })
                HostStore.clear()
            },
            async moveHostToResource (directoryId) {
                try {
                    await this.$store.dispatch('hostRelation/transferHostToResourceModule', {
                        params: {
                            bk_biz_id: HostStore.uniqueBusiness.bk_biz_id,
                            bk_host_id: HostStore.getSelected().map(item => item.host.bk_host_id),
                            bk_module_id: directoryId
                        },
                        config: {
                            requestId: this.request.moveToResource
                        }
                    })
                    HostStore.clear()
                    this.$success('转移成功')
                    RouterQuery.set({
                        _t: Date.now(),
                        page: 1
                    })
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .transfer-menu {
        display: inline-block;
    }
    .dropdown-trigger-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #c4c6cc;
        height: 32px;
        min-width: 68px;
        border-radius: 2px;
        padding: 0 15px;
        color: #63656E;
        font-size: 14px;
    }
    .dropdown-trigger-btn.bk-icon {
        font-size: 18px;
    }
    .dropdown-trigger-btn .bk-icon {
        font-size: 22px;
    }
    .dropdown-trigger-btn:hover {
        cursor: pointer;
        border-color: #979ba5;
    }
</style>
