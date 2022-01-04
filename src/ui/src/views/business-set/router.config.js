import Meta from '@/router/meta'
import {
  MENU_RESOURCE_BUSINESS_SET
} from '@/dictionary/menu-symbol.js'

export default [
  {
    name: MENU_RESOURCE_BUSINESS_SET,
    path: 'business-set',
    component: () => import('./index.vue'),
    meta: new Meta({
      menu: {
        i18n: '业务集'
      }
    })
  }
]
