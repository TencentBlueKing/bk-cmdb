import {
    isSameRequest
} from './util.js'

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
    }
}
