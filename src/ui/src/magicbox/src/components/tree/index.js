import bkTree from './tree.vue'

bkTree.install = Vue => {
    Vue.component(bkTree.name, bkTree)
}

export default bkTree
