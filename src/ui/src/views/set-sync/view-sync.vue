<template>
    <div class="view-sync-layout">
        <div class="header">
            <i :class="['bk-icon', ...iconClass]"></i>
            <div class="right-content">
                <div class="operations">
                    <span>{{$t('正在同步中，请耐心等待…')}}</span>
                    <bk-button v-if="syncStatus === 2" theme="default">{{$t('返回列表')}}</bk-button>
                    <template v-else>
                        <bk-button theme="default">{{$t('失败重试')}}</bk-button>
                        <bk-button theme="default">{{$t('关闭页面')}}</bk-button>
                    </template>
                </div>
                <div class="sync-count">
                    <p v-if="syncStatus === 2">{{$t('已同步10处，剩余5处')}}</p>
                    <i18n v-else
                        path="同步数量提示"
                        tag="div">
                        <span class="success" place="successCount" @click="handleFilterSuccess">0</span>
                        <span class="fail" place="failCount" @click="handleFilterFail">50 </span>
                    </i18n>
                </div>
            </div>
        </div>
        <bk-table class="sync-table"
            :height="tableHeight"
            :data="syncList">
            <bk-table-column prop="topo_path" :label="$t('拓扑结构')"></bk-table-column>
            <bk-table-column prop="status" :label="$t('状态')" width="400">
                <template slot-scope="{ row }">
                    <span v-if="row.status === 0" class="sync-status fail-status">
                        <i class="bk-icon icon-cc-log-02"></i>
                        {{$t('同步失败')}}
                    </span>
                    <span v-else-if="row.status === 1" style="color: #2DCB56;">{{$t('已完成')}}</span>
                    <span v-else-if="row.status === 2" class="sync-status syncing-status">
                        <img src="../../assets/images/icon/loading.svg" alt="">
                        {{$t('同步中')}}
                    </span>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
            <infinite-loading slot="append" v-if="ready"
                force-use-infinite-wrapper=".sync-table .bk-table-body-wrapper"
                :distance="42"
                :identifier="infiniteIdentifier"
                @infinite="infiniteHandler">
                <span slot="no-more"></span>
                <span slot="no-results"></span>
                <span slot="error"></span>
                <div slot="spinner" style="height: 42px;"
                    v-bkloading="{
                        isLoading: $loading()
                    }">
                </div>
            </infinite-loading>
        </bk-table>
    </div>
</template>

<script>
    import infiniteLoading from 'vue-infinite-loading'
    export default {
        components: {
            infiniteLoading
        },
        data () {
            return {
                ready: false,
                infiniteIdentifier: 0,
                syncStatus: 1,
                syncList: [],
                pagination: {
                    current: 1,
                    count: 0,
                    size: 10
                }
            }
        },
        computed: {
            iconClass () {
                let classMap = []
                switch (this.syncStatus) {
                    case 0:
                        classMap = ['sync-fail', 'icon-close']
                        break
                    case 1:
                        classMap = ['sync-success', 'icon-check-1']
                        break
                    case 2:
                        classMap = ['syncing', 'icon-cc-syncing']
                        break
                    default:
                        classMap = []
                }
                return classMap
            },
            tableHeight () {
                return this.$APP.height - 206
            }
        },
        created () {
            this.getSyncList()
        },
        methods: {
            async getSyncList () {
                const delay = (duration) => new Promise(resolve => setTimeout(() => {
                    resolve()
                }, duration))
                await delay(1000)
                const arr = []
                for (let i = 0; i < 10; i++) {
                    arr.push(...[
                        {
                            topo_path: '广东省厅 / 深圳区 / 正式环境集群',
                            status: 0
                        },
                        {
                            topo_path: '广东省厅 / 深圳区 / 正式环境集群',
                            status: 1
                        },
                        {
                            topo_path: '广东省厅 / 深圳区 / 正式环境集群',
                            status: 2
                        },
                        {
                            topo_path: '广东省厅 / 深圳区 / 正式环境集群',
                            status: -1
                        }
                    ])
                }
                this.syncList.push(...arr)
                this.ready = true
            },
            async infiniteHandler (infiniteState) {
                try {
                    await this.getSyncList()
                    infiniteState.loaded()
                    // if (!data) {
                    //     infiniteState.complete()
                    // }
                } catch (e) {
                    infiniteState.error()
                }
            },
            handleFilterFail () {

            },
            handleFilterSuccess () {

            }
        }
    }
</script>

<style lang="scss" scoped>
    .view-sync-layout {
        padding: 0 50px;
    }
    .header {
        display: flex;
        align-items: center;
        margin-bottom: 20px;
        .bk-icon {
            width: 60px;
            height: 60px;
            line-height: 60px;
            text-align: center;
            font-size: 20px;
            font-weight: bold;
            color: #ffffff;
            border-radius: 50%;
            margin-right: 12px;
            &.syncing {
                background-color: #3A84FF;
            }
            &.sync-fail {
                background-color: #FF5656;
            }
            &.sync-success {
                background-color: #51CE71;
            }
            &.icon-cc-syncing {
                font-size: 32px;
                font-weight: normal;
            }
        }
        .operations {
            display: flex;
            align-items: center;
            font-size: 20px;
            color: #313237;
            padding-bottom: 4px;
            .bk-button {
                width: 86px;
                height: 26px;
                line-height: 24px;
                padding: 0;
                margin-left: 10px;
            }
        }
        .sync-count {
            font-size: 12px;
            color: #979BA5;
            span {
                font-weight: bold;
                cursor: pointer;
            }
            .success {
                color: #51CE71;
            }
            .fail {
                color: #FF5656;
            }
        }
    }
    .sync-status {
        display: flex;
        align-items: center;
        &.fail-status {
            color: #FF5656;
            .bk-icon {
                margin-right: 4px;
                color: #63656E;
            }
        }
        &.syncing-status {
            img {
                width: 18px;
                margin-right: 4px;
            }
        }
    }
</style>
