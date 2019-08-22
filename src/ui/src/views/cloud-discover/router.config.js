import Meta from '@/router/meta'
import { NAV_BUSINESS_RESOURCE } from '@/dictionary/menu'
import {
    C_CLOUD_DISCOVER,
    U_CLOUD_DISCOVER,
    D_CLOUD_DISCOVER,
    R_CLOUD_DISCOVER
} from '@/dictionary/auth'

export const OPERATION = {
    R_CLOUD_DISCOVER,
    C_CLOUD_DISCOVER,
    U_CLOUD_DISCOVER,
    D_CLOUD_DISCOVER
}

const path = '/cloud-discover'

export default {
    name: 'cloudDiscover',
    path: path,
    component: () => import('./index.vue'),
    meta: new Meta({
        menu: {
            id: 'cloudDiscover',
            i18n: '云资源发现',
            path: path,
            order: 4,
            parent: NAV_BUSINESS_RESOURCE,
            adminView: false
        },
        auth: {
            operation: Object.values(OPERATION),
            setAuthScope () {
                this.authScope = 'global'
            }
        },
        requireBusiness: true,
        i18nTitle: '云资源发现'
    })
}
