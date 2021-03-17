import Vue from 'vue'
import i18n from '@/i18n'
import store from '@/store'
import PropertySelector from './property-selector.vue'
export default {
  show(props = {}) {
    const vm = new Vue({
      i18n,
      store,
      render(h) {
        return h(PropertySelector, {
          ref: 'propertySelector',
          props,
          on: {
            closed: () => {
              this.$el && this.$el.parentElement && this.$el.parentElement.removeChild(this.$el)
              this.$destroy()
            }
          }
        })
      }
    })
    vm.$mount()
    document.body.appendChild(vm.$el)
    vm.$refs.propertySelector.open()
    return vm
  }
}
