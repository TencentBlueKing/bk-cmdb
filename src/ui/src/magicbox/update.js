import Vue from 'vue'
import {
    bkInput,
    bkRadio,
    bkCheckbox,
    bkDropdownMenu,
    bkDatePicker,
    bkTimePicker,
    bkTagInput,
    bkSearchSelect,
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
    bkPopover
} from 'bk-magic-vue'

const components = [
    bkInput,
    bkRadio,
    bkCheckbox,
    bkDropdownMenu,
    bkDatePicker,
    bkTimePicker,
    bkTagInput,
    bkSearchSelect,
    bkTable,
    bkTableColumn,
    bkPagination,
    bkDialog,
    bkPopover
]
components.forEach(component => {
    Vue.component(`update-${component.name}`, component)
})

Vue.component('bk-button', bkButton)
Vue.component('bk-tab', bkTab)
Vue.component('bk-tab-panel', bkTabPanel)
Vue.component('bk-sideslider', bkSideslider)
Vue.component('bk-select', bkSelect)
Vue.component('bk-option', bkOption)
Vue.component('bk-option-group', bkOptionGroup)
