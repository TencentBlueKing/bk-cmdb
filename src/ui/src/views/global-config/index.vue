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
  <div class="global-config">
    <bk-tab :active.sync="activeTabName">
      <bk-tab-panel name="business-general-config" :label="$t('业务通用')">
        <div class="config-container">
          <BusinessGneralConfig></BusinessGneralConfig>
        </div>
      </bk-tab-panel>
      <bk-tab-panel name="platform-info-config" :label="$t('平台信息')">
        <div class="config-container">
          <PlatformInfoConfig></PlatformInfoConfig>
        </div>
      </bk-tab-panel>
      <bk-tab-panel name="idle-pool-config" :label="$t('业务空闲机池')">
        <div class="config-container">
          <!--加 v-if 是为了每次进入时初始化表单状态 -->
          <IdlePoolConfig v-if="activeTabName === 'idle-pool-config'"></IdlePoolConfig>
        </div>
      </bk-tab-panel>
    </bk-tab>
  </div>
</template>

<script>
  import { ref, watch, onMounted } from 'vue'
  import BusinessGneralConfig from './children/business-general-config.vue'
  import PlatformInfoConfig from './children/platform-info-config.vue'
  import IdlePoolConfig from './children/idle-pool-config.vue'
  import queryStore from '@/router/query'
  import store from '@/store'
  import EventBus from '@/utils/bus'
  import router from '@/router'

  export default {
    name: 'global-config',
    components: {
      BusinessGneralConfig,
      PlatformInfoConfig,
      IdlePoolConfig
    },
    setup() {
      const activeTabName = ref(queryStore.get('tab') || 'business-general-config')

      store.dispatch('globalConfig/fetchDefaultConfig')

      if (!store.state.globalConfig.auth) {
        router.currentRoute.meta.view = 'permission'
      }

      watch(activeTabName, () => {
        queryStore.set('tab', activeTabName.value)
        store.dispatch('globalConfig/fetchConfig')
          .finally(() => {
            EventBus.$emit('globalConfig/fetched')
          })
      })

      onMounted(() => {
        queryStore.set('tab', activeTabName.value)
      })

      return {
        activeTabName
      }
    }
  }
</script>

<style lang="scss" scoped>
  // 覆盖全局设置的内边距，这里因设计要求需要重置
  ::v-deep .bk-tab .bk-tab-header {
    padding: 0;
  }

  .bk-tab{
    background-color: #fff;
    height: auto;
    margin: 20px;
  }

  .config-container {
    display: flex;
    justify-content: center;
    margin-top: 40px;
    padding-bottom: 50px;
  }
</style>
