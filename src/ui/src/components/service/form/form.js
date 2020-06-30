import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import ProcessForm from './form.vue'
const Component = Vue.extend({
    components: {
        ProcessForm
    },
    created () {
        this.unwatch = RouterQuery.watch('*', () => {
            this.handleClose()
        })
    },
    beforeDestroy () {
        this.unwatch()
    },
    methods: {
        handleClose () {
            document.body.removeChild(this.$el)
            this.$destroy()
        }
    },
    render (h) {
        return <process-form ref="form" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></process-form>
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
        vm.$refs.form.show()
    }
}
