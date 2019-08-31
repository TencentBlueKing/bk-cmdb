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
        computed: {
            ...mapGetters('objectBiz', ['authorizedBusiness']),
            ...mapGetters(['isAdminView'])
        },
        methods: {
            handleHostClick () {
                const name = this.isAdminView ? 'resourceHostDetails' : 'businessHostDetails'
                this.$router.push({
                    name,
                    params: {
                        business: this.isAdminView ? '' : this.host['biz'][0]['bk_biz_id'],
                        id: this.host['host']['bk_host_id']
                    },
                    query: {
                        from: this.$route.fullPath
                    }
                })
            },
            checkoutBizAuth (bizId) {
                return this.authorizedBusiness.some(biz => biz['bk_biz_id'] === bizId)
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
