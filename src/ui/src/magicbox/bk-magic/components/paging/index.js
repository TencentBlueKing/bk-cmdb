import bkPaging from './src/paging'

bkPaging.install = Vue => {
    Vue.component(bkPaging.name, bkPaging)
}

export default bkPaging
