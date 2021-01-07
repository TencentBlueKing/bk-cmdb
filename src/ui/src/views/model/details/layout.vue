<template>
    <div class="layout-wrapper">
        <div class="layout-box" v-bkloading="{isLoading: $loading(['searchObjectAttribute', 'searchGroup', 'updatePropertyGroup'])}">
            <div class="hidden-list">
                <div class="hidden-list-title">
                    <i class="bk-icon icon-eye-slash-shape"></i>
                    <span>{{$t('ModelManagement["隐藏字段"]')}}</span>
                </div>
                <vue-draggable ref="draggableHideField" 
                class="hidden-list-content"
                v-model="hiddenField" 
                index="0" 
                :options="{animation: 150, group: 'field'}"
                :move="checkMove">
                    <div v-for="(field, index) in hiddenField" :key="index" class="hidden-list-item">
                        <span class="triple-dot">
                            <i></i><i></i><i></i>
                        </span>
                        <span class="hidden-list-text">
                            <span class="text-name">{{field['bk_property_name']}}</span>
                            <i v-if="field['isrequired'] && !field['isonly']" class="icon-cc-required color-danger"></i>
                            <i v-if="field['isonly']" class="icon-cc-key"></i>
                        </span>
                    </div>
                </vue-draggable>
            </div>
            <ul class="layout-list">
                <li class="layout-group" v-for="(group, groupIndex) in groupFieldList" :key="groupIndex">
                    <div class="layout-title">
                        <span class="layout-title-text" v-if="!group.isEditTitle">{{group['bk_group_name']}}</span>
                        <input v-else v-focus maxlength="20" @blur="changeGroupName(group)" type="text" class="layout-title-text cmdb-form-input" v-model="group['bk_group_name']"
                        >
                        <i class="icon-cc-edit" @click.stop.prevent="editGroupName(group)"></i>
                        <span class="layout-title-icon">
                            <i class="bk-icon icon-arrows-up" v-if="groupIndex" @click="groupMove(groupIndex, groupIndex - 1)"></i>
                            <i class="bk-icon icon-arrows-down" v-if="groupIndex !== groupFieldList.length - 1" @click="groupMove(groupIndex, groupIndex + 1)"></i>
                            <i class="icon-cc-del" v-if="!(group['ispre'] || group['bk_group_id'] === 'default')" @click="deleteGroup(group, groupIndex)"></i>
                        </span>
                    </div>
                    <vue-draggable
                    class="layout-group-field clearfix"
                    v-model="group.properties"
                    :index="groupIndex"
                    :options="{animation: 150, group:'field'}" 
                    :move="checkMove" >
                        <div class="layout-item" 
                        v-for="(property, propertyIndex) in group.properties"
                        :key="propertyIndex">
                            <span class="triple-dot">
                                <i></i><i></i><i></i>
                            </span>
                            <span class="layout-list-text">
                                <span class="text-name">{{property['bk_property_name']}}</span>
                                <i v-if="property['isrequired'] && !property['isonly']" class="icon-cc-required color-danger"></i><i v-if="property['isonly']" class="icon-cc-key"></i>
                            </span>
                            <i class="bk-icon icon-eye-slash-shape" v-if="!property['isonly'] && !property['isrequired']" @click="hideField(property, propertyIndex, groupIndex)"></i>
                        </div>
                    </vue-draggable>
                </li>
                <li class="layout-list-add">
                    <span @click="addGroup">
                        <i class="bk-icon icon-plus"></i>
                        <span>{{$t('ModelManagement["新增字段分组"]')}}</span>
                    </span>
                </li>
            </ul>
        </div>
        <footer class="footer-btn" v-if="!isReadOnly">
            <bk-button type="primary" @click="confirm" :loading="$loading('updatePropertyGroup')">{{$t('Common["确定"]')}}</bk-button>
            <bk-button class="default" type="default" :title="$t('Common[\'取消\']')" @click="cancel">{{$t('Common["取消"]')}}</bk-button>
        </footer>
    </div>
</template>

<script>
    import vueDraggable from 'vuedraggable'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            vueDraggable
        },
        directives: {
            focus: {
                inserted: function (el) {
                    el.focus()
                }
            }
        },
        data () {
            return {
                activeGroupName: '',
                hiddenField: [],
                groupFieldList: [],
                groupFieldListCopy: [],
                hiddenFieldCopy: []
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                return this.activeModel['bk_ispaused']
            },
            isEditTitle () {
                let isEdit = this.groupFieldList.find(({isEditTitle}) => {
                    return isEditTitle
                })
                return !!isEdit
            }
        },
        created () {
            this.getFieldData()
        },
        methods: {
            ...mapActions('objectModelFieldGroup', [
                'searchGroup',
                'updateGroup',
                'createGroup',
                'updatePropertyGroup'
            ]),
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            isCloseConfirmShow () {
                let {
                    hiddenField,
                    hiddenFieldCopy,
                    groupFieldList,
                    groupFieldListCopy
                } = this
                if (hiddenField.length !== hiddenFieldCopy.length || groupFieldList.length !== groupFieldListCopy.length) {
                    return true
                }

                let result = hiddenField.some((field, index) => {
                    return hiddenFieldCopy[index]['bk_property_id'] !== field['bk_property_id']
                })
                if (result) {
                    return true
                }
                return groupFieldList.some((group, groupIndex) => {
                    let curGroupCopy = groupFieldListCopy[groupIndex]
                    if (group.properties.length !== curGroupCopy.properties.length) {
                        return true
                    }
                    let res = group.properties.some((property, propertyIndex) => {
                        return property['bk_property_id'] !== curGroupCopy.properties[propertyIndex]['bk_property_id']
                    })
                    return res
                })
            },
            editGroupName (group) {
                if (!this.isEditTitle) {
                    group.isEditTitle = true
                    this.activeGroupName = group['bk_group_name']
                }
            },
            hideField (property, propertyIndex, groupIndex) {
                this.groupFieldList[groupIndex].properties.splice(propertyIndex, 1)
                this.hiddenField.push(property)
            },
            async getFieldData () {
                const res = await Promise.all([
                    this.searchGroup({
                        objId: this.activeModel['bk_obj_id'],
                        config: {
                            requestId: 'searchGroup'
                        }
                    }),
                    this.searchObjectAttribute({
                        params: {
                            bk_supplier_account: this.supplierAccount,
                            bk_obj_id: this.activeModel['bk_obj_id']
                        },
                        config: {
                            requestId: 'searchObjectAttribute'
                        }
                    })
                ])
                this.setGroupFieldList(res)
            },
            setGroupFieldList (data) {
                let groupList = data[0]
                let fieldList = data[1]
                groupList.sort((groupA, groupB) => groupA['bk_group_index'] - groupB['bk_group_index'])
                groupList.map(group => {
                    Object.assign(group, {isEditTitle: false})
                    if (!group.hasOwnProperty('properties')) {
                        Object.assign(group, {properties: []})
                    }
                })
                let hiddenField = []
                fieldList.map(field => {
                    let group = groupList.find(({bk_group_id: groupId}) => groupId === field['bk_property_group'])
                    if (group) {
                        group.properties.push(field)
                    } else {
                        hiddenField.push(field)
                    }
                })

                groupList.map(group => {
                    group.properties.sort((propertyA, propertyB) => propertyA['bk_property_index'] - propertyB['bk_property_index'])
                })
                hiddenField.sort((propertyA, propertyB) => propertyA['bk_property_index'] - propertyB['bk_property_index'])
                
                this.groupFieldList = groupList
                this.hiddenField = hiddenField
                this.groupFieldListCopy = this.$tools.clone(groupList)
                this.hiddenFieldCopy = this.$tools.clone(hiddenField)
            },
            /**
             * 检查是否可移动到指定区域
             * @param evt {Object} - 拖拽对象的相关属性
             * @return - 返回false会取消移动操作
             */
            checkMove (evt) {
                // 唯一字段、必填字段不能够被隐藏
                return !(evt.to.attributes[2].value === 'hidden-list-content' && (evt.draggedContext.element.isonly || evt.draggedContext.element.isrequired))
            },
            async changeGroupName (group) {
                if (!this.checkGroupParams(group)) {
                    return
                }
                let params = {
                    condition: {
                        id: group.id
                    },
                    data: {
                        bk_group_name: group['bk_group_name']
                    }
                }
                await this.updateGroup({
                    params
                }).then(() => {
                    this.$http.cancel(`post_searchGroup_${this.activeModel['bk_obj_id']}`)
                })
                let activeGroup = this.groupFieldList.find(({id}) => id === group.id)
                activeGroup['bk_group_name'] = group['bk_group_name']
                group.isEditTitle = false
            },
            /**
             * 调整分组位置
             * @param from {Number} - 当前项的index
             * @param to {Number} - 要移动到的项的index
             */
            async groupMove (from, to) {
                let {
                    groupFieldList
                } = this
                await this.updateGroupIndex(groupFieldList[from], groupFieldList[to]);
                [groupFieldList[from], groupFieldList[to]] = [groupFieldList[to], groupFieldList[from]]
                this.$forceUpdate()
            },
            updateGroupIndex (fromGroup, toGroup) {
                return Promise.all([
                    this.updateGroup({
                        params: {
                            condition: {
                                id: fromGroup.id
                            },
                            data: {
                                bk_group_id: toGroup['bk_group_index']
                            }
                        }
                    }),
                    this.updateGroup({
                        params: {
                            condition: {
                                id: toGroup.id
                            },
                            data: {
                                bk_group_id: fromGroup['bk_group_index']
                            }
                        }
                    })
                ])
            },
            deleteGroup (group, groupIndex) {
                let isPropertyExist = group.properties.find(property => {
                    return property['isrequired']
                })
                if (isPropertyExist) {
                    this.$error(this.$t('ModelManagement["该分组中存在必填字段，不可删除"]'))
                    return
                }
                this.$store.dispatch('objectModelFieldGroup/deleteGroup', {
                    id: group.id
                })
                group.properties.map(property => {
                    property['bk_property_group'] = 'none'
                    this.hiddenField.push(property)
                })
                this.groupFieldList.splice(groupIndex, 1)
            },
            checkGroupParams (group) {
                if (group) {
                    if (this.activeGroupName === group['bk_group_name']) {
                        group.isEditTitle = false
                        return false
                    }
                    let isExist = this.groupFieldList.findIndex(({bk_group_name: bkGroupName, bk_group_id: bkGroupId}) => {
                        return bkGroupName === group['bk_group_name'] && bkGroupId !== group['bk_group_id']
                    }) > -1
                    if (isExist) {
                        this.$error(this.$t('ModelManagement["该名字已经存在"]'))
                        return false
                    }
                } else {
                    let isExist = this.groupFieldList.findIndex(({bk_group_name: bkGroupName}) => {
                        return bkGroupName === this.$t('ModelManagement["未命名"]')
                    }) > -1
                    if (isExist) {
                        this.$error(this.$t('ModelManagement["已经存在未命名分组"]'))
                        return false
                    }
                }
                return true
            },
            async addGroup () {
                if (this.isEditTitle || !this.checkGroupParams()) {
                    return
                }
                let reg = /^[0-9]+$/
                let groupId = 0
                let groupIndex = 0
                this.groupFieldList.map(({bk_group_id: bkGroupId, bk_group_index: bkGroupIndex}) => {
                    if (reg.test(bkGroupId)) {
                        groupId = parseInt(bkGroupId) > groupId ? parseInt(bkGroupId) : groupId
                    }
                    groupIndex = bkGroupIndex > groupIndex ? bkGroupIndex : groupIndex
                })
                groupId++
                groupIndex++

                const res = await this.createGroup({
                    params: {
                        bk_group_id: groupId.toString(),
                        bk_group_name: this.$t('ModelManagement["未命名"]'),
                        bk_group_index: groupIndex,
                        bk_obj_id: this.activeModel['bk_obj_id'],
                        bk_supplier_account: this.supplierAccount
                    }
                })
                this.groupFieldList.push({
                    bk_group_id: groupId.toString(),
                    bk_group_index: groupIndex,
                    bk_group_name: this.$t('ModelManagement["未命名"]'),
                    isEditTitle: false,
                    id: res.id,
                    properties: []
                })
            },
            async confirm () {
                let params = []
                this.groupFieldList.map(group => {
                    group.properties.map((property, index) => {
                        params.push({
                            condition: {
                                bk_obj_id: property['bk_obj_id'],
                                bk_property_id: property['bk_property_id'],
                                bk_supplier_account: this.supplierAccount
                            },
                            data: {
                                bk_property_group: group['bk_group_id'],
                                bk_property_index: index
                            }
                        })
                    })
                })
                this.hiddenField.map((property, index) => {
                    params.push({
                        condition: {
                            bk_obj_id: property['bk_obj_id'],
                            bk_property_id: property['bk_property_id'],
                            bk_supplier_account: this.supplierAccount
                        },
                        data: {
                            bk_property_group: 'none',
                            bk_property_index: index
                        }
                    })
                })
                await this.updatePropertyGroup({
                    params,
                    config: {
                        requestId: 'updatePropertyGroup'
                    }
                }).then(() => {
                    this.$http.cancel(`post_searchGroup_${this.activeModel['bk_obj_id']}`)
                })
                this.$success(this.$t('Common["更新成功"]'))
                this.groupFieldListCopy = this.$tools.clone(this.groupFieldList)
                this.hiddenFieldCopy = this.$tools.clone(this.hiddenField)
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .layout-wrapper {
        position: relative;
        height: 100%;
        padding: 30px 30px 20px;
        .layout-box {
            height: calc(100% - 64px);
            border: 1px solid $cmdbBorderLightColor;
        }
        .triple-dot {
            display: inline-block;
            position: absolute;
            left: 6px;
            top: 14px;
            width: 3px;
            height: 15px;
            overflow: hidden;
            i {
                float: left;
                width: 3px;
                height: 3px;
                background: $cmdbBorderLightColor;
                margin: 1px 0;
            }
        }
        .hidden-list {
            float: left;
            width: 143px;
            height: 100%;
            text-align: center;
            border-right: 1px solid $cmdbBorderLightColor;
            .hidden-list-title {
                background: $cmdbPrimaryColor;
                line-height: 42px;
                padding-top: 4px;
                font-weight: bold;
                font-size: 0;
                border-bottom: 1px solid $cmdbBorderLightColor;
                >i,
                >span {
                    font-size: 14px;
                    vertical-align: middle;
                }
                >i {
                    padding-right: 5px;
                }
            }
            .hidden-list-content {
                height: calc(100% - 46px);
                @include scrollbar;
                >div {
                    position: relative;
                    border-bottom: 1px solid $cmdbBorderLightColor;
                    font-size: 12px;
                    line-height: 40px;
                    height: 40px;
                    cursor: move;
                    transition: all .35s;
                    .icon-eye-slash-shape {
                        display: none;
                    }
                }
                .hidden-list-text {
                    padding: 0 10px 0 15px;
                    width: 100%;
                    height: 40px;
                    line-height: 40px;
                    font-size: 12px;
                    @include ellipsis;
                    .text-name {
                        display: inline-block;
                        max-width: calc(100% - 40px);
                        @include ellipsis;
                        vertical-align: middle;
                    }
                }
            }
        }
        .layout-list {
            float: right;
            width: calc(100% - 143px);
            padding: 0 20px;
            height: 100%;
            @include scrollbar;
            .layout-group {
                .layout-title {
                    line-height: 42px;
                    height: 46px;
                    padding-top: 4px;
                    border-bottom: 1px dashed $cmdbBorderLightColor;
                    &:hover {
                        .icon-cc-edit {
                            display: inline-block;
                        }
                        .layout-title-icon i {
                            opacity: 1;
                        }
                    }
                    .layout-title-text {
                        display: inline-block;
                        width: auto;
                        color: #c3cdd7;
                        font-weight: bold;
                        line-height: 24px;
                        height: 26px;
                        &.cmdb-form-input {
                            color: $cmdbTextColor;
                        }
                    }
                    .icon-cc-edit {
                        display: none;
                        cursor: pointer;
                    }
                    .layout-title-icon {
                        float: right;
                        cursor: pointer;
                        .icon-cc-del:hover {
                            color: $cmdbDangerColor;
                        }
                        i {
                            color: #c3cdd7;
                            opacity: .5;
                            transition: all .5s;
                        }
                    }
                }
                .layout-group-field {
                    width: 100%;
                    padding: 10px 0;
                    min-height: 50px;
                    font-size: 0;
                    >div {
                        float: left;
                        // display: inline-block;
                        position: relative;
                        width: 50%;
                        height: 30px;
                        @include ellipsis;
                        cursor: move;
                        &:hover {
                            background: #f1f7ff;
                            .triple-dot,
                            .icon-eye-slash-shape {
                                display: inline-block;
                            }
                        }
                    }
                    .layout-list-text {
                        display: inline-block;
                        padding: 0 32px 0 15px;
                        width: 100%;
                        height: 30px;
                        line-height: 30px;
                        font-size: 14px;
                        @include ellipsis;
                    }
                    .triple-dot {
                        display: none;
                        top: 7px;
                        i {
                            background: $cmdbTextColor;
                        }
                    }
                    .icon-eye-slash-shape {
                        display: none;
                        font-size: 12px;
                        position: absolute;
                        right: 12px;
                        top: 9px;
                        cursor: pointer;
                    }
                }
            }
            .layout-list-add {
                width: 100%;
                border-top:1px solid $cmdbBorderLightColor;
                text-align: center;
                color: $cmdbMainBtnColor;
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
        }
        .icon-cc-key {
            transform: scale(calc(9 / 12));
            color: #ffb400;
        }
        .icon-cc-required {
            font-size: 12px;
            transform: scale(calc(9 / 12));
        }
    }
</style>
