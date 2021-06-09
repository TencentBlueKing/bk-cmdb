/* eslint-disable no-unused-vars */
import { ref } from '@vue/composition-api'
import CombineRequest from '@/api/combine-request.js'

const requestId = Symbol()

// 每一片的大小（一次请求的最大实例个数）
const segment = 10

// 每次并发请求数（建议不超过6个）
const concurrency = 4

export const instanceCounts = ref([])

export default function useInstanceCount(state = {}, root) {
  const { modelIds } = state
  instanceCounts.value = []

  const fetchData = async function () {
    const allResult = await CombineRequest.setup(requestId, params => root.$store.dispatch('objectCommonInst/searchInstanceCount', {
      params: {
        condition: { obj_ids: params }
      },
      config: {
        requestId,
        globalError: false
      }
    }), { segment, concurrency }).add(modelIds)

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
        list.push(...value)
      }
      // 一个批次更新一次
      instanceCounts.value.push(list)
    }
  }

  return {
    fetchData
  }
}
