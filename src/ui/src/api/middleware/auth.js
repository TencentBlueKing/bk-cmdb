/* eslint-disable */
import * as AUTH from '@/dictionary/auth'
import { 
    isSameRequest,
    isRedirectResponse
} from './util.js'

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

const CONFIG = {
    origin: {
        url: 'auth/verify',
        method: 'post'
    },
    redirect: {
        url: `topo/privilege/user/detail/0/${window.User.name}`,
        method: 'get'
    }
}

export default {
    request: config => {
        if (isSameRequest(CONFIG.origin, config)) {
            Object.assign(config, CONFIG.redirect)
        }
        return config
    },
    response: response => {
        if (isRedirectResponse(CONFIG.redirect, response)) {
            Object.assign(response.data, {
                data: []
            })
        }
        return response
    }
}
