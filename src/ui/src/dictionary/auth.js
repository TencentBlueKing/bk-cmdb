// 模型分组
export const C_MODEL_GROUP = 'modelClassification.create'
export const U_MODEL_GROUP = 'modelClassification.update'
export const D_MODEL_GROUP = 'modelClassification.delete'

// 模型
export const C_MODEL = 'model.create'
export const R_MODEL = 'model.findMany'
export const U_MODEL = 'model.update'
export const D_MODEL = 'model.delete'

// 实例
export const C_INST = 'modelInstance.create'
export const U_INST = 'modelInstance.update'
export const D_INST = 'modelInstance.delete'
export const R_INST = 'modelInstance.findMany'

// 动态分组
export const C_CUSTOM_QUERY = 'dynamicGrouping.create'
export const U_CUSTOM_QUERY = 'dynamicGrouping.update'
export const D_CUSTOM_QUERY = 'dynamicGrouping.delete'
export const R_CUSTOM_QUERY = 'dynamicGrouping.findMany'

// 进程管理
export const C_PROCESS = 'process.create'
export const U_PROCESS = 'process.update'
export const D_PROCESS = 'process.delete'
export const R_PROCESS = 'process.findMany'
export const PROCESS_BIND_MODULE = 'process.boundModuleToProcess'
export const PROCESS_UNBIND_MODULE = 'process.unboundModelToProcess'
export const PROCESS_SEARCH_MODULE = 'process.findBoundModuleProcess'

// 业务拓扑
export const C_TOPO = 'mainlineInstanceTopology.create'
export const U_TOPO = 'mainlineInstanceTopology.update'
export const D_TOPO = 'mainlineInstanceTopology.delete'
export const R_TOPO = 'mainlineInstanceTopology.findMany'
export const TOPO_TRANSFER_HOST = 'mainlineInstanceTopology.transferHost'

// 主机管理
export const C_HOST = 'hostInstance.create'
export const R_HOST = 'hostInstance.findMany'
export const U_HOST = 'hostInstance.update'
export const D_HOST = 'hostInstance.delete'
export const HOST_TO_RESOURCE = 'hostInstance.moveHostFromModuleToResPool'
export const HOST_ASSIGN = 'hostInstance.moveResPoolHostToBizIdleModule'

// 关联类型
export const C_RELATION = 'associationType.create'
export const U_RELATION = 'associationType.update'
export const D_RELATION = 'associationType.delete'

// 业务
export const C_BUSINESS = 'business.create'
export const U_BUSINESS = 'business.update'
export const R_BUSINESS = 'business.findMany'
export const BUSINESS_ARCHIVE = 'business.archive'

// 事件推送
export const C_EVENT = 'eventPushing.create'
export const U_EVENT = 'eventPushing.update'
export const D_EVENT = 'eventPushing.delete'
export const R_EVENT = 'eventPushing.findMany'

// 操作审计
export const R_AUDIT = 'auditlog.findMany'

// 系统基础
export const SYSTEM_TOPOLOGY = 'systemBase.modelTopologyOperation'
export const SYSTEM_MANAGEMENT = 'systemBase.adminEntrance'
export const SYSTEM_MODEL_GRAPHICS = 'systemBase.modelTopologyView'

export const HAS_BUSINESS_MODE = [
    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,
    C_MODEL,
    U_MODEL,
    D_MODEL,
    C_INST,
    U_INST,
    D_INST,
    R_INST,
    C_CUSTOM_QUERY,
    U_CUSTOM_QUERY,
    D_CUSTOM_QUERY,
    R_CUSTOM_QUERY,
    C_PROCESS,
    U_PROCESS,
    D_PROCESS,
    R_PROCESS,
    PROCESS_BIND_MODULE,
    PROCESS_UNBIND_MODULE,
    PROCESS_SEARCH_MODULE,
    C_HOST,
    U_HOST,
    D_HOST
]

export const GET_AUTH_META = (auth, options = {}) => {
    const [ type, action ] = auth.split('.')
    const meta = {
        resource_type: type,
        action: action
    }
    // if (!options.noMode) {
    //     const business = store.getters['objectBiz/bizId']
    //     if (HAS_BUSINESS_MODE.includes(auth) && !store.getters.isAdminView && business) {
    //         meta.bk_biz_id = business
    //     }
    // }
    return meta
}

export const GET_MODEL_INST_AUTH_META = (modelId, auth, models = []) => {
    const meta = GET_AUTH_META(auth, { noMode: true })
    const model = models.find(model => model.bk_obj_id === modelId) || {}
    const modelLabel = (model.metadata || {}).label || {}
    if (modelLabel.hasOwnProperty('bk_biz_id')) {
        Object.assign(meta, {
            bk_biz_id: modelLabel.bk_biz_id
        })
    }
    return Object.assign(meta, {
        parent_layers: [{
            resource_type: 'model',
            resource_id: model.id
        }]
    })
}
