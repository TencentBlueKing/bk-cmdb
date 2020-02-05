import Meta from '@/router/meta'
import {
    MENU_BUSINESS,
    MENU_BUSINESS_HOST,
    MENU_BUSINESS_HOST_APPLY,
    MENU_BUSINESS_HOST_APPLY_EDIT,
    MENU_BUSINESS_HOST_APPLY_CONFIRM,
    MENU_BUSINESS_HOST_APPLY_CONFLICT,
    MENU_BUSINESS_HOST_APPLY_FAILED
} from '@/dictionary/menu-symbol'
import {
    U_HOST_APPLY
} from '@/dictionary/auth'

export default [{
    name: MENU_BUSINESS_HOST_APPLY,
    path: 'host-apply',
    component: () => import('./index.vue'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST
        },
        auth: {
            operation: {
                U_HOST_APPLY
            }
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_CONFIRM,
    path: 'host-apply/confirm',
    component: () => import('./property-confirm'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        },
        auth: {
            operation: {
                U_HOST_APPLY
            }
        },
        layout: {
            previous (view) {
                return new Promise((resolve, reject) => {
                    view.leaveConfirmConfig.active = false
                    view.$nextTick(() => {
                        const config = {
                            name: MENU_BUSINESS_HOST_APPLY_EDIT,
                            query: {
                                mid: view.$route.query.mid
                            }
                        }
                        if (view.isBatch) {
                            config.query.batch = 1
                        }
                        resolve(config)
                    })
                })
            }
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_EDIT,
    path: 'host-apply/edit',
    component: () => import('./edit'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        },
        auth: {
            operation: {
                U_HOST_APPLY
            }
        },
        layout: {
            previous (view) {
                const config = {
                    name: MENU_BUSINESS_HOST_APPLY,
                    query: {}
                }
                if (String(view.$route.query.mid).indexOf(',') === -1) {
                    config.query.module = view.$route.query.mid
                }
                return config
            }
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_CONFLICT,
    path: 'host-apply/conflict',
    component: () => import('./conflict-list'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        },
        auth: {
            operation: {
                U_HOST_APPLY
            }
        },
        layout: {
            previous (view) {
                return {
                    name: MENU_BUSINESS_HOST_APPLY,
                    query: {
                        module: view.$route.query.mid
                    }
                }
            }
        }
    })
}, {
    name: MENU_BUSINESS_HOST_APPLY_FAILED,
    path: 'host-apply/failed',
    component: () => import('./failed-list'),
    meta: new Meta({
        owner: MENU_BUSINESS,
        menu: {
            i18n: '主机属性自动应用',
            parent: MENU_BUSINESS_HOST_APPLY
        },
        auth: {
            operation: {
                U_HOST_APPLY
            }
        },
        layout: {
            previous: {
                name: MENU_BUSINESS_HOST_APPLY
            }
        }
    })
}]
