import bkTab from './tab'
import bkTabpanel from './tab-panel'

bkTab.install = Vue => {
    Vue.component(bkTab.name, bkTab)
}

bkTabpanel.install = Vue => {
    Vue.component(bkTabpanel.name, bkTabpanel)
}

export default {
    bkTab,
    bkTabpanel
}
