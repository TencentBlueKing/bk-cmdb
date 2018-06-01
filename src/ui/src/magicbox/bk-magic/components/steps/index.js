import bkSteps from './src/steps'

bkSteps.install = Vue => {
    Vue.component(bkSteps.name, bkSteps)
}

export default bkSteps
