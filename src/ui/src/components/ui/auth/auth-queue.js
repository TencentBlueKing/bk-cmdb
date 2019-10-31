import Vue from 'vue'
import equal from 'deep-equal'
import $http from '@/api'
import { GET_AUTH_META } from '@/dictionary/auth'

export const deepEqual = equal

const debounce = (fn, delay) => {
    let timer = null
    return function () {
        const _this = this
        const args = arguments
        clearTimeout(timer)
        timer = setTimeout(() => {
            fn.apply(_this, args)
        }, delay)
    }
}

export default new Vue({
    data () {
        return {
            queue: [],
            authInstances: [],
            verify: null
        }
    },
    watch: {
        queue () {
            if (!this.queue.length) return
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
                const metas = types.map(type => {
                    const meta = GET_AUTH_META(type, item.data)
                    return meta
                })
                return metas
            })
            const resources = []
            params.forEach(metas => {
                resources.push(...metas)
            })
            const authData = await $http.post('auth/verify', { resources })
            authInstances.forEach(instance => {
                const findIndex = queue.findIndex(item => equal(item.data, instance.data))
                const types = Array.isArray(instance.data.type) ? instance.data.type : [instance.data.type]
                if (findIndex > -1) {
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
