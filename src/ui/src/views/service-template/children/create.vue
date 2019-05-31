<template>
    <div class="create-template-wrapper">
        <div class="info-group">
            <h3>基本属性</h3>
            <div class="form-info clearfix">
                <label class="label-text fl" for="templateName">
                    {{$t('ServiceManagement["模板名称"]')}}
                    <span class="color-danger">*</span>
                </label>
                <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('templateName') }">
                    <input type="text" class="cmdb-form-input" id="templateName"
                        name="templateName"
                        :placeholder="$t('ServiceManagement[\'请输入模版名称\']')"
                        :disabled="!isCreatedType"
                        v-model.trim="formData.tempalteName"
                        v-validate="'required|singlechar'">
                    <p class="form-error">{{errors.first('templateName')}}</p>
                </div>
            </div>
            <div class="form-info clearfix">
                <span class="label-text fl">
                    {{$t('ServiceManagement["服务分类"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('mainClassificationId') }" style="width: auto;">
                    <cmdb-selector
                        class="fl"
                        :placeholder="$t('ServiceManagement[\'请选择一级分类\']')"
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
                        :placeholder="$t('ServiceManagement[\'请选择二级分类\']')"
                        :auto-select="true"
                        :list="secondaryList"
                        v-validate="'required'"
                        name="secondaryClassificationId"
                        v-model="formData['secondaryClassification']">
                    </cmdb-selector>
                    <p class="form-error">{{errors.first('secondaryClassificationId')}}</p>
                </div>
            </div>
        </div>
        <div class="info-group">
            <h3>进程服务</h3>
            <div class="precess-box">
                <div class="process-create">
                    <bk-button class="create-btn" @click="handleCreateProcess">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t("ServiceManagement['新建进程']")}}</span>
                    </bk-button>
                    <span class="create-tips">{{$t("ServiceManagement['新建进程提示']")}}</span>
                </div>
                <process-table
                    :properties="properties"
                    @on-edit="handleUpdateProcess"
                    @on-delete="handleDeleteProcess"
                    :list="processList">
                </process-table>
                <div class="btn-box">
                    <bk-button type="primary" @click="handleSubmit">{{$t("Common['确定']")}}</bk-button>
                    <bk-button @click="handleCancelOperation">{{$t("Common['取消']")}}</bk-button>
                </div>
            </div>
        </div>
        <cmdb-slider :is-show.sync="slider.show" :title="slider.title">
            <template slot="content">
                <process-form
                    :properties="properties"
                    :property-groups="propertyGroups"
                    :inst="attribute.inst.edit"
                    :type="attribute.type"
                    :save-disabled="false"
                    @on-submit="handleSaveProcess"
                    @on-cancel="handleCancelProcess">
                </process-form>
            </template>
        </cmdb-slider>
    </div>
</template>

<script>
    import processForm from './process-form.vue'
    import processTable from './process'
    import { mapActions, mapGetters, mapMutations } from 'vuex'
    export default {
        components: {
            processTable,
            processForm
        },
        data () {
            return {
                properties: [],
                propertyGroups: [],
                objectUnique: [],
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
                mainList: [],
                secondaryList: [],
                allSecondaryList: [],
                formData: {
                    mainClassification: '',
                    secondaryClassification: '',
                    tempalteName: ''
                }
            }
        },
        computed: {
            ...mapGetters('serviceProcess', ['localProcessTemplate']),
            processList () {
                return this.localProcessTemplate
            },
            isCreatedType () {
                return !this.$route.params['template']
            },
            originTemplateValues () {
                return this.$route.params['template']
            }
        },
        async created () {
            const title = this.isCreatedType ? this.$t("ServiceManagement['新建服务模版']") : this.originTemplateValues['name']
            this.$store.commit('setHeaderTitle', title)
            try {
                await this.reload()
                if (!this.isCreatedType) {
                    this.initFill()
                }
            } catch (e) {
                console.log(e)
            }
        },
        methods: {
            ...mapMutations('serviceProcess', ['deleteLocalProcessTemplate', 'clearLocalProcessTemplate']),
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectUnique', ['searchObjectUniqueConstraints']),
            ...mapActions('serviceClassification', ['searchServiceCategory']),
            ...mapActions('serviceTemplate', ['operationServiceTemplate']),
            ...mapActions('processTemplate', ['createProcessTemplate']),
            initFill () {
                this.formData.tempalteName = this.originTemplateValues['name']
                this.formData.mainClassification = this.allSecondaryList.filter(classification => classification['id'] === this.originTemplateValues['service_category_id'])[0]['parent_id']
                this.formData.secondaryClassification = this.originTemplateValues['service_category_id']
            },
            async reload () {
                this.properties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: 'process',
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_process`,
                        fromCache: false
                    }
                })
                await this.getServiceClassification()
                this.getPropertyGroups()
                this.getObjectUnique()
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: 'process',
                    params: this.$injectMetadata(),
                    config: {
                        fromCache: false,
                        requestId: 'post_searchGroup_process'
                    }
                }).then(groups => {
                    this.propertyGroups = groups
                    return groups
                })
            },
            getObjectUnique () {
                this.objectUnique = this.searchObjectUniqueConstraints({
                    objId: 'process',
                    params: {},
                    config: {
                        requestId: 'searchObjectUniqueConstraints'
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
                this.mainList = result.info.filter(classification => !classification['parent_id'])
                this.allSecondaryList = result.info.filter(classification => classification['parent_id'])
            },
            handleSelect (id, data) {
                this.secondaryList = this.allSecondaryList.filter(classification => classification['parent_id'] === id && classification['root_id'] === id)
                if (!this.secondaryList.length) {
                    this.formData.secondaryClassification = ''
                }
            },
            handleSaveProcess (values, changedValues, originValues, type) {
                console.log(values)
                console.log(changedValues)
                console.log(originValues)
                // if (type === 'create') {}
            },
            handleCancelProcess () {
                this.slider.show = false
            },
            handleCreateProcess () {
                this.slider.show = true
                this.slider.title = this.$t("ProcessManagement['添加进程']")
                this.attribute.type = 'create'
                this.attribute.inst.edit = {}
            },
            handleUpdateProcess (template) {
                this.slider.show = true
                this.slider.title = template['bk_func_name']
                this.attribute.type = 'update'
                this.attribute.inst.edit = template
            },
            handleDeleteProcess (template) {
                if (this.isCreatedType) {
                    this.$bkInfo({
                        title: this.$t("ServiceManagement['确认删除模板进程']"),
                        confirmFn: () => {
                            this.deleteClocalProcessTemplate(template)
                        }
                    })
                }
            },
            async handleSubmit () {
                if (!await this.$validator.validateAll()) return
                this.operationServiceTemplate({
                    params: this.$injectMetadata({
                        name: this.formData.tempalteName,
                        service_category_id: this.formData.secondaryClassification
                    })
                }).then(data => {
                    this.createProcessTemplate({
                        params: this.$injectMetadata({
                            service_template_id: data.id,
                            processes: this.processList
                        })
                    }).then(() => {
                        this.$bkMessage({
                            message: this.$t("Common['保存成功']"),
                            theme: 'success'
                        })
                        this.clearLocalProcessTemplate()
                        this.handleCancelOperation()
                    })
                })
            },
            handleCancelOperation () {
                this.$router.push({ name: 'serviceTemplate' })
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
                .bk-selector {
                    width: 254px;
                    margin-right: 10px;
                }
            }
            .precess-box {
                padding-left: 30px;
            }
            .process-create {
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
                    font-size: 12px;
                    padding-left: 10px;
                }
            }
            .btn-box {
                padding-top: 30px;
            }
        }
    }
</style>
