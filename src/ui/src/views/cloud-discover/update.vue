<template>
    <div class="detail-wrapper">
        <div class="detail-box">
            <ul class="event-form" v-model="curPush">
                <li class="form-item">
                    <label for="" class="label-name">
                        任务名称<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <input type="text"
                               v-model="curPush.bk_task_name"
                               class="cmdb-form-input"
                               placeholder="请输入任务名称">
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        账号类型<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <select v-model="curPush.bk_account_type">
                            <option value="腾讯云">腾讯云</option>
                        </select>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                      同步周期
                    </label>
                    <div class="item-content length-short">
                        <select name="period-type"
                                v-model="curPush.bk_period_type"
                                @change="getPeriodSelect">
                            <option :value="period.id"
                                    v-for="period in periodList">{{ period.name }}
                            </option>
                        </select>
                        <input type="text"
                               v-model="curPush.bk_period"
                               class="cmdb-form-input"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        ID<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <input type="text"
                               v-model="curPush.bk_secret_id"
                               class="cmdb-form-input"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        Key<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <input type="password"
                               v-model="curPush.bk_secret_key"
                               class="cmdb-form-input"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">同步资源</label>
                    <div class="item-content">
                        <input type="checkbox"
                               v-model="curPush.bk_obj_id"
                               class="cmdb-checkbox-text"/>
                        <label>主机</label>
                        <input type="checkbox"
                               class="cmdb-checkbox-text"
                               style="margin-left: 30px"/>
                        <label>交换机</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">账号管理员</label>
                    <div class="item-content">
                        <input type="text"
                               v-model="curPush.bk_account_admin"
                               class="cmdb-form-input"/>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">资源确认</label>
                    <div class="item-content">
                        <input type="checkbox" v-model="curPush.bk_confirm"/>
                        <label>新增需要确认</label>
                        <input type="checkbox" v-model="curPush.bk_attr_confirm" style="margin-left: 30px"/>
                        <label>属性变化需要确认</label>
                    </div>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button type="primary" :loading="$loading('savePush')" class="btn" @click="update">{{$t('Common["保存"]')}}</bk-button>
            <bk-button type="default" class="btn vice-btn" @click="cancel">{{$t('Common["取消"]')}}</bk-button>
        </footer>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
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
                    name: '每5分钟'
                }, {
                    id: 'hour',
                    name: '每小时'
                }, {
                    id: 'day',
                    name: '每天'
                }]
            }
        },
        computed: {
            ...mapGetters([
                'language'
            ]),
            params () {
                let params = this.curPush
                return params
            }
        },
        methods: {
            ...mapActions('cloudDiscover', ['updateCloudTask']),
            getPeriodSelect () {

            },
            async update () {
                let res = null
                res = await this.updateCloudTask({params: this.params, config: {requestId: 'savePush'}})
                this.$emit('saveSuccess')
                this.$success(this.$t('EventPush["修改成功"]'))
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .detail-wrapper{
        height: 100%;
        .pop-master{
            position: absolute;
            left: 0;
            top: 0;
            bottom: 0;
            right: 0;
        }
        .detail-box{
            padding: 20px 30px;
            height: calc(100% - 63px);
            overflow-y: auto;
            @include scrollbar;
        }
        .event-form{
            .form-item{
                width: 100%;
                margin-bottom: 20px;
                &:after{
                    display: block;
                    content: "";
                    clear: both;
                }
            }
            .label-name{
                position: relative;
                float: left;
                width: 85px;
                text-align: right;
                line-height: 36px;
                font-size: 14px;
                .color-danger{
                    position: absolute;
                    top: 2px;
                    right: -10px;
                }
            }
            .item-content{
                margin-left: 15px;
                width: calc(100% - 100px);
                float: left;
                span {
                    font-size: 14px;
                }
                &.length-short{
                    position: relative;
                    input{
                        width: 97px;
                        margin-right: 10px;
                    }
                }
            }
        }
        .info{
            background: #fff3da;
            border-radius: 2px;
            width: 100%;
            padding-left: 20px;
            height: 42px;
            line-height: 40px;
            font-size: 0;
            border: 1px solid #ffc947;
            .bk-icon {
                position: relative;
                top: -1px;
                margin-right: 10px;
                color: #ffc947;
                font-size: 20px;
            }
            span {
                font-size: 14px;
                vertical-align: middle;
            }
            .num{
                font-weight: bold;
            }
        }
        .event-wrapper{
            .event-box{
                padding: 13px 0;
                border-bottom: 1px solid #eceef5;
                .event-title{
                    margin-right: 10px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                    height: 32px;
                    line-height: 32px;
                    font-weight: bold;
                    cursor: pointer;
                }
                .bk-icon{
                    line-height: 32px;
                    width: 32px;
                    text-align: center;
                    margin-right: -6px;
                    font-size: 12px;
                    font-weight: bold;
                    transition: all .5s;
                    &.up{
                        transform: rotate(180deg);
                    }
                }
                ul{
                    width: 100%;
                    float: left;
                }
            }
            .event-item{
                padding: 7px 0;
                height: 32px;
                line-height: 18px;
                &:after{
                    display: block;
                    content: "";
                    clear: both;
                }
                .label-name{
                    float: left;
                    width: 85px;
                    text-align: right;
                    font-size: 14px;
                    margin-right: 15px;
                    font-weight: bold;
                    white-space: nowrap;
                    text-overflow: ellipsis;
                    overflow: hidden;
                }
                .options{
                    font-size: 0;
                    label{
                        height: 18px;
                        width: 85px;
                        margin-right: 10px;
                        &:nth-child(1) {
                            width: 120px;
                        }
                        &:nth-child(4) {
                            width: 76px;
                            margin-right: 0;
                        }
                        &.cmdb-form-checkbox{
                            padding: 0;
                            cursor: pointer;
                        }
                        input[type="checkbox"]{
                            margin-right: 6px;
                        }
                        .cmdb-checkbox-text{
                            display: inline-block;
                            width: calc(100% - 20px);
                            overflow: hidden;
                            text-overflow: ellipsis;
                            white-space: nowrap;
                        }
                    }
                }
            }
        }
        .footer{
            height: 63px;
            line-height: 63px;
            background: #f9f9f9;
            font-size: 0;
            padding-left: 130px;
            .btn{
                margin-right: 11px;
            }
        }
    }
</style>
