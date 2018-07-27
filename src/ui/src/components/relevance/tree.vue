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
    <div class="relevance-wrapper" v-bkloading="{isLoading: isLoading}">
        <ul class="relevance-box" v-show="!isLoading">
            <template v-if="ztreeDataSourceList.length">
                <v-tree :list="ztreeDataSourceList"
                    :callback='nodeClick'
                    :expand="expandClick"
                    :is-open="false"
                    :detailCallback="detailClick"
                    ref="tree"
                ></v-tree>
            </template>
            <template v-else>
                <div class="relevance-none">
                    <img src="../../common/images/no-relevance.png" alt="">
                    <div>{{$t("Common['当前还未有关联项']")}}</div>
                </div>
            </template>
        </ul>
        <v-attribute
            class="list-attribute"
            ref="attribute"
            :isShow.sync="attr.isShow"
            :instId="attr.instId"
            :objId="attr.objId"
            :instName="attr.instName"
            :objName="attr.objName"
        ></v-attribute>
    </div>
</template>

<script>
    import vTree from '@/components/tree/tree'
    import vAttribute from './attribute'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            // 当前实例ID
            objId: {
                default: ''
            },
            // 具体某一项的ID
            ObjectID: {
                default: ''
            }
        },
        data () {
            return {
                isLoading: true,
                curNode: {},
                treeItemId: 1,                  // 树形图具体某一项的id 用于区分点击的哪一项 前端递增
                ztreeDataSourceList: [],
                attr: {
                    isShow: false,
                    instId: '',
                    objId: '',
                    instName: '',
                    objName: ''
                }
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount'])
        },
        created () {
            this.getRelationInfo()
        },
        methods: {
            /*
                树形图点击回调
                node: 当前点击的节点
                parent: 当前点击节点的父节点
                list: 树形图列表
            */
            nodeClick (node, parent, list) {
            },
            detailClick (node) {
                this.attr.instId = node['bk_inst_id']
                this.attr.objId = node['bk_obj_id']
                this.attr.instName = node['bk_inst_name']
                this.attr.objName = node['bk_obj_name']
                this.attr.isShow = true
            },
            /*
                节点展开与收起回调事件
            */
            expandClick (node) {
                this.curNode = node
                this.$set(node, 'isFolder', true)
                this.$set(node, 'isExpand', true)
                if (node.hasOwnProperty('bk_obj_icon') && node['bk_obj_icon'] === '') {
                    this.getInstChild(node)
                }
            },
            getInstChild (node) {
                node.loadNode = 1
                this.$axios.post(`inst/search/topo/owner/${this.bkSupplierAccount}/object/${node['bk_obj_id']}/inst/${node['bk_inst_id']}`, {
                    page: {
                        start: 0,
                        limit: 10
                    }
                }).then(res => {
                    if (res.result) {
                        this.$set(node, 'children', [])
                        let tree = []
                        res.data.map(model => {
                            if (model['bk_obj_id'] !== 'plat') {
                                model.name = model['bk_obj_name']
                                // 加入分页相关信息 默认是第一页 每页10条
                                model.page = 1
                                model.pageSize = 10
                                
                                model.isExpand = false
                                model.isFolder = false
                                model.level = node.level + 1
                                model.id = this.treeItemId++
                                if (model.count && model.hasOwnProperty('children') && model.children && model.children.length) {
                                    model.children.map(inst => {
                                        inst.name = inst['bk_inst_name']
                                        inst.level = node.level + 2
                                        inst.id = this.treeItemId++
                                        inst.children = []
                                    })
                                } else {
                                    model.children = []
                                }
                                model.loadNode = 2
                                tree.push(model)
                            }
                        })
                        this.$set(node, 'children', tree)
                        node.loadNode = 2
                    }
                })
            },
            pageSize (node) {
                this.getInstByPage(node)
            },
            pageTurning (node) {
                this.getInstByPage(node)
            },
            getInstByPage (node) {
                let method = 'post'
                this.$axios({
                    url: `inst/search/${this.bkSupplierAccount}/${node['bk_obj_id']}`,
                    method: method,
                    data: {
                        condition: {},
                        fields: [],
                        page: {
                            start: (node.page - 1) * node.pageSize,
                            limit: parseInt(node.pageSize)
                        }
                    }
                }).then(res => {
                    if (res.result) {
                        let children = []
                        res.data.info.map(inst => {
                            children.push({
                                bk_obj_id: inst['bk_obj_id'],
                                bk_inst_id: inst['bk_inst_id'],
                                id: this.treeItemId++,
                                name: inst['bk_inst_name'],
                                pageSize: node.pageSize,
                                page: node.page,
                                level: node.level + 1,
                                isExpand: false,
                                isFolder: false,
                                children: []
                            })
                        })
                        node.children = children
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                获取树形图信息
            */
            async getRelationInfo (ObjId, ObjectID) {
                try {
                    this.isLoading = true
                    let params = {
                        fields: [],
                        page: {
                            limit: 10
                        },
                        condition: {
                        }
                    }
                    let res = await this.$axios.post(`inst/search/topo/owner/${this.bkSupplierAccount}/object/${this.objId}/inst/${this.ObjectID}`, params)
                    let tree = []
                    res.data.map(model => {
                        if (model['bk_obj_id'] !== 'plat') {
                            model.level = 1
                            model.page = 1
                            model.pageSize = 10
                            model.loadNode = 2
                            model.id = this.treeItemId++
                            model.name = model['bk_obj_name']
                            if (!model.count) {
                                model.children = []
                            }
                            if (model.count && model.hasOwnProperty('children') && model.children && model.children.length) {
                                model.children.map(inst => {
                                    inst.name = inst['bk_inst_name']
                                    // 插入层级
                                    inst.level = 2
                                    // 根据添加唯一性id
                                    inst.id = this.treeItemId++
                                })
                            }
                            tree.push(model)
                        }
                    })
                    this.ztreeDataSourceList = tree
                } catch (e) {
                    this.$alertMsg(e.data['bk_error_msg'])
                } finally {
                    this.isLoading = false
                }
            }
        },
        components: {
            vTree,
            vAttribute
        }
    }
</script>

<style lang="scss" scoped>
    .relevance-wrapper{
        position: relative;
        min-height: 300px;
        .relevance-box{
            height: 100%;
        }
        .list-attribute{
            z-index: 99;
        }
        .relevance-none{
            margin-top: 50px;
            text-align: center;
            color: #656b81;
            img{
                margin-bottom: 20px;
                width: 131px;
            }
        }
    }
</style>
