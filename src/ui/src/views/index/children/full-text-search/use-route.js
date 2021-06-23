import { reactive, toRefs, watch } from '@vue/composition-api'

export default function useQuery(root) {
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
