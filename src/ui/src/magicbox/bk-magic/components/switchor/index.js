import Switchor from './src/switchor'

Switchor.install = Vue => {
    Vue.component(Switchor.name, Switchor)
}

export default Switchor
