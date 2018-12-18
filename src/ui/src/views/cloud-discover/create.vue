<template>
    <div class="detail-wrapper">
        <div class="detail-box">
            <ul class="event-form" v-model="taskMap">
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["任务名称"]')}}<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <input type="text"
                               v-model="taskMap.bk_task_name"
                               class="cmdb-form-input"
                               :placeholder="$t('Cloud[\'请输入任务名称\']')"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["账号类型"]')}}<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <bk-selector
                            :list="cloudList"
                            :selected.sync="taskMap.bk_account_type"
                            :placeholder="$t('Cloud[\'请选择账号类型\']')"
                        ></bk-selector>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["ID"]')}}<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <input type="text"
                               v-model="taskMap.bk_secret_id"
                               class="cmdb-form-input"
                               :placeholder="$t('Cloud[\'请输入ID\']')"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{ $t('Cloud["Key"]')}}<span class="color-danger">*</span>
                        <a class="a-set"
                           href="https://cloud.tencent.com/document/api/213/15692"
                           target="_blank">{{ $t('Cloud["如何获取ID和Key?"]')}}
                        </a>
                    </label>
                    <div class="item-content">
                        <input type="password"
                               v-model="taskMap.bk_secret_key"
                               class="cmdb-form-input"
                               :placeholder="$t('Cloud[\'请输入key\']')"/>
                    </div>
                </li>
                <li class="form-item-two">
                    <label for="" class="label-name-two">
                        {{ $t('Cloud["同步周期"]')}}
                    </label>
                    <div class="item-content-two length-short">
                        <bk-selector
                            class="selector"
                            :list="periodList"
                            :selected.sync="taskMap.bk_period_type">
                        </bk-selector>
                        <input type="text"
                               v-model="taskMap.bk_period"
                               class="cmdb-form-input"
                               :disabled = "disabled"
                               :placeholder="placeholder"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">{{ $t('Cloud["任务维护人"]')}}</label>
                    <div class="item-content">
                        <input type="text"
                               v-model="taskMap.bk_account_admin"
                               class="cmdb-form-input"/>
                    </div>
                </li>
                <li>
                    <label class="resource-lable">{{ $t('Cloud["同步资源"]')}}</label>
                    <div>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" value="host" v-model="taskMap.bk_obj_id" disabled>
                            <span class="cmdb-checkbox-text">{{ $t('Hosts["主机"]')}}</span>
                        </label>
                    </div>
                </li>
                <li>
                    <label class="resource-confirm">{{ $t('Cloud["资源自动确认"]')}}
                    <span class="span-text">{{ $t('Cloud["(不勾选，发现实例将不需要确认直接录入主机资源池)"]')}}</span>
                    </label>
                    <div>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" v-model="taskMap.bk_confirm">
                            <span class="cmdb-checkbox-text">{{ $t('Cloud["新增需要确认"]')}}</span>
                        </label>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" v-model="taskMap.bk_attr_confirm">
                            <span class="cmdb-checkbox-text">{{ $t('Cloud["属性变化需要确认"]')}}</span>
                        </label>
                    </div>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button type="primary" :loading="$loading('savePush')" class="btn" @click="save">{{$t('Common["保存"]')}}</bk-button>
            <bk-button type="default" class="btn vice-btn" @click="cancel">{{$t('Common["取消"]')}}</bk-button>
        </footer>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            type: {
                type: String,
                default: 'create'
            }
        },
        data () {
            return {
                disabled: false,
                placeholder: this.$t('Cloud["例如: 19:30"]'),
                cloudList: [{
                    id: 'tencent_cloud',
                    name: this.$t('Cloud["腾讯云"]')
                }],
                periodList: [{
                    id: 'minute',
                    name: this.$t('Cloud["每五分钟"]')
                }, {
                    id: 'hour',
                    name: this.$t('Cloud["每小时"]')
                }, {
                    id: 'day',
                    name: this.$t('Cloud["每天"]')
                }],
                taskMap: {
                    bk_task_name: '',
                    bk_account_type: this.$t('Cloud["腾讯云"]'),
                    bk_period_type: 'day',
                    bk_secret_id: '',
                    bk_secret_key: '',
                    bk_obj_id: 'host',
                    bk_account_admin: '',
                    bk_confirm: false,
                    bk_attr_confirm: false,
                    bk_period: ''
                },
                tempTaskMap: {
                    bk_task_name: '',
                    bk_account_type: this.$t('Cloud["腾讯云"]'),
                    bk_period_type: 'day',
                    bk_secret_id: '',
                    bk_secret_key: '',
                    bk_obj_id: 'host',
                    bk_account_admin: '',
                    bk_confirm: false,
                    bk_attr_confirm: false,
                    bk_period: ''
                }
            }
        },
        methods: {
            ...mapActions('cloudDiscover', ['addCloudTask']),
            async save () {
                let params = this.taskMap
                let res = null
                res = await this.addCloudTask({params: params, config: {requestId: 'savePush'}})
                this.$emit('saveSuccess')
                this.$success(this.$t('Inst["创建成功"]'))
            },
            cancel () {
                this.$emit('cancel')
            },
            isCloseConfirmShow () {
                let tempTaskMap = this.tempTaskMap
                let taskMap = this.taskMap
                for (let key in taskMap) {
                    if (taskMap[key] !== tempTaskMap[key]) {
                        return true
                    }
                }
                return false
            }
        },
        watch: {
            'taskMap.bk_period_type' () {
                if (this.taskMap.bk_period_type === 'minute') {
                    this.disabled = true
                    this.placeholder = ''
                } else if (this.taskMap.bk_period_type === 'hour') {
                    this.disabled = false
                    this.placeholder = this.$t('Cloud["例如: 30"]')
                } else {
                    this.disabled = false
                    this.placeholder = this.$t('Cloud["例如: 19:30"]')
                }
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
        .event-form {
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
                    text-align: right;
                    line-height: 27px;
                    font-size: 14px;
                    .color-danger {
                        position: absolute;
                        top: 2px;
                        right: -10px;
                    }
                    .a-set {
                        font-size: 8px;
                        float: right;
                        color: dodgerblue;
                    }
                }
                .item-content {
                    span {
                        font-size: 14px;
                    }
                    &.length-short {
                         position: relative;
                        input {
                            width: 97px;
                            margin-right: 10px;
                        }
                    }
                }
            }
            .form-item:nth-child(even) {
                margin-left: 35px;
            }
            .form-item-two {
                position: relative;
                width: 300px;
                height: 63px;
                margin-bottom: 17px;
                float: left;
                &:after {
                     display: block;
                     content: "";
                     clear: both;
                 }
                .label-name-two {
                        position: relative;
                        width: 85px;
                        text-align: right;
                        line-height: 27px;
                        font-size: 14px;
                }
                .item-content-two {
                    span {
                        font-size: 14px;
                    }
                    &.length-short {
                         position: relative;
                        .selector {
                            position: absolute;
                            width: 144px;
                            left: 0;
                            top: 0;
                        }
                        input {
                            width: 150px;
                            margin-left: 5px;
                            position: absolute;
                            right: 0;
                            top: 0;
                        }
                    }
                }
            }
            .resource-lable {
                height: 19px;
                width: 56px;
            }
            .resource-confirm {
                width: 416px;
                height: 19px;
                margin-bottom: 11px;
                margin-top: 23px;
                .span-text {
                    opacity:0.5;
                }
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
