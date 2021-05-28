import { computed, isRef, ref } from '@vue/composition-api'
import debounce from 'lodash.debounce'

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
    const { c: queryObjId, ps: limit = 10, p: page = 1 } = query
    const kw = queryKeyword.value
    const nonLetter = /\W/.test(kw)
    // eslint-disable-next-line no-useless-escape
    const singleSpecial = /[!"#$%&'()\*,-\./:;<=>?@\[\\\]^_`{}\|~]{1}/
    const queryString = kw.length === 1 ? kw.replace(singleSpecial, '') : kw

    const params = {
      page: {
        start: typing.value ? 0 : (page - 1) * limit,
        limit: typing.value ? 10 : Number(limit)
      },
      bk_obj_id: '',
      bk_biz_id: '',
      query_string: nonLetter ? `*${queryString}*` : queryString
    }

    if (typing.value) {
      return params
    }

    let objId = ''
    let filter = []
    if (queryObjId) {
      const isNormalModelInst = !['model', 'host', 'biz'].includes(queryObjId)
      objId = isNormalModelInst ? queryObjId : ''
      filter = isNormalModelInst ? ['instance'] : [queryObjId]
    }

    params.bk_obj_id = objId
    params.filter = filter
    return params
  })

  const getSearchResult = async () => {
    if (!queryKeyword.value.length) return
    try {
      fetching.value = true
      result.value = await $store.dispatch('fullTextSearch/search', {
        params: params.value,
        config: {
          requestId: 'search',
          cancelPrevious: true
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
