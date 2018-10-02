<template>
    <div class="host-topo-layout">
        <h3 class="host-topo-header">{{$t("BusinessTopology['业务拓扑']")}}</h3>
        <div class="host-topo-path" v-for="(topo, index) in hostTopo" :key="index">{{topo}}</div>
    </div>
</template>

<script>
    export default {
        props: {
            host: {
                type: Object,
                required: true
            }
        },
        computed: {
            hostTopo () {
                const hostTopo = []
                this.host.module.map(module => {
                    const set = this.host.set.find(set => {
                        return set['bk_set_id'] === module['bk_set_id']
                    })
                    if (set) {
                        let biz = this.host.biz.find(biz => {
                            return biz['bk_biz_id'] === set['bk_biz_id']
                        })
                        if (biz) {
                            hostTopo.push(`${biz['bk_biz_name']} > ${set['bk_set_name']} > ${module['bk_module_name']}`)
                        }
                    }
                })
                return hostTopo
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-topo-layout {
        padding: 28px 18px 7px;
        .host-topo-header {
            padding: 0 0 10px 0;
            font-size: 14px;
            line-height: 14px;
            color: #333948;
        }
        .host-topo-path {
            font-size: 12px;
            line-height: 20px;
        }
    }
</style>