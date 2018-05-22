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
        <template v-if="treeType==='list'">
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
                </span>
            </div>
        </template>
        <template v-else>
            <div class="item" :class="{'item-active': isActive}">
                <span class="tree-icon" v-if="model.children&&!(model.children&&model.children.length===0)&&model.loadNode!==1" :class="model.isFolder ? 'tree-icon-open' : 'tree-icon-close'" @click="open(model)"></span>
                <span class="item-loading" v-show='model.loadNode==1'>
                    <img alt="loading" src="../../common/images/icon/loading_2_16x16.gif">
                </span>
                <span v-if="hasCheckbox" class="bk-form-checkbox"><input type="checkbox" v-model="model.isChecked" @click="checkboxClick(model)"></span>
                <span @click="Func(model)" :class="isActive" :title="model.name" class="item-name"><i class="name-icon" v-if="model.ObjName">{{model.ObjName[0]}}</i><i class="name-icon" v-else-if="model.ObjId==='app'">业</i>{{model.name}}</span>
                <span>{{model.count}}</span>
            </div>
        </template>
        <ul v-show="model.isFolder">
            <tree-item v-for="(item, index) in model.children"
            :key="index"
            :hasCheckbox="hasCheckbox"
            :model.sync="item"
            :trees.sync="trees"
            :callback="callback"
            :expand="expand"
            :treeType="treeType"
            :pageSize="pageSize"
            :pageTurning="pageTurning"
            ></tree-item>
        </ul>
    </li>
</template>

<script>
    export default {
        name: 'treeItem',
        props: {
            // 改变分页信息回调
            pageSize: {
                type: Function
            },
            // 翻页回调
            pageTurning: {
                type: Function
            },
            // 树形图类型 默认为空 可选值 list
            treeType: {
                type: String,
                default: ''
            },
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
            /*
                是否包含有checkbox
            */
            hasCheckbox: {
                default: false,
                type: Boolean
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
            Func (m) {
                // 查找点击的子节点
                var recurFunc = (data, list) => {
                    data.forEach((i) => {
                        if (i.id === m.id) {
                            // i.clickNode = true
                            // console.log(i)
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
                // this.$set(m, 'isFolder', !m.isFolder)
                if (typeof this.expand === 'function' && m.isExpand) {
                    // this.expand.call(null, m)
                    if (m.loadNode !== 2) {  // 节点未加载完毕
                        this.expand.call(null, m)
                    } else {
                        m.isFolder = !m.isFolder
                    }
                } else {
                    m.isFolder = !m.isFolder
                }
            },
            checkboxClick (m) {
                // 判断是否有子级
                if (m.hasOwnProperty('children') && m.children.length) {
                    m.children.map(val => {
                        // 此处拿到的状态值是相反的
                        val.isChecked = !m.isChecked
                    })
                }
            },
            /*
                上一页
            */
            pageUp (m) {
                if (m.page === 1) {
                    return
                }
                m.page--
                this.pageTurning.call(null, m, 'prev')
            },
            /*
                下一页
            */
            pageDown (m) {
                if (m.page === Math.ceil(m.count / m.pageSize)) {
                    return
                }
                m.page++
                this.pageTurning.call(null, m, 'next')
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
                border-right-width: 25px;
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
