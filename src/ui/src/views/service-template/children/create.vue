<template>
    <div class="create-template-wrapper">
        <div class="info-group">
            <h3>基本属性</h3>
            <div class="form-info clearfix">
                <label class="label-text fl" for="templateName">
                    {{$t('ServiceManagement["模板名称"]')}}
                    <span class="color-danger">*</span>
                </label>
                <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('fieldId') }">
                    <input type="text" class="cmdb-form-input" id="templateName"
                        name="fieldId"
                        :placeholder="$t('ServiceManagement[\'请输入模版名称\']')"
                        v-model.trim="formData.tempalteName"
                        v-validate="'required|fieldId'">
                    <p class="form-error">{{errors.first('fieldId')}}</p>
                </div>
            </div>
            <div class="form-info clearfix">
                <span class="label-text fl">
                    {{$t('ServiceManagement["服务分类"]')}}
                    <span class="color-danger">*</span>
                </span>
                <div class="cmdb-form-item fl" :class="{ 'is-error': errors.has('classificationId') }" style="width: auto;">
                    <bk-selector
                        class="fl"
                        placeholder="请选择一级分类"
                        :list="[]"
                        v-validate="'required'"
                        name="classificationId"
                        :selected.sync="formData.primaryClassification">
                    </bk-selector>
                    <bk-selector
                        class="fl"
                        placeholder="请选择二级分类"
                        :list="[]"
                        v-validate="'required'"
                        name="classificationId"
                        :selected.sync="formData.secondaryClassification">
                    </bk-selector>
                    <p class="form-error">{{errors.first('classificationId')}}</p>
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
                <process-table></process-table>
                <div class="btn-box">
                    <bk-button type="primary" @click="handleSubmit">{{$t("Common['确定']")}}</bk-button>
                    <bk-button>{{$t("Common['取消']")}}</bk-button>
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
                    @on-submit="handleSave"
                    @on-cancel="handleCancel">
                </process-form>
            </template>
        </cmdb-slider>
    </div>
</template>

<script>
    import processForm from './process-form.vue'
    import processTable from './process'
    import { mapActions } from 'vuex'
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
                formData: {
                    primaryClassification: '',
                    secondaryClassification: '',
                    tempalteName: ''
                }
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t("ServiceManagement['新建服务模版']"))
            this.reload()
        },
        methods: {
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectUnique', ['searchObjectUniqueConstraints']),
            ...mapActions('procConfig', [
                'searchProcess'
            ]),
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
            async getObjectUnique () {
                this.objectUnique = await this.searchObjectUniqueConstraints({
                    objId: 'process',
                    params: {},
                    config: {
                        requestId: 'searchObjectUniqueConstraints'
                    }
                })
            },
            handleSave () {

            },
            handleCancel () {
                this.slider.show = false
            },
            handleCreateProcess () {
                this.slider.show = true
                this.slider.title = this.$t("ProcessManagement['添加进程']")
            },
            handleSubmit () {
                this.$validator.validateAll().then(result => {
                    console.log(result)
                })
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
