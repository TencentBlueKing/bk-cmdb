import Vue from 'vue'
import {
    bkInput,
    bkDropdownMenu,
    bkDatePicker,
    bkTable,
    bkTableColumn,
    bkPagination,
    bkSideslider,
    bkButton,
    bkSelect,
    bkOption,
    bkOptionGroup,
    bkTab,
    bkTabPanel,
    bkDialog,
    bkPopover,
    bkCheckbox
} from 'bk-magic-vue'

const components = [
    bkDropdownMenu,
    bkPopover
]
components.forEach(component => {
    Vue.component(`update-${component.name}`, component)
})

Vue.component('bk-button', bkButton)
Vue.component('bk-input', bkInput)
Vue.component('bk-tab', bkTab)
Vue.component('bk-tab-panel', bkTabPanel)
Vue.component('bk-sideslider', bkSideslider)
Vue.component('bk-select', bkSelect)
Vue.component('bk-option', bkOption)
Vue.component('bk-option-group', bkOptionGroup)
Vue.component('bk-table', bkTable)
Vue.component('bk-table-column', bkTableColumn)
Vue.component('bk-checkbox', bkCheckbox)
Vue.component('bk-pagination', bkPagination)
Vue.component('bk-date-picker', bkDatePicker)
Vue.component('bk-dialog', bkDialog)
