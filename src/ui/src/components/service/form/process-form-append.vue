<template>
    <span v-if="isLimited"></span>
    <bk-popover class="process-form-tips" v-else>
        <i class="icon-cc-lock-fill"></i>
        <template slot="content">
            <i18n path="进程表单锁定提示">
                <bk-link theme="primary" @click="handleRedirect" place="link">{{$t('跳转服务模板')}}</bk-link>
            </i18n>
        </template>
    </bk-popover>
</template>

<script>
    import Tippy from 'bk-magic-vue/lib/utils/tippy'
    export default {
        props: {
            serviceTemplateId: Number,
            bizId: Number,
            property: {
                type: Object,
                default: () => ({})
            }
        },
        computed: {
            isLimited () {
                return ['bk_func_name', 'bk_process_name'].includes(this.property.bk_property_id)
            }
        },
        mounted () {
            if (this.isLimited) {
                this.setupTips()
            } else {
                this.hackRadius()
            }
        },
        methods: {
            setupTips () {
                const DOM = this.$el.previousElementSibling
                Tippy(DOM, {
                    content: this.$t('系统限定不可修改'),
                    arrow: true,
                    placement: 'top'
                })
            },
            hackRadius () {
                const hackDOM = this.$el.parentElement.querySelectorAll('.bk-form-input,.bk-form-textarea,.bk-textarea-wrapper,.bk-select')
                Array.prototype.forEach.call(hackDOM, dom => {
                    dom.style.borderTopRightRadius = 0
                    dom.style.borderBottomRightRadius = 0
                })
            },
            handleRedirect () {
                this.$routerActions.redirect({
                    name: 'operationalTemplate',
                    params: {
                        bizId: this.bizId,
                        templateId: this.serviceTemplateId
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .process-form-tips {
        width: 24px;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #dcdee5;
        border-left: none;
        background-color: #fafbfd;
        font-size: 14px;
        overflow: hidden;
        /deep/ .bk-tooltip-ref {
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
        }
    }
</style>
