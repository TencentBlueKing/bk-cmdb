import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import DynamicGroupPreview from './preview.vue'
const Component = Vue.extend({
    components: {
        DynamicGroupPreview
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
        return (<dynamic-group-preview ref="preview" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></dynamic-group-preview>)
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
        vm.$refs.preview.show()
    }
}
