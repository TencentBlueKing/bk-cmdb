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
import VueI18n from 'vue-i18n'
import Cookies from 'js-cookie'
import messages from './lang/messages'
import { LANG_COOKIE_NAME, LANG_KEYS, LANG_SET } from './constants'
import { setLocale as setMagicboxLocale } from '@/magicbox'

Vue.use(VueI18n)

const langInCookie = Cookies.get(LANG_COOKIE_NAME)
const matchedLang = LANG_SET.find(lang => lang.id === langInCookie || lang?.alias?.includes(langInCookie))
const locale = matchedLang?.id || LANG_KEYS.ZH_CN

const i18n = new VueI18n({
  locale,
  fallbackLocale: LANG_KEYS.ZH_CN,
  messages,
  missing(locale, path) {
    // eslint-disable-next-line no-underscore-dangle
    const parsedPath = i18n._path.parsePath(path)
    return parsedPath[parsedPath.length - 1]
  }
})

export const changeLocale = (locale) => {
  Cookies.set(LANG_COOKIE_NAME, locale, {
    expires: 365,
    domain: window.location.hostname.split('.').slice(-2)
      .join('.')
  })

  document.body.setAttribute('lang', locale)

  setMagicboxLocale(locale)

  i18n.locale = locale
}

export const language = locale

export const t = (content, ...rest) => i18n.t(content, ...rest)

export default i18n
