<template>
    <div class="sync-set-wrapper">
        <div class="title clearfix">
            <div class="tips fl">
                <p class="mr10">{{$t('请确认以下模版修改信息：')}}</p>
                <span class="mr30">
                    <i class="dot"></i>
                    {{$t('新增模块')}}
                </span>
                <span>
                    <i class="dot red"></i>
                    {{$t('删除模块')}}
                </span>
            </div>
            <bk-checkbox class="expand-all fr" v-model="expandAll">{{$t('全部展开')}}</bk-checkbox>
        </div>
        <div class="instance-list">
            <set-instance class="instance-item" :expand="true"></set-instance>
        </div>
        <div class="footer">
            <bk-button theme="primary" class="mr10" @click="handleConfirmSync">{{$t('确认同步')}}</bk-button>
            <bk-button class="mr10">{{$t('取消')}}</bk-button>
            <span>{{$t('已选择20个集群实例')}}</span>
        </div>
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
                expandAll: false
            }
        },
        methods: {
            handleConfirmSync () {
                this.$bkInfo({
                    type: 'warning',
                    title: '确定同步拓扑？',
                    // subTitle: '即将同步拓扑模版【正式集群模版】，模版的拓扑结构将会更新到选中的集群实例中',
                    subTitle: '即将批量同步拓扑模版【正式集群模版】到选中的20个集群实例中，模版的拓扑结构将会更新到选中的集群实例中',
                    extCls: 'set-confirm-sync'
                    // confirmFn: async () => {}
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .sync-set-wrapper {
        padding: 0 20px;
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
