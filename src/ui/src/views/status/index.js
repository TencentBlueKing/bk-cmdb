import status403 from './403'
import status404 from './404'
import statusError from './error'
import statusRequireBusiness from './require-business'

export default {
    '403': status403,
    '404': status404,
    'error': statusError,
    'requireBusiness': statusRequireBusiness
}
