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

import { showLoginModal as showModal } from '@blueking/login-modal'

export const showLoginModal = () => {
  const successUrl = `${window.location.origin}/static/login_success.html`

  const siteLoginUrl = window.Site.login
  if (!siteLoginUrl) {
    console.error('Login URL not configured!')
    return
  }

  const loginURL = new URL(siteLoginUrl)
  loginURL.searchParams.set('c_url', successUrl)
  const pathname = loginURL.pathname.endsWith('/') ? loginURL.pathname : `${loginURL.pathname}/`
  const loginUrl = `${loginURL.origin}${pathname}plain/${loginURL.search}`

  showModal({ loginUrl })
}

export const gotoLoginPage = () => {
  if (!window.Site.login) {
    console.error('The login URL is not configured!')
    return
  }
  try {
    const loginURL = new URL(window.Site.login)
    loginURL.searchParams.set('c_url', location.href)
    location.href = loginURL.href
  } catch (_) {
    console.error('The login URL invalid!')
  }
}
