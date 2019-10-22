import Vue from 'vue'
import $http from '@/api'
import { GET_AUTH_META } from '@/dictionary/auth'

export default new Vue({
    data () {
        return {
            queue: [],
            authInstances: []
        }
    },
    watch: {
        queue () {
            if (!this.queue.length) return
            const timer = setTimeout(async () => {
                const resources = this.queue.map(item => {
                    const meta = GET_AUTH_META(item.resource.type)
                    if (meta.scope === 'business' && item.bizId) {
                        meta.bk_biz_id = item.bizId
                    }
                    delete meta.scope
                    return meta
                })
                const authData = await $http.post('auth/verify', { resources }).finally(() => {
                    clearTimeout(timer)
                })
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
            }, 100)
        }
    },
    methods: {
        addQueue (auth) {
            this.authInstances.push(auth)
            const id = auth.id
            const repeat = this.queue.findIndex(item => item.id === id)
            repeat < 0 && this.queue.push(auth)
        }
    }
})
