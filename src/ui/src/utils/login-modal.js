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

let isShow = false
let loginWindow = null
let checkWindowTimer = -1

const messageListener = ({ data = {} }) => {
  if (data === null || typeof data !== 'object' || data.target !== 'bk-login' || !this.loginWindow) return

  hideLoginModal()
}

window.addEventListener('message', messageListener, false)

window.addEventListener('beforeunload', () => {
  hideLoginModal()
  window.removeEventListener('message', messageListener, false)
})

const checkWinClose = () => {
  checkWindowTimer && clearTimeout(checkWindowTimer)
  checkWindowTimer = setTimeout(() => {
    if (!loginWindow || loginWindow.closed) {
      hideLoginModal()
      clearTimeout(checkWindowTimer)
      return
    }
    checkWinClose()
  }, 300)
}

export const hideLoginModal = () => {
  isShow = false
  if (loginWindow) {
    loginWindow.close()
  }
  loginWindow = null
}

export const showLoginModal = (data = {}) => {
  if (isShow) {
    window.blur()
    loginWindow.focus()
    return
  }

  isShow = true

  const { width = 700, height = 510 } = data
  const successUrl = `${window.location.origin}/static/login_success.html`

  const siteLoginUrl = window.Site.login || ''
  if (!siteLoginUrl) {
    console.error('Login URL not configured!')
    return
  }

  const [loginBaseUrl] = siteLoginUrl.split('?')
  const loginUrl = `${loginBaseUrl}plain?c_url=${successUrl}`

  const { availHeight, availWidth } = window.screen
  loginWindow = window.open(
    loginUrl,
    '_blank',
    `
      width=${width},
      height=${height},
      left=${(availWidth - width) / 2},
      top=${(availHeight - height) / 2},
      channelmode=0,
      directories=0,
      fullscreen=0,
      location=0,
      menubar=0,
      resizable=1,
      scrollbars=0,
      status=0,
      titlebar=0,
      toolbar=0,
      close=0
    `
  )
  checkWinClose()
}
