import cursor from '@/directives/cursor'
import http from '@/api'
import { language } from '@/i18n'
import { RESOURCE_TYPE_NAME, RESOURCE_ACTION_NAME, GET_AUTH_META } from '@/dictionary/auth'
import store from '@/store'
const SCOPE_NAME = language === 'en' ? {
    global: 'global',
    business: 'business'
} : {
    global: '全局',
    business: '业务'
}

const convertAuth = authList => {
    return http.post('auth/convert', {
        data: authList.map(auth => {
            const { resource_type: type, action, scope } = GET_AUTH_META(auth)
            return { type, action, business_id: scope === 'global' ? 0 : store.getters['objectBiz/bizId'] }
        })
    })
}

export const translateAuth = async (authList = []) => {
    if (!authList.length) {
        return authList
    }
    try {
        const convertedAuth = await convertAuth(authList)
        const business = store.state.objectBiz.authorizedBusiness.find(business => business.bk_biz_id === store.getters['objectBiz/bizId']) || {}
        return authList.map((auth, index) => {
            const { resource_type: resourceType, action, scope } = GET_AUTH_META(auth)
            return {
                action_id: convertedAuth[index].action,
                action_name: RESOURCE_ACTION_NAME[action],
                scope_id: scope === 'global' ? 'bk_cmdb' : String(business.bk_biz_id),
                scope_name: scope === 'global' ? '配置平台' : business.bk_biz_name,
                scope_type: scope === 'global' ? 'system' : 'biz',
                scope_type_name: SCOPE_NAME[scope],
                system_id: 'bk_cmdb',
                system_name: '配置平台',
                resource_type: convertedAuth[index].type,
                resource_type_name: RESOURCE_TYPE_NAME[resourceType],
                resources: [[]]
            }
        })
    } catch (e) {
        console.error(e)
        return []
    }
}

cursor.setOptions({
    globalCallback: async options => {
        const permission = await translateAuth(options.auth)
        const permissionModal = window.permissionModal
        permissionModal && permissionModal.show(permission)
    },
    x: 16,
    y: 8
})
