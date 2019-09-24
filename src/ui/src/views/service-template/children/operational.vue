<template>
    <div class="create-template-wrapper">
        <div class="info-group">
            <h3>{{$t('基本属性')}}</h3>
            <div class="form-info clearfix">
                <label class="label-text fl" for="templateName">
                    {{$t('模板名称')}}
                    <span class="color-danger">*</span>
                </label>
                <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('templateName') }">
                    <bk-input type="text" class="cmdb-form-input" id="templateName"
                        name="templateName"
                        :placeholder="$t('请输入模板名称')"
                        :disabled="!isCreatedType"
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
        </div>
        <div class="info-group">
            <h3>{{$t('服务进程')}}</h3>
            <div class="precess-box">
                <div class="process-create">
                    <span
                        v-cursor="{
                            active: !$isAuthorized(auth),
                            auth: [auth]
                        }">
                        <bk-button class="create-btn" :disabled="!$isAuthorized(auth)" @click="handleCreateProcess">
                            <i class="bk-icon icon-plus"></i>
                            <span>{{$t('新建进程')}}</span>
                        </bk-button>
                    </span>
                    <span class="create-tips">{{$t('新建进程提示')}}</span>
                </div>
                <process-table
                    v-if="processList.length"
                    :loading="processLoading"
                    :properties="properties"
                    @on-edit="handleUpdateProcess"
                    @on-delete="handleDeleteProcess"
                    :list="processList">
                </process-table>
                <div class="btn-box">
                    <span
                        v-cursor="{
                            active: !$isAuthorized(auth),
                            auth: [auth]
                        }">
                        <bk-button theme="primary"
                            :disabled="!$isAuthorized(auth)"
                            @click="handleSubmit">
                            {{isCreatedType ? $t('提交') : $t('保存')}}
                        </bk-button>
                    </span>
                    <bk-button @click="handleReturn">{{$t('取消')}}</bk-button>
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
                    :properties="properties"
                    :property-groups="propertyGroups"
                    :inst="attribute.inst.edit"
                    :type="attribute.type"
                    :is-created-service="isCreatedType"
                    :save-disabled="false"
                    :has-used="hasUsed"
                    @on-submit="handleSaveProcess"
                    @on-cancel="handleCancelProcess">
                </process-form>
            </template>
        </bk-sideslider>
        <bk-dialog
            v-model="createdSucess.show"
            :width="490"
            :close-icon="false"
            :show-footer="false"
            :title="createdSucess.title">
            <div class="created-success">
                <div class="content">
                    <i class="bk-icon icon-check-1"></i>
                    <p>{{$t('服务模板创建成功')}}</p>
                    <span v-if="processList.length">{{$tc('创建成功前往服务拓扑', createdSucess.name, { name: createdSucess.name })}}</span>
                    <span v-else>
                        <i class="bk-icon icon-exclamation"></i>
                        {{$t('创建成功无进程提示')}}
                    </span>
                </div>
                <div class="btn-box">
                    <bk-button v-if="processList.length"
                        theme="primary"
                        class="mr10"
                        @click="handleGoInstance">
                        {{$t('创建服务实例')}}
                    </bk-button>
                    <bk-button @click="handleCancelOperation">{{$t('返回列表')}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import processForm from './process-form.vue'
    import processTable from './process'
    import { mapActions, mapGetters, mapMutations } from 'vuex'
    import { MENU_BUSINESS_SERVICE_TEMPLATE, MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
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
                createdSucess: {
                    show: false,
                    name: ''
                },
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
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('serviceProcess', ['localProcessTemplate']),
            isCreatedType () {
                return !this.$route.params['templateId']
            },
            templateId () {
                return this.$route.params['templateId']
            },
            auth () {
                if (this.isCreatedType) {
                    return this.$OPERATION.C_SERVICE_TEMPLATE
                }
                return this.$OPERATION.U_SERVICE_TEMPLATE
            }
        },
        async created () {
            this.setBreadcrumbs()
            this.processList = this.localProcessTemplate
            try {
                await this.reload()
                if (!this.isCreatedType) {
                    this.initEdit()
                }
            } catch (e) {
                console.log(e)
            }
        },
        beforeDestroy () {
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
            setBreadcrumbs () {
                if (this.isCreatedType) {
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
                this.hasUsed = this.isCreatedType ? false : Boolean(this.originTemplateValues['service_instance_count'])
                this.getProcessList()
            },
            async reload () {
                if (!this.isCreatedType) {
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
                this.originTemplateValues = await this.findServiceTemplate({
                    id: this.templateId,
                    config: {
                        globalError: false,
                        cancelPrevious: true,
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
            },
            async getServiceClassification () {
                const result = await this.searchServiceCategory({
                    params: this.$injectMetadata(),
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
                    })
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
                if (type === 'create') {
                    this.createProcessTemplate({
                        params: this.$injectMetadata({
                            service_template_id: this.originTemplateValues['id'],
                            processes: [{
                                spec: values
                            }]
                        })
                    }).then(() => {
                        this.getProcessList()
                        this.handleCancelProcess()
                    })
                } else {
                    this.updateProcessTemplate({
                        params: this.$injectMetadata({
                            process_template_id: values['process_id'],
                            process_property: changedValues
                        })
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
                this.slider.show = true
                this.slider.title = template['bk_func_name']['value']
                this.attribute.type = 'update'
                this.attribute.inst.edit = template
            },
            handleDeleteProcess (template) {
                this.$bkInfo({
                    title: this.$t('确认删除模板进程'),
                    confirmFn: () => {
                        if (this.isCreatedType) {
                            this.deleteLocalProcessTemplate(template)
                            this.processList = this.localProcessTemplate
                        } else {
                            this.deleteProcessTemplate({
                                params: {
                                    data: this.$injectMetadata({
                                        process_templates: [template['process_id']]
                                    })
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
                            return {
                                spec: process
                            }
                        })
                    })
                }).then(() => {
                    if (this.isCreatedType) {
                        this.createdSucess.show = true
                    } else {
                        this.handleCancelOperation()
                    }
                })
            },
            async handleSubmit () {
                if (!await this.$validator.validateAll()) return
                if (this.formData.templateId) {
                    this.updateServiceTemplate({
                        params: this.$injectMetadata({
                            id: this.formData.templateId,
                            name: this.formData.templateName,
                            service_category_id: this.formData.secondaryClassification
                        })
                    }).then(() => {
                        if (this.isCreatedType) {
                            this.handleSubmitProcessList()
                        } else {
                            this.$success(this.$t('保存成功'))
                            this.handleCancelOperation()
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
                    })
                }).then(data => {
                    this.createdSucess.name = data.name
                    this.formData.templateId = data.id
                    this.handleSubmitProcessList()
                })
            },
            handleGoInstance () {
                this.$router.replace({ name: MENU_BUSINESS_SERVICE_TOPOLOGY })
            },
            handleReturn () {
                const moduleId = this.$route.params['moduleId']
                if (moduleId) {
                    this.$router.replace({
                        name: MENU_BUSINESS_SERVICE_TOPOLOGY,
                        query: {
                            module: this.$route.params.moduleId
                        }
                    })
                } else {
                    this.handleCancelOperation()
                }
            },
            handleCancelOperation () {
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
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-template-wrapper {
        padding: 0 20px;
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
            padding: 14px 0 24px;
        }
        .btn-box {
            padding: 32px 0 36px;
        }
    }
</style>
