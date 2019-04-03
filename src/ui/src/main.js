import Vue from 'vue'
import App from './App.vue'
import router from './router/index.js'
import store from './store'
import i18n from './i18n'
import cmdbRequestMixin from './mixins/request'
import cmdbAuthMixin from './mixins/auth'
import cmdbInjectMixin from './mixins/inject'
import cmdbUI from './components/ui'
import tools from './utils/tools'
import clipboard from 'vue-clipboard2'
import './directives'
import './magicbox'
import api from './api'
import './setup/validate'
import './setup/cookie'
import './assets/scss/common.scss'
import './assets/scss/admin-view.scss'
import './assets/icon/cc-icon/style.css'
import './assets/icon/bk-icon-2.0/iconfont.css'
Vue.use(cmdbUI)
Vue.use(clipboard)
Vue.mixin(cmdbRequestMixin)
Vue.mixin(cmdbAuthMixin)
Vue.mixin(cmdbInjectMixin)
Vue.prototype.$http = api
Vue.prototype.$tools = tools
/* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    store,
    i18n,
    components: { App },
    template: '<App/>'
})
