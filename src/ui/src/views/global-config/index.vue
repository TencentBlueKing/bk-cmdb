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
    <bk-tab :active.sync="activeTabName" :before-toggle="hanldeBeforeToggle">
      <bk-tab-panel name="business-general-config" :label="$t('业务通用')">
        <div class="config-container">
          <BusinessGneralConfig @has-change="handleHasValChange"></BusinessGneralConfig>
        </div>
      </bk-tab-panel>
      <bk-tab-panel name="idle-pool-config" :label="$t('业务空闲机池')">
        <div class="config-container">
          <!--加 v-if 是为了每次进入时初始化表单状态 -->
          <IdlePoolConfig v-if="activeTabName === 'idle-pool-config'" @has-change="handleHasValChange"></IdlePoolConfig>
        </div>
      </bk-tab-panel>
      <bk-tab-panel name="id-generate" :label="$t('ID生成器')">
        <div class="config-container">
          <IDGenerate v-if="activeTabName === 'id-generate'" @has-change="handleHasValChange"></IDGenerate>
        </div>
      </bk-tab-panel>
    </bk-tab>
  </div>
</template>

<script>
  import { ref, watch, onMounted, h } from 'vue'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox/index.js'
  import BusinessGneralConfig from './children/business-general-config.vue'
  import IdlePoolConfig from './children/idle-pool-config.vue'
  import IDGenerate from './children/id-generate.vue'
  import queryStore from '@/router/query'
  import store from '@/store'
  import EventBus from '@/utils/bus'
  import router from '@/router'

  export default {
    name: 'global-config',
    components: {
      BusinessGneralConfig,
      IdlePoolConfig,
      IDGenerate
    },
    setup() {
      const subHeader = h('div', {
        class: 'leave-confirm-content'
      }, t('离开将会导致未保存信息丢失'))
      const activeTabName = ref(queryStore.get('tab') || 'business-general-config')
      const hasChange = ref(false)

      store.dispatch('globalConfig/fetchConfig')

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

      const handleHasValChange = (val) => {
        hasChange.value = val
      }
      const hanldeBeforeToggle = (name) => {
        if (!hasChange.value) {
          activeTabName.value = name
          return
        }
        const confirmFn = () => {
          activeTabName.value = name
          hasChange.value = false
        }
        beforeLeave(confirmFn)
      }
      const beforeLeave = (comfirm, cancel) => {
        $bkInfo({
          title: t('确认离开当前页？'),
          subHeader,
          okText: t('离开'),
          cancelText: t('取消'),
          closeIcon: false,
          confirmFn: () => {
            comfirm?.()
          },
          cancelFn: () => {
            cancel?.()
          }
        })
      }

      return {
        activeTabName,
        hasChange,
        beforeLeave,
        handleHasValChange,
        hanldeBeforeToggle
      }
    },
    beforeRouteLeave(to, from, next) {
      if (!this.hasChange) {
        next()
        return
      }
      this.beforeLeave(() => next(), () => next(false))
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
  .leave-confirm-content {
    text-align: center;
    font-size: 14px;
  }
</style>
