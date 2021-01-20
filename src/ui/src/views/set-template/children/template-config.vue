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
                    v-validate="'required|singlechar|length:256'"
                    v-model.trim="templateName"
                    :placeholder="$t('请输入xx', { name: $t('模板名称') })">
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
                <cmdb-auth :auth="{ type: $OPERATION.U_SET_TEMPLATE, relation: [bizId, templateId] }">
                    <bk-button slot-scope="{ disabled }"
                        class="options-confirm"
                        theme="primary"
                        :disabled="disabled"
                        @click="handleEdit">
                        {{$t('编辑')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth :auth="{ type: $OPERATION.D_SET_TEMPLATE, relation: [bizId, templateId] }">
                    <bk-button slot-scope="{ disabled }"
                        class="options-confirm"
                        :disabled="disabled"
                        @click="handleDelete">
                        {{$t('删除')}}
                    </bk-button>
                </cmdb-auth>
            </template>
            <template v-else>
                <cmdb-auth :auth=" mode === 'create' ? { type: $OPERATION.C_SET_TEMPLATE, relation: [bizId] } : { type: $OPERATION.U_SET_TEMPLATE, relation: [bizId, templateId] }">
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
            width="480"
            :esc-close="false"
            :mask-close="false"
            :show-footer="false"
            :close-icon="false">
            <div class="update-alert-layout">
                <i class="bk-icon icon-check-1"></i>
                <h3>{{$t('修改成功')}}</h3>
                <p class="update-success-tips">{{$t(serviceChange ? '集群模板修改成功并需要同步提示' : '集群模板修改成功提示')}}</p>
                <div class="btns">
                    <bk-button class="mr10" theme="primary"
                        v-if="isApplied && serviceChange"
                        @click="handleToSyncInstance">
                        {{$t('同步集群')}}
                    </bk-button>
                    <bk-button class="mr10"
                        :theme="isApplied && serviceChange ? 'default' : 'primary'"
                        @click="handleToCreateInstance">
                        {{$t('创建集群')}}
                    </bk-button>
                    <bk-button theme="default" @click="handleClose">{{$t('关闭按钮')}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </section>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE, MENU_BUSINESS_HOST_AND_SERVICE } from '@/dictionary/menu-symbol'
    import cmdbSetTemplateTree from './template-tree.vue'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            cmdbSetTemplateTree
        },
        data () {
            return {
                templateInfo: null,
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
            ...mapGetters('objectBiz', ['bizId']),
            mode () {
                return this.insideMode || this.$route.params.mode
            },
            isApplied () {
                return this.templateInfo && this.templateInfo.set_instance_count > 0
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
            '$route.query': {
                immediate: true,
                handler (query) {
                    if (query.edit) {
                        this.handleEdit()
                    }
                }
            },
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
                if (this.mode !== 'create') {
                    this.$store.commit('setTitle', this.templateName)
                }
            },
            async getSetTemplateInfo () {
                try {
                    const data = await this.$store.dispatch('setTemplate/getSetTemplates', {
                        bizId: this.bizId,
                        params: {
                            set_template_ids: [this.templateId]
                        }
                    })
                    const info = data.info[0]
                    this.templateInfo = info
                    this.templateName = info.set_template.name
                    this.originalTemplateName = info.set_template.name
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
                    title: this.$t('确认删除'),
                    subTitle: this.$t('确认删除xx集群模板', { name: this.originalTemplateName }),
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
                            this.$routerActions.redirect({
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
                    width: 480,
                    title: this.$t('创建成功'),
                    subTitle: this.$t('创建集群模板成功提示'),
                    okText: this.$t('创建集群'),
                    cancelText: this.$t('返回列表'),
                    confirmFn: () => {
                        this.$routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
                    },
                    cancelFn: () => {
                        this.$routerActions.redirect({ name: MENU_BUSINESS_SET_TEMPLATE })
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
                    const breadcrumbs = window.CMDB_APP.$children[0].$refs.topView.$refs.breadcrumbs
                    if (breadcrumbs.from) {
                        breadcrumbs.handleClick()
                    } else {
                        this.$routerActions.redirect({
                            name: MENU_BUSINESS_SET_TEMPLATE
                        })
                    }
                }
            },
            handleServiceChange (value) {
                this.serviceChange = value
            },
            handleServiceSelected (services) {
                this.services = services
            },
            handleClose () {
                this.showUpdateInfo = false
                this.insideMode = null
                this.$routerActions.redirect({
                    name: 'setTemplateConfig',
                    params: {
                        mode: 'view',
                        templateId: this.templateId
                    },
                    history: false
                })
            },
            handleToCreateInstance () {
                this.$routerActions.redirect({ name: MENU_BUSINESS_HOST_AND_SERVICE })
            },
            handleToSyncInstance () {
                this.showUpdateInfo = false
                this.insideMode = null
                this.$routerActions.redirect({
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
            font-size: 50px;
            color: #fff;
            border-radius: 50%;
            background-color: #2dcb56;
            margin: 8px 0 15px;
        }
        h3 {
            font-size: 24px;
            color: #313238;
            font-weight: normal;
            padding-bottom: 16px;
        }
        .btns {
            font-size: 0;
            padding-bottom: 20px;
            .bk-button {
                min-width: 86px;
            }
        }
        .update-success-tips {
            padding-bottom: 24px;
        }
    }
</style>
