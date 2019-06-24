import cursor from '@/directives/cursor'
import { RESOURCE_TYPE_NAME, RESOURCE_ACTION_NAME, GET_AUTH_META } from '@/dictionary/auth'

const translateAuth = (authList = []) => {
    const authMap = {}
    authList.forEach(auth => {
        const meta = GET_AUTH_META(auth)
        if (authMap.hasOwnProperty(meta.resource_type)) {
            authMap[meta.resource_type].push(meta.action)
        } else {
            authMap[meta.resource_type] = [meta.action]
        }
    })
    return Object.keys(authMap).map(type => {
        return {
            scope: RESOURCE_TYPE_NAME[type],
            action: authMap[type].map(action => RESOURCE_ACTION_NAME[action]).join('，')
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
