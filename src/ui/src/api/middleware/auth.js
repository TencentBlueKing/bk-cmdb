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
        method: 'post',
        data: []
    },
    redirect: {
        url: `topo/privilege/user/detail/0/${window.User.name}`,
        method: 'get'
    }
}

const isAdmin = window.User.admin === '1'

const defaultMeta = {
    bk_biz_id: 0,
    is_pass: true,
    parent_layers: null,
    reason: '',
    resource_id: 0,
    resource_type: null,
    action: null
}


const transformResponse = data => {
    const payload = CONFIG.origin.data
    return payload.resources.map(resource => {
        const meta = {
            ...defaultMeta,
            ...resource
        }

        const auth = `${meta.resource_type}.${meta.action}`
        switch (auth) {
            case AUTH.SYSTEM_MANAGEMENT:
                setAdminEntranceMeta(meta)
                break
            default:
                break
        }
        return meta
    })
}

const setAdminEntranceMeta = meta => {
    if (isAdmin) {
        meta.is_pass = true
    }
    return meta
}

export default {
    request: config => {
        if (isSameRequest(CONFIG.origin, config)) {
            CONFIG.origin.data = config.data
            Object.assign(config, CONFIG.redirect)
        }
        return config
    },
    response: response => {
        if (isRedirectResponse(CONFIG.redirect, response)) {
            const data = transformResponse(response.data)
            Object.assign(response.data, {
                data: data
            })
        }
        return response
    }
}
