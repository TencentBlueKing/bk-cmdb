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
  <div class="conflict-list">
    <property-confirm-table
      v-bkloading="{ isLoading: $loading(requestIds.applyPreview) }"
      :list="table.list"
      :max-height="$APP.height - 220"
      :total="table.total">
    </property-confirm-table>
    <div class="bottom-actionbar">
      <div class="actionbar-inner">
        <cmdb-auth :auth="{ type: $OPERATION.U_HOST_APPLY, relation: [bizId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="$loading(requestIds.applyPreview) || disabled"
            @click="handleApply">
            {{$t('应用')}}
          </bk-button>
        </cmdb-auth>
        <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
      </div>
    </div>
  </div>
</template>

<script>
  import { mapGetters, mapActions } from 'vuex'
  import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
  import { MENU_BUSINESS_HOST_APPLY_RUN } from '@/dictionary/menu-symbol'
  import { CONFIG_MODE } from '@/service/service-template/index.js'

  export default {
    components: {
      propertyConfirmTable
    },
    data() {
      return {
        table: {
          list: [],
          total: 0
        },
        requestIds: {
          applyPreview: Symbol()
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      mode() {
        return this.$route.params.mode
      },
      targetId() {
        return Number(this.$route.query.id)
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
          [this.requestIds.applyPreview]: {
            [CONFIG_MODE.MODULE]: {
              action: 'getApplyPreview',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  [this.targetIdsKey]: [this.targetId]
                }
              }
            },
            [CONFIG_MODE.TEMPLATE]: {
              action: 'getTemplateApplyPreview',
              payload: {
                params: {
                  bk_biz_id: this.bizId,
                  [this.targetIdsKey]: [this.targetId]
                }
              }
            }
          }
        }
      }
    },
    created() {
      this.setBreadcrumbs()
      this.initData()
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

          this.table.list = previewData.plans ?? []
          this.table.total = previewData.unresolved_conflict_count
        } catch (e) {
          console.error(e)
        }
      },
      setBreadcrumbs() {
        this.$store.commit('setTitle', this.$t('未应用主机'))
      },
      handleApply() {
        this.$bkInfo({
          title: this.$t('确认应用'),
          subTitle: this.$t('自动应用的主机属性统一修改为目标值确认'),
          confirmFn: () => {
            this.gotoApply()
          }
        })
      },
      async gotoApply() {
        // 失效主机只有单模块场景，需要更新的字段每一条都是一致的
        const updateFields = this.table.list[0].update_fields

        // 合入是否更新失效主机选项值
        const propertyConfig = {
          changed: true,
          [this.targetIdsKey]: [this.targetId],
          additional_rules: updateFields.map(item => ({
            bk_attribute_id: item.bk_attribute_id,
            [this.targetIdKey]: this.targetId,
            bk_property_value: item.bk_property_value
          })),
          bk_host_ids: this.table.list.map(item => item.bk_host_id)
        }

        // 更新属性配置，用于后续应用执行
        this.$store.commit('hostApply/setPropertyConfig', propertyConfig)

        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY_RUN,
          query: {
            from: 'conflict-list'
          }
        })
      },
      handleCancel() {
        this.$routerActions.back()
      }
    }
  }
</script>

<style lang="scss" scoped>
  .conflict-list {
    padding: 15px 20px 0;
  }

  .bottom-actionbar {
    width: 100%;
    height: 50px;
    z-index: 100;

    .actionbar-inner {
      padding: 20px 0 0 0;
      display: flex;
      align-items: center;

      .bk-button {
        min-width: 86px;
      }
      & + .bk-button {
        margin-left: 8px;
      }
      .auth-box + .bk-button {
        margin-left: 8px;
      }
    }
  }
</style>
