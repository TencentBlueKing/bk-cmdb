<template>
    <div class="group-layout" v-bkloading="{ isLoading: $loading() }">
        <div class="layout-header">
            <bk-button @click="previewShow = true" :disabled="!preview.properties.length">{{$t('字段预览')}}</bk-button>
        </div>
        <div class="layout-content">
            <div class="group"
                v-for="(group, index) in groupedProperties"
                :key="index">
                <div class="group-header clearfix">
                    <div class="header-title fl">
                        <div class="group-name">
                            {{group.info['bk_group_name']}}
                            <span v-if="!isAdminView && group.info['bk_isdefault']">（{{$t('内置组不支持修改，排序')}}）</span>
                        </div>
                        <div class="title-icon-btn" v-if="updateAuth && isEditable(group.info) && group.info['bk_group_id'] !== 'none'">
                            <i class="title-icon icon icon-cc-edit"
                                @click="handleEditGroup(group)">
                            </i>
                            <i class="title-icon bk-icon icon-cc-delete"
                                :class="{ disabled: ['default'].includes(group.info['bk_group_id']) }"
                                @click="handleDeleteGroup(group, index)">
                            </i>
                        </div>
                    </div>
                </div>
                <vue-draggable class="property-list clearfix"
                    element="ul"
                    v-model="group.properties"
                    :options="{
                        group: 'property',
                        animation: 150,
                        filter: '.no-drag',
                        disabled: !updateAuth || !isEditable(group.info)
                    }"
                    :class="{
                        empty: !group.properties.length,
                        disabled: !updateAuth || !isEditable(group.info)
                    }"
                    @change="handleDragChange"
                    @end="handleDragEnd">
                    <li class="property-item fl"
                        v-for="(property, _index) in group.properties"
                        :class="{ 'only-ready': !isFieldEditable(property) }"
                        :key="_index"
                        :title="property['bk_property_name']"
                        @click="handleFieldDetailsView(!isFieldEditable(property), property)">
                        <span class="drag-logo"></span>
                        <div class="drag-content">
                            <div class="field-name">
                                <span>{{property['bk_property_name']}}</span>
                                <i v-if="property.isrequired">*</i>
                            </div>
                            <p>{{fieldTypeMap[property['bk_property_type']]}}</p>
                        </div>
                        <template v-if="isFieldEditable(property)">
                            <i class="property-icon icon icon-cc-edit mr10"
                                @click="handleEditField(group, property)">
                            </i>
                            <i class="property-icon bk-icon icon-cc-delete"
                                @click="handleDeleteField({ property, index, _index })">
                            </i>
                        </template>
                    </li>
                    <li class="property-add no-drag fl"
                        v-cursor="{
                            active: !updateAuth,
                            auth: [$OPERATION.U_MODEL]
                        }"
                        v-if="isEditable(group.info) && group.info['bk_group_id'] !== 'none'"
                        @click="handleAddField(group)">
                        <i class="bk-icon icon-plus"></i>
                        {{$t('添加')}}
                    </li>
                    <template v-if="!group.properties.length">
                        <li class="property-empty no-drag disabled" v-if="!(updateAuth && isEditable(group.info))">{{$t('暂无字段')}}</li>
                    </template>
                </vue-draggable>
                <template v-if="updateAuth && !activeModel['bk_ispaused']">
                    <div class="add-group" v-if="index === (groupedProperties.length - 1)">
                        <a class="add-group-trigger" href="javascript:void(0)"
                            @click="handleAddGroup">
                            {{$t('添加分组')}}
                            <i class="icon icon-cc-edit"></i>
                        </a>
                    </div>
                </template>
            </div>
        </div>

        <bk-dialog class="bk-dialog-no-padding"
            v-model="dialog.isShow"
            :mask-close="false"
            :width="600"
            @cancel="handleCancelAddProperty"
            @confirm="handleConfirmAddProperty">
            <div class="dialog-title" slot="tools">{{$t('新建字段')}}</div>
            <div class="dialog-content">
                <div class="dialog-filter">
                    <bk-input type="text" class="cmdb-form-input" v-model.trim="dialog.filter" right-icon="bk-icon icon-search"></bk-input>
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

        <bk-dialog class="bk-dialog-no-padding group-dialog"
            v-model="groupDialog.isShow"
            width="480"
            :mask-close="false">
            <div class="group-dialog-header" slot="tools">{{groupDialog.title}}</div>
            <div class="group-dialog-content">
                <label class="label-item">
                    <span>{{$t('分组名称')}}</span>
                    <span class="color-danger">*</span>
                    <div class="cmdb-form-item" :class="{ 'is-error': errors.has('groupName') }">
                        <bk-input v-model.trim="groupForm.groupName"
                            name="groupName"
                            v-validate="'required'">
                        </bk-input>
                        <p class="form-error">{{errors.first('groupName')}}</p>
                    </div>
                </label>
                <div class="label-item">
                    <span>{{$t('是否默认折叠')}}</span>
                    <div class="cmdb-form-item">
                        <bk-switcher theme="primary" v-model="groupForm.isFlod" size="small"></bk-switcher>
                    </div>
                </div>
            </div>
            <div class="group-dialog-footer" slot="footer">
                <bk-button theme="primary" @click="handleCreateGroup" v-if="groupDialog.type === 'create'">{{$t('确定')}}</bk-button>
                <bk-button theme="primary" @click="handleUpdateGroup" v-else>{{$t('确定')}}</bk-button>
                <bk-button @click="handleCancelGroupDialog">{{$t('取消')}}</bk-button>
            </div>
        </bk-dialog>

        <bk-sideslider
            :width="540"
            :title="slider.title"
            :is-show.sync="slider.isShow"
            :before-close="handleSliderBeforeClose">
            <the-field-detail
                ref="fieldForm"
                class="slider-content"
                slot="content"
                v-if="slider.isShow"
                :is-read-only="isReadOnly"
                :is-edit-field="slider.isEditField"
                :field="slider.curField"
                :group="slider.curGroup"
                :property-index="slider.propertyIndex"
                @save="handleFieldSave"
                @cancel="handleSliderBeforeClose">
            </the-field-detail>
        </bk-sideslider>

        <bk-sideslider
            :width="676"
            :title="$t('字段预览')"
            :is-show.sync="previewShow">
            <preview-field v-if="previewShow"
                slot="content"
                :properties="preview.properties"
                :property-groups="preview.groups">
            </preview-field>
        </bk-sideslider>

        <bk-sideslider
            :width="540"
            :title="$t('字段详情')"
            :is-show.sync="fieldDetailsDialog.isShow"
            @hidden="handleHideFieldDetailsView">
            <field-details-view v-if="fieldDetailsDialog.isShow"
                slot="content"
                :field="fieldDetailsDialog.field">
            </field-details-view>
        </bk-sideslider>
    </div>
</template>

<script>
    import vueDraggable from 'vuedraggable'
    import theFieldDetail from './field-detail'
    import previewField from './preview-field'
    import fieldDetailsView from './field-view'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        components: {
            vueDraggable,
            theFieldDetail,
            previewField,
            fieldDetailsView
        },
        props: {
            customObjId: String
        },
        data () {
            return {
                groupedProperties: [],
                shouldUpdatePropertyIndex: false,
                previewShow: false,
                fieldTypeMap: {
                    'singlechar': this.$t('短字符'),
                    'int': this.$t('数字'),
                    'float': this.$t('浮点'),
                    'enum': this.$t('枚举'),
                    'date': this.$t('日期'),
                    'time': this.$t('时间'),
                    'longchar': this.$t('长字符'),
                    'objuser': this.$t('用户'),
                    'timezone': this.$t('时区'),
                    'bool': 'bool'
                },
                dialog: {
                    isShow: false,
                    group: null,
                    filter: '',
                    selectedProperties: [],
                    addedProperties: [],
                    deletedProperties: []
                },
                groupDialog: {
                    isShow: false,
                    type: 'create',
                    title: this.$t('新建分组')
                },
                groupForm: {
                    groupName: '',
                    isFlod: true
                },
                slider: {
                    isShow: false,
                    title: this.$t('新建字段'),
                    isEditField: false,
                    curField: {},
                    curGroup: {},
                    propertyIndex: 0
                },
                preview: {
                    properties: [],
                    groups: []
                },
                fieldDetailsDialog: {
                    isShow: false,
                    field: {}
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'isAdminView', 'isBusinessSelected']),
            ...mapGetters('objectModel', ['isInjectable', 'activeModel']),
            objId () {
                return this.$route.params.modelId || this.customObjId
            },
            isReadOnly () {
                return this.activeModel && this.activeModel['bk_ispaused']
            },
            sortedProperties () {
                const propertiesSorted = this.isAdminView ? this.groupedProperties : this.metadataGroupedProperties
                let properties = []
                propertiesSorted.forEach(group => {
                    properties = properties.concat(group.properties)
                })
                return properties
            },
            groupedPropertiesCount () {
                const count = {}
                this.groupedProperties.forEach(({ info, properties }) => {
                    const groupId = info['bk_group_id']
                    count[groupId] = properties.length
                })
                return count
            },
            metadataGroupedProperties () {
                return this.groupedProperties.filter(group => !!this.$tools.getMetadataBiz(group.info))
            },
            updateAuth () {
                const cantEdit = ['process', 'plat']
                if (cantEdit.includes(this.$route.params.modelId)) {
                    return false
                }
                const editable = this.isAdminView || (this.isBusinessSelected && this.isInjectable)
                return editable && this.$isAuthorized(this.$OPERATION.U_MODEL)
            }
        },
        async created () {
            const [properties, groups] = await Promise.all([this.getProperties(), this.getPropertyGroups()])
            this.preview.properties = properties
            this.preview.groups = groups
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
            isFieldEditable (item) {
                if (item.ispre || this.isReadOnly || !this.updateAuth) {
                    return false
                }
                if (!this.isAdminView) {
                    return !!this.$tools.getMetadataBiz(item)
                }
                return true
            },
            isEditable (group) {
                if (this.isReadOnly) {
                    return false
                }
                if (this.isAdminView) {
                    return true
                }
                return !!this.$tools.getMetadataBiz(group)
            },
            canRiseGroup (index, group) {
                if (this.isAdminView) {
                    return index !== 0 && !['none'].includes(group.info['bk_group_id'])
                }
                const metadataIndex = this.metadataGroupedProperties.indexOf(group)
                return metadataIndex !== 0
            },
            canDropGroup (index, group) {
                if (this.isAdminView) {
                    return index !== (this.groupedProperties.length - 2) && !['none'].includes(group.info['bk_group_id'])
                }
                const metadataIndex = this.metadataGroupedProperties.indexOf(group)
                return metadataIndex !== (this.metadataGroupedProperties.length - 1)
            },
            async resetData () {
                const [properties, groups] = await Promise.all([this.getProperties(), this.getPropertyGroups()])
                this.preview.properties = properties
                this.preview.groups = groups
                this.init(properties, groups)
            },
            init (properties, groups) {
                properties = this.setPropertIndex(properties)
                groups = this.separateMetadataGroups(groups)
                groups = this.setGroupIndex(groups)
                const groupedProperties = groups.map(group => {
                    return {
                        info: group,
                        properties: properties.filter(property => {
                            if (['default', 'none'].includes(property['bk_property_group']) && group['bk_group_id'] === 'default') {
                                return true
                            }
                            return property['bk_property_group'] === group['bk_group_id']
                        })
                    }
                })
                this.groupedProperties = groupedProperties
            },
            getPropertyGroups () {
                return this.searchGroup({
                    objId: this.objId,
                    params: this.$injectMetadata({}, { inject: this.isInjectable }),
                    config: {
                        requestId: `get_searchGroup_${this.objId}`,
                        cancelPrevious: true
                    }
                })
            },
            getProperties () {
                return this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        'bk_obj_id': this.objId,
                        'bk_supplier_account': this.supplierAccount
                    }, {
                        inject: this.isInjectable
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_${this.objId}`,
                        cancelPrevious: true
                    }
                })
            },
            separateMetadataGroups (groups) {
                const publicGroups = []
                const metadataGroups = []
                groups.forEach(group => {
                    if (this.$tools.getMetadataBiz(group)) {
                        metadataGroups.push(group)
                    } else {
                        publicGroups.push(group)
                    }
                })
                publicGroups.sort((groupA, groupB) => {
                    return groupA['bk_group_index'] - groupB['bk_group_index']
                })
                metadataGroups.sort((groupA, groupB) => {
                    return groupA['bk_group_index'] - groupB['bk_group_index']
                })
                return [...publicGroups, ...metadataGroups]
            },
            setGroupIndex (groups) {
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
            handleEditGroup (group) {
                this.groupDialog.isShow = true
                this.groupDialog.type = 'update'
                this.groupDialog.title = this.$t('编辑分组')
                this.groupDialog.group = group
                this.groupForm.groupName = group.info['bk_group_name']
            },
            async handleUpdateGroup () {
                const isExist = this.groupedProperties.some(originalGroup => originalGroup !== this.groupDialog.group && originalGroup.info['bk_group_name'] === this.groupForm.groupName)
                if (isExist) {
                    this.$error(this.$t('该名字已经存在'))
                    return
                }
                await this.updateGroup({
                    params: this.$injectMetadata({
                        condition: {
                            id: this.groupDialog.group.info.id
                        },
                        data: {
                            'bk_group_name': this.groupForm.groupName
                        }
                    }, { inject: this.isInjectable }),
                    config: {
                        requestId: `put_updateGroup_name_${this.groupDialog.group.info.id}`,
                        cancelPrevious: true
                    }
                })
                this.groupDialog.group.info['bk_group_name'] = this.groupForm.groupName
                this.handleCancelGroupDialog()
            },
            handleAddGroup () {
                this.groupDialog.isShow = true
                this.groupDialog.type = 'create'
                this.groupDialog.title = this.$t('新建分组')
            },
            handleCancelGroupDialog () {
                this.groupDialog.isShow = false
                this.groupDialog.group = {}
                this.groupForm.groupName = ''
                this.groupForm.isFlod = true
            },
            async handleCreateGroup () {
                const groupedProperties = this.groupedProperties
                const isExist = groupedProperties.some(group => group.info['bk_group_name'] === this.groupForm.groupName)
                if (isExist) {
                    this.$error(this.$t('该名字已经存在'))
                    return
                } else if (!await this.$validator.validateAll()) {
                    return
                }
                const groupId = Date.now().toString()
                this.createGroup({
                    params: this.$injectMetadata({
                        'bk_group_id': groupId,
                        'bk_group_index': groupedProperties.length - 1,
                        'bk_group_name': this.groupForm.groupName,
                        'bk_obj_id': this.objId,
                        'bk_supplier_account': this.supplierAccount
                    }, {
                        inject: this.isInjectable
                    }),
                    config: {
                        requestId: `post_createGroup_${groupId}`
                    }
                }).then(group => {
                    groupedProperties.splice(groupedProperties.length - 1, 0, {
                        info: group,
                        properties: []
                    })
                    this.handleCancelGroupDialog()
                })
            },
            handleDeleteGroup (group, index) {
                if (['default', 'none'].includes(group.info['bk_group_id'])) {
                    return
                }
                if (group.properties.length) {
                    this.$error(this.$t('请先清空该分组下的字段'))
                    return
                }
                this.deleteGroup({
                    id: group.info.id,
                    config: {
                        requestId: `delete_deleteGroup_${group.info.id}`,
                        fromCache: true,
                        data: this.$injectMetadata({}, {
                            inject: this.isInjectable
                        })
                    }
                }).then(() => {
                    this.groupedProperties.splice(index, 1)
                    this.$success(this.$t('删除成功'))
                })
            },
            resortGroups () {
                this.groupedProperties.sort((groupA, groupB) => groupA.info['bk_group_index'] - groupB.info['bk_group_index'])
            },
            updateGroupIndex () {
                const groupToUpdate = this.groupedProperties.filter((group, index) => group.info['bk_group_index'] !== index && group.info['bk_group_id'] !== 'none')
                groupToUpdate.forEach(group => {
                    this.updateGroup({
                        params: this.$injectMetadata({
                            condition: {
                                id: group.info.id
                            },
                            data: {
                                'bk_group_index': group.info['bk_group_index']
                            }
                        }, {
                            inject: this.isInjectable
                        }),
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
                    params: this.$injectMetadata({
                        data: properties.map(property => {
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
                        })
                    }, { inject: this.isInjectable }),
                    config: {
                        requestId: `put_updatePropertyGroup_${this.objId}`,
                        cancelWhenRouteChange: false
                    }
                })
            },
            handleAddField (group) {
                if (!this.updateAuth) return
                this.slider.isEditField = false
                this.slider.curField = {}
                this.slider.curGroup = group.info
                this.slider.propertyIndex = group.properties.length
                this.slider.title = this.$t('新建字段')
                this.slider.isShow = true
            },
            handleEditField (group, property) {
                this.slider.isEditField = true
                this.slider.curField = property
                this.slider.curGroup = group.info
                this.slider.title = this.$t('编辑字段')
                this.slider.isShow = true
            },
            handleFieldSave () {
                this.resetData()
                this.slider.isShow = false
                this.slider.curField = {}
                this.slider.curGroup = {}
            },
            handleDeleteField ({ property: field, index, fieldIndex }) {
                this.$bkInfo({
                    title: this.$tc('确定删除字段？', field['bk_property_name'], { name: field['bk_property_name'] }),
                    confirmFn: async () => {
                        await this.$store.dispatch('objectModelProperty/deleteObjectAttribute', {
                            id: field.id,
                            config: {
                                data: this.$injectMetadata({}, {
                                    inject: this.isInjectable
                                }),
                                requestId: 'deleteObjectAttribute',
                                originalResponse: true
                            }
                        }).then(res => {
                            this.$http.cancel(`post_searchObjectAttribute_${this.activeModel['bk_obj_id']}`)
                            if (res.data.bk_error_msg === 'success' && res.data.bk_error_code === 0) {
                                this.groupedProperties[index].properties.splice(fieldIndex, 1)
                            }
                        })
                    }
                })
            },
            handleSliderBeforeClose () {
                const hasChanged = Object.keys(this.$refs.fieldForm.changedValues).length
                if (hasChanged) {
                    return new Promise((resolve, reject) => {
                        this.$bkInfo({
                            title: this.$t('确认退出'),
                            subTitle: this.$t('退出会导致未保存信息丢失'),
                            extCls: 'bk-dialog-sub-header-center',
                            confirmFn: () => {
                                this.slider.isShow = false
                                resolve(true)
                            },
                            cancelFn: () => {
                                resolve(false)
                            }
                        })
                    })
                }
                this.slider.isShow = false
                return true
            },
            handleFieldDetailsView (show, field) {
                if (!show) return
                this.fieldDetailsDialog.isShow = true
                this.fieldDetailsDialog.field = field
            },
            handleHideFieldDetailsView () {
                this.fieldDetailsDialog.isShow = false
                this.fieldDetailsDialog.field = {}
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
    .layout-header {
        margin: 10px 0 14px;
    }
    .group {
        margin-bottom: 19px;
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
                background-color: $cmdbBorderColor;
            }
            .title-input {
                width: 180px;
                display: inline-block;
                top: -5px;
                /deep/ .bk-form-input {
                    height: 28px;
                    line-height: 28px;
                }
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
                span {
                    @include inlineBlock;
                    font-size: 12px;
                    color: #c4c6cc;
                }
            }
            .group-count {
                font-size: 16px;
                display: inline-block;
                vertical-align: middle;
                color: #c3cdd7;
            }
            .title-icon-btn {
                @include inlineBlock;
                .icon-cc-edit {
                    margin: 0 4px;
                }
            }
            .title-icon {
                display: none;
                vertical-align: middle;
                width: 21px;
                height: 21px;
                line-height: 24px;
                text-align: center;
                font-size: 16px;
                color: $modelHighlightColor;
                cursor: pointer;
                &.disabled {
                    color: #C4C6CC;
                    cursor: not-allowed;
                }
            }
            &:hover .title-icon {
                display: inline-block;
            }
        }
    }
    .property-list {
        width: calc(100% + 10px);
        margin: 0 0 0 -5px;
        font-size: 14px;
        position: relative;
        &.empty {
            min-height: 70px;
        }
        &.disabled {
            .property-item {
                cursor: pointer;
            }
        }
        .property-item {
            display: flex;
            align-items: center;
            flex-wrap: wrap;
            position: relative;
            width: calc(20% - 10px);
            height: 59px;
            padding: 10px 10px 10px 14px;
            margin: 10px 5px;
            border: 1px solid #dcdee5;
            border-radius: 2px;
            background-color: #ffffff;
            user-select: none;
            cursor: move;
            &.only-ready {
                background-color: #f4f6f9;
            }
            &:hover {
                .drag-logo {
                    display: block;
                }
                .property-icon {
                    display: inline-block;
                }
                &::before {
                    display: block;
                }
            }
            &.sortable-ghost {
                height: 59px;
                background: #fff;
                color: #fff;
                border: 1px dashed $cmdbBorderFocusColor;
                &::before, .drag-content, .property-icon {
                    display: none;
                }
            }
            &::before {
                content: '';
                display: none;
                position: absolute;
                top: 25px;
                left: 5px;
                width: 2px;
                height: 2px;
                border-radius: 50%;
                background-color: #666666;
                box-shadow: 0 6px 0 0 #666666,
                    0 12px 0 0 #666666,
                    0 -6px 0 0 #666666,
                    6px 0 0 0 #3b3b3b,
                    6px 6px 0 0 #3b3b3b,
                    6px -6px 0 0 #3b3b3b,
                    6px 12px 0 0 #3b3b3b;
            }
            .drag-content {
                flex: 1;
                width: 0;
                color: #737987;
                margin-left: 6px;
                .field-name {
                    display: flex;
                    span {
                        line-height: 21px;
                        @include ellipsis;
                    }
                    i {
                        font-size: 16px;
                        font-style: normal;
                        font-weight: bold;
                        margin: 0 4px;
                    }
                }
                p {
                    font-size: 12px;
                    color: #c4c6cc;
                    @include ellipsis;
                }
            }
            .property-icon {
                color: #3a84ff;
                display: none;
                cursor: pointer;
            }
        }
        .property-add {
            width: calc(20% - 10px);
            height: 59px;
            line-height: 59px;
            margin: 10px 5px;
            text-align: center;
            border: 1px dashed #dcdee5;
            color: #979ba5;
            cursor: pointer;
            .icon-plus {
                font-weight: bold;
                margin-top: -2px;
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
            color: #3a84ff;
            font-size: 16px;
            .icon {
                margin: -2px 0 0 5px;
                display: inline-block;
                vertical-align: middle;
            }
        }
        .add-group-input {
            font-size: 0;
            display: inline-block;
            vertical-align: middle;
            width: 180px;
            /deep/ .bk-form-input {
                font-size: 14px;
                height: 30px;
                line-height: 30px;
            }
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
        padding: 0 0 20px 0;
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
    .group-dialog-header {
        color: #313237;
        font-size: 20px;
        padding: 18px 24px 14px;
    }
    .group-dialog-content {
        padding: 0 24px;
        .cmdb-form-item {
            margin: 10px 0 20px;
            &.is-error {
                /deep/ .bk-form-input {
                    border-color: #ff5656;
                }
            }
        }
    }
</style>
