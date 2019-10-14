export default {
    computed: {
        $OPERATION () {
            const { authScope, operation, view } = this.$route.meta.auth
            const operationMap = {
                ...operation,
                ...view
            }
            Object.keys(operationMap).forEach(key => {
                operationMap[key] = `${operationMap[key]}.${authScope}`
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
