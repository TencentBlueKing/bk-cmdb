<template>
    <bk-popover class="process-form-tips">
        <i class="icon icon-cc-lock"></i>
        <template slot="content">
            <span v-if="isLimited">{{$t('系统限定不可修改')}}</span>
            <i18n v-else path="进程表单锁定提示">
                <bk-link theme="primary" @click="handleRedirect" place="link">{{$t('跳转服务模板')}}</bk-link>
            </i18n>
        </template>
    </bk-popover>
</template>

<script>
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
        methods: {
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
        display: inline-block;
        width: 16px;
        height: 16px;
        line-height: 16px;
        text-align: center;
        vertical-align: middle;
        color: #c3cdd7;
        /deep/ {
            .bk-tooltip-ref {
                font-size: 0;
            }
        }
        .icon {
            font-size: 16px;
        }
    }
</style>
