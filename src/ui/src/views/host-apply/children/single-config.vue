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
  <cmdb-sticky-layout class="config-wrapper">
    <template #header="{ sticky }">
      <top-steps :class="{ 'is-sticky': sticky }" :current="1" />
    </template>
    <div :class="['single-config', { 'is-loading': $loading(requestIds.rules) }]"
      v-bkloading="{ isLoading: $loading(requestIds.rules) }">

      <service-template-tips class="config-tips" v-if="isTemplateMode" :service-template-ids="ids" />

      <div class="config-body">
        <div :class="['choose-field', { 'not-choose': !checkedPropertyIdList.length }]">
          <div class="choose-hd">
            <span class="label">{{$t('自动应用字段')}}</span>
            <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
              <bk-button
                icon="plus"
                slot-scope="{ disabled }"
                :disabled="disabled"
                @click="handleChooseField">
                {{$t('选择字段')}}
              </bk-button>
            </cmdb-auth>
          </div>
          <div class="choose-bd" v-show="checkedPropertyIdList.length">
            <property-config-table
              ref="propertyConfigTable"
              :mode="mode"
              :checked-property-id-list.sync="checkedPropertyIdList"
              :rule-list="initRuleList"
              @property-value-change="handlePropertyValueChange">
            </property-config-table>
          </div>
        </div>
      </div>
    </div>
    <template #footer="{ sticky }">
      <div :class="['wrapper-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="nextButtonDisabled || disabled"
            @click="handleNextStep">
            {{$t('下一步')}}
          </bk-button>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </template>
    <host-property-modal
      :visible.sync="propertyModalVisible"
      :checked-list.sync="checkedPropertyIdList">
    </host-property-modal>
    <leave-confirm
      reverse
      :id="leaveConfirmConfig.id"
      :active="leaveConfirmConfig.active"
      :title="$t('是否退出配置')"
      :content="$t('启用步骤未完成，退出将会丢失当前配置')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')">
    </leave-confirm>
  </cmdb-sticky-layout>
</template>

<script>
  import { mapGetters, mapState } from 'vuex'
  import leaveConfirm from '@/components/ui/dialog/leave-confirm'
  import topSteps from './top-steps.vue'
  import hostPropertyModal from './host-property-modal'
  import propertyConfigTable from './property-config-table'
  import serviceTemplateTips from './service-template-tips.vue'
  import { MENU_BUSINESS_HOST_APPLY_CONFIRM } from '@/dictionary/menu-symbol'
  import { CONFIG_MODE } from '@/services/service-template/index.js'

  export default {
    name: 'single-config',
    components: {
      leaveConfirm,
      topSteps,
      hostPropertyModal,
      propertyConfigTable,
      serviceTemplateTips
    },
    props: {
      mode: {
        type: String,
        required: true
      },
      // 模块或模板id
      ids: {
        type: Array,
        default: () => ([])
      }
    },
    data() {
      return {
        initRuleList: [],
        checkedPropertyIdList: [],
        nextButtonDisabled: true,
        propertyModalVisible: false,
        leaveConfirmConfig: {
          id: 'singleConfig',
          active: true
        },
        requestIds: {
          rules: Symbol('rules')
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapState('hostApply', ['ruleDraft']),
      isModuleMode() {
        return this.mode === CONFIG_MODE.MODULE
      },
      isTemplateMode() {
        return this.mode === CONFIG_MODE.TEMPLATE
      },
      targetId() {
        return this.ids[0]
      },
      hasRuleDraft() {
        return Object.keys(this.ruleDraft).length > 0
      },
      targetIdsKey() {
        const targetIdsKeys = {
          [CONFIG_MODE.MODULE]: 'bk_module_ids',
          [CONFIG_MODE.TEMPLATE]: 'service_template_ids'
        }
        return targetIdsKeys[this.mode]
      },
      targetIdKey() {
        const targetIdKeys = {
          [CONFIG_MODE.MODULE]: 'bk_module_id',
          [CONFIG_MODE.TEMPLATE]: 'service_template_id'
        }
        return targetIdKeys[this.mode]
      },
      requestConfigs() {
        return {
          [this.requestIds.rules]: {
            [CONFIG_MODE.MODULE]: {
              action: 'hostApply/getRules',
              payload: {
                bizId: this.bizId,
                params: {
                  bk_module_ids: [this.targetId]
                }
              }
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/getTemplateRules',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  service_template_ids: [this.targetId]
                }
              }
            }
          }
        }
      }
    },
    watch: {
      checkedPropertyIdList() {
        this.toggleNextButtonDisabled()
      }
    },
    created() {
      this.initData()
    },
    methods: {
      async initData() {
        try {
          const ruleData = await this.getRules()
          this.initRuleList = ruleData.info || []
          const checkedPropertyIdList = this.initRuleList.map(item => item.bk_attribute_id)
          this.checkedPropertyIdList = this.hasRuleDraft ? [...this.checkedPropertyIdList] : checkedPropertyIdList
        } catch (e) {
          console.log(e)
        }
      },
      getRules() {
        const requestConfig = this.requestConfigs[this.requestIds.rules][this.mode]
        return this.$store.dispatch(requestConfig.action, {
          config: {
            requestId: this.requestIds.rules
          },
          ...requestConfig.payload
        })
      },
      toggleNextButtonDisabled() {
        this.$nextTick(() => {
          if (this.$refs.propertyConfigTable) {
            const { propertyRuleList } = this.$refs.propertyConfigTable
            const everyTruthy = propertyRuleList.every((property) => {
              // eslint-disable-next-line no-underscore-dangle
              const validTruthy = property.__extra__.valid !== false
              // eslint-disable-next-line no-underscore-dangle
              let valueTruthy = property.__extra__.value
              if (property.bk_property_type === 'bool') {
                valueTruthy = true
              } else if (property.bk_property_type === 'int') {
                valueTruthy = valueTruthy !== null && String(valueTruthy)
              }
              return valueTruthy && validTruthy
            })
            this.nextButtonDisabled = !this.checkedPropertyIdList.length || !everyTruthy
          }
        })
      },
      async handleNextStep() {
        const { propertyRuleList, removeRuleIds } = this.$refs.propertyConfigTable
        const additionalRules = propertyRuleList.map(property => ({
          bk_attribute_id: property.id,
          [this.targetIdKey]: this.targetId,
          // eslint-disable-next-line no-underscore-dangle
          bk_property_value: property.__extra__.value
        }))

        const savePropertyConfig = {
          // 配置对象列表
          [this.targetIdsKey]: [this.targetId],
          // 附加的规则
          additional_rules: additionalRules,
          // 删除的规则，来源于编辑表格删除
          remove_rule_ids: removeRuleIds
        }

        this.$store.commit('hostApply/setPropertyConfig', savePropertyConfig)
        this.$store.commit('hostApply/setRuleDraft', {
          rules: propertyRuleList
        })

        // 使离开确认失活
        this.leaveConfirmConfig.active = false
        this.$nextTick(function () {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_APPLY_CONFIRM,
            history: true
          })
        })
      },
      handlePropertyValueChange() {
        this.toggleNextButtonDisabled()
      },
      handleChooseField() {
        this.propertyModalVisible = true
      },
      handleCancel() {
        this.$routerActions.back()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .config-wrapper {
    max-height: 100%;
    @include scrollbar-y;
    .wrapper-footer {
      display: flex;
      align-items: center;
      height: 52px;
      padding: 0 20px;
      .bk-button {
        min-width: 86px;

        & + .bk-button {
          margin-left: 8px;
        }
      }
      .auth-box + .bk-button {
        margin-left: 8px;
      }
      &.is-sticky {
        background-color: #fff;
        border-top: 1px solid $borderColor;
      }
    }
  }
  .single-config {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 0 20px;

    &.is-loading {
      min-height: 160px;
      width: 100%;
    }

    .config-tips {
      margin-top: 12px;
    }

    .config-head,
    .config-foot {
      flex: none;
    }
    .config-body {
      width: 1066px;
      flex: auto;
    }
  }

  .choose-field {
    padding: 16px 2px;
    .choose-hd {
      .label {
        font-size: 14px;
        color: #63656e;
        margin-right: 8px;
      }
    }
    .choose-bd {
      margin-top: 20px;

      .form-element-content {
        padding: 4px 0;
      }
    }
  }
</style>
