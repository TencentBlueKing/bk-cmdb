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
            ...mapGetters(['bkPrivBizList'])
        },
        methods: {
            handleHostClick () {
                const path = this.host['biz'][0]['default'] === 1 ? this.pathMap.resource : this.pathMap.hosts
                if (path === this.pathMap.hosts) {
                    const bizId = this.host['biz'][0]['bk_biz_id']
                    if (this.checkoutBizAuth(bizId)) {
                        this.$store.commit('setBkBizId', bizId)
                    } else {
                        this.$alertMsg(this.$t('Hosts["权限不足"]'))
                        return
                    }
                }
                this.$store.commit('setHostSearch', {
                    ip: this.host['host']['bk_host_innerip'],
                    outerip: false,
                    exact: 1
                })
                this.$router.push(path)
            },
            checkoutBizAuth (bizId) {
                return this.bkPrivBizList.some(biz => biz['bk_biz_id'] === bizId)
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