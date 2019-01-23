<template>
    <div class="search-item-host clearfix" @click="handleHostClick" :title="getHostTitle(host)">
        <div class="host-ip fl">{{host['host']['bk_host_innerip']}}</div>
        <div class="host-biz fr">{{host['biz'][0]['bk_biz_name']}}</div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            host: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                pathMap: {
                    resource: '/resource',
                    hosts: '/hosts'
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['privilegeBusiness'])
        },
        methods: {
            handleHostClick () {
                const resourceName = 'resource'
                const hostName = 'hosts'
                const name = this.host['biz'][0]['default'] === 1 ? resourceName : hostName
                const bizId = this.host['biz'][0]['bk_biz_id']
                if (name === hostName) {
                    if (!this.checkoutBizAuth(bizId)) {
                        this.$error(this.$t('Hosts["权限不足"]'))
                        return
                    }
                }
                this.$router.push({
                    name,
                    query: {
                        business: bizId,
                        ip: this.host['host']['bk_host_innerip'],
                        outer: false,
                        inner: true,
                        exact: 1
                    }
                })
            },
            checkoutBizAuth (bizId) {
                return this.privilegeBusiness.some(biz => biz['bk_biz_id'] === bizId)
            },
            getHostTitle (host) {
                return `${host['host']['bk_host_innerip']}—${host['biz'][0]['bk_biz_name']}`
            }
        }
    }
</script>

<style lang="scss" scoped>
.search-item-host{
    padding: 0 4px;
}
.search-item-host:hover{
    .host-ip{
        color: #3c96ff;
    }
}
.host-ip {
    width: 60%;
    @include ellipsis;
}
.host-biz{
    width: 35%;
    text-align: right;
    color: #c3cdd7;
    @include ellipsis;
}
</style>