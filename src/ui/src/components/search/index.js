import bool from './bool'
import date from './date'
import enumComponent from './enum'
import float from './float'
import int from './int'
import list from './list'
import longchar from './longchar'
import objuser from './objuser'
import organization from './organization'
import singlechar from './singlechar'
import time from './time'
import timezone from './timezone'
import serviceTemplate from './service-template'

export default {
    install (Vue, ops = {}) {
        const components = [
            bool,
            date,
            enumComponent,
            float,
            int,
            list,
            longchar,
            objuser,
            organization,
            singlechar,
            time,
            timezone,
            serviceTemplate
        ]
        components.forEach(component => {
            Vue.component(component.name, component)
        })
    }
}
