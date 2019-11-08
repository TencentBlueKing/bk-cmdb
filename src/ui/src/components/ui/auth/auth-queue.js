import Vue from 'vue'
import equal from 'deep-equal'
import debounce from 'lodash.debounce'
import $http from '@/api'
import { GET_AUTH_META } from '@/dictionary/auth'

export const deepEqual = equal

export default new Vue({
    data () {
        return {
            queue: [],
            authInstances: [],
            verify: null
        }
    },
    watch: {
        queue (queue) {
            if (!queue.length) return
            this.verify()
        }
    },
    created () {
        this.verify = debounce(this.getAuth, 20)
    },
    methods: {
        pushQueue (auth) {
            this.authInstances.push(auth)
            const repeat = this.queue.some(item => equal(item.data, auth.data))
            !repeat && this.queue.push(auth)
        },
        async getAuth () {
            const queue = [...this.queue]
            const authInstances = [...this.authInstances]
            this.queue = []
            this.authInstances = []
            const params = queue.map(item => {
                const types = Array.isArray(item.data.type) ? item.data.type : [item.data.type]
                return types.map(type => GET_AUTH_META(type, item.data))
            })
            const resources = params.reduce((acc, metas) => acc.concat(metas), [])
            const authData = await $http.post('auth/verify', { resources })
            authInstances.forEach(instance => {
                const findIndex = queue.findIndex(item => equal(item.data, instance.data))
                if (findIndex > -1) {
                    const types = Array.isArray(instance.data.type) ? instance.data.type : [instance.data.type]
                    const auths = []
                    types.forEach((type, index) => {
                        const authIndex = findIndex + index
                        authData[authIndex] && auths.push(authData[authIndex])
                    })
                    instance.component.updateAuth(auths)
                }
            })
        }
    }
})
