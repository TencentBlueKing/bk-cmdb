<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div
    class="field-group"
    :class="{
      'is-dragging': isDragging,
      'is-readonly': !updateAuth
    }"
    v-bkloading="{ isLoading: $loading(Object.values(requestIds)), extCls: 'field-loading' }">
    <div class="field-options">
      <cmdb-auth v-if="isShowOptionBtn" :auth="authResources" @update-auth="handleReceiveAuth">
        <template #default="{ disabled }">
          <bk-button theme="primary" :disabled="disabled || activeModel.bk_ispaused"
            @click="handleAddField(displayGroupedProperties[0])">{{$t('新建字段')}}</bk-button>
        </template>
      </cmdb-auth>
      <bk-dropdown-menu class="dropdown-menu" :align="'left'" v-if="canBeImport">
        <template slot="dropdown-trigger">
          <bk-button class="import-btn">{{$t('导入')}}
            <i :class="['bk-icon icon-angle-down']"></i>
          </bk-button>
        </template>
        <ul class="bk-dropdown-list" slot="dropdown-content">
          <li class="bk-dropdown-item" v-if="hideImport">
            <cmdb-auth tag="label" class="label-btn"
              :auth="importAuth"
              :class="{ 'disabled': isReadOnlyImport }"
              @click="handleImport">
              <span class="label-btn-text">{{$t('导入字段')}}</span>
            </cmdb-auth>
          </li>
          <li class="bk-dropdown-item" @click="handleExport">
            <label class="label-btn">
              <span class="label-btn-text">{{$t('导出字段')}}</span>
            </label>
          </li>
        </ul>
      </bk-dropdown-menu>
      <cmdb-auth v-if="isShowOptionBtn" :auth="authResources" @update-auth="handleReceiveAuth">
        <template #default="{ disabled }">
          <bk-button :disabled="disabled || activeModel.bk_ispaused" @click="handleAddGroup">{{$t('新建分组')}}</bk-button>
        </template>
      </cmdb-auth>
      <bk-button @click="previewShow = true" :disabled="!properties.length">{{
        $t("字段预览")
      }}</bk-button>
      <bk-input
        class="filter-input"
        clearable
        right-icon="icon-search"
        :placeholder="$t('请输入关键字')"
        v-model.trim="keyword"
      >
      </bk-input>
      <div class="setting-btn" v-if="canEditSort" @click="configProperty.show = true">
        <i class="bk-icon left-icon icon-cog"></i>
      </div>
    </div>
    <div class="group-wrapper">
      <draggable
        class="group-list"
        tag="div"
        draggable=".group-item"
        ghost-class="group-item-ghost"
        group="group-list"
        handle=".collapse-group-title"
        :animation="200"
        :disabled="!updateAuth"
        @start="handleGroupDragStart"
        @end="handleGroupDragEnd"
        @change="handleGroupDragChange"
        v-model="displayGroupedProperties">
        <div
          class="group-item"
          v-for="(group, groupIndex) in displayGroupedProperties"
          :key="group.bk_classification_id"
          :class="[{ 'is-collapse': !groupCollapseState[group.info.bk_group_id] }]">
          <div class="group-header" slot="title">
            <collapse-group-title
              :drag-icon="updateAuth"
              :dropdown-menu="isEditable(group.info)"
              :collapse="groupCollapseState[group.info.bk_group_id]"
              :title="`${group.info.bk_group_name} ( ${group.properties.length} )`"
              @click.native="toggleGroup(group)"
              :commands="[
                {
                  text: $t('编辑分组'),
                  auth: authResources,
                  onUpdateAuth: handleReceiveAuth,
                  disabled: !isEditable(group.info) || group.info['bk_isdefault'],
                  handler: () => handleEditGroup(group)
                },
                {
                  text: $t('删除分组'),
                  disabled: !isEditable(group.info) || group.info['bk_isdefault'],
                  auth: authResources,
                  handler: () => handleDeleteGroup(group, groupIndex)
                }
              ]"
              v-bk-tooltips="{
                disabled: !isBuiltInGroup(group.info),
                content: $t('全局配置不可以在业务内调整'),
                placement: 'right'
              }">
            </collapse-group-title>
          </div>
          <bk-transition name="collapse" duration-type="ease">
            <draggable
              class="field-list"
              v-show="!groupCollapseState[group.info.bk_group_id]"
              tag="ul"
              v-model="group.properties"
              ghost-class="field-item-ghost"
              draggable=".field-item"
              group="field-list"
              :animation="150"
              :disabled="!updateAuth || !isEditable(group.info)"
              :class="{
                empty: !group.properties.length,
                disabled: !updateAuth || !isEditable(group.info)
              }"
              @start="handleModelDragStart"
              @end="handleModelDragEnd"
              @change="handleModelDragChange">
              <li
                class="field-item"
                v-for="(property, fieldIndex) in group.properties"
                :key="fieldIndex"
                @click="handleFieldDetailsView({ group, index: groupIndex, fieldIndex, property })">
                <field-card
                  :class="['field-card-container',{ 'only-ready': (!updateAuth || !isFieldEditable(property)) }]"
                  :field="property"
                  :field-unique="getFieldUnique(property)"
                  :deletable="false"
                  :only-ready="!updateAuth || !isFieldEditable(property)"
                  @remove-field="handleDeleteField({ property, index: groupIndex, fieldIndex })">
                  <template #action-append>
                    <cmdb-auth
                      class="mr10"
                      :auth="authResources"
                      @update-auth="handleReceiveAuth"
                      @click.native.stop>
                      <bk-button
                        slot-scope="{ disabled }"
                        class="field-button"
                        :text="true"
                        :disabled="disabled || !isFieldEditable(property, false)"
                        @click.stop="handleEditField(group, property)">
                        <i class="field-button-icon icon-cc-edit-shape"></i>
                      </bk-button>
                    </cmdb-auth>
                    <cmdb-auth
                      class="mr10"
                      @update-auth="handleReceiveAuth"
                      :auth="authResources"
                      @click.native.stop
                      v-if="!property.ispre">
                      <bk-button
                        slot-scope="{ disabled }"
                        class="field-button"
                        :text="true"
                        :disabled="disabled || !isFieldEditable(property) || isTemplateField(property)"
                        @click.stop="handleDeleteField({ property, index: groupIndex, fieldIndex })">
                        <i class="field-button-icon bk-icon icon-cc-del"></i>
                      </bk-button>
                    </cmdb-auth>
                  </template>
                  <template #tag-append v-if="isTemplateField(property)">
                    <div @mouseenter="(event) => handleTemplateTagHover(event, property)">
                      <mini-tag :text="$t('模板')" hover-bg-color="#C9F5E2" />
                    </div>
                  </template>
                </field-card>
              </li>
              <li class="field-add fl" v-if="isEditable(group.info)">
                <cmdb-auth v-if="isShowOptionBtn" @update-auth="handleReceiveAuth" :auth="authResources" tag="div">
                  <bk-button
                    slot-scope="{ disabled }"
                    class="field-add-btn"
                    :text="true"
                    :disabled="disabled"
                    @click.stop="handleAddField(group)">
                    <i class="bk-icon icon-plus"></i>
                    {{customObjId ? $t('新建业务字段') : $t('添加字段')}}
                  </bk-button>
                </cmdb-auth>
              </li>
              <li class="property-empty" v-if="!isEditable(group.info) && !group.properties.length">{{$t('暂无字段')}}</li>
            </draggable>
          </bk-transition>
        </div>
        <div class="add-group">
          <cmdb-auth v-if="isShowOptionBtn" @update-auth="handleReceiveAuth" :auth="authResources">
            <bk-button slot-scope="{ disabled }"
              class="add-group-trigger"
              :text="true"
              :disabled="disabled || activeModel.bk_ispaused"
              @click.stop="handleAddGroup">
              <i class="bk-icon icon-cc-plus"></i>
              {{customObjId ? $t('新建业务分组') : $t('添加分组')}}
            </bk-button>
          </cmdb-auth>
        </div>
      </draggable>
      <cmdb-data-empty
        v-if="displayGroupedProperties.length === 0"
        slot="empty"
        :stuff="dataEmpty"
        @clear="handleClearFilter">
      </cmdb-data-empty>
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
          <bk-input type="text" class="cmdb-form-input" clearable
            v-model.trim="dialog.filter" right-icon="bk-icon icon-search">
          </bk-input>
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
              :title="property.bk_property_name"
              @click="handleSelectProperty(property)">
              {{property.bk_property_name}}
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
              v-validate="'required|length:128'">
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
        <bk-button theme="primary"
          v-if="groupDialog.type === 'create'"
          :disabled="errors.has('groupName')"
          @click="handleCreateGroup">
          {{$t('提交')}}
        </bk-button>
        <bk-button theme="primary"
          v-else
          :disabled="errors.has('groupName')"
          @click="handleUpdateGroup">
          {{$t('保存')}}
        </bk-button>
        <bk-button @click="groupDialog.isShow = false">{{$t('取消')}}</bk-button>
      </div>
    </bk-dialog>

    <bk-sideslider
      ref="sidesliderComp"
      class="sides-slider"
      v-transfer-dom
      :width="640"
      :title="slider.title"
      :is-show.sync="slider.isShow"
      :before-close="slider.beforeClose"
      @hidden="handleSliderHidden">
      <field-details-view v-if="slider.isShow && slider.view === 'details'"
        slot="content"
        :field="slider.curField"
        :can-edit="updateAuth && isFieldEditable(slider.curField, false)"
        @on-edit="handleEditField(slider.curGroup, slider.curField)"
        @on-delete="handleDeleteField({
          property: slider.curField,
          index: slider.index,
          fieldIndex: slider.fieldIndex
        })">
      </field-details-view>
      <the-field-detail v-else-if="slider.isShow && slider.view === 'operation'"
        ref="fieldForm"
        slot="content"
        :is-main-line-model="isMainLineModel"
        :is-read-only="isReadOnly"
        :is-edit-field="slider.isEditField"
        :properties="properties"
        :field="slider.curField"
        :group="slider.curGroup"
        :groups="groupedProperties.map(item => item.info)"
        :custom-obj-id="customObjId"
        @save="handleFieldSave"
        @cancel="handleSliderBeforeClose">
      </the-field-detail>
    </bk-sideslider>

    <bk-sideslider
      ref="sidesliderComp"
      v-transfer-dom
      :width="676"
      :title="$t('字段预览')"
      :is-show.sync="previewShow">
      <preview-field v-if="previewShow"
        slot="content"
        :properties="properties"
        :property-groups="groups">
      </preview-field>
    </bk-sideslider>

    <bk-sideslider
      ref="sidesliderComp"
      v-transfer-dom
      :is-show.sync="configProperty.show"
      :width="676"
      :title="$t('实例表格字段排序设置')"
      :before-close="handleColumnsConfigSliderBeforeClose">
      <cmdb-columns-config slot="content"
        ref="cmdbColumnsConfig"
        v-if="configProperty.show"
        :properties="properties"
        :selected="configProperty.selected"
        :disabled-columns="disabledConfig"
        :show-reset="false"
        :confirm-text="$t('确定')"
        @on-cancel="handleColumnsConfigSliderBeforeClose"
        @on-apply="handleApplyConfig">
      </cmdb-columns-config>
    </bk-sideslider>
    <div class="field-bind-template-popover-content" ref="fieldBindTemplateRef"
      v-bkloading="{
        isLoading: $loading(bindTemplatePopover.requestId),
        size: 'mini',
        theme: 'primary',
        mode: 'spin'
      }">
      <i18n path="xx模板的字段" v-if="bindTemplatePopover.template.name">
        <template #template>
          <span class="template-name" @click.stop="handleViewTemplate(bindTemplatePopover.template.id)">
            {{ bindTemplatePopover.template.name }}
          </span>
        </template>
      </i18n>
    </div>
  </div>
</template>

<script>
  import Draggable from 'vuedraggable'
  import has from 'has'
  import qs from 'qs'
  import debounce from 'lodash.debounce'
  import theFieldDetail from './field-detail'
  import previewField from './preview-field'
  import fieldDetailsView from './field-view'
  import CmdbColumnsConfig from '@/components/columns-config/columns-config.vue'
  import { mapGetters, mapActions, mapState } from 'vuex'
  import { MENU_BUSINESS, MENU_MODEL_FIELD_TEMPLATE } from '@/dictionary/menu-symbol'
  import { BUILTIN_MODELS } from '@/dictionary/model-constants'
  import { v4 as uuidv4 } from 'uuid'
  import CollapseGroupTitle from '@/views/model-manage/children/collapse-group-title.vue'
  import { PROPERTY_TYPE_NAMES } from '@/dictionary/property-constants'
  import FieldCard from '@/components/model-manage/field-card.vue'
  import useUnique from '@/views/field-template/children/use-unique.js'
  import fieldTemplateService from '@/service/field-template'
  import MiniTag from '@/components/ui/other/mini-tag.vue'
  import { escapeRegexChar } from '@/utils/util'
  import { getUniqueProperties } from '@/components/filters/utils'

  export default {
    name: 'FieldGroup',
    components: {
      Draggable,
      theFieldDetail,
      previewField,
      fieldDetailsView,
      CmdbColumnsConfig,
      CollapseGroupTitle,
      FieldCard,
      MiniTag
    },
    props: {
      customObjId: String,
      canBeImport: Boolean,
      hideImport: Boolean,
      isReadOnlyImport: Boolean,
      importAuth: Object
    },
    data() {
      return {
        updateAuth: false,
        isDragging: false,
        properties: [],
        groups: [],
        groupedProperties: [],
        displayGroupedProperties: [],
        previewShow: false,
        keyword: '',
        groupCollapseState: {},
        initGroupState: {},
        fieldTypeMap: PROPERTY_TYPE_NAMES,
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
          view: 'details',
          isShow: false,
          title: this.$t('新建字段'),
          isEditField: false,
          curField: {},
          curGroup: {},
          group: {},
          beforeClose: null,
          index: null,
          fieldIndex: null,
          backView: ''
        },
        configProperty: {
          show: false,
          selected: []
        },
        requestIds: {
          properties: Symbol(),
          propertyGroups: Symbol()
        },
        dataEmpty: {
          type: 'search',
        },
        uniqueList: [],
        bindTemplatePopover: {
          show: false,
          instance: null,
          template: {},
          requestId: ''
        }
      }
    },
    computed: {
      ...mapState('userCustom', ['globalUsercustom']),
      ...mapGetters(['supplierAccount']),
      ...mapGetters('objectModel', ['activeModel']),
      ...mapGetters('objectMainLineModule', ['isMainLine']),
      isGlobalView() {
        const [topRoute] = this.$route.matched
        return topRoute ? topRoute.name !== MENU_BUSINESS : true
      },
      bizId() {
        if (this.isGlobalView) {
          return null
        }
        return parseInt(this.$route.params.bizId, 10)
      },
      objId() {
        return this.$route.params.modelId || this.customObjId
      },
      isReadOnly() {
        return this.activeModel && this.activeModel.bk_ispaused
      },
      sortedProperties() {
        const propertiesSorted = this.isGlobalView ? this.groupedProperties : this.bizGroupedProperties
        let properties = []
        propertiesSorted.forEach((group) => {
          properties = properties.concat(group.properties)
        })
        return properties
      },
      groupedPropertiesCount() {
        const count = {}
        this.groupedProperties.forEach(({ info, properties }) => {
          const groupId = info.bk_group_id
          count[groupId] = properties.length
        })
        return count
      },
      bizGroupedProperties() {
        return this.groupedProperties.filter(group => this.isBizCustomData(group.info))
      },
      curModel() {
        if (!this.objId) return {}
        return this.$store.getters['objectModelClassify/getModelById'](this.objId)
      },
      modelId() {
        return this.curModel.id || null
      },
      isMainLineModel() {
        return this.isMainLine(this.curModel)
      },
      authResources() {
        if (this.customObjId) { // 业务自定义字段
          return {
            type: this.$OPERATION.U_BIZ_MODEL_CUSTOM_FIELD,
            relation: [this.$store.getters['objectBiz/bizId']]
          }
        }
        return {
          relation: [this.modelId],
          type: this.$OPERATION.U_MODEL
        }
      },
      disabledConfig() {
        const disabled = {
          host: ['bk_host_id', 'bk_host_innerip', 'bk_host_innerip_v6', 'bk_cloud_id'],
          biz: ['bk_biz_name']
        }
        return disabled[this.objId] || ['bk_inst_name']
      },
      curGlobalCustomTableColumns() {
        return this.globalUsercustom[`${this.objId}_global_custom_table_columns`]
      },
      canEditSort() {
        return !this.customObjId
      },
      isShowOptionBtn() {
        return BUILTIN_MODELS.PROJECT !== this.$route.params.modelId
      }
    },
    watch: {
      groupedProperties: {
        handler() {
          this.filterField()
        },
        deep: true
      },
      keyword() {
        this.handleFilter()
      }
    },
    async created() {
      this.handleFilter = debounce(this.filterField, 300)
      const [properties, groups, uniqueList] = await Promise.all([
        this.getProperties(),
        this.getPropertyGroups(),
        this.getVerification()
      ])
      this.properties = properties
      this.groups = groups
      this.uniqueList = uniqueList
      this.init(properties, groups)

      this.$nextTick(() => {
        this.execRouteAction()
      })
    },
    beforeDestroy() {
      // 通过isShow=false在划开页面时仍然会出现未关闭的情况，因此直接调用组件内部方法关闭
      this.$refs?.sidesliderComp?.handleClose?.()
    },
    methods: {
      ...mapActions('objectModelFieldGroup', [
        'searchGroup',
        'updateGroup',
        'switchGroupIndex',
        'deleteGroup',
        'createGroup',
        'updatePropertyGroup',
        'updatePropertySort'
      ]),
      ...mapActions('objectModelProperty', ['searchObjectAttribute']),
      ...mapActions('objectUnique', [
        'searchObjectUniqueConstraints',
      ]),
      toggleGroup(group) {
        this.groupCollapseState[`${group.info.bk_group_id}`] = !this.groupCollapseState[`${group.info.bk_group_id}`]
      },
      isBizCustomData(data) {
        return has(data, 'bk_biz_id') && data.bk_biz_id > 0
      },
      isBuiltInGroup(group) {
        if (this.isGlobalView) {
          return false
        }
        return !this.isBizCustomData(group)
      },
      isFieldEditable(item, checkIspre = true) {
        if ((checkIspre && item.ispre) || this.isReadOnly || !this.updateAuth) {
          return false
        }
        if (!this.isGlobalView) {
          return this.isBizCustomData(item)
        }
        return true
      },
      isEditable(group) {
        if (this.isReadOnly) {
          return false
        }
        if (this.isGlobalView) {
          return true
        }
        return this.isBizCustomData(group)
      },
      isTemplateField(field) {
        return field.bk_template_id > 0
      },
      canRiseGroup(index, group) {
        if (this.isGlobalView) {
          return index !== 0
        }
        const customDataIndex = this.bizGroupedProperties.indexOf(group)
        return customDataIndex !== 0
      },
      canDropGroup(index, group) {
        if (this.isGlobalView) {
          return index !== (this.groupedProperties.length - 1)
        }
        const customDataIndex = this.bizGroupedProperties.indexOf(group)
        return customDataIndex !== (this.bizGroupedProperties.length - 1)
      },
      async resetData(filedId) {
        const [properties, groups, uniqueList] = await Promise.all([
          this.getProperties(),
          this.getPropertyGroups(),
          this.getVerification()
        ])
        if (filedId && this.slider.isShow) {
          const field = properties.find(property => property.bk_property_id === filedId)
          if (field) {
            this.slider.curField = field
          } else {
            this.handleSliderHidden()
          }
        }
        this.properties = properties
        this.groups = groups
        this.uniqueList = uniqueList
        this.init(properties, groups)
      },
      init(properties, groups) {
        properties = this.sortProperties(properties)
        const separatedGroups = this.separateBizCustomGroups(groups)
        const groupCollapseState = {}
        const groupedProperties = separatedGroups.map((group) => {
          groupCollapseState[group.bk_group_id] = group.is_collapse
          return {
            info: group,
            properties: properties.filter((property) => {
              if (['default', 'none'].includes(property.bk_property_group) && group.bk_group_id === 'default') {
                return true
              }
              return property.bk_property_group === group.bk_group_id
            })
          }
        })

        const seletedProperties = this.$tools.getHeaderProperties(properties, [], this.disabledConfig)
        const curGlobalCustomTableColumns = this.curGlobalCustomTableColumns
          ?.map(column => properties.find(prop => prop.bk_property_id === column))
          ?.filter(column => column)

        // 保证固定展示的列一定出现在已选择的数据中，无论之前的配置是什么
        this.configProperty.selected = getUniqueProperties(seletedProperties || [], curGlobalCustomTableColumns || [])
          .map(property => property.bk_property_id)

        this.initGroupState = this.$tools.clone(groupCollapseState)
        this.groupCollapseState = Object.assign({}, groupCollapseState, this.groupCollapseState)
        this.groupedProperties = groupedProperties
      },
      filterField() {
        if (this.keyword) {
          const reg = new RegExp(escapeRegexChar(this.keyword), 'i')
          const displayGroupedProperties = []
          this.groupedProperties.forEach((group) => {
            const matchedProperties = []
            group.properties.forEach((property) => {
              if (reg.test(property.bk_property_name) || reg.test(property.bk_property_id)) {
                matchedProperties.push(property)
              }
            })
            if (matchedProperties.length) {
              displayGroupedProperties.push({
                ...group,
                properties: matchedProperties
              })
            }
          })
          displayGroupedProperties.forEach((group) => {
            this.groupCollapseState[group.info.bk_group_id] = false
          })
          this.displayGroupedProperties = displayGroupedProperties
        } else {
          this.displayGroupedProperties = this.groupedProperties
        }
      },
      getPropertyGroups() {
        return this.searchGroup({
          objId: this.objId,
          params: this.isGlobalView ? {} : { bk_biz_id: this.bizId },
          config: {
            requestId: this.requestIds.propertyGroups,
            cancelPrevious: true
          }
        })
      },
      getProperties() {
        const params = {
          bk_obj_id: this.objId,
          bk_supplier_account: this.supplierAccount
        }
        if (!this.isGlobalView) {
          params.bk_biz_id = this.bizId
        }
        return this.searchObjectAttribute({
          params,
          config: {
            requestId: this.requestIds.properties,
            cancelPrevious: true
          },
          injectId: this.objId === 'host' ? 'host' : false
        })
      },
      getVerification() {
        return this.searchObjectUniqueConstraints({
          objId: this.objId,
          params: {},
          config: {
            requestId: 'searchObjectUniqueConstraints'
          }
        })
      },
      getFieldUnique(property) {
        if (this.uniqueList.length > 0) {
          const uniqueList = this.uniqueList.map(item => ({
            ...item,
            keys: item.keys.map(key => key.key_id)
          }))
          const { getUniqueByField } =  useUnique([], uniqueList)
          const { list: fieldUniqueList, type: fieldUniqueType } = getUniqueByField(property)
          const fieldUniqueWithNameList = fieldUniqueList
            .filter(item => item.keys.every(key => this.properties.find(({ id }) => id === key)))
            .map(item => ({
              ...item,
              names: item.keys.map(key => this.properties.find(field => field.id === key)?.bk_property_name)
            }))
          return {
            list: fieldUniqueWithNameList,
            type: fieldUniqueType
          }
        }
      },
      separateBizCustomGroups(groups) {
        const publicGroups = []
        const bizCustomGroups = []
        groups.forEach((group) => {
          if (this.isBizCustomData(group)) {
            bizCustomGroups.push(group)
          } else {
            publicGroups.push(group)
          }
        })
        publicGroups.sort((groupA, groupB) => groupA.bk_group_index - groupB.bk_group_index)
        bizCustomGroups.sort((groupA, groupB) => groupA.bk_group_index - groupB.bk_group_index)
        return [...publicGroups, ...bizCustomGroups]
      },
      sortProperties(properties) {
        properties.sort((propertyA, propertyB) => propertyA.bk_property_index - propertyB.bk_property_index)
        return properties
      },
      handleCancelAddProperty() {
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
      handleSelectProperty(property) {
        const { selectedProperties } = this.dialog
        const { addedProperties } = this.dialog
        const { deletedProperties } = this.dialog
        const selectedIndex = selectedProperties.indexOf(property)
        const addedIndex = addedProperties.indexOf(property)
        const deletedIndex = deletedProperties.indexOf(property)
        if (selectedIndex !== -1) {
          selectedProperties.splice(selectedIndex, 1)
          const isDeleteFromGroup = property.bk_property_group === this.dialog.group.info.bk_group_id
          if (isDeleteFromGroup && deletedIndex === -1) {
            deletedProperties.push(property)
          }
          if (addedIndex !== -1) {
            addedProperties.splice(addedIndex, 1)
          }
        } else {
          selectedProperties.push(property)
          const isAddFromOtherGroup = property.bk_property_group !== this.dialog.group.info.bk_group_id
          if (isAddFromOtherGroup && addedIndex === -1) {
            addedProperties.push(property)
          }
          if (deletedIndex !== -1) {
            deletedProperties.splice(deletedIndex, 1)
          }
        }
      },
      handleConfirmAddProperty() {
        const {
          selectedProperties,
          addedProperties,
          deletedProperties
        } = this.dialog
        if (addedProperties.length || deletedProperties.length) {
          this.groupedProperties.forEach((group) => {
            if (group === this.dialog.group) {
              // eslint-disable-next-line max-len
              const resortedProperties = [...selectedProperties].sort((propertyA, propertyB) => propertyA.bk_property_index - propertyB.bk_property_index)
              group.properties = resortedProperties
            } else {
              const resortedProperties = group.properties.filter(property => !addedProperties.includes(property))
              if (group.info.bk_group_id === 'none') {
                Array.prototype.push.apply(resortedProperties, deletedProperties)
              }
              resortedProperties.sort((A, B) => A.bk_property_index - B.bk_property_index)
              group.properties = resortedProperties
            }
          })
        }
        this.handleCancelAddProperty()
      },
      filter(property) {
        return property.bk_property_name.toLowerCase().indexOf(this.dialog.filter.toLowerCase()) !== -1
      },
      handleEditGroup(group) {
        this.groupDialog.isShow = true
        this.groupDialog.isShowContent = true
        this.groupDialog.type = 'update'
        this.groupDialog.title = this.$t('编辑分组')
        this.groupDialog.group = group
        this.groupForm.isCollapse = group.info.is_collapse
        this.groupForm.groupName = group.info.bk_group_name
      },
      async handleUpdateGroup() {
        const valid = await this.$validator.validate('groupName')
        if (!valid) {
          return
        }
        const curGroup = this.groupDialog.group
        // eslint-disable-next-line max-len
        const isExist = this.groupedProperties.some(originalGroup => originalGroup !== curGroup && originalGroup.info.bk_group_name === this.groupForm.groupName)
        if (isExist) {
          this.$error(this.$t('该名字已经存在'))
          return
        }
        const params = {
          condition: {
            id: curGroup.info.id
          },
          data: {
            bk_group_name: this.groupForm.groupName,
            is_collapse: this.groupForm.isCollapse
          }
        }
        if (!this.isGlobalView) {
          params.bk_biz_id = this.bizId
        }
        await this.updateGroup({
          params,
          config: {
            requestId: `put_updateGroup_name_${curGroup.info.id}`,
            cancelPrevious: true
          }
        })
        curGroup.info.bk_group_name = this.groupForm.groupName
        curGroup.info.is_collapse = this.groupForm.isCollapse
        this.groupCollapseState[curGroup.info.bk_group_id] = this.groupForm.isCollapse
        this.groupDialog.isShow = false
        this.$success(this.$t('修改成功'))
      },
      handleAddGroup() {
        this.groupDialog.isShow = true
        this.groupDialog.isShowContent = true
        this.groupDialog.type = 'create'
        this.groupDialog.title = this.$t('新建分组')
      },
      handleCancelGroupLeave() {
        this.groupDialog.group = {}
        this.groupForm.groupName = ''
        this.groupForm.isCollapse = false
        this.groupDialog.isShowContent = false
        this.groupDialog.isShow = false
      },
      async handleCreateGroup() {
        try {
          const valid = await this.$validator.validate('groupName')
          if (!valid) {
            return
          }
          const { groupedProperties } = this
          const isExist = groupedProperties.some(group => group.info.bk_group_name === this.groupForm.groupName)
          if (isExist) {
            this.$error(this.$t('该名字已经存在'))
            return
          }
          const latestIndex = Math.max(...groupedProperties.map(group => group.info.bk_group_index))
          const params = {
            bk_group_id: uuidv4(),
            bk_group_index: latestIndex + 1,
            bk_group_name: this.groupForm.groupName,
            bk_obj_id: this.objId,
            bk_supplier_account: this.supplierAccount,
            is_collapse: this.groupForm.isCollapse
          }
          if (!this.isGlobalView) {
            params.bk_biz_id = this.bizId
          }
          const group = await this.createGroup({ params })
          groupedProperties.push({
            info: group,
            properties: []
          })
          this.$set(this.groupCollapseState, group.bk_group_id, group.is_collapse)
          this.groupDialog.isShow = false
          this.$success(this.$t('创建成功'))
        } catch (err) {
          console.log(err)
        }
      },
      async handleDeleteGroup(group, index) {
        if (group.properties.length) {
          this.$error(this.$t('请先清空该分组下的字段'))
          return
        }
        await this.deleteGroup({
          id: group.info.id,
          config: {
            data: this.isGlobalView ? {} : { bk_biz_id: this.bizId }
          }
        })
        this.groupedProperties.splice(index, 1)
        this.$success(this.$t('删除成功'))
      },
      resortGroups() {
        this.groupedProperties.sort((groupA, groupB) => groupA.info.bk_group_index - groupB.info.bk_group_index)
      },
      handleGroupDragChange({ moved }) {
        const groupA = this.displayGroupedProperties[moved.oldIndex]
        const groupB = this.displayGroupedProperties[moved.newIndex]
        this.updateGroupIndex(groupA, groupB)
      },
      updateGroupIndex(groupA, groupB) {
        return this.switchGroupIndex({
          params: {
            condition: {
              id: [groupA.info.id, groupB.info.id],
            },
          },
          config: {
            requestId: 'put_updateGroup_index',
            cancelPrevious: true
          }
        }).then(() => {
          const groupAIndex = groupA.info.bk_group_index
          const groupBIndex = groupB.info.bk_group_index
          groupA.info.bk_group_index = groupBIndex
          groupB.info.bk_group_index = groupAIndex
          this.resortGroups()
          this.$success(this.$t('修改成功'))
        })
          .catch((err) => {
            console.log(err)
          })
      },
      handleModelDragChange(moveInfo) {
        if (has(moveInfo, 'moved') || has(moveInfo, 'added')) {
          const info = moveInfo.moved
            ? { ...moveInfo.moved }
            : { ...moveInfo.added }
          this.updatePropertyIndex(info)
        }
      },
      async updatePropertyIndex({ element: property, newIndex }) {
        let curIndex = 0
        let curGroup = ''
        const { bk_property_id: propertyId } = property
        const group = this.groupedProperties?.
          find(group => group?.properties?.
            find(item => item?.bk_property_id === propertyId))
        const len = group?.properties?.length || 0

        // 取移动字段新位置的前一个字段 index + 1，当给空字段组添加新字段时，curIndex 默认为 0
        if (newIndex > 0 && len !== 1) {
          // 拖拽插件bug 跨组拖动到最后的位置index会多1
          const index = newIndex === len ? newIndex - 2 : newIndex - 1
          curIndex = Number(group.properties[index].bk_property_index) + 1
        }
        curGroup = group.info.bk_group_id

        const params = {
          bk_property_group: curGroup,
          bk_property_index: curIndex
        }

        if (!this.isGlobalView) {
          params.bk_biz_id = this.bizId
        }

        try {
          await this.updatePropertySort({
            objId: this.objId,
            propertyId: property.id,
            params,
            config: {
              requestId: `updatePropertySort_${this.objId}`
            }
          })

          // 重新初始化字段及分组
          this.resetData()

          this.$success(this.$t('修改成功'))
        } catch (error) {
          console.log(error)
        }
      },
      handleAddField(group = {}) {
        this.slider.isEditField = false
        this.slider.curField = {}
        this.slider.curGroup = group.info
        this.slider.title = this.$t('新建字段')
        this.slider.isShow = true
        this.slider.beforeClose = this.handleSliderBeforeClose
        this.slider.view = 'operation'
      },
      handleEditField(group, property) {
        this.slider.isEditField = true
        this.slider.curField = property
        this.slider.curGroup = group.info || group
        this.slider.title = this.$t('编辑字段')
        this.slider.isShow = true
        this.slider.beforeClose = this.handleSliderBeforeClose
        this.slider.view = 'operation'
      },
      handleFieldSave(filedId) {
        this.handleBackView()
        this.resetData(filedId)
      },
      handleDeleteField({ property: field, index, fieldIndex }) {
        this.$bkInfo({
          title: this.$t('确定删除字段？'),
          subTitle: this.$t('删除模型字段提示', { property: field.bk_property_name, model: this.curModel.bk_obj_name }),
          confirmLoading: this.$loading('deleteObjectAttribute'),
          confirmFn: async () => {
            if (this.$loading('deleteObjectAttribute')) return false
            try {
              const res = await this.$store.dispatch('objectModelProperty/deleteObjectAttribute', {
                id: field.id,
                config: {
                  data: this.isGlobalView ? {} : { bk_biz_id: this.bizId },
                  requestId: 'deleteObjectAttribute',
                  originalResponse: true
                }
              })
              this.$http.cancel(`post_searchObjectAttribute_${this.activeModel.bk_obj_id}`)
              if (res.data.bk_error_msg === 'success' && res.data.bk_error_code === 0) {
                this.displayGroupedProperties[index].properties.splice(fieldIndex, 1)
                this.handleSliderHidden()
                this.$success(this.$t('删除成功'))
                this.resetData()
              }
            } catch (error) {
              console.log(error)
            }
          }
        })
      },
      handleBackView() {
        if (this.slider.backView === 'details') {
          this.handleFieldDetailsView({
            group: this.slider.group,
            index: this.slider.index,
            fieldIndex: this.slider.fieldIndex,
            property: this.slider.curField
          })
        } else {
          this.handleSliderHidden()
        }
      },
      handleSliderBeforeClose() {
        return this.$refs.fieldForm.beforeClose(this.handleBackView)
      },
      handleSliderHidden() {
        this.slider.isShow = false
        this.slider.curField = {}
        this.slider.beforeClose = null
        this.slider.backView = ''
      },
      handleFieldDetailsView({ group, index, fieldIndex, property }) {
        this.slider.isShow = true
        this.slider.curField = property
        this.slider.curGroup = group.info
        this.slider.group = group
        this.slider.view = 'details'
        this.slider.backView = 'details'
        this.slider.title = this.$t('字段详情')
        this.slider.index = index
        this.slider.fieldIndex = fieldIndex
        this.slider.beforeClose = null
      },
      handleReceiveAuth(auth) {
        this.updateAuth = auth
      },
      handleApplyConfig(properties) {
        const setProperties = properties.map(property => property.bk_property_id)
        this.$store.dispatch('userCustom/saveGlobalUsercustom', {
          objId: this.objId,
          params: {
            global_custom_table_columns: setProperties
          }
        }).then(() => {
          this.configProperty.selected = setProperties
          this.configProperty.show = false
        })
      },
      handleModelDragStart() {
        this.isDragging = true
      },
      handleModelDragEnd() {
        this.isDragging = false
      },
      handleGroupDragStart() {
        this.isDragging = true
      },
      handleGroupDragEnd() {
        this.isDragging = false
      },
      handleClearFilter() {
        this.keyword = ''
      },
      handleImport() {
        this.$emit('handleImportField')
      },
      handleExport() {
        this.$emit('exportField')
      },
      async handleTemplateTagHover(event, property) {
        this.bindTemplatePopover.instance?.destroy?.()
        this.bindTemplatePopover.template = {}
        this.bindTemplatePopover.instance = this.$bkPopover(event.target, {
          allowHTML: true,
          boundary: 'window',
          trigger: 'mouseenter',
          arrow: true,
          distance: 18,
          theme: 'light',
          interactive: true,
          animateFill: false,
          hideOnClick: false,
          content: this.$refs.fieldBindTemplateRef,
          onShow: () => {
            this.bindTemplatePopover.show = true
          },
          onHidden: () => {
            this.bindTemplatePopover.show = false
          }
        })

        this.bindTemplatePopover.instance.show()

        this.bindTemplatePopover.requestId = `${this.modelId}_${property.id}_${property.bk_template_id}`
        const params = {
          bk_template_id: property.bk_template_id,
          bk_attribute_id: property.id
        }
        const bindTemplate = await fieldTemplateService.getFieldBindTemplate(params, {
          requestId: this.bindTemplatePopover.requestId,
          fromCache: true
        })
        this.bindTemplatePopover.template = bindTemplate || {}
      },
      handleViewTemplate(id) {
        this.$routerActions.open({
          name: MENU_MODEL_FIELD_TEMPLATE,
          query: {
            id,
            action: 'view'
          }
        })
      },
      handleColumnsConfigSliderBeforeClose() {
        const refColumns = this.$refs.cmdbColumnsConfig
        if (!refColumns) {
          return
        }
        const { columnsChangedValues } = refColumns
        if (columnsChangedValues?.()) {
          refColumns.setChanged(true)
          return refColumns.beforeClose(() => {
            this.configProperty.show = false
          })
        }
        this.configProperty.show = false
      },
      execRouteAction() {
        const { action, payload } = this.$route.query
        const params = qs.parse(payload)
        if (action === 'view-field') {
          let fieldIndex
          let fieldGroup
          let groupIndex
          let property
          for (let i = 0; i < this.displayGroupedProperties.length; i++) {
            const groupItem = this.displayGroupedProperties[i]
            for (let j = 0; j < groupItem.properties.length; j++) {
              const fieldItem = groupItem.properties[j]
              if (fieldItem.bk_property_id === params.id) {
                groupIndex = i
                fieldIndex = j
                fieldGroup = groupItem
                property = fieldItem
                break
              }
            }
          }
          this.handleFieldDetailsView({
            group: fieldGroup,
            index: groupIndex,
            fieldIndex,
            property
          })
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
$modelHighlightColor: #3c96ff;
.field-group {
  height: 100%;
  padding: 20px;
  @include scrollbar-y;
}

.field-options {
  display: flex;
  position: sticky;
  top: -20px;
  z-index: 9999;
  background: #F5F7FA;
  margin: -20px -20px 0;
  padding: 20px 20px 14px;
  .bk-button {
    margin-right: 10px;
  }
  .filter-input {
    width: 480px;
    margin-left: auto;
  }
  .setting-btn {
    margin-left: 10px;
    height: 32px;
    color: #979ba5;
    border: 1px solid #c4c6cc;
    border-radius: 2px;
    width: 32px;
    display: flex;
    justify-content: center;
    align-items: center;
    background: white;
    cursor: pointer;
    i {
      font-size: 18px;
    }
  }
}

.group-wrapper {
  position: relative;
}

.group-item + .group-item{
  margin-top: 15px;
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
  display: flex;
  color: #313238;
  font-size: 14px;
}

.field-list {
  $ghostBorderColor: #dcdee5;
  $ghostBackgroundColor:#f5f7fa;
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  width: 100%;
  align-content: flex-start;
  margin-top: 7px;
  font-size: 14px;
  position: relative;
  .field-item {
    height: 60px;
    user-select: none;
    cursor: pointer;
    &-ghost {
      background-color: #f5f7fa !important;
      border: 1px dashed #dcdee5;

      &:hover {
        border-color: #dcdee5;
        background-color: #f5f7fa;
        box-shadow: none;
      }

      > * {
        display: none !important;
      }
    }
  }
  .field-add {
    height: 60px;
    .auth-box{
      display: block;
      height: 100%;
    }
    .field-add-btn {
      width: 100%;
      height: 100%;
      border: 1px dashed #dcdee5;
      background-color: #ffffff;
      border-radius: 2px;
      font-size: 12px;
      &:not(.is-disabled):hover {
        color: #3a84ff;
        border-color: #3a84ff;
      }
    }
    .icon-plus {
      font-weight: bold;
      margin-top: -4px;
      font-size: 16px;
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
  margin: 15px 0 0 0;
  font-size: 0;
  .add-group-trigger {
    color: #3A84FF;
    font-size: 14px;
    height: 30px;
    padding-right: 30px;
    line-height: 30px;
    text-align: left;
    padding-left: 2px;
    border-radius: 2px;
    &.is-disabled {
      color: #c4c6cc;
      .icon {
        color: #63656e;
      }
    }
    .icon-cc-plus {
      margin: -4px 2px 0 0;
      display: inline-block;
      vertical-align: middle;
      font-size: 16px;
    }
    &:not(.is-disabled):hover {
      color: #699df4;
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
  .field-item {
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
        background: #fff url("../../../assets/images/checkbox-sprite.png")
          no-repeat;
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
.group-dialog-footer {
  .bk-button + .bk-button {
    margin-left: 10px;
  }
}
.sides-slider {
  :deep(.slider-main) {
    padding:20px 40px;
  }
}
.dropdown-menu {
  .import-btn {
    font-size: 14px;
  }
}
.bk-dropdown-list {
  .bk-dropdown-item {
    display: block;
    height: 32px;
    line-height: 33px;
    padding: 0 16px;
    color: #63656e;
    font-size: var(--font-size);
    text-decoration: none;
    white-space: nowrap;
    &:hover {
      color: #3a84ff;
      background: #f0f1f5;
      cursor: pointer;
    }
    .label-btn-text {
      cursor: pointer;
      font-size: 14px;
    }
    .disabled {
      color: #c4c6cc;
    }
  }
}

.field-card-container {
  &:hover {
    .field-button {
      visibility: visible;
    }
  }
  .field-button {
    font-size: 0;
    visibility: hidden;
    color: #63656e;
    &:hover {
      color: #3a84ff;
    }
    .field-button-icon {
      font-size: 14px;
    }
    &.is-disabled {
      color: #c4c6cc;
    }
  }
  .flag-append {
    margin-left: 2px;
  }
  &.only-ready {
    background-color: #EAEBF0;
  }
}
.field-bind-template-popover-content {
  min-width: 120px;
  height: 15px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  .template-name {
    color: #3A84FF;
    cursor: pointer;
  }
}
</style>

<style lang="scss">
.field-loading {
  position: sticky !important;
}
</style>
