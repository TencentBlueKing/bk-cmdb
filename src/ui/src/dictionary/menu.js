export const NAV_INDEX = 'index'
export const NAV_BASIC_RESOURCE = 'basicResource'
export const NAV_BUSINESS_RESOURCE = 'businessResource'
export const NAV_AUDIT_ANALYSE = 'auditAnalyse'
export const NAV_MODEL_MANAGEMENT = 'modelManagement'
export const NAV_COLLECT = 'collect'

export default [{
    id: NAV_INDEX,
    i18n: 'Nav["首页"]',
    icon: 'bk-icon icon-home-shape'
}, {
    id: NAV_BASIC_RESOURCE,
    i18n: 'Nav["基础资源"]',
    icon: 'icon-cc-basic-resource',
    submenu: []
}, {
    id: NAV_BUSINESS_RESOURCE,
    i18n: 'Nav["业务资源"]',
    icon: 'icon-cc-nav-resource',
    submenu: []
}, {
    id: NAV_AUDIT_ANALYSE,
    i18n: 'Nav["审计与分析"]',
    icon: 'icon-cc-nav-audit',
    submenu: []
}, {
    id: NAV_MODEL_MANAGEMENT,
    i18n: 'Nav["模型管理"]',
    icon: 'icon-cc-nav-model',
    submenu: []
}, {
    id: NAV_COLLECT,
    i18n: 'Nav["我的收藏"]',
    icon: 'icon-cc-nav-collection',
    submenu: []
}]
