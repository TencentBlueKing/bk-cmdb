export default class Meta {
    constructor (data = {}) {
        const menu = Object.assign({
            id: null,
            i18n: null,
            path: null,
            order: 1,
            parent: null,
            adminView: true,
            onlyAdminView: false,
            requireBusiness: false
        }, data.menu || {})

        const auth = Object.assign({
            view: '',
            operation: []
        })

        return {
            ...data,
            menu,
            auth
        }
    }
}
