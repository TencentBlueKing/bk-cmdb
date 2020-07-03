import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import LabelBatchDialog from './label-batch-dialog.vue'
const Component = Vue.extend({
    components: {
        LabelBatchDialog
    },
    methods: {
        handleClose () {
            document.body.removeChild(this.$el)
            this.$destroy()
        }
    },
    render (h) {
        return <label-batch-dialog ref="dialog" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></label-batch-dialog>
    }
})

export default {
    show (data = {}) {
        const vm = new Component({
            store,
            i18n,
            attrs: data
        })
        vm.$mount()
        document.body.appendChild(vm.$el)
        vm.$refs.dialog.show()
    }
}
