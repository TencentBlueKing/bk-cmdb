<template>
  <div class="conflict-list">
    <property-confirm-table
      ref="propertyConfirmTable"
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
            :disabled="applyButtonDisabled || disabled"
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
  import {
    MENU_BUSINESS_HOST_APPLY,
    MENU_BUSINESS_HOST_APPLY_RUN
  } from '@/dictionary/menu-symbol'

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
        applyRequest: null,
        applyButtonDisabled: false
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      moduleId() {
        return Number(this.$route.query.mid)
      }
    },
    created() {
      this.setBreadcrumbs()
      this.initData()
    },
    methods: {
      ...mapActions('hostApply', [
        'getApplyPreview'
      ]),
      async initData() {
        try {
          const previewData = await this.getApplyPreview({
            params: {
              bk_biz_id: this.bizId,
              bk_module_ids: [this.moduleId]
            },
            config: {
              requestId: 'getHostApplyPreview'
            }
          })
          this.table.list = previewData.plans ?? []
          this.table.total = previewData.count
        } catch (e) {
          this.applyButtonDisabled = true
          console.error(e)
        }
      },
      setBreadcrumbs() {
        this.$store.commit('setTitle', this.$t('策略失效主机'))
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
          // bk_module_ids: [this.moduleId],
          additional_rules: updateFields.map(item => ({
            bk_attribute_id: item.bk_attribute_id,
            bk_module_id: this.moduleId,
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
        this.$routerActions.redirect({
          name: MENU_BUSINESS_HOST_APPLY,
          query: {
            module: this.moduleId
          }
        })
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
      .bk-button {
        min-width: 86px;
      }
    }
  }
</style>
