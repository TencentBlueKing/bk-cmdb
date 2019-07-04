import {
    isSameRequest,
    isRedirectResponse,
    getRedirectId
} from './util.js'
import Cookies from 'js-cookie'

const origin = {
    url: 'biz/with_reduced',
    method: 'get'
}
const redirect = {
    url: 'biz/search/0',
    method: 'post',
    data: {
        'fields': ['bk_biz_id', 'bk_biz_name'],
        'condition': {
            'bk_data_status': {
                '$ne': 'disabled'
            }
        }
    },
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
            let authorizedBusiness = response.data.data.info
            if (window.User.admin !== '1') {
                const cookieBizId = Cookies.get('bk_privi_biz_id')
                const authorizedBizIds = cookieBizId ? cookieBizId.split('-') : []
                authorizedBusiness = authorizedBusiness.filter(business => {
                    return authorizedBizIds.some(id => id === business.bk_biz_id.toString())
                })
            }
            response.data.data = {
                count: authorizedBusiness.length,
                info: authorizedBusiness
            }
        }
        return response
    }
}
