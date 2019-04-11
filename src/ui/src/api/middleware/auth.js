/* eslint-disable */
import * as AUTH from '@/dictionary/auth'

const authActionMap = {
    'findMany': 'search',
    'create': 'update'
}

const modelAuth = [
    AUTH.C_INST,
    AUTH.U_INST,
    AUTH.D_INST,
    AUTH.R_INST
]

const getFullUrl = url => {
    return `${window.API_PREFIX}/${url}`
}

const CONFIG = {
    origin: {
        url: getFullUrl('auth/verify'),
        method: 'post'
    },
    redirect: {
        url: getFullUrl('auth/verify'),
        method: 'post'
        // url: getFullUrl(`topo/privilege/user/detail/0/${window.User.name}`),
        // method: 'get'
    }
}

export default {
    request: config => {
        // if (
        //     getFullUrl(config.url) === CONFIG.origin.url
        //     && config.method === CONFIG.origin.method
        // ) {
        //     Object.assign(config, CONFIG.redirect)
        // }
        return config
    },
    response: response => {
        const { url, baseURL, method } = response.config
        if (
            url === CONFIG.redirect.url
            && method === CONFIG.redirect.method
        ) {
            console.log(response)
        }
        return response
    }
}
