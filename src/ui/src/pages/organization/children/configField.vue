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
    <div class="side-slider-wrapper">
        <div class="content">
            <div class="left-list search-wrapper-hidden">
                <div class="search-wrapper">
                    <div class="the-host">
                        {{$t("Inst['隐藏属性']")}}
                    </div>
                    <div class="search-field">
                        <input type="text" name="" value="" :placeholder="$t('Inst[\'搜索属性\']')" v-model.trim="searchText">
                    </div>
                </div>
                <ul class="list-wrapper">
                    <li v-for="(item, index) in shownList" 
                        @click="addItem(item, index)"
                        :key="index">
                        <span :title="item['bk_property_name']">{{item['bk_property_name']}}</span>
                        <i class="bk-icon icon-angle-right"></i>
                    </li>
                </ul>
            </div>
            <div class="right-list">
                <div class="title">
                    <div class="search-wrapper">
                        {{$t("Inst['已显示属性']")}}
                    </div>
                </div>
                <div :class="['model-content pr20', {'content-left-hidden' : isShow}]" >
                    <div slot="contentRight">
                        <draggable class="content-right" v-model="localHasSelectionList" :options="{animation: 150}">
                            <div v-for="(item, index) in localHasSelectionList" :key="index" class="item">
                                <i class="icon-triple-dot"></i><span :title="item['bk_property_name']">{{item['bk_property_name']}}</span><i class="bk-icon icon-eye-slash-shape" @click="removeItem(item, index)" v-tooltip="$t('Common[\'隐藏\']')"></i>
                            </div>
                        </draggable>
                    </div>
                </div>
            </div>
        </div>
        <div class="bk-form-item bk-form-action content-button">
            <bk-button class="btn apply" type="primary" @click="apply">
                {{$t("Inst['应用']")}}
            </bk-button>
            <bk-button class="btn reinstate vice-btn" type="default" @click="cancel">
                {{$t("Common['取消']")}}
            </bk-button>
            <bk-button class="fr vice-btn" type="default" @click="resetConfirm">
                {{$t("Common['还原默认']")}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import draggable from 'vuedraggable'
    import {mapGetters} from 'vuex'
    import { sortArray } from '@/utils/util'
    export default {
        data () {
            return {
                localHasSelectionList: [],          // 已显示属性 本地操作的结果 还未保存到服务端
                localForSelectionList: [],          // 隐藏属性 本地操作的结果 还未保存到服务端
                hasSelectionList: [],               // 已显示属性 服务端的结果
                forSelectionList: [],               // 隐藏属性 服务端的结果
                searchText: ''
            }
        },
        props: {
            objId: {
                default: ''
            },
            /*
                所有属性列表
            */
            attrList: {
                type: Array,
                default: () => {
                    return []
                }
            },
            /*
                弹窗显示状态
            */
            isShow: {
                type: Boolean,
                default: false
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'usercustom'
            ]),
            shownList () {
                let list = []
                this.localForSelectionList.map(val => {
                    if (val['bk_property_name'].indexOf(this.searchText) !== -1) {
                        list.push(val)
                    }
                })
                return sortArray(list, 'bk_property_name')
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.getUserAttr()
                } else {
                    this.searchText = ''
                }
            }
        },
        methods: {
            resetConfirm () {
                this.$bkInfo({
                    title: this.$t("Common['是否要还原回系统默认显示属性？']"),
                    confirmFn: () => {
                        this.resetFields()
                    }
                })
            },
            async resetFields () {
                try {
                    let params = {}
                    params[`${this.objId}DisplayColumn`] = []
                    await this.$axios.post('usercustom', params)
                    this.$emit('resetFields')
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            addItem (item, index) {
                this.localForSelectionList = this.localForSelectionList.filter(property => property['bk_property_id'] !== item['bk_property_id'])
                this.localHasSelectionList.push(item)
            },
            removeItem (item, index) {
                this.localHasSelectionList = this.localHasSelectionList.filter(property => property['bk_property_id'] !== item['bk_property_id'])
                this.localForSelectionList.push(item)
            },
            /*
                更新用户字段
            */
            updateUsercustom () {
                let usercustom = this.$deepClone(this.usercustom)
                usercustom[`${this.objId}DisplayColumn`] = this.hasSelectionList
                this.$store.commit('setUsercustom', usercustom)
            },
            /*
                取消
            */
            cancel () {
                this.$emit('cancel')
            },
            /*
                保存
            */
            async apply () {
                if (this.localHasSelectionList.length === 0) {
                    this.$alertMsg(this.$t("Inst['请至少选择一项']"), 'primary')
                    return
                }
                let params = {}
                params[`${this.objId}DisplayColumn`] = this.localHasSelectionList
                await this.$axios.post('usercustom', params).then(res => {
                    if (res.result) {
                        this.hasSelectionList = this.$deepClone(this.localHasSelectionList)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
                this.$emit('apply', this.$deepClone(this.hasSelectionList))
                this.updateUsercustom()
            },
            /*
                获取用户定义的字段
            */
            async getUserAttr () {
                if (JSON.stringify(this.usercustom) === '{}') {
                    await this.$axios.post('usercustom/user/search').then(res => {
                        if (res.result) {
                            this.$store.commit('setUsercustom', res.data)
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
                if (this.usercustom.hasOwnProperty(`${this.objId}DisplayColumn`) && this.usercustom[`${this.objId}DisplayColumn`].length) {
                    let selectedList = this.$deepClone(this.usercustom[`${this.objId}DisplayColumn`])
                    selectedList.map(list => {
                        let property = this.attrList.find(attr => {
                            return attr['bk_property_id'] === list['bk_property_id']
                        })
                        if (property) {
                            list['bk_property_name'] = property['bk_property_name']
                        }
                    })
                    this.localHasSelectionList = selectedList
                } else {
                    this.localHasSelectionList = []
                }
                this.hasSelectionList = this.$deepClone(this.localHasSelectionList)
                this.setForSelectionList()
            },
            /*
                获取未选列表
                从所有字段中去掉已选择的即为未选择的列表
            */
            setForSelectionList () {
                let localForSelectionList = []
                this.attrList.map(val => {
                    localForSelectionList.push({
                        bk_property_id: val['bk_property_id'],
                        bk_property_name: val['bk_property_name'],
                        bk_property_type: val['bk_property_type']
                    })
                })
                for (let i = localForSelectionList.length - 1; i >= 0; i--) {
                    let val = localForSelectionList[i]
                    for (let j = 0; j < this.localHasSelectionList.length; j++) {
                        if (val['bk_property_id'] === this.localHasSelectionList[j]['bk_property_id']) {
                            localForSelectionList.splice(i, 1)
                            break
                        }
                    }
                }
                this.localForSelectionList = this.$deepClone(localForSelectionList)
                this.forSelectionList = this.$deepClone(localForSelectionList)
            }
        },
        components: {
            draggable
        }
    }
</script>

<style lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    .side-slider-wrapper {
        height: 100%;
        .content{
            height: calc(100% - 122px);
            border-top: 1px solid #e7e9ef;
        }
        .left-list{
            float: left;
            width: 50%;
            height: 100%;
            border-right: 1px solid #e7e9ef;
            .list-wrapper{
                height: calc(100% - 78px);
                overflow: auto;
                padding: 15px 0 0 0;
                &::-webkit-scrollbar{
                    width: 6px;
                    height: 5px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 20px;
                    background: #a5a5a5;
                }
            }
            &.search-wrapper-hidden{
                .search-wrapper{
                    .text{
                        display:none;
                    }
                    .search{
                        display:none;
                    }
                    .the-host{
                        width:122px;
                        display:inline-block;
                        line-height: 36px;
                    }
                    .search-field{
                        width:120px;
                        // display:inline-block;
                        float: right;
                        input{
                            border:1px solid #e7e9ef;
                            width:100%;
                            height:36px;
                            line-height:36px;
                            outline:none;
                            padding:0 15px;
                        }

                    }
                }
            }
            .search-wrapper{
                width: 100%;
                height: 78px;
                padding: 20px;
                .select-box{
                    float: left;
                    width: 163px;
                    height: 37px;
                    margin-right: 10px;
                    &.open{
                        .bk-selector-icon{
                            top: 17px;
                        }
                    }
                }
                .search{
                    float: left;
                    input{
                        width: 131px;
                        height: 37px;
                        border: 1px solid #e7e9ef;
                        border-radius: 2px;
                        padding: 0 12px;
                        font-size: 14px;
                        color: #bec6de;
                    }
                    &.search2{
                        float: right;
                        input{
                            width: 180px;
                        }
                    }
                }
                .text{
                    float: left;
                    line-height: 37px;
                    margin-left: 9px;
                }
            }
            .list-wrapper{
                border-top: 1px solid #e7e9ef;
                li{
                    height: 42px;
                    line-height: 42px;
                    color: $primaryColor;
                    font-size: 14px;
                    padding-left: 27px;
                    cursor: pointer;
                    &:hover{
                        background: #f9f9f9;
                    }
                    span{
                        display: inline-block;
                        width: 230px;
                        @include ellipsis;
                    }
                    i{
                        float: right;
                        margin-top: 12px;
                        margin-right: 18px;
                        color: #bec6de;
                    }
                }
            }
        }
        .right-list{
            float: left;
            height: 100%;
            width: 50%;
            .bk-tab2{
                border: none !important;
            }
            .content-left-hidden{
                .content-left{
                    display:none;
                }
                .content-right{
                    width:100% !important;
                    padding: 15px 0 0 0;
                    .content-center{
                        width:232px;
                            margin-right:36px;
                        .input-number{
                            width:100%;
                            height:36px;
                            line-height:36px;
                            outline:none;
                            border:1px solid #e7e9ef;
                            background:#f9f9f9;
                            color:#bec6de;
                            padding:0 15px;
                            &.disbale{
                                cursor:not-allowed;
                            }
                            &::-webkit-input-placeholder{
                                font-family: "Microsoft YaHei";
                                color: #c3cdd7;
                            }
                            &:-moz-placeholder{
                                font-family: "Microsoft YaHei";
                                color: #c3cdd7;
                            }
                            &::-moz-placeholder{
                                font-family: "Microsoft YaHei";
                                color: #c3cdd7;
                            }
                            &:-ms-input-placeholder{
                                font-family: "Microsoft YaHei";
                                color: #c3cdd7;
                            }
                        }
                    }
                }
            }
            .title{
                height: 79px;
                line-height: 78px;
                padding: 0 20px;
                width: 430px;
                border-bottom: 1px solid #e7e9ef;
                .list-wrapper{
                    border-top: 1px solid #e7e9ef;
                }
            }
            .model-content{
                color: $primaryColor;
                height: calc(100% - 79px);
                overflow: auto;
                &::-webkit-scrollbar{
                width: 6px;
                height: 5px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 20px;
                    background: #a5a5a5;
                }
                .content-left{
                    float: left;
                    width: 108px;
                    height: 100%;
                    text-align: center;
                    background: #f9f9f9;
                    // padding: 15px 0 0 0;
                    li{
                        height: 43px;
                        line-height: 42px;
                        border-bottom: 1px solid #fff;
                        &.active{
                            background: #fff;
                        }
                        &.add{
                            cursor: pointer;
                            .plus{
                                width:55px;
                                border-right: 1px solid #fff;
                            }
                            .edit{
                                width:53px;
                                text-align:center;
                            }
                            i{
                                color:#bec6de;
                            }
                        }
                    }
                }
                .content-right{
                    // float: left;
                    width: calc(100% - 108px);
                    >.item{
                        height: 43px;
                        line-height: 42px;
                        padding-left: 30px;
                        cursor: move;
                        &:hover{
                            background: #f9f9f9;
                        }
                        span{
                            display: inline-block;
                            width: calc(100% - 46px);
                            @include ellipsis;
                            vertical-align: bottom;
                        }
                        .icon-triple-dot{
                            position: relative;
                            top: -1px;
                            display: inline-block;
                            width: 4px;
                            height: 14px;
                            margin-right: 10px;
                            background: url(../../../common/images/icon/icon-triple-dot.png);
                        }
                        .icon-eye-slash-shape{
                            float: right;
                            font-size: 12px;
                            margin-top: 5px;
                            // margin-right: 20px;
                            padding: 10px;
                            cursor: pointer;
                            &:hover{
                                color: red
                            }
                        }
                    }
                }
            }
        }
        .content-button{
            background: #f9f9f9;
            height: 62px;
            padding: 14px 20px;
            font-size: 0;
            .btn{
                font-size: 14px;
                width: 110px;
                height: 34px;
                line-height: 34px;
                border-radius: 0;
                margin-right: 10px;
                // border: 0;
                border-radius: 2px;
            }
            .info{
                float: right;
                font-size: 14px;
                height: 34px;
                line-height: 34px;
                cursor: pointer;
                input{
                    position: relative;
                    margin-right: 4px;
                }
            }
        }
    }
</style>
