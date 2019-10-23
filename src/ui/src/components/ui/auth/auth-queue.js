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
                const meta = GET_AUTH_META(item.data.type, item.data)
                delete meta.scope
                return meta
            })
            const authData = await $http.post('auth/verify', { resources })
            this.authInstances.forEach(instance => {
                const index = this.queue.findIndex(item => equal(item.resource, instance.resource))
                if (index > -1) {
                    instance.component.updateAuth(authData[index])
                }
            })
            this.queue = []
            this.authInstances = []
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
