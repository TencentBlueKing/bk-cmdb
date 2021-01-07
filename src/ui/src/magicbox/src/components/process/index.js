import bkProcess from './process.vue'

bkProcess.install = Vue => {
    Vue.component(bkProcess.name, bkProcess)
}

export default bkProcess
