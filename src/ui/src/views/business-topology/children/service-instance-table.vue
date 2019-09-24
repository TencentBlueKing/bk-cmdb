<template>
    <div class="table-layout" v-show="show">
        <div class="table-title" @click="localExpanded = !localExpanded"
            @mouseenter="handleShowDotMenu"
            @mouseleave="handleHideDotMenu">
            <bk-checkbox class="title-checkbox"
                :size="16"
                v-model="checked"
                @click.native.stop>
            </bk-checkbox>
            <i class="title-icon bk-icon icon-down-shape" v-if="localExpanded"></i>
            <i class="title-icon bk-icon icon-right-shape" v-else></i>
            <span class="title-label">{{instance.name}}</span>
            <i class="bk-icon icon-exclamation" v-if="withTemplate && !instance.process_count" v-bk-tooltips="tooltips"></i>
            <cmdb-dot-menu class="instance-menu" ref="dotMenu" @click.native.stop>
                <ul class="menu-list"
                    @mouseenter="handleShowDotMenu"
                    @mouseleave="handleHideDotMenu">
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
            <div class="instance-label fr" @click.stop>
                <div :class="['label-title', 'fl', { 'disabled': !$isAuthorized($OPERATION.U_SERVICE_INSTANCE) }]"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.U_SERVICE_INSTANCE),
                        auth: [$OPERATION.U_SERVICE_INSTANCE]
                    }"
                    @click.stop="handleShowEditLabel">
                    <i class="icon-cc-label"></i>
                    <span v-if="!labelShowList.length"> + </span>
                </div>
                <div class="label-list fl">
                    <div class="label-item" :title="`${label.key}：${label.value}`" :key="index" v-for="(label, index) in labelShowList">
                        <span>{{label.key}}</span>
                        <span>:</span>
                        <span>{{label.value}}</span>
                    </div>
                    <bk-popover class="label-item label-tips"
                        v-if="labelTipsList.length"
                        theme="light label-tips"
                        :width="290"
                        placement="bottom-end">
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
            <bk-table-column :label="$t('操作')">
                <template slot-scope="{ row }">
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_SERVICE_INSTANCE),
                            auth: [$OPERATION.U_SERVICE_INSTANCE]
                        }">
                        <bk-button class="mr10"
                            :text="true"
                            :disabled="!$isAuthorized($OPERATION.U_SERVICE_INSTANCE)"
                            @click="handleEditProcess(row)">
                            {{$t('编辑')}}
                        </bk-button>
                    </span>
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.U_SERVICE_INSTANCE),
                            auth: [$OPERATION.U_SERVICE_INSTANCE]
                        }">
                        <bk-button v-if="!withTemplate"
                            :text="true"
                            :disabled="!$isAuthorized($OPERATION.U_SERVICE_INSTANCE)"
                            @click="handleDeleteProcess(row)">
                            {{$t('删除')}}
                        </bk-button>
                    </span>
                </template>
            </bk-table-column>
            <template slot="empty">
                <template v-if="withTemplate">
                    <i18n path="暂无模板进程">
                        <button class="add-process-button text-primary" place="link"
                            @click.stop="handleAddProcessToTemplate">
                            {{$t('模板添加')}}
                        </button>
                    </i18n>
                </template>
                <span style="display: inline-block;" v-else
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.U_SERVICE_INSTANCE),
                        auth: [$OPERATION.U_SERVICE_INSTANCE]
                    }">
                    <button class="add-process-button text-primary"
                        :disabled="!$isAuthorized($OPERATION.U_SERVICE_INSTANCE)"
                        @click.stop="handleAddProcess">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t('添加进程')}}</span>
                    </button>
                </span>
            </template>
        </bk-table>
        <bk-dialog class="bk-dialog-no-padding edit-label-dialog"
            v-model="editLabel.show"
            :width="580"
            :mask-close="false"
            :esc-close="false"
            @confirm="handleSubmitEditLable"
            @cancel="handleCloseEditLable"
            @after-leave="handleSetEditBox">
            <div slot="header">
                {{$t('编辑标签')}}
            </div>
            <template v-if="editLabel.visiable">
                <cmdb-edit-label
                    ref="instanceLabel"
                    :default-list="editLabel.list">
                </cmdb-edit-label>
            </template>
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
                    visiable: false,
                    id: null
                },
                show: true,
                checked: false,
                localExpanded: this.expanded,
                properties: [],
                header: [],
                list: [],
                tooltips: {
                    content: this.$t('模板未添加进程'),
                    placement: 'right'
                }
            }
        },
        computed: {
            currentNode () {
                return this.$store.state.businessTopology.selectedNode
            },
            isModuleNode () {
                const node = this.$store.state.businessTopology.selectedNode
                return node && node.data.bk_obj_id === 'module'
            },
            withTemplate () {
                return this.isModuleNode && !!this.instance.service_template_id
            },
            instanceMenu () {
                const menu = [{
                    name: this.$t('删除'),
                    handler: this.handleDeleteInstance,
                    auth: 'D_SERVICE_INSTANCE'
                }]
                if (!this.withTemplate) {
                    menu.unshift({
                        name: this.$t('添加进程'),
                        handler: this.handleAddProcess,
                        auth: 'U_SERVICE_INSTANCE'
                    }, {
                        name: this.$t('克隆'),
                        handler: this.handleCloneInstance,
                        auth: 'C_SERVICE_INSTANCE'
                    })
                }
                return menu
            },
            module () {
                return this.$store.state.businessTopology.selectedNodeInstance
            },
            flattenList () {
                return this.$tools.flattenList(this.properties, this.list.map(data => data.property || {}))
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
                    this.$success(this.$t('删除成功'))
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
                        setId: this.currentNode.parent.data.bk_inst_id,
                        moduleId: this.module.bk_module_id
                    },
                    query: {
                        title: this.instance.name
                    }
                })
            },
            handleDeleteInstance () {
                this.$bkInfo({
                    title: this.$t('确认删除实例'),
                    subTitle: this.$t('即将删除实例', { name: this.instance.name }),
                    extCls: 'bk-dialog-sub-header-center',
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
                            this.currentNode.data.service_instance_count = this.currentNode.data.service_instance_count - 1
                            this.currentNode.parents.forEach(node => {
                                node.data.service_instance_count = node.data.service_instance_count - 1
                            })
                            this.$success(this.$t('删除成功'))
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
                    }
                })
            },
            handleShowEditLabel () {
                if (!this.$isAuthorized(this.$OPERATION.U_SERVICE_INSTANCE)) return
                this.editLabel.list = this.labelList
                this.editLabel.show = true
                this.editLabel.visiable = true
            },
            handleCloseEditLable () {
                this.editLabel.show = false
            },
            handleSetEditBox () {
                this.editLabel.list = []
                this.editLabel.visiable = false
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
                        this.$success(this.$t('保存成功'))
                        this.$parent.handleCheckALL(false)
                        this.$parent.filter = ''
                        this.$parent.getServiceInstances()
                        this.$parent.getHistoryLabel()
                    }
                    this.handleCloseEditLable()
                    setTimeout(() => {
                        this.handleSetEditBox()
                    }, 200)
                } catch (e) {
                    console.error(e)
                }
            },
            handleShowDotMenu () {
                this.$refs.dotMenu.$el.style.opacity = 1
            },
            handleHideDotMenu () {
                this.$refs.dotMenu.$el.style.opacity = 0
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
        padding: 0 10px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        cursor: pointer;
        .title-checkbox {
            /deep/ .bk-checkbox {
                background-color: #fff;
            }
            &.is-checked {
                /deep/ .bk-checkbox {
                    background-color: #3a84ff !important;
                }
            }
        }
        .title-icon {
            font-size: 14px;
            margin: 0 2px 0 6px;
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
        .instance-menu {
            opacity: 0;
            /deep/ .bk-tooltip-ref {
                width: 100%;
            }
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
        .label-title {
            color: #979ba5;
            &:hover {
                color: #3a84ff;
            }
            &.disabled {
                color: #979ba5 !important;
            }
        }
        .icon-cc-label {
            font-size: 16px;
        }
        .label-list {
            padding-left: 8px;
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
                /deep/ .bk-tooltip-ref {
                    padding: 0 6px;
                    span {
                        line-height: 16px;
                        display: inline-block;
                        vertical-align: top;
                    }
                }
                &:hover {
                    background-color: #e1ecff;
                }
            }
        }
    }
    .edit-label-dialog {
        /deep/ .bk-dialog-header {
            text-align: left !important;
            font-size: 24px;
            color: #444444;
            margin-top: -15px;
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
    .tippy-tooltip.label-tips-theme {
        padding: 8px 6px !important;
    }
</style>
