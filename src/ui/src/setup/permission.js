import cursor from '@/directives/cursor'
import { IAM_ACTIONS } from '@/dictionary/iam-auth'
import { $error } from '@/magicbox'

const SYSTEM_ID = 'bk_cmdb'

// 前端构造的auth结构为：
// [{ type: 'xxx', relation: [xxx] }]
// 为了便于view中书写，其中relation可能存在三种格式:
// relation: [1, 2, ...] 表示该动作只关联一个视图，relation成员为视图拓扑路径上的资源ID，即关联单视图，操作单资源
// relation: [[1, 2], [3, 4], ...] 表示该动作只关联一个视图，relation中的成员为数组，每个数组代表一个视图的拓扑路径上的资源ID，即关联单视图，操作多资源
// relation: [[[1, 2], [3, 4]], [[1, 2], [5, 6]]] 表示该动作关联两个及以上的视图，为第二种情况的多视图场景，即关联多视图，操作多资源
// 因第一、第二种均为第三种的子场景，因此通过简单的类型判断转换为第三种形式
// 类型判断减少复杂度，只判断第一个元素的类型，不合法的混搭写法会报错
function convertRelation (relation = [], type) {
    if (!relation.length) return relation
    try {
        const [levelOne] = relation
        if (!Array.isArray(levelOne)) { // [1, 2, ...]的场景
            return [[relation]]
        }
        const [levelTwo] = levelOne
        if (!Array.isArray(levelTwo)) {
            return relation.map(data => [data])
        }
        return relation
    } catch (error) {
        $error('Convert resource relations fail, wrong params')
        console.error('Convert resource relations fail, wrong params:')
        console.error('auth type:', type)
        console.error('auth relation:', relation)
    }
}

// 将相同动作下的相同视图的实例合并到一起
function mergeSameActions (actions) {
    const actionMap = new Map()
    actions.forEach(action => {
        const viewMap = actionMap.get(action.id) || new Map()
        action.related_resource_types.forEach(({ type, instances }) => {
            const viewInstances = viewMap.get(type) || []
            viewInstances.push(...instances)
            viewMap.set(type, viewInstances)
        })
        actionMap.set(action.id, viewMap)
    })
    const permission = {
        system_id: SYSTEM_ID,
        actions: []
    }
    actionMap.forEach((viewMap, actionId) => {
        const relatedResourceTypes = []
        viewMap.forEach((viewInstances, viewType) => {
            relatedResourceTypes.push({
                type: viewType,
                system_id: SYSTEM_ID,
                instances: viewInstances
            })
        })
        permission.actions.push({
            id: actionId,
            related_resource_types: relatedResourceTypes
        })
    })
    return permission
}

export const translateAuth = auth => {
    const authList = Array.isArray(auth) ? auth : [auth]
    const actions = authList.map(({ type, relation = [] }) => {
        relation = convertRelation(relation, type)
        const definition = IAM_ACTIONS[type]
        const action = {
            id: definition.id,
            related_resource_types: []
        }
        if (!definition.relation) {
            return action
        }
        definition.relation.forEach((viewDefinition, viewDefinitionIndex) => { // 第m个视图的定义n
            const { view, instances } = viewDefinition
            const relatedResource = {
                type: view,
                instances: []
            }
            relation.forEach(resourceViewPaths => { // 第x个资源对应的视图数组
                const viewPathData = resourceViewPaths[viewDefinitionIndex] || [] // 取出第x个资源对应的第m个视图对应的拓扑路径ID数组
                const viewFullPath = viewPathData.map((path, pathIndex) => ({ // 资源x的第m个视图对应的全路径拓扑对象
                    type: instances[pathIndex],
                    id: String(path)
                }))
                relatedResource.instances.push(viewFullPath)
            })
            action.related_resource_types.push(relatedResource)
        })
        return action
    })
    return mergeSameActions(actions)
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
