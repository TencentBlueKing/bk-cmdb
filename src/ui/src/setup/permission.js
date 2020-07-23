import cursor from '@/directives/cursor'
import { IAM_ACTIONS } from '@/dictionary/iam-auth'

const SYSTEM_ID = 'bk_cmdb'

export const translateAuth = auth => {
    const authList = Array.isArray(auth) ? auth : [auth]
    const permission = {
        system_id: SYSTEM_ID,
        actions: authList.map(({ type, relation = [] }) => {
            const definition = IAM_ACTIONS[type]
            const action = {
                id: definition.id,
                related_resource_types: []
            }
            if (!definition.relation) {
                return action
            }
            definition.relation.forEach(({ view, instances }, relationIndex) => {
                const relationIds = relation[relationIndex]
                if (!relationIds) {
                    return false
                }
                const ids = Array.isArray(relationIds) ? relationIds : [relationIds]
                const relatedResource = {
                    system_id: SYSTEM_ID,
                    type: view,
                    instances: []
                }
                instances.forEach((instance, instanceIndex) => {
                    const instanceId = ids[instanceIndex]
                    instanceId && relatedResource.instances.push({
                        type: instance,
                        id: String(instanceId)
                    })
                })
                action.related_resource_types.push([relatedResource])
            })
            return action
        })
    }
    return permission
}

cursor.setOptions({
    globalCallback: options => {
        const permission = translateAuth(options.auth)
        const permissionModal = window.permissionModal
        permissionModal && permissionModal.show(permission)
    },
    x: 16,
    y: 8
})
