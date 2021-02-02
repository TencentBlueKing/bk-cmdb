import Vue from 'vue'
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

function equal (source, target) {
    const {
        resource_type: SResourceType,
        resource_id: SResourceId,
        action: SAction,
        bk_biz_id: SBizId
    } = source
    const {
        resource_type: TResourceType,
        resource_id: TResourceId,
        action: TAction,
        bk_biz_id: TBizId
    } = target
    const SParentLayers = source.parent_layers || []
    const TParentLayers = target.parent_layers || []
    if (
        SResourceType !== TResourceType
        || SResourceId !== TResourceId
        || SAction !== TAction
        || SBizId !== TBizId
        || SParentLayers.length !== TParentLayers.length
    ) {
        return false
    }
    return SParentLayers.every((_, index) => {
        const SParentLayersMeta = SParentLayers[index]
        const TParentLayersMeta = TParentLayers[index]
        return Object.keys(SParentLayersMeta).every(key => SParentLayersMeta[key] === TParentLayersMeta[key])
    })
}

function unique (data) {
    return data.reduce((queue, meta) => {
        const exist = queue.some(exist => equal(exist, meta))
        if (!exist) {
            queue.push(meta)
        }
        return queue
    }, [])
}

export const AuthRequestId = Symbol('auth_request_id')

const authEnable = window.Site.authscheme === 'iam'
let afterVerifyQueue = []
export function afterVerify (func, once = true) {
    if (authEnable) {
        afterVerifyQueue.push({
            handler: func,
            once
        })
    } else {
        func()
    }
}
function execAfterVerify (authData) {
    afterVerifyQueue.forEach(({ handler }) => handler(authData))
    afterVerifyQueue = afterVerifyQueue.filter(({ once }) => !once)
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
            this.queue.push(...authMetas)
        },
        async getAuth () {
            if (!this.queue.length) return
            const queue = unique(this.queue.splice(0))
            const authComponents = this.authComponents.splice(0)
            let authData = []
            try {
                authData = await $http.post('auth/verify', { resources: queue }, { requestId: AuthRequestId })
            } catch (error) {
                console.error(error)
            } finally {
                authData = filterUselssKey(authData, ['resource_id_ex'])
                authComponents.forEach(component => {
                    const authMetas = TRANSFORM_TO_INTERNAL(component.auth)
                    const authResults = []
                    authMetas.forEach(meta => {
                        const result = authData.find(result => {
                            const source = {}
                            const target = {}
                            Object.keys(meta).forEach(key => {
                                source[key] = meta[key]
                                target[key] = result[key]
                            })
                            return equal(source, target)
                        })
                        if (result) {
                            authResults.push(result)
                        }
                    })
                    component.updateAuth(Object.freeze(authResults), Object.freeze(authMetas))
                })
                execAfterVerify(authData)
            }
        }
    }
})
