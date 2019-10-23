import Vue from 'vue'
import equal from 'deep-equal'
import $http from '@/api'
import { GET_AUTH_META } from '@/dictionary/auth'

export const deepEqual = equal

export default new Vue({
    data () {
        return {
            queue: [],
            authInstances: []
        }
    },
    watch: {
        async queue () {
            if (!this.queue.length) return
            const resources = this.queue.map(item => {
                const authResource = item.resource
                const meta = GET_AUTH_META(authResource.type)
                if (meta.scope === 'business' && authResource.bizId) {
                    meta.bk_biz_id = authResource.bizId
                }
                delete meta.scope
                return meta
            })
            const authData = await $http.post('auth/verify', { resources })
            this.authInstances.forEach(task => {
                const auth = authData.find(item => {
                    const compose = [item.resource_type, item.action]
                    if (item.bk_biz_id) {
                        compose.push('business')
                    } else {
                        compose.push('global')
                    }
                    const id = `${compose.join('.')}-${item.resource_id}`
                    return task.id === id
                })
                auth && task.component.updateAuth(auth)
            })
            this.queue = []
            this.authInstances = []
        }
    },
    methods: {
        pushQueue (auth) {
            this.authInstances.push(auth)
            const repeat = this.queue.some(item => equal(item.resource, auth.resource))
            !repeat && this.queue.push(auth)
        }
    }
})
