<template>
    <div class="create-layout clearfix">
        <label class="create-label fl">{{$t('businessTopology["添加主机"]')}}</label>
        <div class="create-hosts">
            <bk-button class="select-host-button" type="default"
                @click="hostSelectorVisible = true">
                <i class="bk-icon icon-plus"></i>
                {{$t('businessTopology["添加主机"]')}}
            </bk-button>
            <div class="create-tables">
                <service-instance-table class="service-instance-table"
                    v-for="(host, index) in hosts"
                    deletable
                    expanded
                    :key="index"
                    :index="index"
                    :id="host.bk_host_id"
                    :name="host.bk_host_innerip"
                    @delete-instance="handleDeleteInstance">
                </service-instance-table>
                <div class="buttons">
                    <bk-button type="primary"
                        :disabled="!hosts.length">
                        {{$t('Common["确定"]')}}
                    </bk-button>
                    <bk-button>{{$t('Common["取消"]')}}</bk-button>
                </div>
            </div>
        </div>
        <host-selector
            :visible.sync="hostSelectorVisible"
            :module-instance="moduleInstance"
            @host-selected="handleSelectHost">
        </host-selector>
    </div>
</template>

<script>
    import hostSelector from '@/components/ui/selector/host.vue'
    import serviceInstanceTable from '@/components/service/instance-table.vue'
    export default {
        components: {
            hostSelector,
            serviceInstanceTable
        },
        data () {
            return {
                hostSelectorVisible: false,
                moduleInstance: {},
                hosts: []
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            moduleId () {
                return parseInt(this.$route.params.moduleId)
            },
            setId () {
                return parseInt(this.$route.params.setId)
            }
        },
        created () {
            this.getModuleInstance()
        },
        methods: {
            async getModuleInstance () {
                try {
                    const data = await this.$store.dispatch('objectModule/searchModule', {
                        bizId: this.business,
                        setId: this.setId,
                        params: {
                            page: { start: 0, limit: 1 },
                            fields: [],
                            condition: {
                                bk_module_id: this.moduleId,
                                bk_supplier_account: this.$store.getters.supplierAccount
                            }
                        },
                        config: {
                            requestId: 'getModuleInstance'
                        }
                    })
                    if (!data.count) {
                        this.$router.replace({ name: '404' })
                    } else {
                        this.moduleInstance = data.info[0]
                    }
                } catch (e) {
                    console.error(e)
                }
            },
            handleSelectHost (checked, hosts) {
                this.hosts.push(...hosts)
                this.hostSelectorVisible = false
            },
            handleDeleteInstance (index) {
                this.hosts.splice(index, 1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-layout {
        height: 100%;
        padding: 32px 23px 0;
        font-size: 14px;
        color: #63656E;
        background-color: #FAFBFD;
    }
    .create-label{
        display: block;
        position: relative;
        line-height: 32px;
        &:after {
            content: "*";
            margin: 0 0 0 4px;
            color: $cmdbDangerColor;
            @include inlineBlock;
        }
    }
    .create-hosts {
        padding-left: 10px;
        height: 100%;
        overflow: hidden;
    }
    .select-host-button {
        height: 32px;
        line-height: 30px;
        font-size: 0;
        .bk-icon {
            position: static;
            height: 30px;
            line-height: 30px;
            font-size: 12px;
            font-weight: bold;
            @include inlineBlock(top);
        }
        /deep/ span {
            font-size: 14px;
        }
    }
    .create-tables {
        height: calc(100% - 54px);
        margin: 22px 0 0 0;
        @include scrollbar-y;
        .buttons {
            padding: 18px 0 0 0;
        }
    }
    .service-instance-table {
        margin-bottom: 12px;
    }
</style>
