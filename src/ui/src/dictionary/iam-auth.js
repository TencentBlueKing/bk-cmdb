export const IAM_VIEWS = {
    // 模型分组
    MODEL_GROUP: 'sys_model_group',
    // 模型列表
    MODEL: 'sys_model',
    // 通用模型实例列表
    INSTANCE: 'sys_instance',
    // 通用模型列表
    INSTANCE_MODEL: 'sys_instance_model',
    // 动态分组列表
    CUSTOM_QUERY: 'biz_custom_query',
    // 业务列表
    BIZ: 'biz',
    // 跨业务转主机选择的主机所属业务的列表
    BIZ_FOR_HOST_TRANS: 'biz_for_host_trans',
    // 主机列表
    HOST: 'host',
    // 主机池目录列表(作为源目录时使用的视图)
    RESOURCE_SOURCE_POOL_DIRECTORY: 'sys_host_rsc_pool_directory',
    // 主机池目录列表(作为目标目录时使用的视图)
    RESOURCE_TARGET_POOL_DIRECTORY: 'sys_resource_pool_directory',
    // 关联类型列表
    ASSOCIATION_TYPE: 'sys_association_type',
    // 事件订阅列表
    EVENT_PUSHING: 'sys_event_pushing',
    // 服务模板列表
    SERVICE_TEMPLATE: 'biz_process_service_template',
    // 集群模板列表
    SET_TEMPLATE: 'biz_set_template',
    // 云区域列表
    CLOUD_AREA: 'sys_cloud_area',
    // 云账户列表
    CLOUD_ACCOUNT: 'sys_cloud_account',
    // 云发现任务
    CLOUD_RESOURCE_TASK: 'sys_cloud_resource_task'
}

export const IAM_VIEWS_NAME = {
    [IAM_VIEWS.MODEL_GROUP]: ['模型分组', 'Model Group'],
    [IAM_VIEWS.MODEL]: ['模型', 'Model'],
    [IAM_VIEWS.INSTANCE]: ['实例', 'Instance'],
    [IAM_VIEWS.INSTANCE_MODEL]: ['实例模型', 'Instance Model'],
    [IAM_VIEWS.CUSTOM_QUERY]: ['动态分组', 'Custom Query'],
    [IAM_VIEWS.BIZ]: ['业务', 'Business'],
    [IAM_VIEWS.BIZ_FOR_HOST_TRANS]: ['业务', 'Business'],
    [IAM_VIEWS.HOST]: ['主机', 'Host'],
    [IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY]: ['主机池目录', 'Resource Pool Directory'],
    [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY]: ['主机池目录', 'Resource Pool Directory'],
    [IAM_VIEWS.ASSOCIATION_TYPE]: ['关联类型', 'Association Type'],
    [IAM_VIEWS.EVENT_PUSHING]: ['事件订阅', 'Event Pushing'],
    [IAM_VIEWS.SERVICE_TEMPLATE]: ['服务模板', 'Service Template'],
    [IAM_VIEWS.SET_TEMPLATE]: ['集群模板', 'Set Template'],
    [IAM_VIEWS.CLOUD_AREA]: ['云区域', 'Cloud Area'],
    [IAM_VIEWS.CLOUD_ACCOUNT]: ['云账户', 'Cloud Account'],
    [IAM_VIEWS.CLOUD_RESOURCE_TASK]: ['云资源发现任务', 'Cloud Resource Task']
}

function basicTransform (cmdbAction, meta = {}) {
    const [internalType, internalAction] = cmdbAction.split('.')
    const inejctedMeta = {
        resource_type: internalType,
        action: internalAction,
        ...meta
    }
    Object.keys(inejctedMeta).forEach(key => {
        const value = inejctedMeta[key]
        if (value === null || value === undefined) {
            delete inejctedMeta[key]
        }
    })
    return inejctedMeta
}

// relation数组表示的是视图拓扑的定义
export const IAM_ACTIONS = {
    // 模型分组
    C_MODEL_GROUP: {
        id: 'create_model_group',
        name: ['模型分组创建', 'Create Model Group'],
        cmdb_action: 'modelClassification.create'
    },
    U_MODEL_GROUP: {
        id: 'edit_model_group',
        name: ['模型分组编辑', 'Update Model Group'],
        cmdb_action: 'modelClassification.update',
        relation: [{
            view: IAM_VIEWS.MODEL_GROUP,
            instances: [IAM_VIEWS.MODEL_GROUP]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    D_MODEL_GROUP: {
        id: 'delete_model_group',
        name: ['模型分组删除', 'Delete Model Group'],
        cmdb_action: 'modelClassification.delete',
        relation: [{
            view: IAM_VIEWS.MODEL_GROUP,
            instances: [IAM_VIEWS.MODEL_GROUP]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },

    // 模型
    C_MODEL: {
        id: 'create_sys_model',
        name: ['模型创建', 'Create Model'],
        cmdb_action: 'model.create',
        relation: [{
            view: IAM_VIEWS.MODEL_GROUP,
            instances: [IAM_VIEWS.MODEL_GROUP]
        }],
        transform: (cmdbAction, relationIds) => {
            const [modelGroupId] = relationIds
            const verifyMeta = basicTransform(cmdbAction, {})
            if (modelGroupId) {
                verifyMeta.parent_layers = [{
                    resource_type: 'modelClassification',
                    resource_id: modelGroupId
                }]
            }
            return verifyMeta
        }
    },
    U_MODEL: {
        id: 'edit_sys_model',
        name: ['模型编辑', 'Update Model'],
        cmdb_action: 'model.update',
        relation: [{
            view: IAM_VIEWS.MODEL,
            instances: [IAM_VIEWS.MODEL]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    D_MODEL: {
        id: 'delete_sys_model',
        name: ['模型删除', 'Delete Model'],
        cmdb_action: 'model.delete',
        relation: [{
            view: IAM_VIEWS.MODEL,
            instances: [IAM_VIEWS.MODEL]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },

    // 实例
    C_INST: {
        id: 'create_sys_instance',
        name: ['实例创建', 'Create Instance'],
        cmdb_action: 'modelInstance.create',
        relation: [{
            view: IAM_VIEWS.INSTANCE_MODEL,
            instances: [IAM_VIEWS.INSTANCE_MODEL]
        }],
        transform: (cmdbAction, relationIds = []) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            if (relationIds.length) {
                const [modelId] = relationIds
                verifyMeta.parent_layers = [{
                    resource_type: 'model',
                    resource_id: modelId
                }]
            }
            return verifyMeta
        }
    },
    U_INST: {
        id: 'edit_sys_instance',
        name: ['实例编辑', 'Update Instance'],
        cmdb_action: 'modelInstance.update',
        relation: [{
            view: IAM_VIEWS.INSTANCE,
            instances: [IAM_VIEWS.INSTANCE_MODEL, IAM_VIEWS.INSTANCE]
        }],
        transform: (cmdbAction, relationIds = []) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            if (relationIds.length) {
                const [modelId, instanceId] = relationIds
                verifyMeta.parent_layers = [{
                    resource_type: 'model',
                    resource_id: modelId
                }]
                if (instanceId) {
                    verifyMeta.resource_id = instanceId
                }
            }
            return verifyMeta
        }
    },
    D_INST: {
        id: 'delete_sys_instance',
        name: ['实例删除', 'Delete Instance'],
        cmdb_action: 'modelInstance.delete',
        relation: [{
            view: IAM_VIEWS.INSTANCE,
            instances: [IAM_VIEWS.INSTANCE_MODEL, IAM_VIEWS.INSTANCE]
        }],
        transform: (cmdbAction, relationIds = []) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            if (relationIds.length) {
                const [modelId, instanceId] = relationIds
                verifyMeta.parent_layers = [{
                    resource_type: 'model',
                    resource_id: modelId
                }]
                if (instanceId) {
                    verifyMeta.resource_id = instanceId
                }
            }
            return verifyMeta
        }
    },
    R_INST: {
        id: 'find_sys_instance',
        name: ['实例查询', 'Search Instance'],
        cmdb_action: 'modelInstance.findMany',
        relation: [{
            view: IAM_VIEWS.INSTANCE,
            instances: [IAM_VIEWS.INSTANCE_MODEL, IAM_VIEWS.INSTANCE]
        }],
        transform: (cmdbAction, relationIds = []) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            if (relationIds.length) {
                const [modelId, instanceId] = relationIds
                verifyMeta.parent_layers = [{
                    resource_type: 'model',
                    resource_id: modelId
                }]
                if (instanceId) {
                    verifyMeta.resource_id = instanceId
                }
            }
            return verifyMeta
        }
    },

    // 动态分组
    C_CUSTOM_QUERY: {
        id: 'create_biz_dynamic_query',
        name: ['动态分组创建', 'Create Custom Query'],
        cmdb_action: 'dynamicGrouping.create',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },
    U_CUSTOM_QUERY: {
        id: 'edit_biz_dynamic_query',
        name: ['动态分组编辑', 'Update Custom Query'],
        cmdb_action: 'dynamicGrouping.update',
        relation: [{
            view: IAM_VIEWS.CUSTOM_QUERY,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.CUSTOM_QUERY]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, customQueryId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id_ex: customQueryId // resource_id需要int，用resource_id_ex传string
            })
        }
    },
    D_CUSTOM_QUERY: {
        id: 'delete_biz_dynamic_query',
        name: ['动态分组删除', 'Delete Custom Query'],
        cmdb_action: 'dynamicGrouping.delete',
        relation: [{
            view: IAM_VIEWS.CUSTOM_QUERY,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.CUSTOM_QUERY]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, customQueryId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id_ex: customQueryId // resource_id需要int，用resource_id_ex传string
            })
        }
    },
    R_CUSTOM_QUERY: {
        id: 'find_biz_dynamic_query',
        name: ['动态分组查询', 'Search Custom Query'],
        cmdb_action: 'dynamicGrouping.findMany',
        relation: [{
            view: IAM_VIEWS.CUSTOM_QUERY,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.CUSTOM_QUERY]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, customQueryId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id_ex: customQueryId // resource_id需要int，用resource_id_ex传string
            })
        }
    },

    // 业务拓扑
    C_TOPO: {
        id: 'create_biz_topology',
        name: ['业务拓扑新建', 'Create Business Topology'],
        cmdb_action: 'mainlineInstance.create',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },
    U_TOPO: {
        id: 'edit_biz_topology',
        name: ['业务拓扑编辑', 'Update Business Topology'],
        cmdb_action: 'mainlineInstance.update',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },
    D_TOPO: {
        id: 'delete_biz_topology',
        name: ['业务拓扑删除', 'Delete Business Topology'],
        cmdb_action: 'mainlineInstance.delete',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },

    U_HOST: {
        id: 'edit_biz_host',
        name: ['业务主机编辑', 'Update Business Host'],
        cmdb_action: 'hostInstance.update',
        relation: [{
            view: IAM_VIEWS.HOST,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.HOST]
        }],
        transform: (cmdbAction, relationIds) => {
            const isBatch = Array.isArray(relationIds[0])
            if (isBatch) { // 批量编辑的场景
                const metas = relationIds.map(([bizId, hostId]) => {
                    const verifyMeta = basicTransform(cmdbAction, {
                        bk_biz_id: bizId,
                        parent_layers: [{
                            resource_type: 'biz',
                            resource_id: bizId
                        }]
                    })
                    if (hostId) {
                        verifyMeta.resource_id = hostId
                    }
                    return verifyMeta
                })
                return metas
            } else { // 单个编辑的场景
                const [bizId, hostId] = relationIds
                const verifyMeta = basicTransform(cmdbAction, {
                    bk_biz_id: bizId,
                    parent_layers: [{
                        resource_type: 'biz',
                        resource_id: bizId
                    }]
                })
                if (hostId) {
                    verifyMeta.resource_id = hostId
                }
                return verifyMeta
            }
        }
    },
    HOST_TO_RESOURCE: {
        id: 'unassign_biz_host',
        name: ['主机归还主机池', 'Transfer Host To Resource Pool'],
        cmdb_action: 'hostInstance.moveHostFromModuleToResPool',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }, {
            view: IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY]
        }],
        transform: (cmdbAction, relationIds) => {
            const [[[bizId], [directoryId]]] = relationIds
            const verifyMeta = basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
            verifyMeta.parent_layers = [{
                resource_type: 'biz',
                resource_id: bizId
            }, {
                resource_type: 'ResourcePoolDirectory',
                resource_id: directoryId
            }]
            return verifyMeta
        }
    },
    // 跨业务转主机
    HOST_TRANSFER_ACROSS_BIZ: {
        id: 'host_transfer_across_business',
        name: ['主机转移到其他业务', 'Transfer Host To Other Business'],
        cmdb_action: 'hostInstance.moveHostToAnotherBizModule',
        relation: [{
            view: IAM_VIEWS.BIZ_FOR_HOST_TRANS,
            instances: [IAM_VIEWS.BIZ_FOR_HOST_TRANS]
        }, {
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [[[currentBizId], [targetBizId]]] = relationIds
            const verifyMeta = basicTransform(cmdbAction)
            verifyMeta.parent_layers = [{
                resource_id: currentBizId,
                resource_type: 'biz'
            }, {
                resource_id: targetBizId,
                resource_type: 'biz'
            }]
            return verifyMeta
        }
    },

    // 主机池主机
    C_RESOURCE_HOST: {
        id: 'create_resource_pool_host',
        name: ['主机池主机创建', 'Create Resource Pool Host'],
        cmdb_action: 'hostInstance.create',
        relation: [{
            view: IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY]
        }],
        transform: (cmdbAction, relationIds = []) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            const [directoryId = 1] = relationIds
            if (directoryId) {
                verifyMeta.parent_layers = [{
                    resource_type: 'resourcePoolDirectory',
                    resource_id: directoryId
                }]
            }
            return verifyMeta
        }
    },
    U_RESOURCE_HOST: {
        id: 'edit_resource_pool_host',
        name: ['主机池主机编辑', 'Update Resource Pool Host'],
        cmdb_action: 'hostInstance.update',
        relation: [{
            view: IAM_VIEWS.HOST,
            instances: [IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY, IAM_VIEWS.HOST]
        }],
        transform: (cmdbAction, relationIds) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            const [directoryId, hostId] = relationIds
            if (hostId) {
                verifyMeta.resource_id = hostId
            }
            if (directoryId) {
                verifyMeta.parent_layers = [{
                    resource_type: 'resourcePoolDirectory',
                    resource_id: directoryId
                }]
            }
            return verifyMeta
        }
    },
    D_RESOURCE_HOST: {
        id: 'delete_resource_pool_host',
        name: ['主机池主机删除', 'Delete Resource Pool Host'],
        cmdb_action: 'hostInstance.delete',
        relation: [{
            view: IAM_VIEWS.HOST,
            instances: [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY, IAM_VIEWS.HOST]
        }],
        transform: (cmdbAction, relationIds) => {
            const verifyMeta = basicTransform(cmdbAction, {})
            const [directoryId, hostId] = relationIds
            if (hostId) {
                verifyMeta.resource_id = hostId
            }
            if (directoryId) {
                verifyMeta.parent_layers = [{
                    resource_type: 'resourcePoolDirectory',
                    resource_id: directoryId
                }]
            }
            return verifyMeta
        }
    },
    TRANSFER_HOST_TO_BIZ: {
        id: 'assign_host_to_biz',
        name: ['主机池主机分配到业务', 'Transfer Resource Pool Host To Business'],
        cmdb_action: 'hostInstance.moveResPoolHostToBizIdleModule',
        relation: [{
            view: IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY]
        }, {
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [[hostRelationIds, bizRelationIds]] = relationIds
            const [directoryId] = hostRelationIds
            const [bizId] = bizRelationIds
            const hostVerifyMeta = basicTransform(cmdbAction, {
                parent_layers: [{
                    resource_type: 'resourcePoolDirectory',
                    resource_id: directoryId
                }, {
                    resource_type: 'business',
                    resource_id: bizId
                }]
            })
            return hostVerifyMeta
        }
    },
    TRANSFER_HOST_TO_DIRECTORY: {
        id: 'host_transfer_in_resource_pool',
        name: ['主机池主机分配到目录', 'Change Resource Pool Host\'s Directory'],
        cmdb_action: 'hostInstance.moveResPoolHostToDirectory',
        relation: [{
            view: IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_SOURCE_POOL_DIRECTORY]
        }, {
            view: IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY]
        }],
        transform: (cmdbAction, relationIds) => {
            const [[[currentDirectoryId], [targetDirectoryId]]] = relationIds
            const hostVerifyMeta = basicTransform(cmdbAction, {
                parent_layers: [{
                    resource_type: 'resourcePoolDirectory',
                    resource_id: currentDirectoryId
                }, {
                    resource_type: 'resourcePoolDirectory',
                    resource_id: targetDirectoryId
                }]
            })
            return hostVerifyMeta
        }
    },

    // 主机池目录
    C_RESOURCE_DIRECTORY: {
        id: 'create_resource_pool_directory',
        name: ['主机池目录创建', 'Create Resource Pool Directory'],
        cmdb_action: 'resourcePoolDirectory.create'
    },
    U_RESOURCE_DIRECTORY: {
        id: 'edit_resource_pool_directory',
        name: ['主机池目录编辑', 'Update Resource Pool Directory'],
        cmdb_action: 'resourcePoolDirectory.update',
        relation: [{
            view: IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    D_RESOURCE_DIRECTORY: {
        id: 'delete_resource_pool_directory',
        name: ['主机池目录删除', 'Delete Resource Pool Directory'],
        cmdb_action: 'resourcePoolDirectory.delete',
        relation: [{
            view: IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY,
            instances: [IAM_VIEWS.RESOURCE_TARGET_POOL_DIRECTORY]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },

    // 关联类型
    C_RELATION: {
        id: 'create_association_type',
        name: ['关联类型创建', 'Create Association Type'],
        cmdb_action: 'associationType.create'
    },
    U_RELATION: {
        id: 'edit_association_type',
        name: ['关联类型编辑', 'Update Association Type'],
        cmdb_action: 'associationType.update',
        relation: [{
            view: IAM_VIEWS.ASSOCIATION_TYPE,
            instances: [IAM_VIEWS.ASSOCIATION_TYPE]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    D_RELATION: {
        id: 'delete_association_type',
        name: ['关联类型删除', 'Delete Association Type'],
        cmdb_action: 'associationType.delete',
        relation: [{
            view: IAM_VIEWS.ASSOCIATION_TYPE,
            instances: [IAM_VIEWS.ASSOCIATION_TYPE]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },

    // 业务
    C_BUSINESS: {
        id: 'create_business',
        name: ['业务创建', 'Create Business'],
        cmdb_action: 'business.create'
    },
    U_BUSINESS: {
        id: 'edit_business',
        name: ['业务编辑', 'Update Business'],
        cmdb_action: 'business.update',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    R_BUSINESS: {
        id: 'find_business',
        name: ['业务查询', 'Search Business'],
        cmdb_action: 'business.findMany',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    BUSINESS_ARCHIVE: {
        id: 'archive_business',
        name: ['业务归档', 'Business Archive'],
        cmdb_action: 'business.archive',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },

    // 事件推送
    C_EVENT: {
        id: 'create_event_subscription',
        name: ['事件订阅创建', 'Create Event Subscription'],
        cmdb_action: 'eventPushing.create'
    },
    U_EVENT: {
        id: 'edit_event_subscription',
        name: ['事件订阅编辑', 'Update Event Subscription'],
        cmdb_action: 'eventPushing.update',
        relation: [{
            view: IAM_VIEWS.EVENT_PUSHING,
            instances: [IAM_VIEWS.EVENT_PUSHING]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    D_EVENT: {
        id: 'delete_event_subscription',
        name: ['事件订阅删除', 'Delete Event Subscription'],
        cmdb_action: 'eventPushing.delete',
        relation: [{
            view: IAM_VIEWS.EVENT_PUSHING,
            instances: [IAM_VIEWS.EVENT_PUSHING]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },
    R_EVENT: {
        id: 'find_event_subscription',
        name: ['事件订阅查询', 'Search Event Subscription'],
        cmdb_action: 'eventPushing.findMany',
        relation: [{
            view: IAM_VIEWS.EVENT_PUSHING,
            instances: [IAM_VIEWS.EVENT_PUSHING]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                resource_id: relationIds[0]
            })
        }
    },

    // 操作审计
    R_AUDIT: {
        id: 'find_audit_log',
        name: ['操作审计查询', 'Search Audit Logs'],
        cmdb_action: 'auditlog.findMany'
    },

    // 拓扑层级新增
    SYSTEM_TOPOLOGY: {
        id: 'edit_business_layer',
        name: ['业务层级编辑', 'Update Business Topology Layer'],
        cmdb_action: 'systemBase.modelTopologyOperation'
    },

    // 拓扑模型关系图
    SYSTEM_MODEL_GRAPHICS: {
        id: 'edit_model_topology_view',
        name: ['模型拓扑视图编辑', 'Update Model Topology View'],
        cmdb_action: 'systemBase.modelTopologyView'
    },

    // 统计报表
    U_STATISTICAL_REPORT: {
        id: 'edit_operation_statistic',
        name: ['运营统计编辑', 'Update Operation Statistic'],
        cmdb_action: 'operationStatistic.update'
    },
    R_STATISTICAL_REPORT: {
        id: 'find_operation_statistic',
        name: ['运营统计查询', 'Search Operation Statistic'],
        cmdb_action: 'operationStatistic.findMany'
    },

    // 服务分类
    C_SERVICE_CATEGORY: {
        id: 'create_biz_service_category',
        name: ['服务分类新建', 'Create Service Category'],
        cmdb_action: 'processServiceCategory.create',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },
    U_SERVICE_CATEGORY: {
        id: 'edit_biz_service_category',
        name: ['服务分类编辑', 'Update Service Category'],
        cmdb_action: 'processServiceCategory.update',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },
    D_SERVICE_CATEGORY: {
        id: 'delete_biz_service_category',
        name: ['服务分类删除', 'Delete Service Cateogry'],
        cmdb_action: 'processServiceCategory.delete',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },

    // 服务模板
    C_SERVICE_TEMPLATE: {
        id: 'create_biz_service_template',
        name: ['服务模板创建', 'Create Service Template'],
        cmdb_action: 'processServiceTemplate.create',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    },
    U_SERVICE_TEMPLATE: {
        id: 'edit_biz_service_template',
        name: ['服务模板编辑', 'Update Service Template'],
        cmdb_action: 'processServiceTemplate.update',
        relation: [{
            view: IAM_VIEWS.SERVICE_TEMPLATE,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.SERVICE_TEMPLATE]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, templateId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id: templateId
            })
        }
    },
    D_SERVICE_TEMPLATE: {
        id: 'delete_biz_service_template',
        name: ['服务模板删除', 'Delete Service Template'],
        cmdb_action: 'processServiceTemplate.delete',
        relation: [{
            view: IAM_VIEWS.SERVICE_TEMPLATE,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.SERVICE_TEMPLATE]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, templateId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id: templateId
            })
        }
    },

    // 服务实例
    C_SERVICE_INSTANCE: {
        id: 'create_biz_service_instance',
        name: ['服务实例创建', 'Create Service Instance'],
        cmdb_action: 'processServiceInstance.create',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
        }
    },
    U_SERVICE_INSTANCE: {
        id: 'edit_biz_service_instance',
        name: ['服务实例编辑', 'Update Service Instance'],
        cmdb_action: 'processServiceInstance.update',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
        }
    },
    D_SERVICE_INSTANCE: {
        id: 'delete_biz_service_instance',
        name: ['服务实例删除', 'Delete Service Instance'],
        cmdb_action: 'processServiceInstance.delete',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
        }
    },

    // 集群模板
    C_SET_TEMPLATE: {
        id: 'create_biz_set_template',
        name: ['集群模板创建', 'Create Set Template'],
        cmdb_action: 'setTemplate.create',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
        }
    },
    U_SET_TEMPLATE: {
        id: 'edit_biz_set_template',
        name: ['集群模板编辑', 'Update Set Template'],
        cmdb_action: 'setTemplate.update',
        relation: [{
            view: IAM_VIEWS.SET_TEMPLATE,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.SET_TEMPLATE]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, templateId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id: templateId
            })
        }
    },
    D_SET_TEMPLATE: {
        id: 'delete_biz_set_template',
        name: ['集群模板删除', 'Delete Set Template'],
        cmdb_action: 'setTemplate.delete',
        relation: [{
            view: IAM_VIEWS.SET_TEMPLATE,
            instances: [IAM_VIEWS.BIZ, IAM_VIEWS.SET_TEMPLATE]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId, templateId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId,
                resource_id: templateId
            })
        }
    },

    // 主机属性自动应用
    U_HOST_APPLY: {
        id: 'edit_biz_host_apply',
        name: ['属性自动应用编辑', 'Update Host Apply'],
        cmdb_action: 'hostApply.update',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
        }
    },

    // 管理员
    R_CONFIG_ADMIN: {
        id: 'global_settings',
        name: ['全局设置', 'Global Settings'],
        cmdb_action: 'configAdmin.find'
    },
    U_CONFIG_ADMIN: {
        id: 'global_settings',
        name: ['全局设置', 'Global Settings'],
        cmdb_action: 'configAdmin.update'
    },

    // 云区域
    C_CLOUD_AREA: {
        id: 'create_cloud_area',
        name: ['云区域创建', 'Create Cloud Area'],
        cmdb_action: 'plat.create'
    },
    U_CLOUD_AREA: {
        id: 'edit_cloud_area',
        name: ['云区域编辑', 'Update Cloud Area'],
        cmdb_action: 'plat.update',
        relation: [{
            view: IAM_VIEWS.CLOUD_AREA,
            instances: [IAM_VIEWS.CLOUD_AREA]
        }],
        transform: (cmdbAction, relationIds) => {
            const [cloudId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: cloudId
            })
        }
    },
    D_CLOUD_AREA: {
        id: 'delete_cloud_area',
        name: ['云区域删除', 'Delete Cloud Area'],
        cmdb_action: 'plat.delete',
        relation: [{
            view: IAM_VIEWS.CLOUD_AREA,
            instances: [IAM_VIEWS.CLOUD_AREA]
        }],
        transform: (cmdbAction, relationIds) => {
            const [cloudId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: cloudId
            })
        }
    },

    // 云账户
    R_CLOUD_ACCOUNT: {
        id: 'find_cloud_account',
        name: ['云账户查询', 'Search Cloud Account'],
        cmdb_action: 'cloudAccount.find',
        relation: [{
            view: IAM_VIEWS.CLOUD_ACCOUNT,
            instances: [IAM_VIEWS.CLOUD_ACCOUNT]
        }],
        transform: (cmdbAction, relationIds) => {
            const [accountId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: accountId
            })
        }
    },
    C_CLOUD_ACCOUNT: {
        id: 'create_cloud_account',
        name: ['云账户创建', 'Create Cloud Account'],
        cmdb_action: 'cloudAccount.create'
    },
    U_CLOUD_ACCOUNT: {
        id: 'edit_cloud_account',
        name: ['云账户编辑', 'Update Cloud Account'],
        cmdb_action: 'cloudAccount.update',
        relation: [{
            view: IAM_VIEWS.CLOUD_ACCOUNT,
            instances: [IAM_VIEWS.CLOUD_ACCOUNT]
        }],
        transform: (cmdbAction, relationIds) => {
            const [accountId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: accountId
            })
        }
    },
    D_CLOUD_ACCOUNT: {
        id: 'delete_cloud_account',
        name: ['云账户删除', 'Delete Cloud Account'],
        cmdb_action: 'cloudAccount.delete',
        relation: [{
            view: IAM_VIEWS.CLOUD_ACCOUNT,
            instances: [IAM_VIEWS.CLOUD_ACCOUNT]
        }],
        transform: (cmdbAction, relationIds) => {
            const [accountId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: accountId
            })
        }
    },

    // 云资源任务
    R_CLOUD_RESOURCE_TASK: {
        id: 'find_cloud_resource_task',
        name: ['云资源任务查询', 'Search Cloud Resource Task'],
        cmdb_action: 'cloudResourceTask.find',
        relation: [{
            view: IAM_VIEWS.CLOUD_RESOURCE_TASK,
            instances: [IAM_VIEWS.CLOUD_RESOURCE_TASK]
        }],
        transform: (cmdbAction, relationIds) => {
            const [taskId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: taskId
            })
        }
    },
    C_CLOUD_RESOURCE_TASK: {
        id: 'create_cloud_resource_task',
        name: ['云资源任务创建', 'Create Cloud Resource Task'],
        cmdb_action: 'cloudResourceTask.create'
    },
    U_CLOUD_RESOURCE_TASK: {
        id: 'edit_cloud_resource_task',
        name: ['云资源任务编辑', 'Update Cloud Resource Task'],
        cmdb_action: 'cloudResourceTask.update',
        relation: [{
            view: IAM_VIEWS.CLOUD_RESOURCE_TASK,
            instances: [IAM_VIEWS.CLOUD_RESOURCE_TASK]
        }],
        transform: (cmdbAction, relationIds) => {
            const [taskId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: taskId
            })
        }
    },
    D_CLOUD_RESOURCE_TASK: {
        id: 'delete_cloud_resource_task',
        name: ['云资源任务删除', 'Delete Cloud Resource Task'],
        cmdb_action: 'cloudResourceTask.delete',
        relation: [{
            view: IAM_VIEWS.CLOUD_RESOURCE_TASK,
            instances: [IAM_VIEWS.CLOUD_RESOURCE_TASK]
        }],
        transform: (cmdbAction, relationIds) => {
            const [taskId] = relationIds
            return basicTransform(cmdbAction, {
                resource_id: taskId
            })
        }
    },

    // 业务自定义字段
    U_BIZ_MODEL_CUSTOM_FIELD: {
        id: 'edit_biz_custom_field',
        name: ['业务自定义字段编辑', 'Update Business Custom Field'],
        cmdb_action: 'modelAttribute.update',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            const [bizId] = relationIds
            return basicTransform(cmdbAction, {
                bk_biz_id: bizId
            })
        }
    },

    // 业务资源查看 (用于控制业务导航下业务选择器的数据)
    R_BIZ_RESOURCE: {
        id: 'find_business_resource',
        name: ['业务资源查看', 'View Business Resource'],
        cmdb_action: 'business.viewBusinessResource',
        relation: [{
            view: IAM_VIEWS.BIZ,
            instances: [IAM_VIEWS.BIZ]
        }],
        transform: (cmdbAction, relationIds) => {
            return basicTransform(cmdbAction, {
                bk_biz_id: relationIds[0]
            })
        }
    }
}

export const OPERATION = {}
Object.keys(IAM_ACTIONS).forEach(key => (OPERATION[key] = key))

// 将配置的权限数据，转换为内部鉴权需要的数据格式, 转换函数均定义在IAM_ACTIONS[someType].transform中
// {
//     bk_biz_id: 1,                    // 业务下的鉴权要带业务id
//     action: 'create',                // 动作
//     resource_type: 'modelInstance',  // 资源类型
//     resource_id: 1,                  // 资源id
//     parent_layers: [{
//         resource_type: 'model',      // 父级依赖，例如实例依赖模型
//         resource_id: 12              // 父级依赖id
//     }]
// }
export function TRANSFORM_TO_INTERNAL (authList) {
    try {
        // 类似导入的鉴权，需要新增、编辑两种权限，统一转成数组处理
        authList = Array.isArray(authList) ? authList : [authList]
        const internalAuthList = []
        authList.forEach(auth => {
            const definition = IAM_ACTIONS[auth.type]
            const customTransform = definition.transform
            let internalAuth
            if (customTransform) {
                internalAuth = customTransform(definition.cmdb_action, auth.relation || [])
            } else {
                internalAuth = basicTransform(definition.cmdb_action)
            }
            // 部分资源可能存在多重权限依赖，因此一个IAM对应的操作，可能需要转换成N种内置的操作，所以均以数组形式返回
            Array.isArray(internalAuth) ? internalAuthList.push(...internalAuth) : internalAuthList.push(internalAuth)
        })
        return internalAuthList
    } catch (error) {
        console.error(error, authList)
        return []
    }
}
