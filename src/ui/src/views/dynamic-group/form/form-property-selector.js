import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import FormPropertySelector from './form-property-selector.vue'
const Component = Vue.extend({
    components: {
        FormPropertySelector
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
            // magicbox实现相关，多个侧滑同时存在，后面的不会挂在到body中，而是挂在到popmanager中，此处不手动移出
            // document.body.removeChild(this.$el)
            this.$destroy()
        }
    },
    render (h) {
        return (<form-property-selector ref="selector" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></form-property-selector>)
    }
})

export default {
    show (data = {}, dynamicGroupForm) {
        const vm = new Component({
            store,
            i18n,
            attrs: data,
            provide: () => {
                return {
                    dynamicGroupForm
                }
            }
        })
        vm.$mount()
        // magicbox实现相关，多个侧滑同时存在，后面的不会挂在到body中，而是挂在到popmanager中，此处不手动添加
        // document.body.appendChild(vm.$el)
        vm.$refs.selector.show()
    }
}
