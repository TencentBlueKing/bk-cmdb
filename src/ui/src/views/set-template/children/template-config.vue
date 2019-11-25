<template>
    <section class="layout" v-bkloading="{ isLoading: $loading('createSetTemplate') }">
        <div class="layout-row">
            <label class="row-label" :title="$t('模板名称')">{{$t('模板名称')}}</label>
            <template v-if="isViewMode">
                <i class="row-symbol">：</i>
                <span class="row-content row-view-content">{{templateName}}</span>
            </template>
            <template v-else>
                <i class="row-symbol required">*</i>
                <bk-input class="row-content"
                    data-vv-name="name"
                    font-size="medium"
                    v-validate="'required|singlechar|length:20'"
                    v-model.trim="templateName"
                    :placeholder="$t('集群模板名称占位符')">
                </bk-input>
            </template>
            <p class="row-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
        </div>
        <div class="layout-row template-row">
            <label class="row-label" :title="$t('集群拓扑')">{{$t('集群拓扑')}}</label>
            <i class="row-symbol" v-if="isViewMode">：</i>
            <i class="row-symbol required" v-else>*</i>
            <cmdb-set-template-tree class="row-content" ref="templateTree"
                :mode="mode"
                :template-id="templateId"
                @service-change="handleServiceChange"
                @service-selected="handleServiceSelected">
            </cmdb-set-template-tree>
            <input type="hidden" :value="services.length" v-validate="'min_value:1'" data-vv-name="service">
            <p class="row-error static" v-if="errors.has('service')">{{$t('请添加服务模板')}}</p>
        </div>
        <div class="template-options">
            <template v-if="isViewMode">
                <cmdb-auth :auth="$authResources({
                    resource_id: templateId,
                    type: $OPERATION.U_SET_TEMPLATE
                })">
                    <bk-button slot-scope="{ disabled }"
                        class="options-confirm"
                        theme="primary"
                        :disabled="disabled"
                        @click="handleEdit"
                    >
                        {{$t('编辑')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth :auth="$authResources({
                    resource_id: templateId,
                    type: $OPERATION.D_SET_TEMPLATE
                })">
                    <bk-button slot-scope="{ disabled }"
                        class="options-confirm"
                        :disabled="disabled"
                        @click="handleDelete"
                    >
                        {{$t('删除')}}
                    </bk-button>
                </cmdb-auth>
            </template>
            <template v-else>
                <cmdb-auth :auth="$authResources({
                    resource_id: templateId,
                    type: $OPERATION[mode === 'create' ? 'C_SET_TEMPLATE' : 'U_SET_TEMPLATE']
                })">
                    <bk-button slot-scope="{ disabled }"
                        class="options-confirm"
                        theme="primary"
                        :disabled="disabled || !hasChange"
                        @click="handleConfirm">
                        {{mode === 'create' ? $t('提交') : $t('保存')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button class="options-cancel" @click="handleCancel">{{$t('取消')}}</bk-button>
            </template>
        </div>
        <bk-dialog v-model="showUpdateInfo"
            :esc-close="false"
            :mask-close="false"
            :show-footer="false"
            :close-icon="false">
            <div class="update-alert-layout">
                <i class="bk-icon icon-check-1"></i>
                <h3>{{$t('修改成功')}}</h3>
                <div class="btns">
                    <bk-button class="mr10" theme="primary" v-if="isApplied" @click="handleToSyncInstance">{{$t('同步集群')}}</bk-button>
                    <bk-button class="mr10" :theme="isApplied ? 'default' : 'primary'" @click="handleToCreateInstance">{{$t('创建集群')}}</bk-button>
                    <bk-button theme="default" @click="handleBackToList">{{$t('返回列表')}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </section>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE, MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import cmdbSetTemplateTree from './template-tree.vue'
    export default {
        components: {
            cmdbSetTemplateTree
        },
        data () {
            return {
                templateName: '',
                originalTemplateName: '',
                serviceChange: false,
                services: [],
                validateServices: false,
                insideMode: null,
                showUpdateInfo: false
            }
        },
        computed: {
            mode () {
                return this.insideMode || this.$route.params.mode
            },
            isApplied () {
                return this.$route.params.isApplied
            },
            isViewMode () {
                return this.mode === 'view'
            },
            templateId () {
                return Number(this.$route.params.templateId) || null
            },
            hasChange () {
                if (this.mode !== 'edit') {
                    return true
                }
                return this.templateName !== this.originalTemplateName || this.serviceChange
            }
        },
        watch: {
            mode (mode) {
                this.errors.clear()
            },
            services (services) {
                if (!this.isViewMode) {
                    this.$nextTick(() => {
                        this.$validator.validate('service')
                    })
                }
            }
        },
        created () {
            this.setBreadcrumbs()
            if (['edit', 'view'].includes(this.mode)) {
                this.getSetTemplateInfo()
            }
        },
        methods: {
            setBreadcrumbs () {
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('集群模板'),
                    route: {
                        name: MENU_BUSINESS_SET_TEMPLATE
                    }
                }, {
                    label: this.mode === 'create' ? this.$t('创建') : this.templateName
                }])
                if (this.mode !== 'create') {
                    this.$store.commit('setTitle', this.templateName)
                }
            },
            async getSetTemplateInfo () {
                try {
                    const info = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        setTemplateId: this.templateId
                    })
                    this.templateName = info.name
                    this.originalTemplateName = info.name
                    this.setBreadcrumbs()
                } catch (e) {
                    console.error(e)
                }
            },
            handleEdit () {
                this.insideMode = 'edit'
            },
            async handleDelete () {
                const instance = this.$bkInfo({
                    title: this.$t('确认删除xx集群模板', { name: this.originalTemplateName }),
                    confirmFn: async () => {
                        try {
                            instance.buttonLoading = true
                            await this.$store.dispatch('setTemplate/deleteSetTemplate', {
                                bizId: this.$store.getters['objectBiz/bizId'],
                                config: {
                                    data: {
                                        set_template_ids: [this.templateId]
                                    }
                                }
                            })
                            this.$router.replace({
                                name: MENU_BUSINESS_SET_TEMPLATE
                            })
                        } catch (e) {
                            console.error(e)
                            instance.buttonLoading = false
                        }
                    }
                })
            },
            async handleConfirm () {
                try {
                    const validateResult = await this.$validator.validateAll()
                    if (!validateResult) {
                        return false
                    }
                    if (this.mode === 'create') {
                        await this.createSetTemplate()
                        this.alertCreateInfo()
                    } else {
                        await this.updateSetTemplate()
                        this.showUpdateInfo = true
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            createSetTemplate () {
                const services = this.$refs.templateTree.services
                const bizId = this.$store.getters['objectBiz/bizId']
                return this.$store.dispatch('setTemplate/createSetTemplate', {
                    bizId,
                    params: {
                        bk_biz_id: bizId,
                        name: this.templateName,
                        service_template_ids: services.map(template => template.id)
                    },
                    config: {
                        requestId: 'createSetTemplate'
                    }
                })
            },
            alertCreateInfo () {
                this.$bkInfo({
                    type: 'success',
                    title: this.$t('创建成功'),
                    okText: this.$t('创建集群'),
                    cancelText: this.$t('返回列表'),
                    confirmFn: () => {
                        this.$router.replace({ name: MENU_BUSINESS_HOST_AND_SERVICE })
                    },
                    cancelFn: () => {
                        this.$router.replace({ name: MENU_BUSINESS_SET_TEMPLATE })
                    }
                })
            },
            updateSetTemplate () {
                const services = this.$refs.templateTree.services
                const bizId = this.$store.getters['objectBiz/bizId']
                return this.$store.dispatch('setTemplate/updateSetTemplate', {
                    bizId,
                    setTemplateId: this.templateId,
                    params: {
                        bk_biz_id: bizId,
                        name: this.templateName,
                        service_template_ids: services.map(template => template.id)
                    },
                    config: {
                        requestId: 'updateSetTemplate'
                    }
                })
            },
            handleCancel () {
                if (this.insideMode) {
                    this.insideMode = null
                    this.$refs.templateTree.recoveryService()
                } else {
                    this.$router.replace({
                        name: MENU_BUSINESS_SET_TEMPLATE
                    })
                }
            },
            handleServiceChange (value) {
                this.serviceChange = value
            },
            handleServiceSelected (services) {
                this.services = services
            },
            handleBackToList () {
                this.$router.replace({ name: MENU_BUSINESS_SET_TEMPLATE })
            },
            handleToCreateInstance () {
                this.$router.replace({ name: MENU_BUSINESS_HOST_AND_SERVICE })
            },
            handleToSyncInstance () {
                this.showUpdateInfo = false
                this.insideMode = null
                this.$router.replace({
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'view',
                        templateId: this.templateId
                    },
                    query: {
                        tab: 'instance'
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        color: #63656E;
    }
    .layout-row {
        position: relative;
        font-size: 0;
        .row-label {
            width: 120px;
            text-align: right;
            font-size: 14px;
            line-height: 32px;
            padding: 0 0 0 10px;
            @include ellipsis;
            @include inlineBlock(top);
        }
        .row-symbol {
            line-height: 32px;
            font-size: 14px;
            width: 24px;
            text-indent: 4px;
            font-style: normal;
            @include inlineBlock(top);
            &.required {
                color: #EA3636;
            }
        }
        .row-content {
            width: 520px;
            line-height: 32px;
            color: #313238;
            font-size: 0px;
            @include inlineBlock(top);
            &.row-view-content {
                font-size: 14px;
            }
            /deep/ .bk-form-input {
                line-height: 32px;
            }
        }
        .row-error {
            position: absolute;
            color: #EA3636;
            padding-left: 145px;
            font-size: 12px;
            top: 100%;
            &.static {
                position: static;
            }
        }
    }
    .template-row {
        margin-top: 39px;
    }
    .template-options {
        padding: 20px 0 0 144px;
        font-size: 0;
        .options-confirm {
            margin-right: 10px;
        }
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
