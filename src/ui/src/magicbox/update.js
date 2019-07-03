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
    bkSelect,
    bkTab,
    bkTabPanel,
    bkDialog,
    bkPopover
]
components.forEach(component => {
    Vue.component(`update-${component.name}`, component)
})

Vue.component('bk-button', bkButton)
