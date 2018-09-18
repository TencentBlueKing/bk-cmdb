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
    <bk-dialog
        :is-show.sync="isShow" 
        :has-header="false" 
        :has-footer="false" 
        :quick-close="false" 
        :width="565" 
        :padding="0">
        <div class="pop-box" slot="content">
            <div class="pop-info">
                <div class="title" v-if="!isEdit">{{$t('ModelManagement["新增分组"]')}}</div>
                <div class="title" v-else>{{$t('ModelManagement["编辑分组"]')}}</div>
                <div class="content">
                    <ul class="content-left">
                        <li class="content-item">
                            <label for="">{{$t('ModelManagement["中文名"]')}}<span class="color-danger">*</span></label>
                            <div class="input-box">
                                <input type="text" class="cmdb-form-input fr" 
                                v-focus
                                v-model.trim="localValue['bk_classification_name']"
                                :disabled="classification['bk_classification_type'] === 'inner'"
                                @blur="validate"
                                :data-vv-name="$t('Common[\'中文名\']')"
                                v-validate="'required|classifyName'">
                                <span v-show="errors.has($t('Common[\'中文名\']'))" class="error-msg color-danger">{{ errors.first($t('Common[\'中文名\']')) }}</span>
                            </div>
                        </li> 
                        <li class="content-item">
                            <label for="">{{$t('ModelManagement["英文名"]')}}<span class="color-danger">*</span></label>
                            <div class="input-box">
                                <input type="text" class="cmdb-form-input fr" v-model.trim="localValue['bk_classification_id']"
                                :data-vv-name="$t('ModelManagement[\'英文名\']')"
                                :disabled="classification['bk_classification_type'] === 'inner' || isEdit"
                                v-validate="'required|classifyId'">
                                <span v-show="errors.has($t('ModelManagement[\'英文名\']'))" class="error-msg color-danger">{{ errors.first($t('ModelManagement[\'英文名\']')) }}</span>
                            </div>
                        </li> 
                    </ul>
                    <div class="content-right" @click="isIconListShow = true">
                        <div class="icon-wrapper">
                            <i :class="localValue['bk_classification_icon']"></i>
                        </div>
                        <div class="text">{{$t('ModelManagement["图标选择"]')}}</div>
                    </div>
                </div>
                <div class="footer">
                    <div class="btn-group">
                        <bk-button type="primary" :loading="$loading('saveClassify')" class="confirm-btn" @click="saveClassify">{{$t('Common["确定"]')}}</bk-button>
                        <bk-button type="default" @click="cancel">{{$t('Common["取消"]')}}</bk-button>
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
                    <ul></ul>
                    <span class="info">{{$t('ModelManagement["单击选择对应图标"]')}}</span>
                </div>
            </div>
        </div>
    </bk-dialog>
</template>

<script>
    import iconList from '@/assets/json/class-icon.json'
    import { mapGetters, mapActions, mapMutations } from 'vuex'
    export default {
        props: {
            isEdit: {
                type: Boolean,
                default: false
            },
            classification: {
                type: Object,
                default: () => {
                    return {
                        bk_classification_icon: '',
                        bk_classification_name: '',
                        bk_classification_id: ''
                    }
                }
            }
        },
        data () {
            return {
                isShow: true,
                iconList,
                isIconListShow: false,
                localValue: JSON.parse(JSON.stringify(this.classification))
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            classificationId () {
                return this.$route.params.classifyId
            },
            activeClassify () {
                if (!this.isEdit) {
                    return {
                        bk_classification_id: '',
                        bk_classification_icon: '',
                        bk_classification_name: '',
                        bk_classification_type: ''
                    }
                }
                let activeClassify = this.classifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === this.classificationId)
                return activeClassify
            }
        },
        methods: {
            ...mapActions('objectModelClassify', [
                'createClassification',
                'updateClassification'
            ]),
            ...mapMutations('objectModelClassify', [
                'updateClassify'
            ]),
            async createClassify () {
                let params = {
                    bk_classification_icon: this.localValue['bk_classification_icon'],
                    bk_classification_id: this.localValue['bk_classification_id'],
                    bk_classification_name: this.localValue['bk_classification_name']
                }
                const res = await this.createClassification({params}).then(res => {
                    this.$http.cancel('post_searchClassificationsObjects')
                    return res
                })
                Object.assign(params, {bk_supplier_account: this.supplierAccount, id: res.id})
                this.updateClassify(params)
                this.$emit('closePop')
                this.$router.push(`/model/${this.localValue['bk_classification_id']}`)
            },
            async editClassify () {
                let params = {
                    bk_classification_icon: this.localValue['bk_classification_icon'],
                    bk_classification_name: this.localValue['bk_classification_name']
                }
                await this.updateClassification({
                    id: this.activeClassify['id'],
                    params
                }).then(() => {
                    this.$http.cancel('post_searchClassificationsObjects')
                })
                this.updateClassify({
                    ...params,
                    ...{
                        bk_classification_id: this.activeClassify['bk_classification_id']
                    }
                })
                this.$emit('closePop')
            },
            saveClassify () {
                if (this.isEdit) {
                    this.editClassify()
                } else {
                    this.createClassify()
                }
            },
            validate () {
                this.$validator.validateAll()
            },
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
                this.isShow = false
                setTimeout(() => {
                    this.$emit('closePop')
                }, 300)
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
        created () {
            if (this.isEdit) {
                this.localValue = {
                    bk_classification_icon: this.activeClassify['bk_classification_icon'],
                    bk_classification_name: this.activeClassify['bk_classification_name'],
                    bk_classification_id: this.activeClassify['bk_classification_id']
                }
            } else {
                this.localValue = {
                    bk_classification_icon: 'icon-cc-default',
                    bk_classification_name: '',
                    bk_classification_id: ''
                }
            }
        },
        directives: {
            focus: {
                inserted: function (el) {
                    el.focus()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .error-danger{
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
                        .input-box {
                            width: 259px;
                            float: left;
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
                        .cmdb-form-input{
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
                        width: 100%;
                        height: 63px;
                        font-size: 38px;
                        text-align: center;
                        color: #63abff;
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
                height: 236px;
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
</style>
