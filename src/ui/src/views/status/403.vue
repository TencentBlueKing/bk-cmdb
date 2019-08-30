<template>
    <div class="permisson-apply">
        <div class="apply-content">
            <h3>{{i18n.resourceTitle}}</h3>
            <p>{{i18n.resourceContent}}</p>
            <div class="operation-btns">
                <bk-button theme="primary" @click="handleApplyPermission">{{i18n.apply}}</bk-button>
            </div>
        </div>
    </div>
</template>
<script>
    export default {
        name: 'PermissionApply',
        data () {
            return {
                i18n: {
                    resourceTitle: this.$t('无权限访问'),
                    resourceContent: this.$t('你没有相应资源的访问权限，请申请权限或联系管理员授权'),
                    apply: this.$t('去申请')
                }
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', '')
        },
        methods: {
            handleApplyPermission () {
                const topWindow = window.top
                const isPaasConsole = topWindow !== window
                const authCenter = window.Site.authCenter || {}
                if (isPaasConsole) {
                    topWindow.postMessage(JSON.stringify({
                        action: 'open_other_app',
                        app_code: authCenter.appCode,
                        app_url: 'apply-by-system'
                    }), '*')
                } else {
                    window.open(authCenter.url)
                }
            }
        }
    }
</script>
<style lang="scss" scoped>
    .apply-content {
        margin-top: 240px;
        text-align: center;
        & > h3 {
            margin: 0 0 30px;
            color: #313238;
            font-size: 20px;
        }
        & > p {
            margin: 0 0 30px;
            color: #979ba5;
            font-size: 14px;
        }
        .bk-button {
            height: 32px;
            line-height: 30px;
        }
    }
</style>
