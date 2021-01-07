/**
 * @file bk-magic main entry
 * @author ielgnaw <wuji0223@gmail.com>
 */

// 引入公共样式，这里引入公共样式，只是为全局引入 bk-magic 时起作用
// 因为只有全局引入 bk-magic 时，才会引用到这个 index.js 文件
import './bk-magic-ui/src/common.scss'

import Vue from 'vue'
import bkButton from './components/button'
import bkDropdownMenu from './components/dropdown-menu'
import bkSideslider from './components/sideslider'
import bkSwitcher from './components/switcher'
import bkTagInput from './components/tag-input'
import bkDialog from './components/dialog'
import bkIconButton from './components/icon-button'
import bkLoading from './components/loading'
import bkSteps from './components/steps'
import bkBadge from './components/badge'
import bkMessage from './components/message'
import bkTooltips from './components/tooltips'
import bkDropdown from './components/dropdown'
import bkInfoBox from './components/info-box'
import bkTooltip from './components/tooltip'
import bkDatePicker from './components/date-picker'
import bkDateRange from './components/date-range'
import bkSelector from './components/selector'
import Tab from './components/tab'
import bkPaging from './components/paging'
import bkTransfer from './components/transfer'
import bkTree from './components/tree'
import bkCollapse from './components/collapse'
import bkCollapseItem from './components/collapse-item'
import bkRound from './components/round-progress'
import bkUpload from './components/upload'
import bkTimeline from './components/timeline'
import bkProcess from './components/process'
import bkCombox from './components/combox'

import bkException from './components/exception'
import locale from './locale'
import localeMixin from './mixins/locale'
import enUS from './locale/lang/en-US'
import zhCN from './locale/lang/zh-CN'

const langPkg = {
    enUS,
    zhCN
}

const {
    bkTab,
    bkTabpanel
} = Tab

const components = {
    bkBadge,
    bkButton,
    bkCollapse,
    bkCollapseItem,
    bkCombox,
    bkDatePicker,
    bkDateRange,
    bkDialog,
    bkDropdown,
    bkDropdownMenu,
    bkIconButton,
    bkPaging,
    bkProcess,
    bkRound,
    bkSelector,
    bkSideslider,
    bkSteps,
    bkSwitcher,
    bkTab,
    bkTabpanel,
    bkTagInput,
    bkTimeline,
    bkTooltip,
    bkTransfer,
    bkTree,
    bkUpload,
    bkException
}

export function install (Vue, opts = {}) {
    locale.use(opts.locale)
    locale.i18n(opts.i18n)

    Object.keys(components).forEach(key => {
        Vue.component(components[key].name, components[key])
    })

    // locale.use({
    //     datePicker: {
    //         // 选择日期
    //         selectDate: '选择日期sdsd'
    //     }
    // })
}

Vue.use(bkLoading.directive)
Vue.prototype.$bkLoading = bkLoading.Loading

Vue.prototype.$bkMessage = bkMessage

Vue.prototype.$bkInfo = bkInfoBox

// 注册指令使用
Vue.use(bkTooltips.directive)
// 工具函数调用
Vue.prototype.$tooltips = bkTooltips.tooltips

// Vue.use(bkTooltips)

if (typeof window !== 'undefined' && window.Vue) {
    install(window.Vue)
}

const bkMagic = {
    bkException,
    ...components,
    bkLoading,
    bkMessage,
    bkTooltips,
    bkInfoBox,
    locale,
    langPkg,
    localeMixin,
    install
}

// 组件单独抛出，方便 import {bkBadge} from '@tencent/bk-magic' 这样引用
export {
    bkBadge,
    bkButton,
    bkCollapse,
    bkCollapseItem,
    bkCombox,
    bkDatePicker,
    bkDateRange,
    bkDialog,
    bkDropdown,
    bkDropdownMenu,
    bkIconButton,
    bkPaging,
    bkProcess,
    bkRound,
    bkSelector,
    bkSideslider,
    bkSteps,
    bkSwitcher,
    bkTab,
    bkTabpanel,
    bkTagInput,
    bkTimeline,
    bkTooltip,
    bkTransfer,
    bkTree,
    bkUpload,
    bkLoading,
    bkMessage,
    bkTooltips,
    bkInfoBox,
    locale,
    bkException,
    langPkg,
    localeMixin
}

// export default Object.assign(bkMagic, {install})
export default bkMagic
