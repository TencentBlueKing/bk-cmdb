/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { computed, set, ref, watch, unref } from 'vue'
import CombineRequest from '@/api/combine-request.js'
import fieldTemplateService from '@/service/field-template'
import { useHttp } from '@/api'

export const statusList = ref([])
export const loadingMap = ref({})

export default function useModelSyncStatus(templateId, modelIdList) {
  // 去重，滚动加载时每次都是全量的数据
  const modelIds = computed(() => [...new Set(unref(modelIdList))])

  // 需要先重置
  statusList.value = []
  loadingMap.value = {}

  const requestIds = []

  let allResult = null

  const searchModelIds = computed(() => modelIds.value
    .filter(id => !statusList.value.some(item => item.object_id === id)))

  const pollingModelIds = computed(() => statusList.value
    .filter(item => isSyncing(item.status))
    .map(item => item.object_id))

  let pollingTimer = null

  const fetchStatus = async (ids) => {
    if (!ids.length) {
      return
    }

    ids.forEach((id) => {
      set(loadingMap.value, id, true)
    })

    const requestId = `searchModelSyncStatus_${templateId}_${ids.join()}`
    requestIds.push(requestId)

    allResult = await CombineRequest.setup(Symbol(), objIds => fieldTemplateService.getModelSyncStatus({
      bk_template_id: templateId,
      object_ids: objIds
    }, {
      requestId,
      globalError: true
    }), { segment: 5, concurrency: 5 }).add(ids)

    for (const result of allResult) {
      // 一个分组的执行结果
      const results = await result
      const list = []
      for (let i = 0; i < results.length; i++) {
        // 分组中的每一个执行结果
        const { status, reason, value } = results[i]
        if (status === 'rejected') {
          console.error(reason?.message)
          continue
        }
        list.push(...value ?? [])
      }

      // 因为会轮询，每次的查询结果可能在statusList已经存在，所以存在就更新否则才追加
      list.forEach((item) => {
        const index = statusList.value.findIndex(sta => sta.object_id === item.object_id)
        if (~index) {
          set(statusList.value, index, item)
        } else {
          statusList.value.push(item)
        }
      })
    }

    ids.forEach((id) => {
      set(loadingMap.value, id, false)
    })
  }

  const unwatchSearchModelIds = watch(searchModelIds, (ids) => {
    fetchStatus(ids)
  }, { immediate: true })

  const unwatchPollingModelIds = watch(pollingModelIds, (ids) => {
    if (pollingTimer) {
      clearTimeout(pollingTimer)
    }
    pollingTimer = setTimeout(() => {
      fetchStatus(ids)
    }, 5000)
  })

  const statusMap = computed(() => {
    const data = {}
    statusList.value.forEach((item) => {
      data[item.object_id] = item
    })
    return data
  })

  const http = useHttp()
  const clear = () => {
    allResult?.return()
    http.cancelRequest(requestIds)
    if (pollingTimer) {
      clearTimeout(pollingTimer)
    }
    unwatchSearchModelIds?.()
    unwatchPollingModelIds?.()
  }

  return {
    statusList,
    statusMap,
    loadingMap,
    clear
  }
}

export const isSyncing = status => ['new', 'waiting', 'executing'].includes(status)
export const isSynced = status => ['finished'].includes(status)
