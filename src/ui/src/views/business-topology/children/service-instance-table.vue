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
            <i class="bk-icon icon-exclamation" v-if="withTemplate && !flattenList.length" v-bk-tooltips="tooltips"></i>
            <cmdb-dot-menu class="instance-menu" @click.native.stop>
                <ul class="menu-list">
                    <li class="menu-item"
                        v-for="(menu, index) in instanceMenu"
                        :key="index">
                        <span class="menu-span"
                            v-cursor="{
                                active: !$isAuthorized($OPERATION[menu.auth]),
                                auth: [$OPERATION[menu.auth]]
                            }">
                            <bk-button class="menu-button"
                                :text="true"
                                :disabled="!$isAuthorized($OPERATION[menu.auth])"
                                @click="menu.handler">
                                {{menu.name}}
                            </bk-button>
                        </span>
                    </li>
                </ul>
            </cmdb-dot-menu>
            <div class="instance-label clearfix" @click.stop v-if="labelShowList.length">
                <div class="label-title fl">
                    <i class="icon-cc-label"></i>
                    <span>{{$t('BusinessTopology["标签"]')}}：</span>
                </div>
                <div class="label-list fl">
                    <div class="label-item" :title="`${label.key}：${label.value}`" :key="index" v-for="(label, index) in labelShowList">
                        <span>{{label.key}}</span>
                        <span>:</span>
                        <span>{{label.value}}</span>
                    </div>
                    <bk-popover class="label-item label-tips"
                        v-if="labelTipsList.length"
                        theme="light"
                        :width="290"
                        placement="top-end">
                        <span>...</span>
                        <div class="tips-label-list" slot="content">
                            <span class="label-item" :title="`${label.key}：${label.value}`" :key="index" v-for="(label, index) in labelTipsList">
                                <span>{{label.key}}</span>
                                <span>:</span>
                                <span>{{label.value}}</span>
                            </span>
                        </div>
                    </bk-popover>
                </div>
            </div>
        </div>
        <bk-table
            v-show="localExpanded"
            v-bkloading="{ isLoading: $loading(Object.values(requestId)) }"
            :data="flattenList">
            <bk-table-column v-for="column in header"
                :key="column.id"
                :prop="column.id"
                :label="column.name">
            </bk-table-column>
            <bk-table-column :label="$t('Common[\'操作\']')">
                <template slot-scope="{ row }">
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_PROCESS),
                            auth: [$OPERATION.U_PROCESS]
                        }">
                        <bk-button class="mr10"
                            :text="true"
                            :disabled="!$isAuthorized($OPERATION.U_PROCESS)"
                            @click="handleEditProcess(row)">
                            {{$t('Common["编辑"]')}}
                        </bk-button>
                    </span>
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.D_PROCESS),
                            auth: [$OPERATION.D_PROCESS]
                        }">
                        <bk-button v-if="!withTemplate"
                            :text="true"
                            :disabled="!$isAuthorized($OPERATION.D_PROCESS)"
                            @click="handleDeleteProcess(row)">
                            {{$t('Common["删除"]')}}
                        </bk-button>
                    </span>
                </template>
            </bk-table-column>
            <template slot="empty">
                <template v-if="withTemplate">
                    <i18n path="BusinessTopology['暂无模板进程']">
                        <button class="add-process-button text-primary" place="link"
                            @click.stop="handleAddProcessToTemplate">
                            {{$t('BusinessTopology["模板添加"]')}}
                        </button>
                    </i18n>
                </template>
                <span style="display: inline-block;" v-else
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_PROCESS),
                        auth: [$OPERATION.C_PROCESS]
                    }">
                    <button class="add-process-button text-primary"
                        :disabled="!$isAuthorized($OPERATION.C_PROCESS)"
                        @click.stop="handleAddProcess">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t('BusinessTopology["添加进程"]')}}</span>
                    </button>
                </span>
            </template>
            <div class="add-process-options" v-if="!withTemplate && list.length" slot="append">
                <span style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.C_PROCESS),
                        auth: [$OPERATION.C_PROCESS]
                    }">
                    <button class="add-process-button text-primary"
                        :disabled="!$isAuthorized($OPERATION.C_PROCESS)"
                        @click="handleAddProcess">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t('BusinessTopology["添加进程"]')}}</span>
                    </button>
                </span>
            </div>
        </bk-table>
        <bk-dialog class="bk-dialog-no-padding"
            v-model="editLabel.show"
            :width="580"
            @confirm="handleSubmitEditLable"
            @cancel="handleCloseEditLable">
            <div slot="header">
                {{$t("BusinessTopology['编辑标签']")}}
            </div>
            <cmdb-edit-label
                ref="instanceLabel"
                :default-list="editLabel.list">
            </cmdb-edit-label>
        </bk-dialog>
    </div>
</template>

<script>
    import cmdbEditLabel from './edit-label.vue'
    export default {
        components: { cmdbEditLabel },
        props: {
            instance: {
                type: Object,
                required: true
            },
            expanded: Boolean
        },
        data () {
            return {
                editLabel: {
                    show: false,
                    list: []
                },
                tipsLabel: {
                    show: false,
                    id: null
                },
                show: true,
                checked: false,
                localExpanded: this.expanded,
                properties: [],
                header: [],
                list: [],
                tooltips: {
                    content: this.$t('BusinessTopology["模板未添加进程"]'),
                    placement: 'right'
                }
            }
        },
        computed: {
            isModuleNode () {
                const node = this.$store.state.businessTopology.selectedNode
                return node && node.data.bk_obj_id === 'module'
            },
            withTemplate () {
                return this.isModuleNode && !!this.instance.service_template_id
            },
            instanceMenu () {
                const menu = [{
                    name: this.$t('BusinessTopology["编辑标签"]'),
                    handler: this.handleShowEditLabel,
                    auth: 'U_SERVICE_INSTANCE'
                }, {
                    name: this.$t('Common["删除"]'),
                    handler: this.handleDeleteInstance,
                    auth: 'D_SERVICE_INSTANCE'
                }]
                if (!this.withTemplate) {
                    menu.unshift({
                        name: this.$t('BusinessTopology["添加进程"]'),
                        handler: this.handleAddProcess,
                        auth: 'C_PROCESS'
                    }, {
                        name: this.$t('BusinessTopology["克隆"]'),
                        handler: this.handleCloneInstance,
                        auth: 'C_PROCESS'
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
            },
            labelList () {
                const list = []
                const labels = this.instance.labels
                labels && Object.keys(labels).forEach((key, index) => {
                    list.push({
                        id: index,
                        key: key,
                        value: labels[key]
                    })
                })
                return list
            },
            labelShowList () {
                return this.labelList.slice(0, 3)
            },
            labelTipsList () {
                return this.labelList.slice(3)
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
                        from: {
                            name: this.$route.name,
                            query: {
                                module: this.module.bk_module_id
                            }
                        },
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
                        from: {
                            name: this.$route.name,
                            query: {
                                module: this.module.bk_module_id
                            }
                        }
                    }
                })
            },
            handleShowEditLabel () {
                this.editLabel.list = this.labelList
                this.editLabel.show = true
            },
            handleCloseEditLable () {
                this.editLabel.list = []
                this.editLabel.show = false
            },
            async handleSubmitEditLable () {
                try {
                    let status = ''
                    const validator = this.$refs.instanceLabel.$validator
                    const removeKeysList = this.$refs.instanceLabel.removeKeysList
                    const list = this.$refs.instanceLabel.submitList
                    const originList = this.$refs.instanceLabel.originList
                    const hasChange = JSON.stringify(this.$refs.instanceLabel.list) !== JSON.stringify(originList)

                    if (list.length && !await validator.validateAll()) {
                        return
                    }

                    if (removeKeysList.length) {
                        status = await this.$store.dispatch('instanceLabel/deleteInstanceLabel', {
                            config: {
                                data: this.$injectMetadata({
                                    instance_ids: [this.instance.id],
                                    keys: removeKeysList
                                }),
                                requestId: 'deleteInstanceLabel',
                                transformData: false
                            }
                        })
                    }

                    if (list.length && hasChange) {
                        const labelSet = {}
                        list.forEach(label => {
                            labelSet[label.key] = label.value
                        })
                        status = await this.$store.dispatch('instanceLabel/createInstanceLabel', {
                            params: this.$injectMetadata({
                                instance_ids: [this.instance.id],
                                labels: labelSet
                            }),
                            config: {
                                requestId: 'createInstanceLabel',
                                transformData: false
                            }
                        })
                    }
                    if (status && status.bk_error_msg === 'success') {
                        this.$success(this.$t('Common["保存成功"]'))
                        this.$parent.filter = ''
                        this.$parent.getServiceInstances()
                        this.$parent.getHistoryLabel()
                    }
                    this.handleCloseEditLable()
                } catch (e) {
                    console.error(e)
                }
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
        .icon-exclamation {
            width: 16px;
            height: 16px;
            line-height: 16px;
            font-size: 14px;
            text-align: center;
            color: #ffffff;
            background: #f0b659;
            border-radius: 50%;
        }
        .title-label {
            font-size: 14px;
            color: #313238;
            @include inlineBlock;
        }
    }
    .add-process-options {
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
            .menu-span {
                display: block;
            }
            .menu-button {
                display: block;
                width: 100%;
                height: 32px;
                line-height: 32px;
                color: #63656E;
                font-size: 12px;
                height: 32px;
                padding: 0 13px;
                text-align: left;
                &:hover {
                    background-color: #E1ECFF;
                    color: #3A84FF;
                }
                &:disabled {
                    color: #dcdee5;
                }
            }
        }
    }
    .instance-label {
        @include inlineBlock;
        font-size: 12px;
        .icon-cc-label {
            color: #979ba5;
            font-size: 16px;
        }
        .label-list {
            padding-left: 4px;
            line-height: 38px;
            font-size: 0;
            .label-item {
                @include inlineBlock;
                font-size: 12px;
                height: 20px;
                line-height: 20px;
                margin-right: 4px;
                padding: 0 6px;
                color: #979ba5;
                background-color: #fafbfd;
                border-radius: 2px;
                &>span {
                    @include ellipsis;
                    display: inline-block;
                    max-width: 54px;
                }
            }
            .label-tips {
                padding: 0;
                .bk-tooltip-ref span {
                    padding: 0 6px;
                }
                &:hover {
                    background-color: #f0f1f5;
                }
            }
        }
    }
</style>

<style lang="scss">
    .tips-label-list {
        .label-item {
            @include inlineBlock;
            font-size: 12px;
            height: 20px;
            line-height: 18px;
            margin: 5px 2px;
            padding: 0 6px;
            color: #979ba5;
            background-color: #fafbfd;
            border: 1px solid #dcdee5;
            border-radius: 2px;
            &>span {
                @include ellipsis;
                display: inline-block;
                max-width: 54px;
            }
        }
    }
</style>
