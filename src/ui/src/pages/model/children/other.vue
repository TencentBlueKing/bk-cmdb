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
    <div class="tab-content other-content">
        <div class="other-list">
            <h3>{{$t('ModelManagement["模型停用"]')}}</h3>
            <p><span v-if="!item['bk_ispaused']">{{$t('ModelManagement["保留模型和相应实例，隐藏关联关系"]')}}</span></p>
            <div class="bottom-contain">
                <bk-button type="primary" :loading="$loading('restartModel')" :disabled="$loading('restartModel')" :title="$t('ModelManagement[\'启用模型\']')" v-if="isReadOnly" class="bk-button main-btn mr10 button-on" @click="restartModelConfirm">
                    {{$t('ModelManagement["启用模型"]')}}
                </bk-button>
                <bk-button type="primary" :loading="$loading('stopModel')" v-else class="mr10" :title="$t('ModelManagement[\'停用模型\']')" @click="showConfirmDialog('stop')" :class="['bk-button bk-default', {'is-disabled': item['ispre'] || parentClassificationId === 'bk_biz_topo'}]" :disabled="item['ispre'] || parentClassificationId === 'bk_biz_topo'">
                    {{$t('ModelManagement["停用模型"]')}}
                </bk-button>
                <span class="btn-tip-content" v-show="isShowTipStop=item['ispre']">
                    <i class="icon-cc-attribute"></i>
                    <span class="btn-tip" :class="{'en': language === 'en'}">
                        <i class="right-triangle"></i>
                        <i class="left-triangle"></i>
                        {{$t('ModelManagement["系统内建模型不可停用"]')}}
                    </span>
                </span>
            </div>
        </div>
        <div class="other-list mt50">
            <h3>{{$t('ModelManagement["模型删除"]')}}</h3>
            <p>{{$t('ModelManagement["删除模型和其下所有实例，此动作不可逆，请谨慎操作"]')}}</p>
            <div class="bottom-contain">
                <bk-button type="primary" :loading="$loading('deleteModel')" class="mr10" :title="$t('ModelManagement[\'删除模型\']')" @click="showConfirmDialog('delete')" :class="['bk-button bk-default', {'is-disabled':item['ispre']}]" :disabled="item['ispre']">
                    <span>{{$t('ModelManagement["删除模型"]')}}</span>
                </bk-button>
                <span class="btn-tip-content" v-show="isShowTipStop=item['ispre']">
                    <i class="icon-cc-attribute"></i>
                    <span class="btn-tip" :class="{'en': language === 'en'}">
                        <i class="right-triangle"></i>
                        <i class="left-triangle"></i>
                        {{$t('ModelManagement["系统内建模型不可删除"]')}}
                    </span>
                </span>
            </div>
        </div>
    </div>
</template>

<script type="text/javascript">
    import {mapGetters} from 'vuex'
    export default {
        props: {
            isMainLine: {
                default: false
            },
            isReadOnly: {
                default: false,
                type: Boolean
            },
            id: {
                default: 0
            },
            item: {
                default () {
                    return {}
                },
                type: Object
            },
            parentClassificationId: ''
        },
        data () {
            return {
                isShowTipStop: false
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'language'
            ])
        },
        methods: {
            /*
                重新启用模型确认弹框
            */
            restartModelConfirm () {
                let self = this
                this.$bkInfo({
                    title: this.$t('ModelManagement["确认要启用该模型？"]'),
                    confirmFn () {
                        self.restartModel()
                    }
                })
            },
            /*
                重新启用模型
            */
            restartModel () {
                let params = {
                    bk_ispaused: false
                }
                this.$axios.put(`object/${this.id}`, params, {id: 'restartModel'}).then(res => {
                    if (res.result) {
                        this.$emit('closeSideSlider')
                        this.$store.dispatch('navigation/getClassifications', true)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
               停用模型
            */
            stopModel () {
                let params = {
                    bk_ispaused: true
                }
                this.$axios.put(`object/${this.id}`, params, {id: 'stopModel'}).then(res => {
                    if (res.result) {
                        this.$emit('stopModel')
                        this.$store.dispatch('navigation/getClassifications', true)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
               删除模型
            */
            deleteModel () {
                if (this.isMainLine) {
                    this.$axios.delete(`topo/model/mainline/owners/${this.bkSupplierAccount}/objectids/${this.item['bk_obj_id']}`, {id: 'deleteModel'}).then(res => {
                        if (res.result) {
                            this.$emit('deleteModel', this.item)
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else {
                    this.$axios.delete(`object/${this.id}`, {id: 'deleteModel'}).then(res => {
                        if (res.result) {
                            this.$emit('deleteModel', this.item)
                            this.$store.dispatch('navigation/getClassifications', true)
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            /*
               确认弹框 确认事件
            */
            dialogConfirm () {
                this.confirmInfo.isConfirmShow = false
                if (this.confirmInfo.confirmType === 'stop') {
                    this.confirmInfo.text = this.$t('ModelManagement["确认要停用该模型？"]')
                    this.stopModel()
                } else if (this.confirmInfo.confirmType === 'delete') {
                    this.confirmInfo.text = this.$t('ModelManagement["确认要删除该模型？"]')
                    this.deleteModel()
                }
            },
            /*
               显示二次确认弹窗
               type: 类型
            */
            showConfirmDialog (type) {
                let self = this
                switch (type) {
                    case 'stop':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要停用该模型？"]'),
                            confirmFn () {
                                self.stopModel()
                            }
                        })
                        break
                    case 'delete':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要删除该模型？"]'),
                            confirmFn () {
                                self.deleteModel()
                            }
                        })
                        break
                }
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    .other-content{
        padding:54px 34px 0 34px;
        .other-list{
            >h3{
                font-size:14px;
                font-weight:bold;
                border-left:4px solid #4d597d;
                line-height:1;
                color:#4d597d;
                padding-left:5px;
                margin:0;
            }
            >p{
                line-height:1;
                margin-top:9px;
                margin-bottom:20px;
            }
        }
        .bottom-contain{
            .btn-tip-content{
                .icon-cc-attribute{
                    cursor: pointer;
                    color: #ffb400;
                    font-size: 16px;
                    +span{
                        display: none;
                    }
                    &:hover{
                        +span{
                            display: inline-block;
                        }
                    }
                }
                .btn-tip{
                    display:inline-block;
                    min-width:170px;
                    height:33px;
                    line-height:33px;
                    text-align:center;
                    box-shadow: 0 0 5px #ebedef;
                    margin-left: 8px;
                    position:relative;
                    background: #333333;
                    color: #fff;
                    border-radius: 2px;
                    font-size: 12px;
                    &.en{
                        min-width: 300px;
                    }
                    .left-triangle{
                        width: 0;
                        height: 0;
                        border-top: 7px solid transparent;
                        border-right: 7px solid #333;
                        border-bottom: 7px solid transparent;
                        position: absolute;
                        left: -7px;
                        top: 9px;
                    }
                }
            }
        }
        .is-disabled{
            cursor: not-allowed !important;
        }
    }
</style>
