/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */
<template>
    <div class="pop-wrapper" v-if="isShow">
        <div class="pop-box">
            <div class="pop-info">
                <div class="title">
                    <template v-if="!isEdit">新增分组</template>
                    <template v-else>编辑分组</template>
                </div>
                <div class="content">
                    <ul class="content-left">
                        <li class="content-item">
                            <label for="">中文名<span class="color-danger">*</span></label>
                            <input type="text" class="bk-form-input" 
                            v-focus
                            v-model.trim="localValue['bk_classification_name']"
                            data-vv-name="中文名"
                            v-validate="'required|name'">
                            <span v-show="errors.has('中文名')" class="help is-danger">{{ errors.first('中文名') }}</span>
                        </li> 
                        <li class="content-item">
                            <label for="">英文名<span class="color-danger">*</span></label>
                            <input type="text" class="bk-form-input" v-model="localValue['bk_classification_id']"
                            name="id"
                            :disabled="isEdit"
                            v-validate="'required|id'">
                            <span v-show="errors.has('id')" class="help is-danger">{{ errors.first('id') }}</span>
                        </li> 
                    </ul>
                    <div class="content-right" @click="isIconListShow = true">
                        <div class="icon-wrapper">
                            <i :class="localValue['bk_classification_icon']"></i>
                        </div>
                        <div class="text">图标选择</div>
                    </div>
                </div>
                <div class="footer">
                    <div class="btn-group">
                        <bk-button type="primary" class="confirm-btn" @click="confirm">确定</bk-button>
                        <bk-button type="default" @click="cancel">取消</bk-button>
                    </div>
                </div>
            </div>
            <div class="pop-icon-list" v-show="isIconListShow">
                <span class="back" @click="closeIconList">
                    <i class="bk-icon icon-back2"></i>
                </span>
                <ul class="icon-box clearfix">
                    <li class="icon" 
                        @click="chooseIcon(icon)"
                        v-for="(icon, index) in iconList" 
                        :class="{'active': localValue['bk_classification_icon'] === icon.value}"
                        :key="index">
                        <i :class="icon.value"></i>
                    </li>
                </ul>
                <div class="page">
                    <span class="info">单击选择对应图标</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import iconList from '@/common/json/modelIcon.json'
    export default {
        props: {
            /*
                弹窗显示状态
            */
            isShow: {
                type: Boolean,
                default: false
            },
            /*
                编辑时参数
            */
            classification: {
                type: Object,
                default: () => {
                    return {
                        bk_classification_icon: '',
                        bk_classification_name: '',
                        bk_classification_id: ''
                    }
                }
            },
            /*
                当前状态 编辑or新增
            */
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                isIconListShow: false,          // 图标弹窗
                iconList: [],
                localValue: {
                    bk_classification_icon: 'icon-cc-business',
                    bk_classification_name: '',
                    bk_classification_id: ''
                }
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    if (this.isEdit) {
                        this.localValue = {
                            bk_classification_icon: this.classification['bk_classification_icon'],
                            bk_classification_name: this.classification['bk_classification_name'],
                            bk_classification_id: this.classification['bk_classification_id']
                        }
                    } else {
                        this.localValue = {
                            bk_classification_icon: 'icon-cc-business',
                            bk_classification_name: '',
                            bk_classification_id: ''
                        }
                    }
                }
            }
        },
        methods: {
            /*
                确认按钮
            */
            confirm () {
                this.$validator.validateAll().then(res => {
                    if (res) {
                        this.$emit('confirm', this.localValue)
                    }
                })
            },
            /*
                取消
            */
            cancel () {
                this.$emit('update:isShow', false)
            },
            /*
                关闭选择图标弹窗
            */
            closeIconList () {
                this.isIconListShow = false
            },
            /*
                选择图标
            */
            chooseIcon (item) {
                this.localValue['bk_classification_icon'] = item.value
                this.closeIconList()
            }
        },
        directives: {
            focus: {
                inserted: function (el) {
                    el.focus()
                }
            }
        },
        mounted () {
            this.iconList = iconList
        }
    }
</script>

<style lang="scss" scoped>
    .pop-wrapper{
        position: fixed;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        background: rgba(0, 0, 0, .6);
        z-index: 2000;
        .is-danger{
            color: #ff5656;
            font-size: 12px;
            margin-left: 58px;
        }
        .pop-box{
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 565px;
            height: 310px;
            border-radius: 2px;
            background-color: #fff;
            box-shadow: 0px 3px 7px 0px rgba(0, 0, 0, 0.1);
            .pop-info{
                .title{
                    margin: 50px auto 40px;
                    text-align: center;
                    font-size: 18px;
                    color: #333948;
                    line-height: 1;
                }
                .content{
                    height: 92px;
                    .content-left{
                        float: left;
                        width: 388px;
                        padding-left: 70px;
                        line-height: 1;
                        .content-item{
                            height: 36px;
                            &:first-child{
                                margin-bottom: 20px;
                            }
                            label{
                                line-height: 36px;
                                float: left;
                                margin-right: 8px;
                                font-size: 14px;
                                .color-danger{
                                    color: #ff5656;
                                }
                            }
                            .bk-form-input{
                                width: 259px;
                                vertical-align: baseline;
                            }
                        }
                    }
                    .content-right{
                        width: 88px;
                        height: 92px;
                        float: right;
                        margin-right: 69px;
                        border: 1px solid #c3cdd7;
                        border-radius: 2px;
                        cursor: pointer;
                        .icon-wrapper{
                            padding-top: 4px;
                            height: 63px;
                            font-size: 38px;
                            text-align: center;
                            color: #c7d0d9;
                        }
                        .text{
                            height: 27px;
                            color: #737987;
                            background: #f9fafb;
                            font-size: 12px;
                            line-height: 27px;
                            text-align: center;
                        }
                    }
                }
                .footer{
                    padding: 12px 18px;
                    margin-top: 50px;
                    height: 60px;
                    border-top: 1px solid #e5e5e5;
                    background: #fafbfd;
                    text-align: right;
                    font-size: 0;
                    .btn-group{
                        .confirm-btn{
                            margin-right: 10px;
                        }
                    }
                }
            }
            .pop-icon-list{
                position: absolute;
                width: 565px;
                height: 310px;
                background: #fff;
                top: 0;
                left: 0;
                padding: 20px 13px 0;
                .icon-box{
                    .icon{
                        float: left;
                        width: 77px;
                        height: 49px;
                        padding: 5px;
                        font-size: 30px;
                        text-align: center;
                        margin-bottom: 10px;
                        cursor: pointer;
                        &:hover,
                        &.active{
                            background: #e2efff;
                            color: #3c96ff;
                        }
                    }
                }
                .back{
                    position: absolute;
                    right: -47px;
                    top: 0;
                    width: 44px;
                    height: 44px;
                    padding: 7px;
                    cursor: pointer;
                    font-size: 18px;
                    text-align: center;
                    background: #2f2f2f;
                    color: #fff;
                }
                .page{
                    height: 52px;
                    .info{
                        padding-right: 25px;
                        line-height: 52px;
                        float: right;
                        color: #c3cdd7;
                        font-size: 16px;
                    }
                }
            }
        }
    }
</style>
