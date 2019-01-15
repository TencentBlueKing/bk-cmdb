<template>
        <div class="host-topo-layout" >
            <h3 class="host-topo-header" @click="collapse = !collapse">
                <i class="bk-icon icon-angle-down" :class="{collapse}"></i>
                {{$t("BusinessTopology['业务拓扑']")}}
            </h3>
            <cmdb-collapse-transition>
                <div class="host-topo-list" v-show="!collapse">
                    <div class="host-topo-path" v-for="(topo, index) in hostTopo" :key="index">
                        {{topo}}
                    </div>
                </div>
            </cmdb-collapse-transition>
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
        data () {
            return {
                collapse: false
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
        padding: 28px 18px 0px 0px;
        .icon-angle-down {
            color: #333948;
            font-size: 12px;
            font-weight: bold;
            transition: transform .2s ease-in-out;
            &.collapse {
                transform: rotate(-90deg);
            }
        }
        .host-topo-header {
            display: inline-block;
            vertical-align: middle;
            padding: 0 0 10px 0;
            font-size: 14px;
            line-height: 14px;
            color: #333948;
            cursor: pointer;
        }
        .host-topo-list {
            padding-bottom: 7px;
        }
        .host-topo-path {
            padding: 0 0 0 18px;
            font-size: 12px;
            line-height: 20px;
        }
    }
</style>