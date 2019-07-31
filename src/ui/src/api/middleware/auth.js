/* eslint-disable */
import * as AUTH from '@/dictionary/auth'
import { 
    isSameRequest,
    isRedirectResponse,
    getRedirectId
} from './util.js'

const authActionMap = {
    'findMany': 'search',
    'create': 'update',
    'archive': 'delete'
}

const businessAuth = [
    AUTH.C_BUSINESS,
    AUTH.U_BUSINESS,
    AUTH.R_BUSINESS,
    AUTH.BUSINESS_ARCHIVE
]

const businessResourceAuth = [
    AUTH.U_HOST,
    AUTH.R_HOST,
    AUTH.HOST_TO_RESOURCE,

    AUTH.R_PROCESS,
    AUTH.C_PROCESS,
    AUTH.U_PROCESS,
    AUTH.D_PROCESS,
    AUTH.PROCESS_BIND_MODULE,
    AUTH.PROCESS_UNBIND_MODULE,

    AUTH.R_CUSTOM_QUERY,
    AUTH.C_CUSTOM_QUERY,
    AUTH.U_CUSTOM_QUERY,
    AUTH.D_CUSTOM_QUERY
]

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

const origin = {
    url: 'auth/verify',
    method: 'post',
    data: []
}
const redirect = {
    url: `topo/privilege/user/detail/0/${window.User.name}`,
    method: 'get',
    redirectId: getRedirectId()
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
    const payload = origin.data
    const modelConfig = flattenModelConfig(data.model_config)
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
            } else if (businessAuth.includes(auth)) {
                meta.parent_layers = [{
                    resource_model: 'biz'
                }]
                setModelMeta(meta, modelConfig)
            } else if (businessResourceAuth.includes(auth) && resource.bk_biz_id) {
                meta.is_pass = true
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

const flattenModelConfig = (modelConfig = {}) => {
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
        if (isSameRequest(origin, config)) {
            origin.data = config.data
            Object.assign(config, redirect)
        }
        return config
    },
    response: response => {
        if (isRedirectResponse(redirect, response)) {
            const data = transformResponse(response.data.data)
            Object.assign(response.data, {
                data: data
            })
        }
        return response
    }
}
