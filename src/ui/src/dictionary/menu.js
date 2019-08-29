export const NAV_INDEX = 'index'
export const NAV_BASIC_RESOURCE = 'basicResource'
export const NAV_BUSINESS_RESOURCE = 'businessResource'
export const NAV_AUDIT_ANALYSE = 'auditAnalyse'
export const NAV_MODEL_MANAGEMENT = 'modelManagement'
export const NAV_PERMISSION = 'permission'
export const NAV_COLLECT = 'collect'
export const NAV_SERVICE_MANAGEMENT = 'serviceManagement'

export default [{
    id: NAV_INDEX,
    i18n: '首页',
    icon: 'bk-icon icon-home-shape'
}, {
    id: NAV_BASIC_RESOURCE,
    i18n: '基础资源',
    icon: 'icon-cc-resource',
    submenu: []
}, {
    id: NAV_BUSINESS_RESOURCE,
    i18n: '业务资源',
    icon: 'icon-cc-nav-resource',
    submenu: []
}, {
    id: NAV_SERVICE_MANAGEMENT,
    i18n: '服务管理',
    icon: 'icon-cc-template-management',
    submenu: []
}, {
    id: NAV_AUDIT_ANALYSE,
    i18n: '审计与分析',
    icon: 'icon-cc-nav-audit',
    submenu: []
}, {
    id: NAV_PERMISSION,
    i18n: '权限控制',
    icon: 'icon-cc-nav-authority',
    submenu: []
}, {
    id: NAV_MODEL_MANAGEMENT,
    i18n: '模型管理',
    icon: 'icon-cc-nav-model',
    submenu: []
}, {
    id: NAV_COLLECT,
    i18n: '我的收藏',
    icon: 'icon-cc-nav-collection',
    submenu: []
}]
