<template>
    <div class="tips-wrapper">
        <div class="content-wrapper">
            <bk-exception type="403">
                <div class="title">
                    <h2>{{$t('无功能权限')}}</h2>
                    <p>{{$t('点击下方按钮申请')}}</p>
                </div>
                <div class="btns">
                    <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
                        {{$t('申请功能权限')}}
                    </bk-button>
                </div>
            </bk-exception>
        </div>
    </div>
</template>
<script>
    import { translateAuth } from '@/setup/permission'
    export default {
        methods: {
            async handleApplyPermission () {
                try {
                    const { view, permission } = this.$route.meta.auth || {}
                    const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
                        params: view ? translateAuth(view) : permission,
                        config: {
                            requestId: 'getSkipUrl'
                        }
                    })
                    window.open(skipUrl)
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tips-wrapper {
        overflow: hidden;
    }
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
