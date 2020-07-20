import Vue from 'vue'
import equal from 'deep-equal'
import debounce from 'lodash.debounce'
import $http from '@/api'
import { TRANSFORM_TO_INTERNAL } from '@/dictionary/iam-auth'

const validCompareKey = ['resource_type', 'action', 'bk_biz_id', 'parent_layers', 'resource_id']
const uselessKey = ['resource_id_ex']
function filterUselessKey (data) {
    return JSON.parse(JSON.stringify(data), (key, value) => {
        if (key === '') return value
        if (uselessKey.includes(key)) return undefined
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
                const response = await $http.post('auth/verify', { resources: queue })
                authData = filterUselessKey(response)
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
                    component.updateAuth(Object.freeze(authResults), Object.freeze(authMetas))
                })
            }
        }
    }
})
