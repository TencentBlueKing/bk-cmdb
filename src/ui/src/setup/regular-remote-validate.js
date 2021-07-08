import http from '@/api'
import { language } from '@/i18n'
const message = {
    en: 'Please input valid regular expression',
    'zh_CN': '请输入合法的正则表达式'
}
export default {
    validate: async (regular) => {
        try {
            if (!regular.trim().length) {
                return { valid: true }
            }
            const data = await http.post(`${window.API_HOST}regular/verify_regular_express`, { regular }, { globalError: false })
            const valid = data.is_valid
            return {
                valid,
                message: valid ? '' : message[language]
            }
        } catch (error) {
            return {
                valid: false,
                message: error.bk_error_msg || error.message
            }
        }
    }
}
