import Vue from 'vue'
import i18n from '@/i18n'
import store from '@/store'
import FilterForm from './filter-form.vue'
import FilterStore from './store'
export default {
  show(props = {}) {
    const exist = FilterStore.getComponent('FilterForm')
    if (exist) {
      exist.$refs.FilterForm.focusIP()
      return exist
    }
    const vm = new Vue({
      i18n,
      store,
      created() {
        FilterStore.setComponent('FilterForm', this)
      },
      beforeDestroy() {
        FilterStore.setComponent('FilterForm', null)
      },
      render(h) {
        return h(FilterForm, {
          ref: 'FilterForm',
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
    vm.$refs.FilterForm.open()
    vm.$watch(() => window.CMDB_APP.$route.name, vm.$refs.FilterForm.close)
    return vm
  }
}
