import Meta from '@/router/meta'
import {
  MENU_BUSINESS_SET,
  MENU_BUSINESS_HOST_AND_SERVICE,
  MENU_BUSINESS_SET_TOPOLOGY,
  MENU_BUSINESS_SET_HOST_DETAILS
} from '@/dictionary/menu-symbol'

export default [{
  name: MENU_BUSINESS_SET_TOPOLOGY,
  path: 'index',
  component: () => import('./index.vue'),
  meta: new Meta({
    owner: MENU_BUSINESS_SET,
    menu: {
      i18n: '业务集拓扑'
    },
    customInstanceColumn: 'business_set_topology_table_column_config'
  }),
  children: [{
    name: MENU_BUSINESS_SET_HOST_DETAILS,
    path: 'host/:id',
    component: () => import('@/views/host-details/index.vue'),
    meta: new Meta({
      owner: MENU_BUSINESS_SET,
      readonly: true,
      menu: {
        i18n: '主机详情',
        relative: MENU_BUSINESS_HOST_AND_SERVICE,
      }
    })
  }]
}]
