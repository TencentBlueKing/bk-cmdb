import router from './index'
import { Base64 } from 'js-base64'
export default function ({ name, params = {}, query = {}, history = false }) {
    const queryBackup = { ...query }
    if (history) {
        const currentRoute = router.app.$route
        const data = {
            name: currentRoute.name,
            params: currentRoute.params,
            query: currentRoute.query
        }
        const base64 = Base64.encode(JSON.stringify(data))
        queryBackup['_f'] = base64
    }
    router.replace({
        name,
        params,
        query: queryBackup
    })
}
