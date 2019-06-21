<template>
    <div class="table-layout" v-show="show">
        <div class="table-title" @click="localExpanded = !localExpanded">
            <cmdb-form-bool class="title-checkbox"
                :size="16"
                v-model="checked"
                @click.native.stop>
            </cmdb-form-bool>
            <i class="title-icon bk-icon icon-down-shape" v-if="localExpanded"></i>
            <i class="title-icon bk-icon icon-right-shape" v-else></i>
            <span class="title-label">{{instance.name}}</span>
            <cmdb-dot-menu class="instance-menu" @click.native.stop>
                <ul class="menu-list">
                    <li class="menu-item"
                        v-for="(menu, index) in instanceMenu"
                        :key="index">
                        <button class="menu-button"
                            @click="menu.handler">
                            {{menu.name}}
                        </button>
                    </li>
                </ul>
            </cmdb-dot-menu>
        </div>
        <cmdb-table
            v-show="localExpanded"
            :loading="$loading(Object.values(requestId))"
            :header="header"
            :list="flattenList"
            :empty-height="42"
            :visible="localExpanded"
            :sortable="false"
            :reference-document-height="false">
            <template slot="data-empty">
                <template v-if="withTemplate">
                    <i18n path="BusinessTopology['暂无模板进程']">
                        <button class="add-process-button text-primary" place="link"
                            @click.stop="handleAddProcessToTemplate">
                            {{$t('BusinessTopology["模板添加"]')}}
                        </button>
                    </i18n>
                </template>
                <button class="add-process-button text-primary" v-else
                    @click.stop="handleAddProcess">
                    <i class="bk-icon icon-plus"></i>
                    <span>{{$t('BusinessTopology["添加进程"]')}}</span>
                </button>
            </template>
            <template slot="__operation__" slot-scope="{ item }">
                <button class="text-primary mr10"
                    @click="handleEditProcess(item)">
                    {{$t('Common["编辑"]')}}
                </button>
                <button class="text-primary" v-if="!withTemplate"
                    @click="handleDeleteProcess(item)">
                    {{$t('Common["删除"]')}}
                </button>
            </template>
        </cmdb-table>
        <div class="add-process-options" v-if="!withTemplate && localExpanded && list.length">
            <button class="add-process-button text-primary" @click="handleAddProcess">
                <i class="bk-icon icon-plus"></i>
                <span>{{$t('BusinessTopology["添加进程"]')}}</span>
            </button>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            instance: {
                type: Object,
                required: true
            },
            expanded: Boolean
        },
        data () {
            return {
                show: true,
                checked: false,
                localExpanded: this.expanded,
                properties: [],
                header: [],
                list: []
            }
        },
        computed: {
            isModuleNode () {
                const node = this.$store.state.businessTopology.selectedNode
                return node && node.bk_obj_id === 'module'
            },
            withTemplate () {
                return this.isModuleNode && !!this.instance.service_template_id
            },
            instanceMenu () {
                const menu = [{
                    name: this.$t('Common["删除"]'),
                    handler: this.handleDeleteInstance
                }]
                if (!this.withTemplate) {
                    menu.unshift({
                        name: this.$t('BusinessTopology["添加进程"]'),
                        handler: this.handleAddProcess
                    }, {
                        name: this.$t('BusinessTopology["克隆"]'),
                        handler: this.handleCloneInstance
                    })
                }
                return menu
            },
            module () {
                return this.$store.state.businessTopology.selectedNodeInstance
            },
            flattenList () {
                return this.$tools.flattenList(this.properties, this.list.map(data => data.property))
            },
            requestId () {
                return {
                    processList: `get_service_instance_${this.instance.id}_processes`,
                    properties: 'get_service_process_properties',
                    deleteProcess: 'delete_service_process'
                }
            }
        },
        watch: {
            localExpanded (expanded) {
                if (expanded) {
                    this.getServiceProcessList()
                }
            },
            checked (checked) {
                this.$emit('check-change', checked, this.instance)
            }
        },
        async created () {
            await this.getProcessProperties()
            if (this.expanded) {
                this.getServiceProcessList()
            }
        },
        methods: {
            async getProcessProperties () {
                const action = 'objectModelProperty/searchObjectAttribute'
                const properties = await this.$store.dispatch(action, {
                    params: {
                        bk_obj_id: 'process',
                        bk_supplier_account: this.$store.getters.supplierAccount
                    },
                    config: {
                        requestId: this.requestId.properties,
                        fromCache: true
                    }
                })
                this.properties = properties
                this.setHeader()
            },
            async getServiceProcessList () {
                try {
                    this.list = await this.$store.dispatch('processInstance/getServiceInstanceProcesses', {
                        params: this.$injectMetadata({
                            service_instance_id: this.instance.id
                        }),
                        config: {
                            requestId: this.requestId.processList
                        }
                    })
                } catch (e) {
                    this.list = []
                    console.error(e)
                }
            },
            setHeader () {
                const display = [
                    'bk_func_name',
                    'bk_process_name',
                    'bk_start_param_regex',
                    'bind_ip',
                    'port',
                    'work_path'
                ]
                const header = display.map(id => {
                    const property = this.properties.find(property => property.bk_property_id === id) || {}
                    return {
                        id: property.bk_property_id,
                        name: property.bk_property_name
                    }
                })
                header.push({
                    id: '__operation__',
                    name: this.$t('Common["操作"]')
                })
                this.header = header
            },
            handleAddProcess () {
                this.$emit('create-process', this)
            },
            async handleEditProcess (item) {
                const processInstance = this.list.find(data => data.relation.bk_process_id === item.bk_process_id)
                this.$emit('update-process', processInstance, this)
            },
            async handleDeleteProcess (item) {
                try {
                    await this.$store.dispatch('processInstance/deleteServiceInstanceProcess', {
                        serviceInstanceId: this.instance.id,
                        config: {
                            data: this.$injectMetadata({
                                process_instance_ids: [item.bk_process_id]
                            })
                        }
                    })
                    this.getServiceProcessList()
                } catch (e) {
                    console.error(e)
                }
            },
            handleCloneInstance () {
                this.$router.push({
                    name: 'cloneServiceInstance',
                    params: {
                        instanceId: this.instance.id,
                        hostId: this.instance.bk_host_id,
                        setId: this.module.bk_set_id,
                        moduleId: this.module.bk_module_id
                    },
                    query: {
                        from: this.$route.fullPath,
                        title: this.instance.name
                    }
                })
            },
            handleDeleteInstance () {
                this.$bkInfo({
                    title: this.$t('BusinessTopology["确认删除实例"]'),
                    content: this.$t('BusinessTopology["即将删除实例"]', { name: this.instance.name }),
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                                config: {
                                    data: this.$injectMetadata({
                                        service_instance_ids: [this.instance.id]
                                    }),
                                    requestId: this.requestId.deleteProcess
                                }
                            })
                            this.$emit('delete-instance', this.instance.id)
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            handleAddProcessToTemplate () {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: this.instance.service_template_id
                    },
                    query: {
                        from: this.$route.fullPath
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-layout {
        padding: 0 0 12px 0;
    }
    .table-title {
        height: 40px;
        padding: 0 11px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        cursor: pointer;
        .title-icon {
            font-size: 12px;
            color: #63656E;
            @include inlineBlock;
        }
        .title-label {
            font-size: 14px;
            color: #313238;
            @include inlineBlock;
        }
    }
    .add-process-options {
        border: 1px solid $cmdbTableBorderColor;
        border-top: none;
        line-height: 42px;
        font-size: 12px;
        text-align: center;
    }
    .add-process-button {
        line-height: 32px;
        .bk-icon,
        span {
            @include inlineBlock;
        }
    }
    .menu-list {
        min-width: 74px;
        padding: 6px 0;
        .menu-item {
            .menu-button {
                display: block;
                width: 100%;
                height: 32px;
                padding: 0 13px;
                line-height: 32px;
                outline: 0;
                border: none;
                text-align: left;
                color: #63656E;
                font-size: 12px;
                background-color: #fff;
                &:hover {
                    background-color: #E1ECFF;
                    color: #3A84FF;
                }
            }
        }
    }
</style>
