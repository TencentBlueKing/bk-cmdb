import bkSelect from './src/select'
import bkSelectOption from './src/select-option'
import bkOptionGroup from './src/option-group'

bkSelect.install = Vue => {
    Vue.component(bkSelect.name, bkSelect)
}

bkSelectOption.install = Vue => {
    Vue.component(bkSelectOption.name, bkSelectOption)
}

bkOptionGroup.install = Vue => {
    Vue.component(bkOptionGroup.name, bkOptionGroup)
}

export default {
    bkSelect,
    bkSelectOption,
    bkOptionGroup
}
