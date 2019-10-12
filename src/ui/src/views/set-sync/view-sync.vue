<template>
    <div class="view-sync-layout">
        <div class="header">
            <i :class="['bk-icon', ...iconClass]"></i>
            <div class="right-content">
                <div class="operations">
                    <span>{{$t('正在同步中请等待')}}</span>
                    <!-- <span>{{$t('同步失败')}}</span>
                    <span>{{$t('同步完成')}}</span> -->
                    <bk-button v-if="syncStatus === 2" theme="default">{{$t('返回列表')}}</bk-button>
                    <template v-else>
                        <bk-button theme="default">{{$t('失败重试')}}</bk-button>
                        <bk-button theme="default" @click="handleClose">{{$t('关闭页面')}}</bk-button>
                    </template>
                </div>
                <div class="sync-count">
                    <i18n v-if="syncStatus === 2"
                        path="同步进度提示"
                        tag="div">
                        <span place="completedCount">10</span>
                        <span place="remainingCount">5</span>
                    </i18n>
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
                    <span v-if="[200].includes(row.status)" style="color: #2DCB56;">{{$t('已完成')}}</span>
                    <span v-else-if="[0, 1, 100].includes(row.status)" class="sync-status syncing-status">
                        <img src="../../assets/images/icon/loading.svg" alt="">
                        {{$t('同步中')}}
                    </span>
                    <span v-else class="sync-status fail-status">
                        <i class="bk-icon icon-cc-log-02"></i>
                        {{$t('同步失败')}}
                    </span>
                </template>
            </bk-table-column>
        </bk-table>
    </div>
</template>

<script>
    export default {
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
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            setTemplateId () {
                return this.$route.params['setTemplateId']
            },
            setInstancesId () {
                const id = `${this.business}_${this.setTemplateId}`
                let syncIdMap = this.$store.state.setFeatures.syncIdMap
                const sessionSyncIdMap = sessionStorage.getItem('setSyncIdMap')
                if (!Object.keys(syncIdMap).length && sessionSyncIdMap) {
                    syncIdMap = JSON.parse(sessionSyncIdMap)
                    this.$store.commit('setFeatures/resetSyncIdMap', syncIdMap)
                }
                return syncIdMap[id] || []
            },
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
        async created () {
            await this.getSyncStatus()
            // this.polling()
        },
        methods: {
            polling () {
                let timer = null
                try {
                    timer = setInterval(async () => {
                        await this.getSyncStatus()
                    }, 5000)
                } catch (e) {
                    console.error(e)
                    clearInterval(timer)
                }
            },
            async getSyncStatus () {
                try {
                    const data = await this.$store.dispatch('setSync/getInstancesSyncStatus', {
                        bizId: this.business,
                        setTemplateId: this.setTemplateId,
                        params: {
                            bk_set_ids: this.setInstancesId
                        },
                        config: {
                            requestId: 'getInstancesSyncStatus',
                            cancelPrevious: true
                        }
                    })
                    const moduleList = []
                    if (data) {
                        const sets = Object.keys(data).map(key => data[key]).filter(set => set)
                        sets.forEach(set => {
                            const modules = (set.detail || []).filter(module => {
                                const moduleData = module.data
                                return moduleData && moduleData.module_diff && moduleData.module_diff.diff_type !== 'unchanged'
                            })
                            moduleList.push(...modules)
                        })
                    }
                    this.syncList = moduleList
                } catch (e) {
                    console.error(e)
                    this.syncList = []
                }
            },
            handleFilterFail () {

            },
            handleFilterSuccess () {

            },
            handleClose () {
                this.$router.replace({
                    name: 'setTemplateInfo',
                    params: {
                        templateId: this.setTemplateId
                    }
                })
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
