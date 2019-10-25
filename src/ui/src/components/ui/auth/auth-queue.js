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
            const queue = [...this.queue]
            const authInstances = [...this.authInstances]
            this.queue = []
            this.authInstances = []
            const resources = queue.map(item => {
                const meta = GET_AUTH_META(item.data.type, item.data)
                delete meta.scope
                return meta
            })
            const authData = await $http.post('auth/verify', { resources })
            authInstances.forEach(instance => {
                const index = queue.findIndex(item => equal(item.resource, instance.resource))
                if (index > -1) {
                    instance.component.updateAuth(authData[index])
                }
            })
        }
    },
    methods: {
        pushQueue (auth) {
            this.authInstances.push(auth)
            const repeat = this.queue.some(item => equal(item.data, auth.data))
            !repeat && this.queue.push(auth)
        }
    }
})
