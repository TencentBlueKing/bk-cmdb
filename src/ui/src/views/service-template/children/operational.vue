<template>
    <cmdb-sticky-layout class="create-template-wrapper" v-bkloading="{ isLoading: $loading(Object.values(request)) }">
        <div class="info-group">
            <h3>{{$t('基本属性')}}</h3>
            <div class="template-info clearfix"
                v-if="isFormMode"
                :class="{
                    'is-edit': insideMode === 'edit'
                }">
                <div class="form-info clearfix">
                    <label class="label-text fl" for="templateName">
                        {{$t('模板名称')}}
                        <span class="color-danger" v-if="isCreateMode">*</span>
                        <span v-else>：</span>
                    </label>
                    <template v-if="isEditNameLoading">
                        <i class="form-loading fl"></i>
                    </template>
                    <template v-else>
                        <div class="cmdb-form-item clearfix fl" :class="{ 'is-error': errors.has('templateName') }">
                            <template v-if="isCreateMode || isEditName">
                                <bk-input type="text" class="cmdb-form-input fl" id="templateName"
                                    name="templateName"
                                    :placeholder="$t('请输入模板名称')"
                                    :class="{ 'is-edit-name': isEditName }"
                                    v-model.trim="formData.templateName"
                                    v-validate="'required|singlechar|length:256'">
                                </bk-input>
                                <p class="form-error">{{errors.first('templateName')}}</p>
                            </template>
                            <template v-if="isEditName">
                                <i class="form-confirm edit-icon bk-icon icon-check-1 fl" @click="handleSaveName"></i>
                                <i class="form-cancel edit-icon bk-icon icon-close fl" @click="handleCancelEditName"></i>
                            </template>
                            <template v-else-if="!isCreateMode">
                                <span class="template-name" :title="formData.templateName">{{formData.templateName}}</span>
                                <i class="icon-cc-edit" v-if="!isCreateMode" @click="handleEditName"></i>
                            </template>
                        </div>
                    </template>
                </div>
                <div class="form-info clearfix">
                    <span class="label-text fl">
                        {{$t('服务分类')}}
                        <span class="color-danger" v-if="isCreateMode">*</span>
                        <span v-else>：</span>
                    </span>
                    <template v-if="isCreateMode || isEditCategory">
                        <template v-if="isEditCategoryLoading">
                            <i class="form-loading fl"></i>
                        </template>
                        <template v-else>
                            <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('mainClassificationId') }" style="width: auto;">
                                <cmdb-selector
                                    class="fl"
                                    display-key="displayName"
                                    :placeholder="$t('请选择一级分类')"
                                    :searchable="true"
                                    :list="mainList"
                                    :popover-options="{
                                        boundary: 'window'
                                    }"
                                    v-validate="'required'"
                                    name="mainClassificationId"
                                    v-model="formData['mainClassification']"
                                    @on-selected="handleSelect">
                                    <template slot-scope="{ name, id }">
                                        <div class="bk-option-content-default" :title="`${name}（#${id}）`">
                                            <div class="bk-option-name medium-font">
                                                {{name}}<span class="category-id">（#{{id}}）</span>
                                            </div>
                                        </div>
                                    </template>
                                </cmdb-selector>
                                <p class="form-error">{{errors.first('mainClassificationId')}}</p>
                            </div>
                            <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('secondaryClassificationId') }" style="width: auto;">
                                <cmdb-selector
                                    class="fl"
                                    display-key="displayName"
                                    :placeholder="$t('请选择二级分类')"
                                    :auto-select="true"
                                    :searchable="true"
                                    :list="secondaryList"
                                    v-validate="'required'"
                                    name="secondaryClassificationId"
                                    v-model="formData['secondaryClassification']">
                                    <template slot-scope="{ name, id }">
                                        <div class="bk-option-content-default" :title="`${name}（#${id}）`">
                                            <div class="bk-option-name medium-font">
                                                {{name}}<span class="category-id">（#{{id}}）</span>
                                            </div>
                                        </div>
                                    </template>
                                </cmdb-selector>
                                <p class="form-error">{{errors.first('secondaryClassificationId')}}</p>
                            </div>
                            <template v-if="isEditCategory">
                                <i class="form-confirm edit-icon bk-icon icon-check-1" @click="handleSaveCategory"></i>
                                <i class="form-cancel edit-icon bk-icon icon-close" @click="handleCancelEditCategory"></i>
                            </template>
                        </template>
                    </template>
                    <template v-else>
                        <span class="info-content" :title="getServiceCategory()">
                            {{getServiceCategory()}}
                        </span>
                        <i class="icon-cc-edit" @click="handleEditCategory"></i>
                    </template>
                </div>
            </div>

            <div class="view-group clearfix" v-else>
                <div class="view-info fl clearfix">
                    <label class="info-label fl">
                        {{$t('模板名称')}}
                    </label>
                    <span class="info-content" :title="formData.templateName">{{formData.templateName}}</span>
                </div>
                <div class="view-info fl clearfix">
                    <label class="info-label fl">
                        {{$t('服务分类')}}
                    </label>
                    <span class="info-content" :title="getServiceCategory()">
                        {{getServiceCategory()}}
                    </span>
                </div>
            </div>
        </div>
        <div class="info-group">
            <h3>{{$t('服务进程')}}</h3>
            <div class="precess-box">
                <div class="process-create" v-if="isFormMode">
                    <cmdb-auth :auth="auth">
                        <bk-button slot-scope="{ disabled }"
                            class="create-btn"
                            theme="default"
                            :disabled="disabled"
                            @click="handleCreateProcess">
                            <i class="bk-icon icon-plus"></i>
                            <span>{{$t('新建进程')}}</span>
                        </bk-button>
                    </cmdb-auth>
                    <span class="create-tips">{{$t('新建进程提示')}}</span>
                </div>
                <process-table
                    v-if="processList.length"
                    :loading="processLoading"
                    :properties="properties"
                    :auth="auth"
                    :show-operation="isFormMode"
                    @on-edit="handleUpdateProcess"
                    @on-delete="handleDeleteProcess"
                    :list="processList">
                </process-table>
                <div v-else-if="!isFormMode" class="process-empty">{{$t('暂未配置进程')}}</div>
            </div>
        </div>
        <div v-if="insideMode !== 'edit'"
            slot="footer"
            slot-scope="{ sticky }"
            :class="{
                'info-footer': true,
                'is-sticky': sticky
            }">
            <cmdb-auth class="mr5" :auth="auth">
                <bk-button slot-scope="{ disabled }"
                    theme="primary"
                    :disabled="disabled"
                    @click="handleSubmit">
                    {{getButtonText()}}
                </bk-button>
            </cmdb-auth>
            <bk-button @click="handleReturn" v-show="isFormMode">{{$t('取消')}}</bk-button>
        </div>
        <bk-sideslider
            v-transfer-dom
            :is-show.sync="slider.show"
            :title="slider.title"
            :width="800"
            :before-close="handleSliderBeforeClose">
            <template slot="content" v-if="slider.show">
                <process-form
                    ref="processForm"
                    :auth="auth"
                    :properties="properties"
                    :property-groups="propertyGroups"
                    :inst="attribute.inst.edit"
                    :type="attribute.type"
                    :is-created-service="isCreateMode"
                    :data-index="attribute.dataIndex"
                    :save-disabled="false"
                    :has-used="hasUsed"
                    :submit-format="formatSubmitData"
                    @on-submit="handleSaveProcess"
                    @on-cancel="handleCancelProcess">
                </process-form>
            </template>
        </bk-sideslider>
        <bk-dialog v-model="showUpdateInfo"
            :esc-close="false"
            :mask-close="false"
            :show-footer="false"
            :close-icon="false">
            <div class="update-alert-layout">
                <i class="bk-icon icon-check-1"></i>
                <h3>{{$t('修改成功')}}</h3>
                <div class="btns">
                    <bk-button class="mr10" theme="primary" @click="handleToSyncInstance">{{$t('同步实例')}}</bk-button>
                    <bk-button theme="default" @click="handleCancelOperation">{{$t('返回列表')}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </cmdb-sticky-layout>
</template>

<script>
    import processForm from './process-form.vue'
    import processTable from './process'
    import { mapActions, mapGetters, mapMutations } from 'vuex'
    import {
        MENU_BUSINESS_SERVICE_TEMPLATE,
        MENU_BUSINESS_HOST_AND_SERVICE
    } from '@/dictionary/menu-symbol'
    import Bus from '@/utils/bus'
    export default {
        components: {
            processTable,
            processForm
        },
        data () {
            return {
                processLoading: false,
                properties: [],
                propertyGroups: [],
                mainList: [],
                secondaryList: [],
                allSecondaryList: [],
                processList: [],
                originTemplateValues: {},
                hasUsed: false,
                attribute: {
                    type: null,
                    dataIndex: null,
                    inst: {
                        details: {},
                        edit: {}
                    }
                },
                slider: {
                    show: false,
                    title: ''
                },
                formData: {
                    mainClassification: '',
                    originMainClassification: '',
                    secondaryClassification: '',
                    originSecondaryClassification: '',
                    templateName: '',
                    templateId: ''
                },
                insideMode: this.$route.params.templateId ? 'edit' : 'view',
                isEditCategory: false,
                isEditCategoryLoading: false,
                isEditName: false,
                isEditNameLoading: false,
                deletable: false,
                showUpdateInfo: false,
                request: {
                    template: Symbol('template'),
                    category: Symbol('category'),
                    properties: Symbol('properties'),
                    propertyGroups: Symbol('propertiyGroups'),
                    processList: Symbol('processList')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('serviceProcess', ['localProcessTemplate']),
            ...mapGetters('objectBiz', ['bizId']),
            templateId () {
                return this.$route.params.templateId
            },
            isCreateMode () {
                return this.templateId === undefined
            },
            isFormMode () {
                return this.isCreateMode || this.insideMode === 'edit'
            },
            auth () {
                if (this.isCreateMode) {
                    return {
                        type: this.$OPERATION.C_SERVICE_TEMPLATE,
                        relation: [this.bizId]
                    }
                }
                return {
                    type: this.$OPERATION.U_SERVICE_TEMPLATE,
                    relation: [this.bizId, Number(this.templateId)]
                }
            },
            setActive () {
                return this.$route.params.active
            }
        },
        created () {
            Bus.$on('module-loaded', count => {
                this.deletable = !count
            })
            this.setBreadcrumbs()
            this.processList = this.localProcessTemplate
            this.refresh()
        },
        beforeDestroy () {
            Bus.$off('module-loaded')
            this.clearLocalProcessTemplate()
        },
        methods: {
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('serviceClassification', ['searchServiceCategory']),
            ...mapActions('serviceTemplate', [
                'createServiceTemplate',
                'updateServiceTemplate',
                'findServiceTemplate',
                'deleteServiceTemplate'
            ]),
            ...mapActions('processTemplate', [
                'createProcessTemplate',
                'updateProcessTemplate',
                'deleteProcessTemplate',
                'getProcessTemplate',
                'getBatchProcessTemplate'
            ]),
            ...mapMutations('serviceProcess', [
                'deleteLocalProcessTemplate',
                'clearLocalProcessTemplate'
            ]),
            async refresh () {
                try {
                    await this.reload()
                    if (this.setActive) {
                        Bus.$emit('active-change', 'instance')
                        this.$route.params.active = null
                    }
                    if (this.$route.params.isEdit) {
                        this.insideMode = 'edit'
                        this.$route.params.isEdit = null
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            setBreadcrumbs () {
                if (this.isCreateMode) {
                    this.$store.commit('setTitle', this.$t('新建模板'))
                } else {
                    this.$store.commit('setTitle', this.$t('模板详情'))
                }
            },
            async reload () {
                try {
                    const request = [
                        this.getProperties(),
                        this.getPropertyGroups(),
                        this.getServiceClassification()
                    ]
                    if (!this.isCreateMode) {
                        request.push(this.getSingleServiceTemplate())
                        request.push(this.getProcessList())
                    }
                    const [properties, groups, { info: categories }, templateResponse] = await Promise.all(request)
                    this.properties = properties
                    this.propertyGroups = groups
                    const categoryList = categories.map(item => {
                        return {
                            ...item.category,
                            displayName: `${item.category.name}（#${item.category.id}）`
                        }
                    })
                    this.allSecondaryList = categoryList.filter(classification => classification.bk_parent_id)
                    this.mainList = categoryList.filter(classification => !classification.bk_parent_id)
                    this.allSecondaryList = categoryList.filter(classification => classification.bk_parent_id)
                    if (!this.isCreateMode) {
                        const { result, data } = templateResponse
                        if (!result) {
                            this.$router.replace({ name: '404' })
                            return
                        }
                        // 原始数据
                        this.originTemplateValues = {
                            service_instance_count: data.service_instance_count,
                            process_instance_count: data.process_instance_count,
                            ...data.template
                        }
                        // 表单数据
                        const secondCategoryId = data.template.service_category_id
                        const secondCategory = this.allSecondaryList.find(category => category.id === secondCategoryId) || {}
                        this.formData.templateId = data.template.id
                        this.formData.templateName = data.template.name
                        this.formData.mainClassification = secondCategory.bk_parent_id
                        this.formData.secondaryClassification = secondCategoryId
                        // 备份，用于取消编辑
                        this.formData.originMainClassification = secondCategory.bk_parent_id
                        this.formData.originSecondaryClassification = secondCategoryId
                        this.hasUsed = data.service_instance_count > 0
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            getProperties () {
                return this.searchObjectAttribute({
                    params: {
                        bk_biz_id: this.bizId,
                        bk_obj_id: 'process',
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: this.request.properties,
                        fromCache: true
                    }
                })
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: 'process',
                    params: { bk_biz_id: this.bizId },
                    config: {
                        requestId: this.request.propertyGroups,
                        fromCache: true
                    }
                })
            },
            getSingleServiceTemplate () {
                return this.findServiceTemplate({
                    id: this.templateId,
                    config: {
                        requestId: this.request.template,
                        globalError: false,
                        transformData: false
                    }
                })
            },
            getServiceClassification () {
                return this.searchServiceCategory({
                    params: { bk_biz_id: this.bizId },
                    config: {
                        requestId: this.request.category
                    }
                })
            },
            getProcessList () {
                this.processLoading = true
                this.getBatchProcessTemplate({
                    params: {
                        bk_biz_id: this.bizId,
                        service_template_id: Number(this.templateId)
                    },
                    config: {
                        requestId: this.request.processList
                    }
                }).then(data => {
                    this.processList = data.info.map(template => {
                        return {
                            process_id: template.id,
                            ...template['property']
                        }
                    })
                }).finally(() => {
                    this.processLoading = false
                })
            },
            handleSelect (id, data) {
                this.secondaryList = this.allSecondaryList.filter(classification => classification['bk_parent_id'] === id)
                if (!this.secondaryList.length) {
                    this.formData.secondaryClassification = ''
                }
            },
            formatSubmitData (data = {}) {
                const keys = Object.keys(data)
                for (const key of keys) {
                    const property = this.properties.find(property => property.bk_property_id === key)
                    if (property
                        && ['enum', 'int', 'float', 'list'].includes(property.bk_property_type)
                        && (!data[key].value || !data[key].as_default_value)) {
                        data[key].value = null
                    } else if (property && property.bk_property_type === 'bool' && !data[key].as_default_value) {
                        data[key].value = null
                    } else if (property && property.bk_property_type === 'table') {
                        // 过滤掉无效数据行
                        const values = data[key].value || []
                        const formattedValues = values.map(row => {
                            const rawValues = {}
                            Object.keys(row).forEach(rowKey => {
                                rawValues[rowKey] = row[rowKey].value
                            })
                            return this.$tools.formatValues(rawValues, property.option)
                        })
                        const newValues = values.map((row, rowIndex) => {
                            const newRowValues = {}
                            Object.keys(row).forEach(rowKey => {
                                if (typeof row[rowKey] === 'object') {
                                    const originalValue = values[rowIndex][rowKey].value
                                    const formattedValue = formattedValues[rowIndex][rowKey]
                                    newRowValues[rowKey] = {
                                        value: originalValue === formattedValue ? originalValue : formattedValue,
                                        as_default_value: values[rowIndex][rowKey].as_default_value
                                    }
                                } else {
                                    newRowValues[rowKey] = row[rowKey]
                                }
                            })
                            return newRowValues
                        })
                        data[key].value = newValues
                    } else if (!data[key].as_default_value) {
                        data[key].value = ''
                    }
                }
                return data
            },
            handleSaveProcess (values, changedValues, type) {
                const data = type === 'create' ? values : changedValues
                const processValues = this.formatSubmitData(data)
                if (type === 'create') {
                    this.createProcessTemplate({
                        params: {
                            bk_biz_id: this.bizId,
                            service_template_id: this.originTemplateValues['id'],
                            processes: [{
                                spec: processValues
                            }]
                        }
                    }).then(() => {
                        this.getProcessList()
                        this.handleCancelProcess()
                    })
                } else {
                    this.updateProcessTemplate({
                        params: {
                            bk_biz_id: this.bizId,
                            process_template_id: values['process_id'],
                            process_property: processValues
                        }
                    }).then(() => {
                        this.getProcessList()
                        this.handleCancelProcess()
                    })
                }
            },
            handleCancelProcess () {
                this.slider.show = false
            },
            handleCreateProcess () {
                this.slider.show = true
                this.slider.title = this.$t('添加进程')
                this.attribute.type = 'create'
                this.attribute.inst.edit = {}
            },
            handleUpdateProcess (template, index) {
                try {
                    this.slider.show = true
                    this.slider.title = template['bk_func_name']['value']
                    this.attribute.type = 'update'
                    this.attribute.inst.edit = template
                    this.attribute.dataIndex = index
                } catch (e) {
                    console.error(e)
                }
            },
            handleDeleteProcess (template, index) {
                this.$bkInfo({
                    title: this.$t('确认删除模板进程'),
                    confirmFn: () => {
                        if (this.isCreateMode) {
                            this.deleteLocalProcessTemplate({ process: template, index })
                            this.processList = this.localProcessTemplate
                        } else {
                            this.deleteProcessTemplate({
                                params: {
                                    data: {
                                        bk_biz_id: this.bizId,
                                        process_templates: [template['process_id']]
                                    }
                                }
                            }).then(() => {
                                this.$success(this.$t('删除成功'))
                                this.getProcessList()
                            })
                        }
                    }
                })
            },
            handleSubmitProcessList () {
                this.createProcessTemplate({
                    params: {
                        bk_biz_id: this.bizId,
                        service_template_id: this.formData.templateId,
                        processes: this.processList.map(process => {
                            delete process.sign_id
                            return {
                                spec: this.formatSubmitData(process)
                            }
                        })
                    }
                }).then(() => {
                    this.$success(this.$t('创建成功'))
                    this.handleCancelOperation()
                }).catch(async e => {
                    // 新建进程失败静默删除服务模板
                    await this.deleteServiceTemplate({
                        params: {
                            data: {
                                bk_biz_id: this.bizId,
                                service_template_id: this.formData.templateId
                            }
                        }
                    })
                    this.formData.templateId = ''
                })
            },
            async handleSubmit () {
                if (!this.isCreateMode && this.insideMode === 'view') {
                    this.insideMode = 'edit'
                    return false
                }
                if (!await this.$validator.validateAll()) return
                if (this.formData.templateId) {
                    this.updateServiceTemplate({
                        params: {
                            bk_biz_id: this.bizId,
                            id: this.formData.templateId,
                            name: this.formData.templateName,
                            service_category_id: this.formData.secondaryClassification
                        }
                    }).then(() => {
                        if (this.isCreateMode) {
                            this.handleSubmitProcessList()
                        } else {
                            this.showUpdateInfo = true
                        }
                    })
                } else {
                    if (!this.processList.length) {
                        this.$bkInfo({
                            title: this.$t('确认提交'),
                            subTitle: this.$t('服务模板创建没进程提示'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.handleCreateTemplate()
                            }
                        })
                        return
                    }
                    this.handleCreateTemplate()
                }
            },
            handleCreateTemplate () {
                this.createServiceTemplate({
                    params: {
                        name: this.formData.templateName,
                        bk_biz_id: this.bizId,
                        service_category_id: this.formData.secondaryClassification
                    }
                }).then(data => {
                    if (this.processList.length) {
                        this.formData.templateId = data.id
                        this.handleSubmitProcessList()
                    } else {
                        this.$success(this.$t('创建成功'))
                        this.handleCancelOperation()
                    }
                })
            },
            handleReturn () {
                if (this.insideMode === 'edit') {
                    this.insideMode = 'view'
                    this.refresh()
                    return
                }
                const moduleId = this.$route.params['moduleId']
                if (moduleId) {
                    this.$routerActions.redirect({
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: {
                            node: 'module-' + this.$route.params.moduleId
                        }
                    })
                } else {
                    this.handleCancelOperation()
                }
            },
            handleCancelOperation () {
                this.showUpdateInfo = false
                this.$routerActions.redirect({ name: MENU_BUSINESS_SERVICE_TEMPLATE })
            },
            handleSliderBeforeClose () {
                const hasChanged = this.$refs.processForm && this.$refs.processForm.hasChange()
                if (hasChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                return true
            },
            getServiceCategory () {
                const first = this.mainList.find(first => first.id === this.formData.mainClassification) || {}
                const second = this.allSecondaryList.find(second => second.id === this.formData.secondaryClassification) || {}
                return `${first.name || '--'} / ${second.name || '--'}`
            },
            getButtonText () {
                if (this.isCreateMode) {
                    return this.$t('提交')
                } else if (this.insideMode === 'view') {
                    return this.$t('编辑')
                }
                return this.$t('保存')
            },
            handleToSyncInstance () {
                Bus.$emit('active-change', 'instance')
                this.showUpdateInfo = false
            },
            handleEditName () {
                this.isEditName = true
            },
            handleCancelEditName () {
                this.formData.templateName = this.originTemplateValues.name
                this.isEditName = false
            },
            async handleSaveName () {
                try {
                    const isValid = await this.$validator.validate('templateName')
                    if (!isValid) {
                        return false
                    }
                    this.isEditNameLoading = true
                    await this.updateServiceTemplate({
                        params: {
                            bk_biz_id: this.bizId,
                            id: this.formData.templateId,
                            name: this.formData.templateName
                        }
                    })
                    this.isEditName = false
                    this.isEditNameLoading = false
                } catch (e) {
                    console.error(e)
                    this.isEditNameLoading = false
                }
            },
            async handleSaveCategory () {
                try {
                    const isValid = await this.$validator.validateAll()
                    if (!isValid) {
                        return false
                    }
                    this.isEditCategoryLoading = true
                    await this.updateServiceTemplate({
                        params: {
                            id: this.formData.templateId,
                            bk_biz_id: this.bizId,
                            service_category_id: this.formData.secondaryClassification
                        }
                    })
                    this.isEditCategory = false
                    this.isEditCategoryLoading = false
                } catch (e) {
                    console.error(e)
                    this.isEditCategoryLoading = false
                }
            },
            handleCancelEditCategory () {
                this.formData.mainClassification = this.formData.originMainClassification
                this.formData.secondaryClassification = this.formData.originSecondaryClassification
                this.isEditCategory = false
            },
            handleEditCategory () {
                this.isEditCategory = true
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-template-wrapper {
        max-height: 100%;
        padding-top: 15px;
        @include scrollbar-y;
        .template-info {
            &.is-edit {
                .form-info {
                    float: left;
                    .cmdb-form-item{
                        width: 200px;
                    }
                    .info-content {
                        display: inline-block;
                        vertical-align: middle;
                        line-height: 36px;
                    }
                    .icon-cc-edit {
                        display: inline-block;
                        vertical-align: middle;
                        margin-left: 5px;
                        color: $primaryColor;
                        cursor: pointer;
                        &:hover {
                            color: #1964e1;
                        }
                    }
                }
                .edit-icon {
                    display: inline-block;
                    vertical-align: middle;
                    width: 32px;
                    height: 32px;
                    margin: 0 0 0 6px;
                    border-radius: 2px;
                    border: 1px solid #c4c6cc;
                    line-height: 30px;
                    font-size: 12px;
                    text-align: center;
                    cursor: pointer;
                    &.form-confirm {
                        color: #0082ff;
                        font-size: 20px;
                        &:before {
                            display: inline-block;
                        }
                    }
                    &.form-cancel {
                        color: #979ba5;
                        font-size: 20px;
                        &:before {
                            display: inline-block;
                        }
                    }
                    &:hover {
                        font-weight: bold;
                    }
                }
            }
        }
        .form-loading {
            width: 16px;
            height: 36px;
            margin: 2px 0;
            background-image: url("../../../assets/images/icon/loading.svg");
            background-position: center center;
            background-repeat: no-repeat;
        }
        .info-group {
            padding: 0 60px 20px 40px;
            h3 {
                color: #63656e;
                font-size: 14px;
                padding-bottom: 26px;
            }
            .form-info {
                font-size: 14px;
                padding-left: 30px;
                padding-bottom: 22px;
                .label-text {
                    line-height: 36px;
                    padding-right: 6px;
                }
                .cmdb-form-item {
                    width: 520px;
                }
                .template-name {
                    @include inlineBlock;
                    @include ellipsis;
                    max-width: calc(100% - 20px);
                    line-height: 36px;
                }
                .bk-select {
                    width: 254px;
                    margin-right: 10px;
                }
            }
            .precess-box {
                padding-left: 30px;
            }
            .process-create {
                display: flex;
                align-items: center;
                padding-bottom: 14px;
                .create-btn {
                    padding: 0 16px;
                    span {
                        margin-left: 0;
                    }
                }
                .icon-plus {
                    font-size: 20px;
                    line-height: normal;
                    font-weight: bold;
                    margin: 0 -4px;
                }
                .create-tips {
                    color: #979Ba5;
                    font-size: 14px;
                    padding-left: 10px;
                }
            }
            .process-empty {
                font-size: 14px;
            }
        }
        .info-footer {
            padding: 0 0 0 70px;
            display: flex;
            align-items: center;
            &.is-sticky {
                background-color: #fff;
                border-top: 1px solid $borderColor;
                padding: 10px 0 10px 70px;
            }
        }
    }
    .created-success {
        font-size: 14px;
        text-align: center;
        color: #444444;
        word-break: break-all;
        .icon-check-1 {
            width: 60px;
            height: 60px;
            line-height: 60px;
            font-size: 50px;
            font-weight: bold;
            color: #ffffff;
            border-radius: 50%;
            background-color: #2dcb56;
            margin-top: 12px;
        }
        .icon-exclamation {
            width: 18px;
            height: 18px;
            line-height: 17px;
            font-size: 12px;
            border: 1px solid #444444;
            border-radius: 50%;
            margin-top: -4px;
        }
        p {
            font-size: 24px;
            padding: 14px 0 18px;
        }
        .btn-box {
            padding: 20px 0 16px;
        }
    }
    .view-group {
        margin-bottom: 22px;
        padding-left: 30px;
        .view-info {
            width: 250px;
            margin-right: 20px;
            font-size: 14px;
            .info-label:after {
                content: "：";
            }
            .info-content {
                display: block;
                @include ellipsis;
            }
        }
    }
    .delete-disabled-btn {
        display: inline-block;
        vertical-align: middle;
        height: 32px;
        font-size: 14px;
        line-height: 32px;
        padding: 0 15px;
        background-color: #fff;
        border: 1px solid #dcdee5;
        border-radius: 2px;
        color: #c4c6cc;
        cursor: not-allowed;
    }
    .update-alert-layout {
        text-align: center;
        .bk-icon {
            width: 58px;
            height: 58px;
            line-height: 58px;
            font-size: 30px;
            color: #fff;
            border-radius: 50%;
            background-color: #2dcb56;
            margin: 8px 0 15px;
        }
        h3 {
            font-size: 24px;
            color: #313238;
            font-weight: normal;
            padding-bottom: 32px;
        }
        .btns {
            font-size: 0;
            padding-bottom: 20px;
        }
    }
    .cmdb-form-input.is-edit-name {
        width: 120px;
    }
    .bk-option-name {
        .category-id {
            color: #c4c6cc;
        }
    }
</style>
