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

// 服务拓扑
export const C_TOPO = 'mainlineInstance.create'
export const U_TOPO = 'mainlineInstance.update'
export const D_TOPO = 'mainlineInstance.delete'
export const R_TOPO = 'mainlineObjectTopology.find'
export const TOPO_TRANSFER_HOST = 'mainlineInstanceTopology.transferHost'

// 业务主机
export const C_HOST = 'hostInstance.create'
export const R_HOST = 'hostInstance.findMany'
export const U_HOST = 'hostInstance.update'
export const D_HOST = 'hostInstance.delete'
export const HOST_TO_RESOURCE = 'hostInstance.moveHostFromModuleToResPool'

// 资源池主机
export const C_RESOURCE_HOST = 'hostInstance.create'
export const U_RESOURCE_HOST = 'hostInstance.update'
export const D_RESOURCE_HOST = 'hostInstance.delete'
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
export const SYSTEM_MODEL_GRAPHICS = 'systemBase.modelTopologyView'

// 云资源发现
export const C_CLOUD_DISCOVER = 'cloudDiscover.create'
export const U_CLOUD_DISCOVER = 'cloudDiscover.update'
export const D_CLOUD_DISCOVER = 'cloudDiscover.delete'
export const R_CLOUD_DISCOVER = 'cloudDiscover.findMany'

// 云资源确认
export const C_CLOUD_CONFIRM = 'cloudConfirm.create'
export const U_CLOUD_CONFIRM = 'cloudConfirm.update'
export const D_CLOUD_CONFIRM = 'cloudConfirm.delete'
export const R_CLOUD_CONFIRM = 'cloudConfirm.findMany'

// 确认历史
export const R_CONFIRM_HISTORY = 'cloudConfirmHistory.findMany'

// 服务分类
export const C_SERVICE_CATEGORY = 'processServiceCategory.create'
export const U_SERVICE_CATEGORY = 'processServiceCategory.update'
export const D_SERVICE_CATEGORY = 'processServiceCategory.delete'
export const R_SERVICE_CATEGORY = 'processServiceCategory.findMany'

// 服务模板
export const C_SERVICE_TEMPLATE = 'processServiceTemplate.create'
export const U_SERVICE_TEMPLATE = 'processServiceTemplate.update'
export const D_SERVICE_TEMPLATE = 'processServiceTemplate.delete'
export const R_SERVICE_TEMPLATE = 'processServiceTemplate.findMany'

// 服务实例
export const C_SERVICE_INSTANCE = 'processServiceInstance.create'
export const U_SERVICE_INSTANCE = 'processServiceInstance.update'
export const D_SERVICE_INSTANCE = 'processServiceInstance.delete'
export const R_SERVICE_INSTANCE = 'processServiceInstance.findMany'

export const STATIC_BUSINESS_MODE = [
    C_MODEL,
    R_MODEL,

    C_MODEL_GROUP,
    U_MODEL_GROUP,
    D_MODEL_GROUP,

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
    D_HOST,
    HOST_TO_RESOURCE,

    C_SERVICE_CATEGORY,
    U_SERVICE_CATEGORY,
    D_SERVICE_CATEGORY,
    R_SERVICE_CATEGORY,

    C_SERVICE_TEMPLATE,
    U_SERVICE_TEMPLATE,
    D_SERVICE_TEMPLATE,
    R_SERVICE_TEMPLATE,

    C_SERVICE_INSTANCE,
    U_SERVICE_INSTANCE,
    D_SERVICE_INSTANCE,
    R_SERVICE_INSTANCE,

    C_TOPO,
    U_TOPO,
    D_TOPO,
    R_TOPO
]

export const DYNAMIC_BUSINESS_MODE = [
    C_INST,
    U_INST,
    D_INST,
    R_INST
]

export const RESOURCE_TYPE_NAME = {
    modelClassification: '模型分类',
    model: '模型',
    modelInstance: '实例',
    dynamicGrouping: '动态分组',
    process: '进程',
    mainlineInstanceTopology: '业务拓扑',
    hostInstance: '主机',
    associationType: '关联类型',
    business: '业务',
    eventPushing: '事件推送',
    auditlog: '操作审计',
    systemBase: '系统基础',
    cloudDiscover: '云资源发现',
    cloudConfirm: '云资源确认',
    cloudConfirmHistory: '云资源确认历史',
    processServiceCategory: '服务分类',
    processServiceTemplate: '服务模板',
    processServiceInstance: '服务实例',
    mainlineInstance: '服务拓扑'
}

export const RESOURCE_ACTION_NAME = {
    create: '新建',
    update: '编辑',
    delete: '删除',
    findMany: '查询',
    boundModuleToProcess: '绑定到模块',
    unboundModelToProcess: '解绑模块',
    findBoundModuleProcess: '查询已绑定模块',
    transferHost: '转移主机',
    moveHostFromModuleToResPool: '删除/归还',
    moveResPoolHostToBizIdleModule: '分配主机到业务空闲机',
    archive: '归档',
    modelTopologyOperation: '拓扑层级管理',
    adminEntrance: '管理页面入口',
    modelTopologyView: '模型拓扑视图'
}

const AUTH_META_KEYS = ['bk_biz_id', 'parent_layers', 'resource_id']

export const GET_AUTH_META = (auth, options = {}) => {
    const [type, action, scope] = auth.split('.')
    const meta = {
        scope: scope || 'global',
        resource_type: type,
        action: action
    }
    Object.keys(options).forEach(key => {
        if (AUTH_META_KEYS.includes(key)) {
            meta[key] = options[key]
        }
    })
    return meta
}
