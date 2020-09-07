import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import ProcessFormMultiple from './form-multiple.vue'
const Component = Vue.extend({
    components: {
        ProcessFormMultiple
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
        return <process-form-multiple ref="form" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></process-form-multiple>
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
