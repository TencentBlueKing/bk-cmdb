import {
    isSameRequest,
    isRedirectResponse
} from './util.js'
import Cookies from 'js-cookie'

const CONFIG = {
    origin: {
        url: 'auth/business-list',
        method: 'get'
    },
    redirect: {
        url: `biz/search/0`,
        method: 'post',
        data: {
            'fields': ['bk_biz_id', 'bk_biz_name'],
            'condition': {
                'bk_data_status': {
                    '$ne': 'disabled'
                }
            }
        }
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
            const cookieBizId = Cookies.get('bk_privi_biz_id')
            const authorizedBizIds = cookieBizId ? cookieBizId.split('-') : []
            const authorizedBusiness = response.data.data.info.filter(business => {
                return authorizedBizIds.some(id => id === business.bk_biz_id.toString())
            })
            response.data.data = {
                count: authorizedBusiness.length,
                info: authorizedBusiness
            }
        }
        return response
    }
}
