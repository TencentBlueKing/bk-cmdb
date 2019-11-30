import $http from '@/api'
import { $error } from '@/magicbox'

const actions = {
    async getViewAuth ({ rootGetters }, viewAuthData) {
        if (rootGetters.site.authscheme !== 'iam') {
            return Promise.resolve(true)
        }
        const result = await $http.post('auth/verify', {
            resources: [viewAuthData]
        })
        return Promise.resolve(result.every(data => data.is_pass))
    },
    async getSkipUrl (context, { params, config = {} }) {
        try {
            const url = await $http.post('auth/skip_url', params, Object.assign(config, { globalError: false }))
            if (url.indexOf('tid') === -1) {
                return url + '?system_id=bk_cmdb&apply_way=custom'
            }
            return url
        } catch (e) {
            const url = (window.Site.authCenter || {}).url
            if (url) {
                return url + '?system_id=bk_cmdb&apply_way=custom'
            }
            $error(e.message)
            throw e
        }
    }
}

export default {
    namespaced: true,
    actions
}
