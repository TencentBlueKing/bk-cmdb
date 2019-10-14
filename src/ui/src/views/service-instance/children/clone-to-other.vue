<template>
    <div class="create-layout clearfix" v-bkloading="{ isLoading: $loading() }">
        <label class="create-label fl">{{$t('添加主机')}}</label>
        <div class="create-hosts">
            <bk-button class="select-host-button" theme="default"
                @click="hostSelectorVisible = true">
                <i class="bk-icon icon-plus"></i>
                {{$t('添加主机')}}
            </bk-button>
            <div class="create-tables">
                <service-instance-table class="service-instance-table"
                    v-for="(host, index) in hosts"
                    ref="serviceInstanceTable"
                    deletable
                    expanded
                    :key="index"
                    :index="index"
                    :id="host.bk_host_id"
                    :name="host.bk_host_innerip"
                    :source-processes="sourceProcesses"
                    @delete-instance="handleDeleteInstance">
                </service-instance-table>
                <div class="buttons">
                    <span
                        v-cursor="{
                            active: !$isAuthorized($OPERATION.C_SERVICE_INSTANCE),
                            auth: [$OPERATION.C_SERVICE_INSTANCE]
                        }">
                        <bk-button theme="primary"
                            :disabled="!hosts.length || !$isAuthorized($OPERATION.C_SERVICE_INSTANCE)"
                            @click="handleConfirm">
                            {{$t('确定')}}
                        </bk-button>
                    </span>
                    <bk-button @click="handleBackToModule">{{$t('取消')}}</bk-button>
                </div>
            </div>
        </div>
        <host-selector
            :visible.sync="hostSelectorVisible"
            :module-instance="module"
            :exclude="[hostId]"
            @host-selected="handleSelectHost">
        </host-selector>
    </div>
</template>

<script>
    import hostSelector from '@/components/ui/selector/host.vue'
    import serviceInstanceTable from '@/components/service/instance-table.vue'
    import { MENU_BUSINESS_SERVICE_TOPOLOGY } from '@/dictionary/menu-symbol'
    export default {
        name: 'clone-to-other',
        components: {
            hostSelector,
            serviceInstanceTable
        },
        props: {
            module: {
                type: Object,
                default () {
                    return {}
                }
            },
            sourceProcesses: {
                type: Array,
                default () {
                    return {}
                }
            }
        },
        data () {
            return {
                hostSelectorVisible: false,
                hosts: []
            }
        },
        computed: {
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            hostId () {
                return parseInt(this.$route.params.hostId)
            },
            moduleId () {
                return parseInt(this.$route.params.moduleId)
            },
            setId () {
                return parseInt(this.$route.params.setId)
            }
        },
        methods: {
            handleSelectHost (checked, hosts) {
                this.hosts.push(...hosts)
                this.hostSelectorVisible = false
            },
            handleDeleteInstance (index) {
                this.hosts.splice(index, 1)
            },
            async handleConfirm () {
                try {
                    const serviceInstanceTables = this.$refs.serviceInstanceTable
                    await this.$store.dispatch('serviceInstance/createProcServiceInstanceWithRaw', {
                        params: this.$injectMetadata({
                            name: this.module.bk_module_name,
                            bk_module_id: this.moduleId,
                            instances: serviceInstanceTables.map(table => {
                                return {
                                    bk_host_id: table.id,
                                    processes: table.processList.map(item => {
                                        return {
                                            process_info: item
                                        }
                                    })
                                }
                            })
                        })
                    })
                    this.$success(this.$t('克隆成功'))
                    this.handleBackToModule()
                } catch (e) {
                    console.error(e)
                }
            },
            handleBackToModule () {
                this.$router.replace({
                    name: MENU_BUSINESS_SERVICE_TOPOLOGY,
                    query: {
                        module: this.moduleId
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-layout {
        margin: 35px 0 0 0;
        font-size: 14px;
        color: #63656E;
    }
    .create-label{
        display: block;
        width: 100px;
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
