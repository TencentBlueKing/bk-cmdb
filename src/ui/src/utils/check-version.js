/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { bkInfoBox } from 'bk-magic-vue'
import i18n from '@/i18n'
import useTimeoutPoll from '@/hooks/use-timeout-poll'

const STATIC_PATH = `${window.Site.publicPath ?? '/'}static/`
const VERSION_FILE = 'build-hash.txt'

let localVersion = ''
let isShown = false
let dialog = null
let checkVersionPoll = null

const checkVersionChannel = new BroadcastChannel('check-version')
checkVersionChannel.addEventListener('message', ({ data }) => {
  if (data?.type !== 'toggle') {
    return
  }

  // 共享弹出状态，视之后有无需求多个标签页只允许一次弹出提示使用
  // isShown = data.payload.isShown;
})

const hideDialog = (silent = false) => {
  if (dialog) {
    dialog.close()
  }

  isShown = false
  dialog = null

  checkVersionPoll.resume()

  if (!silent) {
    checkVersionChannel.postMessage({ type: 'toggle', payload: { isShown: false } })
  }
}

const fetchBuildHash = async () => {
  const response = await fetch(`${STATIC_PATH}${VERSION_FILE}?_=${Date.now()}`)

  if (!response.ok) {
    throw new Error('Failed to get build hash')
  }

  const newVersion = await response.text()

  // localVersion 还不存在，认为是第一次加载
  if (!localVersion) {
    localVersion = newVersion
  }

  if (newVersion !== localVersion) {
    showVersionDialog(newVersion)
  }
}

const checkVersion = async () => {
  try {
    await fetchBuildHash()
  } catch (error) {
    console.error(error)
  }
}

const handleVisibilityChange = () => {
  if (!isShown && document.visibilityState === 'visible') {
    checkVersionPoll.resume()
  } else {
    checkVersionPoll.pause()
  }
}

const showVersionDialog = (newVersion) => {
  if (isShown) {
    return
  }

  dialog = bkInfoBox({
    width: 650,
    maskClose: false,
    title: i18n.t('版本更新'),
    subTitle: i18n.t('建议「刷新页面」体验新的特性，「暂不刷新」可能会遇到未知异常，可手动刷新解决。'),
    okText: i18n.t('刷新'),
    cancelText: i18n.t('取消'),
    confirmFn: () => {
      window.location.reload()
    },
    cancelFn: () => {
      hideDialog()
    },
  })

  if (dialog) {
    isShown = true
    localVersion = newVersion
    checkVersionPoll.pause()
    checkVersionChannel.postMessage({ type: 'toggle', payload: { isShown: true } })
  }
}

const addVisibilityChangeListener = () => {
  document.addEventListener('visibilitychange', handleVisibilityChange)
}

const removeVisibilityChangeListener = () => {
  document.removeEventListener('visibilitychange', handleVisibilityChange)
}

const addBeforeunloadListener = () => {
  window.addEventListener('beforeunload', () => {
    hideDialog()
    removeVisibilityChangeListener()
    checkVersionPoll.reset()
    checkVersionChannel.close()
  })
}

export const watchVersion = () => {
  checkVersionPoll = useTimeoutPoll(checkVersion, 5 * 60 * 1000, { immediate: true, max: -1 })
  addVisibilityChangeListener()
  addBeforeunloadListener()
}
