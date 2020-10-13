import Vue from 'vue'
import App from './App.vue'
import router from './router/index.js'
import store from './store'
import i18n from './i18n'
import cmdbRequestMixin from './mixins/request'
import cmdbAuthMixin from './mixins/auth'
import cmdbAppMixin from './mixins/app.js'
import cmdbFormatter from './filters/formatter.js'
import cmdbUnitFilter from './filters/unit.js'
import cmdbUI from './components/ui'
import cmdbSearchComponent from './components/search/index'
import routerActions from './router/actions'
import tools from './utils/tools'
import clipboard from 'vue-clipboard2'
import './magicbox'
import './directives'
import api from './api'
import './setup/cookie'
import './setup/permission'
import '@icon-cool/bk-icon-cmdb'
import '@icon-cool/bk-icon-cmdb-colorful/src/index'
import './assets/scss/common.scss'

Vue.use(cmdbUI)
Vue.use(cmdbSearchComponent)
Vue.use(clipboard)
Vue.mixin(cmdbRequestMixin)
Vue.mixin(cmdbAuthMixin)
Vue.mixin(cmdbAppMixin)
Vue.filter('formatter', cmdbFormatter)
Vue.filter('unit', cmdbUnitFilter)
Vue.prototype.$http = api
Vue.prototype.$tools = tools
Vue.prototype.$routerActions = routerActions
/* eslint-disable no-new */
window.CMDB_APP = new Vue({
    el: '#app',
    router,
    store,
    i18n,
    components: { App },
    template: '<App/>'
})

if (process.env.COMMIT_ID) {
    window.CMDB_COMMIT_ID = process.env.COMMIT_ID
}
