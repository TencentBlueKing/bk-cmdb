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
  <bk-table
    header-cell-class-name="header-cell"
    ext-cls="property-config-table"
    :data="propertyRuleList"
    @selection-change="handleSelectionChange">
    <bk-table-column type="selection" width="60" align="center" v-if="deletable"></bk-table-column>
    <bk-table-column
      :width="readonly ? 300 : 200"
      :label="$t('字段名称')"
      prop="bk_property_name">
      <div slot-scope="{ row }" :class="{ ignore: row.__extra__.ignore }">
        <span>{{row.bk_property_name}}</span>
        <i class="property-name-tooltips icon-cc-tips"
          v-if="row.placeholder && isExclmationProperty(row.bk_property_type)"
          v-bk-tooltips.top="{
            theme: 'light',
            trigger: 'mouseenter',
            content: row.placeholder
          }">
        </i>
      </div>
    </bk-table-column>
    <bk-table-column
      v-if="multiple"
      :label="isModuleMode ? $t('已配置的模块') : $t('已配置的模板')"
      class-name="table-cell-module-path">
      <div slot-scope="{ row }" :class="{ ignore: row.__extra__.ignore }">
        <template v-if="isModuleMode">
          <div v-for="(id, index) in idList"
            v-show="showMore.expanded[row.id] || index < showMore.max"
            class="path-item"
            :key="index"
            :title="getModulePath(id)">
            {{getModulePath(id)}}
          </div>
        </template>
        <template v-else-if="isTemplateMode">
          <div v-for="(id, index) in idList"
            v-show="showMore.expanded[row.id] || index < showMore.max"
            class="path-item"
            :key="index"
            :title="getTemplateName(id)">
            {{getTemplateName(id)}}
          </div>
        </template>
        <div
          v-show="idList.length > showMore.max"
          :class="['show-more', { expanded: showMore.expanded[row.id] }]"
          @click="handleToggleExpanded(row.id)">
          {{showMore.expanded[row.id] ? $t('收起') : $t('展开更多')}}<i class="bk-cc-icon icon-cc-arrow-down"></i>
        </div>
      </div>
    </bk-table-column>
    <bk-table-column
      v-if="multiple || readonly"
      :width="readonly ? null : 200"
      :label="$t('当前值')"
      :show-overflow-tooltip="false">
      <template slot-scope="{ row }">
        <div v-if="multiple" :class="{ ignore: row.__extra__.ignore }">
          <template v-for="(id, index) in idList">
            <cmdb-property-value
              class="value-item"
              v-show="showMore.expanded[row.id] || index < showMore.max"
              :tag="'div'"
              :key="index"
              :value="getRuleValue(row.id, id)"
              :show-unit="false"
              :is-show-overflow-tips="isShowOverflowTips(row)"
              :property="row">
            </cmdb-property-value>
          </template>
          <div v-show="idList.length > showMore.max" class="show-more">&nbsp;</div>
        </div>
        <template v-else>
          <cmdb-property-value
            :class="['property-value', { disabled: !row.host_apply_enabled }]"
            :show-unit="false"
            :is-show-overflow-tips="isShowOverflowTips(row)"
            :value="row.__extra__.value"
            :property="row">
          </cmdb-property-value>
          <i class="bk-cc-icon icon-cc-tips disabled-tips"
            v-if="!row.host_apply_enabled"
            v-bk-tooltips="$t('该字段已被设为不可编辑，自动应用规则失效')"></i>
        </template>
      </template>
    </bk-table-column>
    <bk-table-column
      v-if="!readonly"
      :label="$t(multiple ? '修改后' : '值')"
      class-name="table-cell-form-element">
      <template slot-scope="{ row }">
        <div class="form-element-content">
          <property-form-element
            :property="row"
            v-bk-tooltips.top="{
              disabled: !row.placeholder || isExclmationProperty(row.bk_property_type),
              theme: 'light',
              trigger: 'click',
              content: row.placeholder
            }"
            @value-change="handlePropertyValueChange"
            @valid-change="handlePropertyValidChange">
          </property-form-element>
        </div>
      </template>
    </bk-table-column>
    <bk-table-column
      v-if="!readonly"
      width="180"
      :label="$t('操作')"
      :render-header="multiple ? (h, data) => renderColumnHeader(h, data, $t('忽略的字段在本次操作不会进行变更')) : null">
      <template slot-scope="{ row }">
        <bk-button theme="primary" text @click="handlePropertyRowDel(row)">
          <span v-if="multiple">{{$t(row.__extra__.ignore ? '恢复' : '忽略')}}</span>
          <span v-else>{{$t('删除')}}</span>
        </bk-button>
      </template>
    </bk-table-column>
    <bk-table-column
      v-if="readonly && showDelColumn"
      width="280"
      :label="$t('操作')">
      <template slot-scope="{ row }">
        <bk-button theme="primary" text @click="handleDeletePropertyRule(row)">
          {{$t('删除')}}
        </bk-button>
      </template>
    </bk-table-column>
    <cmdb-table-empty
      slot="empty"
      :stuff="table.stuff">
    </cmdb-table-empty>
  </bk-table>
</template>
<script>
  /* eslint-disable no-underscore-dangle */
  import { mapGetters, mapState } from 'vuex'
  import has from 'has'
  import propertyFormElement from '@/components/host-apply/property-form-element'
  import { CONFIG_MODE } from '@/service/service-template/index.js'
  import { PROPERTY_TYPES } from '@/dictionary/property-constants'
  import { getPropertyDefaultValue } from '@/utils/tools.js'
  import { isExclmationProperty } from '@/utils/util'

  export default {
    components: {
      propertyFormElement
    },
    props: {
      mode: {
        type: String,
        required: true
      },
      checkedPropertyIdList: {
        type: Array,
        default: () => ([])
      },
      idList: {
        type: Array,
        default: () => ([])
      },
      ruleList: {
        type: Array,
        default: () => ([])
      },
      multiple: {
        type: Boolean,
        default: false
      },
      readonly: {
        type: Boolean,
        default: false
      },
      deletable: {
        type: Boolean,
        default: false
      },
      showDelColumn: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        propertyRuleList: [],
        removeRuleIds: [],
        ignoreRuleIds: [],
        showMore: {
          max: 10,
          expanded: {}
        },
        table: {
          stuff: {
            type: 'default',
            payload: {
              emptyText: this.$t('bk.table.emptyText')
            }
          }
        }
      }
    },
    inject: {
      getModulePath: { default: null },
      getTemplateName: { default: null }
    },
    computed: {
      ...mapGetters('hostApply', ['configPropertyList']),
      ...mapState('hostApply', ['ruleDraft']),
      isModuleMode() {
        return this.mode === CONFIG_MODE.MODULE
      },
      isTemplateMode() {
        return this.mode === CONFIG_MODE.TEMPLATE
      },
      hasRuleDraft() {
        return Object.keys(this.ruleDraft).length > 0
      }
    },
    watch: {
      checkedPropertyIdList() {
        this.setPropertyRuleList()
      },
      configPropertyList() {
        this.setPropertyRuleList()
      },
      propertyRuleList() {
        this.setConfigData()
      }
    },
    created() {
      if (this.hasRuleDraft) {
        this.propertyRuleList = this.$tools.clone(this.ruleDraft.rules)
        const checkedPropertyIdList = this.propertyRuleList.map(item => item.id)
        this.$emit('update:checkedPropertyIdList', checkedPropertyIdList)
      } else {
        this.setPropertyRuleList()
      }
    },
    methods: {
      isExclmationProperty(type) {
        return isExclmationProperty(type)
      },
      setPropertyRuleList() {
        // 当前属性列表中不存在，则添加
        this.checkedPropertyIdList.forEach((id) => {
          // 原始主机属性对象
          const index = this.propertyRuleList.findIndex(property => id === property.id)
          if (index === -1) {
            const findProperty = this.configPropertyList.find(item => id === item.id)
            if (findProperty) {
              const property = this.$tools.clone(findProperty)
              // 初始化值
              if (this.multiple) {
                const rules = this.ruleList.filter(item => item.bk_attribute_id === property.id)
                property.__extra__.ruleList = rules
                if (rules?.length) {
                  // 配置值全部相同（未配置的在ruleList不存在即不参与判断）则初始化“修改后”的值为相同的那个配置值
                  const firstValue = rules?.[0]?.bk_property_value
                  const isSameValue = rules.every(rule => rule?.bk_property_value === firstValue)
                  property.__extra__.value = isSameValue ? firstValue : getPropertyDefaultValue(property)
                } else {
                  property.__extra__.value = getPropertyDefaultValue(property)
                }
              } else {
                const rule = this.ruleList.find(item => item.bk_attribute_id === property.id) || {}
                property.__extra__.ruleId = rule.id
                property.__extra__.value = has(rule, 'bk_property_value') ? rule.bk_property_value : getPropertyDefaultValue(property)
              }
              this.propertyRuleList.push(property)
            }
          }
        })

        // 删除或取消选择的，则去除
        // eslint-disable-next-line max-len
        this.propertyRuleList = this.propertyRuleList.filter(property => this.checkedPropertyIdList.includes(property.id))
      },
      setConfigData() {
        this.removeRuleIds = []
        this.ignoreRuleIds = []
        // 找出不存在于初始数据中的规则
        this.ruleList.forEach((rule) => {
          const findIndex = this.propertyRuleList.findIndex(property => property.id === rule.bk_attribute_id)
          if (findIndex === -1) {
            // 批量模式标记为忽略，单个标记为删除
            if (this.multiple) {
              this.ignoreRuleIds.push(rule.id)
            } else {
              this.removeRuleIds.push(rule.id)
            }
          }
        })
      },
      getRuleValue(attrId, targetId) {
        if (this.isModuleMode) {
          return (this.ruleList.find(rule => rule.bk_attribute_id === attrId && rule.bk_module_id === targetId) || {}).bk_property_value || ''
        }
        return (this.ruleList.find(rule => rule.bk_attribute_id === attrId && rule.service_template_id === targetId) || {}).bk_property_value || ''
      },
      reset() {
        if (!this.hasRuleDraft) {
          this.propertyRuleList = []
        }
      },
      renderColumnHeader(h, data, tips) {
        const directive = {
          content: tips,
          placement: 'bottom-start'
        }
        return <span>{ data.column.label } <i class="bk-cc-icon icon-cc-tips" v-bk-tooltips={ directive }></i></span>
      },
      handleToggleExpanded(id) {
        this.showMore.expanded = { ...this.showMore.expanded, ...{ [id]: !this.showMore.expanded[id] } }
      },
      handleSelectionChange(value) {
        this.$emit('selection-change', value)
      },
      handlePropertyRowDel(property) {
        if (this.multiple) {
          const extra = this.propertyRuleList.find(item => item.id === property.id)?.__extra__
          this.$set(extra, 'ignore', !extra.ignore)
        } else {
          const checkedIndex = this.checkedPropertyIdList.findIndex(id => id === property.id)
          this.checkedPropertyIdList.splice(checkedIndex, 1)
        }

        // 清理展开状态
        delete this.showMore.expanded[property.id]

        this.$emit('property-remove', property)
      },
      handlePropertyValueChange(value) {
        this.$emit('property-value-change', value)
      },
      handleDeletePropertyRule(property) {
        this.$emit('property-rule-delete', property)
      },
      handlePropertyValidChange() {
        this.$emit('property-valid-change')
      },
      isShowOverflowTips(property) {
        const complexTypes = [PROPERTY_TYPES.MAP, PROPERTY_TYPES.ENUMQUOTE, PROPERTY_TYPES.ORGANIZATION]
        return !complexTypes.includes(property.bk_property_type)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .form-element-content {
        width: 80%;
        padding-top: 8px;
        padding-bottom: 8px;
    }
    .path-item,
    .value-item {
        padding: 1px 0;
        @include ellipsis;
    }
    .path-item {
      direction: rtl;
      text-align: left;
    }
    .property-value {
      &.disabled {
        text-decoration: line-through;
      }
    }
    .show-more {
        color: #3a84ff;
        margin-top: 2px;
        cursor: pointer;
        .bk-cc-icon {
            font-size: 22px;
            margin-top: -2px;
        }
        &.expanded {
            .bk-cc-icon {
                transform: rotate(180deg);
            }
        }
    }
    .property-config-table {
      .cell {
        .disabled-tips {
          margin-top: 0;
          margin-left: 6px;
        }
        .ignore {
          color: $textDisabledColor !important;
        }
      }
    }
</style>
<style lang="scss">
    .property-config-table {
        overflow: unset !important;
        .bk-table-body-wrapper {
            overflow: unset !important;
        }

        .cell {
            .icon-cc-tips {
                margin-top: -2px;
            }
        }

        .table-cell-module-path:not(.header-cell) {
            .cell {
                padding-top: 8px;
                padding-bottom: 8px;
            }
        }

        .table-cell-form-element {
            .cell {
                overflow: unset !important;
                display: block !important;
            }
            .search-input-wrapper {
                position: relative;
            }
            .form-objuser .suggestion-list {
                z-index: 1000;
            }
        }
    }
</style>
