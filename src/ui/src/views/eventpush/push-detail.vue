<template>
    <div class="detail-wrapper">
        <div class="detail-box">
            <ul class="event-form">
                <li class="form-item">
                    <label for="" class="label-name">
                        {{$t('订阅名称')}}<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <bk-input type="text" class="cmdb-form-input" :placeholder="$t('请输入订阅名称')"
                            maxlength="20"
                            v-model.trim="tempEventData['subscription_name']"
                            :data-vv-name="$t('订阅名称')"
                            v-validate="'required'">
                        </bk-input>
                        <span v-show="errors.has($t('订阅名称'))" class="color-danger">{{ errors.first($t('订阅名称')) }}</span>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{$t('系统名称')}}
                    </label>
                    <div class="item-content">
                        <bk-input type="text" class="cmdb-form-input" :placeholder="$t('请输入系统名称')"
                            v-model.trim="tempEventData['system_name']"
                            v-validate="'singlechar|length:256'"
                            :data-vv-name="$t('系统名称')">
                        </bk-input>
                        <span v-show="errors.has($t('系统名称'))" class="color-danger">{{ errors.first($t('系统名称')) }}</span>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        URL<span class="color-danger">*</span>
                    </label>
                    <div class="item-content url">
                        <div class="url-box">
                            <bk-input type="text" class="cmdb-form-input" :placeholder="$t('请输入URL')"
                                v-model.trim="tempEventData['callback_url']"
                                v-validate="'required|http'"
                                name="http">
                            </bk-input>
                            <span v-show="errors.has('http')" class="color-danger">{{ errors.first('http') }}</span>
                        </div>
                        <bk-button class="test-btn" theme="primary" @click.prevent="testPush">{{$t('测试推送')}}</bk-button>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{$t('成功确认方式')}}<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <label for="http" class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="confimType" id="http" value="httpstatus"
                                v-model="tempEventData['confirm_mode']"
                            >{{$t('HTTP状态')}}
                        </label>
                        <label for="reg" class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="confimType" id="reg" value="regular"
                                v-model="tempEventData['confirm_mode']"
                            >{{$t('正则验证')}}
                        </label>
                        <div class="input-box">
                            <bk-input type="text" class="cmdb-form-input" :placeholder="$t('请输入正则验证')"
                                v-if="tempEventData['confirm_mode'] === 'regular'"
                                v-model.trim="tempEventData['confirm_pattern']['regular']"
                                :data-vv-name="$t('该字段')"
                                v-validate="'required'">
                            </bk-input>
                            <bk-input type="text" class="cmdb-form-input" :placeholder="$t('成功标志')"
                                v-else
                                v-model.trim="tempEventData['confirm_pattern']['httpstatus']"
                                v-validate="{ required: true, regex: /^[0-9]+$/ }"
                                :data-vv-name="$t('该字段')">
                            </bk-input>
                            <i class="tip" :class="{ 'reg': tempEventData['confirm_mode'] === 'regular' }"></i>
                        </div>
                        <span v-show="errors.has($t('该字段'))" class="color-danger">{{ errors.first($t('该字段')) }}</span>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        {{$t('超时时间')}}<span class="color-danger">*</span>
                    </label>
                    <div class="item-content length-short">
                        <bk-input type="text" class="cmdb-form-input" :placeholder="$t('单位：秒')"
                            v-model.trim="tempEventData['time_out']"
                            v-validate="{ required: true, regex: /^[0-9]+$/ }"
                            :data-vv-name="$t('超时时间')"
                            maxlength="10">
                        </bk-input>
                        <span class="unit">S</span>
                        <div v-show="errors.has($t('超时时间'))" class="color-danger">{{ errors.first($t('超时时间')) }}</div>
                    </div>
                </li>
            </ul>
            <div class="info">
                <i class="bk-icon icon-exclamation-circle"></i>
                <span :class="{ 'color-danger': subscriptionFormError }">{{$t('至少选择1个事件')}}</span><i18n path="已选择"><span class="num" place="number">{{selectNum}}</span></i18n>
            </div>
            <ul class="event-wrapper">
                <li class="event-box clearfix"
                    :key="index"
                    v-for="(classify, index) in eventPushList">
                    <div class="event-title" @click="toggleEventList(classify)">
                        <i class="bk-icon icon-angle-down" :class="{ 'up': classify.isHidden }"></i>
                        {{classify.name}}
                    </div>
                    <transition name="slide">
                        <ul v-if="!classify.isHidden" :style="eventHeight(classify.children.length)">
                            <li v-for="(item, idx) in classify.children" :key="idx" class="event-item">
                                <template v-if="item.id === 'resource' && classify.classificationId === 'bk_host_manage'">
                                    <label for="" class="label-name" :title="item.name">{{item.name}}</label>
                                    <div class="options">
                                        <label for="resourceall" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                value="resourceall"
                                                :checked="tempEventData['subscription_form'][item.id].length === 2"
                                                id="resourceall" @change="checkAll('resource')"><i class="cmdb-checkbox-text" :title="$t('全选')">{{$t('全选')}}</i>
                                        </label>
                                        <label for="hostcreate" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                value="hostcreate"
                                                id="hostcreate"><i class="cmdb-checkbox-text" :title="$t('新增主机')">{{$t('新增主机')}}</i>
                                        </label>
                                        <label for="hostdelete" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                value="hostdelete"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                id="hostdelete"><i class="cmdb-checkbox-text" :title="$t('删除主机')">{{$t('删除主机')}}</i>
                                        </label>
                                    </div>
                                </template>
                                <template v-else-if="item.id === 'host'">
                                    <label for="" class="label-name" :title="item.name">{{item.name}}</label>
                                    <div class="options">
                                        <label :for="'hostall'" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                :value="'hostall'"
                                                :checked="tempEventData['subscription_form'][item.id].length === 3"
                                                :id="'hostall'" @change="checkAll('host')"><i class="cmdb-checkbox-text" :title="$t('全选')">{{$t('全选')}}</i>
                                        </label>
                                        <label :for="item.id + 'update'" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                :value="item.id + 'update'"
                                                :id="item.id + 'update'"><i class="cmdb-checkbox-text" :title="$t('编辑')">{{$t('编辑')}}</i>
                                        </label>
                                        <label for="moduletransfer" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                value="moduletransfer"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                id="moduletransfer"><i class="cmdb-checkbox-text" :title="$t('模块转移')">{{$t('模块转移')}}</i>
                                        </label>
                                        <label for="hostidentifier" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                value="hostidentifier"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                id="hostidentifier"><i class="cmdb-checkbox-text" :title="$t('主机身份')">{{$t('主机身份')}}</i>
                                        </label>
                                    </div>
                                </template>
                                <template v-else>
                                    <label for="" class="label-name" :title="item.name">{{item.name}}</label>
                                    <div class="options">
                                        <label :for="`${item.id}all`" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                :value="`${item.id}all`"
                                                @change="checkAll(item.id)"
                                                :checked="tempEventData['subscription_form'][item.id].length === 3"
                                                :id="`${item.id}all`"><i class="cmdb-checkbox-text" :title="$t('全选')">{{$t('全选')}}</i>
                                        </label>
                                        <label :for="`${item.id}create`" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                :value="`${item.id}create`"
                                                :id="`${item.id}create`"><i class="cmdb-checkbox-text" :title="$t('新增')">{{$t('新增')}}</i>
                                        </label>
                                        <label :for="`${item.id}update`" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                :value="`${item.id}update`"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                :id="`${item.id}update`"><i class="cmdb-checkbox-text" :title="$t('编辑')">{{$t('编辑')}}</i>
                                        </label>
                                        <label :for="`${item.id}delete`" class="cmdb-form-checkbox cmdb-checkbox-small">
                                            <input type="checkbox"
                                                :value="`${item.id}delete`"
                                                v-model="tempEventData['subscription_form'][item.id]"
                                                :id="`${item.id}delete`"><i class="cmdb-checkbox-text" :title="$t('删除')">{{$t('删除')}}</i>
                                        </label>
                                    </div>
                                </template>
                            </li>
                        </ul>
                    </transition>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button theme="primary" :loading="$loading('savePush')" class="btn" @click="save">{{type === 'create' ? $t('提交') : $t('保存')}}</bk-button>
            <bk-button theme="default" class="btn vice-btn" @click="cancel">{{$t('取消')}}</bk-button>
        </footer>
        <v-pop
            v-if="isPopShow"
            :callback-u-r-l="tempEventData['callback_url']"
            @closePop="isPopShow = false"
        >
        </v-pop>
    </div>
</template>

<script>
    import vPop from './pop'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            vPop
        },
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
                isPopShow: false,
                subscriptionFormError: false,
                eventPushList: [],
                tempEventData: {
                    subscription_name: '', // 订阅名
                    system_name: '', // 系统名称
                    callback_url: '', // 回调地址
                    confirm_mode: 'httpstatus', // 回调成功确认模式   httpstatus/regular
                    confirm_pattern: {
                        httpstatus: '200',
                        regular: ''
                    },
                    subscription_form: {}, // 订阅单
                    time_out: 60 // 超时时间  单位 秒
                },
                eventData: {
                    subscription_name: '',
                    system_name: '',
                    callback_url: '',
                    confirm_mode: 'httpstatus',
                    confirm_pattern: '200',
                    subscription_form: [],
                    time_out: 60
                }
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', ['classifications']),
            selectNum () {
                let num = 0
                const {
                    subscription_form: subscriptionForm
                } = this.tempEventData
                for (const key in subscriptionForm) {
                    const item = subscriptionForm[key]
                    if (item.length) {
                        num += item.length
                    }
                }
                if (num) {
                    this.subscriptionFormError = false
                }
                return num
            },
            params () {
                const params = this.$tools.clone(this.tempEventData)
                params['confirm_pattern'] = this.tempEventData['confirm_mode'] === 'httpstatus' ? this.tempEventData['confirm_pattern']['httpstatus'] : this.tempEventData['confirm_pattern']['regular']
                const subscriptionForm = []
                for (const key in params['subscription_form']) {
                    if (params['subscription_form'][key].length) {
                        subscriptionForm.push(params['subscription_form'][key].join(','))
                    }
                }
                params['subscription_form'] = subscriptionForm.join(',')
                params['time_out'] = parseInt(params['time_out'])
                return params
            }
        },
        created () {
            this.setEventPushList()
            this.initData()
        },
        methods: {
            ...mapActions('eventSub', [
                'subscribeEvent',
                'updateEventSubscribe'
            ]),
            isCloseConfirmShow () {
                const tempEventData = this.tempEventData
                const eventData = this.eventData
                for (const key in tempEventData) {
                    if (key === 'confirm_pattern') {
                        if (tempEventData[key][tempEventData['confirm_mode']] !== eventData[key][eventData['confirm_mode']]) {
                            return true
                        }
                    } else if (key === 'subscription_form') {
                        if (this.type === 'create') {
                            if (this.selectNum) {
                                return true
                            }
                        } else {
                            const tempList = JSON.stringify(tempEventData[key])
                            const list = JSON.stringify(eventData[key])
                            if (tempList !== list) {
                                return true
                            }
                        }
                    } else {
                        if (tempEventData[key] !== eventData[key]) {
                            return true
                        }
                    }
                }
                return false
            },
            testPush () {
                this.$validator.validate('http').then(res => {
                    if (res) {
                        this.isPopShow = true
                    }
                })
            },
            async checkParams () {
                let result = false
                await this.$validator.validateAll().then(res => {
                    result = res
                })
                if (!result) {
                    return false
                }
                if (this.selectNum === 0) {
                    this.subscriptionFormError = true
                    return false
                }
                return true
            },
            async save () {
                if (!await this.checkParams()) {
                    return
                }
                // eslint-disable-next-line
                let res = null
                if (this.type === 'create') {
                    res = await this.subscribeEvent({ bkBizId: 0, params: this.params, config: { requestId: 'savePush' } })
                } else {
                    res = await this.updateEventSubscribe({ bkBizId: 0, subscriptionId: this.curPush['subscription_id'], params: this.params, config: { requestId: 'savePush' } })
                }
                this.$emit('saveSuccess')
                this.$success(this.$t('保存成功'))
                this.eventData = { ...this.tempEventData }
            },
            cancel () {
                this.$emit('cancel')
            },
            checkAll (objId) {
                if (event.target.checked) {
                    if (objId === 'resource') {
                        this.tempEventData['subscription_form'][objId] = ['hostcreate', 'hostdelete']
                    } else if (objId === 'host') {
                        this.tempEventData['subscription_form'][objId] = ['moduletransfer', 'hostupdate', 'hostidentifier']
                    } else {
                        this.tempEventData['subscription_form'][objId] = [`${objId}create`, `${objId}update`, `${objId}delete`]
                    }
                } else {
                    this.tempEventData['subscription_form'][objId] = []
                }
            },
            toggleEventList (classify) {
                classify.isHidden = !classify.isHidden
            },
            eventHeight (len) {
                return `height: ${len * 32}px`
            },
            setEventPushList () {
                const eventPushList = []
                const subscriptionForm = {}
                this.classifications.map(classify => {
                    if (classify['bk_objects'].length && classify['bk_classification_id'] !== 'bk_host_manage') {
                        const event = {
                            name: classify['bk_classification_name'],
                            isHidden: false,
                            classificationId: classify['bk_classification_id'],
                            children: []
                        }
                        if (classify['bk_classification_id'] === 'bk_biz_topo') {
                            event.children.push({
                                id: 'set',
                                name: this.$t('集群')
                            }, {
                                id: 'module',
                                name: this.$t('模块')
                            })
                            subscriptionForm['set'] = []
                            subscriptionForm['module'] = []
                        } else {
                            classify['bk_objects'].map(model => {
                                event.children.push({
                                    id: model['bk_obj_id'],
                                    name: model['bk_obj_name']
                                })
                                subscriptionForm[model['bk_obj_id']] = []
                            })
                        }
                        eventPushList.push(event)
                    }
                })
                subscriptionForm['resource'] = []
                subscriptionForm['host'] = []
                eventPushList.unshift({
                    isDefault: true,
                    isHidden: false,
                    classificationId: 'bk_host_manage',
                    name: this.$t('主机业务'),
                    children: [{
                        id: 'resource',
                        name: this.$t('主机池')
                    }, {
                        id: 'host',
                        name: this.$t('主机')
                    }]
                })
                this.$set(this.tempEventData, 'subscription_form', subscriptionForm)
                this.eventPushList = eventPushList
            },
            initData () {
                if (this.type === 'edit') {
                    const subscriptionForm = {}
                    this.curPush['subscription_form'].map(item => {
                        switch (item) {
                            case 'hostcreate':
                                if (subscriptionForm.hasOwnProperty('resource')) {
                                    subscriptionForm['resource'].push('hostcreate')
                                } else {
                                    subscriptionForm['resource'] = ['hostcreate']
                                }
                                break
                            case 'hostdelete':
                                if (subscriptionForm.hasOwnProperty('resource')) {
                                    subscriptionForm['resource'].push('hostdelete')
                                } else {
                                    subscriptionForm['resource'] = ['hostdelete']
                                }
                                break
                            case 'hostidentifier':
                                if (subscriptionForm.hasOwnProperty('host')) {
                                    subscriptionForm['host'].push('hostidentifier')
                                } else {
                                    subscriptionForm['host'] = ['hostidentifier']
                                }
                                break
                            case 'moduletransfer':
                                if (subscriptionForm.hasOwnProperty('host')) {
                                    subscriptionForm['host'].push('moduletransfer')
                                } else {
                                    subscriptionForm['host'] = ['moduletransfer']
                                }
                                break
                            default:
                                const key = item.substr(0, item.length - 6)
                                if (subscriptionForm.hasOwnProperty(key)) {
                                    subscriptionForm[key].push(item)
                                } else {
                                    subscriptionForm[key] = [item]
                                }
                        }
                    })
                    this.tempEventData = {
                        subscription_id: this.curPush['subscription_id'],
                        subscription_name: this.curPush['subscription_name'],
                        system_name: this.curPush['system_name'],
                        callback_url: this.curPush['callback_url'],
                        confirm_mode: this.curPush['confirm_mode'],
                        confirm_pattern: {
                            httpstatus: this.curPush['confirm_mode'] === 'httpstatus' ? this.curPush['confirm_pattern'] : '',
                            regular: this.curPush['confirm_mode'] === 'regular' ? this.curPush['confirm_pattern'] : ''
                        },
                        subscription_form: { ...this.tempEventData['subscription_form'], ...subscriptionForm },
                        time_out: this.curPush['time_out']
                    }
                }
                this.eventData = this.$tools.clone(this.tempEventData)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .slide-enter-active, .slide-leave-active{
        transition: height .3s;
        overflow: hidden;
    }
    .slide-enter, .slide-leave-to{
        height: 0 !important;
    }
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
            height: calc(100% - 51px);
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
                @include ellipsis;
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
                &.url {
                    font-size: 0;
                    .url-box {
                        display: inline-block;
                        width: calc(100% - 106px);
                    }
                    .test-btn {
                        vertical-align: top;
                        margin-left: 10px;
                        width: 96px;
                    }
                }
                span {
                    font-size: 14px;
                }
                .input-box {
                    position: relative;
                    .cmdb-form-input:focus {
                        +.tip {
                            border-top: 1px solid $cmdbBorderFocusColor;
                            border-right: 1px solid $cmdbBorderFocusColor;
                        }
                    }
                    .tip {
                        position: absolute;
                        display: inline-block;
                        left: 2px;
                        top: -4px;
                        width: 8px;
                        height: 8px;
                        background: #fff;
                        border-top: 1px solid $cmdbBorderColor;
                        border-right: 1px solid $cmdbBorderColor;
                        transform: rotate(-45deg);
                        &.reg {
                            left: 121px;
                        }
                    }
                }
                .cmdb-form-radio{
                    cursor: pointer;
                    height: 32px;
                    margin-bottom: 10px;
                    input[type="radio"]{
                        margin-top: -2px;
                    }
                }
                &.length-short{
                    position: relative;
                    .cmdb-form-input{
                        width: 97px;
                        margin-right: 10px;
                    }
                    .unit{
                        position: absolute;
                        left: 110px;
                        top: 8px;
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
                    font-size: 20px;
                    font-weight: bold;
                    transition: all .5s;
                    vertical-align: top;
                    &.up{
                        transform: rotate(-90deg);
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
            padding: 9px 20px;
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fafbfd;
            font-size: 0;
            .btn{
                margin-right: 10px;
            }
        }
    }
</style>
