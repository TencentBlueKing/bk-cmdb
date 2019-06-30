import cursor from '@/directives/cursor'
import { language } from '@/i18n'
import { RESOURCE_TYPE_NAME, RESOURCE_ACTION_NAME, GET_AUTH_META } from '@/dictionary/auth'
const SCOPE_NAME = language === 'en' ? {
    global: 'global',
    business: 'business'
} : {
    global: '全局',
    business: '业务'
}
const translateAuth = (authList = []) => {
    const authMap = {}
    authList.forEach(auth => {
        const { resource_type: resourceType, action, scope } = GET_AUTH_META(auth)
        if (authMap.hasOwnProperty(resourceType)) {
            authMap[resourceType].actions.push(action)
        } else {
            authMap[resourceType] = {
                scope: SCOPE_NAME[scope],
                actions: [action]
            }
        }
    })
    return Object.keys(authMap).map(type => {
        return {
            scope: authMap[type].scope,
            resource: RESOURCE_TYPE_NAME[type],
            action: authMap[type].actions.map(action => RESOURCE_ACTION_NAME[action]).join('，')
        }
    })
}

cursor.setOptions({
    globalCallback: options => {
        const permissionModal = window.permissionModal
        permissionModal && permissionModal.show(translateAuth(options.auth), false)
    },
    x: 16,
    y: 8
})
