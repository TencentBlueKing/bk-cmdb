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

<script setup>
  import { computed, ref, watchEffect } from 'vue'
  import { t } from '@/i18n'
  import routerActions from '@/router/actions'
  import {
    MENU_MODEL_FIELD_TEMPLATE,
  } from '@/dictionary/menu-symbol'
  import ResultStatus from './result-status.vue'
  import fieldTemplateService from '@/service/field-template'

  const props = defineProps({
    scene: {
      type: String,
      default: 'bind'
    },
    templateId: {
      type: Number
    },
    modelIds: {
      type: Array,
      default: () => []
    }
  })

  const status = ref('')
  const title = ref('')
  const summary = ref('')
  const counts = ref({
    success: 0,
    fail: 0
  })
  let timer = null

  const requestIds = {
    sync: Symbol()
  }

  const isSyncResultStatus = computed(() => ['success', 'error', 'abnormal'].includes(status.value))
  const isEditEmptyModel = computed(() => props.scene === 'edit' && !props.modelIds.length)

  const getTaskSyncStatus = async (taskIds) => {
    try {
      status.value = 'loading'
      title.value = t('同步中...')

      const statusList = await fieldTemplateService.getTaskSyncStatus({
        task_ids: taskIds
      })

      if (!statusList) {
        title.value = t('查询状态失败')
        status.value = 'apierror'
        summary.value = 'data empty'
        return
      }

      summary.value = t('正在将变更内容同步至绑定模型，请稍等...')

      const undone = statusList?.some(task => ['new', 'waiting', 'executing'].includes(task.status))
      timer && clearTimeout(timer)
      if (undone) {
        timer = setTimeout(() => {
          getTaskSyncStatus(taskIds)
        }, 3000)
      } else {
        statusList.forEach((task) => {
          if (task.status === 'finished') {
            counts.value.success += 1
          }
          if (task.status === 'failure') {
            counts.value.fail += 1
          }
        })
        const success = counts.value.success === statusList.length
        const error = counts.value.fail === statusList.length

        if (success) {
          status.value = 'success'
        } else if (error) {
          status.value = 'error'
        } else {
          status.value = 'abnormal'
        }
        summary.value = ''
      }
    } catch (err) {
      status.value = 'apierror'
      summary.value = err?.message || ''
      title.value = t('查询状态失败')
    }
  }

  watchEffect(async () => {
    // 编辑模板允许绑定的模型为空
    if (isEditEmptyModel.value) {
      status.value = 'success'
      title.value = t('模板保存成功')
      return
    }
    try {
      status.value = 'loading'
      title.value = t(props.scene === 'edit' ? '模板编辑成功，正在同步至模型中...' : '信息正在同步至模型中...')

      const taskIds = await fieldTemplateService.syncModel({
        bk_template_id: props.templateId,
        object_ids: props.modelIds
      }, { requestId: requestIds.sync })

      summary.value = t('接下来，您可以在模板详情页面，查看模型的绑定的状态')

      getTaskSyncStatus(taskIds)
    } catch (err) {
      status.value = 'apierror'
      summary.value = err?.message || ''
      title.value = t('同步失败')
    }
  })

  const handleAction = (action) => {
    if (action === 'view') {
      routerActions.redirect({
        name: MENU_MODEL_FIELD_TEMPLATE,
        query: {
          action: 'view',
          id: props.templateId,
          tab: isEditEmptyModel.value ? 'field' : 'model'
        }
      })
    } else if (action === 'back') {
      routerActions.redirect({
        name: MENU_MODEL_FIELD_TEMPLATE
      })
    }
  }
</script>

<template>
  <div class="sync-results">
    <result-status :status="status">
      <template #title>
        <span v-if="!isSyncResultStatus || isEditEmptyModel">{{ title }}</span>
        <i18n path="字段组合模板同步结果" v-else>
          <template #success><em class="count">{{counts.success}}</em></template>
          <template #fail><em class="count">{{counts.fail}}</em></template>
        </i18n>
      </template>
      <template #summary>{{ summary }}</template>
      <template #actions>
        <bk-button theme="primary"
          :disabled="$loading(requestIds.sync)"
          @click="handleAction('view')">
          {{ $t('立即查看') }}
        </bk-button>
        <bk-button @click="handleAction('back')">{{ $t('返回列表') }}</bk-button>
      </template>
    </result-status>
  </div>
</template>

<style lang="scss" scoped>
  .sync-results {
    padding: 20px 24px;
  }
</style>
