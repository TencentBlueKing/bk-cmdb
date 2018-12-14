<template>
    <div class="detail-wrapper">
        <div class="detail-box">
            <ul class="cloud-form" v-model="curPush">
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["任务名称"]')}} :
                    </label>
                    <div class="item-content">
                        <label>{{curPush.bk_task_name}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["账号类型"]')}} :
                    </label>
                    <div class="item-content">
                        <label>{{curPush.bk_account_type}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["ID"]')}} :
                    </label>
                    <div class="item-content">
                        <span>{{curPush.bk_secret_id}}</span>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["Key"]')}} :
                    </label>
                    <div class="item-content">
                        <span>*************</span>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">{{ $t('Cloud["同步资源"]')}} : </label>
                    <div class="item-content">
                        <label>{{curPush.bk_obj_id}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">{{ $t('Cloud["自动同步"]')}} : </label>
                    <div class="item-content">
                        <label>{{curPush.bk_period_type}}</label>
                        <label>{{curPush.bk_period}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">{{ $t('Cloud["任务维护人"]')}} : </label>
                    <div class="item-content">
                        <label>{{curPush.bk_account_admin}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">{{ $t('Cloud["资源确认"]')}} : </label>
                    <div class="item-content">
                        <span v-if="curPush.bk_confirm">
                            {{ $t('Cloud["新增需要确认"]')}}
                        </span>
                        <span v-if="curPush.bk_attr_confirm">
                            {{ $t('Cloud["属性变化需要确认"]')}}
                        </span>
                    </div>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button type="primary" :loading="$loading('savePush')" class="btn" @click="edit">{{$t('Common["编辑"]')}}</bk-button>
        </footer>
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
        data () {
            return {
                periodList: [{
                    id: 'minute',
                    name: this.$t('Cloud["每五分钟"]')
                }, {
                    id: 'hour',
                    name: this.$t('Cloud["每小时"]')
                }, {
                    id: 'day',
                    name: this.$t('Cloud["每天"]')
                }]
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
    .detail-wrapper {
        height: 100%;
        .detail-box {
            padding: 17px 20px 0 21px;
        }
        .cloud-form {
            .form-item {
                width: 300px;
                height: 63px;
                margin-bottom: 17px;
                float: left;
                &:after {
                    display: block;
                    content: "";
                    clear: both;
                }
                .label-name {
                    position: relative;
                    width: 85px;
                    float: left;
                    text-align: right;
                    line-height: 27px;
                    font-size: 14px;
                }
                .item-content {
                    float: left;
                    line-height: 27px;
                    width: 200px;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    overflow: hidden;
                    word-break: break-all;
                    span {
                        font-size: 14px;
                    }
                }
            }
            .form-item:nth-child(even) {
                margin-left: 35px;
            }
        }
        .footer {
            height: 63px;
            line-height: 63px;
            background: #f9f9f9;
            font-size: 0;
            padding-left: 24px;
            .btn {
                margin-right: 10px;
            }
        }
    }
</style>
