import Vue from 'vue'
import vClickOutside from 'v-click-outside'
import vTooltip from 'v-tooltip'
import cursor from './cursor.js'

Vue.use(vClickOutside)
Vue.use(vTooltip)
Vue.use(cursor)

export default {
    'v-click-outside': vClickOutside,
    'vTooltip': vTooltip
}
