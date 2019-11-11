import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
import * as OPERATION from '@/dictionary/auth'
export default {
    computed: {
        $OPERATION () {
            return OPERATION
        }
    },
    methods: {
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
