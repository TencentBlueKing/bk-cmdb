import auth from './auth.js'
import business from './business.js'
import adminEntrance from './admin-entrance.js'
const middlewares = [
    business,
    auth,
    adminEntrance
]
const exportMiddlewares = window.Site.authscheme === 'internal' ? middlewares : []
export default exportMiddlewares
