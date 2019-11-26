<template>
    <div class="create-layout clearfix" v-bkloading="{ isLoading: $loading() }">
        <label class="create-label fl">{{$t('添加主机')}}</label>
        <div class="create-hosts">
            <div>
                <bk-button class="select-host-button" theme="default"
                    @click="handleSelectHost">
                    <i class="bk-icon icon-plus"></i>
                    {{$t('添加主机')}}
                </bk-button>
                <i18n class="select-host-count" path="已选择N台主机" v-show="hosts.length">
                    <span place="count" class="count-number">{{hosts.length}}</span>
                </i18n>
            </div>
            <div class="create-tables">
                <service-instance-table class="service-instance-table"
                    v-for="(item, index) in hosts"
                    ref="serviceInstanceTable"
                    deletable
                    expanded
                    :key="item.host.bk_host_id"
                    :index="index"
                    :id="item.host.bk_host_id"
                    :name="item.host.bk_host_innerip"
                    :source-processes="sourceProcesses"
                    :templates="templates"
                    :addible="!withTemplate"
                    @delete-instance="handleDeleteInstance">
                </service-instance-table>
            </div>
            <div class="buttons">
                <cmdb-auth class="mr5" :auth="$authResources({ type: $OPERATION.C_SERVICE_INSTANCE })">
                    <bk-button slot-scope="{ disabled }"
                        theme="primary"
                        :disabled="!hosts.length || disabled"
                        @click="handleConfirm">
                        {{$t('确定')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button @click="handleBackToModule">{{$t('取消')}}</bk-button>
            </div>
        </div>
        <cmdb-dialog v-model="dialog.show" v-bind="dialog.props">
            <component
                :is="dialog.component"
                :confirm-text="$t('确定')"
                v-bind="dialog.componentProps"
                @confirm="handleDialogConfirm"
                @cancel="handleDialogCancel">
            </component>
        </cmdb-dialog>
    </div>
</template>

<script>
    import HostSelector from '@/views/business-topology/host/host-selector.vue'
    import serviceInstanceTable from '@/components/service/instance-table.vue'
    import { MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import { mapGetters } from 'vuex'
    export default {
        name: 'create-service-instance',
        components: {
            [HostSelector.name]: HostSelector,
            serviceInstanceTable
        },
        data () {
            return {
                seed: 0,
                hostSelectorVisible: false,
                moduleInstance: {},
                hosts: [],
                templates: [],
                dialog: {
                    show: false,
                    props: {
                        width: 850,
                        height: 460,
                        showCloseIcon: false
                    },
                    component: null,
                    componentProps: {}
                },
                request: {
                    hostInfo: Symbol('hostInfo')
                }
            }
        },
        computed: {
            ...mapGetters('businessHost', ['getDefaultSearchCondition']),
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            moduleId () {
                return parseInt(this.$route.params.moduleId)
            },
            setId () {
                return parseInt(this.$route.params.setId)
            },
            withTemplate () {
                return this.moduleInstance.service_template_id
            },
            sourceProcesses () {
                return this.templates.map(template => {
                    const value = {}
                    const ip = ['127.0.0.1', '0.0.0.0']
                    Object.keys(template.property).forEach(key => {
                        if (key === 'bind_ip') {
                            value[key] = ip[template.property[key].value - 1]
                        } else {
                            value[key] = template.property[key].value
                        }
                    })
                    return value
                })
            }
        },
        watch: {
            withTemplate (withTemplate) {
                if (withTemplate) {
                    this.getTemplate()
                }
            }
        },
        async created () {
            this.$store.commit('setBreadcrumbs', [{
                label: this.$t('服务拓扑'),
                route: {
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: 'module-' + this.$route.params.moduleId
                    }
                }
            }, {
                label: this.$route.query.title
            }])
            await this.getModuleInstance()
            this.initSelectedHost()
        },
        methods: {
            async getModuleInstance () {
                try {
                    const data = await this.$store.dispatch('objectModule/searchModule', {
                        bizId: this.business,
                        setId: this.setId,
                        params: {
                            page: { start: 0, limit: 1 },
                            fields: [],
                            condition: {
                                bk_module_id: this.moduleId,
                                bk_supplier_account: this.$store.getters.supplierAccount
                            }
                        },
                        config: {
                            requestId: 'getModuleInstance'
                        }
                    })
                    if (!data.count) {
                        this.$router.replace({ name: '404' })
                    } else {
                        this.moduleInstance = data.info[0]
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            async initSelectedHost () {
                try {
                    const resources = this.$route.query.resources
                    if (resources) {
                        const hostIds = resources.split(',').map(id => Number(id))
                        const result = await this.getHostInfo(hostIds)
                        this.hosts = result.info
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            getHostInfo (hostIds) {
                const params = {
                    bk_biz_id: this.business,
                    ip: { data: [], exact: 0, flag: 'bk_host_innerip|bk_host_outerip' },
                    page: {},
                    condition: this.getDefaultSearchCondition()
                }
                const hostCondition = params.condition.find(target => target.bk_obj_id === 'host')
                hostCondition.condition.push({
                    field: 'bk_host_id',
                    operator: '$in',
                    value: hostIds
                })
                return this.$store.dispatch('hostSearch/searchHost', {
                    params,
                    config: {
                        requestId: this.request.hostInfo
                    }
                })
            },
            async getTemplate () {
                try {
                    const data = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: this.$injectMetadata({
                            service_template_id: this.moduleInstance.service_template_id,
                            page: { sort: 'id' }
                        }, { injectBizId: true }),
                        config: {
                            requestId: 'getBatchProcessTemplate'
                        }
                    })
                    this.templates = data.info
                } catch (e) {
                    console.error(e)
                }
            },
            handleSelectHost () {
                this.dialog.componentProps.exist = this.hosts
                this.dialog.component = HostSelector.name
                this.dialog.show = true
            },
            handleDialogConfirm (selected) {
                this.hosts = selected
                this.dialog.show = false
            },
            handleDialogCancel () {
                this.dialog.show = false
            },
            handleDeleteInstance (index) {
                this.hosts.splice(index, 1)
            },
            async handleConfirm () {
                try {
                    const serviceInstanceTables = this.$refs.serviceInstanceTable
                    if (this.withTemplate) {
                        await this.$store.dispatch('serviceInstance/createProcServiceInstanceByTemplate', {
                            params: this.$injectMetadata({
                                name: this.moduleInstance.bk_module_name,
                                bk_module_id: this.moduleId,
                                instances: serviceInstanceTables.map(table => {
                                    return {
                                        bk_host_id: table.id,
                                        processes: table.processList.map((item, index) => {
                                            return {
                                                process_info: item,
                                                process_template_id: table.templates[index].id
                                            }
                                        })
                                    }
                                })
                            }, { injectBizId: true })
                        })
                    } else {
                        await this.$store.dispatch('serviceInstance/createProcServiceInstanceWithRaw', {
                            params: this.$injectMetadata({
                                name: this.moduleInstance.bk_module_name,
                                bk_module_id: this.moduleId,
                                instances: serviceInstanceTables.map(table => {
                                    return {
                                        bk_host_id: table.id,
                                        processes: table.processList.map(item => {
                                            return {
                                                process_info: item
                                            }
                                        })
                                    }
                                })
                            }, { injectBizId: true })
                        })
                    }
                    this.$success(this.$t('添加成功'))
                    this.handleBackToModule()
                } catch (e) {
                    console.error(e)
                }
            },
            handleBackToModule () {
                this.$router.replace({
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: {
                        node: 'module-' + this.moduleId
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-layout {
        padding: 0 23px 0;
        font-size: 14px;
        color: #63656E;
    }
    .create-label{
        display: block;
        position: relative;
        line-height: 32px;
        &:after {
            content: "*";
            margin: 0 0 0 4px;
            color: $cmdbDangerColor;
            @include inlineBlock;
        }
    }
    .create-hosts {
        padding-left: 10px;
        height: 100%;
        overflow: hidden;
    }
    .select-host-button {
        height: 32px;
        line-height: 30px;
        font-size: 0;
        .bk-icon {
            position: static;
            height: 30px;
            line-height: 30px;
            font-size: 12px;
            font-weight: bold;
            @include inlineBlock(top);
        }
        /deep/ span {
            font-size: 14px;
        }
    }
    .select-host-count {
        color: $textColor;
        .count-number {
            font-weight: bold;
        }
    }
    .create-tables {
        max-height: calc(100% - 120px);
        margin: 22px 0 0 0;
        @include scrollbar-y;
    }
    .service-instance-table +  .service-instance-table {
        margin-top: 12px;
    }
    .buttons {
        padding: 20px 0;
    }
</style>
