<template>
    <div class="group-layout" v-bkloading="{isLoading: $loading()}">
        <div class="group"
            v-for="(group, index) in groupedProperties"
            :key="index">
            <div class="group-header clearfix">
                <div class="header-title fl">
                    <template v-if="group.info['bk_group_id'] !== 'none' && group === groupInEditing">
                        <input type="text" class="title-input cmdb-form-input"
                            ref="titleInput"
                            v-model.trim="groupNameInEditing">
                        <a class="title-input-button" href="javascript:void(0)" @click="handleUpdateGroupName(group)">{{$t('Common["保存"]')}}</a>
                        <a class="title-input-button" href="javascript:void(0)" @click="handleCancelEditGroupName">{{$t('Common["取消"]')}}</a>
                    </template>
                    <template v-else>
                        <span class="group-name">{{group.info['bk_group_name']}}</span>
                        <span class="group-count">({{group.properties.length}})</span>
                        <i class="title-icon icon icon-cc-edit"
                            v-if="authority.includes('update') && !isReadOnly && group.info['bk_group_id'] !== 'none'"
                            @click="handleEditGroupName(group)">
                        </i>
                    </template>
                </div>
                <div class="header-options fr" v-if="authority.includes('update') && !isReadOnly">
                    <i class="options-icon bk-icon icon-arrows-up"
                        v-tooltip="$t('ModelManagement[\'上移\']')"
                        :class="{disabled: index === 0 || ['none'].includes(group.info['bk_group_id'])}"
                        @click="handleRiseGroup(index, group)">
                    </i>
                    <i class="options-icon bk-icon icon-arrows-down"
                        v-tooltip="$t('ModelManagement[\'下移\']')"
                        :class="{disabled: index === (groupedProperties.length - 2) || ['none'].includes(group.info['bk_group_id'])}"
                        @click="handleDropGroup(index, group)">
                    </i>
                    <i class="options-icon bk-icon icon-plus-circle-shape"
                        v-tooltip="$t('ModelManagement[\'新建字段\']')"
                        @click="handleAddProperty(group)">
                    </i>
                    <i class="options-icon bk-icon icon-delete"
                        v-tooltip="$t('ModelManagement[\'删除分组\']')"
                        :class="{disabled: ['none', 'default'].includes(group.info['bk_group_id'])}"
                        @click="handleDeleteGroup(group, index)">
                    </i>
                </div>
            </div>
            <vue-draggable class="property-list clearfix"
                element="ul"
                v-model="group.properties"
                :options="{
                    group: 'property',
                    animation: 150,
                    filter: '.filter-empty',
                    disabled: !authority.includes('update') || isReadOnly
                }"
                :class="{empty: !group.properties.length}"
                @change="handleDragChange"
                @end="handleDragEnd">
                <li class="property-item fl"
                    v-for="(property, index) in group.properties"
                    :key="index"
                    :title="property['bk_property_name']">
                    {{property['bk_property_name']}}
                </li>
                <template v-if="!group.properties.length">
                    <li class="property-empty" v-if="authority.includes('update') && !isReadOnly" @click="handleAddProperty(group)">{{$t('ModelManagement["立即添加"]')}}</li>
                    <li class="property-empty disabled" v-else>{{$t('ModelManagement["暂无字段"]')}}</li>
                </template>
            </vue-draggable>
            <template v-if="authority.includes('update') && !isReadOnly">
                <div class="add-group" v-if="index === (groupedProperties.length - 2)">
                    <a class="add-group-trigger" href="javascript:void(0)"
                        v-if="!showAddGroup"
                        @click="handleAddGroup">
                        {{$t('ModelManagement["新建分组"]')}}
                        <i class="icon icon-cc-edit"></i>
                    </a>
                    <template v-else>
                        <input type="text" class="add-group-input cmdb-form-input"
                            ref="addGroupInput"
                            v-model.trim="newGroupName">
                        <a class="add-group-button" href="javascript:void(0)" @click="handleCreateGroup">{{$t('Common["保存"]')}}</a>
                        <a class="add-group-button" href="javascript:void(0)" @click="handleCancelCreateGroup">{{$t('Common["取消"]')}}</a>
                    </template>
                </div>
            </template>
        </div>
        <bk-dialog
            :is-show.sync="dialog.isShow"
            :has-header="false"
            :quick-close="false"
            :width="600"
            @cancel="handleCancelAddProperty"
            @confirm="handleConfirmAddProperty">
            <div class="dialog-title" slot="tools">{{$t('ModelManagement["新建字段"]')}}</div>
            <div class="dialog-content" slot="content">
                <div class="dialog-filter">
                    <input type="text" class="cmdb-form-input" v-model.trim="dialog.filter">
                    <i class="bk-icon icon-search"></i>
                </div>
                <ul class="dialog-property clearfix" ref="dialogProperty">
                    <li class="property-item fl"
                        v-for="(property, index) in sortedProperties"
                        v-show="filter(property)"
                        :key="index">
                        <label class="property-label"
                            :class="{
                                checked: dialog.selectedProperties.includes(property)
                            }"
                            :title="property['bk_property_name']"
                            @click="handleSelectProperty(property)">
                            {{property['bk_property_name']}}
                        </label>
                    </li>
                </ul>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import vueDraggable from 'vuedraggable'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            vueDraggable
        },
        data () {
            return {
                groupedProperties: [],
                groupInEditing: null,
                groupNameInEditing: '',
                shouldUpdatePropertyIndex: false,
                showAddGroup: false,
                newGroupName: '',
                dialog: {
                    isShow: false,
                    group: null,
                    filter: '',
                    selectedProperties: [],
                    addedProperties: [],
                    deletedProperties: []
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isReadOnly () {
                if (this.activeModel) {
                    return this.activeModel['bk_ispaused']
                }
                return false
            },
            objId () {
                return this.$route.params.modelId
            },
            sortedProperties () {
                const properties = []
                this.groupedProperties.forEach(group => {
                    group.properties.forEach(property => {
                        properties.push(property)
                    })
                })
                return properties
            },
            authority () {
                const cantEdit = ['process', 'plat']
                if (cantEdit.includes(this.objId)) {
                    return []
                }
                return this.$store.getters.admin ? ['search', 'update', 'delete'] : []
            }
        },
        async created () {
            const [properties, groups] = await Promise.all([this.getProperties(), this.getPropertyGroups()])
            this.init(properties, groups)
        },
        methods: {
            ...mapActions('objectModelFieldGroup', [
                'searchGroup',
                'updateGroup',
                'deleteGroup',
                'createGroup',
                'updatePropertyGroup'
            ]),
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            init (properties, groups) {
                properties = this.setPropertIndex(properties)
                groups = this.setGroupIndex(groups.concat({
                    'bk_group_index': Infinity,
                    'bk_group_id': 'none',
                    'bk_group_name': this.$t('Common["更多属性"]')
                }))
                const groupedProperties = groups.map(group => {
                    return {
                        info: group,
                        properties: properties.filter(property => property['bk_property_group'] === group['bk_group_id'])
                    }
                })
                this.groupedProperties = groupedProperties
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: this.objId,
                    config: {
                        requestId: `get_searchGroup_${this.objId}`,
                        cancelPrevious: true
                    }
                })
            },
            getProperties () {
                return this.searchObjectAttribute({
                    params: {
                        'bk_obj_id': this.objId,
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: `post_searchObjectAttribute_${this.objId}`,
                        cancelPrevious: true
                    }
                })
            },
            setGroupIndex (groups) {
                groups.sort((groupA, groupB) => groupA['bk_group_index'] - groupB['bk_group_index'])
                groups.forEach((group, index) => {
                    group['bk_group_index'] = index
                })
                return groups
            },
            setPropertIndex (properties) {
                properties.sort((propertyA, propertyB) => propertyA['bk_property_index'] - propertyB['bk_property_index'])
                properties.forEach((property, index) => {
                    property['bk_property_index'] = index
                })
                return properties
            },
            handleEditGroupName (group) {
                this.groupNameInEditing = group.info['bk_group_name']
                this.groupInEditing = group
                this.$nextTick(() => {
                    this.$refs.titleInput[0].focus()
                })
            },
            handleCancelEditGroupName () {
                this.groupInEditing = null
            },
            handleUpdateGroupName (group) {
                const isExist = this.groupedProperties.some(originalGroup => originalGroup !== group && originalGroup.info['bk_group_name'] === this.groupNameInEditing)
                if (isExist) {
                    this.$error(this.$t('ModelManagement["该名字已经存在"]'))
                    return
                }
                this.updateGroup({
                    params: {
                        condition: {
                            id: this.groupInEditing.info.id
                        },
                        data: {
                            'bk_group_name': this.groupNameInEditing
                        }
                    },
                    config: {
                        requestId: `put_updateGroup_name_${this.groupInEditing.info.id}`,
                        cancelPrevious: true
                    }
                })
                group.info['bk_group_name'] = this.groupNameInEditing
                this.groupInEditing = null
            },
            handleRiseGroup (index, group) {
                if (!index || ['none'].includes(group.info['bk_group_id'])) {
                    return
                }
                this.groupedProperties[index - 1]['info']['bk_group_index'] = index
                group['info']['bk_group_index'] = index - 1
                this.updateGroupIndex()
                this.resortGroups()
                this.updatePropertyIndex()
            },
            handleDropGroup (index, group) {
                if (index === (this.groupedProperties.length - 2) || ['none'].includes(group.info['bk_group_id'])) {
                    return
                }
                this.groupedProperties[index + 1]['info']['bk_group_index'] = index
                group.info['bk_group_index'] = index + 1
                this.updateGroupIndex()
                this.resortGroups()
                this.updatePropertyIndex()
            },
            handleAddProperty (group) {
                this.dialog.group = group
                this.dialog.selectedProperties = [...group.properties]
                this.dialog.isShow = true
                this.$nextTick(() => {
                    const $dialogProperty = this.$refs.dialogProperty
                    $dialogProperty.style.height = $dialogProperty.getBoundingClientRect().height + 'px'
                })
            },
            handleCancelAddProperty () {
                this.dialog.isShow = false
                this.dialog.selectedProperties = []
                this.dialog.addedProperties = []
                this.dialog.deletedProperties = []
                this.dialog.filter = ''
                this.dialog.group = null
                this.$nextTick(() => {
                    this.$refs.dialogProperty.style.height = 'auto'
                })
            },
            handleSelectProperty (property) {
                const selectedProperties = this.dialog.selectedProperties
                const addedProperties = this.dialog.addedProperties
                const deletedProperties = this.dialog.deletedProperties
                const selectedIndex = selectedProperties.indexOf(property)
                const addedIndex = addedProperties.indexOf(property)
                const deletedIndex = deletedProperties.indexOf(property)
                if (selectedIndex !== -1) {
                    selectedProperties.splice(selectedIndex, 1)
                    const isDeleteFromGroup = property['bk_property_group'] === this.dialog.group.info['bk_group_id']
                    if (isDeleteFromGroup && deletedIndex === -1) {
                        deletedProperties.push(property)
                    }
                    if (addedIndex !== -1) {
                        addedProperties.splice(addedIndex, 1)
                    }
                } else {
                    selectedProperties.push(property)
                    const isAddFromOtherGroup = property['bk_property_group'] !== this.dialog.group.info['bk_group_id']
                    if (isAddFromOtherGroup && addedIndex === -1) {
                        addedProperties.push(property)
                    }
                    if (deletedIndex !== -1) {
                        deletedProperties.splice(deletedIndex, 1)
                    }
                }
            },
            handleConfirmAddProperty () {
                const {
                    selectedProperties,
                    addedProperties,
                    deletedProperties
                } = this.dialog
                if (addedProperties.length || deletedProperties.length) {
                    this.groupedProperties.forEach(group => {
                        if (group === this.dialog.group) {
                            const resortedProperties = [...selectedProperties].sort((propertyA, propertyB) => propertyA['bk_property_index'] - propertyB['bk_property_index'])
                            group.properties = resortedProperties
                        } else {
                            const resortedProperties = group.properties.filter(property => !addedProperties.includes(property))
                            if (group.info['bk_group_id'] === 'none') {
                                Array.prototype.push.apply(resortedProperties, deletedProperties)
                            }
                            resortedProperties.sort((propertyA, propertyB) => propertyA['bk_property_index'] - propertyB['bk_property_index'])
                            group.properties = resortedProperties
                        }
                    })
                    this.updatePropertyIndex()
                }
                this.handleCancelAddProperty()
            },
            filter (property) {
                return property['bk_property_name'].toLowerCase().indexOf(this.dialog.filter.toLowerCase()) !== -1
            },
            handleDeleteGroup (group, index) {
                if (['default', 'none'].includes(group.info['bk_group_id'])) {
                    return
                }
                if (group.properties.length) {
                    this.$error('请先清空该分组下的字段')
                    return
                }
                this.deleteGroup({
                    id: group.info.id,
                    config: {
                        requestId: `delete_deleteGroup_${group.info.id}`,
                        fromCache: true
                    }
                }).then(() => {
                    this.groupedProperties.splice(index, 1)
                    this.$success(this.$t('Common["删除成功"]'))
                })
            },
            resortGroups () {
                this.groupedProperties.sort((groupA, groupB) => groupA.info['bk_group_index'] - groupB.info['bk_group_index'])
            },
            updateGroupIndex () {
                const groupToUpdate = this.groupedProperties.filter((group, index) => group.info['bk_group_index'] !== index && group.info['bk_group_id'] !== 'none')
                groupToUpdate.forEach(group => {
                    this.updateGroup({
                        params: {
                            condition: {
                                id: group.info.id
                            },
                            data: {
                                'bk_group_index': group.info['bk_group_index']
                            }
                        },
                        config: {
                            requestId: `put_updateGroup_index_${group.info.id}`,
                            cancelWhenRouteChange: false,
                            cancelPrevious: true
                        }
                    })
                })
            },
            handleDragChange (changeInfo) {
                if (changeInfo.hasOwnProperty('moved')) {
                    this.shouldUpdatePropertyIndex = changeInfo.moved.newIndex !== changeInfo.moved.oldIndex
                } else {
                    this.shouldUpdatePropertyIndex = true
                }
            },
            handleDragEnd () {
                if (this.shouldUpdatePropertyIndex) {
                    this.updatePropertyIndex()
                    this.shouldUpdatePropertyIndex = false
                }
            },
            updatePropertyIndex () {
                const properties = []
                let propertyIndex = 0
                this.groupedProperties.forEach(group => {
                    group.properties.forEach(property => {
                        if (property['bk_property_index'] !== propertyIndex || property['bk_property_group'] !== group.info['bk_group_id']) {
                            property['bk_property_index'] = propertyIndex
                            property['bk_property_group'] = group.info['bk_group_id']
                            properties.push(property)
                        }
                        propertyIndex++
                    })
                })
                if (!properties.length) {
                    return false
                }
                this.updatePropertyGroup({
                    params: properties.map(property => {
                        return {
                            condition: {
                                'bk_obj_id': this.objId,
                                'bk_property_id': property['bk_property_id'],
                                'bk_supplier_account': property['bk_supplier_account']
                            },
                            data: {
                                'bk_property_group': property['bk_property_group'],
                                'bk_property_index': property['bk_property_index']
                            }
                        }
                    }),
                    config: {
                        requestId: `put_updatePropertyGroup_${this.objId}`,
                        cancelWhenRouteChange: false
                    }
                })
            },
            handleAddGroup () {
                this.showAddGroup = true
                this.$nextTick(() => {
                    this.$refs.addGroupInput[0].focus()
                })
            },
            handleCreateGroup () {
                const groupedProperties = this.groupedProperties
                const isExist = groupedProperties.some(group => group.info['bk_group_name'] === this.newGroupName)
                if (isExist) {
                    this.$error(this.$t('ModelManagement["该名字已经存在"]'))
                    return
                }
                const groupId = Date.now().toString()
                this.createGroup({
                    params: {
                        'bk_group_id': groupId,
                        'bk_group_index': groupedProperties.length - 1,
                        'bk_group_name': this.newGroupName,
                        'bk_obj_id': this.objId,
                        'bk_supplier_account': this.supplierAccount
                    },
                    config: {
                        requestId: `post_createGroup_${groupId}`
                    }
                }).then(group => {
                    groupedProperties.splice(groupedProperties.length - 1, 0, {
                        info: group,
                        properties: []
                    })
                    this.handleCancelCreateGroup()
                })
            },
            handleCancelCreateGroup () {
                this.showAddGroup = false
                this.newGroupName = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    $modelHighlightColor: #3c96ff;
    .group-layout {
        height: 100%;
        padding: 0 20px 20px;
        @include scrollbar-y;
    }
    .group {
        margin: 28px 0 19px 0;
    }
    .group-header {
        .header-title {
            height: 21px;
            padding: 0 21px 0 13px;
            line-height: 21px;
            color: #333948;
            position: relative;
            font-size: 0;
            &:before {
                content: "";
                position: absolute;
                left: 0;
                top: 3px;
                width: 4px;
                height: 16px;
                background-color: $modelHighlightColor;
            }
            .title-input {
                width: 180px;
                height: 29px;
                line-height: 27px;
                margin: -4px 0 0 0;
            }
            .title-input-button {
                display: inline-block;
                margin: 0 0 0 14px;
                font-size: 14px;
                color: $modelHighlightColor;
            }
            .group-name {
                font-size: 16px;
                display: inline-block;
                vertical-align: middle;
            }
            .group-count {
                font-size: 16px;
                display: inline-block;
                vertical-align: middle;
                color: #c3cdd7;
            }
            .title-icon {
                display: none;
                vertical-align: middle;
                width: 21px;
                height: 21px;
                line-height: 24px;
                text-align: center;
                font-size: 12px;
                color: $modelHighlightColor;
                cursor: pointer;
            }
            &:hover .title-icon {
                display: inline-block;
            }
        }
        .header-options {
            font-size: 0;
            margin-right: -6px;
            .options-icon {
                display: inline-block;
                vertical-align: middle;
                width: 21px;
                height: 21px;
                margin: 0 0 0 6px;
                line-height: 21px;
                text-align: center;
                font-size: 12px;
                cursor: pointer;
                &.disabled {
                    color: #dde4eb;
                    cursor: not-allowed;
                }
            }
        }
    }
    .property-list {
        width: calc(100% + 10px);
        margin: 0 0 0 -5px;
        font-size: 14px;
        line-height: 36px;
        position: relative;
        &.empty {
            min-height: 70px;
        }
        .property-item {
            position: relative;
            width: calc(20% - 10px);
            padding: 0 0 0 17px;
            margin: 10px 5px;
            border: 1px solid #dde4eb;
            background-color: #f6f6f6;
            user-select: none;
            cursor: pointer;
            @include ellipsis;
            &:hover {
                &:before {
                    display: inline-block;
                }
            }
            &.sortable-ghost {
                background: #fff;
                color: #fff;
                border: 1px dashed $cmdbBorderFocusColor;
                &:before {
                    display: none;
                }
            }
            &:before {
                display: none;
                position: absolute;
                left: 7px;
                top: 16px;
                width: 4px;
                height: 4px;
                background-color: #8d9093;
                content: '';
                box-shadow: 0 -6px 0 0 #8d9093, 0 6px 0 0 #8d9093;
            }
        }
        .property-empty {
            position: absolute;
            top: 10px;
            left: 5px;
            width: calc(100% - 10px);
            height: 60px;
            line-height: 60px;
            border: 1px dashed #dde4eb;
            text-align: center;
            font-size: 14px;
            color: $modelHighlightColor;
            cursor: pointer;
            &.disabled {
                cursor: default;
                color: #aaa;
            }
        }
    }
    .add-group {
        margin: 20px 0 0 0;
        line-height: 29px;
        font-size: 0;
        .add-group-trigger {
            display: inline-block;
            vertical-align: middle;
            color: $modelHighlightColor;
            font-size: 16px;
            .icon {
                margin-left: 5px;
                display: inline-block;
                vertical-align: middle;
                font-size: 12px;
            }
        }
        .add-group-input {
            font-size: 16px;
            display: inline-block;
            vertical-align: middle;
            width: 180px;
            height: 29px;
            line-height: 27px;
        }
        .add-group-button {
            display: inline-block;
            vertical-align: middle;
            margin: 0 0 0 14px;
            font-size: 14px;
            color: $modelHighlightColor;
        }
    }
    .dialog-title {
        padding: 20px 13px;
        font-size: 20px;
        color: #333948;
    }
    .dialog-content {
        width: 470px;
        margin: 0 auto;
    }
    .dialog-filter {
        position: relative;
        input {
            padding-right: 40px;
        }
        .icon-search {
            position: absolute;
            right: 11px;
            top: 9px;
            font-size: 18px;
        }
    }
    .dialog-property {
        padding: 3px 29px;
        margin: 28px 0 0 0;
        max-height: 300px;
        @include scrollbar-y;
        .property-item {
            width: 50%;
            margin: 0 0 22px 0;
            .property-label {
                float: left;
                max-width: 100%;
                padding: 0 0 0 4px;
                line-height: 18px;
                cursor: pointer;
                @include ellipsis;
                &:before {
                    content: "";
                    display: inline-block;
                    vertical-align: -4px;
                    width: 18px;
                    height: 18px;
                    background: #fff url("../../../assets/images/checkbox-sprite.png") no-repeat;
                    background-position: 0 -62px;
                }
                &.checked:before {
                    background-position: -33px -62px;
                }
            }
        }
    }
</style>