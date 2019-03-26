// B = business 业务资源类型
// B_(C/R/U/D) 对应的CRUD权限

// 模型分组
export const B_C_MODEL_GROUP = 'bizModelGroup.create'
export const B_U_MODEL_GROUP = 'bizModelGroup.edit'
export const B_D_MODEL_GROUP = 'bizModelGroup.delete'

// 模型
export const B_C_MODEL = 'bizModel.create'
export const B_U_MODEL = 'bizModel.edit'
export const B_D_MODEL = 'bizModel.delete'

// 实例
export const B_C_INST = 'bizInstance.create'
export const B_U_INST = 'bizInstance.edit'
export const B_D_INST = 'bizInstance.delete'
export const B_R_INST = 'bizInstance.get'

// 动态分组
export const B_C_CUSTOM_QUERY = 'dynamicGrouping.create'
export const B_U_CUSTOM_QUERY = 'dynamicGrouping.update'
export const B_D_CUSTOM_QUERY = 'dynamicGrouping.delete'
export const B_R_CUSTOM_QUERY = 'dynamicGrouping.findMany'

// 进程管理
export const B_C_PROCESS = 'bizProcessInstance.create'
export const B_U_PROCESS = 'bizProcessInstance.edit'
export const B_D_PROCESS = 'bizProcessInstance.delete'
export const B_R_PROCESS = 'bizProcessInstance.get'
export const B_PROCESS_BIND_MODULE = 'bizProcessInstance.bindModule'
export const B_PROCESS_SEARCH_MODULE = 'bizProcessInstance.bindModuleQuery'

// 业务拓扑
export const B_C_TOPO = 'bizTopoInstance.create'
export const B_U_TOPO = 'bizTopoInstance.edit'
export const B_D_TOPO = 'bizTopoInstance.delete'
export const B_R_TOPO = 'bizTopoInstance.get'
export const B_TOPO_TRANSFER_HOST = 'bizTopoInstance.hostTransfer'

// 主机管理
export const B_U_HOST = 'bizHostInstance.edit'
export const B_R_HOST = 'bizHostInstance.get'
export const B_HOST_TO_RESOURCE = 'bizHostInstance.moduleTransfer'

// G = global 全局资源类型
// G_(C/R/U/D) 对应的CRUD权限

// 模型分组
export const G_C_MODEL_GROUP = 'sysModelGroup.create'
export const G_U_MODEL_GROUP = 'sysModelGroup.edit'
export const G_D_MODEL_GROUP = 'sysModelGroup.delete'

// 模型
export const G_C_MODEL = 'sysModel.create'
export const G_U_MODEL = 'sysModel.edit'
export const G_D_MODEL = 'sysModel.delete'

// 实例
export const G_C_INST = 'sysInstance.create'
export const G_U_INST = 'sysInstance.edit'
export const G_D_INST = 'sysInstance.delete'
export const G_R_INST = 'sysInstance.get'

// 关联类型
export const G_C_RELATION = 'sysAssociationType.create'
export const G_U_RELATION = 'sysAssociationType.edit'
export const G_D_RELATION = 'sysAssociationType.delete'

// 业务
export const G_C_BUSINESS = 'business.create'
export const G_U_BUSINESS = 'business.update'
export const G_D_BUSINESS = 'business.delete'
export const G_R_BUSINESS = 'business.findMany'

// 主机
export const G_C_HOST = 'sysHostInstance.create'
export const G_U_HOST = 'sysHostInstance.edit'
export const G_D_HOST = 'sysHostInstance.delete'
export const G_R_HOST = 'sysHostInstance.get'
export const G_HOST_ASSIGN = 'sysHostInstance.moduleTransfer'

// 事件推送
export const G_C_EVENT = 'sysEventPushing.create'
export const G_U_EVENT = 'sysEventPushing.edit'
export const G_D_EVENT = 'sysEventPushing.delete'
export const G_R_EVENT = 'sysEventPushing.get'

// 操作审计
export const G_R_AUDIT = 'xxxx'

// 系统基础
export const G_SYSTEM_TOPOLOGY = 'mainlineObjectTopology.find'
export const G_SYSTEM_MANAGEMENT = 'sysSystemBase.adminEntrance'
export const G_SYSTEM_MODEL_GRAPHICS = 'sysSystemBase.modelTopologyView'

// D 动态组合权限,例如具体某个模型的权限

// 模型
export const D_C_MODEL = model => {
    return `${model}.create`
}

export const D_R_MODEL = model => {
    return `${model}.get`
}

export const D_U_MODEL = model => {
    return `${model}.edit`
}

export const D_D_MODEL = model => {
    return `${model}.delete`
}
