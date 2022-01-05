import Meta from '@/router/meta'
import {
  MENU_BUSINESS_SET,
  MENU_BUSINESS_SET_TOPOLOGY
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
  })
}]
