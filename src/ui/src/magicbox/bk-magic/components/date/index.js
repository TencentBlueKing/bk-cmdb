import bkDate from './src/date'

bkDate.install = Vue => {
    Vue.component(bkDate.name, bkDate)
}

export default bkDate
