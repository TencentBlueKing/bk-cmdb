import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_TRANSFER_HOST
} from '@/dictionary/menu-symbol'

export default [{
    name: 'businessTopologyNew',
    path: 'topologyNew',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务拓扑'
        },
        customInstanceColumn: 'business_topology_table_column_config'
    })
}, {
    name: MENU_BUSINESS_TRANSFER_HOST,
    path: 'host/transfer/:type',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS
    })
}]
