/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import has from 'has'
import {
  MENU_INDEX,
  MENU_BUSINESS,
  MENU_BUSINESS_SET,
  MENU_RESOURCE,
  MENU_MODEL,
  MENU_ANALYSIS,
  MENU_PLATFORM_MANAGEMENT,

  MENU_BUSINESS_SET_TOPOLOGY,

  MENU_BUSINESS_HOST_AND_SERVICE,
  MENU_BUSINESS_CUSTOM_QUERY,
  MENU_BUSINESS_SERVICE_TEMPLATE,
  MENU_BUSINESS_SET_TEMPLATE,
  MENU_BUSINESS_SERVICE_CATEGORY,
  MENU_BUSINESS_CUSTOM_FIELDS,
  MENU_BUSINESS_HOST_APPLY,

  MENU_RESOURCE_MANAGEMENT,
  MENU_RESOURCE_CLOUD_AREA,
  MENU_RESOURCE_CLOUD_ACCOUNT,
  MENU_RESOURCE_CLOUD_RESOURCE,

  MENU_MODEL_MANAGEMENT,
  MENU_MODEL_TOPOLOGY,
  MENU_MODEL_TOPOLOGY_NEW,
  MENU_MODEL_BUSINESS_TOPOLOGY,
  MENU_MODEL_ASSOCIATION,

  MENU_ANALYSIS_AUDIT,
  MENU_ANALYSIS_OPERATION,

  MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG,
} from './menu-symbol'
import {
  businessSetViews,
  businessViews,
  resourceViews,
  modelViews,
  analysisViews,
  platformManagementViews
} from '@/views'

const getMenuRoute = (views, symbol) => {
  const menuView = Array.isArray(views)
    ? views.find(view => view.name === symbol)
    : views
  if (menuView) {
    return {
      name: menuView.name,
      path: menuView.path,
      available: menuView.meta.available
    }
  }
  return {}
}

const menus = [{
  id: MENU_INDEX,
  i18n: '首页'
}, {
  id: MENU_BUSINESS_SET,
  i18n: '业务集',
  visibility: false, // 在一级菜单中不显示
  menu: [{
    id: MENU_BUSINESS_SET_TOPOLOGY,
    i18n: '业务集拓扑',
    icon: 'icon-cc-host',
    route: getMenuRoute(businessSetViews, MENU_BUSINESS_SET_TOPOLOGY)
  }]
}, {
  id: MENU_BUSINESS,
  i18n: '业务',
  menu: [{
    id: MENU_BUSINESS_HOST_AND_SERVICE,
    i18n: '业务拓扑',
    icon: 'icon-cc-host',
    route: getMenuRoute(businessViews, MENU_BUSINESS_HOST_AND_SERVICE)
  }, {
    id: MENU_BUSINESS_SERVICE_TEMPLATE,
    i18n: '服务模板',
    icon: 'icon-cc-service-template',
    route: getMenuRoute(businessViews, MENU_BUSINESS_SERVICE_TEMPLATE)
  }, {
    id: MENU_BUSINESS_SET_TEMPLATE,
    i18n: '集群模板',
    icon: 'icon-cc-set-template',
    route: getMenuRoute(businessViews, MENU_BUSINESS_SET_TEMPLATE)
  }, {
    id: MENU_BUSINESS_SERVICE_CATEGORY,
    i18n: '服务分类',
    icon: 'icon-cc-nav-service-topo',
    route: getMenuRoute(businessViews, MENU_BUSINESS_SERVICE_CATEGORY)
  }, {
    id: MENU_BUSINESS_HOST_APPLY,
    i18n: '主机自动应用',
    icon: 'icon-cc-host-apply',
    route: getMenuRoute(businessViews, MENU_BUSINESS_HOST_APPLY)
  }, {
    id: MENU_BUSINESS_CUSTOM_QUERY,
    i18n: '动态分组',
    icon: 'icon-cc-custom-query',
    route: getMenuRoute(businessViews, MENU_BUSINESS_CUSTOM_QUERY)
  }, {
    id: MENU_BUSINESS_CUSTOM_FIELDS,
    i18n: '自定义字段',
    icon: 'icon-cc-custom-field',
    route: getMenuRoute(businessViews, MENU_BUSINESS_CUSTOM_FIELDS)
  }]
}, {
  id: MENU_RESOURCE,
  i18n: '资源',
  menu: [{
    id: MENU_RESOURCE_MANAGEMENT,
    i18n: '资源目录',
    icon: 'icon-cc-square',
    route: getMenuRoute(resourceViews, MENU_RESOURCE_MANAGEMENT)
  }, {
    id: MENU_RESOURCE_CLOUD_AREA,
    i18n: '管控区域',
    icon: 'icon-cc-network-segment',
    route: getMenuRoute(resourceViews, MENU_RESOURCE_CLOUD_AREA, 'resource')
  }, {
    id: MENU_RESOURCE_CLOUD_ACCOUNT,
    i18n: '云账户',
    icon: 'icon-cc-cloud-account',
    route: getMenuRoute(resourceViews, MENU_RESOURCE_CLOUD_ACCOUNT, 'resource')
  }, {
    id: MENU_RESOURCE_CLOUD_RESOURCE,
    i18n: '云资源发现',
    icon: 'icon-cc-cloud-discover',
    route: getMenuRoute(resourceViews, MENU_RESOURCE_CLOUD_RESOURCE, 'resource')
  }]
}, {
  id: MENU_MODEL,
  i18n: '模型',
  menu: [{
    id: MENU_MODEL_MANAGEMENT,
    i18n: '模型管理',
    icon: 'icon-cc-nav-model-02',
    route: getMenuRoute(modelViews, MENU_MODEL_MANAGEMENT)
  }, {
    id: MENU_MODEL_TOPOLOGY,
    i18n: '模型拓扑',
    icon: 'icon-cc-nav-model-topo',
    route: getMenuRoute(modelViews, MENU_MODEL_TOPOLOGY)
  }, {
    id: MENU_MODEL_TOPOLOGY_NEW,
    i18n: '模型关系',
    icon: 'icon-cc-nav-model-topo',
    route: getMenuRoute(modelViews, MENU_MODEL_TOPOLOGY_NEW)
  }, {
    id: MENU_MODEL_BUSINESS_TOPOLOGY,
    i18n: '业务层级',
    icon: 'icon-cc-tree',
    route: getMenuRoute(modelViews, MENU_MODEL_BUSINESS_TOPOLOGY)
  }, {
    id: MENU_MODEL_ASSOCIATION,
    i18n: '关联类型',
    icon: 'icon-cc-nav-associated',
    route: getMenuRoute(modelViews, MENU_MODEL_ASSOCIATION)
  }]
}, {
  id: MENU_ANALYSIS,
  i18n: '运营分析',
  menu: [{
    id: MENU_ANALYSIS_AUDIT,
    i18n: '操作审计',
    icon: 'icon-cc-nav-audit-02',
    route: getMenuRoute(analysisViews, MENU_ANALYSIS_AUDIT)
  }, {
    id: MENU_ANALYSIS_OPERATION,
    i18n: '运营统计',
    icon: 'icon-cc-statistics',
    route: getMenuRoute(analysisViews, MENU_ANALYSIS_OPERATION)
  }]
}, {
  id: MENU_PLATFORM_MANAGEMENT,
  i18n: '平台管理',
  visibility: componentContext => componentContext.$store.state.globalConfig.auth,
  menu: [{
    id: MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG,
    i18n: '全局配置',
    icon: 'icon-cc-setting',
    route: getMenuRoute(platformManagementViews, MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG)
  }]
}]

// 移除未被激活的menu
;(() => {
  menus.forEach((top) => {
    if (has(top, 'menu')) {
      top.menu.forEach((menu) => {
        if (has(menu, 'submenu')) {
          menu.submenu = menu.submenu.filter(submenu => submenu.route.available)
        }
      })
      top.menu = top.menu.filter((menu) => {
        if (has(menu, 'route')) {
          return menu.route.available
        }
        return menu.submenu.length
      })
    }
  })
})()

export default menus
