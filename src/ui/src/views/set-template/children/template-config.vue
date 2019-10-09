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
                @service-change="handleServiceChange">
            </cmdb-set-template-tree>
        </div>
        <div class="template-options">
            <template v-if="isViewMode">
                <bk-button class="options-confirm" theme="primary"
                    @click="handleEdit">
                    {{$t('编辑')}}
                </bk-button>
            </template>
            <template v-else>
                <bk-button class="options-confirm" theme="primary"
                    :disabled="!hasChange"
                    @click="handleConfirm">
                    {{$t('确定')}}
                </bk-button>
                <bk-button class="options-cancel" @click="handleCancel">{{$t('取消')}}</bk-button>
            </template>
        </div>
    </section>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE, MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
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
                insideMode: null
            }
        },
        computed: {
            mode () {
                return this.insideMode || this.$route.params.mode
            },
            isViewMode () {
                return this.mode === 'view'
            },
            templateId () {
                return this.$route.params.templateId
            },
            hasChange () {
                if (this.mode !== 'edit') {
                    return true
                }
                return this.templateName !== this.originalTemplateName || this.serviceChange
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
                const titleMap = {
                    create: this.$t('创建'),
                    edit: this.$t('编辑'),
                    view: this.$t('查看')
                }
                this.$store.commit('setBreadcrumbs', [{
                    label: this.$t('集群模板'),
                    route: {
                        name: MENU_BUSINESS_SET_TEMPLATE
                    }
                }, {
                    label: titleMap[this.mode]
                }])
            },
            async getSetTemplateInfo () {
                try {
                    const info = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        setTemplateId: this.templateId
                    })
                    this.templateName = info.name
                    this.originalTemplateName = info.name
                } catch (e) {
                    console.error(e)
                }
            },
            handleEdit () {
                this.insideMode = 'edit'
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
                        this.alertUpdateInfo()
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
                    okText: this.$t('创建服务实例'),
                    cancelText: this.$t('返回列表'),
                    confirmFn: () => {
                        this.$router.replace({ name: MENU_BUSINESS_SERVICE_TOPOLOGY })
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
            alertUpdateInfo () {
                if (this.insideMode) {
                    this.insideMode = null
                    this.$success(this.$t('修改成功'))
                } else {
                    this.$bkInfo({
                        type: 'success',
                        title: this.$t('修改成功'),
                        okText: this.$t('同步实例'),
                        cancelText: this.$t('返回列表'),
                        confirmFn: () => {
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
                        },
                        cancelFn: () => {
                            this.$router.replace({ name: MENU_BUSINESS_SET_TEMPLATE })
                        }
                    })
                }
            },
            handleCancel () {
                if (this.insideMode) {
                    this.insideMode = null
                } else {
                    this.$router.replace({
                        name: MENU_BUSINESS_SET_TEMPLATE
                    })
                }
            },
            handleServiceChange (value) {
                this.serviceChange = value
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
            width: 14px;
            text-align: center;
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
        }
        .row-error {
            position: absolute;
            color: #EA3636;
            padding-left: 145px;
            font-size: 12px;
            top: 100%;
        }
    }
    .template-row {
        margin-top: 39px;
    }
    .template-options {
        padding:23px 0 0 134px;
        font-size: 0;
        .options-confirm {
            margin-right: 10px;
        }
    }
</style>
