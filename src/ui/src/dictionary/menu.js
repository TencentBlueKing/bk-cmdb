import {
    MENU_BUSINESS_HOST,
    MENU_BUSINESS_SERVICE,
    MENU_BUSINESS_ADVANCED,
    MENU_RESOURCE_EVENTPUSH,
    MENU_MODEL_MANAGEMENT,
    MENU_MODEL_TOPOLOGY,
    MENU_MODEL_ASSOCIATION,
    MENU_ANALYSIS_AUDIT,
    MENU_ANALYSIS_STATISTICS
} from './menu-symbol'
import {
    businessViews,
    resourceViews,
    modelViews,
    analysisViews
} from '@/views'

export const NAV_INDEX = 'index'
export const NAV_BASIC_RESOURCE = 'basicResource'
export const NAV_BUSINESS_RESOURCE = 'businessResource'
export const NAV_AUDIT_ANALYSE = 'auditAnalyse'
export const NAV_MODEL_MANAGEMENT = 'modelManagement'
export const NAV_PERMISSION = 'permission'
export const NAV_COLLECT = 'collect'
export const NAV_SERVICE_MANAGEMENT = 'serviceManagement'

export const HEADER_NAV = [{
    name: 'index',
    i18n: '首页'
}, {
    name: 'business',
    i18n: '业务'
}, {
    name: 'resource',
    i18n: '资源'
}, {
    name: 'model',
    i18n: '模型'
}, {
    name: 'analysis',
    i18n: '运营分析'
}]

const getSubmenu = (views, parentSymbol, pathPrefix = '') => {
    const submenuViews = views.filter(view => {
        return view.meta && view.meta.menu && view.meta.menu.parent === parentSymbol
    })
    const submenu = submenuViews.map(view => {
        const menu = view.meta.menu
        return {
            i18n: menu.i18n,
            route: getMenuRoute(view, parentSymbol, pathPrefix)
        }
    })
    return submenu
}

const getMenuRoute = (views, parentSymbol, pathPrefix = '') => {
    const menuView = Array.isArray(views)
        ? views.find(view => view.meta && view.meta.menu && view.meta.menu.parent === parentSymbol)
        : views
    if (menuView) {
        return {
            name: menuView.name,
            path: menuView.path ? `/${pathPrefix}/${menuView.path}` : undefined
        }
    }
    return {}
}

export const BUSINESS_MENU = [{
    i18n: '主机',
    icon: 'icon-cc-resource',
    submenu: getSubmenu(businessViews, MENU_BUSINESS_HOST, 'business')
}, {
    i18n: '服务',
    icon: 'icon-cc-template-management',
    submenu: getSubmenu(businessViews, MENU_BUSINESS_SERVICE, 'business')
}, {
    i18n: '高级功能',
    icon: 'icon-cc-plus-circle',
    submenu: getSubmenu(businessViews, MENU_BUSINESS_ADVANCED, 'business')
}]

export const RESOURCE_MENU = [{
    i18n: '事件推送',
    icon: 'cc-square',
    route: getMenuRoute(resourceViews, MENU_RESOURCE_EVENTPUSH, 'resource')
}]

export const MODEL_MENU = [{
    i18n: '模型管理',
    icon: 'icon-cc-nav-model',
    submenu: getSubmenu(modelViews, MENU_MODEL_MANAGEMENT, 'model')
}, {
    i18n: '模型关系',
    icon: 'icon-cc-resources',
    submenu: getSubmenu(modelViews, MENU_MODEL_TOPOLOGY, 'model')
}, {
    i18n: '关联分类',
    icon: 'icon-cc-network-manage',
    submenu: getSubmenu(modelViews, MENU_MODEL_ASSOCIATION, 'model')
}]

export const ANALYSIS_MENU = [{
    i18n: '操作审计',
    icon: 'icon-cc-statement',
    route: getMenuRoute(analysisViews, MENU_ANALYSIS_AUDIT, 'analysis')
}, {
    i18n: '运营统计',
    icon: 'icon-cc-statement',
    route: getMenuRoute(analysisViews, MENU_ANALYSIS_STATISTICS, 'analysis')
}]

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
