import bkPagination from './pagination'

bkPagination.install = Vue => {
    Vue.component(bkPagination.name, bkPagination)
}

export default bkPagination
