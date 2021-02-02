<template>
    <div class="tips-wrapper">
        <div class="content-wrapper">
            <bk-exception type="403">
                <div class="title">
                    <h2>{{$t('业务不存在或无权限')}}</h2>
                </div>
                <div class="btns">
                    <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
                        {{$t('申请业务访问权限')}}
                    </bk-button>
                    <bk-button theme="primary" @click="handleCreate">
                        {{$t('创建业务')}}
                    </bk-button>
                </div>
            </bk-exception>
        </div>
    </div>
</template>
<script>
    import { translateAuth } from '@/setup/permission'
    import { MENU_RESOURCE_BUSINESS } from '@/dictionary/menu-symbol'
    export default {
        computed: {
            bizId () {
                return this.$route.params.bizId
            }
        },
        methods: {
            async handleApplyPermission () {
                try {
                    const permission = translateAuth({
                        type: this.$OPERATION.R_BIZ_RESOURCE,
                        relation: this.bizId ? [this.bizId] : []
                    })
                    const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
                        params: permission,
                        config: {
                            requestId: 'getSkipUrl'
                        }
                    })
                    window.open(skipUrl)
                } catch (e) {
                    console.error(e)
                }
            },
            handleCreate () {
                this.$routerActions.redirect({ name: MENU_RESOURCE_BUSINESS })
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
