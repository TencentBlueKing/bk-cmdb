<template>
    <div class="create-wrapper">
        <div class="create-box">
            <ul class="create-form">
                <li class="create-form-item">
                    <label for="" class="label-name">
                        {{ $t('任务名称')}}<span class="color-danger">*</span>
                    </label>
                    <div class="create-item-content">
                        <bk-input type="text"
                            v-model="taskMap.bk_task_name"
                            class="cmdb-form-input"
                            name="taskName"
                            v-validate="'required|singlechar|length:256'"
                            :placeholder="$t('请输入任务名称')">
                        </bk-input>
                    </div>
                    <span v-show="errors.has('taskName')" class="error-info color-danger">{{ errors.first('taskName') }}</span>
                </li>
                <li class="create-form-item">
                    <label for="" class="label-name">
                        {{ $t('账号类型')}}<span class="color-danger">*</span>
                    </label>
                    <div class="create-item-content">
                        <cmdb-selector
                            :list="cloudList"
                            v-model="taskMap.bk_account_type"
                            name="accountType"
                            v-validate="'required'"
                            :placeholder="$t('请选择账号类型')"
                        ></cmdb-selector>
                    </div>
                    <span v-show="errors.has('accountType')" class="error-info color-danger">{{ errors.first('accountType') }}</span>
                </li>
                <li class="create-form-item">
                    <label for="" class="label-name">
                        {{ $t('ID')}}<span class="color-danger">*</span>
                    </label>
                    <div class="create-item-content">
                        <bk-input
                            type="text"
                            v-model="taskMap.bk_secret_id"
                            class="cmdb-form-input"
                            name="ID"
                            v-validate="'required|singlechar|length:256'"
                            :placeholder="$t('请输入ID')">
                        </bk-input>
                    </div>
                    <span v-show="errors.has('ID')" class="error-info color-danger">{{ errors.first('ID') }}</span>
                </li>
                <li class="create-form-item">
                    <label for="" class="label-name">
                        {{ $t('Key')}}<span class="color-danger">*</span>
                        <a class="a-set"
                            href="https://cloud.tencent.com/document/api/213/15692"
                            target="_blank">{{ $t('如何获取ID和Key?')}}
                        </a>
                    </label>
                    <div class="create-item-content">
                        <bk-input
                            type="password"
                            v-model="taskMap.bk_secret_key"
                            class="cmdb-form-input"
                            name="Key"
                            v-validate="'required|singlechar|length:256'"
                            :placeholder="$t('请输入key')">
                        </bk-input>
                    </div>
                    <span v-show="errors.has('Key')" class="error-info color-danger">{{ errors.first('Key') }}</span>
                </li>
                <li class="form-item-two">
                    <label for="" class="label-name-two">
                        {{ $t('同步周期')}}
                    </label>
                    <div class="item-content-two length-short">
                        <bk-select v-model="taskMap.bk_period_type">
                            <bk-option v-for="(option, index) in periodList"
                                :key="index"
                                :id="option.id"
                                :name="option.name">
                            </bk-option>
                        </bk-select>
                        <bk-input
                            type="text"
                            class="cmdb-form-input"
                            v-model="taskMap.bk_period"
                            v-if="taskMap.bk_period_type === 'day'"
                            name="day"
                            v-validate="'required|dayFormat'"
                            :placeholder="$t('例如: 19:30')">
                        </bk-input>
                        <bk-input type="text"
                            class="cmdb-form-input"
                            v-model="taskMap.bk_period"
                            v-if="taskMap.bk_period_type === 'hour'"
                            name="hour"
                            v-validate="'required|hourFormat'"
                            :placeholder="$t('例如: 30')">
                        </bk-input>
                        <div v-show="errors.has('hour')" class="error-info-two color-danger">{{ errors.first('hour') }}</div>
                        <div v-show="errors.has('day')" class="error-info-two color-danger">{{ errors.first('day') }}</div>
                    </div>
                </li>
                <li class="create-form-item">
                    <label for="" class="label-name">{{ $t('任务维护人')}}</label>
                    <cmdb-form-objuser
                        class="fl maintain-selector"
                        v-model="taskMap.bk_account_admin"
                        :multiple="true"
                        name="maintain"
                        v-validate="'required|singlechar|length:256'">
                    </cmdb-form-objuser>
                    <div v-show="errors.has('maintain')" class="error-info color-danger">{{ errors.first('maintain') }}</div>
                </li>
                <li>
                    <label class="resource-lable">{{ $t('同步资源')}}</label>
                    <div>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" value="host" v-model="taskMap.bk_obj_id" disabled>
                            <span class="cmdb-checkbox-text">{{ $t('主机')}}</span>
                        </label>
                    </div>
                </li>
                <li>
                    <div class="resource-confirm">{{ $t('资源自动确认')}}
                        <span class="span-text">{{ $t('(不勾选，发现实例将不需要确认直接录入主机资源池)')}}</span>
                    </div>
                    <div>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" v-model="taskMap.bk_confirm">
                            <span class="cmdb-checkbox-text">{{ $t('新增需要确认')}}</span>
                        </label>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" v-model="taskMap.bk_attr_confirm">
                            <span class="cmdb-checkbox-text">{{ $t('属性变化需要确认')}}</span>
                        </label>
                    </div>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button theme="primary" :loading="$loading('savePush')" class="btn" @click="save">{{$t('保存')}}</bk-button>
            <bk-button theme="default" class="btn vice-btn" @click="cancel">{{$t('取消')}}</bk-button>
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
                timeShow: true,
                placeholder: this.$t('例如: 19:30'),
                cloudList: [{
                    id: 'tencent_cloud',
                    name: this.$t('腾讯云')
                }],
                periodList: [{
                    id: 'minute',
                    name: this.$t('每五分钟')
                }, {
                    id: 'hour',
                    name: this.$t('每小时')
                }, {
                    id: 'day',
                    name: this.$t('每天')
                }],
                taskMap: {
                    bk_task_name: '',
                    bk_account_type: this.$t('腾讯云'),
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
                    bk_account_type: this.$t('腾讯云'),
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
                const isValidate = await this.$validator.validateAll()
                if (!isValidate) {
                    return
                }
                const params = this.taskMap
                await this.addCloudTask({ params: params, config: { requestId: 'savePush' } })
                this.$emit('saveSuccess')
                this.$success(this.$t('创建成功'))
            },
            cancel () {
                this.$emit('cancel')
            },
            isCloseConfirmShow () {
                const tempTaskMap = this.tempTaskMap
                const taskMap = this.taskMap
                for (const key in taskMap) {
                    if (taskMap[key] !== tempTaskMap[key]) {
                        return true
                    }
                }
                return false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-wrapper {
        height: 100%;
        .create-box {
            padding: 17px 20px 0 21px;
        }
        .create-form {
            .create-form-item {
                width: 300px;
                height: 63px;
                margin-bottom: 17px;
                float: left;
                &:after {
                    display: block;
                    content: "";
                    clear: both;
                }
                .error-info {
                    font-size: 12px;
                }
                .maintain-selector{
                    width: 300px;
                    line-height: initial;
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
                .create-item-content {
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
            .create-form-item:nth-child(even) {
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
                .error-info-two {
                    position:absolute;
                    top:100%;
                    font-size: 12px;
                    padding-left: 150px;
                 }
                .label-name-two {
                        position: relative;
                        width: 85px;
                        text-align: right;
                        line-height: 27px;
                        font-size: 14px;
                }
                .item-content-two {
                    height: 36px;
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
                height: 19px;
                margin-top: 20px;
                .span-text {
                    opacity:0.5;
                }
            }
        }
        .footer {
            height: 63px;
            line-height: 63px;
            font-size: 0;
            padding-left: 24px;
            .btn {
                margin-right: 10px;
            }
        }
    }
</style>
