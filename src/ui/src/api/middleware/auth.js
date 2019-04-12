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

const resourceAuth = [
    AUTH.C_RESOURCE_HOST,
    AUTH.U_RESOURCE_HOST,
    AUTH.D_RESOURCE_HOST
]

const eventAuth = [
    AUTH.C_EVENT,
    AUTH.U_EVENT,
    AUTH.D_EVENT,
    AUTH.R_EVENT
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
    is_pass: false,
    parent_layers: null,
    reason: '',
    resource_id: 0,
    resource_type: null,
    action: null
}


const transformResponse = data => {
    const payload = CONFIG.origin.data
    const modelConfig = flatternModelConfig(data.model_config)
    const backConfig = data.sys_config.back_config || []
    const globalBusi = data.sys_config.global_busi || []
    return payload.resources.map(resource => {
        const meta = {
            ...defaultMeta,
            ...resource
        }
        const auth = `${meta.resource_type}.${meta.action}`
        if (isAdmin) {
            meta.is_pass = true
        } else {
            if (modelAuth.includes(auth)) {
                setModelMeta(meta, modelConfig)
            } else if (resourceAuth.includes(auth)) {
                setSystemMeta('resource', meta, globalBusi)
            } else if (eventAuth.includes(auth)) {
                setSystemMeta('event', meta, backConfig)
            } else if (auth === AUTH.R_AUDIT) {
                setSystemMeta('audit', meta, backConfig)
            }
        }
        return meta
    })
}

// 模型实例操作转换
const setModelMeta = (meta, config) => {
    const parentResource = meta.parent_layers[0] || {}
    const model = parentResource.resource_model
    const action = meta.action
    meta.is_pass = (config[model] || []).includes(authActionMap[action] || action)
}

// 系统权限
const setSystemMeta = (type, meta, config) => {
    meta.is_pass = config.includes(type)
}

const flatternModelConfig = (modelConfig = {}) => {
    const config = {}
    Object.values(modelConfig).forEach(group => {
        Object.keys(group).forEach(model => {
            config[model] = group[model]
        })
    })
    return config
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
            const data = transformResponse(response.data.data)
            Object.assign(response.data, {
                data: data
            })
        }
        return response
    }
}
