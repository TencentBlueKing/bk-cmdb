/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import Vue from 'vue'
import Axios from 'axios'
import bkMessage from '@/magicbox/bk-magic/components/message'

const alertMsg = (message, theme = 'error', delay = 3000) => {
    bkMessage({
        message,
        theme,
        delay
    })
}

const catchErrorMsg = (response) => {
    let msg = '系统出现异常, 请记录下错误场景并与开发人员联系, 谢谢!'
    if (response.data && response.data['bk_error_msg']) {
        msg = response.data['bk_error_msg']
    } else if (response.statusText) {
        msg = response.statusText
    }
    alertMsg(msg)
}

let axios = Axios.create({
    baseURL: `${window.siteUrl}api/${window.version}/`,
    xsrfCookieName: 'data_csrftoken',
    xsrfHeaderName: 'X-CSRFToken',
    withCredentials: true,
    headers: {
        'bkcclanguage': 'cn'   // 取值 cn/en
    }
})

const updateLoadingStatus = (config) => {
    if (config.hasOwnProperty('id')) {
        let queue = [...window.CMDB_APP.$store.state.common.axiosQueue]
        queue.splice(queue.indexOf(config.id), 1)
        window.CMDB_APP.$store.commit('updateAxiosQueue', queue)
    }
}

axios.interceptors.request.use(config => {
    const axiosQueue = window.CMDB_APP.$store.state.common.axiosQueue
    if (config.hasOwnProperty('id') && !axiosQueue.some(id => config.id === id)) {
        window.CMDB_APP.$store.commit('updateAxiosQueue', [...axiosQueue, config.id])
    }
    return config
})
axios.interceptors.response.use(
    response => {
        const config = response.config
        updateLoadingStatus(config)
        return response.data
    },
    error => {
        const config = error.config
        updateLoadingStatus(config)
        const globalError = config.hasOwnProperty('globalError') ? !!config.globalError : true
        if (globalError && error.response) {
            switch (error.response.status) {
                case 401:
                    window.location.href = window.loginUrl
                    break
                case 403:
                    alertMsg(error.response.statusText)
                    break
                case 500:
                    // 异常
                    alert('系统出现异常, 请记录下错误场景并与开发人员联系, 谢谢!')
                    break
                default:
                    catchErrorMsg(error.response)
            }
        }
        return Promise.reject(error)   // 返回接口返回的错误信息
    }
)

Vue.prototype.$axios = axios
Vue.prototype.$Axios = Axios

export const $axios = axios
export const $Axios = Axios
export const $alertMsg = alertMsg
