/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <div class="tab-content">
        <form class="bk-form" id="validate_form">
            <div class="bk-form-item">
                <div class="form-common-item">
                    <label class="form-common-label">{{$t('ModelManagement["图标选择"]')}}<span class=""> * </span></label>
                    <div class="form-contain">
                        <div class="select-icon-content">
                            <div class="select-icon-show" @click.stop.prevent="toggleDrop" :class="{'active':isIconDrop}">
                                <div class="icon-content" >
                                    <i :class="baseInfo['bk_obj_icon']"></i>
                                </div>
                                <span class="arrow"><i class="bk-icon icon-angle-down"></i></span>
                            </div>
                            <div class="select-icon-mask" v-show="isIconDrop" @click="closeDrop"></div>
                            <div class="select-icon-list" v-show="isIconDrop">
                                <ul class="clearfix icon-list">
                                    <li v-tooltip="{content: language === 'zh-cn' ? item.nameZh : item.nameEn}" v-for="(item,index) in curIconList" :class="{'active': item.value === baseInfo['bk_obj_icon']}" @click.stop.prevent="chooseIcon(index, item)">
                                        <i :class="item.value"></i>
                                    </li>
                                </ul>
                                <div class="page-wrapper clearfix">
                                    <div class="input-wrapper">
                                        <input type="text" v-model="icon.searchText" :placeholder="$t('ModelManagement[\'请输入关键词\']')">
                                        <i class="bk-icon icon-search"></i>
                                    </div>
                                    <ul class="clearfix page">
                                        <li v-for="page in icon.totalPage"
                                        class="page-item" :class="{'cur-page': icon.curPage === page}"
                                        @click="icon.curPage = page"
                                        >
                                            {{page}}
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="form-common-item">
                    <label class="form-common-label">{{$t('ModelManagement["中文名称"]')}}<span class=""> * </span></label>
                    <div class="form-common-content interior-width-control">
                        <input type="text" name="" value=""
                            maxlength="20"
                            :disabled="isReadOnly || baseInfo.ispre"
                            name="validation_name"
                            :placeholder="$t('ModelManagement[\'请填写模型名\']')"
                            :data-parsley-required="true"
                            :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                            data-parsley-maxlength="20"
                            data-parsley-trigger="input blur"
                            v-model.trim="baseInfo['bk_obj_name']">
                    </div>
                </div>
                <div class="form-common-item">
                    <label class="form-common-label">{{$t('ModelManagement["英文名称"]')}}<span class=""> * </span></label>
                    <div class="form-common-content interior-width-control">
                        <input type="text" name="" value=""
                            :placeholder="$t('ModelManagement[\'下划线，数字，英文小写的组合\']')"
                            :class="{'is-danger': !baseInfoVerify['bk_obj_id']}" :disabled="type==='change'"
                            :data-parsley-required="true"
                            maxlength="20"
                            :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                            data-parsley-maxlength="20"
                            data-parsley-pattern="[a-z\d_]+"
                            :data-parsley-pattern-message="$t('ModelManagement[\'格式不正确，只能包含下划线，数字，英文小写\']')"
                            data-parsley-trigger="input blur"
                            v-model.trim="baseInfo['bk_obj_id']">
                    </div>
                </div>
            </div>
        </form>
        <div class="base-info" v-if="!isReadOnly">
            <bk-button type="primary" @click="saveBaseInfo" :loading="$loading('saveBaseInfo')">{{$t('Common["确定"]')}}</bk-button>
            <button class="btn vice-btn cancel-btn-sider" type="default" :title="$t('Common[\'取消\']')" @click="cancel">{{$t('Common["取消"]')}}</button>
        </div>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    import Parsley from 'parsleyjs'
    import '@/common/js/parsley_locale'
    import {mapGetters} from 'vuex'
    const iconList = require('@/common/json/modelIcon.json')
    export default {
        filters: {
            formatTime (val) {
                if (val) {
                    let res = val.replace('T', ' ')
                    res = res.replace('Z', '')
                    return res
                }
            }
        },
        props: {
            /*
                是否创建主线模型
            */
            isMainLine: {
                default: false
            },
            /*
                创建主线模型时才有值 为父对象模型的ID
            */
            associationId: {
                default: ''
            },
            isReadOnly: {
                default: false,
                type: Boolean
            },
            classificationId: {
                default: 0
            },
            /*
                操作类型 编辑 or 新增
                new 新增 change 编辑
            */
            type: {
                default: 'new'
            },
            isShow: {
                type: Boolean,
                default: false
            },
            objId: {
                default: ''
            }
        },
        data () {
            return {
                baseInfo: {
                    bk_obj_name: '',              // 模型名
                    bk_obj_id: '',                // API标识
                    bk_classification_id: '',
                    bk_supplier_account: 0,
                    bk_obj_icon: 'icon-cc-default'
                },
                baseInfoCopy: {},
                baseInfoVerify: {
                    bk_obj_name: true,             // true: 成功 false 失败
                    bk_obj_id: true
                },
                nowIndex: 0,                   // 选择图标下拉框当前index
                isIconDrop: false,             // 选择图标下拉框
                isChoose: true,                // 判断编辑分类的时候是否选择了icon
                iconValue: 'icon-cc-default', // 选择icon的值
                icon: {
                    searchText: '',
                    list: [],
                    count: 0,
                    curPage: 1,
                    totalPage: 0,
                    size: 24
                }
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'language'
            ]),
            curIconList () {
                let list = this.icon.list
                if (this.icon.searchText.length) {
                    list = this.icon.list.filter(icon => {
                        return icon.nameZh.toLowerCase().indexOf(this.icon.searchText.toLowerCase()) > -1 || icon.nameEn.toLowerCase().indexOf(this.icon.searchText.toLowerCase()) > -1
                    })
                }
                this.icon.count = list.length
                this.icon.totalPage = Math.ceil(list.length / this.icon.size)
                return list.slice((this.icon.curPage - 1) * this.icon.size, this.icon.curPage * this.icon.size)
            }
        },
        watch: {
            isShow (val) {
                if (val) {
                    $('#validate_form').parsley().reset()
                    this.isIconDrop = false
                    if (this.type === 'new') {
                        this.clearData()
                    } else {
                        this.getBaseInfo(this.objId)
                    }
                }
            },
            'icon.searchText' () {
                this.icon.curPage = 1
            }
        },
        methods: {
            isCloseConfirmShow () {
                if (this.type === 'new') {
                    if (this.baseInfo['bk_obj_id'] !== '' || this.baseInfo['bk_obj_name'] !== '' || this.baseInfo['bk_obj_icon'] !== 'icon-cc-default') {
                        return true
                    }
                } else {
                    if (this.baseInfo['bk_obj_name'] !== this.baseInfoCopy['bk_obj_name'] || this.baseInfo['bk_obj_icon'] !== this.baseInfoCopy['bk_obj_icon']) {
                        return true
                    }
                }
                return false
            },
            /*
                点击出现选中icon下拉框
            */
            toggleDrop () {
                this.isIconDrop = !this.isIconDrop
            },
            closeDrop () {
                this.isIconDrop = false
            },
            /*
                选中的icon
            */
            chooseIcon (index, item) {
                this.nowIndex = index
                this.isIconDrop = false
                this.iconValue = item.value
                this.isChoose = false
                this.baseInfo['bk_obj_icon'] = item.value
            },
            /*
                取消按钮
            */
            cancel () {
                this.clearData()
                this.$emit('cancel')
            },
            clearData () {
                this.baseInfo = {
                    bk_obj_name: '',
                    bk_obj_id: '',
                    bk_classification_id: '',
                    bk_supplier_account: this.bkSupplierAccount,
                    bk_obj_icon: 'icon-cc-default'
                }
                this.iconValue = 'icon-cc-default'
                this.baseInfoVerify = {
                    bk_obj_name: true,          // true: 成功 false 失败
                    bk_obj_id: true
                }
            },
            /*
                获取基本信息
            */
            getBaseInfo (ObjId) {
                let params = {
                    bk_obj_id: ObjId
                }
                this.$axios.post('objects', params).then(res => {
                    if (res.result) {
                        this.baseInfo = res.data[0]
                        this.baseInfoCopy = this.$deepClone(res.data[0])
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                保存基本信息
            */
            saveBaseInfo () {
                $('#validate_form').parsley().validate()
                if (!$('#validate_form').parsley().isValid()) return
                let params = {
                    bk_creator: window.userName,
                    bk_modifier: window.userName,
                    bk_classification_id: this.classificationId,
                    bk_obj_name: this.baseInfo['bk_obj_name'],
                    bk_supplier_account: this.bkSupplierAccount,
                    bk_obj_icon: this.iconValue
                }
                if (this.type === 'new') {
                    params['bk_obj_id'] = this.baseInfo['bk_obj_id']
                    if (this.isMainLine) { // 创建主线模型
                        this.$axios.post('topo/model/mainline', {
                            bk_classification_id: this.classificationId,
                            bk_obj_id: this.baseInfo['bk_obj_id'],
                            bk_obj_name: this.baseInfo['bk_obj_name'],
                            bk_supplier_account: this.bkSupplierAccount,
                            bk_asst_obj_id: this.associationId,
                            bk_obj_icon: this.iconValue
                        }, {id: 'saveBaseInfo'}).then(res => {
                            if (res.result) {
                                this.$emit('baseInfoSuccess', {
                                    bk_obj_name: this.baseInfo['bk_obj_name'],
                                    bk_obj_id: this.baseInfo['bk_obj_id'],
                                    id: res.data['id']
                                })
                            } else {
                                this.$alertMsg(this.$t('ModelManagement["创建模型失败"]'))
                            }
                        })
                    } else {
                        this.$axios.post('object', params, {id: 'saveBaseInfo'}).then(res => {
                            if (res.result) {
                                this.$emit('baseInfoSuccess', {
                                    bk_obj_name: this.baseInfo['bk_obj_name'],
                                    bk_obj_id: this.baseInfo['bk_obj_id'],
                                    id: res.data['id'],
                                    bk_obj_icon: this.iconValue
                                })
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                } else if (this.type === 'change') {
                    if (this.baseInfo['bk_obj_name'] === this.baseInfoCopy['bk_obj_name'] && this.baseInfo['bk_obj_icon'] === this.baseInfoCopy['bk_obj_icon']) {
                        this.cancel()
                    } else {
                        params['bk_ispre'] = this.baseInfo['bk_ispre']
                        this.$axios.put(`object/${this.baseInfo['id']}`, params, {id: 'saveBaseInfo'}).then(res => {
                            if (res.result) {
                                this.$alertMsg(this.$t('ModelManagement["修改成功"]'), 'success')
                                this.$emit('confirm', {
                                    bk_obj_name: this.baseInfo['bk_obj_name'],
                                    bk_obj_id: this.baseInfo['bk_obj_id']
                                })
                                this.$store.commit('navigation/updateModel', {
                                    bk_classification_id: this.classificationId,
                                    bk_obj_id: this.baseInfo['bk_obj_id'],
                                    bk_obj_name: this.baseInfo['bk_obj_name']
                                })
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                }
            }
        },
        mounted () {
            this.list = iconList
            this.icon = {
                list: iconList,
                count: this.list.length,
                searchText: '',
                curPage: 1,
                totalPage: Math.ceil(this.list.length / 24),
                size: 24
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    $primaryHoverColor: #6b7baa; // 主要颜色
    $borderColor: #e7e9ef; //边框色
    .tab-content{
        .bk-form{
            padding: 20px 0;
            // display: table;
            .input{
                width:425px;
            }
            .bk-form-item{
                display: table;
                width: 100%;
            }
            .form-common-item{
                position: relative;
                margin-right: 46px;
                float: left;
                font-size: 0;
                &:last-child{
                    margin-right: 0;
                }
                .form-common-label{
                    display: inline-block;
                    width: 70px;
                    vertical-align: top;
                    text-align:right;
                    font-size: 14px;
                    line-height: 36px;
                    span{
                        display: inline-block;
                        color: #f05d5d;
                        padding-left: 3px;
                    }
                }
                .form-common-content{
                    width:158px;
                    height: 36px;
                    display: inline-block;
                    padding-left: 8px;
                    input{
                        width:100%;
                        height: 100%;
                        padding: 0 13px;
                        color: #6b7baa;
                        font-size: 12px;
                    }
                }
                .selcet-width-control{
                    text-align:left;
                }
            }
            .bk-label{
                width:95px;
                color:$primaryHoverColor;
                padding: 10px 12px 10px 0;
            }
            .bk-form-content{
                margin-left:95px;
                .text{
                    display:inline-block;
                    height: 36px;
                    line-height: 1;
                    color: #666;
                    background-color: #fff;
                    border-radius: 2px;
                    width: 100%;
                    box-sizing: border-box;
                    padding: 0 10px;
                    font-size: 14px;
                    vertical-align: middle;
                }
            }
        }
        .bk-form-input, .bk-form-password, .bk-form-select, .bk-form-textarea{
            border:1px solid $borderColor;
            border-radius:2px;
        }
        .base-info{
            width: 100%;
            // margin-left: 97px;
            position: absolute;
            left: 0;
            bottom: 0;
            padding: 14px 10px;
            background: #f9f9f9;
            button{
                height: 36px;
                line-height: 34px;
                border-radius: 2px;
                display: inline-block;
                min-width: 110px;
                margin-left: 10px;
                vertical-align: bottom;
            }
        }
        .parsley-errors-list,.parsley-errors-list{
            margin-top: 8px;
            padding-left: 10px;
        }
    }
    .form-contain{
        display: inline-block;
        margin-left: 8px;
        .bk-select{
            width: 100px !important;
        }
        .from-input{
            width:270px;
            height: 40px;
            line-height: 38px;
            border: 1px solid $borderColor;
            outline: none;
            padding: 0 15px;
        }
        .select-icon-content{ /* 图标选择下拉 */
            display: inline-block;
            position: relative;
            border: 1px solid #c3cdd7;
            .select-icon-mask{
                position: fixed;
                top: 0;
                bottom: 0;
                left: 0;
                right: 0;
            }
            .select-icon-show{
                height: 34px;
                line-height: 34px;
                /*width: 70px;*/
                cursor: pointer;
                &.active{
                    border-bottom: 0;
                }
                .icon-content{
                    // width: 37px;
                    padding-left: 15px;
                    padding-right: 6px;
                    line-height: 36px;
                    height: 100%;
                    // border-right: 1px solid $borderColor;
                    text-align: center;
                    >i{
                        color: #3c96ff;
                        vertical-align: middle;
                        font-size: 24px;
                    }
                    float: left;
                }
                .arrow{
                    float: left;
                    width: 26px;
                    vertical-align: bottom;
                    i{
                        color: #6b7baa;
                        font-size: 12px;
                        padding-left: 2px;
                    }
                }
            }
            .select-icon-list{
                padding: 10px;
                position: absolute;
                top: 44px;
                left: 0;
                width: 382px;
                height: 248px;
                border: 1px solid #bec6de;
                z-index: 500;
                background: #fff;
                box-shadow: 0 2px 2px rgba(0,0,0,.1);
                overflow: auto;
                @include scrollbar;
                .icon-list{
                    padding: 0;
                    margin: 0;
                    width: 360px;
                    height: 184px;
                    li{
                        width: 60px;
                        height: 46px;
                        text-align: center;
                        line-height: 46px;
                        float: left;
                        cursor: pointer;
                        &.active{
                            color: #3c96ff;
                            background: #e2efff;
                        }
                        i{
                            font-size: 24px;
                        }
                        &:hover{
                            background: #e2efff;
                        }
                        &:nth-child(6n){
                            margin-right: 0;
                        }
                    }
                }
                .page-wrapper {
                    padding: 15px 18px 5px;
                    .input-wrapper {
                        float: left;
                        position: relative;
                        vertical-align: bottom;
                        font-size: 12px;
                        color: #c3cdd7;
                        input {
                            width: 116px;
                            height: 22px;
                            padding: 0 25px 0 5px;
                            border: 1px solid #c3cdd7;
                            border-radius: 2px;
                        }
                        .bk-icon {
                            position: absolute;
                            top: 5px;
                            right: 8px;
                        }
                    }
                }
                .page{
                    float: right;
                    li{
                        text-align: center;
                        float: left;
                        margin-right: 5px;
                        width: 22px;
                        height: 22px;
                        line-height: 20px;
                        border-radius: 2px;
                        font-size: 12px;
                        cursor: pointer;
                        color: #737987;
                        border: 1px solid #c3cdd7;
                        &.cur-page{
                            color: #fff;
                            background: #3c96ff;
                            border-color: #3c96ff;
                        }
                        &:last-child{
                            margin: 0;
                        }
                    }
                }
            }
        }
    }
</style>
