<template>
    <div class="layout" v-bkloading="{ isLoading: $loading('getModuleServiceInstances') }">
        <template v-if="instances.length">
            <div class="options">
                <bk-button class="options-button" type="primary"
                    @click="handleCreateServiceInstance">
                    {{$t('BusinessTopology["添加服务实例"]')}}
                </bk-button>
                <bk-button class="options-button" type="default"
                    v-if="withTemplate"
                    @click="handleSyncTemplate">
                    {{$t('BusinessTopology["同步模板"]')}}
                </bk-button>
                <bk-dropdown-menu trigger="click">
                    <bk-button class="options-button clipboard-trigger" type="default" slot="dropdown-trigger">
                        {{$t('Common["更多"]')}}
                        <i class="bk-icon icon-angle-down"></i>
                    </bk-button>
                    <ul class="clipboard-list" slot="dropdown-content">
                        <li v-for="(item, index) in menuItem"
                            :class="['clipboard-item', { 'is-disabled': item.disabled }]"
                            :key="index"
                            @click="item.handler(item.disabled)">
                            {{item.name}}
                        </li>
                    </ul>
                </bk-dropdown-menu>
                <cmdb-form-bool class="options-checkbox"
                    :size="16"
                    @change="handleCheckALL">
                    <span class="checkbox-label">{{$t('Common["全选本页"]')}}</span>
                </cmdb-form-bool>
                <cmdb-form-bool class="options-checkbox"
                    :size="16"
                    @change="handleExpandAll">
                    <span class="checkbox-label">{{$t('Common["全部展开"]')}}</span>
                </cmdb-form-bool>
                <cmdb-form-singlechar class="options-search fr"></cmdb-form-singlechar>
            </div>
            <div class="tables">
                <service-instance-table
                    v-for="(instance, index) in instances"
                    ref="serviceInstanceTable"
                    :key="instance.id"
                    :instance="instance"
                    :expanded="index === 0"
                    @create-process="handleCreateProcess"
                    @update-process="handleUpdateProcess"
                    @delete-instance="handleDeleteInstance">
                </service-instance-table>
            </div>
        </template>
        <service-instance-empty v-else
            @create-instance-success="handleCreateInstanceSuccess">
        </service-instance-empty>
        <cmdb-slider
            :title="processForm.title"
            :is-show.sync="processForm.show">
            <cmdb-form slot="content" v-if="processForm.show"
                :type="processForm.type"
                :inst="processForm.instance"
                :uneditable-properties="processForm.uneditableProperties"
                :properties="processForm.properties"
                :property-groups="processForm.propertyGroups"
                @on-submit="handleSaveProcess"
                @on-cancel="handleCloseProcessForm">
            </cmdb-form>
        </cmdb-slider>
    </div>
</template>

<script>
    import serviceInstanceTable from './service-instance-table.vue'
    import serviceInstanceEmpty from './service-instance-empty.vue'
    export default {
        components: {
            serviceInstanceTable,
            serviceInstanceEmpty
        },
        data () {
            return {
                instances: [],
                processForm: {
                    type: 'create',
                    show: false,
                    title: '',
                    instance: null,
                    referenceService: null,
                    uneditableProperties: [],
                    properties: [],
                    propertyGroups: []
                }
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            currentNode () {
                return this.$store.state.businessTopology.selectedNode
            },
            currentModule () {
                if (this.currentNode && this.currentNode.data.bk_obj_id === 'module') {
                    return this.currentNode.data
                }
                return null
            },
            processTemplateMap () {
                return this.$store.state.businessTopology.processTemplateMap
            },
            withTemplate () {
                return this.currentModule && this.currentModule.service_template_id !== 2
            },
            menuItem () {
                return [{
                    name: this.$t('BusinessTopology["批量编辑"]'),
                    handler: this.batchEdit,
                    disabled: true
                }, {
                    name: this.$t('BusinessTopology["批量删除"]'),
                    handler: this.batchDelete,
                    disabled: true
                }, {
                    name: this.$t('BusinessTopology["复制IP"]'),
                    handler: this.copyIp,
                    disabled: true
                }]
            }
        },
        watch: {
            currentModule (module) {
                if (module) {
                    this.getServiceInstances()
                }
            }
        },
        created () {
            this.getProcessProperties()
            this.getProcessPropertyGroups()
        },
        methods: {
            async getProcessProperties () {
                try {
                    const action = 'objectModelProperty/searchObjectAttribute'
                    this.processForm.properties = await this.$store.dispatch(action, {
                        params: {
                            bk_obj_id: 'process',
                            bk_supplier_account: this.$store.getters.supplierAccount
                        },
                        config: {
                            requestId: 'get_service_process_properties',
                            fromCache: true
                        }
                    })
                } catch (e) {
                    console.error(e)
                    this.processForm.properties = []
                }
            },
            async getProcessPropertyGroups () {
                try {
                    const action = 'objectModelFieldGroup/searchGroup'
                    this.processForm.propertyGroups = await this.$store.dispatch(action, {
                        objId: 'process',
                        params: {},
                        config: {
                            requestId: 'get_service_process_property_groups',
                            fromCache: true
                        }
                    })
                } catch (e) {
                    this.processForm.propertyGroups = []
                    console.error(e)
                }
            },
            async getServiceInstances () {
                try {
                    const data = await this.$store.dispatch('serviceInstance/getModuleServiceInstances', {
                        params: this.$injectMetadata({
                            bk_module_id: this.currentModule.bk_inst_id,
                            with_name: true
                        }),
                        config: {
                            requestId: 'getModuleServiceInstances',
                            cancelPrevious: true
                        }
                    })
                    this.instances = data.info
                } catch (e) {
                    console.error(e)
                    this.instances = []
                }
            },
            handleCreateProcess (referenceService) {
                this.processForm.referenceService = referenceService
                this.processForm.type = 'create'
                this.processForm.title = `${this.$t('BusinessToplogy["创建进程"]')}(${referenceService.instance.name})`
                this.processForm.instance = {}
                this.processForm.show = true
            },
            async handleUpdateProcess (processInstance, referenceService) {
                this.processForm.referenceService = referenceService
                this.processForm.type = 'update'
                this.processForm.title = this.$t('BusinessToplogy["编辑进程"]')
                this.processForm.instance = processInstance.property
                this.processForm.show = true

                const processTemplateId = processInstance.relation.process_template_id
                if (processTemplateId) {
                    const template = await this.getProcessTemplate(processTemplateId)
                    const uneditableProperties = []
                    Object.keys(template).forEach(propertyId => {
                        const value = template[propertyId]
                        if (value.as_default_value) {
                            uneditableProperties.push(propertyId)
                        }
                    })
                    this.processForm.uneditableProperties = uneditableProperties
                }
            },
            async getProcessTemplate (processTemplateId) {
                if (this.processTemplateMap.hasOwnProperty(processTemplateId)) {
                    return Promise.resolve(this.processTemplateMap[processTemplateId])
                }
                const data = await this.$store.dispatch('processTemplate/getProcessTemplate', {
                    params: { processTemplateId },
                    config: {
                        requestId: 'getProcessTemplate'
                    }
                })
                this.$store.commit('businessTopology/setProcessTemplate', {
                    id: processTemplateId,
                    template: data.property
                })
                return Promise.resolve(data.property)
            },
            handleDeleteInstance (id) {
                this.instances = this.instances.filter(instance => instance.id !== id)
            },
            async handleSaveProcess (values, changedValues) {
                try {
                    if (this.processForm.type === 'create') {
                        await this.createProcess(values)
                    } else {
                        await this.updateProcess(changedValues)
                    }
                    this.processForm.referenceService.getServiceProcessList()
                    this.processForm.show = false
                    this.processForm.instance = null
                    this.processForm.referenceService = null
                    this.processForm.uneditableProperties = null
                } catch (e) {
                    console.error(e)
                }
            },
            createProcess (values) {
                return this.$store.dispatch('processInstance/createServiceInstanceProcess', {
                    params: this.$injectMetadata({
                        service_instance_id: this.processForm.referenceService.instance.id,
                        processes: [{
                            process_info: values
                        }]
                    })
                })
            },
            updateProcess (values) {
                return this.$store.dispatch('processInstance/updateServiceInstanceProcess', {
                    business: this.business,
                    processInstanceId: this.processForm.instance.bk_process_id,
                    params: values
                })
            },
            handleCloseProcessForm () {
                this.processForm.show = false
                this.processForm.referenceService = null
                this.processForm.instance = null
                this.processForm.uneditableProperties = []
            },
            handleCreateServiceInstance () {
                this.$router.push({
                    name: 'createServiceInstance',
                    params: {
                        moduleId: this.currentNode.data.bk_inst_id,
                        setId: this.currentNode.parent.data.bk_inst_id
                    }
                })
            },
            handleSyncTemplate () {
                this.$router.push({
                    name: 'synchronous',
                    params: {
                        moduleId: this.currentModule.bk_inst_id,
                        setId: this.currentNode.parent.data.bk_inst_id
                    },
                    query: {
                        path: 'xxxxxxxxxxxx'
                    }
                })
            },
            handleCreateInstanceSuccess () {
                this.getServiceInstances()
            },
            handleCheckALL () {
                // 全选
            },
            handleExpandAll (expanded) {
                this.$refs.serviceInstanceTable.forEach(table => {
                    table.localExpanded = expanded
                })
            },
            batchEdit (disabled) {
                if (disabled) {
                    return false
                }
            },
            batchDelete (disabled) {
                if (disabled) {
                    return false
                }
            },
            copyIp () {}
        }
    }
</script>

<style lang="scss" scoped>
    .options {
        padding: 15px 0;
    }
    .options-button {
        height: 32px;
        padding: 0 8px;
        margin: 0 6px 0 0;
        line-height: 30px;
    }
    .options-checkbox {
        margin: 0 19px 0 10px;
        .checkbox-label {
            padding: 0 0 0 9px;
            line-height: 1.5;
        }
    }
    .options-search {
        /deep/ {
            .cmdb-form-input {
                height: 32px;
                line-height: 30px;
            }
        }
    }
    .clipboard-trigger{
        padding: 0 16px;
        .icon-angle-down {
            font-size: 12px;
            top: 0;
        }
    }
    .clipboard-list{
        width: 100%;
        font-size: 14px;
        line-height: 40px;
        max-height: 160px;
        @include scrollbar-y;
        &::-webkit-scrollbar{
            width: 3px;
            height: 3px;
        }
        .clipboard-item{
            padding: 0 15px;
            cursor: pointer;
            @include ellipsis;
            &:not(.is-disabled):hover{
                background-color: #ebf4ff;
                color: #3c96ff;
            }
            &.is-disabled {
                color: #c4c6cc;
                cursor: not-allowed;
            }
        }
    }
    .tables {
        height: calc(100% - 42px);
        @include scrollbar-y;
    }
</style>
