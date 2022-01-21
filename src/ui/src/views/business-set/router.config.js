import Meta from '@/router/meta'
import {
  MENU_RESOURCE_BUSINESS_SET,
  MENU_RESOURCE_BUSINESS_SET_DETAILS
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
  },
  {
    name: MENU_RESOURCE_BUSINESS_SET_DETAILS,
    path: 'business-set/details/:bizSetId',
    component: () => import('./details.vue'),
    meta: new Meta({
      menu: {
        i18n: '业务集详情',
        relative: MENU_RESOURCE_BUSINESS_SET
      },
      layout: {}
    })
  }
]
