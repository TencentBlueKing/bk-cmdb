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
  <div class="apply-layout">
    <cmdb-tips
      :tips-style="{
        background: 'none',
        border: 'none',
        fontSize: '12px',
        lineHeight: '30px',
        padding: 0
      }"
      :icon-style="{
        color: '#63656E',
        fontSize: '14px',
        lineHeight: '30px'
      }">
      {{$t('转移属性变化确认提示')}}
    </cmdb-tips>
    <div class="options-row">
      <div class="option-label">
        <span class="has-tips" v-bk-tooltips="$t('属性值与目标模块配置不一致的主机')">{{$t('是否更新主机属性')}}</span>
      </div>
      <bk-radio-group class="option-content"
        v-model="updateOption.changed"
        @change="handleUpdateOptionChange">
        <bk-radio :value="true">
          {{$t('是将把转移的主机更新为目标模块配置')}}
        </bk-radio>
        <bk-radio :value="false">
          {{$t('否将保留主机原有配置')}}
        </bk-radio>
      </bk-radio-group>
    </div>
    <div class="options-row conflict" v-if="hasMultipleConflictModule && updateOption.changed">
      <div class="option-label">
        <span class="has-tips" v-bk-tooltips="$t('目标模块配置了不同的自动应用属性，需要重新指定配置')">{{$t('冲突字段配置')}}</span>
      </div>
      <div class="option-content conflict-property-list">
        <template v-for="item in conflictPropertyList">
          <div class="conflict-property-item" :key="item.bk_attribute_id" v-if="item.host_apply_rules.length > 1">
            <div class="property-name" :title="item.propertyName">{{item.propertyName}}</div>
            <bk-select v-model="item.selectedRuleId" class="property-value-selector"
              ext-popover-cls="host-apply-property-value-selector-popover"
              :clearable="false"
              :placeholder="$t('请选择目标模块配置')"
              :searchable="item.host_apply_rules.length > 10"
              @change="handleModuleRuleChange">
              <bk-option v-for="rule in item.host_apply_rules"
                :key="rule.id"
                :id="rule.id"
                :name="getModulePath(rule.bk_module_id)"
                @mouseenter.native="(event) => handleModuleRuleHover(event, rule, item.bk_attribute_id)"
                @mouseleave.native="(event) => handleModuleRuleHover(event, rule, item.bk_attribute_id)">
                <div class="bk-option-content-default">
                  <div class="bk-option-name medium-font">
                    {{getModulePath(rule.bk_module_id)}}
                  </div>
                </div>
              </bk-option>
            </bk-select>
          </div>
        </template>
      </div>
    </div>
    <property-confirm-table class="mt10"
      ref="confirmTable"
      max-height="auto"
      :list="list">
    </property-confirm-table>
    <div class="module-popover" ref="modulePopoverEl" v-show="moduleRulePopover.show">
      <div class="path">{{moduleRulePopover.path}}</div>
      <div :class="['rule-list', { 'flex-col': moduleRulePopover.ruleList.length < 10 }]">
        <div v-for="rule in moduleRulePopover.ruleList" :key="rule.id"
          :class="['rule-item', { current: moduleRulePopover.currentPropertyId === rule.bk_attribute_id }]">
          <div class="property-name" :title="rule.property.bk_property_name">{{rule.property.bk_property_name}}</div>
          <cmdb-property-value
            class="property-value"
            :show-unit="false"
            :value="rule.bk_property_value"
            :property="rule.property">
          </cmdb-property-value>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
  export default {
    name: 'host-attrs-auto-apply',
    components: {
      propertyConfirmTable
    },
    props: {
      info: {
        type: Array,
        required: true
      }
    },
    data() {
      return {
        updateOption: {
          changed: true,
          final_rules: []
        },
        ruleList: [],
        hostPropertyList: [],
        moduleRulePopover: {
          path: '',
          currentPropertyId: null,
          ruleList: [],
          instance: null,
          show: false
        }
      }
    },
    inject: {
      getModulePath: { default: null }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      list() {
        return this.info.filter(item => item.unresolved_conflict_count > 0)
      },
      conflictPropertyList() {
        const conflictPropertyList = []

        // 去重归并所有冲突字段
        this.list.forEach((row) => {
          row.conflicts.forEach((conflict) => {
            if (!conflictPropertyList.some(item => item.bk_property_id === conflict.bk_property_id)) {
              const property = this.hostPropertyList.find(item => item.bk_property_id === conflict.bk_property_id)

              // 取第一个配置规则作为默认选中项
              const defaultRuleId = conflict.host_apply_rules[0].id

              conflictPropertyList.push({
                ...conflict,
                defaultRuleId,
                propertyName: property?.bk_property_name,
                selectedRuleId: ''
              })
            }
          })
        })

        return conflictPropertyList
      },
      conflictRuleList() {
        const ruleList = []

        this.conflictPropertyList.forEach((conflict) => {
          conflict.host_apply_rules.forEach((rule) => {
            if (!ruleList.some(id => id === rule.id)) {
              ruleList.push(rule)
            }
          })
        })

        return ruleList
      },
      moduleIdList() {
        const moduleIdList = []

        this.conflictRuleList.forEach((rule) => {
          if (!moduleIdList.some(id => id === rule.bk_module_id)) {
            moduleIdList.push(rule.bk_module_id)
          }
        })

        return moduleIdList
      },
      hasMultipleConflictModule() {
        // 冲突规则配置存在于多个模块间且每个规则多于一条配置
        const hasMoreThanOne = this.conflictPropertyList.some(conflict => conflict.host_apply_rules.length > 1)
        return this.moduleIdList.length > 1 && hasMoreThanOne
      }
    },
    watch: {
      moduleIdList() {
        this.getModuleFinalRules()
      },
      conflictPropertyList() {
        this.initUpdateOption()
      }
    },
    created() {
      this.getHostPropertyList()
    },
    methods: {
      initUpdateOption() {
        this.conflictPropertyList.forEach((item) => {
          // 将默认项赋值给选中项，触发一次change方法确保目标值与选中的规则值保持一致
          item.selectedRuleId = item.defaultRuleId

          const currentRule = item.host_apply_rules.find(rule => rule.id === item.selectedRuleId)
          this.updateOption.final_rules.push({
            id: item.selectedRuleId, // 规则id
            bk_attribute_id: item.bk_attribute_id,
            bk_property_value: currentRule.bk_property_value
          })
        })
      },
      async getModuleFinalRules() {
        try {
          const rules = await this.$store.dispatch('hostApply/getModuleFinalRules', {
            params: {
              bk_biz_id: this.bizId,
              bk_module_ids: this.moduleIdList
            }
          })
          this.ruleList = rules ?? []
        } catch (e) {
          console.error(e)
        }
      },
      async getHostPropertyList() {
        try {
          const data = await this.$store.dispatch('hostApply/getProperties', {
            params: { bk_biz_id: this.bizId },
            config: {
              requestId: 'getHostPropertyList',
              fromCache: true
            }
          })
          this.hostPropertyList = data ?? []
        } catch (e) {
          console.error(e)
        }
      },
      handleModuleRuleChange(ruleId) {
        // 通过规则id找到完整的规则配置
        const rule = this.conflictRuleList.find(rule => rule.id === ruleId)

        // 更新列表中的目标值
        this.list.forEach((item) => {
          // 通过规则绑定的字段id找到要更新的目标字段
          const targetField = item.update_fields.find(field => field.bk_attribute_id === rule.bk_attribute_id)

          // 更新目标值，值为规则配置的值
          targetField.bk_property_value = rule.bk_property_value
        })

        // 更新冲突字段配置选项，将字段对应的规则id同步为所选规则
        const targetOption = this.updateOption.final_rules.find(item => item.bk_attribute_id === rule.bk_attribute_id)
        targetOption.id = ruleId
        targetOption.bk_property_value = rule.bk_property_value
      },
      handleModuleRuleHover(event, rule, propertyId) {
        const { bk_module_id: moduleId } = rule
        const moduleRuleList = this.ruleList.filter(item => item.bk_module_id === moduleId)

        this.moduleRulePopover.path = this.getModulePath(moduleId)
        this.moduleRulePopover.currentPropertyId = propertyId
        this.moduleRulePopover.ruleList = moduleRuleList.map((rule) => {
          const property = this.hostPropertyList.find(item => item.id === rule.bk_attribute_id)
          // 将完整的属性数据加入
          return { ...rule, property }
        })

        if (this.moduleRulePopover.instance) {
          this.moduleRulePopover.instance.destroy()
        }

        this.moduleRulePopover.instance = this.$bkPopover(event.target, {
          content: this.$refs.modulePopoverEl,
          delay: 300,
          hideOnClick: true,
          placement: 'right',
          animateFill: false,
          theme: 'light',
          boundary: 'window',
          trigger: 'manual',
          arrow: true,
          onShow: () => {
            this.moduleRulePopover.show = true
          },
          onHide: () => {
            this.moduleRulePopover.show = false
          }
        })

        if (event.type === 'mouseenter') {
          this.moduleRulePopover.instance.show()
        } else {
          this.moduleRulePopover.instance.hide()
        }
      },
      handleUpdateOptionChange(changed) {
        if (changed) {
          // 更新列表中的目标值为上次选择的规则配置值
          this.list.forEach((item) => {
            item.update_fields.forEach((field) => {
              const rule = this.updateOption.final_rules.find(rule => rule.bk_attribute_id === field.bk_attribute_id)
              field.bk_property_value = rule.bk_property_value
            })
          })
        } else {
          // 更新列表中的目标值为主机当前值
          this.list.forEach((item) => {
            item.update_fields.forEach((field) => {
              const hostCurrent = item.conflicts.find(conflict => conflict.bk_attribute_id === field.bk_attribute_id)
              field.bk_property_value = hostCurrent.bk_property_value
            })
          })
        }
      },
      getHostApplyConflictResolvers() {
        if (this.updateOption.changed) {
          return {
            changed: true,
            final_rules: this.updateOption.final_rules?.slice()
          }
        }
        return {
          changed: false
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
  .options-row {
    display: flex;
    align-items: center;
    margin: 24px 0;

    .option-label {
      position: relative;
      flex: none;
      width: 130px;
      font-size: 14px;
      white-space: nowrap;
      margin-right: 18px;
      text-align: right;

      .has-tips {
        display: inline-block;
        border-bottom: 1px dashed #979ba5;
        line-height: 20px;
        cursor: default;
      }

      &::after {
        content: "*";
        position: absolute;
        top: 5px;
        right: -10px;
        color: red;
        font-size: 12px;
      }
    }
    .option-content {
      font-size: 14px;

      ::v-deep .bk-form-radio {
        & + .bk-form-radio {
          margin-left: 16px;
        }
      }
    }

    &.conflict {
      align-items: baseline;
      margin-bottom: 16px;
    }
  }

  .conflict-property-list {
    display: flex;
    width: 910px;
    flex-wrap: wrap;
  }
  .conflict-property-item {
    display: flex;
    margin-bottom: 8px;
    .property-name {
      font-size: 12px;
      color: #63656E;
      background: #f5f7fa;
      border: 1px solid #c4c6cc;
      border-radius: 2px 0px 0px 2px;
      line-height: 30px;
      padding: 0 16px;
      width: 140px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .property-value-selector {
      width: 300px;
      margin-left: -1px;
      border-top-left-radius: 0;
      border-bottom-left-radius: 0;

      ::v-deep {
        .bk-select-name {
          direction: rtl;
          text-align: left;
        }
      }
    }

    &:nth-child(2n) {
      margin-left: 12px;
    }
  }

  .module-popover {
    max-width: 460px;
    padding: 8px;

    .path {
      font-size: 12px;
      color: #313238;
      margin-bottom: 4px;
    }
    .rule-list {
      display: flex;
      flex-wrap: wrap;

      &.flex-col {
        flex-direction: column;
        .rule-item {
          .property-value {
            width: 260px;
          }
        }
      }

      .rule-item {
        display: flex;
        margin: 2px 0;

        .property-name {
          font-size: 12px;
          color: #979ba5;
          width: 110px;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;

          &::after {
            content: '：';
          }
        }
        .property-value {
          font-size: 12px;
          color: #313238;
          width: 110px;
          text-overflow: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }

        &.current {
          .property-name,
          .property-value {
            color: #ff9c01;
          }
        }
      }
    }
  }
</style>
<style lang="scss">
  .host-apply-property-value-selector-popover {
    .bk-option-name {
      direction: rtl;
      text-align: left;
    }
  }
</style>
