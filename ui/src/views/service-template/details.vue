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

<script lang="ts">
  import { computed, defineComponent, ref, watch, onMounted } from '@vue/composition-api'
  import router from '@/router/index.js'
  import store from '@/store'
  import RouterQuery from '@/router/query'
  import { getValue } from '@/utils/tools'
  import ServiceTemplateConfig from './children/template-config.vue'
  import ServiceTemplateInstance from './children/template-instance.vue'

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
        handleActiveChange
      }
    }
  })
</script>

<template>
  <div class="details-layout">
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
        <service-template-instance
          :active="active === 'instance'"
          @sync-change="handleSyncStatusChange">
        </service-template-instance>
      </bk-tab-panel>
    </bk-tab>
  </div>
</template>

<style lang="scss" scoped>
.details-layout {
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
}
</style>
