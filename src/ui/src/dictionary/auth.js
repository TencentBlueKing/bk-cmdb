import i18n from '@/i18n'

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

// 统计报表
export const C_STATISTICAL_REPORT = 'operationStatistic.create'
export const U_STATISTICAL_REPORT = 'operationStatistic.update'
export const D_STATISTICAL_REPORT = 'operationStatistic.delete'
export const R_STATISTICAL_REPORT = 'operationStatistic.findMany'

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

// 集群模板
export const C_SET_TEMPLATE = 'setTemplate.create'
export const U_SET_TEMPLATE = 'setTemplate.update'
export const D_SET_TEMPLATE = 'setTemplate.delete'

export const RESOURCE_TYPE_NAME = {
    modelClassification: i18n.t('模型分类'),
    model: i18n.t('模型'),
    modelInstance: i18n.t('实例'),
    dynamicGrouping: i18n.t('动态分组'),
    process: i18n.t('进程'),
    mainlineInstanceTopology: i18n.t('业务拓扑'),
    hostInstance: i18n.t('主机'),
    associationType: i18n.t('关联类型'),
    business: i18n.t('业务'),
    eventPushing: i18n.t('事件推送'),
    auditlog: i18n.t('操作审计'),
    systemBase: i18n.t('系统基础'),
    cloudDiscover: i18n.t('云资源发现'),
    cloudConfirm: i18n.t('云资源确认'),
    cloudConfirmHistory: i18n.t('云资源确认历史'),
    processServiceCategory: i18n.t('服务分类'),
    processServiceTemplate: i18n.t('服务模板'),
    processServiceInstance: i18n.t('服务实例'),
    mainlineInstance: i18n.t('服务拓扑'),
    operationStatistic: i18n.t('运营统计'),
    setTemplate: i18n.t('集群模板')
}

export const RESOURCE_ACTION_NAME = {
    create: i18n.t('新建'),
    update: i18n.t('编辑'),
    delete: i18n.t('删除'),
    findMany: i18n.t('查询'),
    boundModuleToProcess: i18n.t('绑定到模块'),
    unboundModelToProcess: i18n.t('解绑模块'),
    findBoundModuleProcess: i18n.t('查询已绑定模块'),
    transferHost: i18n.t('转移主机'),
    moveHostFromModuleToResPool: i18n.t('删除/归还'),
    moveResPoolHostToBizIdleModule: i18n.t('分配主机到业务空闲机'),
    archive: i18n.t('归档'),
    modelTopologyOperation: i18n.t('拓扑层级管理'),
    adminEntrance: i18n.t('管理页面入口'),
    modelTopologyView: i18n.t('模型拓扑视图')
}

const AUTH_META_KEYS = ['bk_biz_id', 'parent_layers', 'resource_id']

export const GET_AUTH_META = (auth, options = {}) => {
    const [type, action] = auth.split('.')
    const meta = {
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
