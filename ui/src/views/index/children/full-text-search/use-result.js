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

import { computed, isRef, ref, unref } from '@vue/composition-api'
import debounce from 'lodash.debounce'
import { currentSetting as advancedSetting, allModelIds } from './use-advanced-setting.js'

const requestId = Symbol('fullTextSearch')

export default function useResult(state, root) {
  const { $store } = root

  const { route, keyword } = state

  const result = ref({})
  const fetching = ref(-1)

  // 如注入 keyword 则为输入联想模式
  const typing = computed(() => isRef(keyword))

  const queryKeyword = computed(() => (typing.value ? keyword.value : route.value.query.keyword))

  const params = computed(() => {
    const { query } = route.value
    const {
      c: queryObjId,
      k: kind,
      ps: limit = 10,
      p: page = 1
    } = query

    const kw = queryKeyword.value
    const nonLetter = /\W/.test(kw)
    // eslint-disable-next-line no-useless-escape
    const singleSpecial = /[!"#$%&'()\*,-\./:;<=>?@\[\\\]^_`{}\|~]{1}/
    const queryString = kw.length === 1 ? kw.replace(singleSpecial, '') : kw

    const filter = {}
    advancedSetting.targets.forEach((target) => {
      const key = `${target}s`
      filter[key] = advancedSetting[key].length ? advancedSetting[key] : unref(allModelIds)
    })
    const params = {
      filter,
      query_string: nonLetter ? `*${queryString}*` : queryString,
      page: {
        start: typing.value ? 0 : (page - 1) * limit,
        limit: typing.value ? 10 : Number(limit)
      }
    }

    if (queryObjId) {
      params.sub_resource = {
        [`${kind}s`]: queryObjId.split(',')
      }
    }

    return params
  })

  const getSearchResult = async () => {
    if (!params.value.query_string.length || !allModelIds.value.length) {
      return
    }

    try {
      fetching.value = true
      result.value = await $store.dispatch('fullTextSearch/search', {
        params: params.value,
        config: {
          requestId
        }
      })
    } finally {
      fetching.value = false
    }
  }

  const getSearchResultDebounce = debounce(getSearchResult, 200)

  return {
    result,
    fetching,
    getSearchResult: getSearchResultDebounce
  }
}
