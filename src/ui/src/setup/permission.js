import cursor from '@/directives/cursor'
import http from '@/api'
import { language } from '@/i18n'
import { RESOURCE_TYPE_NAME, RESOURCE_ACTION_NAME, GET_AUTH_META } from '@/dictionary/auth'
import { MENU_BUSINESS } from '@/dictionary/menu-symbol'
import store from '@/store'
const SCOPE_NAME = language === 'en' ? {
    global: 'global',
    business: 'business'
} : {
    global: '全局',
    business: '业务'
}

const getScope = () => {
    return window.CMDB_APP.$route.meta.owner === MENU_BUSINESS ? 'business' : 'global'
}

const convertAuth = authList => {
    const scope = getScope()
    return http.post('auth/convert', {
        data: authList.map(auth => {
            const { resource_type: type, action } = GET_AUTH_META(auth)
            return {
                scope: scope === 'global' ? 'system' : 'biz',
                attribute: {
                    type,
                    action
                }
            }
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
        const scope = getScope()
        return authList.map((auth, index) => {
            const { resource_type: resourceType, action } = GET_AUTH_META(auth)
            return {
                action_id: convertedAuth[index].action,
                action_name: RESOURCE_ACTION_NAME[action],
                scope_id: scope === 'global' ? 'bk_cmdb' : business.bk_biz_id ? String(business.bk_biz_id) : '',
                scope_name: scope === 'global' ? '配置平台' : business.bk_biz_name,
                scope_type: scope === 'global' ? 'system' : 'biz',
                scope_type_name: SCOPE_NAME[scope],
                system_id: 'bk_cmdb',
                system_name: '配置平台',
                resource_type: convertedAuth[index].type,
                resource_type_name: RESOURCE_TYPE_NAME[resourceType],
                resources: [[{
                    resource_type_name: RESOURCE_TYPE_NAME[resourceType],
                    resource_type: convertedAuth[index].type,
                    resource_id: '',
                    resource_name: ''
                }]]
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
