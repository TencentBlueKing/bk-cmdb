import Meta from '@/router/meta'
import {
    MENU_BUSINESS
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
        tableColumnsConfigKey: 'business_topology_table_column_config'
    })
}, {
    name: 'removeHost',
    path: 'topologyNew/removeHost',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '移除主机'
        },
        type: 'remove'
    })
}]
