import Vue from 'vue'
import i18n from '@/i18n'
import store from '@/store'
import ColumnsConfig from './columns-config.vue'
export default {
  open({ props = {}, handler = {} }) {
    const vm = new Vue({
      i18n,
      store,
      data() {
        return {
          isShow: false
        }
      },
      render(h) {
        return h('bk-sideslider', {
          ref: 'sideslider',
          props: {
            title: i18n.t('列表显示属性配置'),
            width: 600,
            isShow: this.isShow
          },
          on: {
            'update:isShow': (isShow) => {
              this.isShow = isShow
            },
            'animation-end': () => {
              this.$el && this.$el.parentElement && this.$el.parentElement.removeChild(this.$el)
              this.$destroy()
            }
          }
        }, [h(ColumnsConfig, {
          props,
          slot: 'content',
          on: {
            cancel: () => {
              this.isShow = false
            },
            apply: (properties) => {
              this.isShow = false
              handler.apply && handler.apply(properties)
            },
            reset: () => {
              this.isShow = false
              handler.reset && handler.reset()
            }
          }
        })])
      }
    })
    vm.$mount()
    document.body.appendChild(vm.$el)
    vm.isShow = true
    return vm
  }
}
