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
                    <p v-else>
                        {{$t('已经成功同步实例模块')}}
                        <span class="success" @click="handleFilterSuccess">0 </span>{{$t('个')}},
                        {{$t('失败')}}
                        <span class="fail" @click="handleFilterFail">50 </span>{{$t('个')}}
                    </p>
                </div>
            </div>
        </div>
        <bk-table :data="syncList">
            <bk-table-column prop="topo_path" :label="$t('拓扑结构')"></bk-table-column>
            <bk-table-column prop="status" :label="$t('状态')" width="400">
                <template slot-scope="{ row }">
                    <span v-if="row.status === 0" style="color: #FF5656;">{{$t('同步失败')}}</span>
                    <span v-else-if="row.status === 1" style="color: #2DCB56;">{{$t('已完成')}}</span>
                    <span v-else-if="row.status === 2">{{$t('同步中')}}</span>
                    <span v-else>--</span>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                syncStatus: 1,
                syncList: [
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
                ]
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
                        classMap = ['syncing']
                        break
                    default:
                        classMap = []
                }
                return classMap
            }
        },
        methods: {
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
            font-size: 28px;
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
</style>
