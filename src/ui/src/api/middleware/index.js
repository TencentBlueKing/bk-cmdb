import auth from './auth.js'
import business from './business.js'
const middlewares = [
    business,
    auth
]
const exportMiddlewares = window.Site.authscheme === 'internal' ? middlewares : []
export default exportMiddlewares
