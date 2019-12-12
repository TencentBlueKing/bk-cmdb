import businessSelector from './selector/business.vue'
import clipboardSelector from './selector/clipboard.vue'
import selector from './selector/selector.vue'
import cloudSelector from './selector/cloud.vue'
import details from './details/details.vue'
import form from './form/form.vue'
import formMultiple from './form/form-multiple.vue'
import bool from './form/bool.vue'
import boolInput from './form/bool-input.vue'
import date from './form/date.vue'
import dateRange from './form/date-range.vue'
import time from './form/time.vue'
import int from './form/int.vue'
import float from './form/float.vue'
import longchar from './form/longchar.vue'
import singlechar from './form/singlechar.vue'
import timezone from './form/timezone.vue'
import enumeration from './form/enum.vue'
import objuser from './form/objuser.vue'
import resize from './other/resize.vue'
import collapseTransition from './transition/collapse.js'
import collapse from './collapse/collapse'
import dotMenu from './dot-menu/dot-menu.vue'
import input from './form/input.vue'
import searchInput from './form/search-input.vue'
import inputSelect from './selector/input-select.vue'
import iconButton from './button/icon-button.vue'
import tips from './other/tips.vue'
import dialog from './dialog/dialog.vue'
import auth from './auth/auth.vue'
import tableEmpty from './table-empty/table-empty.vue'
import list from './form/list.vue'
import leaveConfirm from './dialog/leave-confirm.vue'
import user from './user/user.vue'
const install = (Vue, opts = {}) => {
    const components = [
        businessSelector,
        clipboardSelector,
        selector,
        details,
        form,
        formMultiple,
        bool,
        boolInput,
        date,
        dateRange,
        time,
        int,
        float,
        longchar,
        singlechar,
        timezone,
        enumeration,
        objuser,
        resize,
        collapseTransition,
        collapse,
        dotMenu,
        input,
        searchInput,
        inputSelect,
        iconButton,
        tips,
        dialog,
        cloudSelector,
        auth,
        tableEmpty,
        list,
        leaveConfirm,
        user
    ]
    components.forEach(component => {
        Vue.component(component.name, component)
    })
}

export default {
    install,
    businessSelector,
    clipboardSelector,
    selector,
    details,
    form,
    formMultiple,
    bool,
    boolInput,
    date,
    dateRange,
    time,
    int,
    float,
    longchar,
    singlechar,
    timezone,
    enumeration,
    objuser,
    resize,
    collapseTransition,
    dotMenu,
    input,
    searchInput,
    inputSelect,
    iconButton,
    tips,
    dialog,
    cloudSelector,
    auth,
    tableEmpty,
    list,
    leaveConfirm
}
