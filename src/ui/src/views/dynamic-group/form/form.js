import Vue from 'vue'
import store from '@/store'
import i18n from '@/i18n'
import RouterQuery from '@/router/query'
import dynamicGroupForm from './form.vue'
const Component = Vue.extend({
  components: {
    dynamicGroupForm
  },
  created() {
    this.unwatch = RouterQuery.watch('*', () => {
      this.handleClose()
    })
  },
  beforeDestroy() {
    this.unwatch()
  },
  methods: {
    handleClose() {
      document.body.removeChild(this.$el)
      this.$destroy()
    }
  },
  // eslint-disable-next-line no-unused-vars
  render(h) {
    return (<dynamic-group-form ref="form" { ...{ props: this.$options.attrs }} on-close={ this.handleClose }></dynamic-group-form>)
  }
})

export default {
  show(data = {}) {
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
