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
    bkCheckbox,
    bkMessage,
    bkLoading,
    bkBigTree
} from 'bk-magic-vue'

Vue.use(bkButton)
Vue.use(bkInput)
Vue.use(bkTab)
Vue.use(bkTabPanel)
Vue.use(bkSideslider)
Vue.use(bkSelect)
Vue.use(bkOption)
Vue.use(bkOptionGroup)
Vue.use(bkTable)
Vue.use(bkTableColumn)
Vue.use(bkCheckbox)
Vue.use(bkPagination)
Vue.use(bkDatePicker)
Vue.use(bkDialog)
Vue.use(bkPopover)
Vue.use(bkDropdownMenu)
Vue.use(bkLoading)
Vue.use(bkBigTree)

export const $error = (message, delay = 3000) => {
    bkMessage({
        message,
        delay,
        theme: 'error'
    })
}

export const $success = (message, delay = 3000) => {
    bkMessage({
        message,
        delay,
        theme: 'success'
    })
}

export const $info = (message, delay = 3000) => {
    bkMessage({
        message,
        delay,
        theme: 'primary'
    })
}

export const $warn = (message, delay = 3000) => {
    bkMessage({
        message,
        delay,
        theme: 'warning',
        hasCloseIcon: true
    })
}

Vue.prototype.$error = $error
Vue.prototype.$success = $success
Vue.prototype.$info = $info
Vue.prototype.$warn = $warn
