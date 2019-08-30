<template>
    <div class="task-detail-wrapper">
        <div class="task-detail-box">
            <ul class="cloud-form clearfix">
                <li class="detail-form-item">
                    <label for="" class="label-name">
                        {{ $t('任务名称')}} ：
                    </label>
                    <div class="detail-item-content">
                        <span>{{curPush.bk_task_name}}</span>
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">
                        {{ $t('账号类型')}} ：
                    </label>
                    <div class="detail-item-content">
                        <span v-if="curPush.bk_account_type === 'tencent_cloud'">
                            {{$t('腾讯云')}}
                        </span>
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">
                        {{ $t('ID')}} ：
                    </label>
                    <div class="detail-item-content">
                        <span>{{curPush.bk_secret_id}}</span>
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">
                        {{ $t('Key')}} ：
                    </label>
                    <div class="detail-item-content">
                        *************
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">{{ $t('同步资源')}} ：</label>
                    <div class="detail-item-content">
                        <span>{{curPush.bk_obj_id}}</span>
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">{{ $t('自动同步')}} ：</label>
                    <div class="detail-item-content">
                        <span v-if="curPush.bk_period_type === 'minute'">
                            {{$t('每五分钟')}}
                        </span>
                        <span v-else-if="curPush.bk_period_type === 'hour'">
                            {{this.$t('每小时')}} {{curPush.bk_period}}
                        </span>
                        <span v-else>
                            {{this.$t('每天')}} {{curPush.bk_period}}
                        </span>
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">{{ $t('任务维护人')}} ：</label>
                    <div class="detail-item-content">
                        <span>{{curPush.bk_account_admin}}</span>
                    </div>
                </li>
                <li class="detail-form-item">
                    <label for="" class="label-name">{{ $t('资源确认')}} ：</label>
                    <div class="detail-item-content">
                        <span v-if="curPush.bk_attr_confirm && curPush.bk_confirm">
                            {{ $t('新增需要确认、属性变化需要确认')}}
                        </span>
                        <span v-else-if="curPush.bk_confirm">
                            {{ $t('新增需要确认')}}
                        </span>
                        <span v-else-if="curPush.bk_attr_confirm">
                            {{ $t('属性变化需要确认')}}
                        </span>
                        <span class="text-opacity" v-else>
                            {{ $t('直接入库，不需要确认')}}
                        </span>
                    </div>
                </li>
            </ul>
        </div>
        <div class="task-detail-btn">
            <bk-button theme="primary" :loading="$loading('savePush')" class="btn" @click="edit">{{$t('编辑')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            curPush: {
                type: Object
            },
            type: {
                type: String,
                default: 'create'
            }
        },
        methods: {
            ...mapActions('cloudDiscover', ['addCloudTask']),
            edit () {
                this.$emit('edit', this.curPush)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .task-detail-wrapper {
        height: 100%;
        .task-detail-box {
            padding: 17px 20px 0 21px;
        }
        .cloud-form {
            .detail-form-item {
                width: 300px;
                height: 24px;
                margin-bottom: 15px;
                float: left;
                &:after {
                    display: block;
                    content: "";
                    clear: both;
                }
                .label-name {
                    position: relative;
                    width: 90px;
                    float: left;
                    text-align: right;
                    line-height: 24px;
                    font-size: 14px;
                }
                .detail-item-content {
                    float: left;
                    line-height: 27px;
                    width: 200px;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    overflow: hidden;
                    word-break: break-all;
                    .text-opacity {
                        opacity:0.5;
                    }
                    span {
                        font-size: 14px;
                    }
                }
            }
            .detail-form-item:nth-child(even) {
                margin-left: 35px;
            }
        }
        .task-detail-btn {
            height: 63px;
            line-height: 63px;
            font-size: 0;
            padding-left: 100px;
            .btn {
                margin-right: 10px;
            }
        }
    }
</style>
