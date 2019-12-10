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
        <div class="template-info mb10 clearfix" v-if="( isSetNode || isModuleNode) && type === 'details'">
            <template v-if="isModuleNode">
                <div class="info-item fl" :title="`${$t('服务模板')} : ${templateInfo.serviceTemplateName}`">
                    <span class="name fl">{{$t('服务模板')}}</span>
                    <div class="value fl">
                        <div class="template-value" v-if="withTemplate" @click="goServiceTemplate">
                            <span class="text link">{{templateInfo.serviceTemplateName}}</span>
                            <i class="icon-cc-share"></i>
                        </div>
                        <span class="text" v-else>{{templateInfo.serviceTemplateName}}</span>
                    </div>
                </div>
                <div class="info-item fl" :title="`${$t('服务分类')} : ${templateInfo.serviceCategory || '--'}`">
                    <span class="name fl">{{$t('服务分类')}}</span>
                    <div class="value fl">
                        <span class="text">{{templateInfo.serviceCategory || '--'}}</span>
                    </div>
                </div>
            </template>
            <template v-else-if="isSetNode">
                <div class="info-item fl" :title="`${$t('集群模板')} : ${templateInfo.setTemplateName}`">
                    <span class="name fl">{{$t('集群模板')}}</span>
                    <div class="value fl">
                        <template v-if="withSetTemplate">
                            <div class="template-value set-template fl" @click="goSetTemplate">
                                <span class="text link">{{templateInfo.setTemplateName}}</span>
                                <i class="icon-cc-share"></i>
                            </div>
                            <cmdb-auth :auth="$authResources({ type: $OPERATION.U_TOPO })">
                                <bk-button slot-scope="{ disabled }"
                                    :class="['sync-set-btn', 'ml5', { 'has-change': hasChange }]"
                                    :disabled="!hasChange || disabled"
                                    @click="handleSyncSetTemplate">
                                    {{$t('同步集群')}}
                                </bk-button>
                            </cmdb-auth>
                        </template>
                        <span class="text" v-else>{{templateInfo.setTemplateName}}</span>
                    </div>
                </div>
            </template>
        </div>
        <cmdb-details class="topology-details"
            v-if="type === 'details'"
            :properties="properties"
            :property-groups="propertyGroups"
            :inst="flattenedInstance"
            :show-options="modelId !== 'biz' && !isBlueking">
            <template slot="details-options">
                <cmdb-auth :auth="$authResources({ type: $OPERATION.U_TOPO })">
                    <template slot-scope="{ disabled }">
                        <bk-button class="button-edit"
                            theme="primary"
                            :disabled="disabled"
                            @click="handleEdit">
                            {{$t('编辑')}}
                        </bk-button>
                    </template>
                </cmdb-auth>
                <cmdb-auth :auth="$authResources({ type: $OPERATION.D_TOPO })">
                    <template slot-scope="{ disabled }">
                        <span class="inline-block-middle" v-if="moduleFromSetTemplate"
                            v-bk-tooltips="$t('由集群模板创建的模块无法删除')">
                            <bk-button class="btn-delete" disabled>
                                {{$t('删除节点')}}
                            </bk-button>
                        </span>
                        <bk-button class="btn-delete" v-else
                            theme="default"
                            :disabled="disabled"
                            @click="handleDelete">
                            {{$t('删除节点')}}
                        </bk-button>
                    </template>
                </cmdb-auth>
            </template>
        </cmdb-details>
        <template v-else-if="type === 'update'">
            <div class="service-category" v-if="!withTemplate && isModuleNode">
                <span class="title">{{$t('服务分类')}}</span>
                <div class="selector-item mt10 clearfix">
                    <cmdb-selector class="category-selector fl"
                        :list="firstCategories"
                        v-model="first"
                        @on-selected="handleChangeFirstCategory">
                    </cmdb-selector>
                    <cmdb-selector class="category-selector fl"
                        :list="secondCategories"
                        name="secondCategory"
                        v-validate="'required'"
                        v-model="second"
                        @on-selected="handleChangeCategory">
                    </cmdb-selector>
                    <span class="second-category-errors" v-if="errors.has('secondCategory')">{{errors.first('secondCategory')}}</span>
                </div>
            </div>
            <cmdb-form class="topology-form"
                ref="form"
                :properties="properties"
                :property-groups="propertyGroups"
                :disabled-properties="disabledProperties"
                :inst="instance"
                :type="type"
                @on-submit="handleSubmit"
                @on-cancel="handleCancel">
            </cmdb-form>
        </template>
    </div>
</template>

<script>
    import debounce from 'lodash.debounce'
    export default {
        props: {
            active: Boolean
        },
        data () {
            return {
                type: 'details',
                properties: [],
                disabledProperties: [],
                propertyGroups: [],
                instance: {},
                first: '',
                second: '',
                hasChange: false,
                templateInfo: {
                    serviceTemplateName: this.$t('无'),
                    serviceCategory: '',
                    setTemplateName: this.$t('无')
                },
                refresh: null
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            propertyMap () {
                return this.$store.state.businessHost.propertyMap
            },
            propertyGroupMap () {
                return this.$store.state.businessHost.propertyGroupMap
            },
            categoryMap () {
                return this.$store.state.businessHost.categoryMap
            },
            firstCategories () {
                return this.categoryMap[this.business] || []
            },
            secondCategories () {
                const firstCategory = this.firstCategories.find(category => category.id === this.first) || {}
                return firstCategory.secondCategory || []
            },
            selectedNode () {
                return this.$store.state.businessHost.selectedNode
            },
            isModuleNode () {
                return this.selectedNode && this.selectedNode.data.bk_obj_id === 'module'
            },
            isSetNode () {
                return this.selectedNode && this.selectedNode.data.bk_obj_id === 'set'
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
            withSetTemplate () {
                return this.isSetNode && !!this.instance.set_template_id
            },
            flattenedInstance () {
                return this.$tools.flattenItem(this.properties, this.instance)
            },
            moduleFromSetTemplate () {
                return this.isModuleNode && !!this.selectedNode.parent.data.set_template_id
            }
        },
        watch: {
            modelId: {
                immediate: true,
                handler (modelId) {
                    if (modelId && this.active) {
                        this.initProperties()
                    }
                }
            },
            selectedNode: {
                immediate: true,
                handler (node) {
                    if (node && this.active) {
                        this.getNodeDetails()
                    }
                }
            },
            active: {
                immediate: true,
                handler (active) {
                    if (active) {
                        this.refresh()
                    }
                }
            }
        },
        created () {
            this.refresh = debounce(() => {
                this.initProperties()
                this.getNodeDetails()
            }, 10)
        },
        methods: {
            async initProperties () {
                this.type = 'details'
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
            async getNodeDetails () {
                this.type = 'details'
                await this.getInstance()
                if (this.withSetTemplate) {
                    this.getDiffTemplateAndInstances()
                }
                this.disabledProperties = this.selectedNode.data.bk_obj_id === 'module' && this.withTemplate ? ['bk_module_name'] : []
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
                    this.$store.commit('businessHost/setProperties', {
                        id: modelId,
                        properties: properties
                    })
                }
                return Promise.resolve(properties)
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
                    this.$store.commit('businessHost/setPropertyGroups', {
                        id: modelId,
                        groups: groups
                    })
                }
                return Promise.resolve(groups)
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
                    this.$store.commit('businessHost/setSelectedNodeInstance', this.instance)
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
                const setTemplateId = data.info[0].set_template_id
                let setTemplateInfo = {}
                if (setTemplateId) {
                    setTemplateInfo = await this.getSetTemplateInfo(setTemplateId)
                }
                this.templateInfo.setTemplateName = setTemplateInfo.name || this.$t('无')
                return data.info[0]
            },
            getSetTemplateInfo (id) {
                return this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
                    bizId: this.business,
                    setTemplateId: id,
                    config: {
                        requestId: 'getSingleSetTemplateInfo'
                    }
                })
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
                this.getServiceInfo(instance)
                return instance
            },
            async getServiceInfo (instance) {
                this.templateInfo.serviceTemplateName = instance.service_template_id ? instance.bk_module_name : this.$t('无')
                const categories = await this.getServiceCategories()
                const firstCategory = categories.find(category => category.secondCategory.some(second => second.id === instance.service_category_id)) || {}
                const secondCategory = (firstCategory.secondCategory || []).find(second => second.id === instance.service_category_id) || {}
                this.templateInfo.serviceCategory = `${firstCategory.name || '--'} / ${secondCategory.name || '--'}`
            },
            async getServiceCategories () {
                if (this.categoryMap.hasOwnProperty(this.business)) {
                    return this.categoryMap[this.business]
                } else {
                    try {
                        const data = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
                            params: this.$injectMetadata({}, { injectBizId: true })
                        })
                        const categories = this.collectServiceCategories(data.info)
                        this.$store.commit('businessHost/setCategories', {
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
                }
                this.type = 'update'
            },
            handleChangeFirstCategory (id, category) {
                if (!this.secondCategories.length) {
                    this.second = ''
                    this.$set(this.$refs.form.values, 'service_category_id', '')
                }
            },
            handleChangeCategory (id, category) {
                this.$set(this.$refs.form.values, 'service_category_id', id)
            },
            async handleSubmit (value) {
                if (!await this.$validator.validateAll()) return
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
                    this.getServiceInfo(this.instance)
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
            async handleDelete () {
                const count = await this.getSelectedNodeHostCount()
                if (count) {
                    this.$error(this.$t('目标包含主机, 不允许删除'))
                    return
                }
                this.$bkInfo({
                    title: `${this.$t('确定删除')} ${this.selectedNode.name}?`,
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
                        innerHTML: this.$tc('解除模板影响', this.templateInfo.serviceTemplateName, { name: this.templateInfo.serviceTemplateName })
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
                                }, { injectBizId: true })
                            }
                        })
                        this.selectedNode.data.service_template_id = 0
                        this.instance.service_template_id = null
                        this.templateInfo.serviceTemplateName = this.$t('无')
                        this.disabledProperties = []
                    }
                })
            },
            async getDiffTemplateAndInstances () {
                try {
                    const data = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
                        bizId: this.business,
                        setTemplateId: this.instance.set_template_id,
                        params: {
                            bk_set_ids: [this.instance.bk_set_id]
                        },
                        config: {
                            requestId: 'diffTemplateAndInstances'
                        }
                    })
                    const diff = data[0] ? data[0].module_diffs : []
                    const len = diff.filter(_module => _module.diff_type !== 'unchanged').length
                    this.hasChange = !!len
                } catch (e) {
                    console.error(e)
                }
            },
            handleSyncSetTemplate () {
                this.$store.commit('setFeatures/setSyncIdMap', {
                    id: `${this.business}_${this.instance.set_template_id}`,
                    instancesId: [this.instance.bk_set_id]
                })
                this.$router.push({
                    name: 'setSync',
                    params: {
                        setTemplateId: this.instance.set_template_id,
                        moduleId: this.selectedNode.data.bk_inst_id
                    }
                })
            },
            goServiceTemplate () {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: this.instance.service_template_id,
                        moduleId: this.selectedNode.data.bk_inst_id
                    }
                })
            },
            goSetTemplate () {
                this.$router.push({
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'view',
                        templateId: this.instance.set_template_id,
                        moduleId: this.selectedNode.data.bk_inst_id
                    }
                })
            },
            async getSelectedNodeHostCount () {
                const defaultModel = ['biz', 'set', 'module', 'host', 'object']
                const modelInstKey = {
                    biz: 'bk_biz_id',
                    set: 'bk_set_id',
                    module: 'bk_module_id',
                    host: 'bk_host_id',
                    object: 'bk_inst_id'
                }
                const conditionParams = {
                    condition: defaultModel.map(model => {
                        return {
                            bk_obj_id: model,
                            condition: [],
                            fields: []
                        }
                    })
                }
                const selectedNode = this.selectedNode
                const selectedModel = defaultModel.includes(selectedNode.data.bk_obj_id) ? selectedNode.data.bk_obj_id : 'object'
                const selectedModelCondition = conditionParams.condition.find(model => model.bk_obj_id === selectedModel)
                selectedModelCondition.condition.push({
                    field: modelInstKey[selectedModel],
                    operator: '$eq',
                    value: selectedNode.data.bk_inst_id
                })
                const data = await this.$store.dispatch('hostSearch/searchHost', {
                    params: {
                        ...conditionParams,
                        bk_biz_id: this.business,
                        ip: {
                            flag: 'bk_host_innerip|bk_host_outer',
                            exact: 0,
                            data: []
                        },
                        page: {
                            start: 0,
                            limit: 1,
                            sort: ''
                        }
                    },
                    config: {
                        requestId: 'searchHosts',
                        cancelPrevious: true
                    }
                })
                return data && data.count
            }
        }
    }
</script>

<style lang="scss" scoped>
    .template-info {
        font-size: 14px;
        color: #63656e;
        padding: 20px 0 20px 36px;
        border-bottom: 1px solid #F0F1F5;
        .info-item {
            width: 50%;
            max-width: 400px;
            line-height: 26px;
        }
        .name {
            position: relative;
            padding-right: 16px;
            &::after {
                content: ":";
                position: absolute;
                right: 10px;
            }
        }
        .value {
            width: calc(80% - 10px);
            padding-right: 10px;
            .text {
                @include inlineBlock;
                @include ellipsis;
                max-width: calc(100% - 16px);
                font-size: 14px;
            }
            .template-value {
                width: 100%;
                font-size: 0;
                color: #3a84ff;
                cursor: pointer;
                &.set-template {
                    width: auto;
                    max-width: calc(100% - 75px);
                }
            }
            .icon-cc-share {
                @include inlineBlock;
                font-size: 12px;
                margin-left: 4px;
            }
        }
    }
    .topology-details {
        padding: 0 !important;
        /deep/ {
            .property-list {
                margin-left: 36px;
            }
            .details-options {
                padding: 28px 18px 0 36px;
            }
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
    .service-category {
        font-size: 12px;
        padding: 20px 0 24px 36px;
        border-bottom: 1px solid #dcdee5;
        .selector-item {
            position: relative;
            width: 50%;
            max-width: 554px;
            padding-right: 54px;
        }
        .category-selector {
            width: calc(50% - 5px);
            & + .category-selector {
                margin-left: 10px;
            }
        }
        .second-category-errors {
            position: absolute;
            top: 100%;
            left: 0;
            margin-left: calc(50% - 18px);
            line-height: 14px;
            font-size: 12px;
            color: #ff5656;
            max-width: 100%;
            @include ellipsis;
        }
    }
    .topology-form {
        /deep/ {
            .form-groups {
                padding: 0;
            }
            .property-list {
                margin-left: 36px;
            }
            .property-item {
                max-width: 554px !important;
            }
            .form-options {
                padding: 10px 0 0 36px;
            }
        }
    }
    .button-edit {
        min-width: 76px;
        margin-right: 4px;
    }
    .btn-delete{
        min-width: 76px;
        &:not(.is-disabled):hover {
            color: #ffffff;
            border-color: #ff5656;
            background-color: #ff5656;
        }
    }
    .sync-set-btn {
        position: relative;
        height: 26px;
        line-height: 24px;
        padding: 0 10px;
        font-size: 12px;
        margin-top: -2px;
        &.has-change::before {
            content: '';
            position: absolute;
            top: -4px;
            right: -4px;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background-color: #EA3636;
        }
    }
</style>
