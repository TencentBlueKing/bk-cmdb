import { camelize } from '@/utils/util.js'
import testIds from '@/dictionary/test-id.js'

const isString = val => Object.prototype.toString.call(val) === '[object String]'

const testIdKeys = Object.keys(testIds)
const tagName = name => camelize(name.toLowerCase())
const routeName = name => camelize(name.replace(/^menu_/, ''), '_')

const tagNames = ['nav', 'header', 'button', 'form', 'ul', 'li', 'div', 'section']
const availableTagName = el => el.tagName && tagNames.includes(el.tagName.toLowerCase())

function setDataTestId(el, binding, vnode) {
  const { value, modifiers } = binding
  const { context, componentInstance, componentOptions } = vnode
  const moduleId = testIdKeys.find(key => modifiers[key]) || routeName(context.$route.name)
  const moduleSetting = testIds[moduleId]

  if (value && !isString(value)) {
    console.warn('only accept string values!')
    return
  }

  // 定义文件中的优先级最高，其次是在常用标签上写了value的，最后是自定义组件的形式
  let blockId = 'unknown'
  if (moduleSetting?.[value]) {
    blockId = moduleSetting[value]
  } else if (availableTagName(el) && value) {
    blockId = `${tagName(el.tagName)}_${camelize(value)}`
  } else if (componentInstance) {
    blockId = `comp_${tagName(componentOptions.tag)}`
  }

  el.setAttribute('data-test-id', `${moduleId}_${blockId}`)
}

const directive = {
  bind(el, binding, vnode) {
    setDataTestId(el, binding, vnode)
  },
  componentUpdated(el, binding, vnode) {
    if (binding.modifiers.dynamic) {
      setDataTestId(el, binding, vnode)
    }
  }
}

export default {
  install: Vue => Vue.directive('test-id', directive)
}
