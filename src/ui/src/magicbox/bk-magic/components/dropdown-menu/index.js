import bkDropdownMenu from './src/dropdown-menu'

bkDropdownMenu.install = Vue => {
    Vue.component(bkDropdownMenu.name, bkDropdownMenu)
}

export default bkDropdownMenu
