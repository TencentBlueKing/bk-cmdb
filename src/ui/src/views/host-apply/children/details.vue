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
  <div class="host-apply-details">
    <template v-if="id">
      <div class="config-panel">
        <div class="config-head">
          <h2 class="config-title">
            <div class="config-name" v-bk-overflow-tips>{{modeValues.title}}</div>
            <small class="last-edit-time" v-if="localHasRule && ruleLastEditTime">
              ( {{$t('上次编辑时间')}}{{ruleLastEditTime}} )
            </small>
          </h2>
        </div>
        <div class="config-body">
          <template v-if="serviceTemplateApplyEnabled">
            <div class="empty">
              <div class="desc">
                <i class="bk-cc-icon icon-cc-tips"></i>
                <span>{{$t('当前模块所属的服务模板已配置了自动应用规则')}}</span>
              </div>
              <div class="action">
                <bk-button
                  outline
                  theme="primary"
                  @click="handleGoTemplate">
                  {{$t('跳转查看')}}
                </bk-button>
              </div>
            </div>
          </template>

          <template v-else-if="applyEnabled">
            <div class="view-field">
              <div class="view-bd">
                <div class="field-list">
                  <div class="field-list-table">
                    <property-config-table
                      ref="propertyConfigTable"
                      :mode="mode"
                      :readonly="true"
                      :show-del-column="true"
                      :checked-property-id-list.sync="checkedPropertyIdList"
                      :rule-list="ruleList"
                      @property-rule-delete="handleDeletePropertyRule">
                    </property-config-table>
                  </div>
                </div>
              </div>
              <div class="view-ft">
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                  <template #default="{ disabled }">
                    <bk-button
                      theme="primary"
                      :disabled="disabled"
                      @click="$emit('edit')">
                      {{$t('编辑')}}
                    </bk-button>
                  </template>
                </cmdb-auth>
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                  <template #default="{ disabled }">
                    <bk-button
                      :disabled="!hasConflict || disabled"
                      @click="$emit('view-conflict')">
                      <span v-bk-tooltips="{ content: $t('无失效需处理') }" v-if="!hasConflict">
                        {{$t('失效主机')}}<em class="conflict-num">{{conflictNum}}</em>
                      </span>
                      <span v-else>
                        {{$t('失效主机')}}<em class="conflict-num">{{conflictNum}}</em>
                      </span>
                    </bk-button>
                  </template>
                </cmdb-auth>
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                  <template #default="{ disabled }">
                    <bk-button
                      :disabled="disabled"
                      @click="$emit('close')">
                      {{$t('关闭自动应用')}}
                    </bk-button>
                  </template>
                </cmdb-auth>
              </div>
            </div>
          </template>

          <template v-else>
            <div class="empty" v-if="!localHasRule">
              <div class="desc">
                <i class="bk-cc-icon icon-cc-tips"></i>
                <span v-if="mode === CONFIG_MODE.MODULE">{{$t('当前模块未启用自动应用策略')}}</span>
                <span v-else-if="mode === CONFIG_MODE.TEMPLATE">{{$t('当前模板未启用自动应用策略')}}</span>
              </div>
              <div class="action">
                <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                  <template #default="{ disabled }">
                    <bk-button
                      outline
                      theme="primary"
                      :disabled="disabled"
                      @click="$emit('edit')">
                      {{$t('立即启用')}}
                    </bk-button>
                  </template>
                </cmdb-auth>
              </div>
            </div>
            <div class="view-field" v-else>
              <div class="view-bd">
                <div class="field-list">
                  <div class="field-list-table disabled">
                    <property-config-table
                      ref="propertyConfigTable"
                      :mode="mode"
                      :readonly="true"
                      :checked-property-id-list.sync="checkedPropertyIdList"
                      :rule-list="ruleList">
                    </property-config-table>
                  </div>
                  <div class="closed-mask">
                    <div class="empty">
                      <div class="desc">
                        <i class="bk-cc-icon icon-cc-tips"></i>
                        <span v-if="mode === CONFIG_MODE.MODULE">{{$t('该模块已关闭属性自动应用')}}</span>
                        <span v-else-if="mode === CONFIG_MODE.TEMPLATE">{{$t('该模板已关闭属性自动应用')}}</span>
                      </div>
                      <div class="action">
                        <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
                          <template #default="{ disabled }">
                            <bk-button
                              outline
                              theme="primary"
                              :disabled="disabled"
                              @click="$emit('edit')">
                              {{$t('重新启用')}}
                            </bk-button>
                          </template>
                        </cmdb-auth>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </template>
    <div class="empty" v-else>
      <div class="desc">
        <i class="bk-cc-icon icon-cc-tips"></i>
        <span v-if="mode === CONFIG_MODE.MODULE">{{$t('主机属性自动应用暂无业务模块')}}</span>
        <span v-else-if="mode === CONFIG_MODE.TEMPLATE">{{$t('主机属性自动应用暂无服务模板')}}</span>
      </div>
      <div class="action">
        <bk-button
          outline
          theme="primary"
          @click="$routerActions.redirect({ name: modeValues.createLink })">
          {{$t('跳转创建')}}
        </bk-button>
      </div>
    </div>
  </div>
</template>

<script>
  import { computed, defineComponent, ref, toRefs } from '@vue/composition-api'
  import {
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_SERVICE_TEMPLATE
  } from '@/dictionary/menu-symbol'
  import router from '@/router/index.js'
  import { formatTime } from '@/utils/tools.js'
  import { CONFIG_MODE } from '@/services/service-template/index.js'
  import propertyConfigTable from './property-config-table'

  export default defineComponent({
    components: {
      propertyConfigTable
    },
    props: {
      id: Number,
      bizId: Number,
      currentNode: Object,
      ruleList: Array,
      hasRule: Boolean,
      checkedPropertyIdList: Array,
      conflictNum: Number
    },
    setup(props, { emit }) {
      const {
        id,
        hasRule,
        ruleList,
        currentNode,
        conflictNum
      } = toRefs(props)

      const mode = computed(() => router.app.$route.params?.mode)

      const modeValues = computed(() => {
        const values = {
          [CONFIG_MODE.MODULE]: {
            targetIdKey: 'bk_module_id',
            createLink: MENU_BUSINESS_HOST_AND_SERVICE,
            title: currentNode.value.bk_inst_name
          },
          [CONFIG_MODE.TEMPLATE]: {
            targetIdKey: 'service_template_id',
            createLink: MENU_BUSINESS_SERVICE_TEMPLATE,
            title: currentNode.value.name
          }
        }
        return values[mode.value] ?? {}
      })

      const propertyConfigTable = ref(null)

      const ruleLastEditTime = computed(() => {
        const lastTimeList = ruleList.value
          ?.filter(rule => rule[modeValues.value.targetIdKey] === id.value)
          ?.map(rule => new Date(rule.last_time).getTime())
        if (lastTimeList?.length) {
          const latestTime = Math.max(...lastTimeList)
          return formatTime(latestTime, 'YYYY-MM-DD HH:mm:ss')
        }
        return ''
      })

      const applyEnabled = computed(() => currentNode.value?.host_apply_enabled)

      const serviceTemplateApplyEnabled = computed(() => currentNode.value?.service_template_host_apply_enabled)

      const localHasRule = computed(() => hasRule.value && ruleList.value?.length > 0)

      const hasConflict = computed(() => conflictNum.value > 0)

      const handleDeletePropertyRule = (property) => {
        emit('delete-rule', property)
      }

      const handleGoTemplate = () => {
        router.push({
          params: {
            mode: CONFIG_MODE.TEMPLATE
          },
          query: {
            id: currentNode.value.service_template_id
          }
        })
      }

      const reset = () => {
        propertyConfigTable.value?.reset()
      }

      return {
        CONFIG_MODE,
        mode,
        modeValues,
        propertyConfigTable,
        ruleLastEditTime,
        applyEnabled,
        serviceTemplateApplyEnabled,
        localHasRule,
        hasConflict,
        reset,
        handleDeletePropertyRule,
        handleGoTemplate,
      }
    }
  })
</script>

<style lang="scss" scoped>
  .host-apply-details {
    height: 100%;

    .empty {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 80%;

      .desc {
        font-size: 14px;
        color: #63656e;

        .icon-cc-tips {
          margin-top: -2px;
        }
      }
      .action {
        margin-top: 18px;
      }
    }
  }

  .config-panel {
    display: flex;
    flex-direction: column;
    height: 100%;

    .config-head,
    .config-foot {
      flex: none;
    }
    .config-body {
      flex: auto;
    }

    .config-title {
      display: flex;
      align-items: center;
      height: 32px;
      font-size: 14px;
      color: #313238;
      font-weight: 700;
      margin-top: 8px;

      .config-name {
        @include ellipsis;
      }

      .last-edit-time {
        flex: none;
        font-size: 12px;
        font-weight: 400;
        color: #979ba5;
        margin-left: .2em;
      }
    }

    .view-field {
      .field-list {
        position: relative;

        .field-list-table {
          &.disabled {
            opacity: 0.2;
          }
        }
        .closed-mask {
          position: absolute;
          width: 100%;
          height: 100%;
          min-height: 210px;
          left: 0;
          top: 0;
        }
      }
      .view-bd,
      .view-ft {
        margin: 20px 0;
        .bk-button {
          margin-right: 4px;
          min-width: 86px;
        }
      }
    }

    .conflict-num {
      font-size: 12px;
      color: #fff;
      background: #c4c6cc;
      border-radius: 8px;
      font-style: normal;
      padding: 0px 4px;
      font-family: arial;
      margin-left: 4px;
    }
  }

  .close-apply-confirm-modal {
    .content {
      font-size: 14px;
    }
    .tips {
      margin: 12px 0;
    }
  }
</style>
<style lang="scss">
  .close-apply-confirm-modal {
    .bk-dialog-sub-header {
      padding-left: 32px !important;
      padding-right: 32px !important;
    }
  }
</style>
