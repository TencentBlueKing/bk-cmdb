import Vue from 'vue'
import Button from './components/button/index'
import IconButton from './components/button-icon/index'
import Dialog from './components/dialog/index'
import InfoBox from './components/info-box/index'
import sideSlider from './components/side-slider/index'
import Steps from './components/steps/index'
import Badge from './components/badge/index'
import Paging from './components/paging/index'
import Table from './components/table/index'
import Switchor from './components/switchor/index'
import DropDown from './components/dropdown/index'
import Selector from './components/selector/index'
import Loading from './components/loading/index'
import Message from './components/message/index'
import Tooltips from './components/tooltips/index'
import Tab from './components/tab/index'
import dropdownMenu from './components/dropdown-menu/index'
import Toolbox from './components/toolbox/index'
import DateSelect from './components/date/index'
import DateRangeSelect from './components/date-range/index'
import Select from './components/select/index'
import NumberInput from './components/number/index'
import locale from './locale'

const install = (Vue, opts = {}) => {
    const components = [
        Dialog,
        sideSlider,
        Steps,
        Badge,
        Paging,
        DropDown,
        Selector,
        NumberInput,
        Tab.bkTab,
        Tab.bkTabPanel,
        dropdownMenu,
        Toolbox,
        DateSelect,
        DateRangeSelect,
        Select.bkSelect,
        Select.bkSelectOption,
        Select.bkOptionGroup
    ]

    const formComponents = [
        Button,
        IconButton,
        Table,
        Switchor
    ]

    components.map(component => {
        Vue.component(component.name, component)
    })

    formComponents.map(component => {
        Vue.component(component.name, component)
    })

    Vue.use(Loading.directive)
    Vue.use(Tooltips)
    Vue.prototype.$bkInfo = InfoBox
    Vue.prototype.$bkLoading = Loading.Loading
    Vue.prototype.$bkMessage = Message

    Vue.prototype.t = locale.t

    Vue.prototype.setLang = (lang) => {
        locale.use(lang)
    }
}

export default {
    version: '1.0.0',
    install: install,
    Button,
    IconButton,
    Dialog,
    sideSlider,
    Steps,
    Badge,
    Paging,
    DropDown,
    Table,
    Switchor,
    Message,
    Tab,
    dropdownMenu,
    Toolbox,
    Select,
    NumberInput
}
