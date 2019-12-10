<template>
    <div class="group-layout" v-bkloading="{ isLoading: $loading(), extCls: 'field-loading' }">
        <div class="layout-header">
            <bk-button @click="previewShow = true" :disabled="!preview.properties.length">{{$t('字段预览')}}</bk-button>
        </div>
        <div class="layout-content">
            <div class="group"
                v-for="(group, index) in groupedProperties"
                :key="index">
                <cmdb-collapse
                    :collapse.sync="groupState[group.info['bk_group_id']]">
                    <div class="group-header clearfix" slot="title">
                        <div class="header-title fl">
                            <div class="group-name">
                                {{group.info['bk_group_name']}}
                                <span v-if="!isAdminView && isBuiltInGroup(group.info)">（{{$t('全局配置不可以在业务内调整')}}）</span>
                            </div>
                            <div class="title-icon-btn" v-if="!(!isAdminView && isBuiltInGroup(group.info))" @click.stop>
                                <cmdb-auth class="ml10" :auth="authResources" @update-auth="handleReceiveAuth">
                                    <bk-button slot-scope="{ disabled }"
                                        class="icon-btn"
                                        :text="true"
                                        :disabled="disabled || !isEditable(group.info)"
                                        @click.stop="handleEditGroup(group)">
                                        <i class="title-icon icon icon-cc-edit"></i>
                                    </bk-button>
                                </cmdb-auth>
                                <cmdb-auth class="ml5" :auth="authResources">
                                    <bk-button slot-scope="{ disabled }"
                                        class="icon-btn"
                                        :text="true"
                                        :disabled="disabled || !isEditable(group.info) || group.info['bk_isdefault']"
                                        @click.stop="handleDeleteGroup(group, index)">
                                        <i class="title-icon bk-icon icon-cc-delete"></i>
                                    </bk-button>
                                </cmdb-auth>
                                <cmdb-auth class="ml5" :auth="authResources">
                                    <bk-button slot-scope="{ disabled }"
                                        class="icon-btn"
                                        :text="true"
                                        :disabled="disabled || !isEditable(group.info) || !canRiseGroup(index, group)"
                                        @click.stop="handleRiseGroup(index, group)">
                                        <i class="title-icon bk-icon icon-arrows-up"></i>
                                    </bk-button>
                                </cmdb-auth>
                                <cmdb-auth class="ml5" :auth="authResources">
                                    <bk-button slot-scope="{ disabled }"
                                        class="icon-btn"
                                        :text="true"
                                        :disabled="disabled || !isEditable(group.info) || !canDropGroup(index, group)"
                                        @click.stop="handleDropGroup(index, group)">
                                        <i class="title-icon bk-icon icon-arrows-down"></i>
                                    </bk-button>
                                </cmdb-auth>
                            </div>
                        </div>
                    </div>
                    <template>
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
                                v-for="(property, fieldIndex) in group.properties"
                                :class="{ 'only-ready': !updateAuth || !isFieldEditable(property) }"
                                :key="fieldIndex"
                                @click="handleFieldDetailsView(!updateAuth || !isFieldEditable(property), property)">
                                <span class="drag-logo"></span>
                                <div class="drag-content">
                                    <div class="field-name">
                                        <span :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                                        <i v-if="property.isrequired">*</i>
                                    </div>
                                    <p>{{fieldTypeMap[property['bk_property_type']]}}</p>
                                </div>
                                <template v-if="!property['ispre']">
                                    <cmdb-auth class="mr10" :auth="authResources" @click.native.stop>
                                        <bk-button slot-scope="{ disabled }"
                                            class="property-icon-btn"
                                            :text="true"
                                            :disabled="disabled || !isFieldEditable(property)"
                                            @click.stop="handleEditField(group, property)">
                                            <i class="property-icon icon-cc-edit"></i>
                                        </bk-button>
                                    </cmdb-auth>
                                    <cmdb-auth class="mr10" :auth="authResources" @click.native.stop>
                                        <bk-button slot-scope="{ disabled }"
                                            class="property-icon-btn"
                                            :text="true"
                                            :disabled="disabled || !isFieldEditable(property)"
                                            @click.stop="handleDeleteField({ property, index, fieldIndex })">
                                            <i class="property-icon bk-icon icon-cc-delete"></i>
                                        </bk-button>
                                    </cmdb-auth>
                                </template>
                            </li>
                            <li class="property-add no-drag fl" v-if="isEditable(group.info)">
                                <cmdb-auth :auth="authResources" style="display: block;">
                                    <bk-button slot-scope="{ disabled }"
                                        class="property-add-btn"
                                        :text="true"
                                        :disabled="disabled"
                                        @click.stop="handleAddField(group)">
                                        <i class="bk-icon icon-plus"></i>
                                        {{customObjId ? $t('新建业务字段') : $t('添加')}}
                                    </bk-button>
                                </cmdb-auth>
                            </li>
                            <li class="property-empty" v-if="!isEditable(group.info) && !group.properties.length">{{$t('暂无字段')}}</li>
                        </vue-draggable>
                    </template>
                </cmdb-collapse>
                <div class="add-group" v-if="index === (groupedProperties.length - 1)">
                    <cmdb-auth :auth="authResources">
                        <bk-button slot-scope="{ disabled }"
                            class="add-group-trigger"
                            :text="true"
                            :disabled="disabled || activeModel['bk_ispaused']"
                            @click.stop="handleAddGroup">
                            <i class="bk-icon icon-plus-circle"></i>
                            {{customObjId ? $t('新建业务分组') : $t('添加分组')}}
                        </bk-button>
                    </cmdb-auth>
                </div>
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
                    <bk-input type="text" class="cmdb-form-input" clearable v-model.trim="dialog.filter" right-icon="bk-icon icon-search"></bk-input>
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
            :mask-close="false"
            @after-leave="handleCancelGroupLeave">
            <div class="group-dialog-header" slot="tools">{{groupDialog.title}}</div>
            <div class="group-dialog-content" v-if="groupDialog.isShowContent">
                <label class="label-item">
                    <span>{{$t('分组名称')}}</span>
                    <span class="color-danger">*</span>
                    <div class="cmdb-form-item" :class="{ 'is-error': errors.has('groupName') }">
                        <bk-input v-model.trim="groupForm.groupName"
                            :placeholder="$t('请输入xx', { name: $t('分组名称') })"
                            name="groupName"
                            v-validate="'required'">
                        </bk-input>
                        <p class="form-error">{{errors.first('groupName')}}</p>
                    </div>
                </label>
                <div class="label-item">
                    <span>{{$t('是否默认折叠')}}</span>
                    <div class="cmdb-form-item">
                        <bk-switcher theme="primary" v-model="groupForm.isCollapse" size="small"></bk-switcher>
                    </div>
                </div>
            </div>
            <div class="group-dialog-footer" slot="footer">
                <bk-button theme="primary" @click="handleCreateGroup" v-if="groupDialog.type === 'create'">{{$t('提交')}}</bk-button>
                <bk-button theme="primary" @click="handleUpdateGroup" v-else>{{$t('保存')}}</bk-button>
                <bk-button @click="groupDialog.isShow = false">{{$t('取消')}}</bk-button>
            </div>
        </bk-dialog>

        <bk-sideslider
            v-transfer-dom
            :width="540"
            :title="slider.title"
            :is-show.sync="slider.isShow"
            :before-close="handleSliderBeforeClose">
            <the-field-detail
                ref="fieldForm"
                slot="content"
                v-if="slider.isShow"
                :is-main-line-model="isMainLineModel"
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
            v-transfer-dom
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
            v-transfer-dom
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
                updateAuth: false,
                groupedProperties: [],
                shouldUpdatePropertyIndex: false,
                previewShow: false,
                groupState: {},
                initGroupState: {},
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
                    'bool': 'bool',
                    'list': this.$t('列表')
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
                    isShowContent: false,
                    isShow: false,
                    type: 'create',
                    title: this.$t('新建分组')
                },
                groupForm: {
                    groupName: '',
                    isCollapse: false
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
            curModel () {
                if (!this.objId) return {}
                return this.$store.getters['objectModelClassify/getModelById'](this.objId)
            },
            modelId () {
                return this.curModel.id || null
            },
            isMainLineModel () {
                return ['bk_host_manage', 'bk_biz_topo', 'bk_organization'].includes(this.curModel.bk_classification_id)
            },
            authResources () {
                return this.$authResources({
                    resource_id: this.modelId,
                    type: this.$OPERATION.U_MODEL
                })
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
            isBuiltInGroup (group) {
                if (group) {
                    return !(group.metadata && group.metadata.label && group.metadata.label.bk_biz_id)
                }
                return false
            },
            isFieldEditable (item) {
                if (item.ispre || this.isReadOnly) {
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
                    return index !== 0
                }
                const metadataIndex = this.metadataGroupedProperties.indexOf(group)
                return metadataIndex !== 0
            },
            canDropGroup (index, group) {
                if (this.isAdminView) {
                    return index !== (this.groupedProperties.length - 1)
                }
                const metadataIndex = this.metadataGroupedProperties.indexOf(group)
                return metadataIndex !== (this.metadataGroupedProperties.length - 1)
            },
            handleRiseGroup (index, group) {
                this.groupedProperties[index - 1]['info']['bk_group_index'] = index
                group['info']['bk_group_index'] = index - 1
                this.updateGroupIndex()
                this.resortGroups()
                this.updatePropertyIndex()
            },
            handleDropGroup (index, group) {
                this.groupedProperties[index + 1]['info']['bk_group_index'] = index
                group.info['bk_group_index'] = index + 1
                this.updateGroupIndex()
                this.resortGroups()
                this.updatePropertyIndex()
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
                const groupState = {}
                const groupedProperties = groups.map(group => {
                    groupState[group['bk_group_id']] = group['is_collapse']
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
                this.initGroupState = this.$tools.clone(groupState)
                this.groupState = Object.assign({}, groupState, this.groupState)
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
                this.groupDialog.isShowContent = true
                this.groupDialog.type = 'update'
                this.groupDialog.title = this.$t('编辑分组')
                this.groupDialog.group = group
                this.groupForm.isCollapse = group.info['is_collapse']
                this.groupForm.groupName = group.info['bk_group_name']
            },
            async handleUpdateGroup () {
                const curGroup = this.groupDialog.group
                const isExist = this.groupedProperties.some(originalGroup => originalGroup !== curGroup && originalGroup.info['bk_group_name'] === this.groupForm.groupName)
                if (isExist) {
                    this.$error(this.$t('该名字已经存在'))
                    return
                }
                await this.updateGroup({
                    params: this.$injectMetadata({
                        condition: {
                            id: curGroup.info.id
                        },
                        data: {
                            'bk_group_name': this.groupForm.groupName,
                            'is_collapse': this.groupForm.isCollapse
                        }
                    }, { inject: this.isInjectable }),
                    config: {
                        requestId: `put_updateGroup_name_${curGroup.info.id}`,
                        cancelPrevious: true
                    }
                })
                curGroup.info['bk_group_name'] = this.groupForm.groupName
                curGroup.info['is_collapse'] = this.groupForm.isCollapse
                this.groupState[curGroup.info.bk_group_id] = this.groupForm.isCollapse
                this.groupDialog.isShow = false
            },
            handleAddGroup () {
                this.groupDialog.isShow = true
                this.groupDialog.isShowContent = true
                this.groupDialog.type = 'create'
                this.groupDialog.title = this.$t('新建分组')
            },
            handleCancelGroupLeave () {
                this.groupDialog.group = {}
                this.groupForm.groupName = ''
                this.groupForm.isCollapse = false
                this.groupDialog.isShowContent = false
                this.groupDialog.isShow = false
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
                        'bk_group_index': groupedProperties.length + 1,
                        'bk_group_name': this.groupForm.groupName,
                        'bk_obj_id': this.objId,
                        'bk_supplier_account': this.supplierAccount,
                        'is_collapse': this.groupForm.isCollapse
                    }, {
                        inject: this.isInjectable
                    }),
                    config: {
                        requestId: `post_createGroup_${groupId}`
                    }
                }).then(group => {
                    groupedProperties.push({
                        info: group,
                        properties: []
                    })
                    this.groupState[group.bk_group_id] = group.is_collapse
                    this.groupDialog.isShow = false
                })
            },
            handleDeleteGroup (group, index) {
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
                const groupToUpdate = this.groupedProperties.filter((group, index) => group.info['bk_group_index'] !== index)
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
                const updateProperties = []
                const flatProperties = this.groupedProperties.reduce((acc, group) => acc.concat(group.properties), [])
                flatProperties.forEach((property, index) => {
                    if (property['bk_property_index'] !== index) {
                        property['bk_property_index'] = index
                        updateProperties.push(property)
                    }
                })
                if (!updateProperties.length) return
                this.updatePropertyGroup({
                    params: this.$injectMetadata({
                        data: updateProperties.map(property => {
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
            },
            handleReceiveAuth (auth) {
                this.updateAuth = auth
            }
        }
    }
</script>

<style lang="scss" scoped>
    $modelHighlightColor: #3c96ff;
    .group-layout {
        height: 100%;
        padding: 20px;
        @include scrollbar-y;
    }
    .layout-header {
        margin: 0 0 14px;
    }
    .group {
        margin-bottom: 19px;
    }
    /deep/ .collapse-layout {
        width: 100%;
        .collapse-trigger {
            display: flex;
            align-items: center;
        }
        .collapse-arrow {
            margin-right: 8px;
            color: #63656e;
        }
    }
    .group-header {
        .header-title {
            height: 21px;
            padding: 0 21px 0 0;
            line-height: 21px;
            color: #313237;
            position: relative;
            font-size: 0;
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
                font-weight: normal;
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
            }
            .icon-btn {
                @include inlineBlock;
                display: none;
                vertical-align: middle;
                font-size: 0;
                height: 21px;
                color: $modelHighlightColor;
                &.is-disabled {
                    color: #c4c6cc;
                }
            }
            .title-icon {
                font-size: 16px;
                width: 16px;
                height: 16px;
            }
            &:hover .icon-btn {
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
                &::before {
                    display: none !important;
                }
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
                border-color: #3a84ff;
                background-color: #f0f5ff;
                .drag-logo {
                    display: block;
                }
                .property-icon-btn {
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
                &::before, .drag-content, .property-icon-btn {
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
                    align-items: center;
                    span {
                        line-height: 21px;
                        @include ellipsis;
                    }
                    i {
                        font-size: 16px;
                        font-style: normal;
                        font-weight: bold;
                        margin: 4px 4px 0;
                        line-height: 7px;
                    }
                }
                p {
                    font-size: 12px;
                    color: #c4c6cc;
                    @include ellipsis;
                }
            }
            .property-icon-btn {
                font-size: 0;
                color: #3a84ff;
                display: none;
                .property-icon {
                    font-size: 14px;
                }
                &.is-disabled {
                    color: #C4C6CC;
                }
            }
        }
        .property-add {
            width: calc(20% - 10px);
            margin: 10px 5px;
            .property-add-btn {
                display: block;
                width: 100%;
                height: 59px;
                line-height: 59px;
                text-align: center;
                border: 1px dashed #dcdee5;
                color: #979ba5;
                background-color: #ffffff;
                &:not(.is-disabled):hover {
                    color: #3a84ff;
                    border-color: #3a84ff;
                }
            }
            .icon-plus {
                font-weight: bold;
                margin-top: -2px;
            }
        }
        .property-empty {
            width: calc(100% - 10px);
            height: 60px;
            line-height: 60px;
            border: 1px dashed #dde4eb;
            text-align: center;
            font-size: 14px;
            color: #aaaaaa;
            margin: 10px 0 10px 5px;
        }
    }
    .add-group {
        margin: 20px 0 0 0;
        line-height: 29px;
        font-size: 0;
        .add-group-trigger {
            color: #3a84ff;
            font-size: 16px;
            &.is-disabled {
                color: #C4C6cc;
                .icon {
                    color: #63656E;
                }
            }
            .icon-plus-circle {
                margin: -4px 2px 0 0;
                display: inline-block;
                vertical-align: middle;
            }
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

<style lang="scss">
    .field-loading {
        position: sticky !important;
    }
</style>
