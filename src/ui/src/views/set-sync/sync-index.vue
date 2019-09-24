<template>
    <div class="sync-set-wrapper" v-bkloading="{ isLoading: $loading('diffTemplateAndInstances') }">
        <template v-if="noInfo">
            <div class="no-content">
                <img src="../../assets/images/no-content.png" alt="no-content">
                <p>{{$t('无集群模板更新信息')}}</p>
                <bk-button theme="primary" @click="handleGoback">{{$t('返回')}}</bk-button>
            </div>
        </template>
        <template v-else-if="isLatestInfo">
            <div class="no-content">
                <img src="../../assets/images/latest-data.png" alt="no-content">
                <p>{{$t('最新集群模板信息')}}</p>
                <bk-button theme="primary" @click="handleGoback">{{$t('返回')}}</bk-button>
            </div>
        </template>
        <template v-else-if="diffList.length">
            <div class="title clearfix">
                <div class="tips fl">
                    <p class="mr10">{{$t('请确认以下模板修改信息')}}：</p>
                    <span class="mr30">
                        <i class="dot"></i>
                        {{$t('新增模块')}}
                    </span>
                    <span class="mr30">
                        <i class="dot blue"></i>
                        {{$t('变更模块')}}
                    </span>
                    <span>
                        <i class="dot red"></i>
                        {{$t('删除模块')}}
                    </span>
                </div>
                <bk-checkbox class="expand-all fr" v-model="expandAll" @change="handleExpandAll">{{$t('全部展开')}}</bk-checkbox>
            </div>
            <div class="instance-list">
                <set-instance class="instance-item"
                    ref="setInstance"
                    v-for="(instance, index) in diffList"
                    :key="instance.bk_set_id"
                    :instance="instance"
                    :expand="index === 0">
                </set-instance>
            </div>
            <div class="footer">
                <bk-button theme="primary" class="mr10" @click="handleConfirmSync">{{$t('确认同步')}}</bk-button>
                <bk-button class="mr10">{{$t('取消')}}</bk-button>
                <span>{{$tc('已选集群实例', 20, { count: 20 })}}</span>
            </div>
        </template>
    </div>
</template>

<script>
    import setInstance from './set-instance'
    export default {
        components: {
            setInstance
        },
        data () {
            return {
                expandAll: false,
                diffList: [],
                noInfo: false,
                isLatestInfo: false
            }
        },
        computed: {
            bizId () {
                return this.$store.getters['objectBiz/bizId']
            }
        },
        async created () {
            await this.getDiffData()
        },
        methods: {
            async getDiffData () {
                try {
                    this.diffList = await this.$store.dispatch('setSync/diffTemplateAndInstances', {
                        bizId: 3 || this.bizId,
                        setTemplateId: 1,
                        params: {
                            bk_set_ids: [14]
                        },
                        config: {
                            requestId: 'diffTemplateAndInstances'
                        }
                    })
                    this.isLatestInfo = false
                    this.noInfo = false
                } catch (e) {
                    console.error(e)
                    this.noInfo = true
                }
            },
            handleConfirmSync () {
                this.$bkInfo({
                    type: 'warning',
                    title: '确定同步拓扑？',
                    // subTitle: '即将同步拓扑模版【正式集群模版】，模版的拓扑结构将会更新到选中的集群实例中',
                    subTitle: '即将批量同步拓扑模版【正式集群模版】到选中的20个集群实例中，模版的拓扑结构将会更新到选中的集群实例中',
                    extCls: 'set-confirm-sync',
                    confirmFn: async () => {
                        try {
                            this.$router.push({
                                name: 'viewSync'
                            })
                            await this.$store.dispatch('setSync/syncTemplateToInstances', {
                                bizId: this.bizId,
                                setTemplateId: '',
                                params: {
                                    bk_set_ids: []
                                },
                                config: {
                                    requestId: 'syncTemplateToInstances'
                                }
                            })
                            this.$success(this.$t('同步成功'))
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            },
            handleExpandAll (expand) {
                this.$refs.setInstance.forEach(instance => {
                    instance.localExpand = expand
                })
            },
            handleGoback () {

            }
        }
    }
</script>

<style lang="scss" scoped>
    .sync-set-wrapper {
        padding: 0 20px;
    }
    .no-content {
        position: absolute;
        top: 50%;
        left: 50%;
        font-size: 16px;
        color: #63656e;
        text-align: center;
        transform: translate(-50%, -50%);
        img {
            width: 130px;
        }
        p {
            padding: 20px 0 30px;
        }
    }
    .tips {
        display: flex;
        align-items: center;
        font-size: 14px;
        color: #63656E;
        .dot {
            display: inline-block;
            width: 10px;
            height: 10px;
            border-radius: 50%;
            background-color: #2DCB56;
            margin-right: 2px;
            &.red {
                background-color: #FF5656;
            }
            &.blue {
                background-color: #3A84FF;
            }
        }
    }
    .expand-all {
        color: #888991;
    }
    .instance-list {
        padding: 20px 0 10px;
        .instance-item {
            margin-bottom: 10px;
        }
    }
    .footer {
        display: flex;
        align-items: center;
        > span {
            color: #979BA5;
            font-size: 14px;
        }
    }
</style>

<style lang="scss">
    .set-confirm-sync {
        .bk-dialog-content {
            width: 440px !important;
        }
        .bk-dialog-type-body {
            padding: 2px 24px 5px !important;
        }
        .bk-dialog-type-sub-header {
            padding: 3px 40px 24px !important;
            .header {
                white-space: unset !important;
                text-align: left !important;
            }
        }
        .bk-dialog-footer {
            padding-bottom: 32px !important;
        }
    }
</style>
