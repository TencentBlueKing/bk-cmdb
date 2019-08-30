import {
    isSameRequest,
    isRedirectResponse,
    getRedirectId
} from './util.js'

const origin = {
    url: 'auth/admin_entrance',
    method: 'get'
}
const redirect = {
    url: `topo/privilege/user/detail/0/${window.User.name}`,
    method: 'get',
    redirectId: getRedirectId()
}

export default {
    request: config => {
        if (isSameRequest(origin, config)) {
            Object.assign(config, redirect)
        }
        return config
    },
    response: response => {
        if (isRedirectResponse(redirect, response)) {
            response.data.data = {
                is_pass: window.User.admin === '1'
            }
        }
        return response
    }
}
