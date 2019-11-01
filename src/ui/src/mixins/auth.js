import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
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
        },
        $authResources (resources = {}) {
            if (typeof resources !== 'object') return resources
            const auth = {}
            const isAdminview = this.$route.matched.length && this.$route.matched[0].name !== MENU_BUSINESS
            const bizId = this.$store.getters['objectBiz/bizId']
            if (bizId && !isAdminview) {
                auth.bk_biz_id = bizId
            }
            return Object.assign(auth, resources)
        }
    }
}
