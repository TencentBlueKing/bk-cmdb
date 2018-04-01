import bkSelector from './src/selector'

bkSelector.install = Vue => {
    Vue.component(bkSelector.name, bkSelector)
}

export default bkSelector
