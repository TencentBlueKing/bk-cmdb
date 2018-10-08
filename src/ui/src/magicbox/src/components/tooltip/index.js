/**
 * @file tooltip entry
 * @author ielgnaw <wuji0223@gmail.com>
 */

import bkTooltip from './tooltip'

bkTooltip.install = Vue => {
    Vue.component(bkTooltip.name, bkTooltip)
}

export default bkTooltip
