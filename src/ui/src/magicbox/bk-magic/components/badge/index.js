import bkBadge from './src/badge'

bkBadge.install = Vue => {
    Vue.component(bkBadge.name, bkBadge)
}

export default bkBadge
