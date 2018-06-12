/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'
import VueI18n from 'vue-i18n'

import VeeValidate from 'vee-validate'
import dictionary from './common/js/Validator'
import bkMagic from './magicbox/bk-magic'

import i18nConfig from './common/js/i18n'
import VTooltip from 'v-tooltip'
import vClickOutside from 'v-click-outside'
import Cookies from 'js-cookie'
import moment from 'moment'
import vDrag from './directive/drag'
import axiosQueue from './mixins/axios-queue'
import '@/api/axios'

const languageTranslate = {
    'zh_cn': 'zh_CN',
    'zh-cn': 'zh_CN',
    'zh': 'zh_CN'
}
const bkMagicLang = {
    'zh_CN': 'zh',
    'en': 'en'
}
let language = Cookies.get('blueking_language') || 'zh_CN'
language = languageTranslate.hasOwnProperty(language) ? languageTranslate[language] : language
document.body.setAttribute('lang', language)
Vue.use(vDrag)
Vue.use(VTooltip)
Vue.use(vClickOutside)
Vue.use(VueI18n)
Vue.use(bkMagic)
Vue.use(VeeValidate, {locale: language})

Vue.directive('focus', {
    update (el, binding) {
        if (binding.value !== binding.oldVal) {
            el.focus()
        }
    }
})

Vue.mixin(axiosQueue)

Vue.config.productionTip = false

Vue.prototype.$alertMsg = (msg, theme = 'error') => {
    Vue.prototype.$bkMessage({
        message: msg,
        theme: theme,
        delay: 3000
    })
}

Vue.prototype.$formatTime = (data, formatStr = 'YYYY-MM-DD HH:mm:ss') => {
    if (!data) {
        return ''
    }
    let time = moment(data).format(formatStr)
    if (time === 'Invalid date') {
        return data
    } else {
        return time
    }
}

Vue.prototype.$deepClone = (data) => {
    return JSON.parse(JSON.stringify(data))
}

/* eslint-disable no-new */
window.CMDB_APP = new Vue({
    el: '#app',
    router,
    template: '<App/>',
    i18n: new VueI18n({
        locale: language,
        messages: i18nConfig,
        fallbackLocale: 'zh_CN',
        missing: function (locale, path) {
            let parsedPath = window.CMDB_APP.$i18n._path.parsePath(path)
            return parsedPath[parsedPath.length - 1]
        }
    }),
    store,
    components: { App },
    created () {
        this.setLang(bkMagicLang[language])
    }
})
