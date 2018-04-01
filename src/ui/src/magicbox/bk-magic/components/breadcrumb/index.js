import bkBreadcrumb from './src/breadcrumb'

bkBreadcrumb.install = Vue => {
  Vue.component(bkBreadcrumb.name, bkBreadcrumb)
}

export default bkBreadcrumb
