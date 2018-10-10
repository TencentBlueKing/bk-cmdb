import Vue from 'vue'
import businessSelector from './selector/business.vue'
import clipboardSelector from './selector/clipboard.vue'
import selector from './selector/selector.vue'
import table from './table/table.vue'
import tableSelector from './table/table-selector.vue'
import slider from './slider/slider.vue'
import details from './details/details.vue'
import form from './form/form.vue'
import formMultiple from './form/form-multiple.vue'
import bool from './form/bool.vue'
import boolInput from './form/bool-input.vue'
import date from './form/date.vue'
import dateRange from './form/date-range.vue'
import time from './form/time.vue'
import int from './form/int.vue'
import longchar from './form/longchar.vue'
import singlechar from './form/singlechar.vue'
import timezone from './form/timezone.vue'
import enumeration from './form/enum.vue'
import objuser from './form/objuser.vue'
import associateInput from './form/associate-input.vue'
import tree from './tree/tree.vue'
import resize from './other/resize.vue'
import collapseTransition from './transition/collapse.js'
const install = (Vue, opts = {}) => {
    const components = [
        businessSelector,
        clipboardSelector,
        selector,
        table,
        tableSelector,
        slider,
        details,
        form,
        formMultiple,
        bool,
        boolInput,
        date,
        dateRange,
        time,
        int,
        longchar,
        singlechar,
        timezone,
        enumeration,
        objuser,
        associateInput,
        tree,
        resize,
        collapseTransition
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
    table,
    tableSelector,
    slider,
    details,
    form,
    formMultiple,
    bool,
    boolInput,
    date,
    dateRange,
    time,
    int,
    longchar,
    singlechar,
    timezone,
    enumeration,
    objuser,
    associateInput,
    tree,
    resize,
    collapseTransition
}
