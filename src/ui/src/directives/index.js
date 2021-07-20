import Vue from 'vue'
import vClickOutside from 'v-click-outside'
import cursor from './cursor.js'
import transferDom from './transfer-dom.js'

Vue.use(vClickOutside)
Vue.use(cursor)
Vue.directive('transfer-dom', transferDom)

export default {
  'v-click-outside': vClickOutside,
  'v-transfer-dom': transferDom
}
