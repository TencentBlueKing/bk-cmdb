import Vue from 'vue'
import Axios from 'axios'
import md5 from 'md5'
import CachedPromise from './_cached-promise'
import RequestQueue from './_request-queue'
import { $error, $warn } from '@/magicbox'

const Site = window.Site
window.API_HOST = Site.buildVersion.indexOf('dev') !== -1 ? Site.url : (window.location.origin + '/')
window.API_PREFIX = window.API_HOST + 'api/' + Site.version
// axios实例
const axiosInstance = Axios.create({
    baseURL: window.API_PREFIX,
    xsrfCookieName: 'data_csrftoken',
    xsrfHeaderName: 'X-CSRFToken',
    withCredentials: true
})

// axios实例拦截器
axiosInstance.interceptors.response.use(
    response => {
        return response.data
    },
    error => {
        return Promise.reject(error)
    }
)

const $http = {
    queue: new RequestQueue(),
    cache: new CachedPromise(),
    cancelRequest: requestId => {
        return $http.queue.cancel(requestId)
    },
    cancelCache: requestId => {
        return $http.cache.delete(requestId)
    },
    cancel: requestId => {
        return Promise.all([$http.cancelRequest(requestId), $http.cancelCache(requestId)])
    },
    setHeader: (key, value) => {
        axiosInstance.defaults.headers[key] = value
    },
    deleteHeader: key => {
        delete axiosInstance.defaults.headers[key]
    }
}

const methodsWithoutData = ['delete', 'get', 'head', 'options']
const methodsWithData = ['post', 'put', 'patch']
const allMethods = [...methodsWithoutData, ...methodsWithData]

// 在自定义对象$http上添加各请求方法
allMethods.forEach(method => {
    Object.defineProperty($http, method, {
        get () {
            return getRequest(method)
        }
    })
})

/**
 * 获取http不同请求方式对应的函数
 * @param {method} http method 与 axios实例中的method保持一致
 * @return {Function} 实际调用的请求函数
 */
function getRequest (method) {
    if (methodsWithData.includes(method)) {
        return (url, data, config) => {
            return getPromise(method, url, data, config)
        }
    }
    return (url, config) => {
        return getPromise(method, url, null, config)
    }
}

/**
 * 实际发起http请求的函数，根据配置调用缓存的promise或者发起新的请求
 * @param {method} http method 与 axios实例中的method保持一致
 * @param {url} 请求地址
 * @param {data} 需要传递的数据, 仅 post/put/patch 三种请求方式可用
 * @param {userConfig} 用户配置，包含axios的配置与本系统自定义配置
 * @return {Promise} 本次http请求的Promise
 */
async function getPromise (method, url, data, userConfig = {}) {
    const config = initConfig(method, url, userConfig)
    let promise
    if (config.cancelPrevious) {
        await $http.cancel(config.requestId)
    }
    if (config.clearCache) {
        $http.cache.delete(config.requestId)
    } else {
        promise = $http.cache.get(config.requestId)
    }
    if (config.fromCache && promise) {
        return promise
    }
    promise = new Promise((resolve, reject) => {
        const axiosRequest = methodsWithData.includes(method) ? axiosInstance[method](url, data, config) : axiosInstance[method](url, config)
        axiosRequest.then(response => {
            Object.assign(config, response.config)
            handleResponse({config, response, resolve, reject})
        }).catch(error => {
            Object.assign(config, error.config)
            reject(error)
        })
    }).catch(error => {
        return handleReject(error, config)
    }).finally(() => {
        $http.queue.delete(config.requestId)
    })
    // 添加请求队列
    $http.queue.set(config)
    // 添加请求缓存
    $http.cache.set(config.requestId, promise, config)
    return promise
}

/**
 * 处理http请求成功结果
 * @param {config} 请求配置
 * @param {response} cgi原始返回数据
 * @param {resolve} promise完成函数
 * @param {reject} promise拒绝函数
 * @return
 */
function handleResponse ({config, response, resolve, reject}) {
    if (!response.result && config.globalError) {
        reject({message: response['bk_error_msg']})
    } else {
        resolve(config.originalResponse ? response : response.data, config)
    }
}

/**
 * 处理http请求失败结果
 * @param {error} Error 对象
 * @param {config} 请求配置
 * @return Promise.reject
 */
function handleReject (error, config) {
    if (Axios.isCancel(error)) {
        return Promise.reject(error)
    }
    if (config.globalError && error.response) {
        const {status, data} = error.response
        const nextError = {message: error.message}
        if (status === 401) {
            window.location.href = Site.login
        } else if (data && data['bk_error_msg']) {
            nextError.message = data['bk_error_msg']
        } else if (status === 500) {
            nextError.message = '系统出现异常, 请记录下错误场景并与开发人员联系, 谢谢!'
        }
        $error(nextError.message)
        return Promise.reject(nextError)
    }
    $error(error.message)
    return Promise.reject(error)
}

/**
 * 初始化本系统http请求的各项配置
 * @param {method} http method 与 axios实例中的method保持一致
 * @param {url} 请求地址, 结合method 生成md5 requestId
 * @param {userConfig} 用户配置，包含axios的配置与本系统自定义配置
 * @return {Promise} 本次http请求的Promise
 */
function initConfig (method, url, userConfig) {
    if (userConfig.hasOwnProperty('requestGroup')) {
        userConfig.requestGroup = userConfig.requestGroup instanceof Array ? userConfig.requestGroup : [userConfig.requestGroup]
    }
    const defaultConfig = {
        ...getCancelToken(),
        // http请求默认id
        requestId: md5(method + url),
        requestGroup: [],
        // 是否全局捕获异常
        globalError: true,
        // 是否直接复用缓存的请求
        fromCache: false,
        // 是否在请求发起前清楚缓存
        clearCache: false,
        // 响应结果是否返回原始数据
        originalResponse: false,
        // 当路由变更时取消请求
        cancelWhenRouteChange: true,
        // 取消上次请求
        cancelPrevious: false
    }
    return Object.assign(defaultConfig, userConfig)
}

/**
 * 生成http请求的cancelToken，用于取消尚未完成的请求
 * @return {Object}
 *      cancelToken: axios实例使用的cancelToken
 *      cancelExcutor: 取消http请求的可执行函数
 */
function getCancelToken () {
    let cancelExcutor
    const cancelToken = new Axios.CancelToken(excutor => {
        cancelExcutor = excutor
    })
    return {
        cancelToken,
        cancelExcutor
    }
}

Vue.prototype.$http = $http

export default $http
