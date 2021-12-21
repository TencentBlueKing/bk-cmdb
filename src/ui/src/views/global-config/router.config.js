import Meta from '@/router/meta'
import { MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG, MENU_PLATFORM_MANAGEMENT } from '@/dictionary/menu-symbol.js'
import { OPERATION } from '@/dictionary/iam-auth'

export default [{
  name: MENU_PLATFORM_MANAGEMENT_GLOBAL_CONFIG,
  path: 'global-config',
  component: () => import('./index.vue'),
  meta: new Meta({
    menu: {
      i18n: '全局配置',
      parent: MENU_PLATFORM_MANAGEMENT
    },
    auth: {
      view: { type: OPERATION.U_CONFIG_ADMIN }
    },
    layout: {
      breadcrumbs: true
    },
  })
}]
