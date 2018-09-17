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
    <div class="field-content">
        <div class="list-wrapper clearfix">
            <div class="left-list search-wrapper-hidden">
                <div class="search-wrapper">
                    <div class="the-host">
                        <bk-select :selected.sync="object.selected" @on-selected="object.filter = ''" :disabled="objectOptions.length === 1">
                            <bk-select-option v-for="(option, index) in objectOptions" 
                                :key="index"
                                :value="option.bkObjId"
                                :label="option.bkObjName">
                            </bk-select-option>
                        </bk-select>
                    </div>
                    <div class="search-field">
                        <input type="text" :placeholder="$t('Inst[\'搜索属性\']')" v-model.trim="object.filter">
                    </div>
                </div>

                <ul class="list-wrapper">
                    <li v-for="(property, index) in hideList" @click="addItem(property)" ref="hideItem" :key="index">
                        <span :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                        <i class="bk-icon icon-angle-right"></i>
                    </li>
                </ul>
            </div>
            <div class="right-list">
                <div class="title">
                    <div class="search-wrapper">
                        {{$t('Inst[\'已显示属性\']')}}
                    </div>
                </div>
                <draggable class="content-right" v-model="shownList" :options="{animation: 150}">
                    <div v-for="(property, index) in shownList" :key="index" class="item">
                        <i class="icon-triple-dot"></i><span :title="property['bk_property_name']">{{property['bk_property_name']}}</span><i class="bk-icon icon-eye-slash-shape" @click="removeItem(index)" v-tooltip="$t('Common[\'隐藏\']')"></i>
                    </div>
                </draggable>
            </div>
        </div>
        <div class="bk-form-item bk-form-action content-button">
            <bk-button class="btn" type="primary" :loading="$loading('userCustom')" @click="apply">
                {{$t('Inst[\'应用\']')}}
            </bk-button>
            <bk-button class="vice-btn btn reinstate cancel" type="default" @click="cancel">
                {{$t('Common[\'取消\']')}}
            </bk-button>
        </div>
    </div>
</template>

<script type="text/javascript">
    import draggable from 'vuedraggable'
    import sortUnicode from '@/common/js/sortUnicode'
    import { sortArray } from '@/utils/util'
    export default {
        props: {
            shownFields: {
                type: Array,
                required: true
            },
            fieldOptions: {
                type: Array,
                required: true
            },
            isShow: {
                type: Boolean,
                required: true
            },
            isShowExclude: {
                type: Boolean,
                default: true
            },
            minField: {
                type: Number,
                default: 1
            }
        },
        data () {
            return {
                object: {
                    selected: '',
                    filter: '',
                    list: []
                },
                shownList: [],
                excludeFields: ['bk_host_innerip', 'bk_host_outerip']
            }
        },
        computed: {
            objectOptions () {
                return this.fieldOptions.map(({bk_obj_id: bkObjId, bk_obj_name: bkObjName}, index) => {
                    if (index === 0) {
                        this.object.selected = bkObjId
                    }
                    return {bkObjId, bkObjName}
                })
            },
            shownProperty () {
                return this.shownList.map(property => {
                    return property['bk_property_id']
                })
            },
            hideList () {
                let hideList = []
                this.object.list.map(property => {
                    let {
                        bk_isapi: bkIsapi,
                        bk_property_id: bkPropertyId,
                        bk_obj_id: bkObjId
                    } = property
                    if (!bkIsapi) {
                        const isCurrentShownProperty = this.shownList.some(property => property['bk_obj_id'] === bkObjId && this.shownProperty.indexOf(bkPropertyId) !== -1)
                        if (!isCurrentShownProperty && property['bk_property_name'].toLowerCase().indexOf(this.object.filter.toLowerCase()) !== -1) {
                            if (this.isShowExclude) {
                                hideList.push(property)
                            } else if (this.excludeFields.indexOf(bkPropertyId) === -1) {
                                hideList.push(property)
                            }
                        }
                    }
                })
                return sortArray(hideList, 'bk_property_name')
            }
        },
        watch: {
            shownFields (header) {
                this.setShownList()
            },
            isShow (isShow) {
                if (isShow) {
                    this.setShownList()
                } else {
                    this.object.selected = ''
                }
            },
            'object.selected' (selectedObjId) {
                let targetFieldOption = this.fieldOptions.find(({bk_obj_id: bkObjId}) => selectedObjId === bkObjId)
                this.object.list = targetFieldOption ? targetFieldOption['properties'] : []
                this.object.filter = ''
            }
        },
        methods: {
            setSortKey (data) {
                data.map(item => {
                    let str = item['bk_property_name']
                    let sortKey = ''
                    for (let i = 0; i < str.length; i++) {
                        let code = str.charCodeAt(i)
                        if (code < 40869 && code >= 19968) {
                            sortKey += sortUnicode.strChineseFirstPY.charAt(code - 19968)
                        } else {
                            sortKey += str[i]
                        }
                    }
                    item.sortKey = sortKey
                })
            },
            addItem (property) {
                this.shownList.push(property)
            },
            removeItem (index) {
                if (this.shownList.length <= this.minField) {
                    this.$alertMsg(this.$t('Common[\'至少选择N个字段\']', {N: this.minField}))
                } else {
                    this.shownList.splice(index, 1)
                }
            },
            setShownList () {
                this.shownList = this.shownFields.slice(0)
                this.object.selected = this.objectOptions.length ? this.objectOptions[0]['bkObjId'] : ''
            },
            apply () {
                this.$emit('apply', this.shownList.slice(0))
            },
            cancel () {
                this.$emit('cancel')
            }
        },
        components: {
            draggable
        }
    }
</script>
<style media="screen" lang="scss" scoped>
    $primaryColor: #6b7baa;
    $lineColor: #e7e9ef;
    .field-content{
        height: calc(100% - 122px);
        .list-wrapper{
            height: 100%;
        }
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
                span{
                    display: inline-block;
                    width: 230px;
                    @include ellipsis;
                }
                &:hover{
                    background: #f9f9f9;
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
        .content-right{
            height: calc(100% - 78px);
            width:100% !important;
            padding: 15px 0 0 0;
            overflow-y: auto;
            @include scrollbar;
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
        .item{
            height: 43px;
            line-height: 42px;
            padding-left: 30px;
            cursor: move;
            span{
                display: inline-block;
                width: calc(100% - 66px);
                @include ellipsis;
                vertical-align: bottom;
            }
            &:hover{
                background: #f9f9f9;
            }
            .icon-triple-dot{
                position: relative;
                top: -1px;
                display: inline-block;
                width: 4px;
                height: 14px;
                margin-right: 10px;
                background: url(../../common/images/icon/icon-triple-dot.png);
            }
            .icon-eye-slash-shape{
                float: right;
                font-size: 12px;
                margin-top: 5px;
                margin-right: 20px;
                padding: 10px;
                cursor: pointer;
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
            margin-right: 10px;
            &.apply{
                &:hover{
                    background: #4d597d;
                }
            }
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
</style>
