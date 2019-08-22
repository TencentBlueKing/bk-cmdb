import * as OPERATION from '@/dictionary/auth'

const OPERATION_KEYS = Object.keys(OPERATION)
export default {
    computed: {
        $OPERATION () {
            const { authScope, operation } = this.$route.meta.auth
            const operationMap = {}
            operation.forEach(auth => {
                const key = OPERATION_KEYS.find(key => OPERATION[key] === auth)
                operationMap[key] = `${auth}.${authScope}`
            })
            return operationMap
        }
    },
    methods: {
        $isAuthorized (auth = '', option = { type: 'operation' }) {
            if (!auth) return true
            const types = Array.isArray(auth) ? auth : [auth]
            const authorized = types.map(auth => {
                return this.$store.getters['auth/isAuthorized'](auth, option)
            })
            return !authorized.some(auth => !auth)
        }
    }
}
