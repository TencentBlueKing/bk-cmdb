<template>
    <div class="tips-wrapper">
        <div class="content-wrapper">
            <div class="title">
                <i class="icon icon-cc-no-authority"></i>
                <h2>{{$t("Common['无业务权限']")}}</h2>
            </div>
            <div class="btns">
                <bk-button type="primary" @click="handleApplyPermission">
                    {{$t("Common['申请业务权限']")}}
                </bk-button>
            </div>
        </div>
    </div>
</template>
<script>
    export default {
        created () {
            this.$store.commit('setHeaderTitle', '')
        },
        methods: {
            handleApplyPermission () {
                const topWindow = window.top
                const isPaasConsole = topWindow !== window && topWindow.BLUEKING
                const authCenter = window.Site.authCenter || {}
                if (isPaasConsole) {
                    topWindow.postMessage(JSON.stringify({
                        action: 'open_other_app',
                        app_code: authCenter.appCode,
                        app_url: 'perm-apply'
                    }), '*')
                } else {
                    window.open(authCenter.url)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .content-wrapper {
        margin-top: 100px;
        text-align: center;
        color: #63656E;
        font-size: 14px;
        .title {
            .icon {
                font-size: 56px;
                color: #979BA5;
            }
            h2 {
                margin-top: 10px;
                margin-bottom: 10px;
                font-size: 22px;
                color: #313238;
                font-weight: normal;
            }
        }
        .btns {
            margin-top: 24px;
            .bk-button {
                border-radius: 3px;
                padding-left: 10px;
                padding-right: 10px;
            }
        }
    }
</style>
