import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import AuditDetailsSlider from './details-slider.vue'
const Component = Vue.extend({
    components: {
        AuditDetailsSlider
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
        return <audit-details-slider ref="details" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></audit-details-slider>
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
        vm.$refs.details.show()
    }
}
