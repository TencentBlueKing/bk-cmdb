<template>
    <div class="layout" v-bkloading="{ isLoading: $loading('getModuleServiceInstances') }">
        <template v-if="instances.length || inSearch">
            <div class="options">
                <cmdb-form-bool class="options-checkall"
                    :size="16"
                    :checked="isCheckAll"
                    :title="$t('Common[\'全选本页\']')"
                    @change="handleCheckALL">
                </cmdb-form-bool>
                <bk-button class="options-button" type="primary"
                    @click="handleCreateServiceInstance">
                    {{$t('BusinessTopology["添加服务实例"]')}}
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
                <div class="options-right fr">
                    <cmdb-form-bool class="options-checkbox"
                        :size="16"
                        :checked="isExpandAll"
                        @change="handleExpandAll">
                        <span class="checkbox-label">{{$t('Common["全部展开"]')}}</span>
                    </cmdb-form-bool>
                    <cmdb-form-singlechar class="options-search"
                        :placeholder="$t('BusinessTopology[\'请输入IP搜索\']')"
                        v-model="filter">
                        <i class="bk-icon icon-close"
                            v-show="filter.length"
                            @click="handleClearFilter">
                        </i>
                        <i class="bk-icon icon-search" @click.stop="handleSearch"></i>
                    </cmdb-form-singlechar>
                </div>
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
                    @delete-instance="handleDeleteInstance"
                    @check-change="handleCheckChange">
                </service-instance-table>
            </div>
            <bk-paging class="pagination"
                v-if="pagination.totalPage > 1"
                pagination-able
                location="left"
                :cur-page="pagination.current"
                :total-page="pagination.totalPage"
                :pagination-count="pagination.size"
                @page-change="handlePageChange"
                @pagination-change="handleSizeChange">
            </bk-paging>
            <div class="filter-empty" v-if="!instances.length">
                <div class="filter-empty-content">
                    <i class="bk-icon icon-empty"></i>
                    <span>{{$t('BusinessTopology["暂无符合条件的实例"]')}}</span>
                </div>
            </div>
        </template>
        <service-instance-empty v-else
            @create-instance-success="handleCreateInstanceSuccess">
        </service-instance-empty>
        <cmdb-slider
            :title="processForm.title"
            :is-show.sync="processForm.show"
            :before-close="handleBeforeClose">
            <cmdb-form slot="content" v-if="processForm.show"
                ref="processForm"
                :type="processForm.type"
                :inst="processForm.instance"
                :uneditable-properties="processForm.uneditableProperties"
                :properties="processForm.properties"
                :property-groups="processForm.propertyGroups"
                @on-submit="handleSaveProcess"
                @on-cancel="handleBeforeClose">
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
                checked: [],
                isCheckAll: false,
                isExpandAll: false,
                filter: '',
                inSearch: false,
                instances: [],
                pagination: {
                    current: 1,
                    totalPage: 0,
                    size: 10
                },
                processForm: {
                    type: 'create',
                    show: false,
                    title: '',
                    instance: null,
                    referenceService: null,
                    uneditableProperties: [],
                    properties: [],
                    propertyGroups: [],
                    unwatch: null
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
                    return this.$store.state.businessTopology.selectedNodeInstance
                }
                return null
            },
            processTemplateMap () {
                return this.$store.state.businessTopology.processTemplateMap
            },
            menuItem () {
                return [{
                    name: this.$t('BusinessTopology["批量删除"]'),
                    handler: this.batchDelete,
                    disabled: !this.checked.length
                }, {
                    name: this.$t('BusinessTopology["复制IP"]'),
                    handler: this.copyIp,
                    disabled: !this.checked.length
                }]
            }
        },
        watch: {
            currentNode (node) {
                if (node && node.data.bk_obj_id === 'module') {
                    this.filter = ''
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
                            bk_module_id: this.currentNode.data.bk_inst_id,
                            with_name: true,
                            page: {
                                start: (this.pagination.current - 1) * this.pagination.size,
                                limit: this.pagination.size
                            },
                            search_key: this.filter
                        }),
                        config: {
                            requestId: 'getModuleServiceInstances',
                            cancelPrevious: true
                        }
                    })
                    this.checked = []
                    this.isCheckAll = false
                    this.isExpandAll = false
                    this.instances = data.info
                    this.pagination.totalPage = Math.ceil(data.count / this.pagination.size)
                } catch (e) {
                    console.error(e)
                    this.instances = []
                }
            },
            handleSearch () {
                this.inSearch = true
                this.handlePageChange(1)
            },
            handleClearFilter () {
                this.filter = ''
                this.handleSearch()
            },
            handlePageChange (page) {
                this.pagination.current = page
                this.getServiceInstances()
            },
            handleSizeChange (size) {
                this.pagination.current = 1
                this.pagination.size = size
                this.getServiceInstances()
            },
            handleCheckChange (checked, instance) {
                if (checked) {
                    this.checked.push(instance)
                } else {
                    this.checked = this.checked.filter(target => target.id !== instance.id)
                }
            },
            handleCreateProcess (referenceService) {
                this.processForm.referenceService = referenceService
                this.processForm.type = 'create'
                this.processForm.title = `${this.$t('BusinessTopology["添加进程"]')}(${referenceService.instance.name})`
                this.processForm.instance = {}
                this.processForm.show = true
                this.$nextTick(() => {
                    const { processForm } = this.$refs
                    this.processForm.unwatch = processForm.$watch(() => {
                        return processForm.values.bk_func_name
                    }, (newVal, oldValue) => {
                        if (processForm.values.bk_process_name === oldValue) {
                            processForm.values.bk_process_name = newVal
                        }
                    })
                })
            },
            async handleUpdateProcess (processInstance, referenceService) {
                this.processForm.referenceService = referenceService
                this.processForm.type = 'update'
                this.processForm.title = this.$t('BusinessTopology["编辑进程"]')
                this.processForm.instance = processInstance.property
                this.processForm.show = true

                const processTemplateId = processInstance.relation.process_template_id
                if (processTemplateId) {
                    const template = await this.getProcessTemplate(processTemplateId)
                    const uneditableProperties = []
                    Object.keys(template).forEach(propertyId => {
                        const value = template[propertyId]
                        if (value.as_default_value && value.value) {
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
            async handleSaveProcess (values, changedValues, instance) {
                try {
                    this.processForm.unwatch && this.processForm.unwatch()
                    if (this.processForm.type === 'create') {
                        await this.createProcess(values)
                    } else {
                        await this.updateProcess(values, instance)
                    }
                    await this.processForm.referenceService.getServiceProcessList()
                    this.updateServiceInstanceName()
                    this.processForm.show = false
                    this.processForm.instance = null
                    this.processForm.referenceService = null
                    this.processForm.uneditableProperties = []
                } catch (e) {
                    console.error(e)
                }
            },
            updateServiceInstanceName () {
                const serviceInstance = this.processForm.referenceService
                const processes = serviceInstance.list
                const instance = this.instances.find(instance => instance === serviceInstance.instance) || {}
                const name = (instance.name || '').split('_').slice(0, 1)
                if (processes.length) {
                    const process = processes[0].property
                    name.push(process.bk_process_name)
                    name.push(process.port)
                }
                instance.name = name.join('_')
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
            updateProcess (values, instance) {
                return this.$store.dispatch('processInstance/updateServiceInstanceProcess', {
                    params: this.$injectMetadata({
                        processes: [{ ...instance, ...values }]
                    })
                })
            },
            handleCloseProcessForm () {
                this.processForm.show = false
                this.processForm.referenceService = null
                this.processForm.instance = null
                this.processForm.uneditableProperties = []
            },
            handleBeforeClose () {
                const changedValues = this.$refs.processForm.changedValues
                if (Object.keys(changedValues).length) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('Common["退出会导致未保存信息丢失，是否确认？"]'),
                            confirmFn: () => {
                                this.handleCloseProcessForm()
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.handleCloseProcessForm()
            },
            handleCreateServiceInstance () {
                this.$router.push({
                    name: 'createServiceInstance',
                    params: {
                        moduleId: this.currentNode.data.bk_inst_id,
                        setId: this.currentNode.parent.data.bk_inst_id
                    },
                    query: {
                        from: this.$route.fullPath,
                        title: this.currentNode.name
                    }
                })
            },
            handleCreateInstanceSuccess () {
                this.getServiceInstances()
            },
            handleCheckALL (checked) {
                this.filter = ''
                this.isCheckAll = checked
                this.$refs.serviceInstanceTable.forEach(table => {
                    table.checked = checked
                })
            },
            handleExpandAll (expanded) {
                this.filter = ''
                this.isExpandAll = expanded
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
                this.$bkInfo({
                    title: this.$t('BusinessTopology["确认删除实例"]'),
                    content: this.$t('BusinessTopology["即将删除选中的实例"]', { count: this.checked.length }),
                    confirmFn: async () => {
                        try {
                            const serviceInstanceIds = this.checked.map(instance => instance.id)
                            await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                                config: {
                                    data: this.$injectMetadata({
                                        service_instance_ids: serviceInstanceIds
                                    }),
                                    requestId: 'batchDeleteServiceInstance'
                                }
                            })
                            this.instances = this.instances.filter(instance => !serviceInstanceIds.includes(instance.id))
                            this.checked = []
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            copyIp () {
                this.$copyText(this.checked.map(instance => instance.name.split('_')[0]).join('\n')).then(() => {
                    this.$success(this.$t('Common["复制成功"]'))
                }, () => {
                    this.$error(this.$t('Common["复制失败"]'))
                })
            }
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
        margin: 0 0 0 6px;
        line-height: 30px;
    }
    .options-checkall {
        width: 36px;
        height: 32px;
        line-height: 30px;
        padding: 0 9px;
        text-align: center;
        border: 1px solid #C4C6CC;
        border-radius: 2px;
    }
    .options-right {
        text-align: right;
        white-space: nowrap;
    }
    .options-checkbox {
        margin: 0 15px 0 0;
        .checkbox-label {
            padding: 0 0 0 4px;
        }
    }
    .options-search {
        @include inlineBlock;
        position: relative;
        width: 240px;
        .icon-search {
            position: absolute;
            top: 9px;
            right: 9px;
            font-size: 14px;
            cursor: pointer;
        }
        .icon-close {
            position: absolute;
            top: 8px;
            right: 30px;
            width: 16px;
            height: 16px;
            line-height: 16px;
            border-radius: 50%;
            text-align: center;
            background-color: #ddd;
            color: #fff;
            font-size: 12px;
            transition: backgroundColor .2s linear;
            cursor: pointer;
            &:before {
                display: block;
                transform: scale(.7);
            }
            &:hover {
                background-color: #ccc;
            }
        }
        /deep/ {
            .cmdb-form-input {
                height: 32px;
                line-height: 30px;
                padding-right: 50px;
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
        max-height: calc(100% - 120px);
        @include scrollbar-y;
    }
    .pagination {
        padding: 10px 0 0 0;
    }
    .filter-empty {
        width: 100%;
        height: calc(100% - 130px);
        display: table;
        .filter-empty-content {
            display: table-cell;
            vertical-align: middle;
            text-align: center;
            .icon-empty {
                display: block;
                margin: 0 0 10px 0;
                font-size: 65px;
                color: #c3cdd7;
            }
        }
    }
</style>
