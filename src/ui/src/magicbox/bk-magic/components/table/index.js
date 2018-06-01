import bkTable from './src/table'
import Paging from '../paging/src/paging'

bkTable.install = Vue => {
    Vue.component(bkTable.name, bkTable)
    Vue.component(Paging.name, Paging)
}

export default bkTable
