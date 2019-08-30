<template>
    <div class="update-wrapper">
        <div class="update-box">
            <ul class="update-event-form">
                <li class="update-form-item">
                    <label for="" class="label-name">
                        {{ $t('任务名称')}}<span class="color-danger">*</span>
                    </label>
                    <div class="update-item-content">
                        <bk-input type="text"
                            v-model="curPush.bk_task_name"
                            name="taskName"
                            v-validate="'required|singlechar|length:256'"
                            class="cmdb-form-input">
                        </bk-input>
                    </div>
                    <span v-show="errors.has('taskName')" class="color-danger">{{ errors.first('taskName') }}</span>
                </li>
                <li class="update-form-item">
                    <label for="" class="label-name">
                        {{ $t('账号类型')}}<span class="color-danger">*</span>
                    </label>
                    <div class="update-item-content">
                        <cmdb-selector
                            :list="cloudList"
                            v-model="curPush.bk_account_type"
                            name="accountType"
                            v-validate="'required'"
                            :placeholder="$t('请选择账号类型')"
                        ></cmdb-selector>
                    </div>
                    <span v-show="errors.has('accountType')" class="error-info color-danger">{{ errors.first('accountType') }}</span>
                </li>
                <li class="update-form-item">
                    <label for="" class="label-name">
                        {{ $t('ID')}}<span class="color-danger">*</span>
                    </label>
                    <div class="update-item-content">
                        <bk-input type="text"
                            v-model="curPush.bk_secret_id"
                            name="ID"
                            v-validate="'required|singlechar|length:256'"
                            class="cmdb-form-input">
                        </bk-input>
                    </div>
                    <span v-show="errors.has('ID')" class="color-danger">{{ errors.first('ID') }}</span>
                </li>
                <li class="update-form-item">
                    <label for="" class="label-name">
                        {{ $t('Key')}}<span class="color-danger">*</span>
                        <a class="set"
                            href="https://cloud.tencent.com/document/api/213/15692"
                            target="_blank">{{$t('如何获取ID和Key?')}}
                        </a>
                    </label>
                    <div class="update-item-content">
                        <bk-input v-model="curPush.bk_secret_key"
                            class="cmdb-form-input"
                            type="password">
                        </bk-input>
                    </div>
                </li>
                <li class="form-item-two">
                    <label for="" class="label-name-two">{{ $t('同步周期')}}</label>
                    <div class="item-content-two length-short">
                        <bk-select class="selector"
                            v-model="curPush.bk_period_type">
                            <bk-option v-for="(option, index) in periodList"
                                :key="index"
                                :id="option.id"
                                :name="option.name">
                            </bk-option>
                        </bk-select>
                        <bk-input type="text"
                            class="cmdb-form-input"
                            v-model="curPush.bk_period"
                            v-if="curPush.bk_period_type === 'day'"
                            name="day"
                            v-validate="'required|dayFormat'"
                            :placeholder="$t('例如: 19:30')">
                        </bk-input>
                        <bk-input type="text"
                            class="cmdb-form-input"
                            v-model="curPush.bk_period"
                            v-if="curPush.bk_period_type === 'hour'"
                            name="hour"
                            v-validate="'required|hourFormat'"
                            :placeholder="$t('例如: 30')">
                        </bk-input>
                        <div v-show="errors.has('day')" class="update-error-info color-danger">{{ errors.first('day') }}</div>
                        <div v-show="errors.has('hour')" class="update-error-info color-danger">{{ errors.first('hour') }}</div>
                    </div>
                </li>
                <li class="update-form-item">
                    <label for="" class="label-name">{{ $t('任务维护人')}}</label>
                    <cmdb-form-objuser
                        class="fl maintain-selector"
                        v-model="curPush.bk_account_admin"
                        :multiple="true"
                        name="maintain"
                        v-validate="'required|singlechar|length:256'">
                    </cmdb-form-objuser>
                    <span v-show="errors.has('maintain')" class="color-danger">{{ errors.first('maintain') }}</span>
                </li>
                <li>
                    <label>{{ $t('同步资源')}}</label>
                    <div>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" value="host" v-model="curPush.bk_obj_id" disabled>
                            <span class="cmdb-checkbox-text">{{ $t('主机')}}</span>
                        </label>
                    </div>
                </li>
                <li>
                    <div class="u-resource-confirm">{{ $t('资源自动确认')}}
                        <span class="span-text">{{ $t('(不勾选，发现实例将不需要确认直接录入主机资源池)')}}</span>
                    </div>
                    <div>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" v-model="curPush.bk_confirm">
                            <span class="cmdb-checkbox-text">{{ $t('新增需要确认')}}</span>
                        </label>
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox" v-model="curPush.bk_attr_confirm">
                            <span class="cmdb-checkbox-text">{{ $t('属性变化需要确认')}}</span>
                        </label>
                    </div>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button theme="primary" :loading="$loading('savePush')" class="btn" @click="update">{{$t('保存')}}</bk-button>
            <bk-button theme="default" class="btn vice-btn" @click="cancel">{{$t('取消')}}</bk-button>
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
                tips: false,
                placeholder: '',
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
                }]
            }
        },
        watch: {
            'curPush.bk_period_type' () {
                if (this.curPush.bk_period_type === 'hour') {
                    this.placeholder = this.$t('例如: 30')
                } else {
                    this.placeholder = this.$t('例如: 19:30')
                }
            },
            'curPush': {
                handler () {
                    this.tips = true
                },
                deep: true
            }
        },
        methods: {
            ...mapActions('cloudDiscover', ['updateCloudTask']),
            async update () {
                const isValidate = await this.$validator.validateAll()
                if (!isValidate) {
                    return
                }
                const params = this.curPush
                await this.updateCloudTask({ params: params, config: { requestId: 'savePush' } })
                this.$emit('saveSuccess')
                this.$success(this.$t('修改成功'))
            },
            cancel () {
                this.$emit('cancel')
            },
            isCloseConfirmShow () {
                if (this.tips) {
                    return true
                }
                return false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .update-wrapper {
        height: 100%;
        .update-box {
            padding: 17px 20px 0 21px;
        }
        .update-event-form {
            .update-form-item {
                width: 300px;
                height: 63px;
                margin-bottom: 17px;
                float: left;
                &:after {
                    display: block;
                    content: "";
                    clear: both;
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
                    .set {
                        font-size: 8px;
                        float: right;
                        color: dodgerblue;
                    }
                }
                .update-item-content {
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
            .update-form-item:nth-child(even) {
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
                .update-error-info {
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
            .u-resource-confirm {
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
