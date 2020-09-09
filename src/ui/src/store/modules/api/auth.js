import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'
const actions = {
    async getViewAuth ({ rootGetters }, viewAuthData) {
        if (rootGetters.site.authscheme !== 'iam') {
            return Promise.resolve(true)
        }
        const result = await $http.post('auth/verify', {
            resources: TRANSFORM_TO_INTERNAL(viewAuthData)
        })
        return Promise.resolve(result.every(data => data.is_pass))
    },
    async getSkipUrl (context, { params, config = {} }) {
        return $http.post('auth/skip_url', params, config)
    }
}

export default {
    namespaced: true,
    actions
}
