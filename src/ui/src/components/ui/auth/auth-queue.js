import Vue from 'vue'
import equal from 'deep-equal'
import debounce from 'lodash.debounce'
import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'

function filterUselssKey (data, uselessKeys) {
    return JSON.parse(JSON.stringify(data), (key, value) => {
        if (key === '') return value
        if (uselessKeys.includes(key)) return undefined
        return value
    })
}

export default new Vue({
    data () {
        return {
            queue: [],
            authComponents: [],
            verify: debounce(this.getAuth, 20)
        }
    },
    watch: {
        queue (queue) {
            this.verify()
        }
    },
    methods: {
        add ({ component, data }) {
            this.authComponents.push(component)
            const authMetas = TRANSFORM_TO_INTERNAL(data)
            authMetas.forEach(meta => {
                const exist = this.queue.some(exist => equal(meta, exist))
                if (!exist) {
                    this.queue.push(meta)
                }
            })
        },
        async getAuth () {
            if (!this.queue.length) return
            const queue = this.queue.splice(0)
            const authComponents = this.authComponents.splice(0)
            let authData = []
            try {
                authData = await $http.post('auth/verify', { resources: queue })
            } catch (error) {
                console.error(error)
            } finally {
                authComponents.forEach(component => {
                    const authMetas = TRANSFORM_TO_INTERNAL(component.auth)
                    const authResults = []
                    authMetas.forEach(meta => {
                        const result = authData.find(result => {
                            const source = {}
                            const target = {}
                            Object.keys(meta).forEach(key => {
                                source[key] = meta[key]
                                if (key === 'parent_layers') {
                                    target[key] = filterUselssKey(result[key], ['resource_id_ex'])
                                } else {
                                    target[key] = result[key]
                                }
                            })
                            return equal(source, target)
                        })
                        if (result) {
                            authResults.push(result)
                        }
                    })
                    component.updateAuth(Object.freeze(authResults), Object.freeze(authMetas))
                })
            }
        }
    }
})
