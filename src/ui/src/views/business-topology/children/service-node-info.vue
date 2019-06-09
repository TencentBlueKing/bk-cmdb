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
            :show-delete="false"
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="flattenedInstance"
            :show-options="modelId !== 'biz'"
            @on-edit="handleEdit">
            <span class="property-value fl" slot="__template_name__">
                {{flattenedInstance.__template_name__}}
                <bk-button v-if="withTemplate" @click="handleRemoveTemplate">{{$t('BusinessTopology["解除模板"]')}}</bk-button>
            </span>
        </cmdb-details>
        <cmdb-form class="topology-details" v-else-if="type === 'update'"
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="instance"
            :type="type"
            @on-submit="handleSubmit"
            @on-cancel="handleCancel">
            <template slot="extra-options">
                <bk-button type="danger" style="margin-left: 4px" @click="handleDelete">{{$t('Common["删除"]')}}
                </bk-button>
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
                propertyGroups: [],
                instance: {}
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
            selectedNode () {
                return this.$store.state.businessTopology.selectedNode
            },
            modelId () {
                if (this.selectedNode) {
                    return this.selectedNode.data.bk_obj_id
                }
                return null
            },
            withTemplate () {
                return this.instance.service_template_id && this.instance.service_template_id !== 2
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
            selectedNode (node) {
                if (node) {
                    this.type = 'details'
                    this.getInstance()
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
                    bk_property_name: this.$t('BusinessTopology["模板名称"]'),
                    bk_property_group: group.bk_group_id,
                    bk_property_index: 1
                }, {
                    bk_property_id: '__service_category__',
                    bk_property_name: this.$t('BusinessTopology["服务分类"]'),
                    bk_property_group: group.bk_group_id,
                    bk_property_index: 2
                }]
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
                    bk_group_index: Infinity,
                    bk_group_name: this.$t('BusinessTopology["服务模板信息"]'),
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
                if (instance.service_template_id !== 2) {
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
                this.type = 'update'
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
                    await (promiseMap[this.modelId] || this.updateCustomInstance)(value)
                    this.selectedNode.data.bk_inst_name = value[nameMap[this.modelId] || 'bk_inst_name']
                    this.instance = Object.assign({}, this.instance, value)
                    this.type = 'details'
                    this.$success(this.$t('Common["修改成功"]'))
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
                    title: `${this.$t('Common["确定删除"]')} ${this.selectedNode.name}?`,
                    content: this.modelId === 'module'
                        ? this.$t('Common["请先转移其下所有的主机"]')
                        : this.$t('Common[\'下属层级都会被删除，请先转移其下所有的主机\']'),
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
                            tree.setSelected(parentId, true, true)
                            tree.removeNode(nodeId)
                            this.$success(this.$t('Common[\'删除成功\']'))
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
                        requestId: 'deleteNodeInstance'
                    }
                })
            },
            deleteModuleInstance () {
                return this.$store.dispatch('objectModule/deleteModule', {
                    bizId: this.business,
                    setId: this.selectedNode.parent.data.bk_inst_id,
                    moduleId: this.selectedNode.data.bk_inst_id,
                    config: {
                        requestId: 'deleteNodeInstance'
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
                this.$bkInfo({
                    title: this.$t('BusinessTopology["确认解除模板"]'),
                    content: this.$t('BusinessTopology["解除模板影响"]'),
                    confirmFn: async () => {
                        await this.$store.dispatch('serviceInstance/removeServiceTemplate', {
                            config: {
                                requestId: 'removeServiceTemplate',
                                data: this.$injectMetadata({
                                    bk_module_id: this.instance.bk_module_id
                                })
                            }
                        })
                        this.instance.service_template_id = null
                        this.instance.__template_name__ = '--'
                    }
                })
            }
        }
    }
</script>
