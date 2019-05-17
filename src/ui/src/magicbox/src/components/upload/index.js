import bkUpload from './upload.vue'

bkUpload.install = Vue => {
    Vue.component(bkUpload.name, bkUpload)
}

export default bkUpload
