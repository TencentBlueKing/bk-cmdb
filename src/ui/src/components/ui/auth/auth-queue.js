import Vue from 'vue'
import equal from 'deep-equal'
import debounce from 'lodash.debounce'
import $http from '@/api'
import { GET_AUTH_META } from '@/dictionary/auth'

export const deepEqual = equal

function transformAuthMetas (data) {
    const authTypes = Array.isArray(data.type) ? data.type : [data.type]
    return authTypes.map(authType => GET_AUTH_META(authType, data))
}

const validCompareKey = ['resource_type', 'action', 'bk_biz_id', 'parent_layers', 'resource_id']

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
        pushQueue ({ component, data }) {
            this.authInstances.push(component)
            const authMetas = transformAuthMetas(data)
            authMetas.forEach(meta => {
                const exist = this.queue.some(exist => equal(meta, exist))
                if (!exist) {
                    this.queue.push(meta)
                }
            })
        },
        async getAuth () {
            const queue = [...this.queue]
            const authInstances = [...this.authInstances]
            this.queue = []
            this.authInstances = []
            const authData = await $http.post('auth/verify', { resources: queue })
            authInstances.forEach(instance => {
                const authMetas = transformAuthMetas(instance.auth)
                const authResults = []
                authMetas.forEach(meta => {
                    const result = authData.find(result => {
                        const source = {}
                        const target = {}
                        validCompareKey.forEach(key => {
                            if (meta.hasOwnProperty(key)) {
                                source[key] = meta[key]
                                target[key] = result[key]
                            }
                        })
                        return equal(source, target)
                    })
                    if (result) {
                        authResults.push(result)
                    }
                })
                instance.updateAuth(authResults)
            })
        }
    }
})
