// B = business 业务资源类型
// B_(C/R/U/D) 对应的CRUD权限

// 模型分组
export const B_C_MODEL_GROUP = 'modelClassification.create'
export const B_U_MODEL_GROUP = 'modelClassification.edit'
export const B_D_MODEL_GROUP = 'modelClassification.delete'

// 模型
export const B_C_MODEL = 'model.create'
export const B_U_MODEL = 'model.edit'
export const B_D_MODEL = 'model.delete'

// 实例
export const B_C_INST = 'modelInstance.create'
export const B_U_INST = 'modelInstance.edit'
export const B_D_INST = 'modelInstance.delete'
export const B_R_INST = 'modelInstance.findMany'

// 动态分组
export const B_C_CUSTOM_QUERY = 'dynamicGrouping.create'
export const B_U_CUSTOM_QUERY = 'dynamicGrouping.update'
export const B_D_CUSTOM_QUERY = 'dynamicGrouping.delete'
export const B_R_CUSTOM_QUERY = 'dynamicGrouping.findMany'

// 进程管理
export const B_C_PROCESS = 'process.create'
export const B_U_PROCESS = 'process.edit'
export const B_D_PROCESS = 'process.delete'
export const B_R_PROCESS = 'process.findMany'
export const B_PROCESS_BIND_MODULE = 'process.boundModuleToProcess'
export const B_PROCESS_UNBIND_MODULE = 'process.unboundModelToProcess'
export const B_PROCESS_SEARCH_MODULE = 'process.findBoundModuleProcess'

// 业务拓扑
export const B_C_TOPO = 'mainlineInstanceTopology.create'
export const B_U_TOPO = 'mainlineInstanceTopology.edit'
export const B_D_TOPO = 'mainlineInstanceTopology.delete'
export const B_R_TOPO = 'mainlineInstanceTopology.findMany'
export const B_TOPO_TRANSFER_HOST = 'mainlineInstanceTopology.transferHost'

// 主机管理
export const B_U_HOST = 'hostInstance.edit'
export const B_R_HOST = 'hostInstance.findMany'
export const B_HOST_TO_RESOURCE = 'hostInstance.moduleTransfer'

// G = global 全局资源类型
// G_(C/R/U/D) 对应的CRUD权限

// 模型分组
export const G_C_MODEL_GROUP = 'modelClassification.create'
export const G_U_MODEL_GROUP = 'modelClassification.edit'
export const G_D_MODEL_GROUP = 'modelClassification.delete'

// 模型
export const G_C_MODEL = 'model.create'
export const G_U_MODEL = 'model.edit'
export const G_D_MODEL = 'model.delete'

// 实例
export const G_C_INST = 'modelInstance.create'
export const G_U_INST = 'modelInstance.edit'
export const G_D_INST = 'modelInstance.delete'
export const G_R_INST = 'modelInstance.findMany'

// 关联类型
export const G_C_RELATION = 'associationType.create'
export const G_U_RELATION = 'associationType.edit'
export const G_D_RELATION = 'associationType.delete'

// 业务
export const G_C_BUSINESS = 'business.create'
export const G_U_BUSINESS = 'business.update'
export const G_R_BUSINESS = 'business.findMany'
export const G_BUSINESS_ARCHIVE = 'business.archive'

// 主机
export const G_C_HOST = 'hostInstance.create'
export const G_U_HOST = 'hostInstance.edit'
export const G_D_HOST = 'hostInstance.delete'
export const G_R_HOST = 'hostInstance.findMany'
export const G_HOST_ASSIGN = 'hostInstance.moduleTransfer'

// 事件推送
export const G_C_EVENT = 'eventPushing.create'
export const G_U_EVENT = 'eventPushing.update'
export const G_D_EVENT = 'eventPushing.delete'
export const G_R_EVENT = 'eventPushing.findMany'

// 操作审计
export const G_R_AUDIT = 'auditlog.findMany'

// 系统基础
export const G_SYSTEM_TOPOLOGY = 'mainlineObjectTopology.find'
export const G_SYSTEM_MANAGEMENT = 'sysSystemBase.adminEntrance'
export const G_SYSTEM_MODEL_GRAPHICS = 'modelTopology.findMany'

// D 动态组合权限,例如具体某个模型的权限

// 模型
export const D_C_MODEL = model => {
    return `modelInstance.create`
}

export const D_R_MODEL = model => {
    return `modelInstance.findMany`
}

export const D_U_MODEL = model => {
    return `modelInstance.edit`
}

export const D_D_MODEL = model => {
    return `modelInstance.delete`
}
