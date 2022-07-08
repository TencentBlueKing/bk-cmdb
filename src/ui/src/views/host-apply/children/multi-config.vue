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
    <template #header="{ sticky }" v-if="!isDel">
      <top-steps :class="{ 'is-sticky': sticky }" :current="1" />
    </template>
    <div class="multi-config" v-bkloading="{ isLoading: $loading(requestIds.rules) }">

      <service-template-tips class="config-tips" v-if="isTemplateMode" :service-template-ids="ids" />

      <div class="config-bd">
        <div class="config-item">
          <div class="item-label">
            <i18n path="已选择N个模块：" v-if="isModuleMode">
              <template #count><span>{{ids.length}}</span></template>
            </i18n>
            <i18n path="已选择N个模板：" v-else-if="isTemplateMode">
              <template #count><span>{{ids.length}}</span></template>
            </i18n>
          </div>
          <div class="item-content">
            <div :class="['target-list', { 'show-more': showMore.isMoreShowed }]" ref="targetList">
              <template v-if="isModuleMode">
                <div
                  v-for="(id, index) in ids" :key="index"
                  class="target-item"
                  v-bk-tooltips="getModulePath(id)">
                  <span class="target-icon">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
                  {{getModuleName(id)}}
                </div>
              </template>
              <template v-else-if="isTemplateMode">
                <div
                  v-for="(id, index) in ids" :key="index"
                  class="target-item"
                  v-bk-tooltips="getTemplateName(id)">
                  <span class="target-icon">{{$i18n.locale === 'en' ? 'M' : '模'}}</span>
                  {{getTemplateName(id)}}
                </div>
              </template>
              <div
                :class="['target-item', 'more', { 'opened': showMore.isMoreShowed }]"
                :style="{ left: `${showMore.linkLeft}px` }"
                v-show="showMore.showLink" @click="handleShowMore">
                {{showMore.isMoreShowed ? $t('收起') : $t('展开更多')}}<i class="bk-cc-icon icon-cc-arrow-down"></i>
              </div>
            </div>
          </div>
        </div>
        <div class="config-item">
          <div class="item-label">
            {{$t(isDel ? '请勾选要删除的字段：' : '已配置的字段：')}}
          </div>
          <div class="item-content">
            <div class="choose-toolbar" v-if="!isDel">
              <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                <bk-button
                  icon="plus"
                  slot-scope="{ disabled }"
                  :disabled="disabled"
                  @click="handleChooseField">
                  {{$t('选择字段')}}
                </bk-button>
              </cmdb-auth>
              <span class="tips"><i class="bk-cc-icon icon-cc-tips"></i>{{$t('批量设置字段的自动应用功能提示')}}</span>
            </div>
            <div class="config-table" v-show="checkedPropertyIdList.length">
              <property-config-table
                ref="propertyConfigTable"
                :mode="mode"
                :multiple="true"
                :readonly="isDel"
                :deletable="isDel"
                :checked-property-id-list.sync="checkedPropertyIdList"
                :rule-list="initRuleList"
                :id-list="ids"
                @property-value-change="handlePropertyValueChange"
                @selection-change="handlePropertySelectionChange"
                @property-remove="handlePropertyRemove">
              </property-config-table>
            </div>
          </div>
        </div>
      </div>
    </div>
    <template #footer="{ sticky }">
      <div :class="['wrapper-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }" v-if="!isDel">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="nextButtonDisabled || disabled"
            @click="handleNextStep">
            {{$t('下一步')}}
          </bk-button>
        </cmdb-auth>
        <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }" v-else>
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="delButtonDisabled || disabled"
            @click="handleDel">
            {{$t('确定删除按钮')}}
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
      v-bind="leaveConfirmConfig"
      reverse
      :title="$t('是否退出配置')"
      :content="$t('启用步骤未完成，退出将会丢失当前配置')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')">
    </leave-confirm>
  </cmdb-sticky-layout>
</template>
<script>
  /* eslint-disable no-underscore-dangle */
  import { mapGetters, mapState } from 'vuex'
  import leaveConfirm from '@/components/ui/dialog/leave-confirm'
  import topSteps from './top-steps.vue'
  import hostPropertyModal from './host-property-modal'
  import propertyConfigTable from './property-config-table'
  import serviceTemplateTips from './service-template-tips.vue'
  import { MENU_BUSINESS_HOST_APPLY_CONFIRM } from '@/dictionary/menu-symbol'
  import { CONFIG_MODE } from '@/services/service-template/index.js'

  export default {
    name: 'multi-config',
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
      },
      action: {
        type: String,
        default: '' // 'batch-del' | 'batch-edit'
      }
    },
    data() {
      return {
        initRuleList: [],
        checkedPropertyIdList: [],
        showMore: {
          listMaxRow: 2,
          showLink: false,
          isMoreShowed: false,
          linkLeft: 0
        },
        selectedPropertyRow: [],
        propertyModalVisible: false,
        nextButtonDisabled: false,
        delButtonDisabled: true,
        leaveConfirmConfig: {
          id: 'multiConfig',
          active: true
        },
        requestIds: {
          rules: Symbol('rules'),
          del: Symbol('del')
        }
      }
    },
    inject: [
      'getModuleName',
      'getModulePath',
      'getTemplateName'
    ],
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      ...mapState('hostApply', ['ruleDraft']),
      isModuleMode() {
        return this.mode === CONFIG_MODE.MODULE
      },
      isTemplateMode() {
        return this.mode === CONFIG_MODE.TEMPLATE
      },
      isDel() {
        return this.action === 'batch-del'
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
                  bk_module_ids: this.ids
                }
              }
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/getTemplateRules',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  service_template_ids: this.ids
                }
              }
            }
          },
          [this.requestIds.del]: {
            [CONFIG_MODE.MODULE]: {
              action: 'hostApply/deleteRules'
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'hostApply/deleteTemplateRules'
            }
          }
        }
      }
    },
    watch: {
      checkedPropertyIdList() {
        this.$nextTick(() => {
          this.toggleNextButtonDisabled()
        })
      }
    },
    created() {
      this.initData()
      this.leaveConfirmConfig.active = !this.isDel
    },
    mounted() {
      this.setShowMoreLinkStatus()
      window.addEventListener('resize', this.setShowMoreLinkStatus)
    },
    beforeDestroy() {
      window.removeEventListener('resize', this.setShowMoreLinkStatus)
    },
    methods: {
      async initData() {
        try {
          const ruleData = await this.getRules()
          this.initRuleList = ruleData.info
          const attrIds = this.initRuleList.map(item => item.bk_attribute_id)
          const checkedPropertyIdList = [...new Set(attrIds)]
          // eslint-disable-next-line max-len
          this.checkedPropertyIdList = this.hasRuleDraft ? [...new Set([...this.checkedPropertyIdList])] : checkedPropertyIdList
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
      setShowMoreLinkStatus() {
        const { targetList } = this.$refs
        // eslint-disable-next-line prefer-destructuring
        const itemEl = targetList.getElementsByClassName('target-item')[0]
        const itemStyle = getComputedStyle(itemEl)
        // eslint-disable-next-line max-len
        const itemWidth = itemEl.offsetWidth + parseInt(itemStyle.marginLeft, 10) + parseInt(itemStyle.marginRight, 10)
        const targetListWidth = targetList.clientWidth
        const maxCountInRow = Math.floor(targetListWidth / itemWidth)
        const rowCount = Math.ceil(this.ids.length / maxCountInRow)
        this.showMore.showLink = rowCount > this.showMore.listMaxRow
        this.showMore.linkLeft = itemWidth * (maxCountInRow - 1)
      },
      getAvailablePropertyRuleList(propertyRuleList) {
        // 被忽略的不需要
        const availableList = propertyRuleList.filter(property => !property.__extra__.ignore)

        return availableList
      },
      toggleNextButtonDisabled() {
        const { propertyRuleList } = this.$refs.propertyConfigTable
        const availableList = this.getAvailablePropertyRuleList(propertyRuleList)

        const everyTruthy = availableList.every((property) => {
          const validTruthy = property.__extra__.valid !== false

          let valueTruthy = property.__extra__.value
          if (property.bk_property_type === 'bool') {
            valueTruthy = true
          } else if (property.bk_property_type === 'int') {
            valueTruthy = valueTruthy !== null && String(valueTruthy)
          }

          return valueTruthy && validTruthy
        })

        this.nextButtonDisabled = !availableList.length || !everyTruthy
      },
      handleNextStep() {
        const { propertyRuleList, ignoreRuleIds } = this.$refs.propertyConfigTable

        const availablePropertyRuleList = this.getAvailablePropertyRuleList(propertyRuleList)
        const additionalRules = []
        this.ids.forEach((id) => {
          availablePropertyRuleList.forEach((property) => {
            additionalRules.push({
              bk_attribute_id: property.id,
              [this.targetIdKey]: id,
              bk_property_value: property.__extra__.value
            })
          })
        })

        const savePropertyConfig = {
          // 配置对象列表
          [this.targetIdsKey]: this.ids,
          // 附加的规则
          additional_rules: additionalRules,
          // 删除的规则，来源于编辑表格删除
          ignore_rule_ids: ignoreRuleIds
        }
        this.$store.commit('hostApply/setPropertyConfig', savePropertyConfig)
        this.$store.commit('hostApply/setRuleDraft', {
          rules: propertyRuleList
        })

        this.leaveConfirmConfig.active = false
        this.$nextTick(function () {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_APPLY_CONFIRM,
            history: true
          })
        })
      },
      handleDel() {
        this.$bkInfo({
          title: this.$t('确认删除自动应用字段？'),
          subTitle: this.$t('删除后将会移除字段在对应模块中的配置'),
          confirmFn: async () => {
            // eslint-disable-next-line max-len
            const ruleIds = this.selectedPropertyRow.reduce((acc, cur) => acc.concat(cur.__extra__.ruleList.map(item => item.id)), [])
            const requestConfig = this.requestConfigs[this.requestIds.del][this.mode]
            try {
              await this.$store.dispatch(requestConfig.action, {
                bizId: this.bizId,
                params: {
                  data: {
                    host_apply_rule_ids: ruleIds,
                    [this.targetIdsKey]: this.ids
                  }
                }
              })

              this.$success(this.$t('删除成功'))
              this.goBack()
            } catch (e) {
              console.log(e)
            }
          }
        })
      },
      goBack() {
        // 删除离开不用确认
        this.leaveConfirmConfig.active = !this.isDel
        this.$nextTick(() => {
          this.$routerActions.back()
        })
      },
      handleCancel() {
        this.$store.commit('hostApply/clearRuleDraft')
        this.goBack()
      },
      handlePropertySelectionChange(value) {
        this.selectedPropertyRow = value
        this.delButtonDisabled = this.selectedPropertyRow.length <= 0
      },
      handlePropertyValueChange() {
        this.toggleNextButtonDisabled()
      },
      handlePropertyRemove() {
        this.toggleNextButtonDisabled()
      },
      handleChooseField() {
        this.propertyModalVisible = true
      },
      handleShowMore() {
        this.showMore.isMoreShowed = !this.showMore.isMoreShowed
      }
    }
  }
</script>
<style lang="scss" scoped>
  .config-wrapper {
    --labelWidth: 180px;
    max-height: 100%;
    @include scrollbar-y;
    .wrapper-footer {
      display: flex;
      align-items: center;
      height: 52px;
      margin-left: calc(var(--labelWidth) + 12px);
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
        padding: 0 20px;
        margin: 0;
        background-color: #fff;
        border-top: 1px solid $borderColor;
      }
    }
  }

  .multi-config {
    padding-top: 12px;

    .config-tips {
      margin: 0 20px 16px 20px;
    }

    .config-item {
      display: flex;
      margin: 8px 0;

      .item-label {
        flex: none;
        width: var(--labelWidth);
        font-size: 14px;
        font-weight: bold;
        color: #63656e;
        text-align: right;
        margin-right: 12px;
      }
      .item-content {
        flex: auto;
      }

      .choose-toolbar {
        margin-bottom: 18px;;
        .tips {
          font-size: 12px;
          margin-left: 8px;
          .icon-cc-tips {
            margin-right: 8px;
            margin-top: -2px;
            font-size: 14px;
          }
        }
      }
    }
  }

  .target-list {
    position: relative;
    max-height: 72px;
    overflow: hidden;
    transition: all .2s ease-out;

    &.show-more {
      max-height: 100%;
    }
  }
  .target-item {
    position: relative;
    display: inline-block;
    vertical-align: middle;
    height: 26px;
    width: 120px;
    margin: 0 10px 10px 0;
    line-height: 24px;
    padding: 0 20px 0 25px;
    border: 1px solid #c4c6cc;
    border-radius: 13px;
    color: $textColor;
    font-size: 12px;
    cursor: default;
    @include ellipsis;

    &:hover {
      border-color: $primaryColor;
      color: $primaryColor;
      .target-icon {
          background-color: $primaryColor;
      }
    }

    .target-icon {
      position: absolute;
      left: 2px;
      top: 2px;
      width: 20px;
      height: 20px;
      border-radius: 50%;
      line-height: 20px;
      text-align: center;
      color: #fff;
      font-size: 12px;
      background-color: #c4c6cc;
    }

    &.more {
      position: absolute;
      left: 0;
      bottom: 0;
      background: #fafbfd;
      border: 0 none;
      border-radius: unset;
      cursor: pointer;
      color: #3a84ff;
      font-size: 14px;
      text-align: left;
      padding: 0 0 0 .1em;
      line-height: 26px;
      .bk-cc-icon {
        font-size: 22px;
      }

      &.opened {
        position: static;
        .bk-cc-icon {
          transform: rotate(180deg);
        }
      }
    }
  }
</style>
