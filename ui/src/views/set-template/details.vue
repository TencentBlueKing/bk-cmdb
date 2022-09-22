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
  import SetTemplateConfig from './children/template-config.vue'
  import SetTemplateInstance from './children/template-instance.vue'
  import { getTemplateSyncStatus } from './children/use-template-data.js'

  export default defineComponent({
    components: {
      SetTemplateConfig,
      SetTemplateInstance
    },
    setup() {
      const $tab = ref(null)

      const query = computed(() => RouterQuery.getAll())

      const active = ref(query.value.tab)

      const bizId = computed(() => store.getters['objectBiz/bizId'])

      const templateId = computed(() => parseInt(router.app.$route.params.templateId, 10))

      const checkSyncStatus = async () => {
        try {
          const needSync = await getTemplateSyncStatus(bizId.value, templateId.value)
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
      <bk-tab-panel :label="$t('集群模板配置')" name="config">
        <set-template-config
          @sync-change="handleSyncStatusChange"
          @active-change="handleActiveChange">
        </set-template-config>
      </bk-tab-panel>
      <bk-tab-panel :label="$t('集群模板实例')" render-directive="if" name="instance">
        <set-template-instance
          :template-id="templateId"
          @sync-change="handleSyncStatusChange">
        </set-template-instance>
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
        height: calc(100% - 50px);
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
