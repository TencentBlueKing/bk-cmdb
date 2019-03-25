// B = business 业务资源类型
// B_(C/R/U/D) 对应的CRUD权限

// 模型分组
export const B_C_MODEL_GROUP = 'modelGroup.create'
export const B_U_MODEL_GROUP = 'modelGroup.update'
export const B_D_MODEL_GROUP = 'modelGroup.delete'

// 模型
export const B_C_MODEL = 'model.create'
export const B_U_MODEL = 'model.update'
export const B_D_MODEL = 'xxxx'

// 实例
export const B_C_INST = 'xxxx'
export const B_U_INST = 'xxxx'
export const B_D_INST = 'xxxx'
export const B_R_INST = 'xxxx'

// 动态分组
export const B_C_CUSTOM_QUERY = 'xxxx'
export const B_U_CUSTOM_QUERY = 'xxxx'
export const B_D_CUSTOM_QUERY = 'xxxx'
export const B_R_CUSTOM_QUERY = 'xxxx'

// 进程管理
export const B_C_PROCESS = 'xxxx'
export const B_U_PROCESS = 'xxxx'
export const B_D_PROCESS = 'xxxx'
export const B_R_PROCESS = 'xxxx'
export const B_PROCESS_BIND_MODULE = 'xxxx'

// 业务拓扑
export const B_C_TOPO = 'xxxx'
export const B_U_TOPO = 'xxxx'
export const B_D_TOPO = 'xxxx'
export const B_R_TOPO = 'xxxx'
export const B_TOPO_TRANSFER_HOST = 'xxxx'

// 主机管理
export const B_U_HOST = 'xxxx'
export const B_R_HOST = 'xxxx'
export const B_HOST_TO_RESOURCE = 'xxxx'

// G = global 全局资源类型
// G_(C/R/U/D) 对应的CRUD权限

// 模型分组
export const G_C_MODEL_GROUP = 'xxxx'
export const G_U_MODEL_GROUP = 'xxxx'
export const G_D_MODEL_GROUP = 'xxxx'

// 模型
export const G_C_MODEL = 'xxxx'
export const G_R_MODEL = 'xxxx'
export const G_U_MODEL = 'xxxx'
export const G_D_MODEL = 'xxxx'

// 实例
export const G_C_INST = 'xxxx'
export const G_U_INST = 'xxxx'
export const G_D_INST = 'xxxx'
export const G_R_INST = 'xxxx'

// 关联类型
export const G_C_RELATION = 'xxxx'
export const G_R_RELATION = 'xxxx'
export const G_U_RELATION = 'xxxx'
export const G_D_RELATION = 'xxxx'

// 业务
export const G_C_BUSINESS = 'xxxx'
export const G_U_BUSINESS = 'xxxx'
export const G_D_BUSINESS = 'xxxx'
export const G_R_BUSINESS = 'xxxx'

// 主机
export const G_C_HOST = 'xxxx'
export const G_U_HOST = 'xxxx'
export const G_D_HOST = 'xxxx'
export const G_R_HOST = 'xxxx'
export const G_HOST_ASSIGN = 'xxxx'

// 事件推送
export const G_C_EVENT = 'xxxx'
export const G_U_EVENT = 'xxxx'
export const G_D_EVENT = 'xxxx'
export const G_R_EVENT = 'xxxx'

// 操作审计
export const G_R_AUDIT = 'xxxx'

// 系统基础
export const G_SYSTEM_TOPOLOGY = 'xxxx'
export const G_SYSTEM_MANAGEMENT = 'xxxx'
export const G_SYSTEM_MODEL_GRAPHICS = 'XXX'

// D 动态组合权限,例如具体某个模型的权限

// 模型
export const D_C_MODEL = model => {
    return `${model}.create`
}

export const D_R_MODEL = model => {
    return `${model}.read`
}

export const D_U_MODEL = model => {
    return `${model}.update`
}

export const D_D_MODEL = model => {
    return `${model}.delete`
}
