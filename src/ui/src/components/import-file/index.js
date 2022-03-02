import { reactive } from '@vue/composition-api'
import Vue from 'vue'
import contentComponent from './import'
import store from '@/store'
import i18n from '@/i18n'
import useStep from './step'
import useFile from './file'
import { bkInfoBox } from 'bk-magic-vue'
let instance = null
const state = reactive({
  visible: false,
  title: null,
  submit: null,
  success: null,
  bk_obj_id: null,
  bk_biz_id: null,
  fileTips: null,
  template: null
})
const setState = (newState) => {
  Object.assign(state, newState)
}

const show = (props = {}) => {
  instance = createSideslider(props)
  instance.$mount()
  state.visible = true
}
const close = () => {
  setState({
    visible: false,
    title: null,
    submit: null,
    success: null,
    bk_biz_id: null,
    bk_obj_id: null,
    fileTips: null,
    template: null
  })
  setTimeout(() => {
    instance.$destroy()
    instance = null
    const [, { reset: resetStep }] = useStep()
    resetStep()
    const [, { clear: clearFile }] = useFile()
    clearFile()
  }, 200)
}

const beforeClose = () => {
  const [{ state }] = useFile()
  const [current] = useStep()
  if (current.value !== 2 || !state.value === 'pending') {
    return true
  }
  return new Promise((resolve) => {
    bkInfoBox({
      title: i18n.t('确认关闭'),
      subTitle: i18n.t('数据导入终止提示'),
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
          is-show={ state.visible }
          width={ 700 }
          title={ state.title }
          before-close={ beforeClose }
          { ...{ on: { 'update:isShow': close } } }>
          <content-component slot="content"></content-component>
        </bk-sideslider>
      )
    }
  })
  return new Component({ store, i18n })
}

export default function () {
  return [state, { close, show, setState }]
}
