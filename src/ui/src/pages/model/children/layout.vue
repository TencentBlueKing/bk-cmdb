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
    <div class="tab-content">
        <div class="table-content" v-bkloading="{isLoading: isLoading}">
            <div class="hidden-list">
                <div class="hidden-list-title">
                    <i class="bk-icon icon-eye-slash-shape"></i>
                    {{$t('ModelManagement["隐藏字段"]')}}
                </div>
                <ul>
                    <draggable ref="draggableHideField" v-model="hideAttr" index="0" class="content-left" :options="{animation: 150, group:'field'}" :move="checkMove" @end="moveEnd">
                        <li v-for="(attr, index) in hideAttr" :key="index">
                            <span class="hidden-list-icon">
                                <i></i><i></i><i></i>
                            </span>
                            <span class="hidden-list-text">
                                <span class="text-name">{{attr['bk_property_name']}}</span>
                                <i v-if="attr['isrequired'] && !attr['isonly']" class="icon-cc-required"></i><i v-if="attr['isonly']" class="icon-cc-key"></i>
                            </span>
                        </li>
                    </draggable>
                </ul>
            </div>
            <div class="layout-list">
                <div class="layout-list-ul" v-for="(group, groupIndex) in groupAttrList" :key="groupIndex">
                    <div class="layout-list-title">
                        <span class="layout-title-text" v-if="!group.isEditTitle">{{group['bk_group_name']}}</span>
                        <input v-else v-focus maxlength="20" @blur="changeGroupName(group)" type="text" class="layout-title-text border" v-model="group['bk_group_name']"
                        >
                        <i class="icon-cc-edit" @click.stop.prevent="editGroupName(group)"></i>
                        <span class="layout-title-icon">
                            <i class="bk-icon icon-arrows-up" v-if="groupIndex" @click="groupMove(groupAttrList, groupIndex, groupIndex - 1)"></i>
                            <i class="bk-icon icon-arrows-down" v-if="groupIndex !== groupAttrList.length - 1" @click="groupMove(groupAttrList, groupIndex, groupIndex + 1)"></i>
                            <i class="icon-cc-del f14 vm" @click="deleteGroup(group, groupIndex)"></i>
                        </span>
                    </div>
                    <ul>
                        <draggable class="content-right" :index="groupIndex" v-model="group.properties" :options="{animation: 150, group:'field'}" :move="checkMove" @end="moveEnd">
                            <li v-for="(property, propertyIndex) in group.properties" :key="propertyIndex">
                                <span class="layout-list-icon">
                                    <i></i><i></i><i></i>
                                </span>
                                <span class="layout-list-text">
                                    <span class="text-name">{{property['bk_property_name']}}</span>
                                    <i v-if="property['isrequired'] && !property['isonly']" class="icon-cc-required"></i><i v-if="property['isonly']" class="icon-cc-key"></i>
                                </span>
                                <i class="bk-icon icon-eye-slash-shape" v-if="!property['isonly'] && !property['isrequired']" @click="deleteAttr(property, propertyIndex, groupIndex)"></i>
                            </li>
                        </draggable>
                    </ul>
                </div>
                <div class="layout-list-add">
                    <span @click="addGroup">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t('ModelManagement["新增字段分组"]')}}</span>
                    </span>
                </div>
            </div>
        </div>
        <div class="base-info">
            <bk-button class="btn main-btn" type="primary" :title="$t('Common[\'确认\']')" :loading="$loading('updatePropertyIndex')" @click="confirm">{{$t('Common["确定"]')}}</bk-button>
            <bk-button class="btn vice-btn cancel-btn-sider" type="default" :title="$t('Common[\'取消\']')" @click="cancel">{{$t('Common["取消"]')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import draggable from 'vuedraggable'
    export default {
        props: {
            // 字段分组tab展显示状态
            isShow: {
                default: false
            },
            activeModel: {
                type: Object
            }
        },
        data () {
            return {
                isLoading: false,
                activeGroupName: '',    // 当前编辑的分组名
                attrGroup: [],          // 属性分组
                attrList: [],           // 全部属性
                groupAttrList: [],      // 按分组排好序的属性
                groupAttrListCopy: [],  // 保存时做比对
                hideAttr: [],           // 隐藏字段
                hideAttrCopy: [],       // 隐藏字段 保存时做比对
                activeAttr: {}          // 当前移动的属性
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            isEditTitle () {
                let isEdit = this.groupAttrList.find(({isEditTitle}) => {
                    return isEditTitle
                })
                return !!isEdit
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.getAttrData()
                }
            }
        },
        methods: {
            deleteAttr (property, propertyIndex, groupIndex) {
                this.groupAttrList[groupIndex].properties.splice(propertyIndex, 1)
                this.hideAttr.push(property)
            },
            /**
             * 检查是否可移动到指定区域
             * @param evt {Object} - 拖拽对象的相关属性
             * @return - 返回false会取消移动操作
             */
            checkMove (evt) {
                this.activeAttr = evt.draggedContext.element
                // 唯一字段、必填字段不能够被隐藏
                return !(evt.to.attributes[2].value === 'content-left' && (evt.draggedContext.element.isonly || evt.draggedContext.element.isrequired))
            },
            /**
             * 移动结束回调
             */
            moveEnd (evt) {
                this.$forceUpdate()
            },
            /**
             * 确认
             * @param isSaveAll {Boolean} - 是否保存全部属性 只在点击确认按钮时需要保存全部
             */
            async confirm (isSaveAll = true) {
                this.isLoading = true
                let params = []
                if (isSaveAll) {
                    this.groupAttrList.map(group => {
                        group.properties.map((property, index) => {
                            params.push({
                                condition: {
                                    bk_obj_id: property['bk_obj_id'],
                                    bk_property_id: property['bk_property_id'],
                                    bk_supplier_account: this.bkSupplierAccount
                                },
                                data: {
                                    bk_property_group: group['bk_group_id'],
                                    bk_property_index: index
                                }
                            })
                        })
                    })
                }
                this.hideAttr.map((property, index) => {
                    params.push({
                        condition: {
                            bk_obj_id: property['bk_obj_id'],
                            bk_property_id: property['bk_property_id'],
                            bk_supplier_account: this.bkSupplierAccount
                        },
                        data: {
                            bk_property_group: 'none',
                            bk_property_index: index
                        }
                    })
                })
                try {
                    await this.$axios.put('/objectatt/group/property', params, {id: 'updatePropertyIndex'})
                    this.$alertMsg(this.$t('Common["更新成功"]'), 'success')
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                } finally {
                    this.isLoading = false
                }
            },
            /**
             * 取消
             */
            cancel () {
                this.$emit('cancel')
            },
            /**
             * 调整分组位置
             * @param groupAttrList {Array} - 分组列表
             * @param from {Number} - 当前项的index
             * @param to {Number} - 要移动到的项的index
             */
            async groupMove (groupAttrList, from, to) {
                await this.updateGroupIndex(groupAttrList[from], groupAttrList[to]);
                [groupAttrList[from], groupAttrList[to]] = [groupAttrList[to], groupAttrList[from]]
                this.$forceUpdate()
            },
            /**
             * 更新属性分组信息
             */
            async updateGroupIndex (fromGroup, toGroup) {
                this.isLoading = true
                let groupList = [fromGroup, toGroup]
                await this.$Axios.all(groupList.map((group, index) => {
                    let params = {
                        condition: {
                            id: group.id
                        },
                        data: {
                            bk_group_index: index ? fromGroup['bk_group_index'] : toGroup['bk_group_index']
                        }
                    }
                    return this.$axios.put('/objectatt/group/update', params)
                }))
                this.isLoading = false
            },
            /**
             * 添加分组
             */
            async addGroup () {
                if (this.isEditTitle || !this.checkGroupParams()) {
                    return
                }
                this.isLoading = true
                // 取 groupId groupIndex
                let reg = /^[0-9]+$/
                let groupId = 0
                let groupIndex = 0
                this.groupAttrList.map(({bk_group_id: bkGroupId, bk_group_index: bkGroupIndex}) => {
                    if (reg.test(bkGroupId)) {
                        groupId = parseInt(bkGroupId) > groupId ? parseInt(bkGroupId) : groupId
                    }
                    groupIndex = bkGroupIndex > groupIndex ? bkGroupIndex : groupIndex
                })
                groupId++
                groupIndex++

                let params = {
                    bk_group_id: groupId.toString(),  // groupID唯一，前端不展示
                    bk_group_name: this.$t('ModelManagement["未命名"]'),
                    bk_group_index: groupIndex,
                    bk_obj_id: this.activeModel['bk_obj_id'],
                    bk_supplier_account: this.bkSupplierAccount
                }
                try {
                    let res = await this.$axios.post('/objectatt/group/new', params)
                    this.groupAttrList.push({
                        bk_group_id: groupId.toString(),
                        bk_group_index: groupIndex,
                        bk_group_name: this.$t('ModelManagement["未命名"]'),
                        isEditTitle: false,
                        id: res.data.id,
                        properties: []
                    })
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                } finally {
                    this.isLoading = false
                }
            },
            /**
             * 将分组名切换为编辑态
             * @param group {Object} - 分组信息
             */
            editGroupName (group) {
                if (!this.isEditTitle) {
                    group.isEditTitle = true
                    this.activeGroupName = group['bk_group_name']
                }
            },
            /**
             * 修改分组名
             * @param group {Object} - 分组信息
             */
            async changeGroupName (group) {
                if (!this.checkGroupParams(group)) {
                    return
                }
                this.isLoading = true
                let params = {
                    condition: {
                        id: group.id
                    },
                    data: {
                        bk_group_name: group['bk_group_name']
                    }
                }
                try {
                    await this.$axios.put('/objectatt/group/update', params)
                    let activeGroup = this.groupAttrList.find(({id}) => {
                        return id === group.id
                    })
                    activeGroup['bk_group_name'] = group['bk_group_name']
                    group.isEditTitle = false
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                } finally {
                    this.isLoading = false
                }
            },
            /**
             * 删除分组
             */
            async deleteGroup (group, groupIndex) {
                if (group['ispre']) {
                    this.$alertMsg(this.$t('ModelManagement["系统内置分组不可删除"]'))
                    return
                }
                if (group['bk_group_id'] === 'default') {
                    this.$alertMsg(this.$t('ModelManagement["默认字段分组不可删除"]'))
                    return
                }
                let property = group.properties.find(property => {
                    return property['isrequired']
                })
                if (property) {
                    this.$alertMsg(this.$t('ModelManagement["该分组中存在必填字段，不可删除"]'))
                    return
                }
                try {
                    await this.$axios.delete(`/objectatt/group/groupid/${group['id']}`)
                    // 该分组下有属性时更新属性分组为none
                    if (group.properties.length) {
                        this.confirm(false)
                        group.properties.map(property => {
                            property['bk_property_group'] = 'none'
                            this.hideAttr.push(property)
                        })
                    } else {
                        this.$alertMsg(this.$t('ModelManagement["删除分组成功"]'), 'success')
                    }
                    this.groupAttrList.splice(groupIndex, 1)
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            /**
             * 获取字段相关信息
             */
            async getAttrData () {
                this.isLoading = true
                await this.getAttrGroup()
                await this.getAttr()
                this.setGroupAttrList()
                this.isLoading = false
            },
            /**
             * 获取属性分组
             */
            async getAttrGroup () {
                try {
                    let res = await this.$axios.post(`/objectatt/group/property/owner/${this.bkSupplierAccount}/object/${this.activeModel['bk_obj_id']}`, {})
                    this.attrGroup = res.data
                    this.attrGroup.sort((groupA, groupB) => {
                        return groupA['bk_group_index'] - groupB['bk_group_index']
                    })
                } catch (e) {
                    this.isLoading = false
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            /**
             * 获取属性
             */
            async getAttr () {
                try {
                    let params = {
                        bk_obj_id: this.activeModel['bk_obj_id'],
                        bk_supplier_account: this.bkSupplierAccount
                    }
                    let res = await this.$store.dispatch('object/getAttribute', {objId: this.activeModel['bk_obj_id'], force: true})
                    this.attrList = res.data
                } catch (e) {
                    this.isLoading = false
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            /**
             * 将属性分组
             */
            setGroupAttrList () {
                let groupAttrList = this.$deepClone(this.attrGroup)
                groupAttrList.map(group => {
                    this.$set(group, 'isEditTitle', false)
                    if (!group.hasOwnProperty('properties')) {
                        group['properties'] = []
                    }
                })
                let hideAttr = []
                this.attrList.map(attr => {
                    let group = groupAttrList.find(({bk_group_id: bkGroupId}) => {
                        return bkGroupId === attr['bk_property_group']
                    })
                    if (group) {
                        group.properties.push(attr)
                    } else {
                        hideAttr.push(attr)
                    }
                })

                // 排序
                groupAttrList.map(group => {
                    group.properties.sort((propertyA, propertyB) => {
                        return propertyA['bk_property_index'] - propertyB['bk_property_index']
                    })
                })
                hideAttr.sort((propertyA, propertyB) => {
                    return propertyA['bk_property_index'] - propertyB['bk_property_index']
                })

                this.groupAttrList = groupAttrList
                this.hideAttr = hideAttr
                this.groupAttrListCopy = this.$deepClone(this.groupAttrList)
                this.hideAttrCopy = this.$deepClone(this.hideAttr)
            },
            /**
             * 检查分组名参数
             * @param group {Object} - 分组信息
             */
            checkGroupParams (group) {
                if (group) {
                    if (this.activeGroupName === group['bk_group_name']) {
                        group.isEditTitle = false
                        return false
                    }
                    let isExist = this.groupAttrList.findIndex(({bk_group_name: bkGroupName, bk_group_id: bkGroupId}) => {
                        return bkGroupName === group['bk_group_name'] && bkGroupId !== group['bk_group_id']
                    }) > -1
                    if (isExist) {
                        this.$alertMsg(this.$t('ModelManagement["该名字已经存在"]'))
                        return false
                    }
                } else {
                    let isExist = this.groupAttrList.findIndex(({bk_group_name: bkGroupName}) => {
                        return bkGroupName === this.$t('ModelManagement["未命名"]')
                    }) > -1
                    if (isExist) {
                        this.$alertMsg(this.$t('ModelManagement["已经存在未命名分组"]'))
                        return false
                    }
                }
                return true
            }
        },
        directives: {
            focus: {
                inserted: function (el) {
                    el.focus()
                }
            }
        },
        components: {
            draggable
        }
    }
</script>

<style lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #f9f9f9; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; // 主要颜色
    .tab-content{
        padding: 30px 30px 85px 30px;
        .table-content{
            width: 100%;
            height: 100%;
            overflow-y: auto;
            border:1px solid $borderColor;
            .hidden-list{
                width: 143px;
                height: 100%;
                float: left;
                border-right:1px solid $borderColor;
                text-align: center;
                &-title{
                    background: $primaryColor;
                    line-height: 42px;
                    height: 46px;
                    padding-top: 4px;
                    font-weight: bold;
                    border-bottom:1px solid $borderColor;
                }
                ul{
                    height: calc(100% - 46px);
                    overflow-y: auto;
                    &::-webkit-scrollbar{
                        width: 6px;
                        height: 5px;
                    }
                    &::-webkit-scrollbar-thumb{
                        border-radius: 20px;
                        background: #a5a5a5;
                    }
                    .content-left{
                        height: 100%;
                    }
                    li{
                        cursor: move;
                        border-bottom:1px solid $borderColor;
                        font-size: 0;
                        position: relative;
                        -webkit-transition: all .35s;
                        transition: all .35s;
                        .hidden-list-text,
                        .layout-list-text{
                            font-size: 12px;
                            width: 100%;
                            display: inline-block;
                            overflow: hidden;
                            white-space: nowrap;
                            text-overflow: ellipsis;
                            height: 40px;
                            line-height: 40px;
                            padding: 0 10px 0 15px;
                            .text-name {
                                display: inline-block;
                                max-width: calc(100% - 40px);
                                overflow: hidden;
                                text-overflow: ellipsis;
                                vertical-align: middle;
                            }
                        }
                        .hidden-list-icon,
                        .layout-list-icon{
                            position: absolute;
                            left: 6px;
                            top: 14px;
                            width: 3px;
                            height: 15px;
                            overflow: hidden;
                            display: inline-block;
                            i{
                                display: inline-block;
                                width: 3px;
                                height: 3px;
                                background: $borderColor;
                                margin: 1px 0;
                            }
                        }
                        &:hover{
                            box-shadow: 0 0 6px #eeeeee;
                            .hidden-list-icon,
                            .layout-list-icon{
                                i{
                                    background: $primaryHoverColor;
                                }
                            }
                        }
                    }
                }
            }
            .layout-list{
                width: calc(100% - 143px);
                float: right;
                padding: 0 20px;
                height: 100%;
                overflow-y: auto;
                &::-webkit-scrollbar{
                    width: 6px;
                    height: 5px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 20px;
                    background: #a5a5a5;
                }
                .content-right{
                    min-height: 30px;
                }
                &-icon{
                    position: absolute;
                    left: 6px;
                    top: 7px;
                    width: 3px;
                    height: 15px;
                    overflow: hidden;
                    display: none;
                    i{
                        display: inline-block;
                        width: 3px;
                        height: 3px;
                        background: $primaryHoverColor;
                        margin: 1px 0;
                        float: left;
                    }
                }
                &-ul{
                    >ul{
                        width: 100%;
                        padding: 11px 0;
                        font-size: 0;
                        li{
                            position: relative;
                            width: 50%;
                            height: 30px;
                            display: inline-block;
                            overflow: hidden;
                            white-space: nowrap;
                            text-overflow: ellipsis;
                            cursor: move;
                            .icon-eye-slash-shape{
                                display: none;
                                font-size: 12px;
                                position: absolute;
                                right: 12px;
                                top: 9px;
                                cursor: pointer;
                            }
                            &:hover{
                                background-color: #f1f7ff;
                                .icon-eye-slash-shape,
                                .layout-list-icon{
                                    display: inline-block;
                                }
                            }
                        }
                    }
                }
                &-add{
                    width: 100%;
                    border-top:1px solid $borderColor;
                    text-align: center;
                    color: #498fe0;
                    line-height: 36px;
                    margin-top: 18px;
                    .icon-plus{
                        font-size: 12px;
                        position: relative;
                        top: -1px;
                        cursor: pointer;
                    }
                    span{
                        display: inline-block;
                        cursor: pointer;
                    }
                }
                &-title{
                    line-height: 42px;
                    height: 46px;
                    padding-top: 4px;
                    border-bottom:1px dashed $borderColor;
                    .layout-title-text{
                        color: #c3cdd7;
                        width: auto;
                        font-weight: bold;
                        display: inline-block;
                        line-height: 24px;
                        height: 26px;
                        border:1px solid $defaultColor;
                        /* padding: 0 10px; */
                        &.border{
                            border:1px solid $fnMainColor;
                            padding: 0 10px;
                        }
                    }
                    .icon-cc-edit{
                        display: none;
                        cursor: pointer;
                    }
                    .layout-title-icon{
                        float: right;
                        /* display: none; */
                        i{
                            color: #b4c1e8;
                            opacity: 0.5;
                            cursor: pointer;
                            -webkit-transition: all .35s;
                            transition: all .35s;
                            padding: 5px 0;
                            &:hover{
                                color: $primaryHoverColor;
                                &.icon-cc-del{
                                    color: #f05d5d;
                                }
                            }
                        }
                    }
                    &:hover{
                        .icon-cc-edit,
                        .layout-title-icon{
                            display: inline-block;
                            i{
                                opacity: 1;
                            }
                        }
                    }
                }
                &-text{
                    font-size: 14px;
                    width: 100%;
                    display: inline-block;
                    overflow: hidden;
                    white-space: nowrap;
                    text-overflow: ellipsis;
                    height: 30px;
                    line-height: 30px;
                    padding: 0 32px 0 15px;
                }
            }
        }
    }
    .base-info{
        width: 100%;
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
            -webkit-transition: all .35s !important;
            transition: all .35s !important;
        }
    }
    .icon-cc-required,
    .icon-cc-key {
        display: inline-block;
        transform: scale(calc(8 / 12));
        letter-spacing: 1px;
        font-size: 12px;
    }
    .icon-cc-required {
        color: #ff5656;
    }
    .icon-cc-key {
        color: #ffb400;
        transform: scale(calc(9 / 12));
    }
</style>
