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
  <div class="template-layout">
    <bk-tab class="template-tab"
      type="unborder-card"
      ref="infoTab"
      :class="{
        'no-header': !isUpdate
      }"
      :show-header="isUpdate"
      :active.sync="active">
      <bk-tab-panel :label="$t('服务模板配置')" name="config">
        <service-template-config @sync-change="handleSyncStatusChange"></service-template-config>
      </bk-tab-panel>
      <bk-tab-panel :label="$t('服务模板实例')" name="instance" v-if="isUpdate">
        <service-template-instance :active="active === 'instance'"></service-template-instance>
      </bk-tab-panel>
    </bk-tab>
  </div>
</template>

<script>
  import { mapActions, mapGetters } from 'vuex'
  import ServiceTemplateConfig from './children/operational'
  import ServiceTemplateInstance from './children/template-instance'
  import Bus from '@/utils/bus'
  import RouterQuery from '@/router/query'
  export default {
    components: {
      ServiceTemplateConfig,
      ServiceTemplateInstance
    },
    data() {
      return {
        active: RouterQuery.get('tab', 'config')
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      templateId() {
        return this.$route.params.templateId
      },
      isUpdate() {
        return this.templateId !== undefined
      }
    },
    watch: {
      active: {
        immediate: true,
        handler(active) {
          RouterQuery.set({
            tab: active
          })
        }
      }
    },
    created() {
      this.checkSyncStatus()
      Bus.$on('active-change', (active) => {
        this.active = active
      })
    },
    beforeDestroy() {
      Bus.$off('active-change')
    },
    methods: {
      ...mapActions('serviceTemplate', [
        'getServiceTemplateSyncStatus'
      ]),
      async checkSyncStatus() {
        try {
          const { service_templates: syncStatusList = [] } = await this.getServiceTemplateSyncStatus({
            bizId: this.bizId,
            params: {
              is_partial: true,
              service_template_ids: [Number(this.templateId)]
            },
            config: {
              cancelPrevious: true
            }
          })
          const needSync = this.$tools.getValue(syncStatusList, '0.need_sync')
          const tabHeader = this.$refs.infoTab.$el.querySelector('.bk-tab-label-item.is-last')
          if (needSync) {
            tabHeader.classList.add('has-tips')
          } else {
            tabHeader.classList.remove('has-tips')
          }
        } catch (error) {
          console.error(error)
        }
      },
      handleSyncStatusChange() {
        this.checkSyncStatus()
      }
    }
  }
</script>

<style lang="scss" scoped>
    .template-tab {
        /deep/ {
            .bk-tab-header {
                padding: 0;
                margin: 0 20px;
            }
            .bk-tab-section {
                padding: 0;
            }
        }
        /deep/ .bk-tab-label-item.has-tips:before {
            content: "";
            position: absolute;
            top: 18px;
            right: 12px;
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background-color: $dangerColor;
        }
        &.no-header {
            /deep/ .bk-tab-section {
                height: 100%;
            }
        }
    }
</style>
