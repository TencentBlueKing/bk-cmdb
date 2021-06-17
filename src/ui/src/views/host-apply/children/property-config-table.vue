<template>
  <bk-table
    header-cell-class-name="header-cell"
    ext-cls="property-config-table"
    :data="modulePropertyList"
    @selection-change="handleSelectionChange"
  >
    <bk-table-column type="selection" width="60" align="center" v-if="deletable"></bk-table-column>
    <bk-table-column
      :width="readonly ? 300 : 200"
      :label="$t('字段名称')"
      prop="bk_property_name"
    >
    </bk-table-column>
    <bk-table-column
      v-if="multiple"
      :label="$t('已配置的模块')"
      class-name="table-cell-module-path"
    >
      <template slot-scope="{ row }">
        <template v-for="(id, index) in moduleIdList">
          <div
            class="path-item"
            :key="index" v-show="showMore.expanded[row.id] || index < showMore.max"
            :title="$parent.getModulePath(id)"
          >
            {{$parent.getModulePath(id)}}
          </div>
        </template>
        <div
          v-show="moduleIdList.length > showMore.max"
          :class="['show-more', { expanded: showMore.expanded[row.id] }]"
          @click="handleToggleExpanded(row.id)"
        >
          {{showMore.expanded[row.id] ? $t('收起') : $t('展开更多')}}<i class="bk-cc-icon icon-cc-arrow-down"></i>
        </div>
      </template>
    </bk-table-column>
    <bk-table-column
      v-if="multiple || readonly"
      :width="readonly ? null : 200"
      :label="$t('当前值')"
      show-overflow-tooltip>
      <template slot-scope="{ row }">
        <template v-if="multiple">
          <template v-for="(id, index) in moduleIdList">
            <cmdb-property-value
              class="value-item"
              v-show="showMore.expanded[row.id] || index < showMore.max"
              :tag="'div'"
              :key="index"
              :value="getRuleValue(row.id, id)"
              :show-unit="false"
              :property="row">
            </cmdb-property-value>
          </template>
          <div v-show="moduleIdList.length > showMore.max" class="show-more">&nbsp;</div>
        </template>
        <template v-else>
          <cmdb-property-value
            :class="['property-value', { disabled: !row.host_apply_enabled }]"
            :show-unit="false"
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
      class-name="table-cell-form-element"
    >
      <template slot-scope="{ row }">
        <div class="form-element-content">
          <property-form-element :property="row" @value-change="handlePropertyValueChange"></property-form-element>
        </div>
      </template>
    </bk-table-column>
    <bk-table-column
      v-if="!readonly"
      width="180"
      :label="$t('操作')"
      :render-header="multiple ? (h, data) => renderColumnHeader(h, data, $t('删除操作不影响原有配置')) : null"
    >
      <template slot-scope="{ row }">
        <bk-button theme="primary" text @click="handlePropertyRowDel(row)">{{$t('删除')}}</bk-button>
      </template>
    </bk-table-column>
  </bk-table>
</template>
<script>
  import { mapGetters, mapState } from 'vuex'
  import has from 'has'
  import propertyFormElement from '@/components/host-apply/property-form-element'
  export default {
    components: {
      propertyFormElement
    },
    props: {
      checkedPropertyIdList: {
        type: Array,
        default: () => ([])
      },
      moduleIdList: {
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
      }
    },
    data() {
      return {
        modulePropertyList: [],
        removeRuleIds: [],
        ignoreRuleIds: [],
        showMore: {
          max: 10,
          expanded: {}
        }
      }
    },
    computed: {
      ...mapGetters('hostApply', ['configPropertyList']),
      ...mapState('hostApply', ['ruleDraft']),
      hasRuleDraft() {
        return Object.keys(this.ruleDraft).length > 0
      }
    },
    watch: {
      checkedPropertyIdList() {
        this.setModulePropertyList()
      },
      configPropertyList() {
        this.setModulePropertyList()
      },
      modulePropertyList() {
        this.setConfigData()
      }
    },
    created() {
      if (this.hasRuleDraft) {
        this.modulePropertyList = this.$tools.clone(this.ruleDraft.rules)
        const checkedPropertyIdList = this.modulePropertyList.map(item => item.id)
        this.$emit('update:checkedPropertyIdList', checkedPropertyIdList)
      } else {
        this.setModulePropertyList()
      }
    },
    methods: {
      setModulePropertyList() {
        // 当前模块属性列表中不存在，则添加
        this.checkedPropertyIdList.forEach((id) => {
          // 原始主机属性对象
          const moduleIndex = this.modulePropertyList.findIndex(property => id === property.id)
          if (moduleIndex === -1) {
            const findProperty = this.configPropertyList.find(item => id === item.id)
            if (findProperty) {
              const property = this.$tools.clone(findProperty)
              // 初始化值
              if (this.multiple) {
                // eslint-disable-next-line no-underscore-dangle
                property.__extra__.ruleList = this.ruleList.filter(item => item.bk_attribute_id === property.id)
                // eslint-disable-next-line no-underscore-dangle
                property.__extra__.value = this.getPropertyDefaultValue(property)
              } else {
                const rule = this.ruleList.find(item => item.bk_attribute_id === property.id) || {}
                // eslint-disable-next-line no-underscore-dangle
                property.__extra__.ruleId = rule.id
                // eslint-disable-next-line no-underscore-dangle
                property.__extra__.value = has(rule, 'bk_property_value') ? rule.bk_property_value : this.getPropertyDefaultValue(property)
              }
              this.modulePropertyList.push(property)
            }
          }
        })

        // 删除或取消选择的，则去除
        // eslint-disable-next-line max-len
        this.modulePropertyList = this.modulePropertyList.filter(property => this.checkedPropertyIdList.includes(property.id))
      },
      setConfigData() {
        this.removeRuleIds = []
        this.ignoreRuleIds = []
        // 找出不存在于初始数据中的规则
        this.ruleList.forEach((rule) => {
          const findIndex = this.modulePropertyList.findIndex(property => property.id === rule.bk_attribute_id)
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
      getRuleValue(attrId, moduleId) {
        return (this.ruleList.find(rule => rule.bk_attribute_id === attrId && rule.bk_module_id === moduleId) || {}).bk_property_value || ''
      },
      getPropertyDefaultValue(property) {
        let value = ''
        if (property.bk_property_type === 'bool') {
          value = false
        }
        return value
      },
      reset() {
        if (!this.hasRuleDraft) {
          this.modulePropertyList = []
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
        const checkedIndex = this.checkedPropertyIdList.findIndex(id => id === property.id)
        this.checkedPropertyIdList.splice(checkedIndex, 1)

        // 清理展开状态
        delete this.showMore.expanded[property.id]
      },
      handlePropertyValueChange(value) {
        this.$emit('property-value-change', value)
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
        height: 20px;
        line-height: 20px;
        @include ellipsis;
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
