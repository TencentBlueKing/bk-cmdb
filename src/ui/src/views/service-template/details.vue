<script lang="ts">
  import { computed, defineComponent, ref, watch, onMounted } from '@vue/composition-api'
  import router from '@/router/index.js'
  import store from '@/store'
  import RouterQuery from '@/router/query'
  import routerActions from '@/router/actions'
  import { getValue } from '@/utils/tools'
  import ServiceTemplateConfig from './children/template-config.vue'
  import ServiceTemplateInstance from './children/template-instance.vue'
  import { MENU_BUSINESS_SERVICE_TEMPLATE_EDIT } from '@/dictionary/menu-symbol'

  export default defineComponent({
    components: {
      ServiceTemplateConfig,
      ServiceTemplateInstance
    },
    setup() {
      const $tab = ref(null)

      const query = computed(() => RouterQuery.getAll())

      const active = ref(query.value.tab)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => parseInt(router.app.$route.params.templateId, 10))

      const checkSyncStatus = async () => {
        try {
          const { service_templates: syncStatusList = [] } = await store.dispatch('serviceTemplate/getServiceTemplateSyncStatus', {
            bizId: bizId.value,
            params: {
              is_partial: true,
              service_template_ids: [templateId.value]
            },
            config: {
              cancelPrevious: true
            }
          })
          const needSync = getValue(syncStatusList, '0.need_sync')
          const tabHeader = $tab.value.$el.querySelector('.bk-tab-label-item.is-last')
          if (needSync) {
            tabHeader.classList.add('has-tips')
          } else {
            tabHeader.classList.remove('has-tips')
          }
        } catch (error) {
          console.error(error)
        }
      }

      const handleSyncStatusChange = () => {
        checkSyncStatus()
      }

      const handleEdit = () => {
        routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE_EDIT,
          params: {
            templateId: templateId.value
          },
          history: true
        })
      }

      const handleActiveChange = (value) => {
        active.value = value
      }

      watch(active, (active) => {
        RouterQuery.set({
          tab: active
        })
      }, { immediate: true })

      onMounted(() => {
        checkSyncStatus()
      })

      return {
        bizId,
        templateId,
        $tab,
        active,
        handleSyncStatusChange,
        handleEdit,
        handleActiveChange
      }
    }
  })
</script>

<template>
  <cmdb-sticky-layout class="details-layout">
    <div class="layout-main">
      <bk-tab class="template-tab"
        type="unborder-card"
        ref="$tab"
        :active.sync="active">
        <bk-tab-panel :label="$t('服务模板配置')" name="config">
          <service-template-config
            @sync-change="handleSyncStatusChange"
            @active-change="handleActiveChange">
          </service-template-config>
        </bk-tab-panel>
        <bk-tab-panel :label="$t('服务模板实例')" name="instance">
          <service-template-instance :active="active === 'instance'"></service-template-instance>
        </bk-tab-panel>
      </bk-tab>
    </div>
    <template #footer="{ sticky }">
      <div :class="['layout-footer', { 'is-sticky': sticky }]">
        <cmdb-auth :auth="{ type: $OPERATION.U_SERVICE_TEMPLATE, relation: [bizId, templateId] }">
          <bk-button
            theme="primary"
            slot-scope="{ disabled }"
            :disabled="disabled"
            @click="handleEdit">
            {{$t('编辑')}}
          </bk-button>
        </cmdb-auth>
      </div>
    </template>
  </cmdb-sticky-layout>
</template>

<style lang="scss" scoped>
.details-layout {
  .layout-footer {
    display: flex;
    align-items: center;
    height: 52px;
    padding: 0 20px;
    margin-top: 8px;
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

.template-tab {
  ::v-deep {
    .bk-tab-header {
      padding: 0;
      margin: 0;
    }
    .bk-tab-section {
      padding: 0;
    }
  }
  ::v-deep .bk-tab-label-item.has-tips:before {
    content: "";
    position: absolute;
    top: 18px;
    right: 12px;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: $dangerColor;
  }
}
</style>
