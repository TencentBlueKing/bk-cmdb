<template>
    <div class="create-template-wrapper" v-bkloading="{ isLoading: $loading('getSingleServiceTemplate') }">
        <div class="info-group">
            <h3>{{$t('基本属性')}}</h3>

            <template v-if="isFormMode">
                <div class="form-info clearfix">
                    <label class="label-text fl" for="templateName">
                        {{$t('模板名称')}}
                        <span class="color-danger">*</span>
                    </label>
                    <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('templateName') }">
                        <bk-input type="text" class="cmdb-form-input" id="templateName"
                            name="templateName"
                            :placeholder="$t('请输入模板名称')"
                            :disabled="!isCreateMode"
                            v-model.trim="formData.templateName"
                            v-validate="'required|singlechar|length:256'">
                        </bk-input>
                        <p class="form-error">{{errors.first('templateName')}}</p>
                    </div>
                </div>
                <div class="form-info clearfix">
                    <span class="label-text fl">
                        {{$t('服务分类')}}
                        <span class="color-danger">*</span>
                    </span>
                    <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('mainClassificationId') }" style="width: auto;">
                        <cmdb-selector
                            class="fl"
                            :placeholder="$t('请选择一级分类')"
                            :auto-select="false"
                            :list="mainList"
                            v-validate="'required'"
                            name="mainClassificationId"
                            v-model="formData['mainClassification']"
                            @on-selected="handleSelect">
                        </cmdb-selector>
                        <p class="form-error">{{errors.first('mainClassificationId')}}</p>
                    </div>
                    <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('secondaryClassificationId') }" style="width: auto;">
                        <cmdb-selector
                            class="fl"
                            :placeholder="$t('请选择二级分类')"
                            :auto-select="true"
                            :list="secondaryList"
                            :empty-text="emptyText"
                            v-validate="'required'"
                            name="secondaryClassificationId"
                            v-model="formData['secondaryClassification']">
                        </cmdb-selector>
                        <p class="form-error">{{errors.first('secondaryClassificationId')}}</p>
                    </div>
                </div>
            </template>

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
                    <cmdb-auth :auth="$authResources(auth)">
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
                    @on-edit="handleUpdateProcess"
                    @on-delete="handleDeleteProcess"
                    :list="processList">
                </process-table>
                <div class="btn-box">
                    <cmdb-auth class="mr5" :auth="$authResources(auth)">
                        <bk-button slot-scope="{ disabled }"
                            theme="primary"
                            :disabled="disabled"
                            @click="handleSubmit">
                            {{getButtonText()}}
                        </bk-button>
                    </cmdb-auth>
                    <cmdb-auth class="mr5"
                        :auth="$authResources(auth)"
                        v-if="!isFormMode && deletable">
                        <bk-button slot-scope="{ disabled }"
                            :disabled="disabled"
                            @click="handleDeleteTemplate">
                            {{$t('删除')}}
                        </bk-button>
                    </cmdb-auth>
                    <span class="delete-disabled-btn"
                        v-else-if="!isFormMode && !deletable"
                        v-bk-tooltips.top="$t('不可删除')">
                        {{$t('删除')}}
                    </span>
                    <bk-button @click="handleReturn" v-show="isFormMode">{{$t('取消')}}</bk-button>
                </div>
            </div>
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
                    :save-disabled="false"
                    :has-used="hasUsed"
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
    </div>
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
                emptyText: this.$t('请选择一级分类'),
                attribute: {
                    type: null,
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
                    secondaryClassification: '',
                    templateName: '',
                    templateId: ''
                },
                insideMode: 'view',
                deletable: false,
                showUpdateInfo: false
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('serviceProcess', ['localProcessTemplate']),
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
                        type: this.$OPERATION.C_SERVICE_TEMPLATE
                    }
                }
                return {
                    resource_id: Number(this.templateId) || null,
                    type: this.$OPERATION.U_SERVICE_TEMPLATE
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
                'findServiceTemplate'
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
                    if (!this.isCreateMode) {
                        this.initEdit()
                    }
                    if (this.setActive) {
                        Bus.$emit('active-change', 'instance')
                        this.$route.params.active = null
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            setBreadcrumbs () {
                if (this.isCreateMode) {
                    this.$store.commit('setTitle', this.$t('新建模板'))
                    this.$store.commit('setBreadcrumbs', [{
                        label: this.$t('服务模板'),
                        route: {
                            name: MENU_BUSINESS_SERVICE_TEMPLATE
                        }
                    }, {
                        label: this.$t('新建模板')
                    }])
                } else {
                    this.$store.commit('setTitle', this.$t('模板详情'))
                }
            },
            initEdit () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('服务模板'),
                    route: {
                        name: MENU_BUSINESS_SERVICE_TEMPLATE
                    }
                }, {
                    label: this.originTemplateValues.name
                }])
                this.formData.templateId = this.originTemplateValues['id']
                this.formData.templateName = this.originTemplateValues['name']
                this.formData.mainClassification = this.allSecondaryList.filter(classification => classification['id'] === this.originTemplateValues['service_category_id'])[0]['bk_parent_id']
                this.formData.secondaryClassification = this.originTemplateValues['service_category_id']
                this.hasUsed = this.isCreateMode ? false : Boolean(this.originTemplateValues['service_instance_count'])
                this.getProcessList()
            },
            async reload () {
                if (!this.isCreateMode) {
                    this.getSingleServiceTemplate()
                }
                this.properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: 'process',
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: `searchObjectAttribute_templateProcess`,
                        fromCache: false
                    }
                })
                await this.getServiceClassification()
                this.getPropertyGroups()
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: 'process',
                    params: this.$injectMetadata(),
                    config: {
                        fromCache: true,
                        requestId: 'post_searchGroup_process'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            async getSingleServiceTemplate () {
                try {
                    this.originTemplateValues = await this.findServiceTemplate({
                        id: this.templateId,
                        config: {
                            requestId: 'getSingleServiceTemplate',
                            globalError: false,
                            transformData: false
                        }
                    }).then(res => {
                        if (!res.result) {
                            this.$router.replace({ name: '404' })
                        } else {
                            return {
                                service_instance_count: res.data.service_instance_count,
                                process_instance_count: res.data.process_instance_count,
                                ...res.data.template
                            }
                        }
                    })
                } catch (e) {
                    console.error(e)
                    this.$router.replace({ name: '404' })
                }
            },
            async getServiceClassification () {
                const result = await this.searchServiceCategory({
                    params: this.$injectMetadata({}, { injectBizId: true }),
                    config: {
                        requestId: 'get_proc_services_categories'
                    }
                })
                const categoryList = result.info.map(item => item['category'])
                this.mainList = categoryList.filter(classification => !classification['bk_parent_id'])
                this.allSecondaryList = categoryList.filter(classification => classification['bk_parent_id'])
            },
            getProcessList () {
                this.processLoading = true
                this.getBatchProcessTemplate({
                    params: this.$injectMetadata({
                        service_template_id: this.originTemplateValues['id']
                    }, { injectBizId: true })
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
                this.emptyText = this.$t('没有二级分类')
                if (!this.secondaryList.length) {
                    this.formData.secondaryClassification = ''
                }
            },
            handleSaveProcess (values, changedValues, type) {
                const processValues = type === 'create' ? values : changedValues
                if (processValues.hasOwnProperty('protocol') && !processValues.protocol.value) {
                    processValues.protocol.value = null
                }
                if (type === 'create') {
                    this.createProcessTemplate({
                        params: this.$injectMetadata({
                            service_template_id: this.originTemplateValues['id'],
                            processes: [{
                                spec: processValues
                            }]
                        }, { injectBizId: true })
                    }).then(() => {
                        this.getProcessList()
                        this.handleCancelProcess()
                    })
                } else {
                    this.updateProcessTemplate({
                        params: this.$injectMetadata({
                            process_template_id: values['process_id'],
                            process_property: processValues
                        }, { injectBizId: true })
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
            handleUpdateProcess (template) {
                try {
                    this.slider.show = true
                    this.slider.title = template['bk_func_name']['value']
                    this.attribute.type = 'update'
                    this.attribute.inst.edit = template
                } catch (e) {
                    console.error(e)
                }
            },
            handleDeleteProcess (template) {
                this.$bkInfo({
                    title: this.$t('确认删除模板进程'),
                    confirmFn: () => {
                        if (this.isCreateMode) {
                            this.deleteLocalProcessTemplate(template)
                            this.processList = this.localProcessTemplate
                        } else {
                            this.deleteProcessTemplate({
                                params: {
                                    data: this.$injectMetadata({
                                        process_templates: [template['process_id']]
                                    }, { injectBizId: true })
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
                    params: this.$injectMetadata({
                        service_template_id: this.formData.templateId,
                        processes: this.processList.map(process => {
                            if (process.hasOwnProperty('protocol') && !process.protocol.value) {
                                process.protocol.value = null
                            }
                            delete process.sign_id
                            return {
                                spec: process
                            }
                        })
                    }, { injectBizId: true })
                }).then(() => {
                    this.$success(this.$t('创建成功'))
                    this.handleCancelOperation()
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
                        params: this.$injectMetadata({
                            id: this.formData.templateId,
                            name: this.formData.templateName,
                            service_category_id: this.formData.secondaryClassification
                        }, { injectBizId: true })
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
                            title: this.$t('服务模板创建没进程提示'),
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
                    params: this.$injectMetadata({
                        name: this.formData.templateName,
                        service_category_id: this.formData.secondaryClassification
                    }, { injectBizId: true })
                }).then(data => {
                    this.formData.templateId = data.id
                    this.handleSubmitProcessList()
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
                    this.$router.replace({
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
                this.$router.replace({ name: MENU_BUSINESS_SERVICE_TEMPLATE })
            },
            handleSliderBeforeClose () {
                const hasChanged = this.$refs.processForm.hasChange()
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
            handleDeleteTemplate () {
                this.$bkInfo({
                    title: this.$t('确认删除模板'),
                    subTitle: this.$tc('即将删除服务模板', name, { name: this.formData.templateName }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        await this.$store.dispatch('serviceTemplate/deleteServiceTemplate', {
                            params: {
                                data: this.$injectMetadata({
                                    service_template_id: this.templateId
                                }, {
                                    injectBizId: true
                                })
                            },
                            config: {
                                requestId: 'delete_proc_service_template'
                            }
                        })
                        this.$success(this.$t('删除成功'))
                        this.handleReturn()
                    }
                })
            },
            getServiceCategory () {
                const first = this.mainList.find(first => first.id === this.formData.mainClassification) || {}
                const second = this.allSecondaryList.find(second => second.id === this.formData.secondaryClassification) || {}
                return `${first.name} / ${second.name}`
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-template-wrapper {
        .info-group {
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
                    font-size: 12px;
                    line-height: normal;
                    font-weight: bold;
                }
                .create-tips {
                    color: #979Ba5;
                    font-size: 14px;
                    padding-left: 10px;
                }
            }
            .btn-box {
                padding-top: 30px;
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
            font-size: 30px;
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
</style>
