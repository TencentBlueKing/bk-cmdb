<template>
    <div class="tips-wrapper">
        <div class="content-wrapper">
            <div class="title">
                <img src="../../assets/images/no-authority.png" alt="no-authority">
                <h2>{{$t("Common['无业务权限']")}}</h2>
                <p>{{$t("Common['点击下方按钮申请']")}}</p>
            </div>
            <div class="btns">
                <bk-button theme="primary" @click="handleApplyPermission">
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
    .content-wrapper {
        margin-top: 100px;
        text-align: center;
        color: #63656E;
        font-size: 14px;
        .title {
            img {
                width: 136px;
            }
            h2 {
                margin-bottom: 10px;
                font-size: 22px;
                color: #313238;
                font-weight: normal;
            }
            p {
                color: #63656e;
                font-size: 14px;
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
