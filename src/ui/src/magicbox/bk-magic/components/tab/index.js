import bkTab from './src/tab'
import bkTabPanel from './src/tab-panel'

bkTab.install = Vue => {
    Vue.component(bkTab.name, bkTab)
}

bkTabPanel.install = Vue => {
    Vue.component(bkTabPanel.name, bkTabPanel)
}

export default {
    bkTab,
    bkTabPanel
}
