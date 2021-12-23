import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
const actions = {
  async getViewAuth(context, viewAuthData) {
    if (window.Site.authscheme !== 'iam') {
      return Promise.resolve(true)
    }
    const result = await $http.post('auth/verify', {
      // eslint-disable-next-line new-cap
      resources: TRANSFORM_TO_INTERNAL(viewAuthData)
    })
    return Promise.resolve(result.every(data => data.is_pass))
  },
  async getSkipUrl(context, { params, config = {} }) {
    return $http.post('auth/skip_url', params, config)
  }
}

export default {
  namespaced: true,
  actions
}
