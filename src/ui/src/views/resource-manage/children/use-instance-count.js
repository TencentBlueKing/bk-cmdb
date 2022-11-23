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

import { ref } from 'vue'
import CombineRequest from '@/api/combine-request.js'

// 每一片的大小（一次请求的最大实例个数）
const segment = 10

// 每次并发请求数（建议不超过6个）
const concurrency = 4

export const instanceCounts = ref([])

export default function useInstanceCount(state = {}, root) {
  const { modelIds } = state
  instanceCounts.value = []

  const fetchData = async function () {
    const requestIds = []
    const allResult = await CombineRequest.setup(Symbol(), (params) => {
      const requestId = `searchInstanceCount_${params.join()}`
      requestIds.push(requestId)
      return root.$store.dispatch('objectCommonInst/searchInstanceCount', {
        params: {
          condition: { obj_ids: params }
        },
        config: {
          requestId,
          globalError: false
        }
      })
    }, { segment, concurrency }).add(modelIds)

    // 关闭迭代器与取消请求
    root.$once('hook:beforeDestroy', () => {
      allResult?.return()
      root.$http.cancelRequest(requestIds)
    })

    // eslint-disable-next-line no-restricted-syntax
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
      // 一个批次更新一次
      instanceCounts.value.push(list)
    }
  }

  return {
    fetchData
  }
}
