import Vue from 'vue'
import contentComponent from './export'
import store from '@/store'
import i18n from '@/i18n'
import useState from './state'
import useTask from './task'
import { bkInfoBox } from 'bk-magic-vue'
import { toRef, watch } from '@vue/composition-api'
let instance = null
const [state, { setState, resetState }] = useState()
const [, { reset: resetTask }] = useTask()
const visible = toRef(state, 'visible')
const title = toRef(state, 'title')
const step = toRef(state, 'step')
const status = toRef(state, 'status')
const show = () => {
  instance = createSideslider()
  instance.$mount()
  visible.value = true
}
const close = () => {
  visible.value = false
}

watch(visible, (value) => {
  if (value) return
  resetState()
  resetTask()
  setTimeout(() => {
    instance.$destroy()
    instance = null
  }, 200)
})

const beforeClose = () => {
  if (step.value !== 2 || !status.value === 'pending') {
    return true
  }
  return new Promise((resolve) => {
    bkInfoBox({
      title: i18n.t('确认关闭'),
      subTitle: i18n.t('数据导出终止提示'),
      confirmFn: () => resolve(true),
      cancelFn: () => resolve(false)
    })
  })
}

/**
 * 生成导入侧滑组件实例
 */
const createSideslider = () => {
  const Component = Vue.extend({
    components: {
      contentComponent
    },
    mounted() {
      document.body.appendChild(this.$el)
    },
    beforeDestroy() {
      try {
        // document.body.removeChild(this.$el)
      } catch (error) {
        console.error(error)
      }
    },
    render() {
      return (
        <bk-sideslider
          is-show={ visible.value }
          width={ 700 }
          title={ title.value }
          before-close={ beforeClose }
          { ...{ on: { 'update:isShow': close } } }>
          <content-component slot="content"></content-component>
        </bk-sideslider>
      )
    }
  })
  return new Component({ store, i18n })
}

export default function (state) {
  resetState()
  setState(state)
  return { close, show }
}
