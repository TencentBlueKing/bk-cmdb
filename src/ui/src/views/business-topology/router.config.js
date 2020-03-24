import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_TRANSFER_HOST,
    MENU_BUSINESS_HOST_AND_SERVICE,
    MENU_BUSINESS_HOST_DETAILS,
    MENU_BUSINESS_DELETE_SERVICE
} from '@/dictionary/menu-symbol'

import {
    U_HOST,
    C_TOPO,
    U_TOPO,
    D_TOPO,
    C_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    HOST_TO_RESOURCE
} from '@/dictionary/auth'

export default [{
    name: MENU_BUSINESS_HOST_AND_SERVICE,
    path: 'index',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '业务拓扑'
        },
        auth: {
            operation: {
                U_HOST,
                C_TOPO,
                U_TOPO,
                D_TOPO,
                C_SERVICE_INSTANCE,
                R_SERVICE_INSTANCE,
                U_SERVICE_INSTANCE,
                D_SERVICE_INSTANCE,
                HOST_TO_RESOURCE
            }
        },
        customInstanceColumn: 'business_topology_table_column_config',
        customFilterProperty: 'business_topology_filter_property_config'
    })
}, {
    name: MENU_BUSINESS_TRANSFER_HOST,
    path: 'host/transfer/:type',
    component: () => import('@/views/host-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        auth: {
            operation: {
                U_HOST
            }
        },
        layout: {
            previous: (view) => {
                if (view.$route.query.from) {
                    return { path: view.$route.query.from }
                }
                return {
                    name: MENU_BUSINESS_HOST_AND_SERVICE,
                    query: { node: view.$route.query.node }
                }
            }
        }
    })
}, {
    name: MENU_BUSINESS_DELETE_SERVICE,
    path: 'service/delete/:moduleId?/:ids',
    component: () => import('@/views/service-operation/index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '删除服务实例',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        layout: {
            previous () {
                const $route = window.CMDB_APP.$route
                if ($route.query.from) {
                    return {
                        path: $route.query.from,
                        query: $route.query.query
                    }
                } else {
                    return {
                        name: MENU_BUSINESS_HOST_AND_SERVICE,
                        query: {
                            tab: 'serviceInstance',
                            node: `module-${$route.params.moduleId}`
                        }
                    }
                }
            }
        }
    })
}, {
    name: MENU_BUSINESS_HOST_DETAILS,
    path: ':business/host/:id',
    component: () => import('@/views/host-details/index'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机详情',
            relative: MENU_BUSINESS_HOST_AND_SERVICE
        },
        layout: {
            previous: (view) => ({
                name: MENU_BUSINESS_HOST_AND_SERVICE,
                query: view.$route.query
            })
        },
        checkAvailable: (to, from, app) => {
            return parseInt(to.params.business) === app.$store.getters['objectBiz/bizId']
        }
    })
}]
