<template>
    <div class="landing-layout" v-bkloading="{ isLoading: loading }">
        <bk-exception type="404" v-if="notFound">
            <span>{{$t('未查询到主机')}}</span>
        </bk-exception>
        <div class="bk-exception" v-else-if="error">
            <img src="../../assets/images/error.png">
            <p class="exception-text">{{$t('查询主机时发生异常')}}</p>
        </div>
    </div>
</template>

<script>
    import { MENU_BUSINESS_HOST_DETAILS, MENU_RESOURCE_HOST, MENU_RESOURCE_HOST_DETAILS } from '@/dictionary/menu-symbol'
    const BK_NO_LIMIT = 999999999
    export default {
        data () {
            return {
                loading: true,
                notFound: false,
                error: false,
                requestId: Symbol('searchHost')
            }
        },
        computed: {
            params () {
                const bizCondition = { bk_obj_id: 'biz', condition: [], fields: ['bk_biz_id', 'default'] }
                const setCondition = { bk_obj_id: 'set', condition: [], fields: ['bk_set_id'] }
                const moduleCondition = { bk_obj_id: 'module', condition: [], fields: ['bk_module_id'] }
                const hostCondition = { bk_obj_id: 'host', condition: [], fields: ['bk_host_id'] }
                const params = {
                    bk_biz_id: -1,
                    condition: [bizCondition, setCondition, moduleCondition, hostCondition],
                    ip: {
                        data: [this.$route.params.ip],
                        exact: 1,
                        flag: 'bk_host_innerip'
                    },
                    page: {
                        limit: BK_NO_LIMIT,
                        start: 0
                    }
                }
                const cloudId = parseInt(this.$route.params.cloudId)
                if (!isNaN(cloudId)) {
                    hostCondition.condition.push({
                        field: 'bk_cloud_id',
                        operator: '$eq',
                        value: cloudId
                    })
                }
                return params
            }
        },
        created () {
            this.searchHost()
        },
        methods: {
            async searchHost () {
                try {
                    const { info } = await this.$store.dispatch('hostSearch/searchHost', {
                        params: this.params,
                        config: {
                            requestId: this.requestId
                        }
                    })
                    if (!info.length) {
                        this.notFound = true
                        this.loading = false
                    } else if (info.length === 1) {
                        const [data] = info
                        this.redirectToDetails(data)
                    } else {
                        this.redirectToResource()
                    }
                } catch (error) {
                    this.error = true
                    this.loading = false
                    console.error(error)
                }
            },
            redirectToDetails (data) {
                const { host, biz } = data
                if (biz[0].default === 1) {
                    this.$routerActions.redirect({
                        name: MENU_RESOURCE_HOST_DETAILS,
                        params: {
                            id: host.bk_host_id
                        }
                    })
                } else {
                    this.$routerActions.redirect({
                        name: MENU_BUSINESS_HOST_DETAILS,
                        params: {
                            bizId: biz[0].bk_biz_id,
                            id: host.bk_host_id
                        }
                    })
                }
            },
            redirectToResource () {
                this.$routerActions.redirect({
                    name: MENU_RESOURCE_HOST,
                    query: {
                        ip: this.$route.params.ip,
                        scope: 'all'
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .landing-layout {
        display: flex;
        height: 100%;
    }
</style>
