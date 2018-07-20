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
    <li class="tree-item-wrapper">
        <div class="item" :class="{'item-active': isActive, 'open': model.isFolder}" :style="{'padding-left': model.level * 20 + 'px'}">
            <i 
                class="tree-icon icon-cc-triangle-sider" 
                v-if="model.children&&!(model.children.length===0&&model.level===1)&&!(model.children&&model.children.length===0&&model.loadNode===2)&&model.loadNode!==1" 
                :class="model.isFolder ? 'tree-icon-open' : 'tree-icon-close'" 
                @click="open(model)" :style="{'left': model.level * 20 + 'px'}">
            </i>
            <span class="item-loading" v-show='model.loadNode==1' :style="{'left': model.level * 20 + 'px'}">
                <img alt="loading" src="../../common/images/icon/loading_2_16x16.gif">
            </span>
            <span @click="Func(model)" :class="isActive" :title="model['bk_inst_name']" class="item-name">
                <i :class="model['bk_obj_icon']" class="icon" v-if="model['bk_obj_icon']"></i><i v-else class="icon-none"></i>{{model['name']}}
                <span class="count" v-if="model.count">{{model.count}}</span>
                <span class="detail" v-if="!(model.level % 2)" @click="detailClick(model)">{{$t('Common["详情信息"]')}}</span>
            </span>
        </div>
        <ul v-show="model.isFolder">
            <tree-item v-for="(item, index) in model.children"
            :key="index"
            :model.sync="item"
            :trees.sync="trees"
            :callback="callback"
            :expand="expand"
            :detailCallback="detailCallback"
            ></tree-item>
        </ul>
    </li>
</template>

<script>
    export default {
        name: 'treeItem',
        props: {
            // 点击展开回调
            expand: {
                type: Function
            },
            model: {
                type: Object
            },
            // 树形列表
            trees: {
                type: Array
            },
            // 节点点击回调
            callback: {
                type: Function
            },
            // 点击详情回调
            detailCallback: {
                type: Function
            }
        },
        computed: {
            isActive () {
                return this.model.clickNode ? 'active' : ''
            }
        },
        watch: {
            'model.pageSize' (val) {
                this.pageSize.call(null, this.model, val)
            }
        },
        methods: {
            detailClick (m) {
                this.detailCallback.call(null, m)
            },
            Func (m) {
                // 查找点击的子节点
                var recurFunc = (data, list) => {
                    data.forEach((i) => {
                        if (i.id === m.id) {
                            this.$set(i, 'clickNode', true)
                            if (typeof this.callback === 'function') {
                                this.callback.call(null, m, list, this.trees)
                            }
                        } else {
                            this.$set(i, 'clickNode', false)
                        }
                        if (i.children) {
                            recurFunc(i.children, i)
                        }
                    })
                }
                recurFunc(this.trees, this.trees)
                this.trees.splice()
            },
            open (m) {
                this.$set(m, 'isExpand', !m.isExpand)
                if (typeof this.expand === 'function' && m.isExpand) {
                    if (m.loadNode !== 2) {  // 节点未加载完毕
                        this.expand.call(null, m)
                    } else {
                        m.isFolder = !m.isFolder
                    }
                } else {
                    m.isFolder = !m.isFolder
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .tree-item-wrapper{
        padding-left: 20px;
        color: #6b7baa;
        font-size: 12px;
        position: relative;
        padding-top: 18px;
        &:before{
            content: "";
            width: 20px;
            height: 1px;
            border-top: 1px dashed #c3cdd7;
            position: absolute;
            left: 6px;
            top: 26px;
            z-index: 1;
        }
        &:after{
            content: "";
            width: 1px;
            height: calc(100% - 34px);
            border-left: 1px dashed #c3cdd7;
            position: absolute;
            left: 26px;
            top: 26px;
            z-index: 1;
        }
        .item{
            z-index: 2;
            position: relative;
            cursor: pointer;
            height: 16px;
            line-height: 16px;
            font-size: 0;
            &:before{
                content: '';
                position: absolute;
                left: 0;
                top: -4px;
                width: 100%;
                height: 24px;
                padding-left: 20px;
                background-clip : content-box;
                background-color: transparent;
                box-sizing: border-box;
                z-index: -1;
            }
            &:hover:before{
                background-color: #e6f7ff;
            }
            &.item-active:before{
                background-color: #bae7ff;
            }
            &:hover{
                color: #498fe0;
                .item-name{
                    .name-icon{
                        background: #498fe0;
                    }
                    .operation-group{
                        display: block;
                    }
                }
            }
            .item-loading{
                display: inline-block;
                height: 16px;
                position: absolute;
                top: 0;
                left: 0;
            }
            .tree-icon{
                position: absolute;
                left: 0;
                top: 0;
                display: inline-block;
                width: 16px;
                height: 16px;
            }
            .tree-icon-close{
                background: url(../../common/images/icon/tree_up.png) no-repeat 0 0 #fff;
            }
            .tree-icon-open{
                background: url(../../common/images/icon/tree_down.png) no-repeat 0 0 #fff;
            }
            .item-name{
                font-size: 12px;
                padding: 0 0 0 25px;
                width: 100%;
                display: inline-block;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
                border: 4px solid transparent;
                // border-right-width: 25px;
                border-left-width: 0;
                margin-top: -4px;
                .operation-group{
                    display: none;
                    margin-top: 6px;
                    height: 22px;
                    line-height: 1;
                    select{
                        margin-top: 1px;
                    }
                    .arrow-btn{
                        display: inline-block;
                        height: 22px;
                        border: 1px solid #e7e9ef;
                        vertical-align: top;
                        line-height: 18px;
                        .left{
                            display: inline-block;
                            transform: rotate(90deg);
                        }
                        .right{
                            display: inline-block;
                            transform: rotate(-90deg);
                        }
                    }
                }
                .name-icon{
                    display: inline-block;
                    font-style: normal;
                    background: #6b7baa;
                    color: #fff;
                    border-radius: 5px;
                    width: 16px;
                    height: 16px;
                    line-height: 15px;
                    text-align: center;
                    vertical-align: top;
                    margin-right: 8px;
                }
            }
            .active{
                color: #498fe0;
                .name-icon{
                    background: #498fe0;
                }
            }
        }
    }
</style>

<style lang="scss">
    $borderColor: #bec6de;
    .tree-wrapper.list{
        .tree-item-wrapper{
            color: #6b7baa;
            padding: 0;
            &:before{
                width: 0;
                height: 0;
            }
            &:after{
                width: 0;
                height: 0;
            }
            .item-active{
                background: #f1f7ff;
                .item-name{
                    color: #6b7baa;
                    &.active{
                        color: #3c96ff;
                        .icon{
                            color: #3c96ff;
                        }
                        .icon-none{
                            background: #3c96ff;
                        }
                    }
                }
            }
        }
        .root{
            >.item{
                &.open{
                    border-bottom: 1px solid $borderColor;
                }
                .tree-icon{
                    left: 0;
                }
            }
            .item{
                height: 36px;
                line-height: 36px;
                &:hover{
                    background: #f1f7ff;
                    >.item-name{
                        >.detail{
                            display: block;
                        }
                    }
                }
                &:before{
                    padding: 0;
                    width: 0;
                    height: 0;
                }
                .item-name{
                    .icon{
                        margin-right: 10px;
                        vertical-align: text-bottom;
                        font-size: 14px;
                        color: #ffb400;
                    }
                    .icon-none{
                        display: inline-block;
                        margin-top: -2px;
                        margin-right: 8px;
                        border-radius: 2px;
                        width: 10px;
                        height: 10px;
                        background: #c3cdd7;
                    }
                    .count{
                        float: right;
                        display: inline-block;
                        vertical-align: text-bottom;
                        margin-top: 11px;
                        margin-left: 10px;
                        margin: 11px 13px 11px 10px;
                        color: #ffb80f;
                        background: #fff7e5;
                        border-radius: 7px;
                        height: 14px;
                        font-size: 12px;
                        line-height: 14px;
                        padding: 0 8px;
                    }
                    .detail{
                        display: none;
                        float: right;
                        margin: 7px 10px;
                        padding: 0 5px;
                        height: 22px;
                        line-height: 20px;
                        border: 1px solid #c3cdd7;
                        border-radius: 2px;
                        color: #737987;
                        background: #fff;
                        &:hover{
                            background: #fafafa;
                        }
                    }
                }
                .tree-icon{
                    top: 10px;
                    font-size: 14px;
                    color: #498fe0;
                    background: none;
                    transition: all .2s;
                }
                .tree-icon-open{
                    transform: rotate(-135deg);
                }
                .tree-icon-close{
                    color: #bec6de;
                    transform: rotate(-180deg);
                }
                .item-loading{
                    top: 10px;
                }
            }
        }
    }
</style>

