import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import i18n from './i18n'
import cmdbRequestMixin from './mixins/request'
import cmdbClassifyMixin from './mixins/classify'
import cmdbUI from './components/ui'
import tools from './utils/tools'
import clipboard from 'vue-clipboard2'
import './directives'
import './magicbox'
import './api'
import './setup/validate'
import './assets/scss/common.scss'
import './assets/icon/cc-icon/style.css'
import './assets/icon/bk-icon-2.0/iconfont.css'
Vue.use(cmdbUI)
Vue.use(clipboard)
Vue.mixin(cmdbRequestMixin)
Vue.mixin(cmdbClassifyMixin)
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
