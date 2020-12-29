<template>
    <div class="module-selector-with-tab">
        <bk-tab :active.sync="tab.active" type="border-card">
            <bk-tab-panel
                v-for="(panel, index) in availableTabList"
                v-bind="panel.props"
                render-directive="if"
                :key="index">
                <div class="tab-content">
                    <div class="content-container" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
                        <cmdb-auth class="auth-component" v-if="['idle', 'business'].includes(panel.props.name)"
                            :ignore="Object.keys(authorized).includes(panel.props.name)"
                            :auth="auth"
                            @update-auth="handleUpdateAuth(...arguments, panel.props.name)">
                        </cmdb-auth>
                        <component v-if="authorized[panel.props.name] !== false"
                            class="selector-component"
                            :is="panel.component.name"
                            v-bind="panel.component.props"
                            @cancel="handleCancel"
                            @confirm="handleConfirm">
                        </component>
                        <no-permission class="no-permission-container" v-else :permission="permission" @cancel="handleCancel" />
                    </div>
                </div>
            </bk-tab-panel>
        </bk-tab>
    </div>
</template>

<script>
    import { translateAuth } from '@/setup/permission'
    import { AuthRequestId } from '@/components/ui/auth/auth-queue.js'
    import ModuleSelector from './module-selector.vue'
    import AcrossBusinessModuleSelector from './across-business-module-selector.vue'
    import NoPermission from './no-permission.vue'
    export default {
        name: 'module-selector-with-tab',
        components: {
            NoPermission,
            [ModuleSelector.name]: ModuleSelector,
            [AcrossBusinessModuleSelector.name]: AcrossBusinessModuleSelector
        },
        props: {
            modules: {
                type: Array,
                default () {
                    return []
                }
            },
            business: {
                type: Object,
                default () {
                    return {}
                }
            },
            confirmLoading: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                tab: {
                    list: [
                        {
                            props: {
                                name: 'idle',
                                label: this.$t('转移到空闲模块'),
                                visible: true
                            },
                            component: {
                                name: ModuleSelector.name,
                                props: {
                                    moduleType: 'idle',
                                    business: {},
                                    confirmText: '',
                                    confirmLoading: false
                                }
                            }
                        },
                        {
                            props: {
                                name: 'business',
                                label: this.$t('转移到业务模块'),
                                visible: true
                            },
                            component: {
                                name: ModuleSelector.name,
                                props: {
                                    moduleType: 'business',
                                    business: {},
                                    confirmLoading: false
                                }
                            }
                        },
                        {
                            props: {
                                name: 'acrossBusiness',
                                label: this.$t('转移到其他业务模块'),
                                visible: true
                            },
                            component: {
                                name: AcrossBusinessModuleSelector.name,
                                props: {
                                    business: {},
                                    confirmLoading: false
                                }
                            }
                        }
                    ],
                    active: 'idle'
                },
                authorized: {},
                request: {
                    auth: AuthRequestId
                }
            }
        },
        computed: {
            bizId () {
                return this.business.bk_biz_id
            },
            isIdleSetModules () {
                return this.modules.every(module => module.default >= 1)
            },
            availableTabList () {
                const availableTabList = []
                this.tab.list.forEach(tab => {
                    tab.component.props.business = this.business
                    if (tab.props.name !== 'acrossBusiness') {
                        const defaultChecked = this.modules.map(module => module.bk_module_id)
                        const firstSelectionModules = this.modules.map(module => module.bk_module_id).sort()
                        tab.component.props.previousModules = firstSelectionModules
                        tab.component.props.defaultChecked = defaultChecked
                        tab.component.props.confirmText = tab.props.name === 'idle' && this.isIdleSetModules ? this.$t('确定') : ''
                        availableTabList.push(tab)
                    } else if (this.isIdleSetModules) {
                        availableTabList.push(tab)
                    }
                })
                return availableTabList
            },
            activeTab () {
                return this.availableTabList.find(tab => tab.props.name === this.tab.active)
            },
            auth () {
                return [
                    { type: this.$OPERATION.C_SERVICE_INSTANCE, relation: [this.bizId] },
                    { type: this.$OPERATION.U_SERVICE_INSTANCE, relation: [this.bizId] },
                    { type: this.$OPERATION.D_SERVICE_INSTANCE, relation: [this.bizId] }
                ]
            },
            permission () {
                return translateAuth(this.auth)
            }
        },
        watch: {
            confirmLoading (value) {
                this.activeTab.component.props.confirmLoading = value
            }
        },
        methods: {
            handleCancel () {
                this.$emit('cancel')
            },
            handleConfirm () {
                const currentTab = this.activeTab
                const tab = { tabName: currentTab.props.name, moduleType: currentTab.component.props.moduleType }
                this.$emit('confirm', tab, ...arguments)
            },
            handleUpdateAuth (isAuthorized, panel) {
                // 已鉴权则不再更新，配合auth组件ignore，在切换tab时不重复鉴权
                if (!this.authorized.hasOwnProperty(panel)) {
                    this.$set(this.authorized, panel, isAuthorized)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .module-selector-with-tab {
        height: var(--height); // defined in dialog el

        .tab-content,
        .selector-component,
        .content-container,
        .no-permission-container {
            height: 100%;
        }

        .auth-component {
            display: none;
        }

        /deep/ .bk-tab {
            height: 100%;
            .bk-tab-header {
                padding: 0;
                height: 43px;
                background-image: linear-gradient(transparent 41px,#dcdee5 0);
                .bk-tab-label-list {
                    height: 42px;
                    .bk-tab-label-item {
                        line-height: 42px;
                        min-width: auto;
                        &.active {
                            color: #313238;
                            background-color: #fff;
                        }
                        &:not(.is-disabled):hover {
                            color: #313238;
                        }
                    }
                }
            }
            .bk-tab-header-setting {
                height: 42px;
                line-height: 42px;
            }
            .bk-tab-section {
                padding: 0;
                height: calc(100% - 43px);
                overflow: visible;
                .bk-tab-content {
                    height: 100%;
                }
            }
        }
    }
</style>
