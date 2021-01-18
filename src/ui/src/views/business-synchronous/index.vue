<template>
    <section class="batch-wrapper" v-bkloading="{ isLoading: $loading() }">
        <cmdb-tips>{{$t('同步模板功能提示')}}</cmdb-tips>
        <h2 class="title">{{$t('将会同步以下信息')}}：</h2>
        <div class="info-layout cleafix">
            <ul class="process-list fl">
                <li class="process-item"
                    v-for="(process, index) in processList"
                    :key="index"
                    :class="{
                        'show-tips': !process.confirmed,
                        'is-active': activeIndex === index,
                        'is-remove': process.type === 'removed'
                    }"
                    @click="handleChangeActive(process, index)">
                    <span class="process-name" :title="process.process_template_name">{{process.process_template_name}}</span>
                    <span class="process-service-count" v-if="process.type !== 'others'">{{getInstanceCount(process)}}</span>
                </li>
            </ul>
            <div class="change-details"
                v-if="current"
                :key="current.process_template_id">
                <cmdb-collapse class="details-info">
                    <div class="collapse-title" slot="title">
                        {{$t('变更内容')}}
                        <span v-if="current.type === 'changed'">（{{changedProperties.length}}）</span>
                    </div>
                    <div class="info-content">
                        <div class="process-info"
                            v-if="current.type === 'added'">
                            <div class="info-item" style="width: auto;">
                                {{$t('模板中新增进程')}}
                                <span class="info-item-value">{{current.process_template_name}}</span>
                            </div>
                        </div>
                        <div class="process-info"
                            v-if="current.type === 'removed'">
                            <div class="info-item" style="width: auto;">
                                <span class="info-item-value" style="font-weight: 700;">{{current.process_template_name}}</span>
                                {{$t('从模板中删除')}}
                            </div>
                        </div>
                        <div class="process-info clearfix"
                            v-else-if="current.type === 'changed'">
                            <div :class="['info-item fl', { table: changed.property.bk_property_type === 'table' }]"
                                v-for="(changed, index) in changedProperties"
                                :key="index"
                                v-bk-overflow-tips>
                                {{changed.property.bk_property_name}}：
                                <span class="info-item-value">
                                    <span v-if="changed.property.bk_property_id === 'bind_info' && !changed.template_property_value.length">
                                        {{$t('移除所有进程监听信息')}}
                                    </span>
                                    <cmdb-property-value v-else
                                        :value="getChangedValue(changed)"
                                        :property="changed.property">
                                    </cmdb-property-value>
                                </span>
                            </div>
                        </div>
                        <div class="process-info"
                            v-else-if="current.type === 'others'">
                            <div class="info-item" style="width: auto;">
                                {{$t('服务分类')}}：
                                <span class="info-item-value">{{current.service_category}}</span>
                            </div>
                        </div>
                    </div>
                </cmdb-collapse>
                <cmdb-collapse class="details-modules"
                    v-for="(module, index) in current.modules"
                    :key="index"
                    :collapse="true"
                    @collapse-change.once="handleModulesCollapseChange(...arguments, module)">
                    <div class="collapse-title" slot="title">
                        {{getModuleTopoPath(module.bk_module_id)}} {{$t('涉及实例')}}
                    </div>
                    <ul class="instance-list">
                        <li class="instance-item"
                            v-for="(instance, instanceIndex) in module.service_instances"
                            :key="instanceIndex"
                            @click="handleViewDiff(instance, module)">
                            <span class="instance-name" v-bk-overflow-tips>{{instance.service_instance.name}}</span>
                            <span class="instance-diff-count"
                                v-if="instance.changed_attributes">
                                ({{instance.changed_attributes.length}})
                            </span>
                        </li>
                    </ul>
                </cmdb-collapse>
            </div>
        </div>
        <div class="batch-options">
            <bk-button class="mr10" theme="primary"
                :disabled="!allConfirmed"
                @click="handleConfirm">
                {{$t('确认并同步')}}
            </bk-button>
            <bk-button @click="handleGoBackModule">{{$t('取消')}}</bk-button>
        </div>
        <bk-sideslider
            v-transfer-dom
            :width="676"
            :is-show.sync="slider.show"
            :title="slider.title">
            <template slot="content" v-if="slider.show">
                <instance-details slot="content"
                    v-if="slider.show"
                    v-bind="slider.props"
                    :properties="properties">
                </instance-details>
            </template>
        </bk-sideslider>
    </section>
</template>

<script>
    import InstanceDetails from './children/details.vue'
    import formatter from '@/filters/formatter'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            InstanceDetails
        },
        data () {
            return {
                processList: [],
                properties: [],
                activeIndex: null,
                topoPath: {},
                categories: [],
                slider: {
                    show: false,
                    title: '',
                    props: {
                        module: null,
                        instance: null,
                        type: ''
                    }
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectBiz', ['bizId']),
            current () {
                if (this.activeIndex !== null) {
                    return this.processList[this.activeIndex]
                }
                return null
            },
            changedProperties () {
                if (this.current && this.current.type === 'changed') {
                    return this.getChangedProperties()
                }
                return []
            },
            templateId () {
                return Number(this.$route.params.template)
            },
            modules () {
                return String(this.$route.params.modules).split(',').map(id => Number(id))
            },
            allConfirmed () {
                return this.processList.every(process => process.confirmed)
            }
        },
        async created () {
            try {
                await this.getProperties()
                this.getTopoPath()
                this.getDifference()
            } catch (e) {
                console.error(e)
            }
        },
        methods: {
            handleChangeActive (process, index) {
                this.activeIndex = index
                process.confirmed = true
            },
            async getProperties () {
                try {
                    this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: {
                            bk_obj_id: 'process',
                            bk_supplier_account: this.supplierAccount,
                            bk_biz_id: this.bizId
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            async getTopoPath () {
                try {
                    const { nodes } = await this.$store.dispatch('objectMainLineModule/getTopoPath', {
                        bizId: this.bizId,
                        params: {
                            topo_nodes: this.modules.map(moduleId => ({ bk_obj_id: 'module', bk_inst_id: moduleId }))
                        }
                    })
                    const topoPath = {}
                    nodes.forEach(node => {
                        topoPath[node.topo_node.bk_inst_id] = node.topo_path.reverse().map(path => path.bk_inst_name).join(' / ')
                    })
                    this.topoPath = topoPath
                } catch (e) {
                    console.error(e)
                }
            },
            async getDifference () {
                try {
                    const differences = await this.$store.dispatch('businessSynchronous/searchServiceInstanceDifferences', {
                        params: {
                            bk_module_ids: this.modules,
                            service_template_id: this.templateId,
                            bk_biz_id: this.bizId
                        }
                    })
                    const processList = []
                    const differenceType = ['changed', 'added', 'removed']
                    let changedCategory = null
                    differences.forEach(difference => {
                        differenceType.forEach(type => {
                            difference[type].forEach(info => {
                                const moduleInfo = { ...info, bk_module_id: difference.bk_module_id }
                                const item = processList.find(item => item.type === type && item.process_template_id === info.process_template_id)
                                if (item) {
                                    item.modules.push(moduleInfo)
                                } else {
                                    const newItem = {
                                        type: type,
                                        process_template_id: info.process_template_id,
                                        process_template_name: info.process_template_name,
                                        modules: [moduleInfo]
                                    }
                                    const length = processList.push(newItem)
                                    newItem.confirmed = length === 1
                                }
                            })
                        })
                        changedCategory = (difference.changed_attributes || []).find(attr => attr.property_id === 'service_category_id')
                    })
                    if (changedCategory) {
                        const categoryInfo = await this.getServiceCategoryDifference(changedCategory)
                        categoryInfo.confirmed = !processList.length
                        processList.push(categoryInfo)
                    }
                    this.processList = processList
                    this.activeIndex = 0
                } catch (e) {
                    console.error(e)
                }
            },
            async getServiceCategoryDifference (changedCategory) {
                try {
                    const categoryInfo = {
                        type: 'others',
                        process_template_id: 'service_category_id',
                        process_template_name: this.$t('服务分类变更'),
                        modules: []
                    }
                    const [{ info: categories }, { info: modules }] = await Promise.all([
                        this.getServiceCategory(),
                        this.getServiceModules()
                    ])
                    this.categories = categories
                    const templateCategoryId = changedCategory.template_property_value
                    categoryInfo.service_category = this.getCagetoryPath(templateCategoryId)
                    modules.forEach(module => {
                        if (module.service_category_id !== templateCategoryId) {
                            categoryInfo.modules.push({
                                bk_module_id: module.bk_module_id,
                                template_service_category: categoryInfo.service_category,
                                current_service_category: this.getCagetoryPath(module.service_category_id),
                                service_instances: []
                            })
                        }
                    })

                    return categoryInfo
                } catch (e) {
                    console.error(e)
                }
            },
            getServiceCategory () {
                return this.$store.dispatch('serviceClassification/searchServiceCategory', {
                    params: { bk_biz_id: this.bizId }
                })
            },
            getCagetoryPath (id) {
                const second = this.categories.find(second => second.category.id === id)
                const first = this.categories.find(first => first.category.id === this.$tools.getValue(second, 'category.bk_parent_id'))
                const firstName = this.$tools.getValue(first, 'category.name') || '--'
                const secondName = this.$tools.getValue(second, 'category.name') || '--'
                return [firstName, secondName].join(' / ')
            },
            getServiceModules () {
                return this.$store.dispatch('serviceTemplate/getServiceTemplateModules', {
                    bizId: this.bizId,
                    serviceTemplateId: this.templateId,
                    params: {
                        bk_module_ids: this.modules
                    }
                })
            },
            getChangedProperties () {
                const changed = []
                this.current.modules.forEach(module => {
                    module.service_instances.forEach(instance => {
                        (instance.changed_attributes || []).forEach(changedProperty => {
                            const isExist = changed.some(exist => exist.property.bk_property_id === changedProperty.property_id)
                            const property = this.properties.find(property => property.bk_property_id === changedProperty.property_id)
                            if (!isExist && property) {
                                changed.push({
                                    property: property,
                                    template_property_value: changedProperty.template_property_value
                                })
                            }
                        })
                    })
                })
                return changed
            },
            getChangedValue (changed) {
                const property = changed.property
                let value = changed.template_property_value
                value = Object.prototype.toString.call(value) === '[object Object]' ? value.value : value
                return formatter(value, property)
            },
            getModuleTopoPath (moduleId) {
                return this.topoPath[moduleId]
            },
            async handleModulesCollapseChange (collapse, module) {
                const loaded = module.service_instances.__loaded__
                if (this.current.type === 'others' && !loaded) {
                    try {
                        module.service_instances.__loaded__ = true
                        const { info: instances } = await this.getModuleServiceInstances(module.bk_module_id)
                        module.service_instances.push(...instances.map(instance => {
                            return {
                                changed_attributes: [{
                                    property_id: 'service_category_id',
                                    property_name: this.$t('服务分类'),
                                    property_value: module.current_service_category,
                                    template_property_value: module.template_service_category
                                }],
                                service_instance: instance
                            }
                        }))
                    } catch (e) {
                        console.error(e)
                        module.service_instances.__loaded__ = false
                    }
                }
            },
            getModuleServiceInstances (moduleId) {
                return this.$store.dispatch('serviceInstance/getModuleServiceInstances', {
                    params: {
                        bk_biz_id: this.bizId,
                        bk_module_id: moduleId,
                        with_name: true
                    }
                })
            },
            handleViewDiff (instance, module) {
                this.slider.title = instance.service_instance.name
                this.slider.props = {
                    module,
                    instance,
                    type: this.current.type
                }
                this.slider.show = true
            },
            handleConfirm () {
                this.$store.dispatch('businessSynchronous/syncServiceInstanceByTemplate', {
                    params: {
                        service_template_id: this.templateId,
                        bk_module_ids: this.modules,
                        service_instances: this.getServiceInstanceIds(),
                        bk_biz_id: this.bizId
                    }
                }).then(() => {
                    this.$success(this.$t('同步成功'))
                    this.handleGoBackModule()
                })
            },
            handleGoBackModule () {
                this.$routerActions.back()
            },
            getServiceInstanceIds () {
                const ids = []
                this.processList.forEach(process => {
                    process.modules.forEach(module => {
                        module.service_instances.forEach(data => {
                            ids.push(data.service_instance.id)
                        })
                    })
                })
                return [...new Set(ids)]
            },
            getInstanceCount (process) {
                return process.modules.reduce((count, module) => count + module.service_instances.length, 0)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .batch-wrapper {
        padding: 10px 20px;
        .title {
            margin-top: 24px;
            font-size: 14px;
            line-height: 20px;
        }
        .collapse-title {
            font-size: 14px;
            color: $textColor;
        }
    }
    .info-layout {
        margin-top: 10px;
        border: 1px solid $borderColor;
        border-bottom: none;
        height: calc(100vh - 350px);
        overflow: hidden;
        .process-list {
            position: relative;
            margin-right: -1px;
            width: 200px;
            height: 100%;
            z-index: 2;
            @include scrollbar-y;
        }
        .change-details {
            position: relative;
            height: 100%;
            padding: 20px;
            background-color: #FFF;
            border-left: 1px solid $borderColor;
            border-bottom: 1px solid $borderColor;
            z-index: 1;
            @include scrollbar-y;
        }
    }
    .process-list {
        border-bottom: 1px solid $borderColor;
        .process-item {
            display: flex;
            padding: 0 12px 0 14px;
            height: 61px;
            align-items: center;
            justify-content: space-between;
            background-color: #FAFBFD;
            border-right: 1px solid $borderColor;
            border-bottom: 1px solid $borderColor;
            cursor: pointer;
            &.is-active {
                background-color: #FFF;
                border-right: none;
                .process-name {
                    font-weight: bold;
                    color: $primaryColor;
                }
                &.is-remove {
                    .process-name {
                        color: $dangerColor;
                    }
                }
            }
            &.is-remove {
                .process-name {
                    text-decoration: line-through;
                }
            }
            &.show-tips {
                .process-name:after {
                    position: absolute;
                    width: 6px;
                    height: 6px;
                    top: 21px;
                    right: 4px;
                    border-radius: 50%;
                    background-color: #FF5656;
                    content: "";
                    z-index: 1;
                }
            }
            .process-name {
                line-height: 60px;
                position: relative;
                padding: 0 14px 0 0;
                @include ellipsis;
            }
            .process-service-count {
                padding: 0 8px;
                height: 16px;
                line-height: 16px;
                font-size: 12px;
                font-style: normal;
                text-align: center;
                background-color: #c4c6cc;
                color: #fff;
                border-radius: 8px;
            }
        }
    }
    .details-info {
        .process-info {
            padding: 0 0 0 22px;
            .info-item {
                width: 200px;
                font-size: 14px;
                margin: 20px 40px 0 0;
                @include ellipsis;
                .info-item-value {
                    color: #313238;
                }

                &.table {
                    width: 100%;
                    /deep/ .table-value {
                        width: 800px;
                    }
                }
            }
        }
    }
    .details-modules {
        margin-top: 60px;
        & ~ .details-modules {
            margin-top: 20px;
        }
    }
    .instance-list {
        padding: 0 0 0 22px;
        .instance-item {
            display: inline-flex;
            align-items: center;
            justify-content: space-between;
            width: 240px;
            margin: 10px 80px 0 0;
            padding: 0 4px;
            height: 22px;
            border: 1px solid $borderColor;
            background-color: #FAFBFD;
            font-size: 12px;
            cursor: pointer;
            &:hover {
                .instance-name,
                .instance-diff-count {
                    color: $primaryColor;
                }
                .instance-diff-count {
                    font-weight: bold;
                }
            }
            .instance-name {
                padding-right: 20px;
                @include ellipsis;
            }
            .instance-diff-count {
                color: #C4C6CC;
            }
        }
    }
    .batch-options {
        margin-top: 20px;
    }
</style>
