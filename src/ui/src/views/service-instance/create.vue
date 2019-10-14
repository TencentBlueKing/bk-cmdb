<template>
    <div class="create-layout clearfix" v-bkloading="{ isLoading: $loading() }">
        <label class="create-label fl">{{$t('添加主机')}}</label>
        <div class="create-hosts">
            <bk-button class="select-host-button" theme="default"
                @click="hostSelectorVisible = true">
                <i class="bk-icon icon-plus"></i>
                {{$t('添加主机')}}
            </bk-button>
            <div class="create-tables">
                <service-instance-table class="service-instance-table"
                    v-for="(host, index) in hosts"
                    ref="serviceInstanceTable"
                    deletable
                    expanded
                    :key="index"
                    :index="index"
                    :id="host.bk_host_id"
                    :name="host.bk_host_innerip"
                    :source-processes="sourceProcesses"
                    :templates="templates"
                    :show-tips="showTips"
                    @delete-instance="handleDeleteInstance">
                </service-instance-table>
                <div class="buttons">
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.C_SERVICE_INSTANCE),
                            auth: [$OPERATION.C_SERVICE_INSTANCE]
                        }">
                        <bk-button theme="primary"
                            :disabled="!hosts.length || !$isAuthorized($OPERATION.C_SERVICE_INSTANCE)"
                            @click="handleConfirm">
                            {{$t('确定')}}
                        </bk-button>
                    </span>
                    <bk-button @click="handleBackToModule">{{$t('取消')}}</bk-button>
                </div>
            </div>
        </div>
        <host-selector
            :visible.sync="hostSelectorVisible"
            :module-instance="moduleInstance"
            @host-selected="handleSelectHost">
        </host-selector>
    </div>
</template>

<script>
    import hostSelector from '@/components/ui/selector/host.vue'
    import serviceInstanceTable from '@/components/service/instance-table.vue'
    import { MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
    export default {
        name: 'create-service-instance',
        components: {
            hostSelector,
            serviceInstanceTable
        },
        data () {
            return {
                seed: 0,
                hostSelectorVisible: false,
                moduleInstance: {},
                hosts: [],
                templates: [],
                showTips: false
            }
        },
        computed: {
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
        created () {
            this.$store.commit('setBreadcrumbs', [{
                label: this.$t('服务拓扑'),
                route: {
                    name: MENU_BUSINESS_SERVICE_TOPOLOGY,
                    query: {
                        module: this.$route.params.moduleId
                    }
                }
            }, {
                label: this.$route.query.title
            }])
            this.getModuleInstance()
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
            async getTemplate () {
                try {
                    const data = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: this.$injectMetadata({
                            service_template_id: this.moduleInstance.service_template_id
                        }),
                        config: {
                            requestId: 'getBatchProcessTemplate',
                            cancelPrevious: true
                        }
                    })
                    this.templates = data.info
                } catch (e) {
                    console.error(e)
                }
            },
            handleSelectHost (checked, hosts) {
                this.hosts.push(...hosts)
                this.hostSelectorVisible = false
            },
            handleDeleteInstance (index) {
                this.hosts.splice(index, 1)
            },
            async handleConfirm () {
                try {
                    const serviceInstanceTables = this.$refs.serviceInstanceTable
                    if (serviceInstanceTables.some(table => !table.processList.length)) {
                        this.showTips = true
                        return
                    }
                    if (this.withTemplate) {
                        await this.$store.dispatch('serviceInstance/createProcServiceInstanceByTemplate', {
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
                            })
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
                            })
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
                    name: MENU_BUSINESS_SERVICE_TOPOLOGY,
                    query: {
                        module: this.moduleId
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
    .create-tables {
        height: calc(100% - 54px);
        margin: 22px 0 0 0;
        @include scrollbar-y;
        .buttons {
            padding: 18px 0 0 0;
        }
    }
    .service-instance-table {
        margin-bottom: 12px;
    }
</style>
