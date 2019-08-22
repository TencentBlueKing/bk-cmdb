<template>
    <div class="tips-wrapper">
        <div class="content-wrapper">
            <div class="title">
                <img src="../../assets/images/no-authority.png" alt="no-authority">
                <h2>{{$t('无业务权限')}}</h2>
                <p>{{$t('点击下方按钮申请')}}</p>
            </div>
            <div class="btns">
                <bk-button theme="primary" @click="handleApplyPermission" :loading="$loading('getSkipUrl')">
                    {{$t('申请业务权限')}}
                </bk-button>
            </div>
        </div>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
    export default {
        computed: {
            ...mapGetters(['permission'])
        },
        created () {
            this.$store.commit('setHeaderTitle', '')
        },
        beforeRouteEnter (to, from, next) {
            if (from.fullPath === '/') {
                next({ name: 'index' })
            } else {
                next()
            }
        },
        methods: {
            async handleApplyPermission () {
                try {
                    const skipUrl = await this.$store.dispatch('auth/getSkipUrl', {
                        params: this.permission,
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
