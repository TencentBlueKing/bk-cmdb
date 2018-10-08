/**
 * @file button-icon entry
 * @author ielgnaw <wuji0223@gmail.com>
 */

import bkIconButton from './button'

bkIconButton.install = Vue => {
    Vue.component(bkIconButton.name, bkIconButton)
}

export default bkIconButton
