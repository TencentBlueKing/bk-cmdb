/**
 * @file datepicker entry
 * @author ielgnaw <wuji0223@gmail.com>
 */

import bkDatePicker from './date-picker.vue'

bkDatePicker.install = Vue => {
    Vue.component(bkDatePicker.name, bkDatePicker)
}

export default bkDatePicker
