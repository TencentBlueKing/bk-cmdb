import { reactive, toRefs, watch } from '@vue/composition-api'

export const pickQuery = (query = {}, include = [], exclude = []) => {
  let queryKeys = Object.keys(query)
  if (include.length) {
    queryKeys = queryKeys.filter(item => include.includes(item))
  }
  if (exclude.length) {
    queryKeys = queryKeys.filter(item => !exclude.includes(item))
  }
  const newQuery = {}
  queryKeys.forEach(key => newQuery[key] = query[key])
  return newQuery
}

export default function useRoute(root) {
  const state = reactive({ route: root.$route })

  watch(
    () => root.$route,
    (route) => {
      state.route = route
    }
  )

  return {
    ...toRefs(state)
  }
}
