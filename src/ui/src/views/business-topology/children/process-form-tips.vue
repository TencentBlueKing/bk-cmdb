<template>
    <bk-popover class="process-form-tips">
        <i class="icon icon-cc-lock"></i>
        <i18n path="进程表单锁定提示" slot="content">
            <bk-link theme="primary" @click="handleRedirect" place="link">{{$t('跳转服务模板')}}</bk-link>
        </i18n>
    </bk-popover>
</template>

<script>
    export default {
        props: {
            serviceTemplateId: Number
        },
        methods: {
            handleRedirect () {
                this.$router.push({
                    name: 'operationalTemplate',
                    params: {
                        templateId: this.serviceTemplateId
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .process-form-tips {
        font-size: 16px;
        vertical-align: middle;
        color: #c3cdd7;
        margin-left: 6px;
    }
</style>
