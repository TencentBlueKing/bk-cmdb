/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import Vue from 'vue'
import { subEnv } from '@blueking/sub-saas'

import App from './App.vue'
import IframeEntry from './IframeEntry.vue'
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
import { gotoLoginPage } from '@/utils/login-helper'
import clipboard from 'vue-clipboard2'
import './magicbox'
import './directives'
import api from './api'
import './setup/cookie'
import './setup/permission'
import './setup/build-in-vars'
import '@/assets/icon/bk-icon-cmdb/style.css'
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

api.get(`${window.API_HOST}is_login`).then(() => {
  window.CMDB_APP = new Vue({
    el: '#app',
    router,
    store,
    i18n,
    render() {
      return !subEnv ? <App /> : <IframeEntry />
    }
  })
})
  .catch(() => {
    gotoLoginPage()
  })

if (process.env.COMMIT_ID) {
  window.CMDB_COMMIT_ID = process.env.COMMIT_ID
}
