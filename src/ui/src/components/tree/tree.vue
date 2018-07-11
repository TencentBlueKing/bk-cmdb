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
    <ul class="tree-wrapper list">
        <tree-item class="root" :key='index' v-for="(item, index) in treeDataSource"
        :model.sync="item"
        :trees.sync="treeDataSource"
        :callback="callback"
        :detailCallback="detailCallback"
        :expand="expand"
        ></tree-item>
    </ul>
</template>

<script type="text/javascript">
    import treeItem from './treeItem'
    export default {
        props: {
            // 点击展开回调
            expand: {
                type: Function
            },
            // 点击节点回调
            callback: {
                type: Function
            },
            list: {
                type: Array
            },
            isOpen: {
                default: false,
                type: Boolean
            },
            detailCallback: {
                type: Function
            }
        },
        data () {
            return {
                treeDataSource: [],
                level: 0
            }
        },
        watch: {
            'list': {
                handler: function () {
                    this.initTreeData()
                },
                deep: true
            }
        },
        methods: {
            getTree () {
                return this.treeDataSource
            },
            initTreeData () {
                var tempList = JSON.parse(JSON.stringify(this.list))
                // 递归操作，增加删除一些属性。比如: 展开/收起
                var recurrenceFunc = (data, level) => {
                    data.forEach(m => {
                        if (this.hasCheckbox) {
                            if (!m.hasOwnProperty('isChecked')) {
                                m.isChecked = false
                            }
                        }
                        if (!m.hasOwnProperty('clickNode')) {
                            // this.$set(m, 'clickNode', false)
                            // m.clickNode = false
                            m.clickNode = m.hasOwnProperty('clickNode') ? m.clickNode : false
                        }
                        m.children = m.children || []

                        if (!m.hasOwnProperty('isFolder')) {
                            m.isFolder = m.hasOwnProperty('open') ? m.open : this.isOpen
                        }
                        if (!m.hasOwnProperty('isExpand')) {
                            m.isExpand = m.hasOwnProperty('open') ? m.open : this.isOpen
                        }
                        m.loadNode = 0

                        m.level = level
                        let tempLevel = level
                        recurrenceFunc(m.children, ++tempLevel)
                    })
                }
                recurrenceFunc(tempList, 1)

                this.treeDataSource = tempList
            }
        },
        update () {
            this.initTreeData()
        },
        mounted () {
            this.$nextTick(() => {
                this.initTreeData()
            })
        },
        components: {
            treeItem
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    ul {
        line-height: 1.5em;
        list-style-type: dot;
    }
    .root{
        &:before{
            width: 0 !important;
            height: 0 !important;
        }
    }
    .tree-wrapper.list{
        border: 1px solid #bec6de;
        border-bottom: 0;
        .tree-item-wrapper{
            border-bottom: 1px solid #bec6de;
        }
    }
</style>

