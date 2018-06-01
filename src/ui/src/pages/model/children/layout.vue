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
    <div class="tab-content" @click="blur" v-bkloading="{isLoading: isGroupLoading}">
        <div class="table-content">
            <div class="hidden-list">
                <div class="hidden-list-title">
                    <i class="bk-icon icon-eye-slash-shape"></i>
                    {{$t('ModelManagement["隐藏字段"]')}}
                </div>
                <ul>
                    <draggable ref="draggableHideField" index="0" class="content-left" v-model="hideField" :options="{animation: 150, group:'field'}" :move="checkMove" @end="checkEnd">
                        <li v-for="item in hideField">
                            <span class="hidden-list-icon">
                                <i></i><i></i><i></i>
                            </span>
                            <span class="hidden-list-text">
                                <span class="text-name">{{item['bk_property_name']}}</span>
                                <i v-if="item['isrequired'] && !item['isonly']" class="icon-cc-required"></i><i v-if="item['isonly']" class="icon-cc-key"></i>
                            </span>
                        </li>
                    </draggable>
                </ul>
            </div>
            <div class="layout-list">
                <div class="layout-list-ul" v-for="(vitems, vindex) in fieldClassification">
                    <div class="layout-list-title">
                        <span class="layout-title-text" v-show="!isEditTitle || vindex !== curFieldIndex">
                            {{vitems.label}}
                        </span>
                        <input type="text" class="layout-title-text border" 
                            v-focus="labelShow" 
                            v-model="labelList[vindex]" 
                            v-show="isEditTitle && vindex === curFieldIndex" 
                            maxlength="20"
                            @click.stop.prevent="focus()"
                            @keyup.enter="renameFieldGroup(vitems, vindex)" 
                            @blur="renameFieldGroup(vitems, vindex)"
                        >
                        <i class="icon-cc-edit" @click.stop.prevent="editTitle(vitems, vindex)" v-show="!isEditTitle"></i>
                        <span class="layout-title-icon">
                            <i class="bk-icon icon-arrows-up" @click="sortByGroupIndex(vitems, vindex - 1)" v-if="vindex !== 0"></i>
                            <i class="bk-icon icon-arrows-down" @click="sortByGroupIndex(vitems, vindex + 1)" v-if="vindex < fieldClassification.length - 1"></i>
                            <i class="icon-cc-del f14 vm" @click="deleteFieldGroup(vitems)"></i>
                        </span>
                    </div>
                    <ul>
                        <draggable class="content-right" :index="vindex" v-model="vitems.modelField" :options="{animation: 150, group:'field'}" :move="checkMove" @end="checkEnd">
                            <li v-for="item in vitems.modelField">
                                <span class="layout-list-icon">
                                    <i></i><i></i><i></i>
                                </span>
                                <span class="layout-list-text">
                                    <span class="text-name">{{item['bk_property_name']}}</span>
                                    <i v-if="item['isrequired'] && !item['isonly']" class="icon-cc-required"></i><i v-if="item['isonly']" class="icon-cc-key"></i>
                                </span>
                                <i class="bk-icon icon-eye-slash-shape" @click="deleteModelField(item)"></i>
                            </li>
                        </draggable>
                    </ul>
                </div>
                <div class="layout-list-add">
                    <span @click="addFieldGroup">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t('ModelManagement["新增字段分组"]')}}</span>
                    </span>
                </div>
            </div>
        </div>
        <!-- <div class="base-info">
            <button class="btn main-btn" type="primary" title="确认">确认</button>
            <button class="btn vice-btn cancel-btn-sider" type="default" title="取消">取消</button>
        </div> -->
    </div>
</template>

<script>
    import {mapGetters} from 'vuex'
    import draggable from 'vuedraggable'
    export default {
        props: {
            // 字段分组tab展显示状态
            isShow: {
                default: false
            },
            id: {               // 模型ID
                default: 0
            },
            objId: {
                default: ''
            },
            isNewField: {
                default: false
            }
        },
        data () {
            return {
                isGroupLoading: false,     // loading
                attrList: [],              // 属性列表
                groupNameId: -1,
                groupFieldList: [],        // 分组
                curMoveFieldId: -1,        // 当前移动的字段的id
                evtToIndex: -1,            // 移动至分组
                isEditTitle: false,        // 是否处于编辑标题状态
                curFieldIndex: -1,         // 当前编辑的字段 索引
                fieldClassification: [],   // 显示字段
                hideField: [],             // 隐藏字段
                isLoading: false,          // 是否处于加载列表状态
                lastGroupIndex: 0,         // 分组groupIndex的最大值
                labelShow: false,          // 是否focus到输入框
                labelList: []
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ])
        },
        watch: {
            isShow (val) {
                if (val) {
                    this.isGroupLoading = true
                    this.hideField = []
                    this.fieldClassification = []
                    this.labelList = []
                    this.getFieldGroup().then(res => {
                        if (res.result) {
                            this.getAttr()
                        }
                    })
                }
            },
            /*
                新增字段更新字段分栏
            */
            isNewField () {
                this.getAttr()
            }
        },
        methods: {
            /*
                点击空白区域取消编辑状态
            */
            blur () {
                this.isEditTitle = false
            },
            /*
                点击输入框保持编辑状态
            */
            focus () {
                this.isEditTitle = true
            },
            /*
                查询全部字段
            */
            getAttr () {
                let params = {
                    bk_obj_id: this.objId,
                    bk_supplier_account: this.bkSupplierAccount
                }
                this.$axios.post('/object/attr/search', params).then((res) => {
                    this.hideField = []
                    if (res.result) {
                        this.attrList = this.$deepClone(res.data)
                        let list = res.data
                        this.fieldClassification = []
                        this.labelList = []

                        list.map((item) => {
                            // 隐藏字段
                            if (item['bk_property_group'] === 'none') {
                                this.hideField.push(item)
                            }
                        })
                        
                        for (let t = 0; t < this.groupFieldList.length; t++) {
                            let modelFieldList = []
                            for (let item of list) {
                                // 根据不同的 PropertyGroup 划分分组
                                if (item['bk_property_group'] === this.groupFieldList[t]['bk_group_id']) {
                                    modelFieldList.push(item)
                                }
                            }
                            // 重新组合结构（分组id，分组名，分组排序值，分组下的所有字段）
                            let jsonGroup = {
                                id: this.groupFieldList[t]['id'],
                                groupID: this.groupFieldList[t]['bk_group_id'],
                                label: this.groupFieldList[t]['bk_group_name'],
                                propertyIndex: this.groupFieldList[t]['bk_group_index'],
                                modelField: this.$deepClone(modelFieldList),
                                isPre: this.groupFieldList[t]['ispre']
                            }
                            this.fieldClassification.push(jsonGroup)
                        }
                        // 根据 propertyIndex 对 fieldClassification 排序
                        this.fieldClassification = this.orderByPropertyIndex(this.fieldClassification)
                        for (const key in this.fieldClassification) {
                            this.labelList[key] = this.fieldClassification[key].label
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.isGroupLoading = false
                })
            },
            /*
                根据 propertyIndex 对 fieldClassification 重新排序
            */
            orderByPropertyIndex (item) {
                for (var i = 0; i < item.length; i++) {
                    for (let j = 0; j < item.length - 1 - i; j++) {
                        if (item[j].propertyIndex > item[j + 1].propertyIndex) {
                            let temp = item[j + 1]
                            item[j + 1] = item[j]
                            item[j] = temp
                        }
                    }
                }
                return item
            },
            /*
                向上向下排序
            */
            sortByGroupIndex (item, index) {
                let params = {}
                if (!this.isEditTitle) {
                    for (let i = 0; i < 2; i++) {
                        // 交换前后两个分组的 groupIndex 值
                        if (i === 0) {
                            params = {
                                condition: {
                                    id: item['id']
                                },
                                data: {
                                    bk_group_index: this.fieldClassification[index].propertyIndex
                                }
                            }
                        } else {
                            params = {
                                condition: {
                                    id: this.fieldClassification[index].id
                                },
                                data: {
                                    bk_group_index: item.propertyIndex
                                }
                            }
                        }
                        // 更新前后两个分组的信息
                        this.$axios.put(`/objectatt/group/update`, params).then((res) => {
                            if (res.result) {
                                this.getFieldGroup().then(res => {
                                    if (res.result) {
                                        this.getAttr()
                                    }
                                })
                            } else {
                                this.$alertMsg(res['bk_error_msg'])
                            }
                        })
                    }
                }
            },
            /*
                编辑标题
            */
            editTitle (item, index) {
                this.isEditTitle = true
                this.curFieldIndex = index
                this.labelShow = true
            },
            /*
                拖动过程获取目标字段信息
            */
            checkMove (evt) {
                this.curMoveFieldId = evt.draggedContext.element['id']   // 目标字段的id
                this.evtToIndex = evt.to.attributes[1].value          // 目标字段移除到分组索引
            },
            /*
                拖动停止
            */
            checkEnd (evt) {
                let params = {}
                // 往'隐藏字段'里移动
                if (evt.to.attributes[2].value === 'content-left') {
                    params = {
                        bk_property_group: 'none'
                    }
                } else {
                    params = {
                        bk_property_group: this.fieldClassification[this.evtToIndex].groupID.toString()
                    }
                }
                // 移动区域发生变化
                if (evt.from !== evt.to) {
                    this.$axios.put(`/object/attr/${this.curMoveFieldId}`, params).then((res) => {
                        if (!res.result) {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            /*
                查询分组信息
            */
            getFieldGroup () {
                return this.$axios.post(`/objectatt/group/property/owner/${this.bkSupplierAccount}/object/${this.objId}`, {}).then((res) => {
                    if (res.result) {
                        this.groupFieldList = res.data
                        // 对当前 groupFieldList 中所有的 groupIndex 值进行排序，获取最大的 groupIndex 值
                        let arr = []
                        this.groupFieldList.map((item) => {
                            arr.push(item['bk_group_index'])
                        })
                        arr.sort()
                        if (arr.length !== 0) {
                            // 后台默认字段的 groupIndex 为 -1，不支持增加分组时的 groupIndex 为 0 ，所以要手动调整
                            if (arr[arr.length - 1] === -1) {
                                arr[arr.length - 1] = 0
                            }
                            this.lastGroupIndex = arr[arr.length - 1] + 1
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res
                })
            },
            /*
                检查是否存在未命名分组
            */
            checkGroupName () {
                for (let i = 0; i < this.groupFieldList.length; i++) {
                    if (this.groupFieldList[i]['bk_group_name'] === this.$t('ModelManagement["未命名"]')) {
                        return false
                    }
                }
                return true
            },
            /*
                新增分组
            */
            addFieldGroup () {
                if (!this.checkGroupName()) {
                    this.$alertMsg(this.$t('ModelManagement["已经存在未命名分组"]'))
                    return
                }
                let rand = Math.random().toString(36).substr(2)
                let params = {
                    bk_group_id: rand,  // groupID唯一，前端不展示
                    bk_group_name: this.$t('ModelManagement["未命名"]'),
                    bk_group_index: this.lastGroupIndex,
                    bk_obj_id: this.objId,
                    bk_supplier_account: this.bkSupplierAccount
                }
                this.$axios.post('/objectatt/group/new', params).then((res) => {
                    if (res.result) {
                        this.getFieldGroup().then(res => {
                            if (res.result) {
                                this.getAttr()
                            }
                        })
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                对分组重命名
            */
            renameFieldGroup (item, index) {
                this.isEditTitle = false    // 不可编辑
                let status = true
                let rename = true   // 是否可以重命名
                // 已存在其他相同名称
                this.groupFieldList.forEach(element => {
                    if (this.labelList[index] === element['bk_group_name'] && item.label !== this.labelList[index]) {
                        status = false
                        this.$alertMsg(this.$t('ModelManagement["该名字已经存在"]'))
                        this.isEditTitle = true     // 编辑状态
                    }
                })
                if (this.labelList[index].length > 20) {
                    this.$alertMsg(this.$t('ModelManagement["分组长度超出限制"]'))
                    this.isEditTitle = true
                }
                if (this.isEditTitle) {
                    return
                }
                // 编辑前后保持一致
                if (item.label === this.labelList[index]) {
                    status = true
                    rename = false
                }
                // 没重名可以保存 || 前后名字相同
                if (status && rename) {
                    let params = {
                        condition: {
                            id: item['id']
                            // bk_group_id: item['bk_group_id']
                        },
                        data: {
                            bk_group_name: this.labelList[index]
                        }
                    }
                    this.$axios.put('/objectatt/group/update', params).then((res) => {
                        if (res.result) {
                            this.$alertMsg(this.$t('ModelManagement["修改成功"]'), 'success')
                            this.getFieldGroup().then(res => {
                                if (res.result) {
                                    this.getAttr()
                                    for (let iterator of this.fieldClassification[index].modelField) {
                                        // 更新分组下面所有字段的PropertyGroup
                                        let iteratorParams = {
                                            bk_property_group: item['bk_group_id'].toString()
                                        }
                                        this.$axios.put(`/object/attr/${iterator['id']}`, iteratorParams).then((res) => {
                                            if (res.result) {
                                                this.getAttr()
                                            } else {
                                                this.$alertMsg(res['bk_error_msg'])
                                            }
                                        })
                                    }
                                }
                            })
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
                this.labelShow = false
            },
            /*
                删除分组
            */
            deleteFieldGroup (item) {
                let status = true
                // 系统内置字段不可删除
                if (item['isPre']) {
                    status = false
                    this.$alertMsg(this.$t('ModelManagement["系统内置分组不可删除"]'))
                    return
                }
                // "默认字段"分组不可删除
                if (item['groupID'] === 'default') {
                    status = false
                    this.$alertMsg(this.$t('ModelManagement["默认字段分组不可删除"]'))
                } else {
                    if (item.modelField) {
                        item.modelField.forEach(ele => {
                            // modelField 存在必填字段
                            if (ele['isrequired']) {
                                status = false
                                this.$alertMsg(this.$t('ModelManagement["该分组中存在必填字段，不可删除"]'))
                            }
                        })
                    }
                }
                
                if (status) {
                    if (item.modelField) {
                        item.modelField.forEach(ele => {
                            let params = {
                                bk_property_group: 'none'
                            }
                            this.$axios.put(`/object/attr/${ele['id']}`, params).then((res) => {
                                if (res.result) {
                                    this.$alertMsg(this.$t('Common["删除成功"]'), 'success')
                                    this.getAttr()
                                } else {
                                    this.$alertMsg(res['bk_error_msg'])
                                }
                            })
                        })
                    }
                    this.$axios.delete(`/objectatt/group/groupid/${item['id']}`).then((res) => {
                        if (res.result) {
                            this.$alertMsg(this.$t('ModelManagement["删除分组成功"]'), 'success')
                            this.getFieldGroup().then(res => {
                                if (res.result) {
                                    this.getAttr()
                                }
                            })
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            /*
                删除某个模型属性
            */
            deleteModelField (item) {
                // 必填字段不可删
                let params = {
                    bk_property_group: 'none'
                }
                this.$axios.put(`/object/attr/${item['id']}`, params).then((res) => {
                    if (res.result) {
                        this.$alertMsg(this.$t('ModelManagement["隐藏成功"]'), 'success')
                        this.getAttr()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            }
        },
        directives: {
            focus: {
                update: function (el, {value}) {
                    if (value) {
                        el.focus()
                    }
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