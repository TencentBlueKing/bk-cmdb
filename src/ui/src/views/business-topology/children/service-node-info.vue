<template>
    <div class="node-info" style="height: 100%;"
        v-bkloading="{
            isLoading: $loading([
                'getModelProperties',
                'getModelPropertyGroups',
                'getNodeInstance',
                'updateNodeInstance',
                'deleteNodeInstance',
                'removeServiceTemplate'
            ])
        }"
    >
        <cmdb-details class="topology-details"
            v-if="type === 'details'"
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="flattenedInstance"
            :show-options="modelId !== 'biz' && !isBlueking">
            <template slot="details-options">
                <span style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.U_TOPO),
                        auth: [$OPERATION.U_TOPO]
                    }">
                    <bk-button class="button-edit"
                        theme="primary"
                        :disabled="!$isAuthorized($OPERATION.U_TOPO)"
                        @click="handleEdit">
                        {{$t('编辑')}}
                    </bk-button>
                </span>
                <span style="display: inline-block;"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.D_TOPO),
                        auth: [$OPERATION.D_TOPO]
                    }">
                    <bk-button class="btn-delete"
                        :disabled="!$isAuthorized($OPERATION.D_TOPO)"
                        @click="handleDelete">
                        {{$t('删除节点')}}
                    </bk-button>
                </span>
            </template>
            <span class="property-value fl" slot="__template_name__">
                <span class="link"
                    v-if="withTemplate"
                    @click="goServiceTemplate">{{flattenedInstance.__template_name__}}</span>
                <span v-else>{{flattenedInstance.__template_name__}}</span>
                <!-- <span style="display: inline-block;"
                    v-if="withTemplate"
                    v-cursor="{
                        active: !$isAuthorized($OPERATION.U_TOPO),
                        auth: [$OPERATION.U_TOPO]
                    }">
                    <bk-button class="unbind-button"
                        :disabled="!$isAuthorized($OPERATION.U_TOPO)"
                        @click="handleRemoveTemplate">
                        {{$t('解除模板')}}
                    </bk-button>
                </span> -->
            </span>
        </cmdb-details>
        <cmdb-form class="topology-form" v-else-if="type === 'update'"
            ref="form"
            :properties="properties"
            :property-groups="propertyGroups"
            :disabled-properties="disabledProperties"
            :inst="instance"
            :type="type"
            @on-submit="handleSubmit"
            @on-cancel="handleCancel">
            <template slot="__service_category__" v-if="!withTemplate">
                <cmdb-selector class="category-selector fl"
                    :list="firstCategories"
                    v-model="first">
                </cmdb-selector>
                <cmdb-selector class="category-selector fl"
                    :list="secondCategories"
                    v-model="second"
                    @on-selected="handleChangeCategory">
                </cmdb-selector>
            </template>
        </cmdb-form>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                type: 'details',
                properties: [],
                disabledProperties: [],
                propertyGroups: [],
                instance: {},
                first: '',
                second: ''
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            propertyMap () {
                return this.$store.state.businessTopology.propertyMap
            },
            propertyGroupMap () {
                return this.$store.state.businessTopology.propertyGroupMap
            },
            categoryMap () {
                return this.$store.state.businessTopology.categoryMap
            },
            firstCategories () {
                return this.categoryMap[this.business] || []
            },
            secondCategories () {
                const firstCategory = this.firstCategories.find(category => category.id === this.first) || {}
                return firstCategory.secondCategory || []
            },
            selectedNode () {
                return this.$store.state.businessTopology.selectedNode
            },
            isModuleNode () {
                return this.selectedNode && this.selectedNode.data.bk_obj_id === 'module'
            },
            isBlueking () {
                let rootNode = this.selectedNode || {}
                if (rootNode.parent) {
                    rootNode = rootNode.parents[0]
                }
                return rootNode.name === '蓝鲸'
            },
            modelId () {
                if (this.selectedNode) {
                    return this.selectedNode.data.bk_obj_id
                }
                return null
            },
            withTemplate () {
                return this.isModuleNode && !!this.instance.service_template_id
            },
            flattenedInstance () {
                return this.$tools.flattenItem(this.properties, this.instance)
            }
        },
        watch: {
            modelId (modelId) {
                if (modelId) {
                    this.type = 'details'
                    this.init()
                }
            },
            async selectedNode (node) {
                if (node) {
                    this.type = 'details'
                    await this.getInstance()
                    this.disabledProperties = node.data.bk_obj_id === 'module' && this.withTemplate ? ['bk_module_name'] : []
                }
            }
        },
        methods: {
            async init () {
                try {
                    const [
                        properties,
                        groups
                    ] = await Promise.all([
                        this.getProperties(),
                        this.getPropertyGroups()
                    ])
                    this.properties = properties
                    this.propertyGroups = groups
                } catch (e) {
                    console.error(e)
                    this.properties = []
                    this.propertyGroups = []
                }
            },
            async getProperties () {
                let properties = []
                const modelId = this.modelId
                if (this.propertyMap.hasOwnProperty(modelId)) {
                    properties = this.propertyMap[modelId]
                } else {
                    const action = 'objectModelProperty/searchObjectAttribute'
                    properties = await this.$store.dispatch(action, {
                        params: this.$injectMetadata({
                            bk_obj_id: modelId,
                            bk_supplier_account: this.$store.getters.supplierAccount
                        }),
                        config: {
                            requestId: 'getModelProperties'
                        }
                    })
                    if (modelId === 'module') {
                        properties.push(...this.getModuleServiceTemplateProperties())
                    }
                    this.$store.commit('businessTopology/setProperties', {
                        id: modelId,
                        properties: properties
                    })
                }
                return Promise.resolve(properties)
            },
            getModuleServiceTemplateProperties () {
                const group = this.getModuleServiceTemplateGroup()
                return [{
                    bk_property_id: '__template_name__',
                    bk_property_name: this.$t('模板名称'),
                    bk_property_group: group.bk_group_id,
                    bk_property_index: 1,
                    bk_isapi: false,
                    editable: false,
                    unit: ''
                }, {
                    bk_property_id: '__service_category__',
                    bk_property_name: this.$t('服务分类'),
                    bk_property_group: group.bk_group_id,
                    bk_property_index: 2,
                    bk_isapi: false,
                    editable: false,
                    unit: ''
                }]
            },
            updateCategoryProperty (state) {
                const serviceCategoryProperty = this.properties.find(property => property.bk_property_id === '__service_category__') || {}
                Object.assign(serviceCategoryProperty, state)
            },
            async getPropertyGroups () {
                let groups = []
                const modelId = this.modelId
                if (this.propertyGroupMap.hasOwnProperty(modelId)) {
                    groups = this.propertyGroupMap[modelId]
                } else {
                    const action = 'objectModelFieldGroup/searchGroup'
                    groups = await this.$store.dispatch(action, {
                        objId: modelId,
                        params: this.$injectMetadata(),
                        config: {
                            requestId: 'getModelPropertyGroups'
                        }
                    })
                    if (modelId === 'module') {
                        groups.push(this.getModuleServiceTemplateGroup())
                    }
                    this.$store.commit('businessTopology/setPropertyGroups', {
                        id: modelId,
                        groups: groups
                    })
                }
                return Promise.resolve(groups)
            },
            getModuleServiceTemplateGroup () {
                return {
                    bk_group_id: '__service_template_info__',
                    bk_group_index: -1,
                    bk_group_name: this.$t('服务模板信息'),
                    bk_obj_id: 'module',
                    ispre: true
                }
            },
            async getInstance () {
                try {
                    const modelId = this.modelId
                    const promiseMap = {
                        biz: this.getBizInstance,
                        set: this.getSetInstance,
                        module: this.getModuleInstance
                    }
                    this.instance = await (promiseMap[modelId] || this.getCustomInstance)()
                    this.$store.commit('businessTopology/setSelectedNodeInstance', this.instance)
                } catch (e) {
                    console.error(e)
                    this.instance = {}
                }
            },
            async getBizInstance () {
                const data = await this.$store.dispatch('objectBiz/searchBusiness', {
                    params: {
                        page: { start: 0, limit: 1 },
                        fields: [],
                        condition: {
                            bk_biz_id: { $eq: this.selectedNode.data.bk_inst_id }
                        }
                    },
                    config: {
                        requestId: 'getNodeInstance',
                        cancelPrevious: true
                    }
                })
                return data.info[0]
            },
            async getSetInstance () {
                const data = await this.$store.dispatch('objectSet/searchSet', {
                    bizId: this.business,
                    params: {
                        page: { start: 0, limit: 1 },
                        fields: [],
                        condition: {
                            bk_set_id: this.selectedNode.data.bk_inst_id
                        }
                    },
                    config: {
                        requestId: 'getNodeInstance',
                        cancelPrevious: true
                    }
                })
                return data.info[0]
            },
            async getModuleInstance () {
                const data = await this.$store.dispatch('objectModule/searchModule', {
                    bizId: this.business,
                    setId: this.selectedNode.parent.data.bk_inst_id,
                    params: {
                        page: { start: 0, limit: 1 },
                        fields: [],
                        condition: {
                            bk_module_id: this.selectedNode.data.bk_inst_id,
                            bk_supplier_account: this.$store.getters.supplierAccount
                        }
                    },
                    config: {
                        requestId: 'getNodeInstance',
                        cancelPrevious: true
                    }
                })
                const instance = data.info[0]
                const serviceInfo = await this.getServiceInfo(instance)
                return {
                    ...instance,
                    ...serviceInfo
                }
            },
            async getServiceInfo (instance) {
                const serviceInfo = {}
                if (instance.service_template_id) {
                    serviceInfo.__template_name__ = instance.bk_module_name
                }
                const categories = await this.getServiceCategories()
                const firstCategory = categories.find(category => category.secondCategory.some(second => second.id === instance.service_category_id)) || {}
                const secondCategory = (firstCategory.secondCategory || []).find(second => second.id === instance.service_category_id) || {}
                serviceInfo.__service_category__ = `${firstCategory.name || '--'} / ${secondCategory.name || '--'}`
                return serviceInfo
            },
            async getServiceCategories () {
                if (this.categoryMap.hasOwnProperty(this.business)) {
                    return this.categoryMap[this.business]
                } else {
                    try {
                        const data = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
                            params: this.$injectMetadata()
                        })
                        const categories = this.collectServiceCategories(data.info)
                        this.$store.commit('businessTopology/setCategories', {
                            id: this.business,
                            categories: categories
                        })
                        return categories
                    } catch (e) {
                        console.error(e)
                        return []
                    }
                }
            },
            collectServiceCategories (data) {
                const categories = []
                data.forEach(item => {
                    if (!item.category.bk_parent_id) {
                        categories.push(item.category)
                    }
                })
                categories.forEach(category => {
                    category.secondCategory = data.filter(item => item.category.bk_parent_id === category.id).map(item => item.category)
                })
                return categories
            },
            async getCustomInstance () {
                const data = await this.$store.dispatch('objectCommonInst/searchInst', {
                    objId: this.modelId,
                    params: {
                        page: { start: 0, limit: 1 },
                        fields: {},
                        condition: {
                            [this.modelId]: [{
                                field: 'bk_inst_id',
                                operator: '$eq',
                                value: this.selectedNode.data.bk_inst_id
                            }]
                        }
                    },
                    config: {
                        requestId: 'getNodeInstance',
                        cancelPrevious: true
                    }
                })
                return data.info[0]
            },
            handleEdit () {
                if (this.modelId === 'module') {
                    if (!this.withTemplate) {
                        const second = this.instance.service_category_id
                        const firstCategory = this.firstCategories.find(({ secondCategory }) => {
                            return secondCategory.some(category => category.id === second)
                        })
                        this.first = firstCategory.id
                        this.second = second
                    }
                    this.updateCategoryProperty({
                        editable: !this.withTemplate
                    })
                }
                this.type = 'update'
            },
            handleChangeCategory (id, category) {
                this.$set(this.$refs.form.values, 'service_category_id', id)
            },
            async handleSubmit (value) {
                const promiseMap = {
                    set: this.updateSetInstance,
                    module: this.updateModuleInstance
                }
                const nameMap = {
                    set: 'bk_set_name',
                    module: 'bk_module_name'
                }
                try {
                    await (promiseMap[this.modelId] || this.updateCustomInstance)(this.$injectMetadata(value))
                    this.selectedNode.data.bk_inst_name = value[nameMap[this.modelId] || 'bk_inst_name']
                    this.instance = Object.assign({}, this.instance, value)
                    this.type = 'details'
                    this.$success(this.$t('修改成功'))
                } catch (e) {
                    console.error(e)
                }
            },
            updateSetInstance (value) {
                return this.$store.dispatch('objectSet/updateSet', {
                    bizId: this.business,
                    setId: this.selectedNode.data.bk_inst_id,
                    params: { ...value },
                    config: {
                        requestId: 'updateNodeInstance'
                    }
                })
            },
            updateModuleInstance (value) {
                delete value.__template_name__
                delete value.__service_category__
                return this.$store.dispatch('objectModule/updateModule', {
                    bizId: this.business,
                    setId: this.selectedNode.parent.data.bk_inst_id,
                    moduleId: this.selectedNode.data.bk_inst_id,
                    params: {
                        bk_supplier_account: this.$store.getters.supplierAccount,
                        ...value
                    },
                    config: {
                        requestId: 'updateNodeInstance'
                    }
                }).then(async () => {
                    const serviceInfo = await this.getServiceInfo({ service_category_id: value.service_category_id || this.instance.service_category_id })
                    Object.assign(this.instance, serviceInfo)
                })
            },
            updateCustomInstance (value) {
                return this.$store.dispatch('objectCommonInst/updateInst', {
                    objId: this.modelId,
                    instId: this.selectedNode.data.bk_inst_id,
                    params: { ...value },
                    config: {
                        requestId: 'updateNodeInstance'
                    }
                })
            },
            handleCancel () {
                this.type = 'details'
            },
            handleDelete () {
                this.$bkInfo({
                    title: `${this.$t('确定删除')} ${this.selectedNode.name}?`,
                    subTitle: this.modelId === 'module'
                        ? this.$t('删除模块提示')
                        : this.$t('下属层级都会被删除，请先转移其下所有的主机'),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        const promiseMap = {
                            set: this.deleteSetInstance,
                            module: this.deleteModuleInstance
                        }
                        try {
                            await (promiseMap[this.modelId] || this.deleteCustomInstance)()
                            const tree = this.selectedNode.tree
                            const parentId = this.selectedNode.parent.id
                            const nodeId = this.selectedNode.id
                            tree.setSelected(parentId, {
                                emitEvent: true
                            })
                            tree.removeNode(nodeId)
                            this.$success(this.$t('删除成功'))
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            deleteSetInstance () {
                return this.$store.dispatch('objectSet/deleteSet', {
                    bizId: this.business,
                    setId: this.selectedNode.data.bk_inst_id,
                    config: {
                        requestId: 'deleteNodeInstance',
                        data: this.$injectMetadata({})
                    }
                })
            },
            deleteModuleInstance () {
                return this.$store.dispatch('objectModule/deleteModule', {
                    bizId: this.business,
                    setId: this.selectedNode.parent.data.bk_inst_id,
                    moduleId: this.selectedNode.data.bk_inst_id,
                    config: {
                        requestId: 'deleteNodeInstance',
                        data: this.$injectMetadata({})
                    }
                })
            },
            deleteCustomInstance () {
                return this.$store.dispatch('objectCommonInst/deleteInst', {
                    objId: this.modelId,
                    instId: this.selectedNode.data.bk_inst_id,
                    config: {
                        requestId: 'deleteNodeInstance',
                        data: this.$injectMetadata()
                    }
                })
            },
            handleRemoveTemplate () {
                const content = this.$createElement('div', {
                    style: {
                        'font-size': '14px'
                    },
                    domProps: {
                        innerHTML: this.$tc('解除模板影响', this.flattenedInstance.__template_name__, { name: this.flattenedInstance.__template_name__ })
                    }
                })
                this.$bkInfo({
                    title: this.$t('确认解除模板'),
                    subHeader: content,
                    confirmFn: async () => {
                        await this.$store.dispatch('serviceInstance/removeServiceTemplate', {
                            config: {
                                requestId: 'removeServiceTemplate',
                                data: this.$injectMetadata({
                                    bk_module_id: this.instance.bk_module_id
                                })
                            }
                        })
                        this.selectedNode.data.service_template_id = 0
                        this.instance.service_template_id = null
                        this.instance.__template_name__ = '--'
                        this.disabledProperties = []
                    }
                })
            },
            goServiceTemplate () {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: this.instance.service_template_id
                    },
                    query: {
                        from: {
                            name: this.$route.name,
                            query: {
                                module: this.instance.bk_module_id
                            }
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-details {
        /deep/ .details-options {
            padding: 28px 18px 0 0;
        }
    }
    .property-value {
        height: 26px;
        line-height: 26px;
        overflow: visible;
        .link {
            color: #3a84ff;
            cursor: pointer;
        }
    }
    .unbind-button {
        height: 26px;
        padding: 0 4px;
        margin: 0 0 0 6px;
        line-height: 24px;
        font-size: 12px;
        color: #63656E;
    }
    .topology-form {
        /deep/ .property-item {
            max-width: 554px !important;
        }
        .category-selector {
            width: calc(50% - 5px);
            & + .category-selector {
                margin-left: 10px;
            }
        }
    }
    .btn-delete{
        min-width: 76px;
        &:hover {
            color: #ff5656;
            border-color: #ff5656;
        }
    }
</style>
