<template>
    <section class="layout" v-bkloading="{ isLoading: $loading('createSetTemplate') }">
        <div class="layout-row">
            <label class="row-label inline-block-middle" :title="$t('模板名称')">{{$t('模板名称')}}</label>
            <i class="row-required">*</i>
            <bk-input class="row-content"
                data-vv-name="name"
                v-validate="'required|singlechar|length:20'"
                v-model.trim="templateName"
                :placeholder="$t('集群模板名称占位符')">
            </bk-input>
            <p class="row-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
        </div>
        <div class="layout-row template-row">
            <label class="row-label inline-block-middle" :title="$t('集群拓扑')">{{$t('集群拓扑')}}</label>
            <i class="row-required">*</i>
            <cmdb-set-template-tree class="row-content" ref="templateTree"></cmdb-set-template-tree>
        </div>
        <div class="template-options">
            <bk-button class="options-confirm" theme="primary" @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button class="options-cancel" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </section>
</template>

<script>
    import { MENU_BUSINESS_SET_TEMPLATE } from '@/dictionary/menu-symbol'
    import cmdbSetTemplateTree from './children/template-tree.vue'
    export default {
        components: {
            cmdbSetTemplateTree
        },
        data () {
            return {
                templateName: ''
            }
        },
        computed: {
            mode () {
                return this.$route.params.mode
            },
            templateId () {
                return this.$route.params.templateId
            }
        },
        created () {
            if (['edit', 'view'].includes(this.mode)) {
                this.getSetTemplateInfo()
            }
        },
        methods: {
            async getSetTemplateInfo () {
                try {
                    const info = await this.$store.dispatch('setTemplate/getSingleSetTemplateInfo', {
                        bizId: this.$store.getters['objectBiz/bizId'],
                        setTemplateId: this.templateId
                    })
                    this.templateName = info.name
                } catch (e) {
                    console.error(e)
                }
            },
            async handleConfirm () {
                try {
                    const validateResult = await this.$validator.validateAll()
                    if (!validateResult) {
                        return false
                    }
                    if (this.mode === 'create') {
                        await this.createSetTemplate()
                    } else {
                        await this.updateSetTemplate()
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
                this.$router.push({
                    name: MENU_BUSINESS_SET_TEMPLATE
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
        .row-label {
            width: 120px;
            text-align: right;
            font-size: 14px;
            line-height: 32px;
            padding: 0 0 0 10px;
            @include ellipsis;
        }
        .row-required {
            position: relative;
            top: 4px;
            font-style: normal;
            color: #EA3636;
        }
        .row-content {
            display: inline-block;
            vertical-align: top;
            width: 520px;
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
        padding:23px 0 0 137px;
        font-size: 0;
        .options-confirm {
            margin-right: 10px;
        }
    }
</style>
