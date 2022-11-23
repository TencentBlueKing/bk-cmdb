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
  <div class="confirm-wrapper">
    <top-steps :current="2" />
    <div class="host-apply-confirm">
      <div class="update-options">
        <div class="option-label">
          <i18n path="同时更新未应用主机">
            <template #invalid>
              <span class="has-tips" v-bk-tooltips="$t('属性当前值与目标值不一致的主机')">{{$t('-未应用主机')}}</span>
            </template>
          </i18n>
        </div>
        <bk-radio-group class="option-content" v-model="updateOption.changed">
          <bk-radio :value="true">
            {{$t('是将把未应用主机更新为当前的配置')}}
          </bk-radio>
          <bk-radio :value="false">
            {{$t('否将保留未应用主机的配置')}}
          </bk-radio>
        </bk-radio-group>
      </div>
      <div class="caption" v-show="updateOption.changed">
        <div class="title">{{$t('请确认以下主机应用信息')}}</div>
        <div class="stat">
          <span class="conflict-item">
            <span v-bk-tooltips="{ content: $t('当一台主机同属多个模块时，且模块配置了不同的属性'), width: 185 }">
              <i class="bk-cc-icon icon-cc-tips"></i>
            </span>
            <i18n path="冲突主机N台">
              <template #num><em class="conflict-num">{{conflictNum}}</em></template>
            </i18n>
          </span>
          <i18n path="主机总数N台">
            <template #num><em class="check-num">{{totalNum}}</em></template>
          </i18n>
        </div>
      </div>
      <!-- max-height 用于控制表格内滚动 -->
      <property-confirm-table v-show="updateOption.changed"
        v-bkloading="{ isLoading: $loading(requestIds.applyPreview) }"
        ref="propertyConfirmTable"
        :max-height="$APP.height - 220 - 120"
        :list="table.list"
        :total="table.total">
      </property-confirm-table>
      <div class="bottom-actionbar">
        <div class="actionbar-inner">
          <bk-button theme="default" @click="handlePrevStep">{{$t('上一步')}}</bk-button>
          <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
            <bk-button
              theme="primary"
              slot-scope="{ disabled }"
              :disabled="applyButtonDisabled || disabled"
              @click="handleApply">
              {{$t('保存并应用')}}
            </bk-button>
          </cmdb-auth>
          <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
      </div>
    </div>

    <leave-confirm
      v-bind="leaveConfirmConfig"
      reverse
      :title="$t('是否退出配置')"
      :content="$t('启用步骤未完成，退出将撤销当前操作')"
      :ok-text="$t('退出')"
      :cancel-text="$t('取消')">
    </leave-confirm>
  </div>
</template>

<script>
  import { mapGetters, mapState, mapActions } from 'vuex'
  import leaveConfirm from '@/components/ui/dialog/leave-confirm'
  import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
  import topSteps from './children/top-steps.vue'
  import {
    MENU_BUSINESS_HOST_APPLY,
    MENU_BUSINESS_HOST_APPLY_EDIT,
    MENU_BUSINESS_HOST_APPLY_RUN
  } from '@/dictionary/menu-symbol'
  import { CONFIG_MODE } from '@/service/service-template/index.js'

  export default {
    components: {
      leaveConfirm,
      topSteps,
      propertyConfirmTable
    },
    data() {
      return {
        updateOption: {
          changed: true
        },
        table: {
          list: [],
          total: 0
        },
        conflictNum: 0,
        totalNum: 0,
        requestIds: {
          applyPreview: Symbol()
        },
        leaveConfirmConfig: {
          id: 'propertyConfirm',
          active: true
        },
        applyButtonDisabled: false
      }
    },
    computed: {
      ...mapState('hostApply', ['propertyConfig']),
      ...mapGetters('objectBiz', ['bizId']),
      mode() {
        return this.$route.params.mode
      },
      isModuleMode() {
        return this.mode === CONFIG_MODE.MODULE
      },
      isTemplateMode() {
        return this.mode === CONFIG_MODE.TEMPLATE
      },
      isBatch() {
        return this.$route.query.batch === 1
      },
      targetIdsKey() {
        const targetIdsKeys = {
          [CONFIG_MODE.MODULE]: 'bk_module_ids',
          [CONFIG_MODE.TEMPLATE]: 'service_template_ids'
        }
        return targetIdsKeys[this.mode]
      },
      targetIds() {
        return this.propertyConfig[this.targetIdsKey]
      },
      requestConfigs() {
        return {
          [this.requestIds.applyPreview]: {
            [CONFIG_MODE.MODULE]: {
              action: 'getApplyPreview',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  ...this.propertyConfig
                }
              }
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'getTemplateApplyPreview',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  ...this.propertyConfig
                }
              }
            }
          }
        }
      }
    },
    beforeRouteLeave(to, from, next) {
      if (to.name !== MENU_BUSINESS_HOST_APPLY_EDIT) {
        this.$store.commit('hostApply/clearRuleDraft')
      }
      next()
    },
    created() {
      // 无配置数据时强制跳转至入口页
      if (!Object.keys(this.propertyConfig).length) {
        this.leaveConfirmConfig.active = false
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY
        })
      } else {
        this.setBreadcrumbs()
        this.initData()
      }
    },
    methods: {
      ...mapActions('hostApply', [
        'getApplyPreview',
        'getTemplateApplyPreview'
      ]),
      async initData() {
        try {
          const requestConfig = this.requestConfigs[this.requestIds.applyPreview][this.mode]
          const previewData = await this[requestConfig.action]({
            ...requestConfig.payload,
            config: {
              requestId: this.requestIds.applyPreview
            }
          })

          // 前端分页
          this.table.list = previewData.plans || []
          this.table.total = this.table.list.length

          // 主机总数
          this.totalNum = previewData.count

          this.conflictNum = previewData.unresolved_conflict_count
        } catch (e) {
          this.applyButtonDisabled = true
          console.error(e)
        }
      },
      setBreadcrumbs() {
        const title = this.isBatch ? '批量应用属性' : '应用属性'
        this.$store.commit('setTitle', this.$t(title))
      },
      goBack() {
        const query = {}
        if (!this.isBatch) {
          // eslint-disable-next-line prefer-destructuring
          query.id = this.propertyConfig[this.targetIdsKey][0]
        }
        this.$store.commit('hostApply/clearRuleDraft')
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY,
          params: {
            mode: this.mode
          },
          query
        })
      },
      async handleApply() {
        // 合入是否更新失效主机选项值
        const propertyConfig = { ...this.propertyConfig, ...{ changed: this.updateOption.changed } }
        // 更新属性配置，用于后续应用执行
        this.$store.commit('hostApply/setPropertyConfig', propertyConfig)

        this.leaveConfirmConfig.active = false
        this.$nextTick(() => {
          this.$routerActions.redirect({
            name: MENU_BUSINESS_HOST_APPLY_RUN
          })
        })
      },
      handleCancel() {
        this.goBack()
      },
      handlePrevStep() {
        this.leaveConfirmConfig.active = false
        this.$nextTick(() => {
          this.$routerActions.back()
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .host-apply-confirm {
    padding: 15px 20px 0;

    .update-options {
      display: flex;
      align-items: center;
      margin: 24px 0;

      .option-label {
        position: relative;
        font-size: 14px;
        white-space: nowrap;
        margin-right: 18px;

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
    }

    .caption {
      display: flex;
      margin-bottom: 14px;
      justify-content: space-between;
      align-items: center;

      .title {
        color: #63656e;
        font-size: 14px;
      }

      .stat {
        color: #313238;
        font-size: 12px;
        margin-right: 8px;

        .conflict-item {
          margin-right: 12px;
        }

        .conflict-num,
        .check-num {
          font-style: normal;
          font-weight: bold;
          margin: 0 .2em;
        }
        .conflict-num {
          color: #ff5656;
        }
        .check-num {
          color: #2dcb56;
        }
      }
    }
  }

  .bottom-actionbar {
    width: 100%;
    height: 50px;
    z-index: 100;

    .actionbar-inner {
      padding: 20px 0 0 0;
      .bk-button {
        margin-right: 4px;
        min-width: 86px;
      }
    }
  }
</style>
