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
    bkSideslider,
    bkTab,
    bkTabPanel,
    bkDialog,
    bkPopover
]
components.forEach(component => {
    Vue.component(`update-${component.name}`, component)
})

Vue.component('bk-button', bkButton)
Vue.component('bk-select', bkSelect)
Vue.component('bk-option', bkOption)
Vue.component('bk-option-group', bkOptionGroup)
