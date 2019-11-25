<template>
    <div class="service-layout" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <template v-if="instances.length || inSearch">
            <div class="options">
                <bk-checkbox class="options-checkall"
                    :size="16"
                    v-model="isCheckAll"
                    :title="$t('全选本页')"
                    @change="handleCheckALL">
                </bk-checkbox>
                <cmdb-auth :auth="$authResources({ type: $OPERATION.C_SERVICE_INSTANCE })">
                    <bk-button slot-scope="{ disabled }" class="options-button" theme="primary"
                        :disabled="disabled"
                        @click="handleCreateServiceInstance">
                        {{$t('新增')}}
                    </bk-button>
                </cmdb-auth>
                <bk-dropdown-menu trigger="click" font-size="medium">
                    <bk-button class="options-button clipboard-trigger" theme="default" slot="dropdown-trigger">
                        {{$t('更多')}}
                        <i class="bk-icon icon-angle-down"></i>
                    </bk-button>
                    <ul class="clipboard-list" slot="dropdown-content">
                        <li v-for="(item, index) in menuItem"
                            class="clipboard-item"
                            :key="index">
                            <cmdb-auth v-if="item.auth" :auth="$authResources({ type: $OPERATION[item.auth] })">
                                <bk-button slot-scope="{ disabled }"
                                    class="item-btn"
                                    text
                                    :disabled="item.disabled || disabled"
                                    @click="item.handler(item.disabled)">
                                    {{item.name}}
                                </bk-button>
                            </cmdb-auth>
                            <bk-button v-else text
                                class="item-btn"
                                :disabled="item.disabled"
                                @click="item.handler(item.disabled)">
                                {{item.name}}
                            </bk-button>
                        </li>
                    </ul>
                </bk-dropdown-menu>
                <cmdb-auth class="options-button sync-template-link"
                    v-show="withTemplate"
                    :auth="$authResources({ type: $OPERATION.U_SERVICE_INSTANCE })">
                    <bk-button slot-scope="{ disabled }"
                        class="topo-sync"
                        :disabled="disabled || !topoStatus"
                        @click="handleSyncTemplate">
                        <i class="bk-icon icon-refresh"></i>
                        {{$t('同步模板')}}
                        <span class="topo-status" v-show="topoStatus"></span>
                    </bk-button>
                </cmdb-auth>
                <div class="options-right fr">
                    <bk-checkbox class="options-checkbox"
                        :size="16"
                        v-model="isExpandAll"
                        @change="handleExpandAll">
                        <span class="checkbox-label">{{$t('全部展开')}}</span>
                    </bk-checkbox>
                    <div class="options-search">
                        <bk-search-select
                            ref="searchSelect"
                            :show-condition="false"
                            :placeholder="$t('请输入实例名称或选择标签')"
                            :data="searchSelect"
                            v-model="searchSelectData"
                            @menu-child-condition-select="handleConditionSelect"
                            @change="handleSearch">
                        </bk-search-select>
                    </div>
                </div>
            </div>
            <div class="tables">
                <service-instance-table
                    v-for="(instance, index) in instances"
                    ref="serviceInstanceTable"
                    :key="instance.id"
                    :instance="instance"
                    :expanded="index === 0"
                    :can-sync="topoStatus"
                    @create-process="handleCreateProcess"
                    @update-process="handleUpdateProcess"
                    @delete-instance="handleDeleteInstance"
                    @check-change="handleCheckChange">
                </service-instance-table>
            </div>
            <bk-pagination class="pagination" v-show="instances.length"
                align="right"
                size="small"
                :current="pagination.current"
                :count="pagination.count"
                :limit="pagination.size"
                @change="handlePageChange"
                @limit-change="handleSizeChange">
            </bk-pagination>
            <div class="filter-empty" v-if="!instances.length">
                <div class="filter-empty-content">
                    <img class="img-empty" src="../../../assets/images/empty-content.png" alt="">
                    <span>{{$t('暂无符合条件的实例')}}</span>
                </div>
            </div>
        </template>
        <service-instance-empty v-else>
        </service-instance-empty>
        <bk-sideslider
            v-transfer-dom
            :width="800"
            :title="processForm.title"
            :is-show.sync="processForm.show"
            :before-close="handleBeforeClose">
            <cmdb-form slot="content" v-if="processForm.show"
                ref="processForm"
                :type="processForm.type"
                :inst="processForm.instance"
                :disabled-properties="processForm.disabledProperties"
                :properties="processForm.properties"
                :property-groups="processForm.propertyGroups"
                @on-submit="handleSaveProcess"
                @on-cancel="handleBeforeClose">
                <template slot="bind_ip">
                    <cmdb-input-select
                        :disabled="checkDisabled"
                        :name="'bindIp'"
                        :placeholder="$t('请选择或输入IP')"
                        :options="processBindIp"
                        :validate="validateRules"
                        v-model="bindIp">
                    </cmdb-input-select>
                </template>
            </cmdb-form>
        </bk-sideslider>

        <bk-dialog class="bk-dialog-no-padding"
            v-model="editLabel.show"
            :mask-close="false"
            :width="580"
            @after-leave="handleSetEditBox">
            <div class="reset-header" slot="header">
                {{$t('批量编辑')}}
                <span>{{$tc('已选择实例', checked.length, { num: checked.length })}}</span>
            </div>
            <batch-edit-label ref="batchLabel"
                v-if="editLabel.visiable"
                :exisiting-label="editLabel.list">
                <cmdb-edit-label
                    ref="instanceLabel"
                    slot="batch-add-label"
                    class="edit-label"
                    :default-list="[]">
                </cmdb-edit-label>
            </batch-edit-label>
            <div class="edit-label-footer" slot="footer">
                <bk-button theme="primary" @click.stop="handleSubmitBatchLabel">{{$t('确定')}}</bk-button>
                <bk-button theme="default" class="ml5" @click.stop="handleCloseBatchLable">{{$t('取消')}}</bk-button>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import serviceInstanceTable from './service-instance-table.vue'
    import serviceInstanceEmpty from './service-instance-empty.vue'
    import batchEditLabel from './batch-edit-label.vue'
    import cmdbEditLabel from './edit-label.vue'
    export default {
        components: {
            serviceInstanceTable,
            serviceInstanceEmpty,
            batchEditLabel,
            cmdbEditLabel
        },
        data () {
            return {
                checked: [],
                isCheckAll: false,
                isExpandAll: false,
                filter: '',
                inSearch: false,
                instances: [],
                searchSelect: [
                    {
                        name: this.$t('服务实例名'),
                        id: 0
                    },
                    {
                        name: `${this.$t('标签')}(value)`,
                        id: 1,
                        children: [{
                            id: '',
                            name: ''
                        }],
                        conditions: []
                    },
                    {
                        name: `${this.$t('标签')}(key)`,
                        id: 2,
                        children: [{
                            id: '',
                            name: ''
                        }]
                    }
                ],
                searchSelectData: [],
                pagination: {
                    current: 1,
                    count: 0,
                    size: 10
                },
                processForm: {
                    type: 'create',
                    show: false,
                    title: '',
                    instance: null,
                    referenceService: null,
                    disabledProperties: [],
                    properties: [],
                    propertyGroups: [],
                    unwatch: null
                },
                editLabel: {
                    show: false,
                    visiable: false,
                    list: []
                },
                topoStatus: false,
                historyLabels: {},
                processBindIp: [],
                bindIp: '',
                templates: [],
                hasInitFilter: false,
                needRefresh: false,
                request: {
                    property: Symbol('property'),
                    propertyGroups: Symbol('propertyGroups'),
                    instance: Symbol('instance'),
                    label: Symbol('label')
                }
            }
        },
        computed: {
            targetInstanceName () {
                return this.$route.query.instanceName || ''
            },
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            currentNode () {
                return this.$store.state.businessHost.selectedNode
            },
            isModuleNode () {
                return this.currentNode && this.currentNode.data.bk_obj_id === 'module'
            },
            withTemplate () {
                return this.isModuleNode && this.currentNode && this.currentNode.data.service_template_id
            },
            currentModule () {
                if (this.currentNode && this.currentNode.data.bk_obj_id === 'module') {
                    return this.$store.state.businessHost.selectedNodeInstance
                }
                return null
            },
            processTemplateMap () {
                return this.$store.state.businessHost.processTemplateMap
            },
            menuItem () {
                return [{
                    name: this.$t('批量删除'),
                    handler: this.batchDelete,
                    disabled: !this.checked.length,
                    auth: 'D_SERVICE_INSTANCE'
                }, {
                    name: this.$t('复制IP'),
                    handler: this.copyIp,
                    disabled: !this.checked.length
                }, {
                    name: this.$t('编辑标签'),
                    handler: this.handleShowBatchLabel,
                    disabled: !this.checked.length,
                    auth: 'U_SERVICE_INSTANCE'
                }]
            },
            bindIpProperty () {
                return this.processForm.properties.find(property => property['bk_property_id'] === 'bind_ip') || {}
            },
            validateRules () {
                const rules = {}
                if (this.bindIpProperty.isrequired) {
                    rules['required'] = true
                }
                rules['regex'] = this.bindIpProperty.option
                return rules
            },
            checkDisabled () {
                const property = this.bindIpProperty
                if (this.processForm.type === 'create') {
                    return false
                }
                return !property.editable || property.isreadonly || this.processForm.disabledProperties.includes('bind_ip')
            }
        },
        watch: {
            async currentNode (node) {
                if (node && node.data.bk_obj_id === 'module') {
                    this.getData()
                }
            },
            bindIp (value) {
                this.$refs.processForm.values.bind_ip = value
            },
            checked () {
                this.isCheckAll = (this.checked.length === this.instances.length) && this.checked.length !== 0
            },
            searchSelectData (searchSelectData) {
                if (!searchSelectData.length && this.inSearch) this.inSearch = false
            }
        },
        async created () {
            await this.getHistoryLabel()
            if (this.targetInstanceName) {
                this.hasInitFilter = true
                this.searchSelectData.push({
                    'name': '服务实例名',
                    'id': 0,
                    'values': [{
                        'id': this.targetInstanceName,
                        'name': this.targetInstanceName
                    }]
                })
                this.searchSelect.shift()
            }
            this.getProcessProperties()
            this.getProcessPropertyGroups()
        },
        methods: {
            refresh () {
                this.inSearch = false
                this.getData()
            },
            async getData () {
                this.needRefresh = false
                const node = this.currentNode
                if (!this.hasInitFilter) {
                    this.searchSelectData = []
                }
                this.pagination.current = 1
                await this.getServiceInstances()
                if (this.withTemplate) {
                    this.getTemplate(node.data.service_template_id)
                }
                if (this.instances.length) {
                    this.getServiceInstanceDifferences()
                }
            },
            async getProcessProperties () {
                try {
                    const action = 'objectModelProperty/searchObjectAttribute'
                    this.processForm.properties = await this.$store.dispatch(action, {
                        params: {
                            bk_obj_id: 'process',
                            bk_supplier_account: this.$store.getters.supplierAccount
                        },
                        config: {
                            requestId: this.request.property,
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
                            requestId: this.request.propertyGroups,
                            fromCache: true
                        }
                    })
                } catch (e) {
                    this.processForm.propertyGroups = []
                    console.error(e)
                }
            },
            async getServiceInstanceDifferences () {
                try {
                    const data = await this.$store.dispatch('businessSynchronous/searchServiceInstanceDifferences', {
                        params: this.$injectMetadata({
                            bk_module_ids: [this.currentNode.data.bk_inst_id],
                            service_template_id: this.withTemplate
                        }, { injectBizId: true })
                    })
                    const difference = data.find(difference => difference.bk_module_id === this.currentNode.data.bk_inst_id)
                    this.topoStatus = !!difference && difference.has_difference
                } catch (error) {
                    console.error(error)
                }
            },
            async getServiceInstances () {
                try {
                    const searchKey = this.searchSelectData.find(item => (item.id === 0 && item.hasOwnProperty('values'))
                        || (![0, 1].includes(item.id) && !item.hasOwnProperty('values')))
                    const data = await this.$store.dispatch('serviceInstance/getModuleServiceInstances', {
                        params: this.$injectMetadata({
                            bk_module_id: this.currentNode.data.bk_inst_id,
                            with_name: true,
                            page: {
                                start: (this.pagination.current - 1) * this.pagination.size,
                                limit: this.pagination.size
                            },
                            search_key: searchKey
                                ? searchKey.hasOwnProperty('values') ? searchKey.values[0].name : searchKey.name
                                : '',
                            selectors: this.getSelectorParams()
                        }, { injectBizId: true }),
                        config: {
                            requestId: this.request.instance,
                            cancelPrevious: true
                        }
                    })
                    if (data.count && !data.info.length) {
                        this.pagination.current -= 1
                        this.getServiceInstances()
                    }
                    this.hasInitFilter = false
                    this.checked = []
                    this.isCheckAll = false
                    this.isExpandAll = false
                    this.instances = data.info
                    this.pagination.count = data.count
                } catch (e) {
                    console.error(e)
                    this.instances = []
                }
            },
            getSelectorParams () {
                try {
                    const labels = this.searchSelectData.filter(item => item.id === 1 && item.hasOwnProperty('values'))
                    const labelsKey = this.searchSelectData.filter(item => item.id === 2 && item.hasOwnProperty('values'))
                    const submitLabel = {}
                    const submitLabelKey = {}
                    labels.forEach(label => {
                        const conditionId = label.condition.id
                        if (!submitLabel[conditionId]) {
                            submitLabel[conditionId] = [label.values[0].id]
                        } else {
                            if (submitLabel[conditionId].indexOf(label.values[0].id) < 0) {
                                submitLabel[conditionId].push(label.values[0].id)
                            }
                        }
                    })
                    labelsKey.forEach(label => {
                        const id = label.values[0].id
                        if (!submitLabelKey[id]) {
                            submitLabelKey[id] = id
                        }
                    })
                    const selectors = Object.keys(submitLabel).map(key => {
                        return {
                            key: key,
                            operator: 'in',
                            values: submitLabel[key]
                        }
                    })
                    const selectorsKey = Object.keys(submitLabelKey).map(key => {
                        return {
                            key: key,
                            operator: 'exists',
                            values: []
                        }
                    })
                    return selectors.concat(selectorsKey)
                } catch (e) {
                    console.error(e)
                    return []
                }
            },
            async getHistoryLabel () {
                try {
                    const historyLabels = await this.$store.dispatch('instanceLabel/getHistoryLabel', {
                        params: this.$injectMetadata({}, { injectBizId: true }),
                        config: {
                            requestId: this.request.label,
                            cancelPrevious: true
                        }
                    })
                    this.historyLabels = historyLabels
                    const keys = Object.keys(historyLabels)
                    const valueOption = keys.map(key => {
                        return {
                            name: key + ' : ',
                            id: key
                        }
                    })
                    const keyOption = keys.map(key => {
                        return {
                            name: key,
                            id: key
                        }
                    })
                    const notRender = this.searchSelect[1].disabled
                    this.$set(this.searchSelect[1], 'disabled', !valueOption.length)
                    this.$set(this.searchSelect[2], 'disabled', !keyOption.length)
                    this.$set(this.searchSelect[1], 'conditions', valueOption)
                    this.$set(this.searchSelect[2], 'children', keyOption)
                    if (notRender && this.$refs.searchSelect) {
                        this.$refs.searchSelect.$forceUpdate()
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            handleSearch () {
                this.inSearch = true
                const instanceName = this.searchSelectData.filter(item => (item.id === 0 && item.hasOwnProperty('values'))
                    || (![0, 1].includes(item.id) && !item.hasOwnProperty('values')))
                if (instanceName.length) {
                    this.searchSelect[0].id === 0 && this.searchSelect.shift()
                } else {
                    this.searchSelect[0].id !== 0 && this.searchSelect.unshift({
                        name: this.$t('服务实例名'),
                        id: 0
                    })
                }
                if (instanceName.length >= 2) {
                    this.searchSelectData.pop()
                    this.$bkMessage({
                        message: this.$t('服务实例名重复'),
                        theme: 'warning'
                    })
                    return
                }
                this.handlePageChange(1)
            },
            handleClearFilter () {
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
                this.getInstanceIpByHost(referenceService.instance.bk_host_id)
                this.processForm.referenceService = referenceService
                this.processForm.type = 'create'
                this.processForm.title = `${this.$t('添加进程')}(${referenceService.instance.name})`
                this.processForm.instance = {}
                this.processForm.show = true
                this.$nextTick(() => {
                    this.bindIp = ''
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
                this.getInstanceIpByHost(processInstance.relation.bk_host_id)
                this.processForm.referenceService = referenceService
                this.processForm.type = 'update'
                this.processForm.title = this.$t('编辑进程')
                this.processForm.instance = processInstance.property
                this.processForm.show = true
                this.$nextTick(() => {
                    this.bindIp = this.$tools.getInstFormValues(this.processForm.properties, processInstance.property)['bind_ip']
                })

                const processTemplateId = processInstance.relation.process_template_id
                if (processTemplateId) {
                    const template = await this.getProcessTemplate(processTemplateId)
                    const disabledProperties = []
                    Object.keys(template).forEach(propertyId => {
                        const value = template[propertyId]
                        if (value.as_default_value) {
                            disabledProperties.push(propertyId)
                        }
                    })
                    this.processForm.disabledProperties = disabledProperties
                }
            },
            async getInstanceIpByHost (hostId) {
                try {
                    const instanceIpMap = this.$store.state.businessHost.instanceIpMap
                    let res = null
                    if (instanceIpMap.hasOwnProperty(hostId)) {
                        res = instanceIpMap[hostId]
                    } else {
                        res = await this.$store.dispatch('serviceInstance/getInstanceIpByHost', {
                            hostId,
                            config: {
                                requestId: 'getInstanceIpByHost'
                            }
                        })
                        this.$store.commit('businessHost/setInstanceIp', { hostId, res })
                    }
                    this.processBindIp = res.options.map(ip => {
                        return {
                            id: ip,
                            name: ip
                        }
                    })
                } catch (e) {
                    this.processBindIp = []
                    console.error(e)
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
                this.$store.commit('businessHost/setProcessTemplate', {
                    id: processTemplateId,
                    template: data.property
                })
                return Promise.resolve(data.property)
            },
            handleDeleteInstance (id) {
                this.getServiceInstances()
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
                    this.processForm.disabledProperties = []
                    this.$success(this.$t('保存成功'))
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
                    }, { injectBizId: true })
                })
            },
            updateProcess (values, instance) {
                return this.$store.dispatch('processInstance/updateServiceInstanceProcess', {
                    params: this.$injectMetadata({
                        processes: [{ ...instance, ...values }]
                    }, { injectBizId: true })
                })
            },
            handleCloseProcessForm () {
                this.processForm.show = false
                this.processForm.referenceService = null
                this.processForm.instance = null
                this.processForm.disabledProperties = []
            },
            handleBeforeClose () {
                const changedValues = this.$refs.processForm.changedValues
                if (Object.keys(changedValues).length) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
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
                        title: this.currentNode.data.bk_inst_name
                    }
                })
            },
            handleCheckALL (checked) {
                this.searchSelectData = []
                this.isCheckAll = checked
                this.$refs.serviceInstanceTable.forEach(table => {
                    table.checked = checked
                })
            },
            handleExpandAll (expanded) {
                this.searchSelectData = []
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
                    title: this.$t('确认删除实例'),
                    subTitle: this.$t('即将删除选中的实例', { count: this.checked.length }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        try {
                            const serviceInstanceIds = this.checked.map(instance => instance.id)
                            const deleteNum = serviceInstanceIds.length
                            await this.$store.dispatch('serviceInstance/deleteServiceInstance', {
                                config: {
                                    data: this.$injectMetadata({
                                        service_instance_ids: serviceInstanceIds
                                    }, { injectBizId: true }),
                                    requestId: 'batchDeleteServiceInstance'
                                }
                            })
                            this.currentNode.data.service_instance_count = this.currentNode.data.service_instance_count - deleteNum
                            this.currentNode.parents.forEach(node => {
                                node.data.service_instance_count = node.data.service_instance_count - deleteNum
                            })
                            this.$success(this.$t('删除成功'))
                            this.getServiceInstances()
                            this.checked = []
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            copyIp () {
                this.$copyText(this.checked.map(instance => instance.name.split('_')[0]).join('\n')).then(() => {
                    this.$success(this.$t('复制成功'))
                }, () => {
                    this.$error(this.$t('复制失败'))
                })
            },
            handleSyncTemplate () {
                this.$router.push({
                    name: 'syncServiceFromModule',
                    params: {
                        moduleId: this.currentNode.data.bk_inst_id,
                        setId: this.currentNode.parent.data.bk_inst_id
                    },
                    query: {
                        path: [...this.currentNode.parents, this.currentNode].map(node => node.name).join(' / ')
                    }
                })
            },
            handleShowBatchLabel (disabled) {
                if (disabled) {
                    return false
                }
                try {
                    this.editLabel.show = true
                    this.editLabel.visiable = true
                    const labelList = []
                    const existingKeys = []
                    for (const instance of this.checked) {
                        const labels = instance.labels
                        labels && Object.keys(labels).forEach(key => {
                            const index = existingKeys.findIndex(exisitingKey => exisitingKey === key)
                            if (index !== -1 && labels[key] === labelList[index].value) {
                                labelList[index].instanceIds.push(instance.id)
                            } else {
                                labelList.push({
                                    instanceIds: [instance.id],
                                    key: key,
                                    value: labels[key]
                                })
                                existingKeys.push(key)
                            }
                        })
                    }
                    this.editLabel.list = labelList
                } catch (e) {
                    console.error(e)
                    this.editLabel.list = []
                }
            },
            async handleSubmitBatchLabel () {
                try {
                    let status = ''
                    const validator = this.$refs.instanceLabel.$validator
                    const list = this.$refs.instanceLabel.submitList
                    if (list.length && !await validator.validateAll()) {
                        return
                    }
                    const removeList = this.$refs.batchLabel.removeList
                    const removeKeys = []
                    const instanceIds = []
                    removeList.forEach(label => {
                        removeKeys.push(label.key)
                        instanceIds.push(...label.instanceIds)
                    })
                    if (removeList.length) {
                        status = await this.$store.dispatch('instanceLabel/deleteInstanceLabel', {
                            config: {
                                data: this.$injectMetadata({
                                    instance_ids: instanceIds,
                                    keys: removeKeys
                                }, { injectBizId: true }),
                                requestId: 'deleteInstanceLabel',
                                transformData: false
                            }
                        })
                    }
                    if (list.length) {
                        const serviceInstanceIds = this.checked.map(instance => instance.id)
                        const labelSet = {}
                        list.forEach(label => {
                            labelSet[label.key] = label.value
                        })
                        status = await this.$store.dispatch('instanceLabel/createInstanceLabel', {
                            params: this.$injectMetadata({
                                instance_ids: serviceInstanceIds,
                                labels: labelSet
                            }, { injectBizId: true }),
                            config: {
                                requestId: 'createInstanceLabel',
                                transformData: false
                            }
                        })
                    }
                    if (status && status.bk_error_msg === 'success') {
                        this.$success(this.$t('保存成功'))
                        this.searchSelectData = []
                        this.getServiceInstances()
                        this.handleCheckALL(false)
                        this.getHistoryLabel()
                    }
                    this.handleCloseBatchLable()
                    setTimeout(() => {
                        this.handleSetEditBox()
                    }, 200)
                } catch (e) {
                    console.error(e)
                }
            },
            handleCloseBatchLable () {
                this.editLabel.show = false
            },
            handleSetEditBox () {
                this.editLabel.list = []
                this.editLabel.visiable = false
            },
            handleConditionSelect (cur, index) {
                const values = this.historyLabels[cur.id]
                const children = values.map(item => {
                    return {
                        id: item,
                        name: item
                    }
                })
                const el = this.$refs.searchSelect
                el.curItem.children = children
                el.updateChildMenu(cur, index, false)
                el.showChildMenu(children)
            },
            async getTemplate (id) {
                try {
                    const data = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: this.$injectMetadata({
                            service_template_id: id
                        }, { injectBizId: true }),
                        config: {
                            requestId: 'getBatchProcessTemplate'
                        }
                    })
                    this.templates = data.info
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .service-layout {
        height: 100%;
        padding: 14px 0 0 0;
    }
    .options {
        padding: 0 0 15px;
    }
    .options-button {
        height: 32px;
        margin: 0 0 0 6px;
        line-height: 30px;
    }
    .topo-sync {
        position: relative;
        .topo-status {
            position: absolute;
            top: -4px;
            right: -4px;
            width: 8px;
            height: 8px;
            background-color: #ea3636;
            border-radius: 50%;
        }
    }
    .options-checkall {
        width: 36px;
        height: 32px;
        line-height: 30px;
        padding: 0 9px;
        text-align: center;
        border: 1px solid #f0f1f5;
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
        min-width: 240px;
        max-width: 280px;
        height: 34px;
        z-index: 99;
        .bk-search-select {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
        }
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
            cursor: pointer;
            @include ellipsis;
            .item-btn {
                display: block;
                width: 100%;
                padding: 0 15px;
                height: 40px;
                line-height: 40px;
                color: #737987;
                text-align: left;
                &:disabled {
                    color: #dcdee5;
                }
                &:not(.is-disabled):hover {
                    background-color: #ebf4ff;
                    color: #3c96ff;
                }
            }
        }
    }
    .sync-template-link {
        display: inline-block;
        position: relative;
        margin-left: 18px;
        padding: 0;
        &::before {
            content: '';
            position: absolute;
            top: 7px;
            left: -11px;;
            width: 1px;
            height: 20px;
            background-color: #dcdee5;
        }
        .icon-refresh {
            top: -1px;
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
            .img-empty {
                display: block;
                margin: 0 auto;
            }
        }
    }
    .reset-header {
        text-align: left;
        span {
            color: #979ba5;
            font-size: 14px;
        }
    }
    .edit-label-footer {
        .bk-button {
            min-width: 76px;
        }
    }
</style>
