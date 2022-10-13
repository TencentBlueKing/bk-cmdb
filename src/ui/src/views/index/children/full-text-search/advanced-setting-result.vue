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
  <div class="setting-tags" v-show="customized">
    <div class="tag-item target">检索对象：{{targetScopes}}</div>
    <div class="tag-item" v-for="(item, index) in targetModels" :key="index" v-bk-overflow-tips>
      {{item.targetName}}：{{item.models}}
    </div>
    <bk-button class="reset" text size="small"
      @click="handleClear">
      {{$t('清空所有')}}
    </bk-button>
  </div>
</template>

<script>
  import { computed, defineComponent } from 'vue'
  import store from '@/store'
  import { t } from '@/i18n'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import { targetMap, finalSetting, handleReset } from './use-advanced-setting.js'
  import { pickQuery } from './use-route.js'

  export default defineComponent({
    setup() {
      const route = computed(() => RouterQuery.route)

      const getModelById = store.getters['objectModelClassify/getModelById']
      const getModelName = id => getModelById(id)?.bk_obj_name ?? '--'

      const targetModels = computed(() => {
        const targetModels = []
        finalSetting.value.targets.forEach((target) => {
          const modelIds = finalSetting.value[`${target}s`]
          targetModels.push({
            targetName: targetMap[target],
            models: modelIds.length ? modelIds.map(id => getModelName(id)).join(' | ') : t('全部')
          })
        })
        return targetModels
      })

      const customized = computed(() => {
        const changedModels = []
        finalSetting.value.targets.forEach((target) => {
          const modelIds = finalSetting.value[`${target}s`]
          changedModels.push(modelIds.length > 0)
        })
        return changedModels.some(changed => changed)
      })

      const targetScopes = computed(() => finalSetting.value.targets.map(target => targetMap[target]).join(' | '))

      const handleClear = () => {
        handleReset()
        const query = pickQuery(route.value.query, ['tab', 'keyword'])
        routerActions.redirect({
          query: {
            ...query,
            t: Date.now()
          }
        })
      }

      return {
        customized,
        targetMap,
        targetScopes,
        targetModels,
        handleClear
      }
    }
  })
</script>

<style lang="scss" scoped>
  .setting-tags {
    display: flex;
    font-size: 12px;
    align-items: center;

    .tag-item {
      background: #dcdee5;
      border-radius: 2px;
      padding: 3px 8px;
      max-width: 500px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;

      & + .tag-item {
        margin-left: 6px;
      }

      &.target {
        flex: none;
      }

      .sub-item {
        & + .sub-item {
          &::before {
            content: "|";
            padding: 0 3px;
          }
        }
      }
    }

    .reset {
      flex: none;
    }
  }
</style>
