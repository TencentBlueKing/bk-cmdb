import bkTimeline from './timeline.vue'

bkTimeline.install = Vue => {
    Vue.component(bkTimeline.name, bkTimeline)
}

export default bkTimeline
