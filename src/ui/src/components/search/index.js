import bool from './bool'
import date from './date'
import enumComponent from './enum'
import float from './float'
import foreignkey from './foreignkey'
import int from './int'
import list from './list'
import longchar from './longchar'
import objuser from './objuser'
import organization from './organization'
import singlechar from './singlechar'
import table from './table'
import time from './time'
import timezone from './timezone'
import serviceTemplate from './service-template'
import module from './module'
import set from './set'

export default {
  install(Vue) {
    const components = [
      bool,
      date,
      enumComponent,
      float,
      foreignkey,
      int,
      list,
      longchar,
      objuser,
      organization,
      singlechar,
      table,
      time,
      timezone,
      serviceTemplate,
      module,
      set
    ]
    components.forEach((component) => {
      Vue.component(component.name, component)
    })
  }
}
